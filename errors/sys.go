package errors

const SysError = 100000

var (
	OK               = NewCodeError(200, "error.sys.success")                         // 成功
	Default          = NewCodeError(SysError+1, "error.sys.otherError")               // 其他错误
	TokenExpired     = NewCodeError(SysError+2, "error.sys.tokenExpired")             // token已经过期
	TokenNotValidYet = NewCodeError(SysError+3, "error.sys.tokenNotValidYet")         // token还未生效
	TokenMalformed   = NewCodeError(SysError+4, "error.sys.tokenMalformed")           // token格式错误
	TokenInvalid     = NewCodeError(SysError+5, "error.sys.tokenInvalid")             // 违法的token
	Parameter        = NewCodeError(SysError+6, "error.sys.parameterError")           // 参数错误
	System           = NewCodeError(SysError+7, "error.sys.systemError")              // 系统错误
	Database         = NewCodeError(SysError+8, "error.sys.databaseError")            // 数据库错误
	NotFind          = NewCodeError(SysError+9, "error.sys.notFound")                 // 未查询到
	Duplicate        = NewCodeError(SysError+10, "error.sys.duplicateParameter")      // 参数重复
	SignatureExpired = NewCodeError(SysError+11, "error.sys.signatureExpired")        // 签名已经过期
	Permissions      = NewCodeError(SysError+12, "error.sys.insufficientPrivileges")  // 权限不足
	Method           = NewCodeError(SysError+13, "error.sys.methodNotSupported")      // method不支持
	Type             = NewCodeError(SysError+14, "error.sys.invalidParameterType")    // 参数的类型不对
	OutRange         = NewCodeError(SysError+15, "error.sys.parameterOutOfRange")     // 参数的值超出范围
	TimeOut          = NewCodeError(SysError+16, "error.sys.timeout")                 // 等待超时
	Server           = NewCodeError(SysError+17, "error.sys.serverError")             // 本实例处理不了该信息
	NotRealize       = NewCodeError(SysError+18, "error.sys.notImplemented")          // 尚未实现
	NotEmpty         = NewCodeError(SysError+19, "error.sys.notEmpty")                // 不为空
	Panic            = NewCodeError(SysError+20, "error.sys.systemPanic")             // 系统异常，请联系开发者
	NotEnable        = NewCodeError(SysError+21, "error.sys.notEnabled")              // 未启用
	Company          = NewCodeError(SysError+22, "error.sys.enterpriseFeature")       // 该功能是企业版功能
	Script           = NewCodeError(SysError+23, "error.sys.scriptExecutionFailed")   // 脚本执行失败
	OnGoing          = NewCodeError(SysError+24, "error.sys.inProgress")              // 正在执行中 - 事务分布式事务中如果返回该错误码,分布式事务会定时重试
	Failure          = NewCodeError(SysError+25, "error.sys.executionFailedRollback") // 执行失败,需要回滚 - 事务分布式事务中如果返回该错误码,分布式事务会进行回滚
	Jump             = NewCodeError(SysError+26, "error.sys.skipExecution")           // 跳过执行
	Limit            = NewCodeError(SysError+27, "error.sys.limit")                   // 已到达上限
)
