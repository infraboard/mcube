package ioc

import "log"

var (
	_debug = true
)

func SetDebug(debug bool) {
	_debug = debug
}

func debug(format string, v ...any) {
	if !_debug {
		return
	}
	log.Printf(format, v...)
}
