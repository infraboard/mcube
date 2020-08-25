package response

import (
	"github.com/infraboard/mcube/bus/event"
)

// ResourceEvent 资源事件
type ResourceEvent interface {
	ResourceType() string
	ResourceAction() string

	ResourceUUID() string
	ResourceDomain() string
	ResourceNamespace() string
	ResourceName() string
	ResourceData() interface{}
}

func newResourceEvent(re ResourceEvent) *event.Event {
	e := event.NewEvent()
	e.Label["domain"] = re.ResourceDomain()
	e.Label["namespace"] = re.ResourceNamespace()
	e.Type = event.OperateEventType

	body := &event.OperateEvent{
		ResourceType: re.ResourceType(),
		ResourceUUID: re.ResourceUUID(),
		ResourceName: re.ResourceName(),
		Action:       re.ResourceAction(),
		Data:         re.ResourceData(),
	}
	e.Body = body

	return e
}
