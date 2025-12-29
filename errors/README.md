# 标准多语言错误系统

## 概述

这是一个标准化的多语言错误处理系统，支持中文和英文两种语言。系统使用 `go-i18n` 库实现国际化功能。

## 主要特性

- ✅ 标准化的多语言键值对
- ✅ 支持中文（zh）和英文（en）
- ✅ 修复了原有的拼写错误（I81n -> I18n）
- ✅ 统一的错误码管理
- ✅ 支持格式化消息
- ✅ 支持多消息组合

## 文件结构

```
errors/
├── locale/
│   ├── zh.json          # 中文语言文件
│   └── en.json          # 英文语言文件
├── i81n.go             # 多语言接口实现（已修复拼写）
├── baseerror.go        # 基础错误结构
├── sys.go              # 系统错误定义
├── user.go             # 用户相关错误
├── device.go           # 设备相关错误
├── media.go            # 媒体相关错误
├── ota.go              # OTA升级错误
├── file.go             # 文件操作错误
├── ud.go               # 用户数据错误
└── example.go          # 使用示例
```

## 使用方法

### 1. 基本错误使用

```go
// 创建系统错误
err := errors.System
zhMsg := err.GetI18nMsg("zh")  // 获取中文消息
enMsg := err.GetI18nMsg("en")  // 获取英文消息
```

### 2. 带消息的错误

```go
// 添加自定义消息
err := errors.Parameter.WithMsg("用户名不能为空")
zhMsg := err.GetI18nMsg("zh")
enMsg := err.GetI18nMsg("en")
```

### 3. 格式化消息

```go
// 使用格式化消息
err := errors.Parameter.WithMsgf("用户 %s 不存在", "testuser")
zhMsg := err.GetI18nMsg("zh")
```

### 4. 多消息组合

```go
// 添加多个消息
err := errors.Parameter.AddMsg("用户名不能为空").AddMsg("密码不能为空")
zhMsg := err.GetI18nMsg("zh")
```

## 错误码分类

- **系统错误**: 100000 - 199999
- **用户错误**: 1000000 - 1999999  
- **设备错误**: 2000000 - 2999999
- **媒体错误**: 3000000 - 3999999
- **OTA错误**: 2100000 - 2199999
- **文件错误**: 1000000 - 1999999
- **用户数据错误**: 4000000 - 4999999

## 多语言键值规范

所有多语言键都使用 `error.${错误文件名}.错误信息的小驼峰格式`，例如：
- `error.sys.success` - 成功
- `error.sys.parameterError` - 参数错误
- `error.user.usernameAlreadyRegistered` - 用户名已注册

## 添加新的错误

1. 在相应的错误文件中定义错误码和键值
2. 在 `locale/zh.json` 中添加中文翻译
3. 在 `locale/en.json` 中添加英文翻译

示例：
```go
// 在错误文件中
var NewError = NewCodeError(ErrorCode+1, "error.sys.newError")

// 在 zh.json 中
"error.sys.newError": "新的错误消息"

// 在 en.json 中  
"error.sys.newError": "New error message"
```

## 测试

运行测试以验证多语言功能：
```bash
go test -v ./errors
```

## 注意事项

1. 所有多语言键必须在中英文文件中都有对应翻译
2. 使用标准化的键名格式：`error.${错误文件名}.错误信息的小驼峰格式`
3. 错误码按模块分类，避免冲突
4. 格式化消息使用 `WithMsgf` 方法
5. 键名使用小驼峰命名法，避免使用下划线
