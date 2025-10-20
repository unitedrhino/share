package errors

const OtaError = 2100000

var (
	OtaRetryStatusError  = NewCodeError(OtaError+1, "error.ota.otaRetryStatusError")
	OtaCancleStatusError = NewCodeError(OtaError+2, "error.ota.otaCancelStatusError")
	OtaDeviceNumError    = NewCodeError(OtaError+3, "error.ota.otaDeviceNumError")
)
