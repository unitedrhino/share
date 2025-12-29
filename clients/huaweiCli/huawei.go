// 华为登录SDK实现
// 参考文档:
// 1. 华为账号UnionID登录指南: https://developer.huawei.com/consumer/cn/doc/harmonyos-guides/account-phone-unionid-login
// 2. 快速登录API: https://developer.huawei.com/consumer/cn/doc/harmonyos-references/account-api-get-user-info-quicklogin-by-code#section2520125725115
// 3. 手机号获取API: https://developer.huawei.com/consumer/cn/doc/harmonyos-references/account-api-get-user-info-get-phone
// 4. 昵称和头像获取API: https://developer.huawei.com/consumer/cn/doc/harmonyos-references/account-api-get-user-info-get-nickname-and-avatar
package huaweiCli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"gitee.com/unitedrhino/share/conf"
	"golang.org/x/oauth2"
)

// 华为OAuth2端点配置
var huaweiEndpoint = oauth2.Endpoint{
	AuthURL:  "https://oauth-login.cloud.huawei.com/oauth2/v3/authorize",
	TokenURL: "https://oauth-login.cloud.huawei.com/oauth2/v3/token",
}

// HuaweiClient 华为登录客户端
type HuaweiClient struct {
	config *oauth2.Config
}

// NewHuaweiClient 创建华为登录客户端实例
func NewHuaweiClient(ctx context.Context, conf *conf.ThirdConf) *HuaweiClient {
	if conf == nil {
		return nil
	}
	// 配置OAuth2
	oauthConfig := &oauth2.Config{
		ClientID:     conf.AppID,
		ClientSecret: conf.AppSecret,
		Endpoint:     huaweiEndpoint,
		Scopes: []string{
			"openid",     // 请求OpenID
			"profile",    // 请求昵称和头像
			"phone.read", // 请求手机号
			"https://www.huawei.com/auth/account/base.profile", // 请求基础用户信息
		},
	}
	return &HuaweiClient{
		config: oauthConfig,
	}
}

// GetAuthCodeURL 获取授权码URL
func (h *HuaweiClient) GetAuthCodeURL(redirectURL, state string) string {
	h.config.RedirectURL = redirectURL
	return h.config.AuthCodeURL(state)
}

// ExchangeToken 使用授权码获取访问令牌
func (h *HuaweiClient) ExchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return h.config.Exchange(ctx, code)
}

// GetUserInfo 获取用户信息（昵称、头像、UnionID、OpenID）
func (h *HuaweiClient) GetUserInfo(ctx context.Context, accessToken string) (*HuaweiUserInfo, error) {
	// 创建HTTP客户端
	client := &http.Client{}
	// 构建请求URL
	userInfoURL := "https://account.cloud.huawei.com/rest.php?nsp_svc=GOpen.User.getInfo"
	// 构建请求参数
	data := url.Values{}
	data.Set("access_token", accessToken)
	data.Set("getNickName", "1") // 返回真实昵称，0返回匿名化账号
	// 构建请求
	req, err := http.NewRequestWithContext(ctx, "POST", userInfoURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 解析响应
	var result struct {
		OpenID         string `json:"openID"`
		UnionID        string `json:"unionID"`
		DisplayName    string `json:"displayName"`
		HeadPictureURL string `json:"headPictureURL"`
		Error          string `json:"error,omitempty"`
		ErrorCode      string `json:"error_code,omitempty"`
		ErrorMsg       string `json:"error_msg,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	// 检查错误
	if result.Error != "" || result.ErrorCode != "" {
		return nil, fmt.Errorf("华为API错误: %s (code: %s, message: %s)", result.Error, result.ErrorCode, result.ErrorMsg)
	}
	// 构建返回结果
	return &HuaweiUserInfo{
		OpenID:         result.OpenID,
		UnionID:        result.UnionID,
		DisplayName:    result.DisplayName,
		HeadPictureURL: result.HeadPictureURL,
	}, nil
}

// GetPhoneNumber 获取用户手机号
func (h *HuaweiClient) GetPhoneNumber(ctx context.Context, accessToken, verifierToken string) (*HuaweiPhoneInfo, error) {
	// 创建HTTP客户端
	client := &http.Client{}
	// 构建请求URL
	phoneURL := "https://account.cloud.huawei.com/rest.php?nsp_svc=GOpen.User.getPhone"
	// 构建请求参数
	data := url.Values{}
	data.Set("access_token", accessToken)
	data.Set("verifier_token", verifierToken)
	// 构建请求
	req, err := http.NewRequestWithContext(ctx, "POST", phoneURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 解析响应
	var result struct {
		PhoneNumber string `json:"phoneNumber"`
		CountryCode string `json:"countryCode"`
		Error       string `json:"error,omitempty"`
		ErrorCode   string `json:"error_code,omitempty"`
		ErrorMsg    string `json:"error_msg,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	// 检查错误
	if result.Error != "" || result.ErrorCode != "" {
		return nil, fmt.Errorf("华为API错误: %s (code: %s, message: %s)", result.Error, result.ErrorCode, result.ErrorMsg)
	}
	// 构建返回结果
	return &HuaweiPhoneInfo{
		PhoneNumber: result.PhoneNumber,
		CountryCode: result.CountryCode,
	}, nil
}

// QuickLoginByCode 快速登录，使用code获取用户信息
func (h *HuaweiClient) QuickLoginByCode(ctx context.Context, code string) (*HuaweiQuickLoginResult, error) {
	// 先获取access token
	token, err := h.ExchangeToken(ctx, code)
	if err != nil {
		return nil, err
	}
	// 然后获取用户信息
	userInfo, err := h.GetUserInfo(ctx, token.AccessToken)
	if err != nil {
		return nil, err
	}
	// 构建返回结果
	return &HuaweiQuickLoginResult{
		Token:    token,
		UserInfo: userInfo,
	}, nil
}

// HuaweiUserInfo 用户信息结构体
type HuaweiUserInfo struct {
	OpenID         string `json:"openID"`
	UnionID        string `json:"unionID"`
	DisplayName    string `json:"displayName"`
	HeadPictureURL string `json:"headPictureURL"`
}

// HuaweiPhoneInfo 用户手机号信息结构体
type HuaweiPhoneInfo struct {
	PhoneNumber string `json:"phoneNumber"`
	CountryCode string `json:"countryCode"`
}

// HuaweiQuickLoginResult 快速登录结果结构体
type HuaweiQuickLoginResult struct {
	Token    *oauth2.Token   `json:"token"`
	UserInfo *HuaweiUserInfo `json:"userInfo"`
}

// SetRedirectURL 设置重定向URL
func (h *HuaweiClient) SetRedirectURL(redirectURL string) {
	h.config.RedirectURL = redirectURL
}

// SetScopes 设置请求的scope
func (h *HuaweiClient) SetScopes(scopes []string) {
	h.config.Scopes = scopes
}
