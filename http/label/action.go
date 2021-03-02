package label

import (
	"github.com/infraboard/mcube/pb/http"
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
func Action(name string) *http.Label {
	return action(name)
}

func action(value string) *http.Label {
	return http.NewLable(ActionLableKey, value)
}
