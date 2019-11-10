package accesslog

import (
	"bytes"

	"net/http"
	"text/template"
	"time"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/response"
)

// LoggerEntry is the structure passed to the template.
type LoggerEntry struct {
	StartTime string
	Status    int
	Duration  time.Duration
	Hostname  string
	Method    string
	Path      string
	Request   *http.Request
}

// LoggerDefaultFormat is the format logged used by the default Logger instance.
var LoggerDefaultFormat = "{{.Status}} | {{.Duration}} | {{.Hostname}} | {{.Method}} {{.Path}}"

// LoggerDefaultDateFormat is the format used for date by the default Logger instance.
var LoggerDefaultDateFormat = time.RFC3339

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	logger.Logger
	dateFormat string
	template   *template.Template
}

// NewLogger returns a new Logger instance
func NewLogger(log logger.Logger) *Logger {
	lm := &Logger{dateFormat: LoggerDefaultDateFormat}
	lm.Logger = log
	lm.SetFormat(LoggerDefaultFormat)
	return lm
}

func (l *Logger) SetFormat(format string) {
	l.template = template.Must(template.New("negroni_parser").Parse(format))
}

func (l *Logger) SetDateFormat(format string) {
	l.dateFormat = format
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, r)

	res := rw.(response.Response)
	log := LoggerEntry{
		StartTime: start.Format(l.dateFormat),
		Status:    res.Status(),
		Duration:  time.Since(start),
		Hostname:  r.Host,
		Method:    r.Method,
		Path:      r.URL.Path,
		Request:   r,
	}

	buff := &bytes.Buffer{}
	l.template.Execute(buff, log)
	l.Debug(buff.String())
}
