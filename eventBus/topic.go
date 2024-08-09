package eventBus

// 服务自己的消息
const (
	SysProjectInfoDelete = "server.core.sys.project.info.delete"
	SysAreaInfoDelete    = "server.core.sys.area.info.delete"

	DmDeviceInfoUnbind         = "server.things.dm.device.info.unbind"
	DmDeviceInfoDelete         = "server.things.dm.device.info.delete"
	DmDeviceOnlineStatusChange = "server.things.dm.device.onlineStatus.change"
	DmDeviceStaticHalfHour     = "server.things.dm.device.static.halfHour" //半小时统计
	DmProductInfoDelete        = "server.things.dm.product.info.delete"
	DmProductCustomUpdate      = "server.things.dm.product.custom.update"   //产品脚本有更新
	DmOtaDeviceUpgradePush     = "server.things.dm.ota.device.upgrade.push" //ota设备推送
	DmOtaJobDelayRun           = "server.things.dm.ota.job.delay.run"       //任务延时启动
	// DmProtocolInfoUpdate 中间的是协议code
	DmProtocolInfoUpdate  = "server.things.dm.protocol.%s.update" //自定义协议配置有更新
	UdRuleTimer           = "server.things.ud.rule.timer"
	UdRuleTimerTenMinutes = "server.things.ud.rule.timer.tenMinutes"
	DgOnlineTimer         = "server.things.dg.online.timer"

	SaleStaticTimer   = "server.sale.static.timer"
	SalePayCheckTimer = "server.sale.pay.check.timer"

	PAliTimer = "server.things.pali.data.timer"

	ServerCacheSync = "server.cache.sync.%s"

	CoreUserDelete    = "server.core.user.delete"
	CoreProjectDelete = "server.core.project.delete"

	CoreApiUserPublish = "server.core.api.user.publish.%v"
)

const (
	ServerCacheKeySysUserInfo          = "cache:sys:user:info"
	ServerCacheKeySysUserTokenInfo     = "cache:sys:userToken:info"
	ServerCacheKeySysProjectInfo       = "cache:sys:project:info"
	ServerCacheKeySysAreaInfo          = "cache:sys:area:info"
	ServerCacheKeySysTenantInfo        = "cache:sys:tenant:info"
	ServerCacheKeySysTenantConfig      = "cache:sys:tenant:config"
	ServerCacheKeySysTenantOpenWebhook = "cache:sys:tenant:open:webhook"
	ServerCacheKeyDmProduct            = "cache:dm:product"
	ServerCacheKeyDmDevice             = "cache:dm:device"
	ServerCacheKeyDmSchema             = "cache:dm:schema"
	ServerCacheKeyDmUserShareDevice    = "cache:dm:user:share:device"
)
