// +build windows nacl plan9

package zap

import (
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

func newSyslog(_ zapcore.Encoder, _ zapcore.LevelEnabler) (zapcore.Core, error) {
	return nil, errors.New("syslog is not supported on this OS")
}
