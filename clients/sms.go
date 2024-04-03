package clients

import (
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"strings"
)

type Sms struct {
	c      conf.Sms
	aliCli *dysmsapi20170525.Client
}

type SendSmsParam struct {
	PhoneNumbers  []string       `json:"phoneNumbers"`
	SignName      string         `json:"signName"`
	TemplateCode  string         `json:"templateCode"`
	TemplateParam map[string]any `json:"templateParam"`
}

func NewSms(c conf.Sms) (*Sms, error) {
	switch c.Mode {
	case conf.SmsAli:
		cli, err := CreateAliClient(c)
		if err != nil {
			return nil, err
		}
		return &Sms{c: c, aliCli: cli}, nil
	default:
		return nil, errors.System.AddMsg("不支持的短信配置类型")
	}
}

func (s *Sms) SendSms(param SendSmsParam) error {
	if !s.c.Enable {
		return errors.NotEnable.AddMsg("未开启短信服务")
	}
	switch s.c.Mode {
	case conf.SmsAli:
		return s.SendSmsAli(param)
	}
	return nil
}
func (s *Sms) SendSmsAli(param SendSmsParam) error {
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(strings.Join(param.PhoneNumbers, ",")),
		SignName:      tea.String(param.SignName),
		TemplateCode:  tea.String(param.TemplateCode),
		TemplateParam: tea.String(utils.MarshalNoErr(param.TemplateParam)),
	}
	_, err := s.aliCli.SendSmsWithOptions(sendSmsRequest, &util.RuntimeOptions{})
	if err != nil {
		return err
	}
	return nil
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateAliClient(sms conf.Sms) (cli *dysmsapi20170525.Client, err error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String(sms.Ali.AccessKeyID),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: tea.String(sms.Ali.AccessKeySecret),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	cli = &dysmsapi20170525.Client{}
	cli, err = dysmsapi20170525.NewClient(config)
	return cli, err
}

//
//func _main(args []*string) (_err error) {
//	client, _err := CreateAliClient()
//	if _err != nil {
//		return _err
//	}
//
//	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
//		PhoneNumbers:  tea.String("17052709767"),
//		SignName:      tea.String("EbelongTool"),
//		TemplateCode:  tea.String("SMS_288215142"),
//		TemplateParam: tea.String(`{"code":123123}`),
//	}
//
//	tryErr := func() (_e error) {
//		defer func() {
//			if r := tea.Recover(recover()); r != nil {
//				_e = r
//			}
//		}()
//		// 复制代码运行请自行打印 API 的返回值
//		_, _err = client.SendSmsWithOptions(sendSmsRequest, &util.RuntimeOptions{})
//		if _err != nil {
//			return _err
//		}
//
//		return nil
//	}()
//
//	if tryErr != nil {
//		var error = &tea.SDKError{}
//		if _t, ok := tryErr.(*tea.SDKError); ok {
//			error = _t
//		} else {
//			error.Message = tea.String(tryErr.Error())
//		}
//		// 此处仅做打印展示，请谨慎对待异常处理，在工程项目中切勿直接忽略异常。
//		// 错误 message
//		fmt.Println(tea.StringValue(error.Message))
//		// 诊断地址
//		var data interface{}
//		d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
//		d.Decode(&data)
//		if m, ok := data.(map[string]interface{}); ok {
//			recommend, _ := m["Recommend"]
//			fmt.Println(recommend)
//		}
//		_, _err = util.AssertAsString(error.Message)
//		if _err != nil {
//			return _err
//		}
//	}
//	return _err
//}
//
//func main() {
//	err := _main(tea.StringSlice(os.Args[1:]))
//	if err != nil {
//		panic(err)
//	}
//}
