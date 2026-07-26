# tools/diff

基于 struct `diff` 标签的结构体字段对比，适合审计日志、变更事件摘要。

## 快速开始

```go
type Spec struct {
    Name  string            `diff:"名称"`
    Items []Item            `diff:"列表"`
    Meta  map[string]string `diff:"元数据"`
}

type Item struct {
    ID   string `diff:"ID,key"`
    Name string `diff:"名称"`
}

before := Spec{Name: "a", Items: []Item{{ID: "1", Name: "x"}}}
after := Spec{Name: "b", Items: []Item{{ID: "1", Name: "y"}}}

// 推荐：先归一化 nil map/slice，再对比（避免 nil vs [] 误报）
records := diff.CompareNormalized(before, after)

// 或手动追加脱敏字段
records = diff.AppendRecords(records,
    diff.MaskedChange("Secret", "密钥", "******", diff.LevelWarn),
)

fmt.Println(diff.Summary(records))
// 名称: a -> b, 列表[1].名称: x -> y
```

错误返回版本（类型不一致时不 panic）：

```go
records, err := diff.CompareE(before, after)
if errors.Is(err, diff.ErrTypeMismatch) {
    // before/after 类型不同
}
```

## 标签约定

| 标签 | 含义 |
|------|------|
| `diff:"描述"` | 参与对比，`FieldDesc` 使用该描述 |
| `diff:"-"` | 忽略该字段 |
| 无标签 | 默认参与对比（`FieldDesc` 用字段名）；`Options.OnlyTaggedFields=true` 时忽略 |
| `diff:"ID,key"` | 切片元素的 identity；同 key 视为同一条，做原地字段 diff |
| `diff:"名称,level=warn"` | 默认级别：`info` / `warn` / `error` |

示例：

```go
type User struct {
    ID       string `diff:"用户ID,key"`
    Name     string `diff:"姓名"`
    Password string `diff:"-"`
    Role     string `diff:"角色,level=warn"`
}
```

## API

| 函数 | 说明 |
|------|------|
| `Compare` / `CompareE` | 直接对比；`Compare` 在类型不一致时 panic |
| `CompareNormalized` / `CompareNormalizedE` | 先 `NormalizeNilContainers` 再对比 |
| `NormalizeNilContainers` | 深拷贝，并把 nil map/slice 变成空容器 |
| `Summary` | 拼成 `字段: old -> new` 文本 |
| `MaskedChange` / `AppendRecords` | 手工构造/追加脱敏变更 |

### Options

```go
diff.Compare(a, b, &diff.Options{
    OnlyTaggedFields: true,                    // 仅对比带 diff 标签的字段
    IgnoreFields:     map[string]bool{"Meta": true},
    FieldLevels:      map[string]diff.DiffLevel{
        "Name": diff.LevelError,               // 覆盖 tag 默认级别（按 FieldPath）
    },
})
```

## 对比规则

- **标量 / `time.Time`**：值不等则记一条 diff。
- **指针**：nil 状态变化记一条；双方非 nil 则解引用继续比。
- **`any` / interface**：比较动态值；类型变化记整值 diff。
- **结构体**：按导出字段递归。
- **切片 / 数组**
  - 元素为 struct（或 `*struct`）：用 `key` 字段匹配；无 `key` 时用全部带 `diff` 标签的标量拼接作 identity。
  - 匹配成功：递归字段 diff，路径带 identity，如 `Items[1].Name` / `列表[1].名称`。
  - 未匹配：记新增（`OldValue=新增`）或移除（`NewValue=移除`）。
- **Map**：按 key 对齐；双方都有则对 value **递归**字段级对比（如 `Meta[x].Name`）。

## DiffRecord

```go
type DiffRecord struct {
    FieldPath string    // 如 Basic.ID、Items[1].Name、Meta[env]
    FieldDesc string    // 人类可读描述
    OldValue  any
    NewValue  any
    Level     DiffLevel // Info / Warn / Error
}
```

常量：

- `MarkerAdd = "新增"`
- `MarkerRemove = "移除"`

## 注意事项

1. `Compare` / `CompareE` 要求两边类型完全一致（`reflect.Type` 相等）。
2. 需要忽略「nil 切片 vs 空切片」时用 `CompareNormalized`。
3. `NormalizeNilContainers` 返回深拷贝，修改结果不会影响原对象。
4. 切片元素路径使用 identity（有 `key` 时）或下标，避免多条元素变更路径撞车。
5. `Options.FieldLevels` 的 key 是完整 `FieldPath`（含切片下标/key、map key）。
