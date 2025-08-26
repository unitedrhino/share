package eventBus

import (
	"context"
	"fmt"
	"time"

	"gitee.com/unitedrhino/share/clients/mqttClient"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/events"
	"github.com/zeromicro/go-zero/core/logx"
)

type mqttEvent struct {
	mqttCli        *mqttClient.MqttClient
	topicConverter *TopicConverter
}

type mqttSubscription struct {
	*mqttClient.MqttSubscription
}

func NewMqttEvent(c conf.EventConf, serverName string, nodeID int64) (*mqttEvent, error) {
	cli, err := mqttClient.NewMqttClient(&c.Mqtt)
	if err != nil {
		return nil, err
	}
	return &mqttEvent{
		mqttCli:        cli,
		topicConverter: NewTopicConverter(),
	}, nil
}

func (m *mqttEvent) Subscribe(subj string, cb events.HandleFunc) (subscription, error) {
	// 将 NATS 格式的主题转换为 MQTT 格式
	mqttTopic := m.topicConverter.ConvertWildcards(subj)
	mqttTopic = "$inner/" + mqttTopic

	sub, err := m.mqttCli.Subscribe(mqttTopic, func(ctx context.Context, ts time.Time, msg []byte) error {
		emsg := events.GetEventMsg(msg)
		if emsg == nil {
			logx.Error(mqttTopic, string(msg))
			return nil
		}
		ctx = emsg.GetCtx()
		ctx, span := ctxs.StartSpan(ctx, mqttTopic, "")
		defer span.End()
		return cb(ctx, emsg.GetTs(), emsg.GetData())
	})
	if err != nil {
		return nil, err
	}
	return &mqttSubscription{sub}, nil
}

// SubscribeWithQoS 支持指定 QoS 级别的订阅
func (m *mqttEvent) SubscribeWithQoS(topic string, qos byte, cb events.HandleFunc) (subscription, error) {
	// 将 NATS 格式的主题转换为 MQTT 格式
	mqttTopic := m.topicConverter.ConvertWildcards(topic)
	mqttTopic = "$inner/" + mqttTopic
	sub, err := m.mqttCli.SubscribeWithQoS(mqttTopic, qos, cb)
	if err != nil {
		return nil, err
	}
	return &mqttSubscription{sub}, nil
}

func (m *mqttEvent) Publish(ctx context.Context, topic string, payload []byte) error {
	// 将 NATS 格式的主题转换为 MQTT 格式
	mqttTopic := m.topicConverter.NatsToMqtt(topic)
	mqttTopic = "$inner/" + mqttTopic
	pubMsg := events.NewEventMsg(ctx, payload)

	return m.mqttCli.Publish(ctx, mqttTopic, pubMsg)
}

// PublishWithQoS 支持指定 QoS 和 retained 标志的发布
func (m *mqttEvent) PublishWithQoS(ctx context.Context, topic string, payload []byte, qos byte, retained bool) error {
	// 将 NATS 格式的主题转换为 MQTT 格式
	mqttTopic := m.topicConverter.NatsToMqtt(topic)
	eventData := events.NewEventMsg(ctx, payload)
	mqttTopic = "$inner/" + mqttTopic
	return m.mqttCli.PublishRaw(mqttTopic, qos, retained, eventData)
}

func (m *mqttEvent) QueueSubscribe(topic, queue string, cb events.HandleFunc) (subscription, error) {
	// 将 NATS 格式的主题转换为 MQTT 格式
	mqttTopic := m.topicConverter.ConvertWildcards(topic)
	// MQTT 的共享订阅使用 $share/group/topic 格式
	mqttTopic = "$inner/" + mqttTopic
	sharedTopic := fmt.Sprintf("$share/%s/%s", queue, mqttTopic)
	sub, err := m.mqttCli.Subscribe(sharedTopic, func(ctx context.Context, ts time.Time, msg []byte) error {
		emsg := events.GetEventMsg(msg)
		if emsg == nil {
			logx.Error(mqttTopic, string(msg))
			return nil
		}
		ctx = emsg.GetCtx()
		ctx, span := ctxs.StartSpan(ctx, mqttTopic, "")
		defer span.End()
		return cb(ctx, emsg.GetTs(), emsg.GetData())
	})
	if err != nil {
		return nil, err
	}
	return &mqttSubscription{sub}, nil
}

// QueueSubscribeWithQoS 支持指定 QoS 级别的队列订阅
func (m *mqttEvent) QueueSubscribeWithQoS(topic, queue string, qos byte, cb events.HandleFunc) (subscription, error) {
	// 将 NATS 格式的主题转换为 MQTT 格式
	mqttTopic := m.topicConverter.ConvertWildcards(topic)
	// MQTT 的共享订阅使用 $share/group/topic 格式
	mqttTopic = "$inner/" + mqttTopic
	sharedTopic := fmt.Sprintf("$share/%s/%s", mqttTopic, queue)
	return m.SubscribeWithQoS(sharedTopic, qos, cb)
}

func (m *mqttSubscription) Unsubscribe() error {
	return m.MqttSubscription.Unsubscribe()
}
