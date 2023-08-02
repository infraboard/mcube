package sense

import "strings"

var (
	DefaultDesenser = NewDefaultDesenser()
	ShortDesenser   = NewShortDesenser()
)

func NewStdDesenser(prefix, subfix int) *StdDesenser {
	return &StdDesenser{
		MaintainPrefixCharLength: prefix,
		MaintainSubfixCharLength: subfix,
	}
}

func NewDefaultDesenser() *StdDesenser {
	return NewStdDesenser(3, 4)
}

func NewShortDesenser() *StdDesenser {
	return NewStdDesenser(2, 0)
}

// 脱敏器
type StdDesenser struct {
	// 保留前缀字符的长度
	MaintainPrefixCharLength int
	// 保留后缀字符的长度
	MaintainSubfixCharLength int
}

func (d *StdDesenser) TotalMaintainCharLen() int {
	return d.MaintainPrefixCharLength + d.MaintainSubfixCharLength
}

func (d *StdDesenser) MaintainSubfixString(value string) string {
	subfix := len(value) - d.MaintainSubfixCharLength
	if subfix <= 0 {
		return value
	}

	return value[subfix:]
}

func (d *StdDesenser) SenseCharNumber(value string) int {
	return len(value) - d.TotalMaintainCharLen()
}

func (d *StdDesenser) MaintainPrefixString(value string) string {
	prefix := len(value) - d.MaintainPrefixCharLength
	if prefix <= 0 {
		return value
	}

	return value[:d.MaintainPrefixCharLength]
}

func (d *StdDesenser) DeSense(value string) string {
	n := d.SenseCharNumber(value)
	if n <= 0 {
		return value
	}

	return d.MaintainPrefixString(value) + strings.Repeat("*", n) + d.MaintainSubfixString(value)
}
