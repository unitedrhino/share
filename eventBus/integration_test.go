package eventBus

import (
	"context"
	"strings"
	"testing"
	"time"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/events"
)

// 全局测试配置变量
var (
	// NATS 配置
	NatsURL   = "nats://localhost:4222"
	NatsUser  = ""
	NatsPass  = ""
	NatsToken = ""

	// MQTT 配置
	MqttBrokers  = []string{"tcp://localhost:1883"}
	MqttUser     = ""
	MqttPass     = ""
	MqttClientID = "share-test-client"
)

// 环境变量约定：
// NATS_URL、NATS_USER、NATS_PASS（可选）
// MQTT_BROKERS（逗号分隔）、MQTT_USER、MQTT_PASS（可选）、MQTT_CLIENT_ID（可选）

// TestHandleFuncCompatibility 测试通用 HandleFunc 的兼容性
func TestHandleFuncCompatibility(t *testing.T) {
	// 测试通用 HandleFunc 定义
	var handleFunc events.HandleFunc = func(ctx context.Context, ts time.Time, msg []byte) error {
		return nil
	}

	// 验证函数签名正确
	if handleFunc == nil {
		t.Fatal("HandleFunc 定义失败")
	}

	// 测试调用
	err := handleFunc(context.Background(), time.Now(), []byte("test"))
	if err != nil {
		t.Fatalf("HandleFunc 调用失败: %v", err)
	}
}

func TestNATS_SendReceive(t *testing.T) {
	url := strings.TrimSpace(NatsURL)
	if url == "" {
		t.Skip("未设置 NatsURL，跳过NATS收发测试")
	}

	c := conf.EventConf{
		Mode: conf.EventModeNats,
		Nats: conf.NatsConf{
			Url:   url,
			User:  NatsUser,
			Pass:  NatsPass,
			Token: NatsToken,
		},
	}

	bus, err := NewFastEvent(c, "share-test", time.Now().UnixNano())
	if err != nil {
		t.Fatalf("NewFastEvent(NATS) 失败: %v", err)
	}

	topic := "server.test.nats.echo"
	recv := make(chan string, 1)

	if err := bus.Subscribe(topic, func(ctx context.Context, t time.Time, body []byte) error {
		recv <- string(body)
		return nil
	}); err != nil {
		t.Fatalf("Subscribe 失败: %v", err)
	}

	if err := bus.Start(); err != nil {
		t.Fatalf("Start 失败: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg := "hello-nats"
	if err := bus.Publish(ctx, topic, msg); err != nil {
		t.Fatalf("Publish 失败: %v", err)
	}

	select {
	case got := <-recv:
		if got != msg {
			t.Fatalf("收消息不一致: got=%q want=%q", got, msg)
		}
	case <-ctx.Done():
		t.Fatalf("等待接收超时")
	}
}
func TestMQTT_SendReceive(t *testing.T) {
	if len(MqttBrokers) == 0 {
		t.Skip("未设置 MqttBrokers，跳过MQTT收发测试")
	}

	c := conf.EventConf{
		Mode: conf.EventModeMqtt,
		Mqtt: conf.MqttConf{
			ClientID: firstNonEmpty(MqttClientID, "share-test-client"),
			Brokers:  MqttBrokers,
			User:     MqttUser,
			Pass:     MqttPass,
			ConnNum:  1,
		},
	}

	bus, err := NewFastEvent(c, "share-test", time.Now().UnixNano())
	if err != nil {
		t.Fatalf("NewFastEvent(MQTT) 失败: %v", err)
	}

	// 使用点分主题，内部会转换为 MQTT 斜杠主题
	topic := "server/test/mqtt/echo"
	recv := make(chan string, 1)

	if err := bus.Subscribe(topic, func(ctx context.Context, t time.Time, body []byte) error {
		recv <- string(body)
		return nil
	}); err != nil {
		t.Fatalf("Subscribe 失败: %v", err)
	}

	if err := bus.Start(); err != nil {
		t.Fatalf("Start 失败: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Second)
	defer cancel()

	msg := "hello-mqtt"
	if err := bus.Publish(ctx, topic, msg); err != nil {
		t.Fatalf("Publish 失败: %v", err)
	}

	select {
	case got := <-recv:
		if got != msg {
			t.Fatalf("收消息不一致: got=%q want=%q", got, msg)
		}
	case <-ctx.Done():
		t.Fatalf("等待接收超时")
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
