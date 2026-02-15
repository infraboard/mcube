package desense

// MaskString 使用指定策略对字符串进行脱敏
// strategy: 脱敏策略名称（default, phone, email, idcard, bankcard, name, address, password）
// value: 待脱敏的字符串
// args: 策略参数（可选）
func MaskString(value, strategy string, args ...string) string {
	desenfor := Get(strategy)
	if desenfor == nil {
		// 策略不存在，使用默认策略
		desenfor = Default()
	}
	return desenfor.DeSense(value, args...)
}

// MaskPhone 手机号脱敏
// 示例: 13812341234 -> 138****1234
func MaskPhone(phone string) string {
	desenfor := Get(PhoneStrategy)
	if desenfor == nil {
		return phone
	}
	return desenfor.DeSense(phone)
}

// MaskEmail 邮箱脱敏
// 示例: test@example.com -> te**@example.com
func MaskEmail(email string) string {
	desenfor := Get(EmailStrategy)
	if desenfor == nil {
		return email
	}
	return desenfor.DeSense(email)
}

// MaskIDCard 身份证脱敏
// 示例: 110101199001011234 -> 110***********1234
func MaskIDCard(idcard string) string {
	desenfor := Get(IDCardStrategy)
	if desenfor == nil {
		return idcard
	}
	return desenfor.DeSense(idcard)
}

// MaskBankCard 银行卡脱敏
// 示例: 6222021234567890123 -> 6222 **** **** 0123
func MaskBankCard(card string) string {
	desenfor := Get(BankCardStrategy)
	if desenfor == nil {
		return card
	}
	return desenfor.DeSense(card)
}

// MaskName 姓名脱敏
// 示例: 张三 -> 张*
func MaskName(name string) string {
	desenfor := Get(NameStrategy)
	if desenfor == nil {
		return name
	}
	return desenfor.DeSense(name)
}

// MaskAddress 地址脱敏
// 示例: 北京市朝阳区某某街道123号 -> 北京市朝阳区****
func MaskAddress(address string) string {
	desenfor := Get(AddressStrategy)
	if desenfor == nil {
		return address
	}
	return desenfor.DeSense(address)
}

// MaskPassword 密码脱敏（完全隐藏）
// 示例: password123 -> ******
func MaskPassword(password string) string {
	desenfor := Get(PasswordStrategy)
	if desenfor == nil {
		return "******"
	}
	return desenfor.DeSense(password)
}
