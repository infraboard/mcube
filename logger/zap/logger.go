package zap

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/infraboard/mcube/logger"
)

// LogOption configures a Logger.
type LogOption = zap.Option

// Logger logs messages to the configured output.
type Logger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

// Meta 日志meta数据
type Meta = map[string]interface{}

// NewLogger returns a new Logger labeled with the name of the selector. This
// should never be used from any global contexts, otherwise you will receive a
// no-op Logger. This is because the logp package needs to be initialized first.
// Instead create new Logger instance that your object reuses. Or if you need to
// log from a static context then you may use logp.L().Infow(), for example.
func NewLogger(selector string, options ...LogOption) *Logger {
	return newLogger(loadLogger().rootLogger, selector, options...)
}

func newLogger(rootLogger *zap.Logger, selector string, options ...LogOption) *Logger {
	log := rootLogger.
		WithOptions(zap.AddCallerSkip(1)).
		WithOptions(options...).
		Named(selector)
	return &Logger{log, log.Sugar()}
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (l *Logger) With(fields ...logger.Field) logger.Logger {
	zap.NewNop()
	sugar := l.sugar.With(transfer(fields)...)
	return &Logger{sugar.Desugar(), sugar}
}

// WithNamespace todo
func (l *Logger) WithNamespace(namespace string) *Logger {
	zap.NewNop()
	sugar := l.sugar.With(Namespace(namespace))
	return &Logger{sugar.Desugar(), sugar}
}

// Named adds a new path segment to the logger's name. Segments are joined by
// periods.
func (l *Logger) Named(name string) logger.Logger {
	logger := l.logger.Named(name)
	return &Logger{logger, logger.Sugar()}
}

// Print uses fmt.Sprint to construct and log a message.
func (l *Logger) Print(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Println todo
func (l *Logger) Println(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Debug uses fmt.Sprint to construct and log a message.
func (l *Logger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Info uses fmt.Sprint to construct and log a message.
func (l *Logger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func (l *Logger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

// Error uses fmt.Sprint to construct and log a message.
func (l *Logger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit(1).
func (l *Logger) Fatal(args ...interface{}) {
	l.sugar.Fatal(args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l *Logger) Panic(args ...interface{}) {
	l.sugar.Panic(args...)
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics.
func (l *Logger) DPanic(args ...interface{}) {
	l.sugar.DPanic(args...)
}

// IsDebug checks to see if the given logger is Debug enabled.
func (l *Logger) IsDebug() bool {
	return l.logger.Check(zapcore.DebugLevel, "") != nil
}

// Printf todo
func (l *Logger) Printf(format string, args ...interface{}) {
	l.sugar.Debugf(format, args...)
}

// Debugf uses fmt.Sprintf to construct and log a message.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.sugar.Debugf(format, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.sugar.Infof(format, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.sugar.Warnf(format, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.sugar.Errorf(format, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit(1).
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.sugar.Fatalf(format, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.sugar.Panicf(format, args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics.
func (l *Logger) DPanicf(format string, args ...interface{}) {
	l.sugar.DPanicf(format, args...)
}

// With context (reflection based)

// Debugw logs a message with some additional context. The additional context
// is added in the form of key-value pairs. The optimal way to write the value
// to the log message will be inferred by the value's type. To explicitly
// specify a type you can pass a Field such as logp.Stringer.
func (l *Logger) Debugw(msg string, fields ...logger.Field) {
	l.sugar.Debugw(msg, transfer(fields)...)
}

// Infow logs a message with some additional context. The additional context
// is added in the form of key-value pairs. The optimal way to write the value
// to the log message will be inferred by the value's type. To explicitly
// specify a type you can pass a Field such as logp.Stringer.
func (l *Logger) Infow(msg string, fields ...logger.Field) {
	l.sugar.Infow(msg, transfer(fields)...)
}

// Warnw logs a message with some additional context. The additional context
// is added in the form of key-value pairs. The optimal way to write the value
// to the log message will be inferred by the value's type. To explicitly
// specify a type you can pass a Field such as logp.Stringer.
func (l *Logger) Warnw(msg string, fields ...logger.Field) {
	l.sugar.Warnw(msg, transfer(fields)...)
}

// Errorw logs a message with some additional context. The additional context
// is added in the form of key-value pairs. The optimal way to write the value
// to the log message will be inferred by the value's type. To explicitly
// specify a type you can pass a Field such as logp.Stringer.
func (l *Logger) Errorw(msg string, fields ...logger.Field) {
	l.sugar.Errorw(msg, transfer(fields)...)
}

// Fatalw logs a message with some additional context, then calls os.Exit(1).
// The additional context is added in the form of key-value pairs. The optimal
// way to write the value to the log message will be inferred by the value's
// type. To explicitly specify a type you can pass a Field such as
// logp.Stringer.
func (l *Logger) Fatalw(msg string, fields ...logger.Field) {
	l.sugar.Fatalw(msg, transfer(fields)...)
}

// Panicw logs a message with some additional context, then panics. The
// additional context is added in the form of key-value pairs. The optimal way
// to write the value to the log message will be inferred by the value's type.
// To explicitly specify a type you can pass a Field such as logp.Stringer.
func (l *Logger) Panicw(msg string, fields ...logger.Field) {
	l.sugar.Panicw(msg, transfer(fields)...)
}

// DPanicw logs a message with some additional context. The logger panics only
// in Development mode.  The additional context is added in the form of
// key-value pairs. The optimal way to write the value to the log message will
// be inferred by the value's type. To explicitly specify a type you can pass a
// Field such as logp.Stringer.
func (l *Logger) DPanicw(msg string, fields ...logger.Field) {
	l.sugar.DPanicw(msg, transfer(fields)...)
}

// Recover stops a panicking goroutine and logs an Error.
func (l *Logger) Recover(msg string) {
	if r := recover(); r != nil {
		msg := fmt.Sprintf("%s. Recovering, but please report this.", msg)
		globalLogger().WithOptions(zap.AddCallerSkip(1)).
			Error(msg, zap.Any("panic", r), zap.Stack("stack"))
	}
}

// SetLevel 动态设置日志等级
func SetLevel(lv Level) {
	loadLogger().atom.SetLevel(lv.zapLevel())
}

// L returns an unnamed global logger.
func L() *Logger {
	return loadLogger().logger
}

// IsDebug returns true if the given selector would be logged.
// Deprecated: Use logp.NewLogger.
func IsDebug(selector string) bool {
	return globalLogger().Named(selector).Check(zap.DebugLevel, "") != nil
}

// transfer 用于方便zap记录
func transfer(m []logger.Field) (ma []interface{}) {
	for i := range m {
		ma = append(ma, zap.Any(m[i].Key, m[i].Value))
	}

	return
}
