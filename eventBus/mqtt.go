package eventBus

//
//import (
//	"context"
//	"gitee.com/unitedrhino/share/clients"
//	"gitee.com/unitedrhino/share/clients/mqttClient"
//	"gitee.com/unitedrhino/share/conf"
//	"gitee.com/unitedrhino/share/events"
//	"github.com/mqtt-io/mqtt.go"
//)
//
//type mqttEvent struct {
//	mqttCli *mqttClient.MqttClient
//}
//type mqttSubscription struct {
//}
//
//func NewMqttEvent(c conf.EventConf, serverName string, nodeID int64) (*mqttEvent, error) {
//	cli, err := mqttClient.NewMqttClient(&c.Mqtt)
//	if err != nil {
//		return nil, err
//	}
//	return &mqttEvent{cli}, nil
//}
//
//func (n *mqttEvent) Subscribe(subj string, cb events.HandleFunc) (*mqttSubscription, error) {
//	sub, err := n.mqttCli.S(subj, cb)
//	if err != nil {
//		return nil, err
//	}
//	return &mqttSubscription{sub}, nil
//}
//
//func (n *mqttEvent) Publish(ctx context.Context, topic string, arg []byte) error {
//	return n.mqttCli.Publish(ctx, topic, arg)
//}
//
//func (n *mqttEvent) QueueSubscribe(subj, queue string, cb events.HandleFunc) (*mqttSubscription, error) {
//	sub, err := n.mqttCli.QueueSubscribe(subj, queue, cb)
//	if err != nil {
//		return nil, err
//	}
//	return &mqttSubscription{sub}, nil
//}
//
//func (n *mqttSubscription) Unsubscribe() error {
//	return n.Unsubscribe()
//}
