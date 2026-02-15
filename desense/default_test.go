package desense_test

import (
	"testing"

	"github.com/infraboard/mcube/v2/desense"
	"github.com/infraboard/mcube/v2/tools/pretty"
)

func TestDeSense(t *testing.T) {
	v := desense.Default().DeSense("abcdefghijkl", "3", "2")
	t.Logf("%s", v)
}

// 测试不带参数的智能默认脱敏
func TestDeSenseWithoutArgs(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"手机号-无参数", "13812341234", "138****1234"},                 // 11位，默认3,4
		{"身份证-无参数", "110101199001011234", "110***********1234"},   // 18位，默认3,4
		{"银行卡-无参数", "6222021234567890123", "6222***********0123"}, // 19位，默认4,4
		{"短字符串", "abc", "a**"},                                    // 3位，默认1,0
		{"中等字符串", "abcdefgh", "ab****gh"},                         // 8位，默认2,2
		{"邮箱式字符串", "test@example.com", "tes*********.com"},        // 16位，默认3,4
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 不传参数，使用智能默认值
			result := desense.Default().DeSense(tt.value)
			if result != tt.expected {
				t.Errorf("DeSense(%s) = %s, want %s", tt.value, result, tt.expected)
			}
		})
	}
}

// 测试tag中使用default不带参数
func TestMaskStructWithDefaultNoArgs(t *testing.T) {
	type User struct {
		Phone  string `mask:"default"`     // 不带参数，智能识别
		IDCard string `mask:"default"`     // 不带参数，智能识别
		Custom string `mask:"default,5,2"` // 带参数，使用自定义
	}

	user := &User{
		Phone:  "13812341234",
		IDCard: "110101199001011234",
		Custom: "customvalue123",
	}

	if err := desense.MaskStruct(user); err != nil {
		t.Fatal(err)
	}

	// 手机号：11位，智能默认3,4
	if user.Phone != "138****1234" {
		t.Errorf("Phone = %s, want 138****1234", user.Phone)
	}

	// 身份证：18位，智能默认3,4
	if user.IDCard != "110***********1234" {
		t.Errorf("IDCard = %s, want 110***********1234", user.IDCard)
	}

	// 自定义参数
	if user.Custom != "custo*******23" {
		t.Errorf("Custom = %s, want custo*******23", user.Custom)
	}

	t.Logf("脱敏结果: %+v", user)
}

type TestMaskStuct struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number" mask:",5,4"`
	TestEmbedMaskStruct
}

type TestEmbedMaskStruct struct {
	Name1        string `json:"name1"`
	PhoneNumber1 string `json:"phone_number1" mask:",5,4"`
}

func TestMusk(t *testing.T) {
	tm := &TestMaskStuct{
		PhoneNumber: "1234567890123",
		Name:        "test",
		TestEmbedMaskStruct: TestEmbedMaskStruct{
			PhoneNumber1: "1234567890123",
			Name1:        "test",
		},
	}

	if err := desense.MaskStruct(tm); err != nil {
		t.Fatal(err)
	}
	t.Log(tm.PhoneNumber1)
}

var set = &TestMaskStuctSet{
	Items: []*TestMaskStuct{
		{
			PhoneNumber: "1234567890123",
			Name:        "test",
		},
	},
}

type TestMaskStuctSet struct {
	T     TestMaskStuct    `json:"t"`
	Items []*TestMaskStuct `json:"items"`
}

func TestMuskSet(t *testing.T) {
	if err := desense.MaskStruct(set); err != nil {
		t.Fatal(err)
	}

	t.Log(pretty.ToJSON(set))
}
