package clients

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/events"
	"gitee.com/i-Things/share/utils"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

type NatsClient struct {
	mode         string
	conf         conf.NatsConf
	nc           *nats.Conn
	js           nats.JetStreamContext
	consumerName string
}

func NewNatsClient2(mode string, ConsumerName string, natsConf conf.NatsConf, nodeID int64) (*NatsClient, error) {
	consumerName := fmt.Sprintf("%s-%d", ConsumerName, nodeID)
	client := NatsClient{
		conf:         natsConf,
		mode:         mode,
		consumerName: consumerName,
	}
	if mode == conf.EventModeNatsJs {
		js, err := NewNatsJetStreamClient(natsConf)
		if err != nil {
			return nil, err
		}
		client.js = js
		return &NatsClient{
			conf:         natsConf,
			js:           js,
			mode:         mode,
			consumerName: consumerName,
		}, nil
	} else {
		nc, err := NewNatsClient(natsConf)
		if err != nil {
			return nil, err
		}
		client.nc = nc
	}
	return &client, nil
}

func (n *NatsClient) QueueSubscribe(subj, queue string, cb events.HandleFunc) (*nats.Subscription, error) {
	if n.mode == conf.EventModeNatsJs {
		return n.js.QueueSubscribe(subj, queue, events.NatsSubscription(cb), nats.Durable("queue_"+events.GenNatsJsDurable(n.consumerName, subj)))
	}
	return n.nc.QueueSubscribe(subj, queue, events.NatsSubscription(cb))
}

func (n *NatsClient) Subscribe(subj string, cb events.HandleFunc) (*nats.Subscription, error) {
	return n.SubscribeWithConsumer(subj, n.consumerName, cb)
}

func (n *NatsClient) SubscribeWithConsumer(subj string, consumer string, cb events.HandleFunc) (*nats.Subscription, error) {
	if n.mode == conf.EventModeNatsJs {
		return n.js.Subscribe(subj, events.NatsSubscription(cb), nats.Durable("normal_"+events.GenNatsJsDurable(consumer, subj)))
	}
	return n.nc.Subscribe(subj, events.NatsSubscription(cb))
}

func (n *NatsClient) SubscribeSync(subj string) (*nats.Subscription, error) {
	if n.mode == conf.EventModeNatsJs {
		subscription, err := n.js.SubscribeSync(subj, nats.Durable("sync_"+events.GenNatsJsDurable(n.consumerName, subj+"--"+uuid.NewString())))
		return subscription, err
	}
	subscription, err := n.nc.SubscribeSync(subj)
	return subscription, err
}

func (n *NatsClient) Publish(ctx context.Context, subj string, data []byte) error {
	pubMsg := events.NewEventMsg(ctx, data)
	if n.mode == conf.EventModeNatsJs {
		_, err := n.js.Publish(subj, pubMsg)
		return err
	}
	err := n.nc.Publish(subj, pubMsg)
	return err
}

func NewNatsClient(conf conf.NatsConf) (nc *nats.Conn, err error) {
	connectOpts := nats.GetDefaultOptions()
	connectOpts.Url = conf.Url
	connectOpts.User = conf.User
	connectOpts.Password = conf.Pass
	connectOpts.Token = conf.Token
	connectOpts.Timeout = 5 * time.Second
	connectOpts.DisconnectedErrCB = func(conn *nats.Conn, err error) {
		logx.Errorf("nats.DisconnectedErrCB  err:%v", err)
	}
	connectOpts.AsyncErrorCB = func(conn *nats.Conn, subscription *nats.Subscription, err error) {
		logx.Errorf("nats.AsyncErrorCB subscription:%v err:%v", utils.Fmt(subscription), err)
	}

	nc, err = connectOpts.Connect()
	if err != nil {
		return
	}
	return nc, err
}

func NewNatsJetStreamClient(conf conf.NatsConf) (nats.JetStreamContext, error) {
	nc, err := NewNatsClient(conf)
	if err != nil {
		return nil, err
	}
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	return js, JetStreamInit(js)
}

func CreateStream(jetStream nats.JetStreamContext, name string, subjects []string) error {
	stream, err := jetStream.StreamInfo(name)
	// stream not found, create it
	if stream == nil {
		_, err = jetStream.AddStream(&nats.StreamConfig{
			Name:      name,
			Subjects:  subjects,
			Retention: nats.InterestPolicy,
			Discard:   nats.DiscardOld,
			MaxAge:    2 * time.Minute,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

var jetStreamInitOnce sync.Once

func JetStreamInit(jetStream nats.JetStreamContext) (err error) {
	jetStreamInitOnce.Do(func() {
		err = CreateStream(jetStream, "server", []string{"server.>"})
		if err != nil {
			return
		}
		err = CreateStream(jetStream, "device", []string{"device.>"})
		if err != nil {
			return
		}
		err = CreateStream(jetStream, "application", []string{"application.>"})
		if err != nil {
			return
		}
	})

	return
}

func CreateConsumer(jetStream nats.JetStreamContext, stream, name string) error {
	ConsumerInfo, err := jetStream.ConsumerInfo(stream, name)
	// stream not found, create it
	if ConsumerInfo == nil {
		_, err = jetStream.AddConsumer(stream, &nats.ConsumerConfig{
			Durable:        name,
			AckPolicy:      nats.AckAllPolicy,
			DeliverSubject: nats.NewInbox(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
