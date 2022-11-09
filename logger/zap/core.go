package zap

import (
	"flag"
	"io/ioutil"
	golog "log"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"unsafe"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/infraboard/mcube/file"
	"github.com/infraboard/mcube/path"
)

var (
	_log unsafe.Pointer // Pointer to a coreLogger. Access via atomic.LoadPointer.
)

// 初始化默认Logger
func init() {
	storeLogger(&coreLogger{
		selectors:    map[string]struct{}{},
		rootLogger:   zap.NewNop(),
		globalLogger: zap.NewNop(),
		logger:       newLogger(zap.NewNop(), ""),
		atom:         zap.NewAtomicLevel(),
	})
}

type coreLogger struct {
	selectors    map[string]struct{}    // Set of enabled debug selectors.
	rootLogger   *zap.Logger            // Root logger without any options configured.
	globalLogger *zap.Logger            // Logger used by legacy global functions (e.g. logp.Info).
	logger       *Logger                // Logger that is the basis for all logp.Loggers.
	observedLogs *observer.ObservedLogs // Contains events generated while in observation mode (a testing mode).
	atom         zap.AtomicLevel        // 动态Level设置
}

// Configure configures the logp package.
func Configure(cfg Config) error {
	var (
		sink         zapcore.Core
		observedLogs *observer.ObservedLogs
		err          error
	)

	atom := zap.NewAtomicLevel()
	atom.SetLevel(cfg.Level.zapLevel())
	// Build a single output (stderr has priority if more than one are enabled).
	switch {
	case cfg.toObserver:
		sink, observedLogs = observer.New(atom)
	case cfg.toIODiscard:
		sink, err = makeDiscardOutput(cfg, atom)
	case cfg.ToStderr:
		sink, err = makeStderrOutput(cfg, atom)
	case cfg.ToSyslog:
		sink, err = makeSyslogOutput(cfg, atom)
	case cfg.ToEventLog:
		sink, err = makeEventLogOutput(cfg, atom)
	case cfg.ToFiles:
		fallthrough
	default:
		sink, err = makeFileOutput(cfg, atom)
	}
	if err != nil {
		return errors.Wrap(err, "failed to build log output")
	}

	// Enabled selectors when debug is enabled.
	selectors := make(map[string]struct{}, len(cfg.Selectors))
	if cfg.Level.Enabled(DebugLevel) && len(cfg.Selectors) > 0 {
		for _, sel := range cfg.Selectors {
			selectors[strings.TrimSpace(sel)] = struct{}{}
		}

		// Default to all enabled if no selectors are specified.
		if len(selectors) == 0 {
			selectors["*"] = struct{}{}
		}

		if _, enabled := selectors["stdlog"]; !enabled {
			// Disable standard logging by default (this is sometimes used by
			// libraries and we don't want their spam).
			golog.SetOutput(ioutil.Discard)
		}

		sink = selectiveWrapper(sink, selectors)
	}

	root := zap.New(sink, makeOptions(cfg)...)
	storeLogger(&coreLogger{
		selectors:    selectors,
		rootLogger:   root,
		globalLogger: root.WithOptions(zap.AddCallerSkip(1)),
		logger:       newLogger(root, ""),
		observedLogs: observedLogs,
		atom:         atom,
	})

	return nil
}

// DevelopmentSetup configures the logger in development mode at debug level.
// By default the output goes to stderr.
func DevelopmentSetup(options ...Option) error {
	cfg := Config{
		Level:       DebugLevel,
		ToStderr:    true,
		development: true,
		addCaller:   true,
	}
	for _, apply := range options {
		apply(&cfg)
	}
	return Configure(cfg)
}

// TestingSetup configures logging by calling DevelopmentSetup if and only if
// verbose testing is enabled (as in 'go test -v').
func TestingSetup(options ...Option) error {
	// Use the flag to avoid a dependency on the testing package.
	f := flag.Lookup("test.v")
	if f != nil && f.Value.String() == "true" {
		return DevelopmentSetup(options...)
	}
	return nil
}

// ObserverLogs provides the list of logs generated during the observation
// process.
func ObserverLogs() *observer.ObservedLogs {
	return loadLogger().observedLogs
}

// Sync flushes any buffered log entries. Applications should take care to call
// Sync before exiting.
func Sync() error {
	return loadLogger().rootLogger.Sync()
}

func makeOptions(cfg Config) []zap.Option {
	var options []zap.Option
	if cfg.addCaller {
		options = append(options, zap.AddCaller())
	}
	if cfg.development {
		options = append(options, zap.Development())
	}

	if len(cfg.Metas) > 0 {
		options = append(options, zap.Fields(cfg.Metas...))
	}

	return options
}

func makeStderrOutput(cfg Config, enab zapcore.LevelEnabler) (zapcore.Core, error) {
	stderr := zapcore.Lock(os.Stderr)
	return zapcore.NewCore(buildEncoder(cfg), stderr, enab), nil
}

func makeDiscardOutput(cfg Config, enab zapcore.LevelEnabler) (zapcore.Core, error) {
	discard := zapcore.AddSync(ioutil.Discard)
	return zapcore.NewCore(buildEncoder(cfg), discard, enab), nil
}

func makeSyslogOutput(cfg Config, enab zapcore.LevelEnabler) (zapcore.Core, error) {
	return newSyslog(buildEncoder(cfg), enab)
}

func makeEventLogOutput(cfg Config, enab zapcore.LevelEnabler) (zapcore.Core, error) {
	return newEventLog(cfg.Name, buildEncoder(cfg), enab)
}

func makeFileOutput(cfg Config, enab zapcore.LevelEnabler) (zapcore.Core, error) {
	name := cfg.Name
	if cfg.Files.Name != "" {
		name = cfg.Files.Name
	}
	filename := path.Resolve(path.Logs, filepath.Join(cfg.Files.Path, name))

	rotator, err := file.NewFileRotator(filename,
		file.MaxSizeBytes(cfg.Files.MaxSize),
		file.MaxBackups(cfg.Files.MaxBackups),
		file.Permissions(os.FileMode(cfg.Files.Permissions)),
		file.Interval(cfg.Files.Interval),
		file.RedirectStderr(cfg.Files.RedirectStderr),
		file.RotateOnStartup(cfg.Files.RotateOnStartup),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create file rotator")
	}

	return zapcore.NewCore(buildEncoder(cfg), rotator, enab), nil
}

func globalLogger() *zap.Logger {
	return loadLogger().globalLogger
}

func loadLogger() *coreLogger {
	p := atomic.LoadPointer(&_log)
	return (*coreLogger)(p)
}

func storeLogger(l *coreLogger) {
	if old := loadLogger(); old != nil {
		old.rootLogger.Sync()
	}
	atomic.StorePointer(&_log, unsafe.Pointer(l))
}

// HasSelector returns true if the given selector was explicitly set.
func HasSelector(selector string) bool {
	_, found := loadLogger().selectors[selector]
	return found
}
