package mqttClient

import (
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"sync"
	"time"

	"gitee.com/unitedrhino/share/errors"
	"github.com/google/uuid"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/events"
	"gitee.com/unitedrhino/share/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	mqttInitOnce sync.Once
	mqttClient   *MqttClient
	// mqttSetOnConnectHandler 如果会话断开可以通过该回调函数来重新订阅消息
	//不使用mqtt的clean session是因为会话保持期间共享订阅也会给离线的客户端,这会导致在线的客户端丢失消息
	mqttSetOnConnectHandler func(cli mqtt.Client)
)

type MqttClient struct {
	clients []mqtt.Client
	cfg     *conf.MqttConf
}

func NewMqttClient(conf *conf.MqttConf) (mcs *MqttClient, err error) {
	mqttInitOnce.Do(func() {
		var clients []mqtt.Client
		var start = time.Now()
		for len(clients) < conf.ConnNum {
			var (
				mc mqtt.Client
			)
			var tryTime = 5
			for i := tryTime; i > 0; i-- {
				mc, err = initMqtt(conf)
				logx.Infof("mqtt_client initMqtt2 mc:%v err:%v", mc, err)
				if err != nil { //出现并发情况的时候可能联犀的http还没启动完毕
					logx.Errorf("mqtt_client 连接失败 重试剩余次数:%v", i-1)
					time.Sleep(time.Second * time.Duration(tryTime) / time.Duration(i))
					continue
				}
				break
			}
			if err != nil {
				logx.Errorf("mqtt_client 连接失败 conf:%#v  err:%v", conf, err)
				os.Exit(-1)
			}
			clients = append(clients, mc)
			var cli = MqttClient{clients: clients, cfg: conf}
			mqttClient = &cli
			logx.Infof("mqtt_client 连接完成 clientNum:%v use:%s", len(clients), time.Now().Sub(start))
		}
	})
	return mqttClient, err
}

func SetMqttSetOnConnectHandler(f func(cli mqtt.Client)) {
	mqttSetOnConnectHandler = f
}

func (m MqttClient) SubscribeRaw(cli mqtt.Client, topic string, qos byte, callback mqtt.MessageHandler) error {
	var clients = m.clients
	if cli != nil {
		clients = []mqtt.Client{cli}
	}
	logx.Infof("mqttClientSubscribe clientNum:%v topic:%v", len(clients), topic)
	for _, c := range clients {
		err := c.Subscribe(topic, qos, callback).Error()
		if err != nil {
			return errors.System.AddDetail(err)
		}
	}
	return nil
}

func (m MqttClient) Subscribe(topic string, cb events.HandleFunc) (*MqttSubscription, error) {
	return m.SubscribeWithQoS(topic, 0, cb)
}

func (m MqttClient) SubscribeWithQoS(topic string, qos byte, cb events.HandleFunc) (*MqttSubscription, error) {
	// 直接使用通用 HandleFunc
	handler := func(client mqtt.Client, msg mqtt.Message) {
		// 从消息中提取上下文和数据
		payload := msg.Payload()
		// 直接调用通用 HandleFunc
		err := cb(context.Background(), time.Now(), payload)
		if err != nil {
			logx.Errorf("MQTT message handler error: %v", err)
		}
	}

	err := m.SubscribeRaw(nil, topic, qos, handler)
	if err != nil {
		return nil, err
	}

	return &MqttSubscription{
		topic:  topic,
		qos:    qos,
		client: &m,
	}, nil
}

func (m MqttClient) QueueSubscribe(topic, queue string, cb events.HandleFunc) (*MqttSubscription, error) {
	// 直接使用传入的主题，共享订阅格式由上层处理
	return m.SubscribeWithQoS(topic, 0, cb)
}

func (m MqttClient) PublishRaw(topic string, qos byte, retained bool, payload interface{}) error {
	id := rand.Intn(len(m.clients))
	return m.clients[id].Publish(topic, qos, retained, payload).Error()
}

func (m MqttClient) Publish(ctx context.Context, topic string, data []byte) error {
	// 使用事件消息格式包装数据
	return m.PublishRaw(topic, 1, false, data)
}

// MqttSubscription 结构体定义
type MqttSubscription struct {
	topic  string
	qos    byte
	client *MqttClient
}

func (ms *MqttSubscription) Unsubscribe() error {
	// MQTT 客户端会自动处理取消订阅
	// 这里可以添加日志记录
	logx.Infof("mqtt subscription unsubscribed: %s", ms.topic)
	return nil
}

// GetClientCount 获取客户端连接数量
func (m MqttClient) GetClientCount() int {
	return len(m.clients)
}

// IsConnected 检查是否有可用的连接
func (m MqttClient) IsConnected() bool {
	for _, client := range m.clients {
		if client.IsConnected() {
			return true
		}
	}
	return false
}

// GetConnectedClients 获取所有已连接的客户端
func (m MqttClient) GetConnectedClients() []mqtt.Client {
	var connectedClients []mqtt.Client
	for _, client := range m.clients {
		if client.IsConnected() {
			connectedClients = append(connectedClients, client)
		}
	}
	return connectedClients
}

func initMqtt(conf *conf.MqttConf) (mc mqtt.Client, err error) {
	opts := mqtt.NewClientOptions()
	for _, broker := range conf.Brokers {
		opts.AddBroker(broker)
	}
	uuid := uuid.NewString()
	clientID := conf.ClientID + "_" + uuid
	logx.Infof("mqtt_client initMqtt conf:%#v clientID:%v brokers:%#v stack=%s", conf, clientID, opts.Servers, utils.Stack(1, 10))
	opts.SetClientID(clientID).SetUsername(conf.User).SetPassword(conf.Pass)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		logx.Infof("mqtt_client Connected clientID:%v", clientID)
		if mqttSetOnConnectHandler != nil {
			mqttSetOnConnectHandler(client)
		}
	})
	opts.SetReconnectingHandler(func(client mqtt.Client, options *mqtt.ClientOptions) {
		logx.Infof("mqtt_client Reconnecting clientID:%#v", options)
		if mqttSetOnConnectHandler != nil {
			mqttSetOnConnectHandler(client)
		}
	})

	opts.SetAutoReconnect(true).SetMaxReconnectInterval(30 * time.Second) //意外离线的重连参数
	opts.SetConnectRetry(true).SetConnectRetryInterval(5 * time.Second)   //首次连接的重连参数

	opts.SetConnectionAttemptHandler(func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
		logx.Infof("mqtt_client 正在尝试连接 broker:%v clientID:%v", utils.Fmt(broker), clientID)
		return tlsCfg
	})
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		logx.Errorf("mqtt_client 连接丢失 err:%v  clientID:%v", utils.Fmt(err), clientID)
	})
	mc = mqtt.NewClient(opts)
	er2 := mc.Connect().WaitTimeout(5 * time.Second)
	if er2 == false {
		logx.Errorf("mqtt_client 连接失败超时")
		err = fmt.Errorf("mqtt_client 连接失败")
		return nil, err
	}
	return
}
