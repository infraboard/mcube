package main

import (
	"fmt"

	"github.com/infraboard/mcube/v2/desense"
)

// 演示智能默认脱敏的强大功能
func main() {
	fmt.Println("=== 智能默认脱敏演示 ===")

	// 示例1: 不带参数，智能识别
	type SmartUser struct {
		Phone    string `mask:"default"` // 11位自动识别为手机号
		IDCard   string `mask:"default"` // 18位自动识别为身份证
		BankCard string `mask:"default"` // 19位自动识别为银行卡
		Email    string `mask:"default"` // 根据长度智能处理
		Short    string `mask:"default"` // 短字符串智能处理
	}

	user1 := &SmartUser{
		Phone:    "13812341234",
		IDCard:   "110101199001011234",
		BankCard: "6222021234567890123",
		Email:    "test@example.com",
		Short:    "abc",
	}

	fmt.Println("原始数据:")
	fmt.Printf("  Phone:    %s\n", user1.Phone)
	fmt.Printf("  IDCard:   %s\n", user1.IDCard)
	fmt.Printf("  BankCard: %s\n", user1.BankCard)
	fmt.Printf("  Email:    %s\n", user1.Email)
	fmt.Printf("  Short:    %s\n", user1.Short)

	desense.MaskStruct(user1)

	fmt.Println("\n脱敏后 (mask:\"default\" 智能识别):")
	fmt.Printf("  Phone:    %s  ✓ 自动保留前3后4\n", user1.Phone)
	fmt.Printf("  IDCard:   %s  ✓ 自动保留前3后4\n", user1.IDCard)
	fmt.Printf("  BankCard: %s  ✓ 自动保留前4后4\n", user1.BankCard)
	fmt.Printf("  Email:    %s  ✓ 根据长度智能处理\n", user1.Email)
	fmt.Printf("  Short:    %s  ✓ 短字符串智能处理\n", user1.Short)

	// 示例2: 混合使用 - 专用策略与默认策略
	fmt.Println("\n\n=== 混合使用演示 ===")

	type MixedUser struct {
		Phone1 string `mask:"phone"`       // 专用策略
		Phone2 string `mask:"default"`     // 智能默认
		Phone3 string `mask:"default,5,2"` // 自定义参数
		Email1 string `mask:"email"`       // 专用策略
		Email2 string `mask:"default"`     // 智能默认
	}

	user2 := &MixedUser{
		Phone1: "13812341234",
		Phone2: "13812341234",
		Phone3: "13812341234",
		Email1: "test@example.com",
		Email2: "test@example.com",
	}

	fmt.Println("原始数据:")
	fmt.Printf("  Phone1: %s\n", user2.Phone1)
	fmt.Printf("  Phone2: %s\n", user2.Phone2)
	fmt.Printf("  Phone3: %s\n", user2.Phone3)
	fmt.Printf("  Email1: %s\n", user2.Email1)
	fmt.Printf("  Email2: %s\n", user2.Email2)

	desense.MaskStruct(user2)

	fmt.Println("\n脱敏后:")
	fmt.Printf("  Phone1: %s  (mask:\"phone\" 专用策略)\n", user2.Phone1)
	fmt.Printf("  Phone2: %s  (mask:\"default\" 智能识别)\n", user2.Phone2)
	fmt.Printf("  Phone3: %s  (mask:\"default,5,2\" 自定义)\n", user2.Phone3)
	fmt.Printf("  Email1: %s  (mask:\"email\" 专用策略)\n", user2.Email1)
	fmt.Printf("  Email2: %s  (mask:\"default\" 智能识别)\n", user2.Email2)

	// 示例3: 向后兼容
	fmt.Println("\n\n=== 向后兼容演示 ===")

	type LegacyUser struct {
		OldStyle string `mask:"default,3,4"` // 旧写法仍然支持
		NewStyle string `mask:"default"`     // 新写法
	}

	user3 := &LegacyUser{
		OldStyle: "13812341234",
		NewStyle: "13812341234",
	}

	fmt.Println("原始数据:")
	fmt.Printf("  OldStyle: %s\n", user3.OldStyle)
	fmt.Printf("  NewStyle: %s\n", user3.NewStyle)

	desense.MaskStruct(user3)

	fmt.Println("\n脱敏后:")
	fmt.Printf("  OldStyle: %s  (旧写法 mask:\"default,3,4\" 继续支持)\n", user3.OldStyle)
	fmt.Printf("  NewStyle: %s  (新写法 mask:\"default\" 智能识别)\n", user3.NewStyle)
	fmt.Println("\n✓ 完全向后兼容，无需修改现有代码！")

	// 示例4: 直接使用函数
	fmt.Println("\n\n=== 直接函数调用 ===")

	phone := "13812341234"
	fmt.Printf("原始手机号: %s\n", phone)
	fmt.Printf("使用 MaskPhone(): %s\n", desense.MaskPhone(phone))
	fmt.Printf("使用 MaskString(default): %s\n", desense.MaskString(phone, "default"))
	fmt.Printf("使用 Default().DeSense(): %s\n", desense.Default().DeSense(phone))
}
