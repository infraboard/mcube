package datasource

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
)

func newGormCustomLogger(log *zerolog.Logger) logger.Interface {
	return &gormCustomLogger{
		log:                       log,
		IgnoreRecordNotFoundError: true,
		SlowThreshold:             200 * time.Millisecond,
	}
}

type gormCustomLogger struct {
	log *zerolog.Logger

	LogLevel                  logger.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
}

// LogMode implements logger.Interface.
func (g *gormCustomLogger) LogMode(l logger.LogLevel) logger.Interface {
	g.LogLevel = l
	return g
}

// Error implements logger.Interface.
func (g *gormCustomLogger) Error(ctx context.Context, msg string, datas ...interface{}) {
	g.log.Info().Fields(datas).Msg(msg)
}

// Info implements logger.Interface.
func (g *gormCustomLogger) Info(ctx context.Context, msg string, datas ...interface{}) {
	g.log.Info().Fields(datas).Msg(msg)
}

// Warn implements logger.Interface.
func (g *gormCustomLogger) Warn(ctx context.Context, msg string, datas ...interface{}) {
	g.log.Warn().Fields(datas).Msg(msg)
}

// Trace implements logger.Interface.
func (g *gormCustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if g.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && g.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !g.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			g.log.Info().Msgf("%s [%.3fms] [rows:%v] %s", err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.log.Info().Msgf("%s [%.3fms] [rows:%v] %s", err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", g.SlowThreshold)
		if rows == -1 {
			g.log.Warn().Msgf("%s [%.3fms] [rows:%v] %s", slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.log.Warn().Msgf("%s [%.3fms] [rows:%v] %s", slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case g.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			g.log.Debug().Msgf("[%.3fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.log.Debug().Msgf("[%.3fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
