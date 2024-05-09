package def

type NotifyCode = string

const (
	NotifyCodeSysUserLoginCaptcha    NotifyCode = "sysUserLoginCaptcha"
	NotifyCodeSysUserRegisterCaptcha NotifyCode = "sysUserRegisterCaptcha"

	NotifyCodeRuleScene   NotifyCode = "ruleScene"       //场景联动通知
	NotifyCodeDeviceAlarm NotifyCode = "ruleDeviceAlarm" //设备告警通知
)

const (
	NotifyGroupCaptcha = "captcha" //验证码通知
	NotifyGroupDevice  = "device"  //设备通知
	NotifyGroupSystem  = "system"  //系统通知
)

type NotifyType = string

const (
	NotifyTypeSms         NotifyType = "sms"
	NotifyTypeEmail       NotifyType = "email"
	NotifyTypeDingTalk    NotifyType = "dingTalk"
	NotifyTypeDingWebhook NotifyType = "dingWebhook"
	NotifyTypeWx          NotifyType = "wx"      //微信推送 todo
	NotifyTypeMessage     NotifyType = "message" //站内信通知
	NotifyTypePhoneCall   NotifyType = "phoneCall"
)
