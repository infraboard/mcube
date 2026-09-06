# 事件总线

支持:

+ kafka
+ nats
+ rabbitmq

## 订阅

`TopicSubscribe` / `QueueSubscribe` 为进程级长订阅。需要按资源动态订阅、断开时退订时，使用 `Subscribe`：

```go
sub, err := bus.GetService().Subscribe(ctx, "flow.agent.cmd."+agentID, func(e *bus.Event) {
    // handle e.Data
})
if err != nil {
    return err
}
defer sub.Unsubscribe()
```

语义与 `TopicSubscribe` 相同（广播：每个订阅者各收一份）。仅持有者订阅某个 subject 时，即点对点投递。

