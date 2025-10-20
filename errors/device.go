package errors

const DeviceError = 2000000

var (
	RespParam       = NewCodeError(DeviceError+1, "error.device.responseParamError")
	DeviceTimeOut   = NewCodeError(DeviceError+2, "error.device.deviceTimeout")
	NotOnline       = NewCodeError(DeviceError+3, "error.device.deviceOffline")
	DeviceResp      = NewCodeError(DeviceError+4, "error.device.deviceResponseError")
	DeviceBound     = NewCodeError(DeviceError+5, "error.device.deviceAlreadyBound")
	DeviceNotBound  = NewCodeError(DeviceError+6, "error.device.deviceNotBound")
	DeviceCantBound = NewCodeError(DeviceError+7, "error.device.deviceCannotBound")
)
