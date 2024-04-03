package def

const (
	NotifyCodeSysUserLoginCaptcha    = "sysUserLoginCaptcha"
	NotifyCodeSysUserRegisterCaptcha = "sysUserRegisterCaptcha"

	NotifyCodeRuleScene   = "ruleScene"       //场景联动通知
	NotifyCodeDeviceAlarm = "ruleDeviceAlarm" //设备告警通知
)

const (
	NotifyGroupCaptcha = "captcha" //验证码通知
	NotifyGroupDevice  = "device"  //设备通知
)

type NotifyType = string

const (
	NotifyTypeSms      NotifyType = "sms"
	NotifyTypeEmail    NotifyType = "email"
	NotifyTypeDingTalk NotifyType = "dingTalk"
	NotifyTypeWx       NotifyType = "wx"      //微信推送 todo
	NotifyTypeMessage  NotifyType = "message" //站内信通知
)
