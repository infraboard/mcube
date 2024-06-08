package desense

import (
	"strconv"
	"strings"
)

func init() {
	Registry(DefaultDesenser, &defaultDesenser{})
}

const (
	DefaultDesenser = "default"
)

// 脱敏器
type defaultDesenser struct {
}

func (d *defaultDesenser) DeSense(value string, args ...string) string {
	MaintainPrefixCharLength, MaintainSubfixCharLength := d.parseArgs(args)
	TotalMaintainCharLen := MaintainPrefixCharLength + MaintainSubfixCharLength
	SenseCharNumber := len(value) - TotalMaintainCharLen
	if SenseCharNumber <= 0 {
		return value
	}

	return d.MaintainPrefixString(value, MaintainPrefixCharLength) +
		strings.Repeat("*", SenseCharNumber) +
		d.MaintainSubfixString(value, MaintainSubfixCharLength)
}

// MaintainPrefixCharLength 保留前缀字符的长度
// MaintainSubfixCharLength 保留后缀字符的长度
func (d *defaultDesenser) parseArgs(args []string) (
	MaintainPrefixCharLength int,
	MaintainSubfixCharLength int,
) {
	if len(args) > 0 {
		MaintainPrefixCharLength, _ = strconv.Atoi(args[0])
	}

	if len(args) > 1 {
		MaintainSubfixCharLength, _ = strconv.Atoi(args[1])
	}
	return
}

func (d *defaultDesenser) MaintainSubfixString(value string, MaintainSubfixCharLength int) string {
	subfix := len(value) - MaintainSubfixCharLength
	if subfix <= 0 {
		return value
	}

	return value[subfix:]
}

func (d *defaultDesenser) MaintainPrefixString(value string, MaintainPrefixCharLength int) string {
	prefix := len(value) - MaintainPrefixCharLength
	if prefix <= 0 {
		return value
	}

	return value[:MaintainPrefixCharLength]
}
