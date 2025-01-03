package def

type NotifyCode = string

const (
	NotifyCodeSysUserLoginCaptcha     NotifyCode = "sysUserLoginCaptcha"
	NotifyCodeSysUserRegisterCaptcha  NotifyCode = "sysUserRegisterCaptcha"
	NotifyCodeSysUserChangePwdCaptcha NotifyCode = "sysUserChangePwdCaptcha"

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
	NotifyTypeSms         NotifyType = "sms"         //短信
	NotifyTypeEmail       NotifyType = "email"       //邮箱
	NotifyTypeDingTalk    NotifyType = "dingTalk"    //钉钉机器人
	NotifyTypeDingWebhook NotifyType = "dingWebhook" //钉钉webhook
	NotifyTypeWxMini      NotifyType = "wxMini"      //微信小程序推送
	NotifyTypeWxEWebhook  NotifyType = "wxEWebHook"  //企业微信webhook
	NotifyTypeMessage     NotifyType = "message"     //站内信通知
	NotifyTypePhoneCall   NotifyType = "phoneCall"   //电话通知
	NotifyTypeWxApp       NotifyType = "wxCorApp"    //企业微信app消息
)
