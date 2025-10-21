package errors

const UserError = 1000000

var (
	DuplicateUsername    = NewCodeError(UserError+1, "error.user.usernameAlreadyRegistered")     // 用户名已经注册
	DuplicateMobile      = NewCodeError(UserError+2, "error.user.mobileAlreadyTaken")            // 手机号已经被占用
	UnRegister           = NewCodeError(UserError+3, "error.user.notRegistered")                 // 未注册
	Password             = NewCodeError(UserError+4, "error.user.accountOrPasswordError")        // 账号或密码错误
	Captcha              = NewCodeError(UserError+5, "error.user.captchaError")                  // 验证码错误
	UidNotRight          = NewCodeError(UserError+6, "error.user.uidIncorrect")                  // uid不对
	NotLogin             = NewCodeError(UserError+7, "error.user.notLoggedIn")                   // 尚未登录
	NotSupportLogin      = NewCodeError(UserError+8, "error.user.loginMethodNotSupported")       // 不支持的登录方式
	RegisterOne          = NewCodeError(UserError+22, "error.user.registrationStepOneFailed")    // 注册第一步未成功
	DuplicateRegister    = NewCodeError(UserError+23, "error.user.duplicateRegistration")        // 重复注册
	NeedUserName         = NewCodeError(UserError+24, "error.user.usernameRequired")             // 需要填入用户名
	PasswordLevel        = NewCodeError(UserError+25, "error.user.passwordStrengthInsufficient") // 密码强度不够
	GetInfoPartFailure   = NewCodeError(UserError+26, "error.user.getUserInfoPartialFailure")    // 获取用户信息有失败
	UsernameFormatErr    = NewCodeError(UserError+27, "error.user.usernameFormatError")          // 账号必须以大小写字母开头，且账号只能包含大小写字母，数字，下划线和减号。 长度为6到20位之间
	AccountOrIpForbidden = NewCodeError(UserError+28, "error.user.accountOrIpForbidden")         // 账号冻结
	UseCaptcha           = NewCodeError(UserError+29, "error.user.useCaptcha")                   // 账号或密码错误
	AccountDisable       = NewCodeError(UserError+30, "error.user.accountDisabled")              // 账号已禁用
	BindAccount          = NewCodeError(UserError+31, "error.user.accountAlreadyBound")          // 账号已绑定
	AccountKickedOut     = NewCodeError(UserError+32, "error.user.accountKickedOut")             // 账号被顶出
	UnBindAccount        = NewCodeError(UserError+33, "error.user.accountNotBound")              // 账号未绑定
	NeedImgCaptcha       = NewCodeError(UserError+34, "error.user.imageCaptchaRequired")         // 请输入图形验证码
)
