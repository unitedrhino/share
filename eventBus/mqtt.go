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
	logx.Infof("MQTT 订阅主题转换: %s -> %s", subj, mqttTopic)

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
func (m *mqttEvent) SubscribeWithQoS(subj string, qos byte, cb events.HandleFunc) (subscription, error) {
	// 将 NATS 格式的主题转换为 MQTT 格式
	mqttTopic := m.topicConverter.ConvertWildcards(subj)
	mqttTopic = "$inner/" + mqttTopic
	logx.Infof("MQTT 订阅主题转换: %s -> %s (QoS: %d)", subj, mqttTopic, qos)

	sub, err := m.mqttCli.SubscribeWithQoS(mqttTopic, qos, cb)
	if err != nil {
		return nil, err
	}
	return &mqttSubscription{sub}, nil
}

func (m *mqttEvent) Publish(ctx context.Context, topic string, arg []byte) error {
	// 将 NATS 格式的主题转换为 MQTT 格式
	mqttTopic := m.topicConverter.NatsToMqtt(topic)
	mqttTopic = "$inner/" + mqttTopic
	logx.Infof("MQTT 发布主题转换: %s -> %s", topic, mqttTopic)

	return m.mqttCli.Publish(ctx, mqttTopic, arg)
}

// PublishWithQoS 支持指定 QoS 和 retained 标志的发布
func (m *mqttEvent) PublishWithQoS(ctx context.Context, topic string, arg []byte, qos byte, retained bool) error {
	// 将 NATS 格式的主题转换为 MQTT 格式
	mqttTopic := m.topicConverter.NatsToMqtt(topic)
	logx.Infof("MQTT 发布主题转换: %s -> %s (QoS: %d, Retained: %v)", topic, mqttTopic, qos, retained)

	eventData := events.NewEventMsg(ctx, arg)
	mqttTopic = "$inner/" + mqttTopic
	return m.mqttCli.PublishRaw(mqttTopic, qos, retained, eventData)
}

func (m *mqttEvent) QueueSubscribe(subj, queue string, cb events.HandleFunc) (subscription, error) {
	// 将 NATS 格式的主题转换为 MQTT 格式
	mqttTopic := m.topicConverter.ConvertWildcards(subj)
	// MQTT 的共享订阅使用 $share/group/topic 格式
	mqttTopic = "$inner/" + mqttTopic
	sharedTopic := fmt.Sprintf("$share/%s/%s", queue, mqttTopic)
	logx.Infof("MQTT 队列订阅主题转换: %s -> %s (Queue: %s)", subj, sharedTopic, queue)

	sub, err := m.mqttCli.Subscribe(sharedTopic, cb)
	if err != nil {
		return nil, err
	}
	return &mqttSubscription{sub}, nil
}

// QueueSubscribeWithQoS 支持指定 QoS 级别的队列订阅
func (m *mqttEvent) QueueSubscribeWithQoS(subj, queue string, qos byte, cb events.HandleFunc) (subscription, error) {
	// 将 NATS 格式的主题转换为 MQTT 格式
	mqttTopic := m.topicConverter.ConvertWildcards(subj)
	// MQTT 的共享订阅使用 $share/group/topic 格式
	mqttTopic = "$inner/" + mqttTopic
	sharedTopic := fmt.Sprintf("$share/%s/%s", mqttTopic, queue)
	logx.Infof("MQTT 队列订阅主题转换: %s -> %s (Queue: %s, QoS: %d)", subj, sharedTopic, queue, qos)

	return m.SubscribeWithQoS(sharedTopic, qos, cb)
}

func (m *mqttSubscription) Unsubscribe() error {
	return m.MqttSubscription.Unsubscribe()
}
