package desense_test

import (
	"testing"

	"github.com/infraboard/mcube/v2/desense"
)

func TestMaskPhone(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected string
	}{
		{"标准手机号", "13812341234", "138****1234"},
		{"短号码", "123456", "123456"},
		{"固定电话", "01012345678", "010****5678"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := desense.MaskPhone(tt.phone)
			if result != tt.expected {
				t.Errorf("MaskPhone(%s) = %s, want %s", tt.phone, result, tt.expected)
			}
		})
	}
}

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected string
	}{
		{"普通邮箱", "test@example.com", "te**@example.com"},
		{"长邮箱", "testuser@example.com", "tes*****@example.com"},
		{"短邮箱", "ab@example.com", "a*@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := desense.MaskEmail(tt.email)
			if result != tt.expected {
				t.Errorf("MaskEmail(%s) = %s, want %s", tt.email, result, tt.expected)
			}
		})
	}
}

func TestMaskIDCard(t *testing.T) {
	tests := []struct {
		name     string
		idcard   string
		expected string
	}{
		{"18位身份证", "110101199001011234", "110***********1234"},
		{"15位身份证", "110101900101123", "110********1123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := desense.MaskIDCard(tt.idcard)
			if result != tt.expected {
				t.Errorf("MaskIDCard(%s) = %s, want %s", tt.idcard, result, tt.expected)
			}
		})
	}
}

func TestMaskBankCard(t *testing.T) {
	tests := []struct {
		name     string
		card     string
		expected string
	}{
		{"19位银行卡", "6222021234567890123", "6222 **** **** 0123"},
		{"16位银行卡", "6222021234567890", "6222 **** **** 7890"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := desense.MaskBankCard(tt.card)
			if result != tt.expected {
				t.Errorf("MaskBankCard(%s) = %s, want %s", tt.card, result, tt.expected)
			}
		})
	}
}

func TestMaskName(t *testing.T) {
	tests := []struct {
		name     string
		username string
		expected string
	}{
		{"两字姓名", "张三", "张*"},
		{"三字姓名", "李小明", "李**"},
		{"四字姓名", "欧阳娜娜", "欧***"},
		{"单字", "李", "李"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := desense.MaskName(tt.username)
			if result != tt.expected {
				t.Errorf("MaskName(%s) = %s, want %s", tt.username, result, tt.expected)
			}
		})
	}
}

func TestMaskAddress(t *testing.T) {
	tests := []struct {
		name     string
		address  string
		expected string
	}{
		{"完整地址", "北京市朝阳区某某街道123号", "北京市朝阳区****"},
		{"短地址", "北京市", "北**"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := desense.MaskAddress(tt.address)
			if result != tt.expected {
				t.Errorf("MaskAddress(%s) = %s, want %s", tt.address, result, tt.expected)
			}
		})
	}
}

func TestMaskPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected string
	}{
		{"普通密码", "password123", "******"},
		{"空密码", "", ""},
		{"长密码", "verylongpassword123456", "******"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := desense.MaskPassword(tt.password)
			if result != tt.expected {
				t.Errorf("MaskPassword(%s) = %s, want %s", tt.password, result, tt.expected)
			}
		})
	}
}

func TestMaskString(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		strategy string
		args     []string
		expected string
	}{
		{"使用phone策略", "13812341234", "phone", nil, "138****1234"},
		{"使用email策略", "test@example.com", "email", nil, "te**@example.com"},
		{"使用default策略", "abcdefghijkl", "default", []string{"3", "2"}, "abc*******kl"},
		{"不存在的策略", "test", "nonexistent", []string{"1", "1"}, "t**t"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := desense.MaskString(tt.value, tt.strategy, tt.args...)
			if result != tt.expected {
				t.Errorf("MaskString(%s, %s) = %s, want %s", tt.value, tt.strategy, result, tt.expected)
			}
		})
	}
}

type TestUserWithStrategies struct {
	Name     string `json:"name" mask:"name"`
	Phone    string `json:"phone" mask:"phone"`
	Email    string `json:"email" mask:"email"`
	IDCard   string `json:"idcard" mask:"idcard"`
	BankCard string `json:"bankcard" mask:"bankcard"`
	Address  string `json:"address" mask:"address"`
	Password string `json:"password" mask:"password"`
}

func TestMaskStructWithStrategies(t *testing.T) {
	user := &TestUserWithStrategies{
		Name:     "张三",
		Phone:    "13812341234",
		Email:    "test@example.com",
		IDCard:   "110101199001011234",
		BankCard: "6222021234567890123",
		Address:  "北京市朝阳区某某街道123号",
		Password: "password123",
	}

	if err := desense.MaskStruct(user); err != nil {
		t.Fatal(err)
	}

	if user.Name != "张*" {
		t.Errorf("Name = %s, want 张*", user.Name)
	}
	if user.Phone != "138****1234" {
		t.Errorf("Phone = %s, want 138****1234", user.Phone)
	}
	if user.Email != "te**@example.com" {
		t.Errorf("Email = %s, want te**@example.com", user.Email)
	}
	if user.IDCard != "110***********1234" {
		t.Errorf("IDCard = %s, want 110***********1234", user.IDCard)
	}
	if user.BankCard != "6222 **** **** 0123" {
		t.Errorf("BankCard = %s, want 6222 **** **** 0123", user.BankCard)
	}
	if user.Address != "北京市朝阳区****" {
		t.Errorf("Address = %s, want 北京市朝阳区****", user.Address)
	}
	if user.Password != "******" {
		t.Errorf("Password = %s, want ******", user.Password)
	}

	t.Logf("脱敏后: %+v", user)
}

func TestBackwardCompatibility(t *testing.T) {
	type OldStyleStruct struct {
		Phone string `json:"phone" mask:"default,3,4"`
	}

	old := &OldStyleStruct{
		Phone: "13812341234",
	}

	if err := desense.MaskStruct(old); err != nil {
		t.Fatal(err)
	}

	if old.Phone != "138****1234" {
		t.Errorf("Phone = %s, want 138****1234", old.Phone)
	}
}

func BenchmarkMaskPhone(b *testing.B) {
	for i := 0; i < b.N; i++ {
		desense.MaskPhone("13812341234")
	}
}

func BenchmarkMaskEmail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		desense.MaskEmail("test@example.com")
	}
}

func BenchmarkMaskStruct(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testUser := &TestUserWithStrategies{
			Name:     "张三",
			Phone:    "13812341234",
			Email:    "test@example.com",
			IDCard:   "110101199001011234",
			BankCard: "6222021234567890123",
			Address:  "北京市朝阳区某某街道123号",
			Password: "password123",
		}
		desense.MaskStruct(testUser)
	}
}
