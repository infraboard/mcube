package logger

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/infraboard/mcube/ioc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	ioc.Config().Registry(&Config{
		CallerDeep: 3,
		Level:      zerolog.DebugLevel,
		Console: Console{
			Enable:  true,
			NoColor: false,
		},
		File: File{
			Enable:     false,
			MaxSize:    100,
			MaxBackups: 6,
		},
		loggers: make(map[string]*zerolog.Logger),
	})
}

type Console struct {
	Enable  bool `toml:"enable" json:"enable" yaml:"enable"  env:"LOG_TO_CONSOLE"`
	NoColor bool `toml:"no_color" json:"no_color" yaml:"no_color"  env:"LOG_CONSOLE_NO_COLOR"`
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
	Enable bool `toml:"enable" json:"enable" yaml:"enable"  env:"LOG_TO_FILE"`
	// 文件的路径
	FilePath string `toml:"file_path" json:"file_path" yaml:"file_path"  env:"LOG_FILE_PATH"`
	// 单位M, 默认100M
	MaxSize int `toml:"max_size" json:"max_size" yaml:"max_size"  env:"LOG_FILE_MAX_SIZE"`
	// 默认保存 6个文件
	MaxBackups int `toml:"max_backups" json:"max_backups" yaml:"max_backups"  env:"LOG_FILE_MAX_BACKUPS"`
	// 保存多久
	MaxAge int `toml:"max_age" json:"max_age" yaml:"max_age"  env:"LOG_FILE_MAX_AGE"`
	// 是否压缩
	Compress bool `toml:"compress" json:"compress" yaml:"compress"  env:"LOG_FILE_COMPRESS"`
}

func (f *File) FileWriter() io.Writer {
	return &lumberjack.Logger{
		Filename:   f.FilePath,
		MaxSize:    f.MaxSize,
		MaxAge:     f.MaxAge,
		MaxBackups: f.MaxBackups,
		Compress:   f.Compress,
	}
}

type Config struct {
	// 0 为打印日志全路径, 默认打印2层路径
	CallerDeep int `toml:"caller_deep" json:"caller_deep" yaml:"caller_deep"  env:"LOG_CALLER_DEEP"`
	// 日志的级别, 默认Debug
	Level zerolog.Level `toml:"level" json:"level" yaml:"level"  env:"LOG_LEVEL"`

	// 控制台日志配置
	Console Console `toml:"console" json:"console" yaml:"console"`
	// 日志文件配置
	File File `toml:"file" json:"file" yaml:"file"`

	ioc.ObjectImpl
	root    *zerolog.Logger
	lock    sync.Mutex
	loggers map[string]*zerolog.Logger
}

func (m *Config) Name() string {
	return LOG
}

func (m *Config) Init() error {
	var writers []io.Writer
	if m.Console.Enable {
		writers = append(writers, m.Console.ConsoleWriter())
	}
	if m.File.Enable {
		writers = append(writers, m.File.FileWriter())
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
	m.SetRoot(root.Logger().Level(m.Level))
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
