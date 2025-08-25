package eventBus

import (
	"testing"
)

func TestTopicConverter_NatsToMqtt(t *testing.T) {
	converter := NewTopicConverter()

	tests := []struct {
		name      string
		natsTopic string
		expected  string
	}{
		{
			name:      "基本转换",
			natsTopic: "server.core.project.info.delete",
			expected:  "server/core/project/info/delete",
		},
		{
			name:      "单个层级",
			natsTopic: "test",
			expected:  "test",
		},
		{
			name:      "带通配符",
			natsTopic: "server.core.*.delete",
			expected:  "server/core/*/delete",
		},
		{
			name:      "带多级通配符",
			natsTopic: "server.core.>",
			expected:  "server/core/>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.NatsToMqtt(tt.natsTopic)
			if result != tt.expected {
				t.Errorf("NatsToMqtt() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTopicConverter_MqttToNats(t *testing.T) {
	converter := NewTopicConverter()

	tests := []struct {
		name      string
		mqttTopic string
		expected  string
	}{
		{
			name:      "基本转换",
			mqttTopic: "server/core/project/info/delete",
			expected:  "server.core.project.info.delete",
		},
		{
			name:      "单个层级",
			mqttTopic: "test",
			expected:  "test",
		},
		{
			name:      "带通配符",
			mqttTopic: "server/core/+/delete",
			expected:  "server.core.+.delete",
		},
		{
			name:      "带多级通配符",
			mqttTopic: "server/core/#",
			expected:  "server.core.#",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.MqttToNats(tt.mqttTopic)
			if result != tt.expected {
				t.Errorf("MqttToNats() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTopicConverter_ConvertWildcards(t *testing.T) {
	converter := NewTopicConverter()

	tests := []struct {
		name      string
		natsTopic string
		expected  string
	}{
		{
			name:      "基本转换",
			natsTopic: "server.core.project.info.delete",
			expected:  "server/core/project/info/delete",
		},
		{
			name:      "单级通配符",
			natsTopic: "server.core.*.delete",
			expected:  "server/core/+/delete",
		},
		{
			name:      "多级通配符",
			natsTopic: "server.core.>",
			expected:  "server/core/#",
		},
		{
			name:      "混合通配符",
			natsTopic: "server.core.*.>",
			expected:  "server/core/+/#",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ConvertWildcards(tt.natsTopic)
			if result != tt.expected {
				t.Errorf("ConvertWildcards() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTopicConverter_ConvertMqttWildcardsToNats(t *testing.T) {
	converter := NewTopicConverter()

	tests := []struct {
		name      string
		mqttTopic string
		expected  string
	}{
		{
			name:      "基本转换",
			mqttTopic: "server/core/project/info/delete",
			expected:  "server.core.project.info.delete",
		},
		{
			name:      "单级通配符",
			mqttTopic: "server/core/+/delete",
			expected:  "server.core.*.delete",
		},
		{
			name:      "多级通配符",
			mqttTopic: "server/core/#",
			expected:  "server.core.>",
		},
		{
			name:      "混合通配符",
			mqttTopic: "server/core/+/#",
			expected:  "server.core.*.>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ConvertMqttWildcardsToNats(tt.mqttTopic)
			if result != tt.expected {
				t.Errorf("ConvertMqttWildcardsToNats() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTopicConverter_IsValidMqttTopic(t *testing.T) {
	converter := NewTopicConverter()

	tests := []struct {
		name     string
		topic    string
		expected bool
	}{
		{"有效主题", "server/core/project/info/delete", true},
		{"根主题", "/", true},
		{"单级主题", "test", true},
		{"带通配符", "server/core/+/delete", true},
		{"多级通配符", "server/core/#", true},
		{"空主题", "", false},
		{"以斜杠开头", "/server/core", false},
		{"以斜杠结尾", "server/core/", false},
		{"连续斜杠", "server//core", false},
		{"通配符位置错误", "server/core/#/delete", false},
		{"通配符混合", "server/core/+delete", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.IsValidMqttTopic(tt.topic)
			if result != tt.expected {
				t.Errorf("IsValidMqttTopic() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTopicConverter_IsValidNatsTopic(t *testing.T) {
	converter := NewTopicConverter()

	tests := []struct {
		name     string
		topic    string
		expected bool
	}{
		{"有效主题", "server.core.project.info.delete", true},
		{"单级主题", "test", true},
		{"带通配符", "server.core.*.delete", true},
		{"多级通配符", "server.core.>", true},
		{"空主题", "", false},
		{"以点开头", ".server.core", false},
		{"以点结尾", "server.core.", false},
		{"连续点", "server..core", false},
		{"通配符位置错误", "server.core.>.delete", false},
		{"通配符混合", "server.core.*delete", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.IsValidNatsTopic(tt.topic)
			if result != tt.expected {
				t.Errorf("IsValidNatsTopic() = %v, want %v", result, tt.expected)
			}
		})
	}
}
