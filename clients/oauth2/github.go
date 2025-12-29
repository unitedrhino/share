package oauth2

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/share/conf"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GithubClient struct {
	config *oauth2.Config
}

func NewGithub(ctx context.Context, conf *conf.ThirdConf) *GithubClient {
	if conf == nil {
		return nil
	}
	// 配置OAuth2
	oauthConfig = &oauth2.Config{
		ClientID:     conf.AppID,
		ClientSecret: conf.AppSecret,
		//RedirectURL:  redirectURL,
		Scopes: []string{
			"user:email", // 请求访问用户邮箱的权限
			"read:user",  // 请求访问用户基本信息的权限
		},
		Endpoint: github.Endpoint,
	}
	return &GithubClient{
		config: oauthConfig,
	}
}

func (g *GithubClient) GetAuthCodeURL() string {
	// 生成并保存随机状态值，用于防止CSRF攻击
	state := "random-state-string" // 实际应用中应该使用更安全的随机生成方式
	// 生成授权URL并重定向
	url := oauthConfig.AuthCodeURL(state)
	return url
}
func (g *GithubClient) GetUserInfo(ctx context.Context, code string) (*GitHubUser, error) {
	// 使用授权码获取访问令牌
	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	// 使用访问令牌创建HTTP客户端
	client := oauthConfig.Client(ctx, token)

	// 获取用户信息
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析用户信息
	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// GitHub用户信息结构体
type GitHubUser struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
}

var (
	oauthConfig *oauth2.Config
)

//
//// 初始化OAuth配置
//func init() {
//	// 从.env文件加载环境变量
//	err := godotenv.Load()
//	if err != nil {
//		log.Printf("警告: 未找到.env文件，将使用系统环境变量: %v", err)
//	}
//
//	// 获取环境变量
//	clientID := os.Getenv("GITHUB_CLIENT_ID")
//	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
//	redirectURL := os.Getenv("GITHUB_REDIRECT_URL")
//
//	if clientID == "" || clientSecret == "" {
//		log.Fatal("请设置GITHUB_CLIENT_ID和GITHUB_CLIENT_SECRET环境变量")
//	}
//
//	if redirectURL == "" {
//		redirectURL = "http://localhost:8080/callback"
//		log.Printf("使用默认重定向URL: %s", redirectURL)
//	}
//
//	// 配置OAuth2
//	oauthConfig = &oauth2.Config{
//		ClientID:     clientID,
//		ClientSecret: clientSecret,
//		RedirectURL:  redirectURL,
//		Scopes: []string{
//			"user:email", // 请求访问用户邮箱的权限
//			"read:user",  // 请求访问用户基本信息的权限
//		},
//		Endpoint: github.Endpoint,
//	}
//}
//
//// 主页路由处理函数
//func homeHandler(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintf(w, `
//		<html>
//			<head>
//				<title>GitHub OAuth2 登录示例</title>
//				<style>
//					body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
//					.login-btn {
//						background-color: #2ea44f;
//						color: white;
//						padding: 10px 20px;
//						text-decoration: none;
//						border-radius: 5px;
//						font-size: 16px;
//					}
//					.login-btn:hover { background-color: #238636; }
//				</style>
//			</head>
//			<body>
//				<h1>欢迎使用GitHub OAuth2登录示例</h1>
//				<a href="/login" class="login-btn">使用GitHub登录</a>
//			</body>
//		</html>
//	`)
//}
//
//// 登录路由处理函数 - 重定向到GitHub授权页面
//func loginHandler(w http.ResponseWriter, r *http.Request) {
//	// 生成并保存随机状态值，用于防止CSRF攻击
//	state := "random-state-string" // 实际应用中应该使用更安全的随机生成方式
//
//	// 生成授权URL并重定向
//	url := oauthConfig.AuthCodeURL(state)
//	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
//}
//
//// 回调路由处理函数 - 处理GitHub的回调
//func callbackHandler(w http.ResponseWriter, r *http.Request) {
//	// 验证状态值，防止CSRF攻击
//	state := r.FormValue("state")
//	if state != "random-state-string" { // 应与登录时生成的状态值比对
//		fmt.Fprintf(w, "状态验证失败: %s", state)
//		return
//	}
//
//	// 从回调中获取授权码
//	code := r.FormValue("code")
//	if code == "" {
//		fmt.Fprintf(w, "未获取到授权码")
//		return
//	}
//
//	// 使用授权码获取访问令牌
//	token, err := oauthConfig.Exchange(r.Context(), code)
//	if err != nil {
//		fmt.Fprintf(w, "获取令牌失败: %v", err)
//		return
//	}
//
//	// 使用访问令牌创建HTTP客户端
//	client := oauthConfig.Client(r.Context(), token)
//
//	// 获取用户信息
//	resp, err := client.Get("https://api.github.com/user")
//	if err != nil {
//		fmt.Fprintf(w, "获取用户信息失败: %v", err)
//		return
//	}
//	defer resp.Body.Close()
//
//	// 解析用户信息
//	var user GitHubUser
//	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
//		fmt.Fprintf(w, "解析用户信息失败: %v", err)
//		return
//	}
//
//	// 显示用户信息
//	fmt.Fprintf(w, `
//		<html>
//			<head>
//				<title>登录成功 - GitHub用户信息</title>
//				<style>
//					body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
//					.profile { display: flex; align-items: center; gap: 20px; margin-bottom: 30px; }
//					.avatar { border-radius: 50%; }
//					.info { line-height: 1.6; }
//				</style>
//			</head>
//			<body>
//				<h1>登录成功！</h1>
//				<div class="profile">
//					<img src="%s" alt="用户头像" class="avatar" width="100" height="100">
//					<div class="info">
//						<h2>%s (%s)</h2>
//						<p>ID: %d</p>
//						<p>邮箱: %s</p>
//						<p>简介: %s</p>
//					</div>
//				</div>
//				<a href="/">返回首页</a>
//			</body>
//		</html>
//	`, user.AvatarURL, user.Name, user.Login, user.ID, user.Email, user.Bio)
//}
//
//func main() {
//	// 注册路由
//	http.HandleFunc("/", homeHandler)
//	http.HandleFunc("/login", loginHandler)
//	http.HandleFunc("/callback", callbackHandler)
//
//	// 启动服务器
//	port := os.Getenv("PORT")
//	if port == "" {
//		port = "8080"
//	}
//	log.Printf("服务器启动在 http://localhost:%s", port)
//	log.Fatal(http.ListenAndServe(":"+port, nil))
//}
