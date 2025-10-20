package errors

const MediaError = 3000000

var (
	MediaCreateError       = NewCodeError(MediaError+1, "error.media.mediaCreateError")
	MediaUpdateError       = NewCodeError(MediaError+2, "error.media.mediaUpdateError")
	MediaNotfoundError     = NewCodeError(MediaError+3, "error.media.mediaNotFoundError")
	MediaActiveError       = NewCodeError(MediaError+4, "error.media.mediaActiveError")
	MediaPullCreateError   = NewCodeError(MediaError+5, "error.media.mediaPullCreateError")
	MediaStreamDeleteError = NewCodeError(MediaError+6, "error.media.mediaStreamDeleteError")
	MediaRecordNotFound    = NewCodeError(MediaError+7, "error.media.mediaRecordNotFound")
	MediaSipUpdateError    = NewCodeError(MediaError+8, "error.media.mediaSipUpdateError")
	MediaSipDevCreateError = NewCodeError(MediaError+9, "error.media.mediaSipDevCreateError")
	MediaSipChnCreateError = NewCodeError(MediaError+10, "error.media.mediaSipChnCreateError")
	MediaSipPlayError      = NewCodeError(MediaError+11, "error.media.mediaSipPlayError")
)
