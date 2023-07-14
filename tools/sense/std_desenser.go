package sense

import "strings"

var (
	DefaultDesenser = NewStdDesenser()
)

func NewStdDesenser() *StdDesenser {
	return &StdDesenser{
		MaintainPrefixCharLength: 4,
	}
}

// 脱敏器
type StdDesenser struct {
	// 保留前缀字符的长度
	MaintainPrefixCharLength int
}

func (d *StdDesenser) DeSense(value string) string {
	delta := len(value) - d.MaintainPrefixCharLength
	if delta <= 0 {
		return value
	}

	return value[:d.MaintainPrefixCharLength] + strings.Repeat("*", delta)
}
