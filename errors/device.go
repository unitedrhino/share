package errors

const DeviceError = 2000000

var (
	RespParam       = NewCodeError(DeviceError+1, "error.device.responseParamError")  // 返回参数不对
	DeviceTimeOut   = NewCodeError(DeviceError+2, "error.device.deviceTimeout")       // 设备回复超时
	NotOnline       = NewCodeError(DeviceError+3, "error.device.deviceOffline")       // 设备离线，请检查电源或设备
	DeviceResp      = NewCodeError(DeviceError+4, "error.device.deviceResponseError") // 设备回复错误
	DeviceBound     = NewCodeError(DeviceError+5, "error.device.deviceAlreadyBound")  // 设备已被绑定
	DeviceNotBound  = NewCodeError(DeviceError+6, "error.device.deviceNotBound")      // 设备未绑定
	DeviceCantBound = NewCodeError(DeviceError+7, "error.device.deviceCannotBound")   // 设备无法绑定
)
