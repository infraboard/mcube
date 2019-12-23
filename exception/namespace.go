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
	// Keyauth 权限系统, 每个子系统有自己的异常空间
	Keyauth = Namespace("keyauth")
)

// ChangeNameSapce 切换异常空间
func ChangeNameSapce(ns Namespace) {
	usedNamespace = ns
}
