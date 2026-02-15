# 依赖可视化使用指南

IOC 包支持两种依赖声明方式，都能在依赖图中正确展示：

## 1️⃣ 声明式依赖（推荐）

通过 `ioc` 标签自动检测，**无需额外代码**：

```go
type UserService struct {
    ioc.ObjectImpl
    Logger *Logger `ioc:"autowire=true;namespace=default"`
    Cache  *Cache  `ioc:"autowire=true;namespace=default"`
}
```

**优点**：
- ✅ 自动检测，零配置
- ✅ 依赖关系一目了然
- ✅ 编译期类型安全

## 2️⃣ 命令式依赖

手动通过 `Get()` 获取依赖时，实现 `DependencyDeclarer` 接口来声明：

```go
type EmailService struct {
    ioc.ObjectImpl
    // 无 ioc 标签，在 Init() 中手动获取
    logger     *Logger
    smtpClient *SMTPClient
}

func (s *EmailService) Init() error {
    ns := ioc.DefaultStore.Namespace("default")
    s.logger = ns.Get("*app.Logger").(*Logger)
    s.smtpClient = ns.Get("*app.SMTPClient").(*SMTPClient)
    return nil
}

// 实现接口来声明命令式依赖
func (s *EmailService) DeclareDependencies() []ioc.DependencyInfo {
    return []ioc.DependencyInfo{
        {
            Name:      "*app.Logger",
            Namespace: "default",
            FieldName: "logger", // 可选，用于文档说明
        },
        {
            Name:      "*app.SMTPClient",
            Namespace: "default",
            FieldName: "smtpClient",
        },
    }
}
```

**使用场景**：
- 动态依赖（运行时决定）
- 条件依赖（某些情况下才加载）
- 循环依赖的临时解决方案

## 3️⃣ 混合使用

两种方式可以共存：

```go
type NotificationService struct {
    ioc.ObjectImpl
    // 声明式依赖
    Logger *Logger `ioc:"autowire=true;namespace=default"`
    
    // 命令式依赖（无标签）
    queue *Queue
}

func (n *NotificationService) Init() error {
    ns := ioc.DefaultStore.Namespace("default")
    n.queue = ns.Get("*app.Queue").(*Queue)
    return nil
}

// 只需声明命令式依赖，标签依赖会自动检测
func (n *NotificationService) DeclareDependencies() []ioc.DependencyInfo {
    return []ioc.DependencyInfo{
        {Name: "*app.Queue", Namespace: "default", FieldName: "queue"},
    }
}
```

## 📊 依赖可视化

### 打印依赖树

```go
// 打印单个命名空间的依赖树
ioc.Default().PrintDependencies()

// 打印所有命名空间
ioc.DefaultStore.PrintAllDependencies()
```

输出示例：
```
=== default Namespace Dependency Tree ===
├─ *app.Logger@1.5.0
├─ *app.Cache@1.2.0
├─ *app.UserRepository@3.0.0 (1 deps)
   └─ *app.Database@2.0.0
├─ *app.UserService@4.1.0 (3 deps)
   ├─ *app.UserRepository@3.0.0 (already shown)
   ├─ *app.Logger@1.5.0 (already shown)
   └─ *app.Cache@1.2.0 (already shown)
```

### 打印依赖统计

```go
ioc.Default().PrintDependencySummary()
```

输出示例：
```
=== default Namespace Summary ===
  📊 Total Objects: 4
  🌿 Leaf Objects (no deps): 2
  🔗 Objects with deps: 2
  ⬆️  Most dependencies: *app.UserService (3 deps)
  ⬇️  Most depended on: *app.UserRepository (used by 1 objects)
```

### 导出 Markdown 文档

```go
markdown := ioc.Default().ExportDependenciesToMarkdown()
os.WriteFile("DEPENDENCIES.md", []byte(markdown), 0644)
```

## 🎯 最佳实践

1. **优先使用声明式依赖**（ioc 标签）
   - 代码更简洁
   - 依赖关系更清晰
   - 无需额外维护

2. **命令式依赖适用于特殊场景**
   - 动态依赖选择
   - 条件加载
   - 复杂初始化逻辑

3. **调试时启用依赖可视化**
   ```go
   if os.Getenv("DEBUG") == "true" {
       ioc.DefaultStore.PrintAllDependencies()
   }
   ```

4. **定期生成依赖文档**
   ```bash
   # 在 CI/CD 中生成
   go run scripts/export_deps.go > docs/IOC_DEPENDENCIES.md
   ```

## ⚠️ 注意事项

1. **DependencyDeclarer 是可选的**
   - 只有使用命令式依赖（手动 Get()）时才需要实现
   - 声明式依赖（ioc 标签）会自动检测

2. **依赖声明不会影响运行时行为**
   - 仅用于可视化和文档生成
   - 不会自动进行依赖注入

3. **循环依赖会被自动检测**
   - 在依赖树中标记为 `⚠️  (circular dependency)`
   - 建议重构代码消除循环依赖

## 📚 更多示例

参考测试文件：[dependency_test.go](dependency_test.go)
- `TestImperativeDependencies` - 命令式依赖完整示例
- `TestMixedDependencies` - 混合依赖示例
- `TestPrintDependencies` - 依赖可视化示例
