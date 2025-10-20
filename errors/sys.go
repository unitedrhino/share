package errors

const SysError = 100000

var (
	OK               = NewCodeError(200, "error.sys.success")
	Default          = NewCodeError(SysError+1, "error.sys.otherError")
	TokenExpired     = NewCodeError(SysError+2, "error.sys.tokenExpired")
	TokenNotValidYet = NewCodeError(SysError+3, "error.sys.tokenNotValidYet")
	TokenMalformed   = NewCodeError(SysError+4, "error.sys.tokenMalformed")
	TokenInvalid     = NewCodeError(SysError+5, "error.sys.tokenInvalid")
	Parameter        = NewCodeError(SysError+6, "error.sys.parameterError")
	System           = NewCodeError(SysError+7, "error.sys.systemError")
	Database         = NewCodeError(SysError+8, "error.sys.databaseError")
	NotFind          = NewCodeError(SysError+9, "error.sys.notFound")
	Duplicate        = NewCodeError(SysError+10, "error.sys.duplicateParameter")
	SignatureExpired = NewCodeError(SysError+11, "error.sys.signatureExpired")
	Permissions      = NewCodeError(SysError+12, "error.sys.insufficientPrivileges")
	Method           = NewCodeError(SysError+13, "error.sys.methodNotSupported")
	Type             = NewCodeError(SysError+14, "error.sys.invalidParameterType")
	OutRange         = NewCodeError(SysError+15, "error.sys.parameterOutOfRange")
	TimeOut          = NewCodeError(SysError+16, "error.sys.timeout")
	Server           = NewCodeError(SysError+17, "error.sys.serverError")
	NotRealize       = NewCodeError(SysError+18, "error.sys.notImplemented")
	NotEmpty         = NewCodeError(SysError+19, "error.sys.notEmpty")
	Panic            = NewCodeError(SysError+20, "error.sys.systemPanic")
	NotEnable        = NewCodeError(SysError+21, "error.sys.notEnabled")
	Company          = NewCodeError(SysError+22, "error.sys.enterpriseFeature")
	Script           = NewCodeError(SysError+23, "error.sys.scriptExecutionFailed")
	OnGoing          = NewCodeError(SysError+24, "error.sys.inProgress")              //事务分布式事务中如果返回该错误码,分布式事务会定时重试
	Failure          = NewCodeError(SysError+25, "error.sys.executionFailedRollback") //事务分布式事务中如果返回该错误码,分布式事务会进行回滚
	Jump             = NewCodeError(SysError+26, "error.sys.skipExecution")
)
