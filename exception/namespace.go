package exception

var (
	usedNamespace = GlobalNamespace
)

// Namespace 异常所属空间
type Namespace string

func (ns Namespace) String() string {
	return string(ns)
}

const (
	// GlobalNamespace 所有服务公用的一些异常
	GlobalNamespace = Namespace("global")
)
