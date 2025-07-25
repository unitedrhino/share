package tools

import (
	"fmt"
	"gitee.com/unitedrhino/share/utils"
	"github.com/maypok86/otter"
	"github.com/parnurzeal/gorequest"
	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

var cityIpCache otter.Cache[string, string]

func init() {
	cache, err := otter.MustBuilder[string, string](10_000).
		CollectStats().
		Cost(func(key string, value string) uint32 {
			return 1
		}).
		WithTTL(time.Hour * 24).
		Build()
	logx.Must(err)
	cityIpCache = cache
}

// GetCityByIp 获取ip所属城市
func GetCityByIp(ip string) string {
	if ip == "" {
		return ""
	}
	if ip == "[::1]" || ip == "127.0.0.1" {
		return "内网IP"
	}
	v, ok := cityIpCache.Get(ip)
	if ok {
		return v
	}
	url := "http://whois.pconline.com.cn/ipJson.jsp?json=true&ip=" + ip
	r := gorequest.New().Retry(1, time.Second*2)
	_, bytes, _ := r.Get(url).EndBytes()
	if len(bytes) == 0 {
		return ""
	}
	tmp, _ := utils.GBKToUTF8(bytes)
	json := gjson.Parse(string(tmp))
	if json.Get("code").Int() == 0 {
		city := fmt.Sprintf("%s %s", json.Get("pro").String(), json.Get("city").String())
		cityIpCache.Set(ip, city)
		return city
	} else {
		return ""
	}
}
