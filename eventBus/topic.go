package eventBus

// 服务自己的消息
const (
	CoreProjectInfoDelete  = "server.core.project.info.delete"
	CoreAreaInfoDelete     = "server.core.area.info.delete"
	CoreOpsWorkOrderFinish = "server.core.ops.workOrder.finish"
	CoreSyncHalfHour       = "server.core.sync.halfHour" //半小时统计
	CoreUserDelete         = "server.core.user.delete"
	CoreUserCreate         = "server.core.user.create"
	CoreUserUpdate         = "server.core.user.update"
	CoreProjectDelete      = "server.core.project.delete"
	CoreApiUserPublish     = "server.core.api.user.publish.%v"

	DmDeviceInfoUnbind         = "server.things.dm.device.info.unbind"
	DmDeviceInfoDelete         = "server.things.dm.device.info.delete"
	DmDeviceOnlineStatusChange = "server.things.dm.device.onlineStatus.change"
	DmDeviceStaticOneHour      = "server.things.dm.device.static.2Hour"     //2小时统计
	DmDeviceStaticHalfHour     = "server.things.dm.device.static.halfHour"  //半小时统计
	DmDeviceStaticOneMinute    = "server.things.dm.device.static.oneMinute" //1分钟统计
	DmProductInfoDelete        = "server.things.dm.product.info.delete"
	DmProductCustomUpdate      = "server.things.dm.product.custom.update"   //产品脚本有更新
	DmOtaDeviceUpgradePush     = "server.things.dm.ota.device.upgrade.push" //ota设备推送
	DmOtaJobDelayRun           = "server.things.dm.ota.job.delay.run"       //任务延时启动
	// DmProtocolInfoUpdate 中间的是协议code
	DmProtocolInfoUpdate  = "server.things.dm.protocol.%s.update" //自定义协议配置有更新
	UdRuleTimer           = "server.things.ud.rule.timer"
	UdRuleTimerTenMinutes = "server.things.ud.rule.timer.tenMinutes"

	//最后一个参数是告警模式
	UdRuleAlarmNotify = "server.things.ud.rule.alarm.%s" //trigger:触发告警 relieve:解除告警

	DgOnlineTimer = "server.things.dg.online.timer"

	SaleStaticTimer   = "server.sale.static.timer"
	SalePayCheckTimer = "server.sale.pay.check.timer"

	PAliTimer = "server.things.pali.data.timer"

	ServerCacheSync = "server.cache.sync.%s"
)

const (
	ServerCacheKeySysUserInfo          = "cache:sys:user:info"
	ServerCacheKeySysUserTokenInfo     = "cache:sys:userToken:info"
	ServerCacheKeySysProjectInfo       = "cache:sys:project:info"
	ServerCacheKeySysAccessApi         = "cache:sys:access:api"
	ServerCacheKeySysRoleAccess        = "cache:sys:role:access"
	ServerCacheKeySysAreaInfo          = "cache:sys:area:info"
	ServerCacheKeySysTenantInfo        = "cache:sys:tenant:info"
	ServerCacheKeySysTenantConfig      = "cache:sys:tenant:config"
	ServerCacheKeySysTenantOpenWebhook = "cache:sys:tenant:open:webhook"
	ServerCacheKeyDmProduct            = "cache:dm:product"
	ServerCacheKeyDmDevice             = "cache:dm:device"
	ServerCacheKeyDmProductSchema      = "cache:dm:product:schema"
	ServerCacheKeyDmDeviceSchema       = "cache:dm:device:schema"
	ServerCacheKeyDmUserShareDevice    = "cache:dm:user:share:device"
	ServerCacheKeyDmMultiDevicesShare  = "cache:dm:user:multishare:devices"
)
