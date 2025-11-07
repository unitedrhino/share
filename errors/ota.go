package errors

const OtaError = 2100000

var (
	OtaRetryStatusError  = NewCodeError(OtaError+1, "error.ota.otaRetryStatusError")  // 升级状态不允许重新升级
	OtaCancleStatusError = NewCodeError(OtaError+2, "error.ota.otaCancelStatusError") // 升级状态已结束
	OtaDeviceNumError    = NewCodeError(OtaError+3, "error.ota.otaDeviceNumError")    // 验证设备数不能超过10个
)
