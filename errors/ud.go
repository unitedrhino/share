package errors

const UdError = 4000000

var (
	TriggerType = NewCodeError(UdError+1, "error.ud.triggerTypeNotSupported") // 触发类型不支持
)
