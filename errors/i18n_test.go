package errors

import (
	"testing"
)

// TestI18nBasic 测试基本的多语言功能
func TestI18nBasic(t *testing.T) {
	tests := []struct {
		name     string
		err      *CodeError
		lang     string
		expected string
	}{
		// 系统错误测试
		{
			name:     "系统错误-中文",
			err:      System,
			lang:     "zh",
			expected: "系统错误",
		},
		{
			name:     "系统错误-英文",
			err:      System,
			lang:     "en",
			expected: "System error",
		},
		{
			name:     "成功-中文",
			err:      OK,
			lang:     "zh",
			expected: "成功",
		},
		{
			name:     "成功-英文",
			err:      OK,
			lang:     "en",
			expected: "Success",
		},
		{
			name:     "参数错误-中文",
			err:      Parameter,
			lang:     "zh",
			expected: "参数错误",
		},
		{
			name:     "参数错误-英文",
			err:      Parameter,
			lang:     "en",
			expected: "Parameter error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.GetI18nMsg(tt.lang)
			if result != tt.expected {
				t.Errorf("GetI18nMsg() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestI18nUserErrors 测试用户相关错误的多语言
func TestI18nUserErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      *CodeError
		lang     string
		expected string
	}{
		{
			name:     "用户名已注册-中文",
			err:      DuplicateUsername,
			lang:     "zh",
			expected: "用户名已经注册",
		},
		{
			name:     "用户名已注册-英文",
			err:      DuplicateUsername,
			lang:     "en",
			expected: "Username already registered",
		},
		{
			name:     "手机号已被占用-中文",
			err:      DuplicateMobile,
			lang:     "zh",
			expected: "手机号已经被占用",
		},
		{
			name:     "手机号已被占用-英文",
			err:      DuplicateMobile,
			lang:     "en",
			expected: "Mobile number already taken",
		},
		{
			name:     "账号或密码错误-中文",
			err:      Password,
			lang:     "zh",
			expected: "账号或密码错误",
		},
		{
			name:     "账号或密码错误-英文",
			err:      Password,
			lang:     "en",
			expected: "Account or password error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.GetI18nMsg(tt.lang)
			if result != tt.expected {
				t.Errorf("GetI18nMsg() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestI18nDeviceErrors 测试设备相关错误的多语言
func TestI18nDeviceErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      *CodeError
		lang     string
		expected string
	}{
		{
			name:     "设备超时-中文",
			err:      DeviceTimeOut,
			lang:     "zh",
			expected: "设备回复超时",
		},
		{
			name:     "设备超时-英文",
			err:      DeviceTimeOut,
			lang:     "en",
			expected: "Device response timeout",
		},
		{
			name:     "设备离线-中文",
			err:      NotOnline,
			lang:     "zh",
			expected: "设备离线，请检查电源或设备",
		},
		{
			name:     "设备离线-英文",
			err:      NotOnline,
			lang:     "en",
			expected: "Device offline, please check power or device",
		},
		{
			name:     "设备已被绑定-中文",
			err:      DeviceBound,
			lang:     "zh",
			expected: "设备已被绑定",
		},
		{
			name:     "设备已被绑定-英文",
			err:      DeviceBound,
			lang:     "en",
			expected: "Device already bound",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.GetI18nMsg(tt.lang)
			if result != tt.expected {
				t.Errorf("GetI18nMsg() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestI18nMediaErrors 测试媒体相关错误的多语言
func TestI18nMediaErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      *CodeError
		lang     string
		expected string
	}{
		{
			name:     "媒体创建错误-中文",
			err:      MediaCreateError,
			lang:     "zh",
			expected: "流服务创建失败",
		},
		{
			name:     "媒体创建错误-英文",
			err:      MediaCreateError,
			lang:     "en",
			expected: "Media service creation failed",
		},
		{
			name:     "媒体更新错误-中文",
			err:      MediaUpdateError,
			lang:     "zh",
			expected: "流服务更新失败",
		},
		{
			name:     "媒体更新错误-英文",
			err:      MediaUpdateError,
			lang:     "en",
			expected: "Media service update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.GetI18nMsg(tt.lang)
			if result != tt.expected {
				t.Errorf("GetI18nMsg() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestI18nOtaErrors 测试OTA相关错误的多语言
func TestI18nOtaErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      *CodeError
		lang     string
		expected string
	}{
		{
			name:     "OTA重试状态错误-中文",
			err:      OtaRetryStatusError,
			lang:     "zh",
			expected: "升级状态不允许重新升级",
		},
		{
			name:     "OTA重试状态错误-英文",
			err:      OtaRetryStatusError,
			lang:     "en",
			expected: "OTA status does not allow retry",
		},
		{
			name:     "OTA取消状态错误-中文",
			err:      OtaCancleStatusError,
			lang:     "zh",
			expected: "升级状态已结束",
		},
		{
			name:     "OTA取消状态错误-英文",
			err:      OtaCancleStatusError,
			lang:     "en",
			expected: "OTA status has ended",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.GetI18nMsg(tt.lang)
			if result != tt.expected {
				t.Errorf("GetI18nMsg() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestI18nFileErrors 测试文件相关错误的多语言
func TestI18nFileErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      *CodeError
		lang     string
		expected string
	}{
		{
			name:     "上传失败-中文",
			err:      Upload,
			lang:     "zh",
			expected: "上传失败",
		},
		{
			name:     "上传失败-英文",
			err:      Upload,
			lang:     "en",
			expected: "Upload failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.GetI18nMsg(tt.lang)
			if result != tt.expected {
				t.Errorf("GetI18nMsg() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestI18nUdErrors 测试用户数据相关错误的多语言
func TestI18nUdErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      *CodeError
		lang     string
		expected string
	}{
		{
			name:     "触发类型不支持-中文",
			err:      TriggerType,
			lang:     "zh",
			expected: "触发类型不支持",
		},
		{
			name:     "触发类型不支持-英文",
			err:      TriggerType,
			lang:     "en",
			expected: "Trigger type not supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.GetI18nMsg(tt.lang)
			if result != tt.expected {
				t.Errorf("GetI18nMsg() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestI18nWithCustomMessage 测试带自定义消息的多语言
func TestI18nWithCustomMessage(t *testing.T) {
	// 测试添加自定义消息
	err := Parameter.WithMsg("用户名不能为空")

	zhMsg := err.GetI18nMsg("zh")
	enMsg := err.GetI18nMsg("en")

	// 由于是自定义消息，应该直接返回消息内容
	if zhMsg != "用户名不能为空" {
		t.Errorf("WithMsg() 中文消息 = %v, want %v", zhMsg, "用户名不能为空")
	}

	if enMsg != "用户名不能为空" {
		t.Errorf("WithMsg() 英文消息 = %v, want %v", enMsg, "用户名不能为空")
	}
}

// TestI18nWithFormatMessage 测试格式化消息的多语言
func TestI18nWithFormatMessage(t *testing.T) {
	// 测试格式化消息
	err := Parameter.WithMsgf("用户 %s 不存在", "testuser")

	zhMsg := err.GetI18nMsg("zh")
	enMsg := err.GetI18nMsg("en")

	expected := "用户 testuser 不存在"
	if zhMsg != expected {
		t.Errorf("WithMsgf() 中文消息 = %v, want %v", zhMsg, expected)
	}

	if enMsg != expected {
		t.Errorf("WithMsgf() 英文消息 = %v, want %v", enMsg, expected)
	}
}

// TestI18nMultipleMessages 测试多消息组合的多语言
func TestI18nMultipleMessages(t *testing.T) {
	// 测试多消息组合
	err := Parameter.AddMsg("用户名不能为空").AddMsg("密码不能为空")

	zhMsg := err.GetI18nMsg("zh")
	enMsg := err.GetI18nMsg("en")

	expected := "用户名不能为空:密码不能为空"
	if zhMsg != expected {
		t.Errorf("AddMsg() 中文消息 = %v, want %v", zhMsg, expected)
	}

	if enMsg != expected {
		t.Errorf("AddMsg() 英文消息 = %v, want %v", enMsg, expected)
	}
}

// TestI18nDefaultLanguage 测试默认语言处理
func TestI18nDefaultLanguage(t *testing.T) {
	// 测试空语言参数
	err := System
	msg := err.GetI18nMsg("")

	// 空语言参数应该返回默认语言（中文）
	if msg == "" {
		t.Error("GetI18nMsg() 空语言参数应该返回默认消息")
	}

	// 测试不支持的语言
	msg = err.GetI18nMsg("fr")
	if msg == "" {
		t.Error("GetI18nMsg() 不支持的语言应该返回默认消息")
	}
}

// TestI18nErrorCodes 测试错误码的正确性
func TestI18nErrorCodes(t *testing.T) {
	tests := []struct {
		name     string
		err      *CodeError
		expected int64
	}{
		{"系统错误码", System, 100007},
		{"用户错误码", DuplicateUsername, 1000001},
		{"设备错误码", DeviceTimeOut, 2000002},
		{"媒体错误码", MediaCreateError, 3000001},
		{"OTA错误码", OtaRetryStatusError, 2100001},
		{"文件错误码", Upload, 1000001},
		{"用户数据错误码", TriggerType, 4000001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := tt.err.GetCode()
			if code != tt.expected {
				t.Errorf("GetCode() = %v, want %v", code, tt.expected)
			}
		})
	}
}

// TestI18nAllKeysExist 测试所有键值是否都存在对应的翻译
func TestI18nAllKeysExist(t *testing.T) {
	// 测试所有系统错误
	systemErrors := []*CodeError{
		OK, Default, TokenExpired, TokenNotValidYet,
		TokenMalformed, TokenInvalid, Parameter, System,
		Database, NotFind, Duplicate, SignatureExpired,
		Permissions, Method, Type, OutRange,
		TimeOut, Server, NotRealize, NotEmpty,
		Panic, NotEnable, Company, Script,
		OnGoing, Failure, Jump,
	}

	for _, err := range systemErrors {
		zhMsg := err.GetI18nMsg("zh")
		enMsg := err.GetI18nMsg("en")

		if zhMsg == "" {
			t.Errorf("系统错误 %v 缺少中文翻译", err.GetCode())
		}
		if enMsg == "" {
			t.Errorf("系统错误 %v 缺少英文翻译", err.GetCode())
		}
	}

	// 测试所有用户错误
	userErrors := []*CodeError{
		DuplicateUsername, DuplicateMobile, UnRegister,
		Password, Captcha, UidNotRight, NotLogin,
		NotSupportLogin, RegisterOne, DuplicateRegister,
		NeedUserName, PasswordLevel, GetInfoPartFailure,
		UsernameFormatErr, AccountOrIpForbidden, UseCaptcha,
		AccountDisable, BindAccount, AccountKickedOut,
		UnBindAccount, NeedImgCaptcha,
	}

	for _, err := range userErrors {
		zhMsg := err.GetI18nMsg("zh")
		enMsg := err.GetI18nMsg("en")

		if zhMsg == "" {
			t.Errorf("用户错误 %v 缺少中文翻译", err.GetCode())
		}
		if enMsg == "" {
			t.Errorf("用户错误 %v 缺少英文翻译", err.GetCode())
		}
	}
}

// BenchmarkI18nGetMessage 性能测试
func BenchmarkI18nGetMessage(b *testing.B) {
	err := System

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.GetI18nMsg("zh")
	}
}

// BenchmarkI18nGetMessageEnglish 英文性能测试
func BenchmarkI18nGetMessageEnglish(b *testing.B) {
	err := System

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.GetI18nMsg("en")
	}
}
