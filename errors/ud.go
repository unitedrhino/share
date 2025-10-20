package errors

const UdError = 4000000

var (
	TriggerType = NewCodeError(UdError+1, "error.ud.triggerTypeNotSupported")
)
