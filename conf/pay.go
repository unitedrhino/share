package conf

type WechatPay struct {
	AppID          string
	MchID          string
	APIv3Key       string
	SerialNo       string
	PrivateKey     string `json:",optional"` //svc里将文件读取到这里
	PrivateKeyFile string `json:",default='etc/apiclient_key.pem'"`
	NotifyUrl      string
}
