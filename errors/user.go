package errors

const UserError = 1000000

var (
	DuplicateUsername    = NewCodeError(UserError+1, "error.user.usernameAlreadyRegistered")
	DuplicateMobile      = NewCodeError(UserError+2, "error.user.mobileAlreadyTaken")
	UnRegister           = NewCodeError(UserError+3, "error.user.notRegistered")
	Password             = NewCodeError(UserError+4, "error.user.accountOrPasswordError")
	Captcha              = NewCodeError(UserError+5, "error.user.captchaError")
	UidNotRight          = NewCodeError(UserError+6, "error.user.uidIncorrect")
	NotLogin             = NewCodeError(UserError+7, "error.user.notLoggedIn")
	NotSupportLogin      = NewCodeError(UserError+8, "error.user.loginMethodNotSupported")
	RegisterOne          = NewCodeError(UserError+22, "error.user.registrationStepOneFailed")
	DuplicateRegister    = NewCodeError(UserError+23, "error.user.duplicateRegistration")
	NeedUserName         = NewCodeError(UserError+24, "error.user.usernameRequired")
	PasswordLevel        = NewCodeError(UserError+25, "error.user.passwordStrengthInsufficient")
	GetInfoPartFailure   = NewCodeError(UserError+26, "error.user.getUserInfoPartialFailure")
	UsernameFormatErr    = NewCodeError(UserError+27, "error.user.usernameFormatError")
	AccountOrIpForbidden = NewCodeError(UserError+28, "error.user.accountOrIpForbidden")
	UseCaptcha           = NewCodeError(UserError+29, "error.user.useCaptcha")
	AccountDisable       = NewCodeError(UserError+30, "error.user.accountDisabled")
	BindAccount          = NewCodeError(UserError+31, "error.user.accountAlreadyBound")
	AccountKickedOut     = NewCodeError(UserError+32, "error.user.accountKickedOut")
	UnBindAccount        = NewCodeError(UserError+33, "error.user.accountNotBound")
	NeedImgCaptcha       = NewCodeError(UserError+34, "error.user.imageCaptchaRequired")
)
