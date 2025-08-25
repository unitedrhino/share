package eventBus

import (
	"strings"
)

// TopicConverter 主题格式转换器
type TopicConverter struct{}

// NewTopicConverter 创建主题转换器实例
func NewTopicConverter() *TopicConverter {
	return &TopicConverter{}
}

// NatsToMqtt 将 NATS 主题格式转换为 MQTT 主题格式
// NATS: server.core.project.info.delete -> MQTT: server/core/project/info/delete
func (tc *TopicConverter) NatsToMqtt(natsTopic string) string {
	// 将点分隔符替换为斜杠分隔符
	return strings.ReplaceAll(natsTopic, ".", "/")
}

// MqttToNats 将 MQTT 主题格式转换为 NATS 主题格式
// MQTT: server/core/project/info/delete -> NATS: server.core.project.info.delete
func (tc *TopicConverter) MqttToNats(mqttTopic string) string {
	// 将斜杠分隔符替换为点分隔符
	return strings.ReplaceAll(mqttTopic, "/", ".")
}

// ConvertWildcards 转换通配符格式
// NATS: server.core.*.delete -> MQTT: server/core/+/delete
// NATS: server.core.> -> MQTT: server/core/#
func (tc *TopicConverter) ConvertWildcards(natsTopic string) string {
	// 先转换分隔符
	mqttTopic := tc.NatsToMqtt(natsTopic)

	// 转换通配符
	// NATS 的 * 对应 MQTT 的 +
	mqttTopic = strings.ReplaceAll(mqttTopic, "*", "+")

	// NATS 的 > 对应 MQTT 的 #
	mqttTopic = strings.ReplaceAll(mqttTopic, ">", "#")

	return mqttTopic
}

// ConvertMqttWildcardsToNats 将 MQTT 通配符转换为 NATS 通配符
func (tc *TopicConverter) ConvertMqttWildcardsToNats(mqttTopic string) string {
	// 先转换分隔符
	natsTopic := tc.MqttToNats(mqttTopic)

	// 转换通配符
	// MQTT 的 + 对应 NATS 的 *
	natsTopic = strings.ReplaceAll(natsTopic, "+", "*")

	// MQTT 的 # 对应 NATS 的 >
	natsTopic = strings.ReplaceAll(natsTopic, "#", ">")

	return natsTopic
}

// IsValidMqttTopic 验证 MQTT 主题格式是否有效
func (tc *TopicConverter) IsValidMqttTopic(topic string) bool {
	if len(topic) == 0 {
		return false
	}

	// MQTT 主题不能以 / 开头或结尾（除非是根主题）
	if topic != "/" && (strings.HasPrefix(topic, "/") || strings.HasSuffix(topic, "/")) {
		return false
	}

	// 检查连续的斜杠
	if strings.Contains(topic, "//") {
		return false
	}

	// 检查通配符位置
	parts := strings.Split(topic, "/")
	for i, part := range parts {
		// + 通配符只能单独使用
		if part == "+" {
			continue
		}

		// # 通配符只能在最后
		if part == "#" {
			return i == len(parts)-1
		}

		// 其他部分不能包含通配符
		if strings.Contains(part, "+") || strings.Contains(part, "#") {
			return false
		}
	}

	return true
}

// IsValidNatsTopic 验证 NATS 主题格式是否有效
func (tc *TopicConverter) IsValidNatsTopic(topic string) bool {
	if len(topic) == 0 {
		return false
	}

	// NATS 主题不能以 . 开头或结尾
	if strings.HasPrefix(topic, ".") || strings.HasSuffix(topic, ".") {
		return false
	}

	// 检查连续的点
	if strings.Contains(topic, "..") {
		return false
	}

	// 检查通配符位置
	parts := strings.Split(topic, ".")
	for i, part := range parts {
		// * 通配符可以单独使用
		if part == "*" {
			continue
		}

		// > 通配符只能在最后
		if part == ">" {
			return i == len(parts)-1
		}

		// 其他部分不能包含通配符
		if strings.Contains(part, "*") || strings.Contains(part, ">") {
			return false
		}
	}

	return true
}
