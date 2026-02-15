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
	MaintainPrefixCharLength, MaintainSubfixCharLength := d.parseArgs(value, args)
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
// 如果没有提供参数，会根据value的长度智能选择合适的默认值
func (d *defaultDesenser) parseArgs(value string, args []string) (
	MaintainPrefixCharLength int,
	MaintainSubfixCharLength int,
) {
	// 如果提供了参数，使用用户指定的值
	if len(args) > 0 && args[0] != "" {
		MaintainPrefixCharLength, _ = strconv.Atoi(args[0])
	}

	if len(args) > 1 && args[1] != "" {
		MaintainSubfixCharLength, _ = strconv.Atoi(args[1])
	}

	// 如果没有提供任何参数，或者两个参数都是0，使用智能默认值
	if len(args) == 0 || (MaintainPrefixCharLength == 0 && MaintainSubfixCharLength == 0) {
		MaintainPrefixCharLength, MaintainSubfixCharLength = d.getSmartDefaults(value)
	}

	return
}

// getSmartDefaults 根据字符串长度智能选择合适的脱敏参数
func (d *defaultDesenser) getSmartDefaults(value string) (prefixLen, suffixLen int) {
	length := len(value)

	switch {
	case length <= 4:
		// 太短，保留第一个字符
		return 1, 0
	case length <= 6:
		// 短字符串，保留前1后1
		return 1, 1
	case length <= 10:
		// 中等长度，保留前2后2
		return 2, 2
	case length == 11:
		// 标准手机号，保留前3后4
		return 3, 4
	case length <= 17:
		// 中长字符串，保留前3后4
		return 3, 4
	case length == 18:
		// 标准身份证，保留前3后4
		return 3, 4
	default:
		// 长字符串（如银行卡），保留前4后4
		return 4, 4
	}
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
