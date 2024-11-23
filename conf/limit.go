package conf

type Limit struct {
	Timeout       int //超时时间:单位秒
	TriggerTime   int // 错误密码次数
	ForbiddenTime int //账号或ip冻结时间
}
