# IOC - 依赖注入容器

[![Go Report Card](https://goreportcard.com/badge/github.com/infraboard/mcube)](https://goreportcard.com/report/github.com/infraboard/mcube)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

一个轻量级、高性能的Go语言依赖注入（IoC）容器，支持命名空间隔离、自动装配、生命周期管理等企业级特性。

---

## 📋 目录

- [简介](#简介)
- [核心特性](#核心特性)
- [快速开始](#快速开始)
- [核心概念](#核心概念)
- [命名空间详解](#命名空间详解)
- [对象注册](#对象注册)
- [依赖注入](#依赖注入)
- [生命周期管理](#生命周期管理)
- [配置管理](#配置管理)
- [高级特性](#高级特性)
- [最佳实践](#最佳实践)
- [API参考](#api参考)
- [常见问题](#常见问题)

---

## 简介

### 什么是IOC？

IOC（Inversion of Control，控制反转）是一种设计模式，它将对象的创建和依赖关系的管理从应用代码中分离出来，交给容器来统一管理。通过IOC容器，你可以：

- **解耦代码**：不需要在代码中硬编码依赖关系
- **集中管理**：统一管理对象的创建、初始化和销毁
- **提高可测试性**：方便进行单元测试和依赖替换

### 为什么需要mcube/ioc？

在构建Go微服务应用时，我们通常需要：

1. **管理复杂依赖**：数据库、缓存、消息队列等基础组件的初始化顺序
2. **配置统一加载**：从配置文件或环境变量加载配置
3. **优雅启停**：按正确的顺序启动和关闭组件
4. **模块化设计**：将应用拆分为独立的模块，便于维护

mcube/ioc 正是为解决这些问题而设计的轻量级依赖注入容器。

---

> **✅ v2.0 已修复死锁问题**
> 
> 在早期版本中，如果在 `Priority()` 方法中调用 `ioc.Get()` 会导致死锁。**v2.0 版本已完全修复此问题**！
> 
> 修复方案：
> - Priority() 在加锁**之前**调用并缓存结果
> - 排序使用缓存值，避免重复调用
> - 完全消除了 Registry-Get 死锁风险
> 
> ```go
> // ✅ v2.0 中这样写也是安全的（但不推荐）
> func (s *Service) Priority() int {
>     other := ioc.Default().Get("other")  // v2.0 已修复，不会死锁
>     return other.Priority() - 1
> }
> 
> // ✅ 推荐写法：简单清晰
> func (s *Service) Priority() int {
>     return 10  // 返回固定值，性能更好
> }
> ```
> 
> **最佳实践建议**：虽然死锁已修复，但仍建议 Priority() 返回常量值，原因：
> - 性能更好（避免重复计算）
> - 逻辑更清晰（优先级应在设计时确定）
> - 避免复杂依赖关系

---

## 核心特性

### ✨ 命名空间隔离

内置4个命名空间，按初始化优先级排序：

| 命名空间 | 优先级 | 用途 | 典型应用 |
|---------|-------|------|---------|
| **configs** | 99 | 配置对象 | 数据库配置、应用配置 |
| **default** | 9 | 工具类 | 日志、缓存、数据库连接 |
| **controllers** | 0 | 业务控制器 | Service层业务逻辑 |
| **apis** | -99 | API处理器 | HTTP/gRPC接口实现 |

### 🔄 自动依赖注入

通过结构体标签（struct tag）声明依赖，容器自动装配：

```go
type UserAPI struct {
    UserService service.UserService `ioc:"autowire=true;namespace=controllers"`
}
```

### ⚡ 生命周期管理

- **Init()**: 对象初始化
- **Priority()**: 启动顺序控制
- **Close(ctx)**: 优雅关闭（倒序）
- **生命周期钩子**: PostConfig、PreInit、PostInit、PreStop、PostStop

### 📦 多配置源支持

- 支持 **TOML**、**YAML**、**JSON** 配置文件
- 支持从**环境变量**加载配置
- 配置优先级：环境变量 > 配置文件

### 🔍 依赖可视化

自动分析依赖关系，生成依赖图，便于理解系统架构。

### 📌 版本控制

支持对象版本管理（语义化版本），同一对象可注册多个版本。

---

## 快速开始

### 安装

```bash
go get github.com/infraboard/mcube/v2
```

### 第一个示例

创建一个简单的Web服务，展示如何使用IOC容器：

**步骤1：定义业务控制器**

```go
package impl

import (
    "github.com/infraboard/mcube/v2/ioc"
)

func init() {
    // 注册到controllers命名空间
    ioc.Controller().Registry(&HelloService{})
}

type HelloService struct {
    ioc.ObjectImpl  // 继承默认实现
}

func (s *HelloService) Init() error {
    // 初始化逻辑
    return nil
}

func (s *HelloService) Hello(name string) string {
    return "Hello, " + name + "!"
}
```

**步骤2：定义API处理器**

```go
package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/infraboard/mcube/v2/ioc"
    ioc_gin "github.com/infraboard/mcube/v2/ioc/config/gin"
)

func init() {
    // 注册到apis命名空间
    ioc.Api().Registry(&HelloAPI{})
}

type HelloAPI struct {
    ioc.ObjectImpl
    // 自动注入依赖
    svc *HelloService `ioc:"autowire=true;namespace=controllers"`
}

func (h *HelloAPI) Init() error {
    // 注册路由
    router := ioc_gin.ObjectRouter(h)
    router.GET("/hello", h.HandleHello)
    return nil
}

func (h *HelloAPI) HandleHello(c *gin.Context) {
    name := c.Query("name")
    c.JSON(http.StatusOK, gin.H{
        "message": h.svc.Hello(name),
    })
}
```

**步骤3：启动应用**

```go
package main

import (
    "context"
    "github.com/infraboard/mcube/v2/ioc/server"
    
    // 导入模块，触发init()注册
    _ "your-project/impl"
    _ "your-project/api"
)

func main() {
    // 启动应用，IOC容器自动管理生命周期
    err := server.Run(context.Background())
    if err != nil {
        panic(err)
    }
}
```

**运行效果**：

```bash
$ go run main.go
# 访问 http://localhost:8080/hello?name=World
# 返回: {"message": "Hello, World!"}
```

---

## 核心概念

### 1️⃣ 命名空间（Namespace）

命名空间是IOC容器中对象的逻辑分组，不同命名空间有不同的初始化优先级：

```go
ioc.Config()      // 配置命名空间，最先初始化（优先级99）
ioc.Default()     // 默认命名空间，工具类（优先级9）
ioc.Controller()  // 控制器命名空间（优先级0）
ioc.Api()         // API命名空间，最后初始化（优先级-99）
```

**初始化顺序**：configs → default → controllers → apis  
**关闭顺序**：apis → controllers → default → configs（倒序）

### 2️⃣ 对象（Object）

对象是注册到IOC容器的组件，必须实现 `Object` 接口：

```go
type Object interface {
    Init() error                  // 初始化
    Name() string                 // 对象名称
    Version() string              // 版本号（默认1.0.0）
    Priority() int                // 优先级（同命名空间内）
    Close(ctx context.Context)    // 优雅关闭
    Meta() ObjectMeta             // 元数据
}
```

**最简实现**：继承 `ioc.ObjectImpl` 获得默认实现，只需覆写需要的方法。

### 3️⃣ 依赖注入（Autowire）

通过结构体标签声明依赖，容器自动装配：

```go
type UserAPI struct {
    // 自动注入控制器
    UserSvc UserService `ioc:"autowire=true;namespace=controllers"`
    
    // 注入指定版本
    Cache CacheService `ioc:"autowire=true;namespace=default;version=2.0.0"`
}
```

**标签说明**：
- `autowire=true`: 启用自动注入
- `namespace=xxx`: 指定从哪个命名空间获取
- `version=x.x.x`: 指定对象版本（可选）

### 4️⃣ 生命周期（Lifecycle）

完整的对象生命周期：

```
注册 → 配置加载 → 依赖注入 → 初始化 → 运行 → 关闭
 ↓         ↓         ↓         ↓       ↓      ↓
Registry  Load   Autowire    Init   Running  Close
```

**生命周期钩子**（可选实现）：

- `OnPostConfig()`: 配置加载后
- `OnPreInit()`: 初始化前
- `OnPostInit()`: 初始化后
- `OnPreStop(ctx)`: 关闭前
- `OnPostStop(ctx)`: 关闭后

---

## 命名空间详解

### 📌 Config命名空间（优先级99）

**用途**：存放各种配置对象，最先初始化

**典型应用**：
- 数据库配置
- Redis配置
- 应用全局配置

**示例**：

```go
import "github.com/infraboard/mcube/v2/ioc"

func init() {
    ioc.Config().Registry(&DatabaseConfig{})
}

type DatabaseConfig struct {
    ioc.ObjectImpl
    Host     string `toml:"host" env:"DB_HOST"`
    Port     int    `toml:"port" env:"DB_PORT"`
    Database string `toml:"database" env:"DB_NAME"`
}

func (c *DatabaseConfig) Init() error {
    // 配置验证
    if c.Host == "" {
        return fmt.Errorf("database host is required")
    }
    return nil
}
```

### 🔧 Default命名空间（优先级9）

**用途**：存放工具类和基础组件

**典型应用**：
- 数据库连接（GORM、MongoDB）
- 缓存客户端（Redis）
- 日志组件
- 消息队列客户端

**示例**：

```go
import "github.com/infraboard/mcube/v2/ioc"

func init() {
    ioc.Default().Registry(&RedisClient{})
}

type RedisClient struct {
    ioc.ObjectImpl
    Config *RedisConfig `ioc:"autowire=true;namespace=configs"`
    client *redis.Client
}

func (r *RedisClient) Init() error {
    r.client = redis.NewClient(&redis.Options{
        Addr: r.Config.Addr,
    })
    return r.client.Ping(context.Background()).Err()
}

func (r *RedisClient) Close(ctx context.Context) {
    if r.client != nil {
        r.client.Close()
    }
}
```

### 🎮 Controller命名空间（优先级0）

**用途**：存放业务逻辑控制器（Service层）

**典型应用**：
- 业务服务实现
- 领域模型
- 业务逻辑处理

**示例**：

```go
import "github.com/infraboard/mcube/v2/ioc"

func init() {
    ioc.Controller().Registry(&UserService{})
}

type UserService struct {
    ioc.ObjectImpl
    DB    *gorm.DB      `ioc:"autowire=true;namespace=default"`
    Cache *RedisClient `ioc:"autowire=true;namespace=default"`
}

func (s *UserService) Init() error {
    // 初始化业务逻辑
    return nil
}

func (s *UserService) GetUser(id string) (*User, error) {
    // 业务实现
    var user User
    err := s.DB.First(&user, "id = ?", id).Error
    return &user, err
}
```

### 🌐 Api命名空间（优先级-99）

**用途**：存放API处理器（HTTP/gRPC接口层），最后初始化

**典型应用**：
- HTTP接口处理器
- gRPC服务实现
- WebSocket处理器

**示例**：

```go
import (
    "github.com/infraboard/mcube/v2/ioc"
    ioc_gin "github.com/infraboard/mcube/v2/ioc/config/gin"
)

func init() {
    ioc.Api().Registry(&UserAPI{})
}

type UserAPI struct {
    ioc.ObjectImpl
    UserService *UserService `ioc:"autowire=true;namespace=controllers"`
}

func (h *UserAPI) Name() string {
    return "user"  // API路径前缀: /api/v1/user
}

func (h *UserAPI) Init() error {
    router := ioc_gin.ObjectRouter(h)
    router.GET("/:id", h.GetUser)
    return nil
}

func (h *UserAPI) GetUser(c *gin.Context) {
    id := c.Param("id")
    user, err := h.UserService.GetUser(id)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, user)
}
```

### 🔄 初始化流程

```
1. [configs]    加载所有配置对象
      ↓
2. [default]    初始化工具类（DB、Redis等）
      ↓
3. [controllers] 初始化业务控制器
      ↓
4. [apis]       注册API路由
      ↓
5. 应用启动
```

---

## 对象注册

### 基本注册

```go
import "github.com/infraboard/mcube/v2/ioc"

func init() {
    // 注册到指定命名空间
    ioc.Controller().Registry(&MyService{})
}

type MyService struct {
    ioc.ObjectImpl
}
```

### 实现Object接口

有两种方式实现Object接口：

**方式1：继承ObjectImpl（推荐）**

```go
type MyService struct {
    ioc.ObjectImpl  // 获得默认实现
}

// 只需覆写需要的方法
func (s *MyService) Init() error {
    // 自定义初始化逻辑
    return nil
}

func (s *MyService) Name() string {
    return "my-service"  // 自定义名称
}
```

**方式2：完全自定义实现**

```go
type MyService struct {
    // 不继承ObjectImpl
}

// 必须实现所有接口方法
func (s *MyService) Init() error { return nil }
func (s *MyService) Name() string { return "my-service" }
func (s *MyService) Version() string { return "1.0.0" }
func (s *MyService) Priority() int { return 0 }
func (s *MyService) Close(ctx context.Context) {}
func (s *MyService) Meta() ioc.ObjectMeta { return ioc.DefaultObjectMeta() }
```

### 对象命名规则

对象名称用于标识和获取对象：

```go
// 1. 默认名称：包名.类型名
type UserService struct {
    ioc.ObjectImpl
}
// 名称: *impl.UserService

// 2. 自定义名称（推荐）
func (s *UserService) Name() string {
    return "user-service"  // 名称: user-service
}
```

### 对象版本控制

支持同一对象的多个版本：

```go
type CacheV1 struct {
    ioc.ObjectImpl
}

func (c *CacheV1) Version() string {
    return "1.0.0"
}

type CacheV2 struct {
    ioc.ObjectImpl
}

func (c *CacheV2) Version() string {
    return "2.0.0"
}

func init() {
    ioc.Default().Registry(&CacheV1{})
    ioc.Default().Registry(&CacheV2{})  // 同名不同版本
}
```

### 优先级控制

控制同一命名空间内的初始化顺序：

```go
type DatabaseService struct {
    ioc.ObjectImpl
}

func (d *DatabaseService) Priority() int {
    return 100  // 数字越大越先初始化
}

type CacheService struct {
    ioc.ObjectImpl
}

func (c *CacheService) Priority() int {
    return 50  // 在Database之后初始化
}
```

**✅ v2.0 死锁问题已修复**

在早期版本中，Priority() 中访问容器会导致死锁。**v2.0 版本已完全修复此问题**！

```go
// ✅ 推荐：返回常量（最简单）
func (s *MyService) Priority() int {
    return -1  // 直接返回固定值
}

// ✅ 支持：相对优先级（v2.0+ 完全安全）
func (s *CacheService) Priority() int {
    // 优先级是相对概念，可以基于依赖关系确定
    db := ioc.Default().Get("database")
    if db != nil {
        return db.Priority() - 10  // 在数据库之后初始化
    }
    return 0
}

// ✅ 支持：条件优先级
func (s *MonitorService) Priority() int {
    if os.Getenv("ENV") == "production" {
        return 100  // 生产环境高优先级
    }
    return 10   // 开发环境低优先级
}
```

**修复说明**：
1. Priority() 现在在 Registry() 加锁**之前**调用
2. Priority 值被缓存到 ObjectWrapper 中
3. 排序使用缓存值，不再重复调用 Priority()
4. 彻底消除了写锁-读锁的冲突

**使用建议**：
- **常量优先级**：适合大多数场景，简单直观
- **相对优先级**：当需要基于依赖关系动态确定顺序时使用
- **注意循环依赖**：确保优先级计算不会形成循环

**性能说明**：Priority() 只在注册时调用一次，结果会被缓存，不影响运行时性能。

**调试工具**：
```bash
# 开启调试模式查看对象注册过程
export IOC_DEBUG=true
go run main.go
```

### 批量注册

```go
func init() {
    ioc.Controller().RegistryAll(
        &UserService{},
        &OrderService{},
        &ProductService{},
    )
}
```

### 对象元数据

为对象添加额外信息：

```go
type MyAPI struct {
    ioc.ObjectImpl
}

func (a *MyAPI) Meta() ioc.ObjectMeta {
    return ioc.ObjectMeta{
        CustomPathPrefix: "/custom/path",  // 自定义API路径前缀
        Extra: map[string]string{
            "description": "My API Handler",
            "author":      "your-name",
        },
    }
}
```

---

## 依赖注入

### 使用标签自动注入（推荐）

通过结构体标签声明依赖，容器自动装配：

```go
type UserAPI struct {
    ioc.ObjectImpl
    
    // 基本注入：从controllers命名空间注入
    UserSvc *UserService `ioc:"autowire=true;namespace=controllers"`
    
    // 指定版本注入
    Cache *CacheService `ioc:"autowire=true;namespace=default;version=2.0.0"`
    
    // 私有字段不会被注入
    logger *log.Logger
}
```

**标签参数说明**：

| 参数 | 说明 | 必填 | 示例 |
|------|------|------|------|
| `autowire` | 是否启用自动注入 | 是 | `autowire=true` |
| `namespace` | 从哪个命名空间获取 | 是 | `namespace=controllers` |
| `version` | 对象版本 | 否 | `version=1.0.0` |

### 手动获取依赖

在Init()方法中手动获取：

```go
type UserAPI struct {
    ioc.ObjectImpl
    userSvc *UserService
}

func (a *UserAPI) Init() error {
    // 方式1：直接Get
    obj := ioc.Controller().Get("user-service")
    a.userSvc = obj.(*UserService)
    
    // 方式2：使用Load（推荐）
    var svc *UserService
    err := ioc.Controller().Load(&svc)
    if err != nil {
        return err
    }
    a.userSvc = svc
    
    return nil
}
```

### 指定版本获取

```go
import "github.com/infraboard/mcube/v2/ioc"

// 获取指定版本
obj := ioc.Default().Get("cache-service", ioc.WithVersion("2.0.0"))
cache := obj.(*CacheService)
```

### 依赖注入执行时机

```
1. Registry()      注册对象
      ↓
2. LoadConfig()    加载配置到对象
      ↓
3. Autowire()      自动注入依赖 ← 这里执行注入
      ↓
4. Init()          对象初始化
```

**重要**：在 `Init()` 方法中，所有标签声明的依赖已经注入完成。

### 循环依赖处理

IOC容器会检测循环依赖：

```go
// ❌ 错误：循环依赖
type ServiceA struct {
    B *ServiceB `ioc:"autowire=true;namespace=controllers"`
}

type ServiceB struct {
    A *ServiceA `ioc:"autowire=true;namespace=controllers"`
}
```

**解决方案**：

1. **延迟获取**：在使用时再获取

```go
type ServiceA struct {
    ioc.ObjectImpl
}

func (a *ServiceA) GetB() *ServiceB {
    obj := ioc.Controller().Get("service-b")
    return obj.(*ServiceB)
}
```

2. **引入中间层**：通过接口解耦

```go
type ServiceA struct {
    Handler BHandler `ioc:"autowire=true;namespace=controllers"`
}

type BHandler interface {
    Handle()
}
```

### 依赖可选性

当依赖不存在时的处理：

```go
type MyService struct {
    ioc.ObjectImpl
    // 必选依赖：不存在会报错
    DB *gorm.DB `ioc:"autowire=true;namespace=default"`
}

func (s *MyService) Init() error {
    // 可选依赖：手动检查
    obj := ioc.Default().Get("optional-cache")
    if obj != nil {
        s.cache = obj.(*Cache)
    }
    return nil
}
```

### 声明依赖关系

当手动获取依赖时，如需在依赖图中展示关系，可实现 `DependencyDeclarer` 接口：

```go
type MyService struct {
    ioc.ObjectImpl
    cache *Cache  // 手动获取
}

func (s *MyService) Init() error {
    obj := ioc.Default().Get("cache")
    s.cache = obj.(*Cache)
    return nil
}

// 声明依赖关系（用于依赖可视化）
func (s *MyService) Dependencies() []ioc.Dependency {
    return []ioc.Dependency{
        {
            Name:      "cache",
            Namespace: "default",
            Version:   "1.0.0",
        },
    }
}
```

---

## 生命周期管理

### 初始化顺序

IOC容器按以下顺序初始化对象：

```
1. 按命名空间优先级：configs(99) → default(9) → controllers(0) → apis(-99)
2. 同一命名空间内按Priority()：数字越大越先初始化
3. 同优先级按注册顺序
```

**示例**：

```go
// configs命名空间：优先级99，最先初始化
type AppConfig struct {
    ioc.ObjectImpl
}
func (c *AppConfig) Priority() int { return 100 }

// default命名空间：优先级9
type Database struct {
    ioc.ObjectImpl
}
func (d *Database) Priority() int { return 50 }

// controllers命名空间：优先级0
type UserService struct {
    ioc.ObjectImpl
}
func (s *UserService) Priority() int { return 0 }

// 初始化顺序：AppConfig → Database → UserService
```

### 完整生命周期

```
┌─────────────────────────────────────────────────────────┐
│  1. Registry(): 注册对象到容器                            │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│  2. LoadConfig(): 从配置文件/环境变量加载配置              │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│  3. OnPostConfig(): 配置加载后钩子（可选）                 │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│  4. Autowire(): 自动注入依赖                              │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│  5. OnPreInit(): 初始化前钩子（可选）                      │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│  6. Init(): 对象初始化                                    │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│  7. OnPostInit(): 初始化后钩子（可选）                     │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│  8. Running: 应用运行中                                   │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│  9. OnPreStop(): 关闭前钩子（可选）                       │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│  10. Close(): 对象关闭（倒序）                            │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│  11. OnPostStop(): 关闭后钩子（可选）                     │
└─────────────────────────────────────────────────────────┘
```

### 生命周期钩子

实现可选的生命周期钩子接口：

**1. PostConfigHook - 配置加载后**

```go
type MyService struct {
    ioc.ObjectImpl
    Host string `toml:"host" env:"HOST"`
}

func (s *MyService) OnPostConfig() error {
    // 配置验证
    if s.Host == "" {
        return fmt.Errorf("host is required")
    }
    // 配置预处理
    s.Host = strings.TrimSpace(s.Host)
    return nil
}
```

**2. PreInitHook - 初始化前**

```go
func (s *MyService) OnPreInit() error {
    // 准备工作
    log.Println("Preparing to initialize MyService...")
    return nil
}
```

**3. PostInitHook - 初始化后**

```go
func (s *MyService) OnPostInit() error {
    // 启动后台任务
    go s.backgroundTask()
    // 注册监听器
    s.registerListeners()
    return nil
}
```

**4. PreStopHook - 关闭前**

```go
func (s *MyService) OnPreStop(ctx context.Context) error {
    // 优雅停机检查
    log.Println("Preparing to stop MyService...")
    // 等待请求完成（带超时）
    return s.waitForRequests(ctx)
}
```

**5. PostStopHook - 关闭后**

```go
func (s *MyService) OnPostStop(ctx context.Context) error {
    // 最终清理
    log.Println("MyService stopped successfully")
    return nil
}
```

### 优雅关闭

容器按**倒序**关闭对象：

```go
type Database struct {
    ioc.ObjectImpl
    conn *sql.DB
}

func (d *Database) Close(ctx context.Context) {
    if d.conn != nil {
        // 等待连接关闭，但不超过context超时
        d.conn.Close()
    }
}
```

**关闭顺序**：apis → controllers → default → configs（与初始化相反）

### 错误处理

```go
func (s *MyService) Init() error {
    // 初始化失败会阻止应用启动
    if err := s.connect(); err != nil {
        return fmt.Errorf("failed to connect: %w", err)
    }
    return nil
}

func (s *MyService) OnPostInit() error {
    // PostInit失败会记录错误但不阻止启动
    if err := s.warmup(); err != nil {
        log.Printf("warmup failed: %v", err)
        return err
    }
    return nil
}

func (s *MyService) Close(ctx context.Context) {
    // Close不返回error，应该做好容错
    if s.conn != nil {
        _ = s.conn.Close()
    }
}
```

### 超时控制

```go
func (s *MyService) Close(ctx context.Context) {
    // 使用context控制关闭超时
    done := make(chan struct{})
    
    go func() {
        s.cleanup()
        close(done)
    }()
    
    select {
    case <-done:
        log.Println("Cleanup completed")
    case <-ctx.Done():
        log.Println("Cleanup timeout")
    }
}
```

---

## 配置管理

### 配置文件加载

支持 **TOML**、**YAML**、**JSON** 三种格式：

<parameter name="startLine">**配置文件示例 (etc/application.toml)**：

```toml
[database]
host = "localhost"
port = 3306
database = "myapp"
username = "root"
password = "secret"

[redis]
addr = "localhost:6379"
db = 0
```

**对象定义**：

```go
type DatabaseConfig struct {
    ioc.ObjectImpl
    Host     string `toml:"host" json:"host" yaml:"host"`
    Port     int    `toml:"port" json:"port" yaml:"port"`
    Database string `toml:"database" json:"database" yaml:"database"`
    Username string `toml:"username" json:"username" yaml:"username"`
    Password string `toml:"password" json:"password" yaml:"password"`
}

func (c *DatabaseConfig) Name() string {
    return "database"  // 对应配置文件中的[database]节点
}

func init() {
    ioc.Config().Registry(&DatabaseConfig{})
}
```

**启用配置文件加载**：

```go
import "github.com/infraboard/mcube/v2/ioc/server"

func main() {
    // 方式1：使用默认配置
    server.DefaultConfig.ConfigFile.Enabled = true
    server.DefaultConfig.ConfigFile.Paths = []string{
        "etc/application.toml",
        "etc/application.yaml",  // 支持多个配置文件
    }
    
    server.Run(context.Background())
}
```

### 环境变量配置

使用 `env` 标签声明环境变量映射：

```go
type DatabaseConfig struct {
    ioc.ObjectImpl
    // 优先从环境变量加载，不存在则使用配置文件
    Host     string `toml:"host" env:"DB_HOST"`
    Port     int    `toml:"port" env:"DB_PORT"`
    Database string `toml:"database" env:"DB_NAME"`
    Username string `toml:"username" env:"DB_USER"`
    Password string `toml:"password" env:"DB_PASS"`
}

func init() {
    ioc.Config().Registry(&DatabaseConfig{})
}
```

**使用**：

```bash
# 设置环境变量
export DB_HOST=production-db.example.com
export DB_PORT=5432
export DB_NAME=prod_database

# 启动应用，环境变量会覆盖配置文件
go run main.go
```

### 配置加载方式

**方式1：通过Server自动加载（推荐）**

```go
import "github.com/infraboard/mcube/v2/ioc/server"

func main() {
    server.DefaultConfig.ConfigFile.Enabled = true
    server.DefaultConfig.ConfigFile.Paths = []string{"etc/app.toml"}
    server.Run(context.Background())
}
```

**方式2：手动加载**

```go
import "github.com/infraboard/mcube/v2/ioc"

func main() {
    // 从文件加载
    content, _ := os.ReadFile("etc/app.toml")
    ioc.Config().LoadFromFileContent(content)
    
    // 从环境变量加载
    ioc.Config().LoadFromEnv("APP")  // 前缀为APP_的环境变量
    
    // 初始化所有对象
    ioc.InitAll()
}
```

### 配置优先级

配置加载遵循以下优先级（高到低）：

```
1. 环境变量（env标签）
   ↓
2. 配置文件（toml/yaml/json标签）
   ↓
3. 字段默认值
```

**示例**：

```go
type Config struct {
    ioc.ObjectImpl
    Host string `toml:"host" env:"HOST"`
    Port int    `toml:"port" env:"PORT"`
}

// 字段默认值
func NewConfig() *Config {
    return &Config{
        Host: "localhost",  // 默认值
        Port: 8080,         // 默认值
    }
}
```

**加载顺序**：
1. 使用默认值：`Host=localhost, Port=8080`
2. 加载配置文件（如果有）：覆盖默认值
3. 加载环境变量（如果有）：覆盖配置文件

### 配置验证

在 `OnPostConfig()` 钩子中验证配置：

```go
type DatabaseConfig struct {
    ioc.ObjectImpl
    Host     string `toml:"host" env:"DB_HOST"`
    Port     int    `toml:"port" env:"DB_PORT"`
    Database string `toml:"database" env:"DB_NAME"`
}

func (c *DatabaseConfig) OnPostConfig() error {
    // 必填项检查
    if c.Host == "" {
        return fmt.Errorf("database host is required")
    }
    
    // 范围检查
    if c.Port <= 0 || c.Port > 65535 {
        return fmt.Errorf("invalid port: %d", c.Port)
    }
    
    // 格式检查
    if c.Database == "" {
        return fmt.Errorf("database name is required")
    }
    
    return nil
}
```

### 多环境配置

**方式1：不同配置文件**

```go
// 根据环境变量选择配置文件
env := os.Getenv("APP_ENV")
if env == "" {
    env = "dev"
}

configFile := fmt.Sprintf("etc/application-%s.toml", env)
server.DefaultConfig.ConfigFile.Paths = []string{configFile}
```

**方式2：配置覆盖**

```toml
# etc/application.toml (基础配置)
[database]
host = "localhost"
port = 3306

# etc/application-prod.toml (生产环境覆盖)
[database]
host = "prod-db.example.com"
port = 5432
```

```go
server.DefaultConfig.ConfigFile.Paths = []string{
    "etc/application.toml",      // 基础配置
    "etc/application-prod.toml", // 环境特定配置（覆盖）
}
```

### 配置热更新

```go
type DatabaseConfig struct {
    ioc.ObjectImpl
    Host string `toml:"host" env:"DB_HOST"`
    mu   sync.RWMutex
}

func (c *DatabaseConfig) GetHost() string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.Host
}

func (c *DatabaseConfig) UpdateHost(host string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.Host = host
}
```

**注意**：IOC容器暂不支持自动配置热更新，需要手动实现。

### 敏感信息处理

```go
type DatabaseConfig struct {
    ioc.ObjectImpl
    Host     string `toml:"host" env:"DB_HOST"`
    Password string `toml:"password" env:"DB_PASS"`
}

func (c *DatabaseConfig) OnPostConfig() error {
    // 从密钥管理系统加载敏感信息
    if c.Password == "" {
        password, err := loadPasswordFromVault()
        if err != nil {
            return err
        }
        c.Password = password
    }
    return nil
}

func (c *DatabaseConfig) String() string {
    // 避免日志泄露密码
    return fmt.Sprintf("DB{host=%s, password=***}", c.Host)
}
```

---

## 高级特性

### 对象版本控制

支持语义化版本（Semantic Versioning）：

```go
type CacheV1 struct {
    ioc.ObjectImpl
}

func (c *CacheV1) Name() string {
    return "cache-service"
}

func (c *CacheV1) Version() string {
    return "1.0.0"
}

type CacheV2 struct {
    ioc.ObjectImpl
}

func (c *CacheV2) Name() string {
    return "cache-service"
}

func (c *CacheV2) Version() string {
    return "2.0.0"
}

func init() {
    ioc.Default().Registry(&CacheV1{})
    ioc.Default().Registry(&CacheV2{})
}
```

**使用指定版本**：

```go
// 注入时指定版本
type MyService struct {
    CacheV1 *CacheV1 `ioc:"autowire=true;namespace=default;version=1.0.0"`
    CacheV2 *CacheV2 `ioc:"autowire=true;namespace=default;version=2.0.0"`
}

// 获取时指定版本
cache := ioc.Default().Get("cache-service", ioc.WithVersion("2.0.0"))
```

### 对象覆盖与替换

允许覆盖已注册的对象（用于测试或自定义实现）：

```go
// 原始实现
type DefaultEmailService struct {
    ioc.ObjectImpl
}

func (s *DefaultEmailService) Send(to, subject, body string) error {
    // 真实邮件发送
    return smtp.SendEmail(to, subject, body)
}

// 测试环境替换为Mock实现
type MockEmailService struct {
    ioc.ObjectImpl
}

func (s *MockEmailService) Send(to, subject, body string) error {
    // Mock实现：只记录日志
    log.Printf("Mock email sent to %s", to)
    return nil
}

func init() {
    if os.Getenv("ENV") == "test" {
        // 覆盖原实现
        ioc.Default().Registry(&MockEmailService{})
    } else {
        ioc.Default().Registry(&DefaultEmailService{})
    }
}
```

### 依赖可视化

查看对象列表和依赖关系：

```go
import "github.com/infraboard/mcube/v2/ioc"

func main() {
    // 列出所有已注册对象
    objects := ioc.Default().List()
    for _, name := range objects {
        fmt.Println(name)
    }
    
    // 统计对象数量
    count := ioc.Controller().Len()
    fmt.Printf("Total controllers: %d\n", count)
    
    // 遍历对象
    ioc.Api().ForEach(func(obj *ioc.ObjectWrapper) {
        fmt.Printf("API: %s (v%s)\n", obj.Name, obj.Version)
    })
}
```

**生成依赖图**（参考 [DEPENDENCY_VISUALIZATION.md](DEPENDENCY_VISUALIZATION.md)）：

```bash
# 使用工具生成依赖关系图
go run tools/dependency-viz/main.go
```

### 条件注册

根据条件决定是否注册对象：

```go
func init() {
    // 只在特性开关启用时注册
    if featureEnabled("new-cache") {
        ioc.Default().Registry(&NewCacheService{})
    } else {
        ioc.Default().Registry(&OldCacheService{})
    }
    
    // 根据环境注册
    if os.Getenv("ENABLE_METRICS") == "true" {
        ioc.Default().Registry(&MetricsCollector{})
    }
}
```

### 自定义命名空间

创建自定义命名空间：

```go
import "github.com/infraboard/mcube/v2/ioc"

// 创建自定义命名空间
func TaskNamespace() ioc.StoreUser {
    return ioc.DefaultStore.Namespace("tasks").
        SetPriority(-50)  // 设置优先级
}

func init() {
    // 注册到自定义命名空间
    TaskNamespace().Registry(&CronTask{})
}
```

### 对象工厂模式

使用工厂函数创建对象：

```go
type ConnectionPool struct {
    ioc.ObjectImpl
    config *PoolConfig `ioc:"autowire=true;namespace=configs"`
    pool   *Pool
}

func (c *ConnectionPool) Init() error {
    // 根据配置创建连接池
    c.pool = NewPool(c.config.MaxConnections, c.config.Timeout)
    return c.pool.Connect()
}

func NewConnectionPool() *ConnectionPool {
    return &ConnectionPool{
        // 初始化默认值
    }
}

func init() {
    ioc.Default().Registry(NewConnectionPool())
}
```

### 延迟初始化

某些对象可能需要延迟初始化：

```go
type HeavyService struct {
    ioc.ObjectImpl
    client *heavy.Client
    once   sync.Once
}

func (s *HeavyService) Init() error {
    // 不在这里初始化重资源
    return nil
}

func (s *HeavyService) GetClient() *heavy.Client {
    // 首次使用时才初始化
    s.once.Do(func() {
        s.client = heavy.NewClient()
    })
    return s.client
}
```

### 对象分组

使用Meta元数据对对象分组：

```go
type UserAPI struct {
    ioc.ObjectImpl
}

func (a *UserAPI) Meta() ioc.ObjectMeta {
    meta := ioc.DefaultObjectMeta()
    meta.Extra["group"] = "user-module"
    meta.Extra["version"] = "v2"
    return meta
}

// 查询特定分组的对象
ioc.Api().ForEach(func(obj *ioc.ObjectWrapper) {
    if obj.Meta().Extra["group"] == "user-module" {
        fmt.Println(obj.Name)
    }
})
```

### 健康检查集成

为对象添加健康检查：

```go
type DatabaseService struct {
    ioc.ObjectImpl
    db *sql.DB
}

func (s *DatabaseService) HealthCheck() error {
    return s.db.Ping()
}

// 在健康检查端点中使用
func healthCheckHandler(c *gin.Context) {
    var unhealthy []string
    
    ioc.Default().ForEach(func(obj *ioc.ObjectWrapper) {
        if checker, ok := obj.Value.(interface{ HealthCheck() error }); ok {
            if err := checker.HealthCheck(); err != nil {
                unhealthy = append(unhealthy, obj.Name)
            }
        }
    })
    
    if len(unhealthy) > 0 {
        c.JSON(500, gin.H{"status": "unhealthy", "services": unhealthy})
    } else {
        c.JSON(200, gin.H{"status": "healthy"})
    }
}
```

### 事件监听

监听对象生命周期事件：

```go
type EventListener struct {
    ioc.ObjectImpl
}

func (l *EventListener) OnPostInit() error {
    log.Println("All objects initialized")
    // 发送初始化完成事件
    publishEvent("ioc.initialized")
    return nil
}

func init() {
    ioc.Default().Registry(&EventListener{})
}
```

---

## 最佳实践

### 项目结构建议

推荐的项目结构：

```
myapp/
├── main.go                 # 应用入口
├── etc/
│   ├── application.toml    # 配置文件
│   └── application-prod.toml
├── configs/                # 配置对象
│   ├── database.go
│   └── redis.go
├── apps/                   # 业务模块
│   ├── user/
│   │   ├── interface.go    # 接口定义
│   │   ├── impl/           # 业务实现 (controllers)
│   │   │   └── impl.go
│   │   └── api/            # API层 (apis)
│   │       └── http.go
│   └── order/
│       ├── interface.go
│       ├── impl/
│       └── api/
└── pkg/                    # 工具包 (default)
    ├── database/
    ├── cache/
    └── logger/
```

**示例代码结构**：

```go
// apps/user/interface.go - 定义接口
package user

type Service interface {
    GetUser(id string) (*User, error)
    CreateUser(user *User) error
}

// apps/user/impl/impl.go - 实现业务逻辑
package impl

import "github.com/infraboard/mcube/v2/ioc"

func init() {
    ioc.Controller().Registry(&UserServiceImpl{})
}

type UserServiceImpl struct {
    ioc.ObjectImpl
    DB *gorm.DB `ioc:"autowire=true;namespace=default"`
}

func (s *UserServiceImpl) Name() string {
    return "user"
}

// apps/user/api/http.go - 实现HTTP接口
package api

import "github.com/infraboard/mcube/v2/ioc"

func init() {
    ioc.Api().Registry(&UserAPI{})
}

type UserAPI struct {
    ioc.ObjectImpl
    UserSvc user.Service `ioc:"autowire=true;namespace=controllers"`
}

func (a *UserAPI) Name() string {
    return "user"
}
```

### 命名约定

**1. 对象命名**

```go
// ✅ 推荐：使用小写连字符
func (s *UserService) Name() string {
    return "user-service"
}

// ❌ 避免：使用Go类型名
func (s *UserService) Name() string {
    return "*impl.UserServiceImpl"  // 太长且不直观
}
```

**2. 配置节点命名**

```toml
# ✅ 推荐：与对象Name()一致
[database]
host = "localhost"

[redis-cache]
addr = "localhost:6379"
```

```go
func (c *DatabaseConfig) Name() string {
    return "database"  // 匹配配置节点
}
```

### 依赖管理原则

**1. 依赖方向**

```
configs → default → controllers → apis
  ↑         ↑           ↑          ↑
 配置      工具类      业务层     接口层
```

- ✅ apis可以依赖controllers
- ✅ controllers可以依赖default
- ❌ default不应该依赖controllers
- ❌ 避免跨层依赖

**2. 接口依赖**

```go
// ✅ 推荐：依赖接口
type UserAPI struct {
    UserSvc user.Service `ioc:"autowire=true;namespace=controllers"`
}

// ❌ 避免：直接依赖实现
type UserAPI struct {
    UserSvc *impl.UserServiceImpl `ioc:"autowire=true;namespace=controllers"`
}
```

**3. Priority()方法的使用（✅ 灵活设计）**

```go
// ✅ 推荐：返回固定值（最简单）
func (s *MyService) Priority() int {
    return -1  // 常量
}

// ✅ 推荐：使用常量定义
const MyServicePriority = 100

func (s *MyService) Priority() int {
    return MyServicePriority
}

// ✅ 支持：相对优先级（v2.0+ 完全安全）
func (s *CacheService) Priority() int {
    // 优先级是相对概念，可以基于其他对象确定
    db := ioc.Default().Get("database")
    if db != nil {
        return db.Priority() - 10  // 在数据库之后初始化
    }
    return 0
}

// ⚠️ 注意：避免循环依赖
func (s *ServiceA) Priority() int {
    b := ioc.Default().Get("service-b")
    // 如果 ServiceB 也依赖 ServiceA 的优先级，会有问题
    return b.Priority() + 1
}
```

**使用场景**：

- **常量优先级**：大多数场景，优先级在设计时就确定
- **相对优先级**：需要基于依赖关系动态确定初始化顺序
- **条件优先级**：根据环境或配置调整优先级

**为什么支持相对优先级？**
- ✅ v2.0 已修复死锁：Priority() 在加锁前调用，可以安全地访问容器
- ✅ 优先级本质是相对的："我应该在 A 之后，B 之前初始化"
- ✅ 更灵活的设计：支持复杂的初始化依赖关系

**性能说明**：Priority() 只在注册时调用一次，结果会被缓存，不影响运行时性能。

**历史说明**：早期版本中，在 Priority() 中调用 Get() 会导致死锁。v2.0 版本已完全修复此问题。

**依赖获取的位置选择**：
```go
// ✅ Priority 中获取：用于确定初始化顺序
func (s *MyService) Priority() int {
    dep := ioc.Default().Get("dependency")
    return dep.Priority() - 1
}

// ✅ Init 中获取：用于实际使用依赖
func (s *MyService) Init() error {
    s.dependency = ioc.Default().Get("dependency")
    return nil
}
```

**4. 避免循环依赖**

```go
// ❌ 错误：循环依赖
type ServiceA struct {
    B *ServiceB `ioc:"autowire=true"`
}
type ServiceB struct {
    A *ServiceA `ioc:"autowire=true"`
}

// ✅ 解决：使用接口或延迟获取
type ServiceA struct {
    ioc.ObjectImpl
}
func (a *ServiceA) GetB() *ServiceB {
    return ioc.Controller().Get("service-b").(*ServiceB)
}
```

### 错误处理

**1. Init()中的错误**

```go
func (s *DatabaseService) Init() error {
    db, err := sql.Open("mysql", s.dsn)
    if err != nil {
        // ✅ 返回带上下文的错误
        return fmt.Errorf("failed to open database: %w", err)
    }
    
    // ✅ 测试连接
    if err := db.Ping(); err != nil {
        return fmt.Errorf("failed to ping database: %w", err)
    }
    
    s.db = db
    return nil
}
```

**2. Close()中的错误处理**

```go
func (s *DatabaseService) Close(ctx context.Context) {
    if s.db == nil {
        return  // ✅ 防御性检查
    }
    
    // ✅ Close不返回error，需要自己处理
    if err := s.db.Close(); err != nil {
        log.Printf("failed to close database: %v", err)
    }
}
```

**3. 配置验证**

```go
func (c *DatabaseConfig) OnPostConfig() error {
    // ✅ 必填项检查
    if c.Host == "" {
        return fmt.Errorf("database.host is required")
    }
    
    // ✅ 范围检查
    if c.MaxConnections < 1 || c.MaxConnections > 1000 {
        return fmt.Errorf("database.max_connections must be between 1 and 1000")
    }
    
    // ✅ 格式验证
    if _, err := url.Parse(c.DSN); err != nil {
        return fmt.Errorf("invalid database.dsn: %w", err)
    }
    
    return nil
}
```

### 测试建议

**1. 单元测试**

```go
func TestUserService(t *testing.T) {
    // 创建测试用的IOC容器
    testStore := ioc.NewNamespaceStore("test", 0)
    
    // 注册Mock对象
    mockDB := &MockDatabase{}
    testStore.Registry(mockDB)
    
    // 创建待测试对象
    svc := &UserServiceImpl{}
    testStore.Registry(svc)
    
    // 执行依赖注入
    testStore.Autowire()
    
    // 初始化
    if err := svc.Init(); err != nil {
        t.Fatal(err)
    }
    
    // 测试
    user, err := svc.GetUser("user-001")
    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

**2. 集成测试**

```go
func TestMain(m *testing.M) {
    // 设置测试环境
    os.Setenv("ENV", "test")
    os.Setenv("DB_HOST", "localhost")
    
    // 注册测试对象
    ioc.Controller().Registry(&UserServiceImpl{})
    ioc.Api().Registry(&UserAPI{})
    
    // 初始化IOC容器
    if err := ioc.InitAll(); err != nil {
        log.Fatal(err)
    }
    
    // 运行测试
    code := m.Run()
    
    // 清理
    ioc.CloseAll(context.Background())
    
    os.Exit(code)
}
```

### 性能优化

**1. 避免过度反射**

```go
// ✅ 推荐：缓存对象引用
type UserAPI struct {
    ioc.ObjectImpl
    userSvc *UserService
}

func (a *UserAPI) Init() error {
    obj := ioc.Controller().Get("user-service")
    a.userSvc = obj.(*UserService)
    return nil
}

// ❌ 避免：每次都Get
func (a *UserAPI) HandleGetUser(c *gin.Context) {
    svc := ioc.Controller().Get("user-service").(*UserService)  // 每次都反射
    // ...
}
```

**2. 合理设置优先级**

```go
// 数据库应该最先初始化
func (d *Database) Priority() int {
    return 100
}

// 依赖数据库的服务优先级较低
func (s *UserService) Priority() int {
    return 50
}
```

### 安全建议

**1. 敏感信息保护**

```go
type DatabaseConfig struct {
    ioc.ObjectImpl
    Host     string `toml:"host"`
    Password string `toml:"password" json:"-"`  // ✅ json标签防止序列化泄露
}

// ✅ 自定义String()方法
func (c *DatabaseConfig) String() string {
    return fmt.Sprintf("DB{host=%s, password=***}", c.Host)
}
```

**2. 配置文件权限**

```bash
# ✅ 限制配置文件权限
chmod 600 etc/application.toml

# ✅ 不要提交包含敏感信息的配置到Git
echo "etc/*.toml" >> .gitignore
```

**3. 环境变量优先**

```go
// ✅ 生产环境使用环境变量而非配置文件
type DatabaseConfig struct {
    Password string `env:"DB_PASSWORD"`  // 优先从环境变量读取
}
```

### 日志规范

```go
func (s *UserService) Init() error {
    log.Printf("[IOC] Initializing UserService...")
    
    if err := s.connect(); err != nil {
        log.Printf("[IOC] UserService init failed: %v", err)
        return err
    }
    
    log.Printf("[IOC] UserService initialized successfully")
    return nil
}

func (s *UserService) Close(ctx context.Context) {
    log.Printf("[IOC] Closing UserService...")
    s.disconnect()
    log.Printf("[IOC] UserService closed")
}
```

---

## API参考

### 核心接口

#### Object接口

所有注册到IOC容器的对象必须实现此接口：

```go
type Object interface {
    Init() error                  // 对象初始化
    Name() string                 // 对象名称
    Version() string              // 版本号（默认1.0.0）
    Priority() int                // 优先级（数字越大越先初始化）
    Close(ctx context.Context)    // 优雅关闭
    Meta() ObjectMeta             // 元数据
}
```

#### StoreUser接口

用户操作接口，用于注册和获取对象：

```go
type StoreUser interface {
    // 注册对象
    Registry(obj Object) StoreUser
    
    // 批量注册对象
    RegistryAll(objs ...Object) StoreUser
    
    // 获取对象
    Get(name string, opts ...GetOption) Object
    
    // 加载对象到变量
    Load(obj any, opts ...GetOption) error
    
    // 列出所有对象名称
    List() []string
    
    // 对象数量
    Len() int
    
    // 遍历所有对象
    ForEach(fn func(*ObjectWrapper))
}
```

#### 生命周期钩子接口

可选实现的钩子接口：

```go
// 配置加载后钩子
type PostConfigHook interface {
    Object
    OnPostConfig() error
}

// 初始化前钩子
type PreInitHook interface {
    Object
    OnPreInit() error
}

// 初始化后钩子
type PostInitHook interface {
    Object
    OnPostInit() error
}

// 停止前钩子
type PreStopHook interface {
    Object
    OnPreStop(ctx context.Context) error
}

// 停止后钩子
type PostStopHook interface {
    Object
    OnPostStop(ctx context.Context) error
}
```

### 命名空间函数

#### ioc.Config()

```go
func Config() StoreUser
```

返回配置命名空间（优先级99），用于注册配置对象。

**示例**：
```go
ioc.Config().Registry(&DatabaseConfig{})
```

#### ioc.Default()

```go
func Default() StoreUser
```

返回默认命名空间（优先级9），用于注册工具类。

**示例**：
```go
ioc.Default().Registry(&DatabaseClient{})
```

#### ioc.Controller()

```go
func Controller() StoreUser
```

返回控制器命名空间（优先级0），用于注册业务控制器。

**示例**：
```go
ioc.Controller().Registry(&UserService{})
```

#### ioc.Api()

```go
func Api() StoreUser
```

返回API命名空间（优先级-99），用于注册API处理器。

**示例**：
```go
ioc.Api().Registry(&UserAPI{})
```

### 常用方法

#### Registry - 注册对象

```go
func (s *NamespaceStore) Registry(obj Object) StoreUser
```

注册对象到命名空间。

**参数**：
- `obj`: 实现Object接口的对象

**返回**：
- 返回自身，支持链式调用

**示例**：
```go
ioc.Default().Registry(&RedisClient{})

// 链式调用
ioc.Controller().
    Registry(&UserService{}).
    Registry(&OrderService{})
```

#### Get - 获取对象

```go
func (s *NamespaceStore) Get(name string, opts ...GetOption) Object
```

根据名称获取对象。

**参数**：
- `name`: 对象名称
- `opts`: 可选参数（如版本）

**返回**：
- Object实例，不存在返回nil

**示例**：
```go
// 获取默认版本
obj := ioc.Default().Get("redis-client")
client := obj.(*RedisClient)

// 获取指定版本
obj := ioc.Default().Get("cache", ioc.WithVersion("2.0.0"))
```

#### Load - 加载对象

```go
func (s *NamespaceStore) Load(obj any, opts ...GetOption) error
```

将对象加载到变量中（通过反射）。

**参数**：
- `obj`: 指向对象的指针
- `opts`: 可选参数

**返回**：
- error: 加载失败返回错误

**示例**：
```go
var db *gorm.DB
err := ioc.Default().Load(&db)
if err != nil {
    return err
}
```

#### List - 列出对象

```go
func (s *NamespaceStore) List() []string
```

返回所有已注册对象的名称列表。

**示例**：
```go
objects := ioc.Controller().List()
for _, name := range objects {
    fmt.Println(name)
}
```

#### ForEach - 遍历对象

```go
func (s *NamespaceStore) ForEach(fn func(*ObjectWrapper))
```

遍历所有已注册的对象。

**示例**：
```go
ioc.Api().ForEach(func(obj *ioc.ObjectWrapper) {
    fmt.Printf("API: %s (v%s)\n", obj.Name, obj.Version)
})
```

### 配置加载

#### LoadFromEnv - 从环境变量加载

```go
func (s *NamespaceStore) LoadFromEnv(prefix string) error
```

从环境变量加载配置。

**参数**：
- `prefix`: 环境变量前缀

**示例**：
```go
// 加载APP_为前缀的环境变量
ioc.Config().LoadFromEnv("APP")

// 例如: APP_DATABASE_HOST -> DatabaseConfig.Host
```

### 选项函数

#### WithVersion - 指定版本

```go
func WithVersion(version string) GetOption
```

获取或注入时指定对象版本。

**示例**：
```go
obj := ioc.Default().Get("cache", ioc.WithVersion("2.0.0"))
```

### 工具函数

#### ObjectUid - 对象唯一标识

```go
func ObjectUid(o *ObjectWrapper) string
```

返回对象的唯一标识（名称.版本）。

**示例**：
```go
uid := ioc.ObjectUid(wrapper)  // "user-service.1.0.0"
```

#### CompareVersion - 版本比较

```go
func CompareVersion(v1, v2 string) int
```

比较两个语义化版本号。

**返回值**：
- `1`: v1 > v2
- `-1`: v1 < v2
- `0`: v1 == v2

**示例**：
```go
result := ioc.CompareVersion("2.0.0", "1.5.0")  // 1
result := ioc.CompareVersion("1.0.0", "2.0.0")  // -1
result := ioc.CompareVersion("1.0.0", "1.0.0")  // 0
```

### 标签语法

依赖注入标签语法：

```go
`ioc:"key1=value1;key2=value2;key3=value3"`
```

**支持的参数**：

| 参数 | 类型 | 说明 | 示例 |
|------|------|------|------|
| `autowire` | bool | 是否自动注入 | `autowire=true` |
| `namespace` | string | 命名空间 | `namespace=controllers` |
| `version` | string | 对象版本 | `version=2.0.0` |

**示例**：

```go
type MyService struct {
    // 完整标签
    DB *gorm.DB `ioc:"autowire=true;namespace=default;version=1.0.0"`
    
    // 最简标签（使用默认版本）
    Cache *Redis `ioc:"autowire=true;namespace=default"`
}
```

### ObjectWrapper

对象包装器，包含对象及其元信息：

```go
type ObjectWrapper struct {
    Name           string      // 对象名称
    Version        string      // 对象版本
    Value          Object      // 对象实例
    Priority       int         // 优先级
    AllowOverwrite bool        // 是否允许覆盖
    Meta           ObjectMeta  // 元数据
}
```

### ObjectMeta

对象元数据：

```go
type ObjectMeta struct {
    CustomPathPrefix string            // 自定义API路径前缀
    Extra            map[string]string // 扩展字段
}
```

---

## 常见问题

### Q1: 为什么我的依赖注入失败了？

**A**: 检查以下几点：

1. **标签语法是否正确**
```go
// ✅ 正确
DB *gorm.DB `ioc:"autowire=true;namespace=default"`

// ❌ 错误：拼写错误
DB *gorm.DB `ioc:"autowire=ture;namespace=default"`

// ❌ 错误：缺少namespace
DB *gorm.DB `ioc:"autowire=true"`
```

2. **对象是否已注册**
```go
// 确保对象已注册
func init() {
    ioc.Default().Registry(&DatabaseClient{})
}
```

3. **命名空间是否正确**
```go
// 确保从正确的命名空间获取
type API struct {
    // 如果Service注册在controllers，必须指定namespace=controllers
    Svc *Service `ioc:"autowire=true;namespace=controllers"`
}
```

4. **字段必须是导出字段（公开）**
```go
// ✅ 正确：首字母大写
DB *gorm.DB `ioc:"autowire=true;namespace=default"`

// ❌ 错误：私有字段不会注入
db *gorm.DB `ioc:"autowire=true;namespace=default"`
```

### Q2: 循环依赖如何解决？

**A**: 有三种解决方案：

**方案1：延迟获取**
```go
type ServiceA struct {
    ioc.ObjectImpl
}

func (a *ServiceA) GetB() *ServiceB {
    return ioc.Controller().Get("service-b").(*ServiceB)
}
```

**方案2：使用接口解耦**
```go
type ServiceA struct {
    Handler BHandler `ioc:"autowire=true;namespace=controllers"`
}

type BHandler interface {
    Handle()
}
```

**方案3：重构设计**
- 通常循环依赖表示设计有问题
- 考虑提取公共逻辑到第三个服务
- 或者调整依赖方向

### Q3: 程序启动时卡住无响应，如何排查？

**A**: 如果程序启动时卡住，按以下步骤排查：

**步骤1：开启调试模式**
```bash
# 开启IOC调试日志
export IOC_DEBUG=true
go run main.go

# 查看日志输出，找到卡住的位置
# 观察对象注册和初始化的顺序
```

**步骤2：检查Init()方法**

最常见的问题是在 Init() 中发生了阻塞或错误：

```go
// ❌ 可能导致卡住的情况
func (s *MyService) Init() error {
    // 1. 无限循环或长时间阻塞
    for {
        // ...
    }
    
    // 2. 等待永远不会完成的操作
    <-s.neverCloseChan
    
    // 3. 循环依赖导致无法获取
    dep := ioc.Default().Get("循环依赖的对象")
    
    return nil
}
```

**步骤3：获取堆栈信息**
```bash
# 方式1：发送SIGQUIT信号（查看所有goroutine）
kill -QUIT <pid>

# 方式2：使用pprof
curl http://localhost:6060/debug/pprof/goroutine?debug=2

# 查找阻塞的goroutine，定位具体位置
```

**步骤4：常见问题和解决方案**

```go
// ✅ Init() 应该快速完成
func (s *MyService) Init() error {
    // 只做简单的初始化
    s.dependency = ioc.Default().Get("dependency")
    s.config = s.loadConfig()
    return nil
}

// ✅ 耗时操作放到 Start() 中
func (s *MyService) Start(ctx context.Context) error {
    // 连接数据库、启动服务等
    return s.connectDB()
}
```

**关于 v2.0 死锁修复**：

在早期版本中，在 Priority() 中调用 Get() 会导致死锁。**v2.0 已完全修复此问题**：
```go
// v2.0 中这样写不会死锁（但不推荐）
func (s *MyService) Priority() int {
    other := ioc.Default().Get("other")  // v2.0 安全，不会死锁
    return other.Priority() - 1
}

// ✅ 仍然推荐返回常量（性能更好，逻辑更清晰）
func (s *MyService) Priority() int {
    return -1
}
```

**最佳实践建议**：
- ✅ Priority() 返回常量（简单清晰）
- ✅ Init() 只做简单初始化，避免耗时操作
- ✅ 检查是否存在循环依赖
- ✅ 使用调试日志追踪执行流程

### Q4: 如何查看所有已注册的对象？

**A**: 使用List()和ForEach()方法：

```go
// 列出所有对象名称
objects := ioc.Controller().List()
for _, name := range objects {
    fmt.Println(name)
}

// 遍历对象详情
ioc.Controller().ForEach(func(obj *ioc.ObjectWrapper) {
    fmt.Printf("%s (v%s) - Priority: %d\n", 
        obj.Name, obj.Version, obj.Priority)
})
```

### Q5: 如何在单元测试中使用IOC？

**A**: 创建独立的测试容器：

```go
func TestMyService(t *testing.T) {
    // 方式1：使用独立的命名空间
    testStore := ioc.NewNamespaceStore("test", 0)
    testStore.Registry(&MockDatabase{})
    testStore.Registry(&MyService{})
    testStore.Autowire()
    
    // 方式2：替换全局对象
    original := ioc.Default().Get("database")
    ioc.Default().Registry(&MockDatabase{})
    
    // 测试完成后恢复
    defer func() {
        ioc.Default().Registry(original)
    }()
}
```

### Q6: 对象初始化顺序如何控制？

**A**: 通过三个维度控制：

1. **命名空间优先级**（最重要）
```go
configs (99) → default (9) → controllers (0) → apis (-99)
```

2. **Priority()方法**（同命名空间内）
```go
func (d *Database) Priority() int {
    return 100  // 数字越大越先初始化
}
```

3. **注册顺序**（同优先级时）
```go
func init() {
    ioc.Default().Registry(&A{})  // 先注册先初始化
    ioc.Default().Registry(&B{})
}
```

### Q7: 配置文件加载失败怎么办？

**A**: 检查以下几点：

1. **文件路径是否正确**
```go
server.DefaultConfig.ConfigFile.Paths = []string{
    "etc/application.toml",  // 相对于项目根目录
}
```

2. **配置节点名称是否匹配**
```toml
# 配置文件
[database]
host = "localhost"
```
```go
// 对象Name()必须匹配
func (c *DatabaseConfig) Name() string {
    return "database"  // 必须匹配[database]
}
```

3. **标签是否正确**
```go
type Config struct {
    // 确保标签格式正确
    Host string `toml:"host" env:"DB_HOST"`
}
```

### Q8: 如何实现配置热更新？

**A**: IOC容器本身不支持自动热更新，但可以手动实现：

```go
type DatabaseConfig struct {
    ioc.ObjectImpl
    Host string `toml:"host"`
    mu   sync.RWMutex
}

// 线程安全的读取
func (c *DatabaseConfig) GetHost() string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.Host
}

// 线程安全的更新
func (c *DatabaseConfig) UpdateHost(host string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.Host = host
}

// 监听配置文件变化
func watchConfig() {
    watcher, _ := fsnotify.NewWatcher()
    watcher.Add("etc/application.toml")
    
    for {
        select {
        case event := <-watcher.Events:
            if event.Op&fsnotify.Write == fsnotify.Write {
                reloadConfig()
            }
        }
    }
}
```

### Q9: 多个相同类型的对象如何区分？

**A**: 使用Name()和Version()区分：

```go
type CachePrimary struct {
    ioc.ObjectImpl
}
func (c *CachePrimary) Name() string {
    return "cache-primary"
}

type CacheSecondary struct {
    ioc.ObjectImpl
}
func (c *CacheSecondary) Name() string {
    return "cache-secondary"
}

// 或者使用版本区分
type CacheV1 struct {
    ioc.ObjectImpl
}
func (c *CacheV1) Version() string {
    return "1.0.0"
}

type CacheV2 struct {
    ioc.ObjectImpl
}
func (c *CacheV2) Version() string {
    return "2.0.0"
}
```

### Q10: Close()方法什么时候被调用？

**A**: 当应用关闭时自动调用：

```go
func main() {
    // server.Run会在收到信号时自动调用所有对象的Close()
    err := server.Run(context.Background())
}
```

手动调用：
```go
// 手动触发关闭
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

ioc.CloseAll(ctx)  // 关闭所有命名空间的对象
```

### Q11: 如何调试IOC容器问题？

**A**: 使用以下方法：

1. **开启调试模式**
```bash
# 设置环境变量开启详细日志
export IOC_DEBUG=true
go run main.go

# 日志会显示：
# - 对象注册过程
# - 锁的获取和释放
# - 依赖注入详情
# - 初始化顺序
```

2. **查看对象列表**
```go
fmt.Println("Registered objects:")
ioc.Controller().ForEach(func(obj *ioc.ObjectWrapper) {
    fmt.Printf("  - %s (v%s)\n", obj.Name, obj.Version)
})
```

3. **检查依赖关系**
```go
// 参考 DEPENDENCY_VISUALIZATION.md 生成依赖图
```

4. **使用GODEBUG追踪init**
```go
GODEBUG=inittrace=1 go run main.go
```

5. **检查初始化错误**
```go
if err := ioc.InitAll(); err != nil {
    log.Printf("Init failed: %v", err)
    // 查看具体是哪个对象初始化失败
}
```

### Q12: 性能考虑

**Q**: IOC容器对性能有影响吗？

**A**: 
- **初始化阶段**：使用反射有一定开销，但只在启动时执行一次
- **运行时**：依赖已注入完成，无额外开销
- **优化建议**：
  - 缓存Get()返回的对象引用，避免重复查找
  - 合理使用Priority()减少不必要的依赖等待
  - 避免在Init()中执行耗时操作，考虑延迟初始化

### Q13: 与其他框架的集成

**Q**: 如何与Gin、GORM等框架集成？

**A**: mcube/ioc已内置集成：

```go
import (
    _ "github.com/infraboard/mcube/v2/ioc/config/gin"      // Gin集成
    _ "github.com/infraboard/mcube/v2/ioc/config/datasource" // GORM集成
    _ "github.com/infraboard/mcube/v2/ioc/config/log"       // 日志集成
)

func main() {
    // 自动配置Gin、GORM等
    server.Run(context.Background())
}
```

---

## 更多资源

- **完整示例**：查看 [examples/](../examples/) 目录
- **依赖可视化**：参考 [DEPENDENCY_VISUALIZATION.md](DEPENDENCY_VISUALIZATION.md)
- **架构评估**：查看 [REVIEW.md](REVIEW.md)
- **死锁修复指南**：参考 [../docs/IOC_DEADLOCK_FIX.md](../docs/IOC_DEADLOCK_FIX.md)
- **问题反馈**：[GitHub Issues](https://github.com/infraboard/mcube/issues)

---

## 许可证

Apache License 2.0

