package response

// ResourceEvent 资源事件
type ResourceEvent interface {
	ResourceType() string
	ResourceUUID() string
	ResourceDomain() string
	ResourceNamespace() string
	ResourceName() string
	ResourceAction() string
}
