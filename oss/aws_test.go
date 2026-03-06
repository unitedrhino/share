package oss

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/oss/common"
)

var (
	bucket    = "tier0-upload-temp"
	key       = "1.jpg"
	accessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	endpoint  = "s3.ap-southeast-1.amazonaws.com"
	region    = "ap-southeast-1"

	timeout time.Duration = time.Second * 500
)

func newTestAwsClient(t *testing.T) *Aws {
	t.Helper()
	ak := os.Getenv("AWS_ACCESS_KEY_ID")
	sk := os.Getenv("AWS_SECRET_ACCESS_KEY")
	ep := os.Getenv("S3_ENDPOINT")
	rg := os.Getenv("S3_REGION")
	bk := os.Getenv("S3_BUCKET")
	if ak == "" || sk == "" {
		t.Skip("跳过：未设置 AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY 环境变量")
	}
	if ep == "" {
		ep = "s3.amazonaws.com"
	}
	if rg == "" {
		rg = "ap-southeast-1"
	}
	if bk == "" {
		bk = "tier0-upload-temp-pre"
	}
	a, err := newAws(conf.AwsConf{OssConf: conf.OssConf{
		OssType:             "aws",
		AccessKeyID:         ak,
		AccessKeySecret:     sk,
		TemporaryBucketName: bk,
		Location:            ep,
		Region:              rg,
		UseSSL:              true,
	}})
	if err != nil {
		t.Fatalf("newAws 失败: %v", err)
	}
	return a
}

// TestSignedPutUrl_StaticCreds 测试静态凭证生成预签名上传 URL
func TestSignedPutUrl_StaticCreds(t *testing.T) {
	a := newTestAwsClient(t)
	ctx := context.Background()

	signedURL, err := a.TemporaryBucket().SignedPutUrl(ctx, "test/unit-test-upload.txt", 3600, common.OptionKv{})
	if err != nil {
		t.Fatalf("SignedPutUrl 失败: %v", err)
	}
	t.Logf("预签名上传 URL:\n%s", signedURL)

	// 验证 URL 包含 X-Amz-Credential（静态凭证时 AK 以 AKIA 开头）
	if !strings.Contains(signedURL, "X-Amz-Credential") {
		t.Error("URL 中未找到 X-Amz-Credential 参数")
	}
	if strings.Contains(signedURL, "AKIA") {
		t.Log("⚠️  当前使用永久凭证（AKIA），建议配置 RoleArn 切换为临时凭证")
	} else if strings.Contains(signedURL, "ASIA") {
		t.Log("✅ 当前使用临时凭证（ASIA），AK 会自动过期")
	}
}

// TestSignedGetUrl_StaticCreds 测试静态凭证生成预签名下载 URL
func TestSignedGetUrl_StaticCreds(t *testing.T) {
	a := newTestAwsClient(t)
	ctx := context.Background()

	signedURL, err := a.TemporaryBucket().SignedGetUrl(ctx, "test/unit-test-upload.txt", 3600, common.OptionKv{})
	if err != nil {
		t.Fatalf("SignedGetUrl 失败: %v", err)
	}
	t.Logf("预签名下载 URL:\n%s", signedURL)

	if !strings.Contains(signedURL, "X-Amz-Credential") {
		t.Error("URL 中未找到 X-Amz-Credential 参数")
	}
}

// TestSignedPutUrl_STS 测试 STS AssumeRole 临时凭证生成预签名上传 URL
func TestSignedPutUrl_STS(t *testing.T) {
	roleArn := os.Getenv("OssRoleArn")
	if roleArn == "" {
		t.Skip("跳过：未设置 AWS_ROLE_ARN 环境变量")
	}

	ak := os.Getenv("AWS_ACCESS_KEY_ID")
	sk := os.Getenv("AWS_SECRET_ACCESS_KEY")
	ep := os.Getenv("S3_ENDPOINT")
	rg := os.Getenv("S3_REGION")
	bk := os.Getenv("S3_BUCKET")
	if ep == "" {
		ep = "s3.amazonaws.com"
	}
	if rg == "" {
		rg = "ap-southeast-1"
	}
	if bk == "" {
		bk = "tier0-upload-temp-pre"
	}

	a, err := newAws(conf.AwsConf{
		OssConf: conf.OssConf{
			OssType:             "aws",
			AccessKeyID:         ak,
			AccessKeySecret:     sk,
			TemporaryBucketName: bk,
			Location:            ep,
			Region:              rg,
			UseSSL:              true,
		},
		RoleArn:         roleArn,
		RoleSessionName: "unit-test-session",
		SessionDuration: 900,
	})
	if err != nil {
		t.Fatalf("newAws（STS 模式）失败: %v", err)
	}

	ctx := context.Background()
	signedURL, err := a.TemporaryBucket().SignedPutUrl(ctx, "test/unit-test-sts.txt", 3600, common.OptionKv{})
	if err != nil {
		t.Fatalf("SignedPutUrl（STS）失败: %v", err)
	}
	t.Logf("STS 预签名上传 URL:\n%s", signedURL)

	if !strings.Contains(signedURL, "ASIA") {
		t.Errorf("STS 模式下 URL 应包含临时 AK（ASIA 开头），实际 URL: %s", signedURL)
	} else {
		t.Log("✅ STS 模式验证通过：URL 使用临时凭证（ASIA），AK 会自动过期")
	}
}

func TestStart(t *testing.T) {
	// 保留空测试入口，具体用例见上方各 TestXxx 函数
}
