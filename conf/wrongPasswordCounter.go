package conf

import (
	"github.com/spf13/cast"
)

type WrongPasswordCounter struct {
	Captcha int `json:",default=5"` // 错误密码次数
	Account []struct {
		Statistics    int `json:",default=1440"` //超时时间 单位秒
		TriggerTimes  int `json:",default=10"`   // 错误密码次数
		ForbiddenTime int `json:",default=10"`   //账号或ip冻结时间 单位:秒
	}
	Ip []struct {
		Timeout       int `json:",default=1440"` //超时时间
		TriggerTimes  int `json:",default=200"`  // 错误密码次数
		ForbiddenTime int `json:",default=60"`   //账号或ip冻结时间
	}
}
type LimitCount struct {
	Times int // 错误密码次数
}

type LoginSafeCtlInfo struct {
	Prefix    string // key前缀
	Key       string // redis key
	Timeout   int    // redis key 超时时间
	Times     int    // 错误密码次数
	Forbidden int    // 账号或ip冻结时间
}

func (counter WrongPasswordCounter) ParseWrongPassConf(userID string, ip string) []*LoginSafeCtlInfo {
	var res []*LoginSafeCtlInfo
	res = append(res, &LoginSafeCtlInfo{
		Prefix:  "login:wrongPassword:captcha:",
		Key:     "login:wrongPassword:captcha:" + userID,
		Timeout: 24 * 3600,
		Times:   counter.Captcha,
	})

	for i, v := range counter.Account {
		res = append(res, &LoginSafeCtlInfo{
			Prefix:    "login:wrongPassword:account:",
			Key:       "login:wrongPassword:account:" + cast.ToString(i+1) + ":" + userID,
			Timeout:   v.Statistics,
			Times:     v.TriggerTimes,
			Forbidden: v.ForbiddenTime,
		})
	}
	for i, v := range counter.Ip {
		res = append(res, &LoginSafeCtlInfo{
			Prefix:    "login:wrongPassword:ip:",
			Key:       "login:wrongPassword:ip:" + cast.ToString(i+1) + ":" + ip,
			Timeout:   v.Timeout,
			Times:     v.TriggerTimes,
			Forbidden: v.ForbiddenTime,
		})
	}

	return res
}
