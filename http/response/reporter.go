package response

import (
	"github.com/infraboard/mcube/bus"
	"github.com/infraboard/mcube/bus/event"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
)

var (
	eReporter bus.Publisher
	log       logger.Logger
)

func getLog() logger.Logger {
	if log == nil {
		log = zap.L().Named("Response")
	}

	return log
}

// SetEventReporter tood
func SetEventReporter(pub bus.Publisher) {
	eReporter = pub
}

// HasEventReporter 是否已经初始化
func HasEventReporter() bool {
	return eReporter != nil
}

func sendEvent(re ResourceEvent) {
	if !HasEventReporter() {
		getLog().Errorf("event reporter is nil")
		return
	}

	e := event.NewEvent()
	eReporter.Pub("xxx", e)
}
