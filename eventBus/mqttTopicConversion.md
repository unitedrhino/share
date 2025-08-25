# MQTT 主题格式转换说明

## 概述

MQTT 事件总线自动处理 NATS 和 MQTT 之间的主题格式转换，使开发者可以继续使用现有的 NATS 格式主题，而无需修改代码。

## 主题格式对比

| 特性 | NATS 格式 | MQTT 格式 |
|------|-----------|-----------|
| 分隔符 | 点 (.) | 斜杠 (/) |
| 单级通配符 | * | + |
| 多级通配符 | > | # |
| 示例 | `server.core.project.info.delete` | `server/core/project/info/delete` |

## 自动转换规则

### 1. 基本主题转换

```go
// NATS 格式 -> MQTT 格式
"server.core.project.info.delete" -> "server/core/project/info/delete"
"test.topic" -> "test/topic"
```

### 2. 通配符转换

```go
// 单级通配符
"server.core.*.delete" -> "server/core/+/delete"

// 多级通配符
"server.core.>" -> "server/core/#"

// 混合通配符
"server.core.*.>" -> "server/core/+/#"
```

### 3. 共享订阅转换

```go
// NATS 队列订阅
QueueSubscribe("server.core.test", "group1", handler)

// 转换为 MQTT 共享订阅
"$share/group1/server/core/test"
```

## 使用示例

### 基本订阅和发布

```go
// 使用 NATS 格式的主题（会自动转换为 MQTT 格式）
mqttEvent, _ := NewMqttEvent(eventConf, "server-name", 1)

// 订阅（内部转换为 MQTT 格式）
sub, err := mqttEvent.Subscribe("server.core.project.info.delete", func(ctx context.Context, msg []byte) error {
    log.Printf("收到消息: %s", string(msg))
    return nil
})

// 发布（内部转换为 MQTT 格式）
err = mqttEvent.Publish(context.Background(), "server.core.project.info.delete", []byte("Hello"))
```

### 通配符订阅

```go
// 订阅所有项目相关的消息
sub, err := mqttEvent.Subscribe("server.core.project.*", func(ctx context.Context, msg []byte) error {
    log.Printf("收到项目相关消息: %s", string(msg))
    return nil
})

// 发布到具体主题
err = mqttEvent.Publish(context.Background(), "server.core.project.info.delete", []byte("删除项目"))
err = mqttEvent.Publish(context.Background(), "server.core.project.info.create", []byte("创建项目"))
// 以上两个消息都会被通配符订阅接收到
```

### 队列订阅（负载均衡）

```go
// 多个服务实例使用相同的队列名称
sub1, _ := mqttEvent.QueueSubscribe("server.core.task.process", "workers", handler1)
sub2, _ := mqttEvent.QueueSubscribe("server.core.task.process", "workers", handler2)

// 发布任务（只会被其中一个 worker 处理）
err = mqttEvent.Publish(context.Background(), "server.core.task.process", taskData)
```

## 主题验证

系统会自动验证主题格式的有效性：

### MQTT 主题验证规则

- 不能为空
- 不能以 `/` 开头或结尾（除非是根主题 `/`）
- 不能包含连续的 `/`
- `+` 通配符只能单独使用
- `#` 通配符只能在最后

### NATS 主题验证规则

- 不能为空
- 不能以 `.` 开头或结尾
- 不能包含连续的 `.`
- `*` 通配符可以单独使用
- `>` 通配符只能在最后

## 日志输出

系统会在转换时输出详细的日志信息：

```
MQTT 订阅主题转换: server.core.project.info.delete -> server/core/project/info/delete
MQTT 发布主题转换: server.core.project.info.delete -> server/core/project/info/delete
MQTT 队列订阅主题转换: server.core.task.process -> $share/workers/server/core/task/process (Queue: workers)
```

## 注意事项

1. **向后兼容**: 现有使用 NATS 格式主题的代码无需修改
2. **通配符行为**: MQTT 和 NATS 的通配符行为略有不同，但基本功能一致
3. **性能影响**: 主题转换的开销很小，对性能影响微乎其微
4. **调试**: 可以通过日志查看实际使用的 MQTT 主题格式

## 手动转换

如果需要手动转换主题格式，可以使用 `TopicConverter`：

```go
converter := NewTopicConverter()

// NATS 到 MQTT
mqttTopic := converter.NatsToMqtt("server.core.project.info.delete")
// 结果: "server/core/project/info/delete"

// MQTT 到 NATS
natsTopic := converter.MqttToNats("server/core/project/info/delete")
// 结果: "server.core.project.info.delete"

// 包含通配符的转换
mqttTopic := converter.ConvertWildcards("server.core.*.delete")
// 结果: "server/core/+/delete"
```
