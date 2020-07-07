package label

import (
	"github.com/infraboard/mcube/http/router"
)

const (
	// ActionLableKey key name
	ActionLableKey = "action"
)

var (
	// Get Label
	Get = action("get")
	// List label
	List = action("list")
	// Create label
	Create = action("create")
	// Update label
	Update = action("update")
	// Delete label
	Delete = action("delete")
)

// Action action Lable
func Action(name string) *router.Label {
	return action(name)
}

func action(value string) *router.Label {
	return router.NewLable(ActionLableKey, value)
}
