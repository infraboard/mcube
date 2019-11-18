package uuid

// Generator uuid 生成器
type Generator interface {
	NextID() (string error)
}
