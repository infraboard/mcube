# 数据脱敏 (Desense)

`desense` 包提供了灵活、易用的数据脱敏功能，支持多种内置脱敏策略和自定义扩展。

## 特性

✨ **多种内置策略**：支持手机号、邮箱、身份证、银行卡、姓名、地址、密码等常见场景  
🔧 **灵活扩展**：支持自定义脱敏策略  
🏷️ **Struct Tag支持**：通过tag标签自动脱敏结构体字段  
🚀 **便捷API**：提供简单易用的函数接口  
🔒 **完全向后兼容**：保留所有旧版API

## 快速开始

### 1. 字符串脱敏

```go
package main

import (
    "fmt"
    "github.com/infraboard/mcube/v2/desense"
)

func main() {
    // 手机号脱敏
    phone := desense.MaskPhone("13812341234")
    fmt.Println(phone) // 输出: 138****1234

    // 邮箱脱敏
    email := desense.MaskEmail("test@example.com")
    fmt.Println(email) // 输出: te**@example.com

    // 身份证脱敏
    idcard := desense.MaskIDCard("110101199001011234")
    fmt.Println(idcard) // 输出: 110***********1234

    // 银行卡脱敏
    bankcard := desense.MaskBankCard("6222021234567890123")
    fmt.Println(bankcard) // 输出: 6222 **** **** 0123

    // 姓名脱敏
    name := desense.MaskName("张三")
    fmt.Println(name) // 输出: 张*

    // 地址脱敏
    address := desense.MaskAddress("北京市朝阳区某某街道123号")
    fmt.Println(address) // 输出: 北京市朝阳区****

    // 密码脱敏
    password := desense.MaskPassword("password123")
    fmt.Println(password) // 输出: ******
}
```

### 2. 结构体脱敏

使用struct tag标签自动脱敏：

```go
type User struct {
    Name     string `json:"name" mask:"name"`
    Phone    string `json:"phone" mask:"phone"`
    Email    string `json:"email" mask:"email"`
    IDCard   string `json:"idcard" mask:"idcard"`
    BankCard string `json:"bankcard" mask:"bankcard"`
    Address  string `json:"address" mask:"address"`
    Password string `json:"password" mask:"password"`
}

func main() {
    user := &User{
        Name:     "张三",
        Phone:    "13812341234",
        Email:    "test@example.com",
        IDCard:   "110101199001011234",
        BankCard: "6222021234567890123",
        Address:  "北京市朝阳区某某街道123号",
        Password: "password123",
    }

    // 自动脱敏所有标记的字段
    if err := desense.MaskStruct(user); err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", user)
    // 输出: {Name:张* Phone:138****1234 Email:te**@example.com ...}
}
```

### 3. 嵌套结构体和切片

支持嵌套结构体和切片的递归脱敏：

```go
type Company struct {
    Name  string `json:"name"`
    Users []*User `json:"users"` // User切片
}

type UserWithCompany struct {
    User    *User    `json:"user"`    // 嵌套结构体
    Company *Company `json:"company"` // 嵌套结构体
}

func main() {
    data := &UserWithCompany{
        User: &User{Phone: "13812341234"},
        Company: &Company{
            Users: []*User{
                {Phone: "13912345678"},
                {Phone: "13612349876"},
            },
        },
    }

    desense.MaskStruct(data) // 自动递归脱敏所有嵌套对象
}
```

### 4. 智能默认脱敏（推荐）

**新特性**：`mask:"default"` 不带参数时会根据字符串长度智能选择合适的脱敏参数！

```go
type AutoUser struct {
    Phone    string `mask:"default"` // 自动识别11位作为手机号，保留前3后4
    IDCard   string `mask:"default"` // 自动识别18位作为身份证，保留前3后4
    ShortStr string `mask:"default"` // 短字符串自动选择合适的保留长度
    Custom   string `mask:"default,5,2"` // 也可以手动指定参数
}

func main() {
    user := &AutoUser{
        Phone:    "13812341234",      // 11位 → 138****1234 (自动3,4)
        IDCard:   "110101199001011234", // 18位 → 110***********1234 (自动3,4)
        ShortStr: "abc",              // 3位  → a** (自动1,0)
        Custom:   "customvalue",      // 自定义 → custo****lue
    }
    
    desense.MaskStruct(user)
}
```

**智能规则**：
- **≤4位**: 保留前1位 (如: `a**`)
- **5-6位**: 保留前1后1 (如: `a***b`)
- **7-10位**: 保留前2后2 (如: `ab****gh`)
- **11位** (手机号): 保留前3后4 (如: `138****1234`)
- **12-18位**: 保留前3后4 (如: `110***********1234`)
- **≥19位** (银行卡等): 保留前4后4 (如: `6222***********0123`)

## 内置脱敏策略

| 策略名称 | Tag值 | 便捷函数 | 效果示例 |
|---------|------|---------|---------|
| 手机号 | `mask:"phone"` | `MaskPhone()` | 138****1234 |
| 邮箱 | `mask:"email"` | `MaskEmail()` | te**@example.com |
| 身份证 | `mask:"idcard"` | `MaskIDCard()` | 110***********1234 |
| 银行卡 | `mask:"bankcard"` | `MaskBankCard()` | 6222 \*\*\*\* \*\*\*\* 0123 |
| 姓名 | `mask:"name"` | `MaskName()` | 张* |
| 地址 | `mask:"address"` | `MaskAddress()` | 北京市朝阳区\*\*\*\* |
| 密码 | `mask:"password"` | `MaskPassword()` | ****** |
| 默认 | `mask:"default,3,4"` | `Default().DeSense()` | abc*******kl |

## 自定义脱敏策略

### 方式1: 实现Desenser接口

```go
// 自定义IP地址脱敏器
type ipDesenser struct{}

func (i *ipDesenser) DeSense(value string, args ...string) string {
    parts := strings.Split(value, ".")
    if len(parts) != 4 {
        return value
    }
    return parts[0] + ".*.*." + parts[3]
}

// 注册自定义策略
func init() {
    desense.Registry("ip", &ipDesenser{})
}

// 使用
func main() {
    // 方式1: 直接调用
    masked := desense.Get("ip").DeSense("192.168.1.1")
    fmt.Println(masked) // 输出: 192.*.*.1

    // 方式2: 使用MaskString
    masked = desense.MaskString("192.168.1.1", "ip")
    
    // 方式3: 在结构体中使用
    type Server struct {
        IP string `mask:"ip"`
    }
}
```

### 方式2: 使用默认策略+参数

```go
type CustomData struct {
    // 保留前5位，后4位
    Data string `mask:"default,5,4"`
}
```

## Tag格式说明

### 基本格式
```
mask:"策略名称,参数1,参数2,..."
```

### 示例
```go
type User struct {
    // 使用phone策略，无需参数
    Phone string `mask:"phone"`
    
    // 使用default策略，不带参数（智能识别，推荐）
    AutoField string `mask:"default"`
    
    // 使用default策略，保留前3位和后2位
    Custom string `mask:"default,3,2"`
    
    // 使用email策略
    Email string `mask:"email"`
    
    // 留空策略名使用default（智能识别）
    Field1 string `mask:""`          // 等同于 mask:"default"
    
    // 留空策略名但指定参数
    Field2 string `mask:",3,4"`      // 等同于 mask:"default,3,4"
}
```

### Default策略说明

`default` 策略支持三种使用方式：

1. **智能模式（推荐）**：`mask:"default"` 或 `mask:""`
   - 根据字符串长度自动选择最合适的脱敏参数
   - 适合大多数场景，无需手动配置

2. **自定义模式**：`mask:"default,3,4"`
   - 手动指定保留前3位、后4位
   - 适合有特殊需求的场景

3. **部分指定**：`mask:",3,4"` 
   - 留空策略名，默认使用default
   - 指定自定义参数

## 向后兼容

所有旧版API完全保留，现有代码无需修改：

```go
// 旧版写法 - 继续支持
type OldUser struct {
    Phone string `mask:"default,3,4"` // 手动指定参数
}

// 新版写法1 - 推荐：使用专用策略
type NewUser1 struct {
    Phone string `mask:"phone"` // 更语义化
}

// 新版写法2 - 推荐：智能默认
type NewUser2 struct {
    Phone string `mask:"default"` // 自动识别，无需配置参数
}
```

**升级建议**：
- ✅ `mask:"default,3,4"` → `mask:"phone"` (使用专用策略)
- ✅ `mask:"default,3,4"` → `mask:"default"` (智能识别)  
- ✅ 混合使用也完全没问题

## API参考

### 便捷函数

```go
// 使用指定策略脱敏字符串
func MaskString(value, strategy string, args ...string) string

// 手机号脱敏
func MaskPhone(phone string) string

// 邮箱脱敏
func MaskEmail(email string) string

// 身份证脱敏
func MaskIDCard(idcard string) string

// 银行卡脱敏
func MaskBankCard(card string) string

// 姓名脱敏
func MaskName(name string) string

// 地址脱敏
func MaskAddress(address string) string

// 密码脱敏
func MaskPassword(password string) string

// 结构体脱敏
func MaskStruct(target any) error
```

### 核心接口

```go
// Desenser 脱敏器接口
type Desenser interface {
    DeSense(value string, args ...string) string
}

// 注册自定义脱敏器
func Registry(name string, d Desenser)

// 获取指定的脱敏器
func Get(name string) Desenser

// 获取默认脱敏器
func Default() Desenser
```

## 性能

运行基准测试：

```bash
cd desense
go test -bench=. -benchmem
```

典型性能（参考）：
- 单次字符串脱敏: ~100-500 ns/op
- 结构体脱敏: ~1-5 μs/op

## 常见问题

### Q: 如何处理不需要脱敏的字段？
A: 不添加`mask` tag即可，或者设置为空字符串。

### Q: 脱敏会修改原始数据吗？
A: 会。`MaskStruct()`会直接修改传入的结构体。如果需要保留原始数据，请先深拷贝。

### Q: `mask:"default"` 和 `mask:"phone"` 有什么区别？
A: 
- `mask:"phone"` 是专门针对手机号设计的策略，总是保留前3后4
- `mask:"default"` 会根据字符串长度智能选择，11位时效果与phone相同
- 推荐使用语义更明确的专用策略如 `phone`、`email` 等

### Q: 如何自定义脱敏规则？
A: 使用 `mask:"default,前缀长度,后缀长度"`，例如：
```go
Field string `mask:"default,5,2"` // 保留前5位和后2位
```

### Q: 支持JSON序列化时自动脱敏吗？
A: 建议在序列化前调用`MaskStruct()`进行脱敏，或者使用自定义的JSON Marshaler。

### Q: 如何临时禁用脱敏？
A: 可以通过条件判断决定是否调用`MaskStruct()`。

## 贡献

欢迎提交Issue和Pull Request！

## 许可证

与mcube项目保持一致。
