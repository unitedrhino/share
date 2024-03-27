package conf

type MqttConf struct {
	ClientID string      //在mqtt中的clientID
	Brokers  []string    //mqtt服务器节点
	User     string      `json:",default=root"` //用户名
	Pass     string      `json:",optional"`     //密码
	ConnNum  int         `json:",default=10"`   //默认连接数
	OpenApi  *DevOpenApi `json:",optional"`     //用来检查设备状态
}

// https://www.emqx.io/docs/zh/latest/admin/api.html#go
type DevOpenApi struct {
	Host      string
	ApiKey    string
	SecretKey string
}
