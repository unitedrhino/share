package eventBus

// 服务自己的消息
const (
	DmDeviceInfoDelete    = "server.things.dm.device.info.delete"
	DmProductInfoDelete   = "server.things.dm.product.info.delete"
	DmProductCustomUpdate = "server.things.dm.product.custom.update" //产品脚本有更新
	DmProductSchemaUpdate = "server.things.dm.product.schema.update" //物模型有更新

	// DmProtocolInfoUpdate 中间的是协议code
	DmProtocolInfoUpdate = "server.things.dm.protocol.%s.update" //自定义协议配置有更新
	UdRuleTiming         = "server.things.ud.rule.timing"
)
