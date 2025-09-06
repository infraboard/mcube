package log

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &Config{
	CallerDeep: 3,
	Level:      zerolog.DebugLevel.String(),
	TraceFiled: "trace_id",
	Console: Console{
		Enable:  true,
		NoColor: false,
	},
	File: File{
		DirPath:    "logs",
		Enable:     false,
		MaxSize:    100,
		MaxBackups: 6,
	},
	root:    &log.Logger,
	loggers: make(map[string]*zerolog.Logger),
}

type Console struct {
	Enable  bool `toml:"enable" json:"enable" yaml:"enable"  env:"ENABLE"`
	NoColor bool `toml:"no_color" json:"no_color" yaml:"no_color"  env:"NO_COLOR"`
}

func (c *Console) ConsoleWriter() io.Writer {
	output := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.NoColor = c.NoColor
		w.TimeFormat = time.RFC3339
	})

	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%-6s", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	return output
}

type File struct {
	// 是否开启文件记录
	Enable bool `toml:"enable" json:"enable" yaml:"enable"  env:"ENABLE"`
	// 日志目录路径
	DirPath string `toml:"dir_path" json:"dir_path" yaml:"dir_path"  env:"DIR_PATH"`
	// 单位M, 默认100M
	MaxSize int `toml:"max_size" json:"max_size" yaml:"max_size"  env:"MAX_SIZE"`
	// 默认保存 6个文件
	MaxBackups int `toml:"max_backups" json:"max_backups" yaml:"max_backups"  env:"MAX_BACKUPS"`
	// 保存多久
	MaxAge int `toml:"max_age" json:"max_age" yaml:"max_age"  env:"MAX_AGE"`
	// 是否压缩
	Compress bool `toml:"compress" json:"compress" yaml:"compress"  env:"COMPRESS"`
}

func (f *File) FileWriter(appLogFileName string) io.Writer {
	return &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s.log", f.DirPath, appLogFileName),
		MaxSize:    f.MaxSize,
		MaxAge:     f.MaxAge,
		MaxBackups: f.MaxBackups,
		Compress:   f.Compress,
	}
}

type Config struct {
	// 0 为打印日志全路径, 默认打印2层路径
	CallerDeep int `toml:"caller_deep" json:"caller_deep" yaml:"caller_deep"  env:"CALLER_DEEP"`
	// 日志的级别, 默认Debug
	Level string `toml:"level" json:"level" yaml:"level"  env:"LEVEL"`
	// 开启Trace时, 记录的TraceId名称, 默认trace_id
	TraceFiled string `toml:"trace_filed" json:"trace_filed" yaml:"trace_filed"  env:"TRACE_FILED"`

	// 控制台日志配置
	Console Console `toml:"console" json:"console" yaml:"console" envPrefix:"CONSOLE_"`
	// 日志文件配置
	File File `toml:"file" json:"file" yaml:"file" envPrefix:"FILE_"`

	ioc.ObjectImpl
	root    *zerolog.Logger
	lock    sync.Mutex
	loggers map[string]*zerolog.Logger
}

func (m *Config) Name() string {
	return AppName
}

// Trace加载后才加载日志
func (i *Config) Priority() int {
	return 997
}

func (m *Config) Init() error {
	var writers []io.Writer
	if m.Console.Enable {
		writers = append(writers, m.Console.ConsoleWriter())
	}
	if m.File.Enable {
		writers = append(writers, m.File.FileWriter(application.Get().GetAppName()))
	}

	if len(writers) == 0 {
		return nil
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	root := zerolog.New(io.MultiWriter(writers...)).With().Timestamp()
	if m.CallerDeep > 0 {
		root = root.Caller()
		zerolog.CallerMarshalFunc = m.CallerMarshalFunc
	}
	level, err := zerolog.ParseLevel(m.Level)
	if err != nil {
		return err
	}
	m.SetRoot(root.Logger().Level(level))
	return nil
}

func (m *Config) CallerMarshalFunc(pc uintptr, file string, line int) string {
	if m.CallerDeep == 0 {
		return file
	}

	short := file
	count := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			count++
		}

		if count >= m.CallerDeep {
			break
		}
	}
	file = short
	return file + ":" + strconv.Itoa(line)
}

func (m *Config) SetRoot(r zerolog.Logger) {
	m.root = &r
}

func (m *Config) Logger(name string) *zerolog.Logger {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.loggers[name]; !ok {
		l := m.root.With().Str(SUB_LOGGER_KEY, name).Logger()
		m.loggers[name] = &l
	}

	return m.loggers[name]
}

func (m *Config) IsMyLogger(logger *zerolog.Logger) bool {
	if logger == nil || m.root == nil {
		return false
	}

	// 快速检查：是否是 root logger
	if logger == m.root {
		return true
	}

	// 遍历所有存储的 logger 进行比较
	for _, storedLogger := range m.loggers {
		if storedLogger == logger {
			return true
		}
	}

	return false
}
