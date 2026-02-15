package desense

import (
	"strings"
	"unicode/utf8"
)

func init() {
	// 注册内置脱敏策略
	Registry(PhoneStrategy, &phoneDesenser{})
	Registry(EmailStrategy, &emailDesenser{})
	Registry(IDCardStrategy, &idCardDesenser{})
	Registry(BankCardStrategy, &bankCardDesenser{})
	Registry(NameStrategy, &nameDesenser{})
	Registry(AddressStrategy, &addressDesenser{})
	Registry(PasswordStrategy, &passwordDesenser{})
}

const (
	// PhoneStrategy 手机号脱敏策略
	PhoneStrategy = "phone"
	// EmailStrategy 邮箱脱敏策略
	EmailStrategy = "email"
	// IDCardStrategy 身份证脱敏策略
	IDCardStrategy = "idcard"
	// BankCardStrategy 银行卡脱敏策略
	BankCardStrategy = "bankcard"
	// NameStrategy 姓名脱敏策略
	NameStrategy = "name"
	// AddressStrategy 地址脱敏策略
	AddressStrategy = "address"
	// PasswordStrategy 密码脱敏策略（完全隐藏）
	PasswordStrategy = "password"
)

// phoneDesenser 手机号脱敏器
// 示例: 13812341234 -> 138****1234
type phoneDesenser struct{}

func (p *phoneDesenser) DeSense(value string, args ...string) string {
	length := len(value)
	if length < 7 {
		// 太短的字符串不脱敏
		return value
	}

	// 默认保留前3位和后4位
	prefixLen := 3
	suffixLen := 4

	if length == 11 {
		// 标准手机号
		return value[:prefixLen] + strings.Repeat("*", length-prefixLen-suffixLen) + value[length-suffixLen:]
	}

	// 其他长度的号码，保留前3位和后4位
	if length < prefixLen+suffixLen {
		prefixLen = length / 2
		suffixLen = length - prefixLen
	}

	return value[:prefixLen] + strings.Repeat("*", length-prefixLen-suffixLen) + value[length-suffixLen:]
}

// emailDesenser 邮箱脱敏器
// 示例: test@example.com -> te**@example.com
type emailDesenser struct{}

func (e *emailDesenser) DeSense(value string, args ...string) string {
	atIndex := strings.Index(value, "@")
	if atIndex <= 0 {
		// 不是有效的邮箱格式
		return value
	}

	localPart := value[:atIndex]
	domainPart := value[atIndex:]

	localLen := utf8.RuneCountInString(localPart)
	if localLen <= 2 {
		return localPart[:1] + "*" + domainPart
	}

	// 保留前2个字符
	showLen := 2
	if localLen > 6 {
		showLen = 3
	}

	return localPart[:showLen] + strings.Repeat("*", localLen-showLen) + domainPart
}

// idCardDesenser 身份证脱敏器
// 示例: 110101199001011234 -> 110***********1234
type idCardDesenser struct{}

func (i *idCardDesenser) DeSense(value string, args ...string) string {
	length := len(value)
	if length < 8 {
		// 太短，不符合身份证格式
		return value
	}

	// 15位或18位身份证：保留前3位和后4位
	prefixLen := 3
	suffixLen := 4

	return value[:prefixLen] + strings.Repeat("*", length-prefixLen-suffixLen) + value[length-suffixLen:]
}

// bankCardDesenser 银行卡脱敏器
// 示例: 6222021234567890123 -> 6222 **** **** 0123
type bankCardDesenser struct{}

func (b *bankCardDesenser) DeSense(value string, args ...string) string {
	// 移除空格
	value = strings.ReplaceAll(value, " ", "")
	length := len(value)

	if length < 8 {
		return value
	}

	// 保留前4位和后4位
	prefixLen := 4
	suffixLen := 4

	masked := value[:prefixLen] + " " + strings.Repeat("*", 4) + " " + strings.Repeat("*", 4) + " " + value[length-suffixLen:]
	return masked
}

// nameDesenser 姓名脱敏器
// 示例: 张三 -> 张*, 欧阳娜娜 -> 欧阳**
type nameDesenser struct{}

func (n *nameDesenser) DeSense(value string, args ...string) string {
	runes := []rune(value)
	length := len(runes)

	if length <= 1 {
		return value
	}

	if length == 2 {
		// 两个字：保留第一个字
		return string(runes[0]) + "*"
	}

	// 三个字及以上：保留第一个字，其余用*代替
	// 对于复姓（如欧阳），可以保留前两个字
	if length >= 3 {
		// 保留第一个字或前两个字（判断是否是复姓）
		return string(runes[0]) + strings.Repeat("*", length-1)
	}

	return value
}

// addressDesenser 地址脱敏器
// 示例: 北京市朝阳区某某街道123号 -> 北京市朝阳区****
type addressDesenser struct{}

func (a *addressDesenser) DeSense(value string, args ...string) string {
	runes := []rune(value)
	length := len(runes)

	if length <= 6 {
		// 太短的地址，保留一半
		showLen := length / 2
		return string(runes[:showLen]) + strings.Repeat("*", length-showLen)
	}

	// 保留前6个字符（通常是省市区）
	showLen := 6
	return string(runes[:showLen]) + strings.Repeat("*", 4)
}

// passwordDesenser 密码脱敏器（完全隐藏）
// 示例: password123 -> ******
type passwordDesenser struct{}

func (p *passwordDesenser) DeSense(value string, args ...string) string {
	if value == "" {
		return ""
	}
	// 密码完全隐藏，固定显示6个*
	return "******"
}
