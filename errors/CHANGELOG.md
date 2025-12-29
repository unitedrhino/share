# 多语言键值格式更新日志

## 更新概述

将多语言定义的key从原来的下划线命名法修改为 `error.${错误文件名}.错误信息的小驼峰格式`。

## 更新内容

### 1. 系统错误 (sys.go)
- `success` → `error.sys.success`
- `other_error` → `error.sys.otherError`
- `token_expired` → `error.sys.tokenExpired`
- `parameter_error` → `error.sys.parameterError`
- `system_error` → `error.sys.systemError`
- `database_error` → `error.sys.databaseError`
- `not_found` → `error.sys.notFound`
- `duplicate_parameter` → `error.sys.duplicateParameter`
- `signature_expired` → `error.sys.signatureExpired`
- `insufficient_privileges` → `error.sys.insufficientPrivileges`
- `method_not_supported` → `error.sys.methodNotSupported`
- `invalid_parameter_type` → `error.sys.invalidParameterType`
- `parameter_out_of_range` → `error.sys.parameterOutOfRange`
- `timeout` → `error.sys.timeout`
- `server_error` → `error.sys.serverError`
- `not_implemented` → `error.sys.notImplemented`
- `not_empty` → `error.sys.notEmpty`
- `system_panic` → `error.sys.systemPanic`
- `not_enabled` → `error.sys.notEnabled`
- `enterprise_feature` → `error.sys.enterpriseFeature`
- `script_execution_failed` → `error.sys.scriptExecutionFailed`
- `in_progress` → `error.sys.inProgress`
- `execution_failed_rollback` → `error.sys.executionFailedRollback`
- `skip_execution` → `error.sys.skipExecution`

### 2. 用户错误 (user.go)
- `username_already_registered` → `error.user.usernameAlreadyRegistered`
- `mobile_already_taken` → `error.user.mobileAlreadyTaken`
- `not_registered` → `error.user.notRegistered`
- `account_or_password_error` → `error.user.accountOrPasswordError`
- `captcha_error` → `error.user.captchaError`
- `uid_incorrect` → `error.user.uidIncorrect`
- `not_logged_in` → `error.user.notLoggedIn`
- `login_method_not_supported` → `error.user.loginMethodNotSupported`
- `registration_step_one_failed` → `error.user.registrationStepOneFailed`
- `duplicate_registration` → `error.user.duplicateRegistration`
- `username_required` → `error.user.usernameRequired`
- `password_strength_insufficient` → `error.user.passwordStrengthInsufficient`
- `get_user_info_partial_failure` → `error.user.getUserInfoPartialFailure`
- `username_format_error` → `error.user.usernameFormatError`
- `account_or_ip_forbidden` → `error.user.accountOrIpForbidden`
- `use_captcha` → `error.user.useCaptcha`
- `account_disabled` → `error.user.accountDisabled`
- `account_already_bound` → `error.user.accountAlreadyBound`
- `account_kicked_out` → `error.user.accountKickedOut`
- `account_not_bound` → `error.user.accountNotBound`
- `image_captcha_required` → `error.user.imageCaptchaRequired`

### 3. 设备错误 (device.go)
- `response_param_error` → `error.device.responseParamError`
- `device_timeout` → `error.device.deviceTimeout`
- `device_offline` → `error.device.deviceOffline`
- `device_response_error` → `error.device.deviceResponseError`
- `device_already_bound` → `error.device.deviceAlreadyBound`
- `device_not_bound` → `error.device.deviceNotBound`
- `device_cannot_bound` → `error.device.deviceCannotBound`

### 4. 媒体错误 (media.go)
- `media_create_error` → `error.media.mediaCreateError`
- `media_update_error` → `error.media.mediaUpdateError`
- `media_not_found_error` → `error.media.mediaNotFoundError`
- `media_active_error` → `error.media.mediaActiveError`
- `media_pull_create_error` → `error.media.mediaPullCreateError`
- `media_stream_delete_error` → `error.media.mediaStreamDeleteError`
- `media_record_not_found` → `error.media.mediaRecordNotFound`
- `media_sip_update_error` → `error.media.mediaSipUpdateError`
- `media_sip_dev_create_error` → `error.media.mediaSipDevCreateError`
- `media_sip_chn_create_error` → `error.media.mediaSipChnCreateError`
- `media_sip_play_error` → `error.media.mediaSipPlayError`

### 5. OTA错误 (ota.go)
- `ota_retry_status_error` → `error.ota.otaRetryStatusError`
- `ota_cancel_status_error` → `error.ota.otaCancelStatusError`
- `ota_device_num_error` → `error.ota.otaDeviceNumError`

### 6. 文件错误 (file.go)
- `upload_failed` → `error.file.uploadFailed`

### 7. 用户数据错误 (ud.go)
- `trigger_type_not_supported` → `error.ud.triggerTypeNotSupported`

## 语言文件更新

### 中文语言文件 (zh.json)
- 完全重写，使用新的键值格式
- 保持所有中文翻译内容不变

### 英文语言文件 (en.json)
- 完全重写，使用新的键值格式
- 保持所有英文翻译内容不变

## 文档更新

1. **README.md** - 更新了多语言键值规范说明
2. **usage_example.md** - 更新了使用示例和错误码分类表
3. **CHANGELOG.md** - 新增变更日志文档

## 新的键值命名规则

- 格式：`error.${错误文件名}.错误信息的小驼峰格式`
- 示例：
  - `error.sys.success` - 系统成功
  - `error.user.usernameAlreadyRegistered` - 用户名已注册
  - `error.device.deviceTimeout` - 设备超时

## 兼容性说明

⚠️ **重要**：此次更新是破坏性更改，所有使用旧键值格式的代码都需要更新。

## 使用方法

```go
// 基本使用
err := errors.System
zhMsg := err.GetI18nMsg("zh")  // 获取中文消息
enMsg := err.GetI18nMsg("en")  // 获取英文消息

// 带自定义消息
err := errors.Parameter.WithMsg("用户名不能为空")

// 格式化消息
err := errors.Parameter.WithMsgf("用户 %s 不存在", "testuser")
```
