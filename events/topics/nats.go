package topics

// 设备交互相关topic
const (
	DeviceUpMsg = "device.up.%s.%s.%s"
	DeviceUpAll = "device.up.>"
	//第一个参数是协议code,第二个参数是handle,第三和第四个参数是产品id和设备名
	DeviceDownMsg = "device.down.%s.%s.%s.%s"
	// DeviceDownAll dd模块订阅以下topic,收到内部的发布消息后向设备推送
	DeviceDownAll = "device.down.%s.>"

	// DeviceUpThing 物模型 最后两个是产品id和设备名称
	DeviceUpThing      = "device.up.thing.%s.%s"
	DeviceUpThingAll   = "device.up.thing.>"
	DeviceDownThing    = "device.down.thing.%s.%s"
	DeviceDownThingAll = "device.down.thing.>"
	// DeviceUpGateway 网关与子设备 最后两个是产品id和设备名称
	DeviceUpGateway      = "device.up.gateway.%s.%s"
	DeviceUpGatewayAll   = "device.up.gateway.>"
	DeviceDownGateway    = "device.down.gateway.%s.%s"
	DeviceDownGatewayAll = "device.down.gateway.>"

	// DeviceUpOta ota升级相关 最后两个是产品id和设备名称
	DeviceUpOta      = "device.up.ota.%s.%s"
	DeviceUpOtaAll   = "device.up.ota.>"
	DeviceDownOta    = "device.down.ota.%s.%s"
	DeviceDownOtaAll = "device.down.ota.>"
	// DeviceUpShadow 设备影子  最后两个是产品id和设备名称
	DeviceUpShadow      = "device.up.shadow.%s.%s"
	DeviceUpShadowAll   = "device.up.shadow.>"
	DeviceDownShadow    = "device.down.shadow.%s.%s"
	DeviceDownShadowAll = "device.down.shadow.>"
	// DeviceUpConfig 设备远程配置 最后两个是产品id和设备名称
	DeviceUpConfig      = "device.up.config.%s.%s"
	DeviceUpConfigAll   = "device.up.config.>"
	DeviceDownConfig    = "device.down.config.%s.%s"
	DeviceDownConfigAll = "device.down.config.>"
	// DeviceUpSDKLog 设备调试日志 最后两个是产品id和设备名称
	DeviceUpSDKLog      = "device.up.log.%s.%s"
	DeviceUpSDKLogAll   = "device.up.log.>"
	DeviceDownSdkLog    = "device.down.log.%s.%s"
	DeviceDownSDKLogAll = "device.down.log.>"

	// DeviceUpExt ext模块(包含ntp) 最后两个是产品id和设备名称
	DeviceUpExt    = "device.up.ext.%s.%s"
	DeviceUpExtAll = "device.up.ext.>"

	// DeviceUpStatusConnected 设备登录后向内部推送以下topic
	DeviceUpStatusConnected = "device.up.status.connected"
	// DeviceUpStatusDisconnected 设备的登出后向内部推送以下topic
	DeviceUpStatusDisconnected = "device.up.status.disconnected"
	DeviceUpStatus             = "device.up.status.>"

	// DeviceDownStatusConnected 设备在线状态修复,第一个参数是协议coe
	DeviceDownStatusConnected = "device.down.%s.status.fix"
)

// 应用事件通知(设备状态变化,设备上报)
const (
	// ApplicationDeviceStatusConnected 设备登录状态推送 中间两个是产品id和设备名称
	ApplicationDeviceStatusConnected = "application.device.%s.%s.status.connected"
	// ApplicationDeviceStatusDisConnected 设备登出状态推送 中间两个是产品id和设备名称
	ApplicationDeviceStatusDisConnected = "application.device.%s.%s.status.disconnected"
	// ApplicationDeviceReportThingProperty 设备物模型属性上报通知 中间两个是产品id和设备名称,最后一个是属性id
	ApplicationDeviceReportThingProperty   = "application.device.%s.%s.report.thing.property.%s"
	ApplicationDeviceReportThingPropertyV2 = "application.v2.device.%s.%s.report.thing.property"
	// ApplicationDeviceReportThingEvent 设备物模型事件上报通知 中间两个是产品id和设备名称,最后两个是事件类型和事件id
	ApplicationDeviceReportThingEvent = "application.device.%s.%s.report.thing.event.%s.%s"
	// ApplicationDeviceReportThingAction 设备物模型事件上报通知 中间两个是产品id和设备名称,最后三个是actionID,请求类型(req resp)和调用方向
	ApplicationDeviceReportThingAction = "application.device.%s.%s.report.thing.action.%s.%s.%s"
	// ApplicationDeviceReportThingPropertyDevice 设备物模型属性上报通知 中间两个是产品id和设备名称
	ApplicationDeviceReportThingPropertyDevice = "application.device.%s.%s.report.thing.property"

	ApplicationDeviceReportThingEventAllDevice    = "application.device.*.*.report.thing.event.>"
	ApplicationDeviceReportThingPropertyAllDevice = "application.device.*.*.report.thing.property.>"
	ApplicationDeviceStatusConnectedAllDevice     = "application.device.*.*.status.connected"
	ApplicationDeviceStatusDisConnectedAllDevice  = "application.device.*.*.status.disconnected"
	ApplicationDeviceStatusAllDevice              = "application.device.*.*.status.>"
)

// 服务自己的消息
const (
	TimedJobClean       = "server.core.timedjob.clean" //定时任务服务缓存及日志定时清理
	SysOssClean         = "server.core.sys.oss.clean"
	DmActionCheckDelay  = "server.things.dm.action.check.delay"
	VidInfoCheckStatus  = "server.video.vid.info.check.status"
	VidInfoInitDatabase = "server.video.vid.info.init.database"
)
