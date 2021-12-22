package inject_tag

import (
	"log"
)

var verbose = false

func logf(format string, v ...interface{}) {
	if !verbose {
		return
	}
	log.Printf(format, v...)
}
