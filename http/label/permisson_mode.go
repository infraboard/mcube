package label

import "github.com/infraboard/mcube/pb/http"

const (
	// Action key name
	PERMISSION_MODE = "permission_mode"
)

var (
	// Get Label
	PERMISSION_MODE_PRBAC = NewPermissionMode("PRBAC")
	// List label
	PERMISSION_MODE_ACL = NewPermissionMode("ACL")
)

// Action action Lable
func NewPermissionMode(name string) *http.Label {
	return permissionMode(name)
}

func permissionMode(value string) *http.Label {
	return http.NewLable(PERMISSION_MODE, value)
}
