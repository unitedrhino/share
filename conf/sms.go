package conf

type SmsType = string

const (
	SmsAli     SmsType = "ali"
	SmsTencent SmsType = "tencent"
)

type Sms struct {
	Mode    string         `json:",default=ali,env=smsMode,options=ali|tencent"`
	Enable  bool           `json:",default=false"`
	Ali     SmsAliConf     `json:",optional"`
	Tencent SmsTencentConf `json:",optional"`
}

type SmsAliConf struct {
	AccessKeyID     string `json:",env=aliAccessKeyID"`
	AccessKeySecret string `json:",env=aliAccessKeySecret"`
}

type SmsTencentConf struct {
	AccessKeyID     string `json:",env=tencentAccessKeyID"`
	AccessKeySecret string `json:",env=tencentAccessKeySecret"`
	AppID           string `json:",env=tencentAppID"`
	AppKey          string `json:",env=tencentAppKey"`
}
