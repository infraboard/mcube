package zap_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
)

func TestLogger(t *testing.T) {
	exerciseLogger := func() {
		meta := map[string]interface{}{"x": 1, "y": 1}
		log := zap.NewLogger("example")
		log.Info("unnamed global logger")
		log.Info("some message")
		log.Infof("some message with parameter x=%v, y=%v", 1, 2)
		log.Infow("some message", logger.NewFieldsFromKV(meta)...)
		log.Infow("", logger.NewAny("empty_message", true))

		// Add context.
		log.With(logger.NewFieldsFromKV(meta)...).Warn("logger with context")

		someStruct := struct {
			X int `json:"x"`
			Y int `json:"y"`
		}{1, 2}
		log.Infow("some message with struct value", logger.NewAny("metrics", someStruct))
	}

	zap.TestingSetup()
	exerciseLogger()
	zap.TestingSetup(zap.AsJSON())
	exerciseLogger()
}

func TestLoggerSelectors(t *testing.T) {
	if err := zap.DevelopmentSetup(zap.WithSelectors("good", " padded "), zap.ToObserverOutput()); err != nil {
		t.Fatal(err)
	}

	assert.True(t, zap.HasSelector("padded"))

	good := zap.NewLogger("good")
	bad := zap.NewLogger("bad")

	good.Debug("is logged")
	logs := zap.ObserverLogs().TakeAll()
	assert.Len(t, logs, 1)

	// Selectors only apply to debug level logs.
	bad.Debug("not logged")
	logs = zap.ObserverLogs().TakeAll()
	assert.Len(t, logs, 0)

	bad.Info("is also logged")
	logs = zap.ObserverLogs().TakeAll()
	assert.Len(t, logs, 1)
}

func TestLoggerLevel(t *testing.T) {
	if err := zap.DevelopmentSetup(zap.ToObserverOutput()); err != nil {
		t.Fatal(err)
	}

	const loggerName = "tester"
	logger := zap.NewLogger(loggerName)

	logger.Debug("debug")
	logs := zap.ObserverLogs().TakeAll()
	if assert.Len(t, logs, 1) {
		assert.Equal(t, loggerName, logs[0].LoggerName)
		assert.Equal(t, "debug", logs[0].Message)
	}

	logger.Info("info")
	logs = zap.ObserverLogs().TakeAll()
	if assert.Len(t, logs, 1) {
		assert.Equal(t, loggerName, logs[0].LoggerName)
		assert.Equal(t, "info", logs[0].Message)
	}

	logger.Warn("warn")
	logs = zap.ObserverLogs().TakeAll()
	if assert.Len(t, logs, 1) {
		assert.Equal(t, loggerName, logs[0].LoggerName)
		assert.Equal(t, "warn", logs[0].Message)
	}

	logger.Error("error")
	logs = zap.ObserverLogs().TakeAll()
	if assert.Len(t, logs, 1) {
		assert.Equal(t, loggerName, logs[0].LoggerName)
		assert.Equal(t, "error", logs[0].Message)
	}
}

func TestRecover(t *testing.T) {
	const recoveryExplanation = "Something went wrong"
	const cause = "unexpected condition"

	zap.DevelopmentSetup(zap.ToObserverOutput())

	defer func() {
		logs := zap.ObserverLogs().TakeAll()

		if assert.Len(t, logs, 1) {
			log := logs[0]
			assert.Equal(t, "zap/core_test.go",
				strings.Split(log.Caller.TrimmedPath(), ":")[0])
			assert.Contains(t, log.Message, recoveryExplanation+
				". Recovering, but please report this.")
			assert.Contains(t, log.ContextMap(), "panic")
		}
	}()

	defer zap.L().Recover(recoveryExplanation)
	panic(cause)
}

func TestHasSelector(t *testing.T) {
	zap.DevelopmentSetup(zap.WithSelectors("*", "config"))
	assert.True(t, zap.HasSelector("config"))
	assert.False(t, zap.HasSelector("publish"))
}

func TestIsDebug(t *testing.T) {
	zap.DevelopmentSetup()
	assert.True(t, zap.IsDebug("all"))

	zap.DevelopmentSetup(zap.WithSelectors("*"))
	assert.True(t, zap.IsDebug("all"))

	zap.DevelopmentSetup(zap.WithSelectors("only_this"))
	assert.False(t, zap.IsDebug("all"))
	assert.True(t, zap.IsDebug("only_this"))
}

func TestL(t *testing.T) {
	if err := zap.DevelopmentSetup(zap.ToObserverOutput()); err != nil {
		t.Fatal(err)
	}

	zap.L().Infow("infow", logger.NewAny("rate", 1))
	logs := zap.ObserverLogs().TakeAll()
	if assert.Len(t, logs, 1) {
		log := logs[0]
		assert.Equal(t, "", log.LoggerName)
		assert.Equal(t, "infow", log.Message)
		assert.Contains(t, log.ContextMap(), "rate")
	}

	const loggerName = "tester"
	zap.L().Named(loggerName).Warnf("warning %d", 1)
	logs = zap.ObserverLogs().TakeAll()
	if assert.Len(t, logs, 1) {
		log := logs[0]
		assert.Equal(t, loggerName, log.LoggerName)
		assert.Equal(t, "warning 1", log.Message)
	}
}

func TestChangeLevel(t *testing.T) {
	if err := zap.DevelopmentSetup(zap.ToObserverOutput()); err != nil {
		t.Fatal(err)
	}

	const loggerName = "tester"
	logger := zap.NewLogger(loggerName)

	logger.Debug("debug")
	logs := zap.ObserverLogs().TakeAll()
	if assert.Len(t, logs, 1) {
		assert.Equal(t, loggerName, logs[0].LoggerName)
		assert.Equal(t, "debug", logs[0].Message)
	}

	zap.SetLevel(zap.ErrorLevel)
	logger.Info("info")
	logs = zap.ObserverLogs().TakeAll()
	assert.Len(t, logs, 0)

	logger.Warn("warn")
	logs = zap.ObserverLogs().TakeAll()
	assert.Len(t, logs, 0)

	logger.Error("error")
	logs = zap.ObserverLogs().TakeAll()
	if assert.Len(t, logs, 1) {
		assert.Equal(t, loggerName, logs[0].LoggerName)
		assert.Equal(t, "error", logs[0].Message)
	}
}
