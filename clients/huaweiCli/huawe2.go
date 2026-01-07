package huaweiCli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

// HuaweiQuickLoginMobileResult 华为一键登录获取手机号结果

type HuaweiQuickLoginMobileResult struct {
	OpenID           string `json:"openId"`           // 用户OpenID
	UnionID          string `json:"unionId"`          // 用户UnionID
	PhoneNumber      string `json:"phoneNumber"`      // 华为账号绑定号码（含国际冠码）
	PhoneNumberValid int    `json:"phoneNumberValid"` // 手机号实时有效性：0需验证，1可直接使用
	PurePhoneNumber  string `json:"purePhoneNumber"`  // 不带国际冠码与国际电话区号的手机号
	PhoneCountryCode string `json:"phoneCountryCode"` // 国际电话区号
	ResultCode       int    `json:"resultCode"`       // 结果码
	ResultDesc       string `json:"resultDesc"`       // 结果描述
}

// GetQuickLoginMobileByCode 一键登录获取华为账号绑定号码和OpenID/UnionID
func (h *HuaweiClient) GetQuickLoginMobileByCode(ctx context.Context, authorizationCode string) (*HuaweiQuickLoginMobileResult, error) {
	// 构建请求参数
	reqBody := map[string]string{
		"code":         authorizationCode,
		"clientId":     h.config.ClientID,
		"clientSecret": h.config.ClientSecret,
	}
	url := "https://account-api.cloud.huawei.com/oauth2/v6/quickLogin/getPhoneNumber"
	// 序列化为JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		logx.WithContext(ctx).Errorf("JSON序列化失败: %v", err)
		return nil, fmt.Errorf("json marshal failed: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logx.WithContext(ctx).Errorf("创建HTTP请求失败: %v", err)
		return nil, fmt.Errorf("create http request failed: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logx.WithContext(ctx).Errorf("发送HTTP请求失败: %v", err)
		return nil, fmt.Errorf("send http request failed: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.WithContext(ctx).Errorf("读取响应失败: %v", err)
		return nil, fmt.Errorf("read response failed: %v", err)
	}
	logx.WithContext(ctx).Infof("华为响应体: %s", strings.NewReplacer("\n", "", "\r\n", "", "\r", "").Replace(string(body)))

	// 解析响应
	var result HuaweiQuickLoginMobileResult
	if err := json.Unmarshal(body, &result); err != nil {
		logx.WithContext(ctx).Errorf("JSON反序列化失败: %v, response: %s", err, string(body))
		return nil, fmt.Errorf("json unmarshal failed: %v", err)
	}

	return &result, nil
}
