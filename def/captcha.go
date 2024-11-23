package def

type Captcha = string

const (
	CaptchaTypePhone = "phone"
	CaptchaTypeImage = "image" //图形验证码
	CaptchaTypeEmail = "email"
)

const (
	CaptchaUseLogin       = "login"
	CaptchaUseRegister    = "register"
	CaptchaUseChangePwd   = "changePwd"
	CaptchaUseBindAccount = "bindAccount"
	CaptchaUseForgetPwd   = "forgetPwd"
)

const (
	CaptchaExpire = 180
)
