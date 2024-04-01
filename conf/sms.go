package conf

const (
	SmsAli = "ali"
)

type Sms struct {
	Mode   string     `json:",default=ali,options=ali"`
	Enable bool       `json:",default=false"`
	Ali    SmsAliConf `json:",optional"`
}

type SmsAliConf struct {
	AccessKeyID     string `json:",env=aliAccessKeyID"`
	AccessKeySecret string `json:",env=aliAccessKeySecret"`
}
