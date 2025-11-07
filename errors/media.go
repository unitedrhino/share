package errors

const MediaError = 3000000

var (
	MediaCreateError       = NewCodeError(MediaError+1, "error.media.mediaCreateError")        // 流服务创建失败
	MediaUpdateError       = NewCodeError(MediaError+2, "error.media.mediaUpdateError")        // 流服务更新失败
	MediaNotfoundError     = NewCodeError(MediaError+3, "error.media.mediaNotFoundError")      // 流服务不存在
	MediaActiveError       = NewCodeError(MediaError+4, "error.media.mediaActiveError")        // 流服务激活失败
	MediaPullCreateError   = NewCodeError(MediaError+5, "error.media.mediaPullCreateError")    // 拉流创建失败
	MediaStreamDeleteError = NewCodeError(MediaError+6, "error.media.mediaStreamDeleteError")  // 流删除错误
	MediaRecordNotFound    = NewCodeError(MediaError+7, "error.media.mediaRecordNotFound")     // 未找到录像列表
	MediaSipUpdateError    = NewCodeError(MediaError+8, "error.media.mediaSipUpdateError")     // ID或channelID不能都为空
	MediaSipDevCreateError = NewCodeError(MediaError+9, "error.media.mediaSipDevCreateError")  // 设备创建失败
	MediaSipChnCreateError = NewCodeError(MediaError+10, "error.media.mediaSipChnCreateError") // 通道创建失败
	MediaSipPlayError      = NewCodeError(MediaError+11, "error.media.mediaSipPlayError")      // 通道播放失败
)
