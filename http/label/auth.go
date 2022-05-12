package label

const (
	// AuthLabelKey 控制是否开启认证
	AuthLabelKey = "auth"

	// PermissionLabelKey 控制是否开启权限判定
	PermissionLabelKey = "permission"

	// AllowLabelKey 控制允许的角色
	AllowLabelKey = "allow"

	// 控制是否开启审计
	AuditLabelKey = "audit"
)

const (
	Enable  = true
	Disable = false
)
