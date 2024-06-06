package eventBus

// 服务自己的消息
const (
	DmDeviceInfoDelete         = "server.things.dm.device.info.delete"
	DmDeviceOnlineStatusChange = "server.things.dm.device.onlineStatus.change"
	DmProductInfoDelete        = "server.things.dm.product.info.delete"
	DmProductCustomUpdate      = "server.things.dm.product.custom.update"   //产品脚本有更新
	DmOtaDeviceUpgradePush     = "server.things.dm.ota.device.upgrade.push" //ota设备推送
	DmOtaJobDelayRun           = "server.things.dm.ota.job.delay.run"       //任务延时启动
	// DmProtocolInfoUpdate 中间的是协议code
	DmProtocolInfoUpdate = "server.things.dm.protocol.%s.update" //自定义协议配置有更新
	UdRuleTimer          = "server.things.ud.rule.timer"
	SaleStaticTimer      = "server.sale.static.timer"
	PAliTimer            = "server.things.pali.data.timer"

	ServerCacheSync = "server.cache.sync.%s"

	CoreUserDelete    = "server.core.user.delete"
	CoreProjectDelete = "server.core.project.delete"

	CoreApiUserPublish = "server.core.api.user.publish.%v"
)

const (
	ServerCacheKeySysUserInfo          = "cache:sys:user:info"
	ServerCacheKeySysProjectInfo       = "cache:sys:project:info"
	ServerCacheKeySysTenantInfo        = "cache:sys:tenant:info"
	ServerCacheKeySysTenantOpenWebhook = "cache:sys:tenant:open:webhook"
	ServerCacheKeyDmProduct            = "cache:dm:product"
	ServerCacheKeyDmDevice             = "cache:dm:device"
	ServerCacheKeyDmSchema             = "cache:dm:schema"
)
