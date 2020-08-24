package response

import (
	"github.com/infraboard/mcube/bus/event"
)

// ResourceEvent 资源事件
type ResourceEvent interface {
	ResourceType() string
	ResourceUUID() string
	ResourceDomain() string
	ResourceNamespace() string
	ResourceName() string
	ResourceAction() string
	ResourceData() interface{}
}

func newEvent(re ResourceEvent) *event.Event {
	e := event.NewEvent()
	e.Label["domain"] = re.ResourceDomain()
	e.Label["namespace"] = re.ResourceNamespace()
	return e
}
