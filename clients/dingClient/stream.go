package dingClient

import (
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
)

type StreamClient = client.StreamClient

// stream 模式钉钉只会随机推一个,所以可以都监听
func NewDingStream(clientId, clientSecret string) *StreamClient {
	cli := client.NewStreamClient(client.WithAppCredential(client.NewAppCredentialConfig(clientId, clientSecret)))
	//cli.RegisterAllEventRouter(OnEventReceived)

	return cli
}

//func OnEventReceived(_ context.Context, df *payload.DataFrame) (*payload.DataFrameResponse, error) {
//	eventHeader := event.NewEventHeaderFromDataFrame(df)
//	if eventHeader.EventType != "chat_update_title" {
//		// ignore events not equals `chat_update_title`; 忽略`chat_update_title`之外的其他事件；
//		// 该示例仅演示 chat_update_title 类型的事件订阅；
//		return event.NewSuccessResponse()
//	}
//
//	logger.GetLogger().Infof("received event, delay=%s, eventType=%s, eventId=%s, eventBornTime=%d, eventCorpId=%s, eventUnifiedAppId=%s, data=%s",
//		time.Duration(time.Now().UnixMilli()-eventHeader.EventBornTime)*time.Millisecond,
//		eventHeader.EventType,
//		eventHeader.EventId,
//		eventHeader.EventBornTime,
//		eventHeader.EventCorpId,
//		eventHeader.EventUnifiedAppId,
//		df.Data)
//	// put your code here; 可以在这里添加你的业务代码，处理事件订阅的业务逻辑；
//
//	return event.NewSuccessResponse()
//}
