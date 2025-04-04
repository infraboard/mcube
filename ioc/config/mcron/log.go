package mcron

import "github.com/infraboard/mcube/v2/ioc/config/log"

type LogWrapper struct {
}

// Info logs routine messages about cron's operation.
func (l *LogWrapper) Info(msg string, keysAndValues ...interface{}) {
	log.Sub(APP_NAME).Info().Msg(msg)
}

// Error logs an error condition.
func (l *LogWrapper) Error(err error, msg string, keysAndValues ...interface{}) {
	log.Sub(APP_NAME).Error().Msgf("%s, %s", err.Error(), msg)
}
