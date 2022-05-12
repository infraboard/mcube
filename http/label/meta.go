package label

const (
	// ResourceLableKey key name
	ResourceLableKey = "resource"

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

type Meta map[string]interface{}

func (m Meta) Resource() string {
	if v, ok := m[ResourceLableKey]; ok {
		return v.(string)
	}

	return ""
}

func (m Meta) Action() string {
	if v, ok := m[ActionLableKey]; ok {
		return v.(string)
	}

	return ""
}

func (m Meta) AuthEnable() bool {
	if v, ok := m[AuthLabelKey]; ok {
		return v.(bool)
	}
	return false
}

func (m Meta) PermissionEnable() bool {
	if v, ok := m[PermissionLabelKey]; ok {
		return v.(bool)
	}
	return false
}

func (m Meta) AuditEnable() bool {
	if v, ok := m[AuditLabelKey]; ok {
		return v.(bool)
	}
	return false
}

func (m Meta) Allow() []string {
	if v, ok := m[AllowLabelKey]; ok {
		return v.([]string)
	}
	return []string{}
}
