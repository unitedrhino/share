# 多语言错误系统使用示例

## 基本用法

```go
package main

import (
    "fmt"
    "gitee.com/unitedrhino/share/errors"
)

func main() {
    // 1. 使用系统错误
    sysErr := errors.System
    fmt.Printf("系统错误 - 中文: %s\n", sysErr.GetI18nMsg("zh"))
    fmt.Printf("系统错误 - 英文: %s\n", sysErr.GetI18nMsg("en"))
    
    // 2. 使用用户错误
    userErr := errors.DuplicateUsername
    fmt.Printf("用户错误 - 中文: %s\n", userErr.GetI18nMsg("zh"))
    fmt.Printf("用户错误 - 英文: %s\n", userErr.GetI18nMsg("en"))
    
    // 3. 添加自定义消息
    customErr := errors.Parameter.WithMsg("用户名不能为空")
    fmt.Printf("自定义错误 - 中文: %s\n", customErr.GetI18nMsg("zh"))
    fmt.Printf("自定义错误 - 英文: %s\n", customErr.GetI18nMsg("en"))
    
    // 4. 使用格式化消息
    formatErr := errors.Parameter.WithMsgf("用户 %s 不存在", "testuser")
    fmt.Printf("格式化错误 - 中文: %s\n", formatErr.GetI18nMsg("zh"))
    fmt.Printf("格式化错误 - 英文: %s\n", formatErr.GetI18nMsg("en"))
    
    // 5. 多消息组合
    multiErr := errors.Parameter.AddMsg("用户名不能为空").AddMsg("密码不能为空")
    fmt.Printf("多消息错误 - 中文: %s\n", multiErr.GetI18nMsg("zh"))
    fmt.Printf("多消息错误 - 英文: %s\n", multiErr.GetI18nMsg("en"))
}
```

## 在HTTP处理中使用

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    // 获取Accept-Language头
    acceptLang := r.Header.Get("Accept-Language")
    
    // 创建错误
    err := errors.Parameter.WithMsg("用户名不能为空")
    
    // 根据语言返回错误消息
    var lang string
    if strings.Contains(acceptLang, "zh") {
        lang = "zh"
    } else {
        lang = "en"
    }
    
    msg := err.GetI18nMsg(lang)
    
    // 返回JSON响应
    response := map[string]interface{}{
        "code": err.GetCode(),
        "msg":  msg,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

## 在gRPC服务中使用

```go
func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
    // 验证用户名
    if req.Username == "" {
        err := errors.Parameter.WithMsg("用户名不能为空")
        return nil, err.ToRpc("zh") // 根据客户端语言设置
    }
    
    // 检查用户名是否已存在
    if userExists(req.Username) {
        err := errors.DuplicateUsername
        return nil, err.ToRpc("zh")
    }
    
    // 创建用户逻辑...
    return &pb.CreateUserResponse{}, nil
}
```

## 错误码分类

| 错误类型 | 错误码范围 | 示例 | 多语言键 |
|---------|-----------|------|---------|
| 系统错误 | 100000-199999 | `errors.System` | `error.sys.systemError` |
| 用户错误 | 1000000-1999999 | `errors.DuplicateUsername` | `error.user.usernameAlreadyRegistered` |
| 设备错误 | 2000000-2999999 | `errors.DeviceTimeOut` | `error.device.deviceTimeout` |
| 媒体错误 | 3000000-3999999 | `errors.MediaCreateError` | `error.media.mediaCreateError` |
| OTA错误 | 2100000-2199999 | `errors.OtaRetryStatusError` | `error.ota.otaRetryStatusError` |
| 文件错误 | 1000000-1999999 | `errors.Upload` | `error.file.uploadFailed` |
| 用户数据错误 | 4000000-4999999 | `errors.TriggerType` | `error.ud.triggerTypeNotSupported` |

## 支持的语言

- `zh` - 中文
- `en` - 英文

## 注意事项

1. 所有错误消息都使用标准化的键名格式：`error.${错误文件名}.错误信息的小驼峰格式`
2. 中英文语言文件必须保持同步
3. 使用 `WithMsgf` 进行格式化消息时，确保参数类型正确
4. 错误码按模块分类，避免冲突
5. 键名使用小驼峰命名法，避免使用下划线
