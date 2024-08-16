package tools

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gcharset"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/maypok86/otter"
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
	bytes := g.Client().GetBytes(context.TODO(), url)
	src := string(bytes)
	srcCharset := "GBK"
	tmp, _ := gcharset.ToUTF8(srcCharset, src)
	json, err := gjson.DecodeToJson(tmp)
	if err != nil {
		return ""
	}
	if json.Get("code").Int() == 0 {
		city := fmt.Sprintf("%s %s", json.Get("pro").String(), json.Get("city").String())
		cityIpCache.Set(ip, city)
		return city
	} else {
		return ""
	}
}
