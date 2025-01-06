package def

type UserSubscribe = string

var (
	UserSubscribeDevicePropertyReport  = "devicePropertyReport"   //设备上报订阅
	UserSubscribeDevicePropertyReport2 = "devicePropertyReportV2" //设备上报订阅
	UserSubscribeDevicePublish         = "devicePublish"          //设备发布消息
	UserSubscribeDeviceActionReport    = "deviceActionReport"     //设备行为消息
	UserSubscribeDeviceConn            = "deviceConn"             //设备连接消息
	UserSubscribeDeviceOtaReport       = "deviceOtaReport"        //设备ota消息推送
	UserSubscribeUserOrderPay          = "userOrderPay"           //用户订单支付通知

	UserSubscribeRuleDebugMsgReport = "ruleDebugMsgReport" //规则引擎调试消息推送

)
