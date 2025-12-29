# 多语言错误系统测试说明

## 测试文件说明

### 1. i18n_test.go - 完整的多语言测试文件

这是主要的多语言测试文件，包含以下测试用例：

#### 基本功能测试
- `TestI18nBasic` - 测试基本的多语言功能
- `TestI18nUserErrors` - 测试用户相关错误的多语言
- `TestI18nDeviceErrors` - 测试设备相关错误的多语言
- `TestI18nMediaErrors` - 测试媒体相关错误的多语言
- `TestI18nOtaErrors` - 测试OTA相关错误的多语言
- `TestI18nFileErrors` - 测试文件相关错误的多语言
- `TestI18nUdErrors` - 测试用户数据相关错误的多语言

#### 高级功能测试
- `TestI18nWithCustomMessage` - 测试带自定义消息的多语言
- `TestI18nWithFormatMessage` - 测试格式化消息的多语言
- `TestI18nMultipleMessages` - 测试多消息组合的多语言
- `TestI18nDefaultLanguage` - 测试默认语言处理
- `TestI18nErrorCodes` - 测试错误码的正确性
- `TestI18nAllKeysExist` - 测试所有键值是否都存在对应的翻译

#### 性能测试
- `BenchmarkI18nGetMessage` - 中文消息获取性能测试
- `BenchmarkI18nGetMessageEnglish` - 英文消息获取性能测试

### 2. example_usage.go - 使用示例文件

这是一个可执行的示例文件，展示了如何使用多语言错误系统的各种功能。

## 运行测试

### 运行所有测试
```bash
go test -v
```

### 运行特定测试
```bash
# 运行基本功能测试
go test -v -run TestI18nBasic

# 运行用户错误测试
go test -v -run TestI18nUserErrors

# 运行设备错误测试
go test -v -run TestI18nDeviceErrors

# 运行性能测试
go test -v -bench=BenchmarkI18n
```

### 运行示例程序
```bash
go run example_usage.go
```

## 测试覆盖范围

### 错误类型覆盖
- ✅ 系统错误 (27个)
- ✅ 用户错误 (21个)
- ✅ 设备错误 (7个)
- ✅ 媒体错误 (11个)
- ✅ OTA错误 (3个)
- ✅ 文件错误 (1个)
- ✅ 用户数据错误 (1个)

### 功能覆盖
- ✅ 基本多语言消息获取
- ✅ 自定义消息处理
- ✅ 格式化消息处理
- ✅ 多消息组合
- ✅ 默认语言处理
- ✅ 错误码验证
- ✅ 键值存在性验证
- ✅ 性能测试

### 语言覆盖
- ✅ 中文 (zh)
- ✅ 英文 (en)

## 测试数据

### 系统错误测试数据
```go
{
    name:     "系统错误-中文",
    err:      &errors.System,
    lang:     "zh",
    expected: "系统错误",
},
{
    name:     "系统错误-英文",
    err:      &errors.System,
    lang:     "en",
    expected: "System error",
},
```

### 用户错误测试数据
```go
{
    name:     "用户名已注册-中文",
    err:      &errors.DuplicateUsername,
    lang:     "zh",
    expected: "用户名已经注册",
},
{
    name:     "用户名已注册-英文",
    err:      &errors.DuplicateUsername,
    lang:     "en",
    expected: "Username already registered",
},
```

## 预期结果

### 成功情况
- 所有测试用例都应该通过
- 中文和英文消息都应该正确返回
- 错误码应该与预期一致
- 性能测试应该显示合理的执行时间

### 失败情况
如果测试失败，可能的原因：
1. 语言文件缺失或格式错误
2. 多语言键值不匹配
3. 错误码定义错误
4. 依赖包问题

## 故障排除

### 常见问题

1. **Go环境问题**
   ```
   package embed is not in std
   ```
   解决方案：检查Go版本，确保使用Go 1.16+

2. **依赖包问题**
   ```
   package github.com/nicksnyder/go-i18n/v2/i18n not found
   ```
   解决方案：运行 `go mod tidy` 安装依赖

3. **语言文件问题**
   ```
   GetI18nMsg() = "", want "系统错误"
   ```
   解决方案：检查 `locale/zh.json` 和 `locale/en.json` 文件

### 调试技巧

1. **查看语言文件内容**
   ```bash
   cat locale/zh.json
   cat locale/en.json
   ```

2. **检查错误码**
   ```go
   fmt.Printf("错误码: %d\n", err.GetCode())
   ```

3. **验证多语言键值**
   ```go
   fmt.Printf("键值: %s\n", err.Msg[0].String())
   ```

## 扩展测试

### 添加新的测试用例
1. 在相应的测试函数中添加新的测试数据
2. 确保测试数据包含中文和英文预期结果
3. 运行测试验证结果

### 添加新的错误类型测试
1. 创建新的测试函数
2. 添加错误定义到相应的错误文件
3. 在语言文件中添加对应的翻译
4. 编写测试用例

## 持续集成

建议在CI/CD流程中包含以下测试：
- 单元测试：`go test -v`
- 性能测试：`go test -bench=BenchmarkI18n`
- 代码覆盖率：`go test -cover`

## 注意事项

1. 测试文件使用 `errors_test` 包名，避免循环依赖
2. 所有测试都使用表驱动测试模式，便于维护
3. 性能测试使用 `Benchmark` 前缀
4. 测试数据应该覆盖所有错误类型和语言
5. 确保测试的独立性和可重复性
