package smsClient

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/samber/lo"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	te "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/zeromicro/go-zero/core/logx"
	"slices"
	"strings"
)

type Sms struct {
	c          conf.Sms
	aliCli     *dysmsapi20170525.Client
	tencentCli *sms.Client
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
	case conf.SmsTencent:
		cli, err := CreateTencentClient(c)
		if err != nil {
			return nil, err
		}
		return &Sms{c: c, tencentCli: cli}, nil
	default:
		return nil, errors.System.AddMsg("不支持的短信配置类型")
	}
}

func (s *Sms) SendSms(ctx context.Context, param SendSmsParam) error {
	if !s.c.Enable {
		return errors.NotEnable.AddMsg("未开启短信服务")
	}
	switch s.c.Mode {
	case conf.SmsAli:
		return s.SendSmsAli(ctx, param)
	case conf.SmsTencent:
		return s.SendSmsTencent(ctx, param)
	}
	return nil
}
func (s *Sms) SendSmsAli(ctx context.Context, param SendSmsParam) error {
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
func (s *Sms) SendSmsTencent(ctx context.Context, param SendSmsParam) error {
	/* 实例化一个请求对象，根据调用的接口和实际情况，可以进一步设置请求参数
	 * 您可以直接查询SDK源码确定接口有哪些属性可以设置
	 * 属性可能是基本类型，也可能引用了另一个数据结构
	 * 推荐使用IDE进行开发，可以方便的跳转查阅各个接口和数据结构的文档说明 */
	request := sms.NewSendSmsRequest()

	/* 基本类型的设置:
	 * SDK采用的是指针风格指定参数，即使对于基本类型您也需要用指针来对参数赋值。
	 * SDK提供对基本类型的指针引用封装函数
	 * 帮助链接：
	 * 短信控制台: https://console.cloud.tencent.com/smsv2
	 * 腾讯云短信小助手: https://cloud.tencent.com/document/product/382/3773#.E6.8A.80.E6.9C.AF.E4.BA.A4.E6.B5.81 */

	/* 短信应用ID: 短信SdkAppId在 [短信控制台] 添加应用后生成的实际SdkAppId，示例如1400006666 */
	// 应用 ID 可前往 [短信控制台](https://console.cloud.tencent.com/smsv2/app-manage) 查看
	request.SmsSdkAppId = common.StringPtr(s.c.Tencent.AppID)

	/* 短信签名内容: 使用 UTF-8 编码，必须填写已审核通过的签名 */
	// 签名信息可前往 [国内短信](https://console.cloud.tencent.com/smsv2/csms-sign) 或 [国际/港澳台短信](https://console.cloud.tencent.com/smsv2/isms-sign) 的签名管理查看
	request.SignName = common.StringPtr(param.SignName)

	/* 模板 ID: 必须填写已审核通过的模板 ID */
	// 模板 ID 可前往 [国内短信](https://console.cloud.tencent.com/smsv2/csms-template) 或 [国际/港澳台短信](https://console.cloud.tencent.com/smsv2/isms-template) 的正文模板管理查看
	request.TemplateId = common.StringPtr(param.TemplateCode)
	keys := lo.Keys(param.TemplateParam)
	slices.Sort(keys)
	var paramSet []string
	for _, k := range keys {
		paramSet = append(paramSet, utils.ToString(param.TemplateParam[k]))
	}
	/* 模板参数: 模板参数的个数需要与 TemplateId 对应模板的变量个数保持一致，若无模板参数，则设置为空*/
	request.TemplateParamSet = common.StringPtrs(paramSet)

	/* 下发手机号码，采用 E.164 标准，+[国家或地区码][手机号]
	 * 示例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号*/
	request.PhoneNumberSet = common.StringPtrs(param.PhoneNumbers)
	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := s.tencentCli.SendSms(request)
	// 处理异常
	if _, ok := err.(*te.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return err
	}
	// 非SDK异常，直接失败。实际代码中可以加入其他的处理。
	if err != nil {
		logx.WithContext(ctx)
		return err
	}
	b, _ := json.Marshal(response.Response)
	// 打印返回的json字符串
	fmt.Printf("%s", b)

	/* 当出现以下错误码时，快速解决方案参考
	 * [FailedOperation.SignatureIncorrectOrUnapproved](https://cloud.tencent.com/document/product/382/9558#.E7.9F.AD.E4.BF.A1.E5.8F.91.E9.80.81.E6.8F.90.E7.A4.BA.EF.BC.9Afailedoperation.signatureincorrectorunapproved-.E5.A6.82.E4.BD.95.E5.A4.84.E7.90.86.EF.BC.9F)
	 * [FailedOperation.TemplateIncorrectOrUnapproved](https://cloud.tencent.com/document/product/382/9558#.E7.9F.AD.E4.BF.A1.E5.8F.91.E9.80.81.E6.8F.90.E7.A4.BA.EF.BC.9Afailedoperation.templateincorrectorunapproved-.E5.A6.82.E4.BD.95.E5.A4.84.E7.90.86.EF.BC.9F)
	 * [UnauthorizedOperation.SmsSdkAppIdVerifyFail](https://cloud.tencent.com/document/product/382/9558#.E7.9F.AD.E4.BF.A1.E5.8F.91.E9.80.81.E6.8F.90.E7.A4.BA.EF.BC.9Aunauthorizedoperation.smssdkappidverifyfail-.E5.A6.82.E4.BD.95.E5.A4.84.E7.90.86.EF.BC.9F)
	 * [UnsupportedOperation.ContainDomesticAndInternationalPhoneNumber](https://cloud.tencent.com/document/product/382/9558#.E7.9F.AD.E4.BF.A1.E5.8F.91.E9.80.81.E6.8F.90.E7.A4.BA.EF.BC.9Aunsupportedoperation.containdomesticandinternationalphonenumber-.E5.A6.82.E4.BD.95.E5.A4.84.E7.90.86.EF.BC.9F)
	 * 更多错误，可咨询[腾讯云助手](https://tccc.qcloud.com/web/im/index.html#/chat?webAppId=8fa15978f85cb41f7e2ea36920cb3ae1&title=Sms)
	 */
	return nil
}
