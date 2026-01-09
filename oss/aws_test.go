package oss

import (
	"os"
	"testing"
	"time"
)

var (
	bucket   = "tier0-upload-temp"
	key      = "1.jpg"
	// 从环境变量读取凭据，避免硬编码敏感信息
	accessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	endpoint  = "s3.ap-southeast-1.amazonaws.com"
	region    = "ap-southeast-1"

	timeout time.Duration = time.Second * 500
)

func TestStart(t *testing.T) {
	//Start()
	//a, err := newAws(conf.AwsConf{OssConf: conf.OssConf{
	//	OssType:             "aws",
	//	AccessKeyID:         accessKey,
	//	AccessKeySecret:     secretKey,
	//	PublicBucketName:    "tier0-upload-public",
	//	TemporaryBucketName: "tier0-upload-temp",
	//	PrivateBucketName:   "tier0-upload-private",
	//	Location:            endpoint,
	//	Region:              region,
	//	UseSSL:              true,
	//}})
	//if err != nil {
	//	t.Error(err)
	//}
	//ctx := context.Background()
	//objs, err := a.PrivateBucket().ListObjects(ctx, "")
	//if err != nil {
	//	t.Error(err)
	//}
	//if len(objs) != 0 {
	//	i, err := a.PrivateBucket().GetObjectInfo(ctx, objs[0].FilePath)
	//	if err != nil {
	//		t.Error(err)
	//	}
	//	t.Log(i)
	//}
	//var filePath = "aa/bb"
	//a.PrivateBucket().CopyFromTempBucket(filePath, filePath)
	//url, err := a.PrivateBucket().SignedPutUrl(ctx, filePath, 1000, common.OptionKv{})
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log(url)
	//i, err := a.PrivateBucket().GetObjectInfo(ctx, filePath)
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log(i)
	//err = a.PrivateBucket().GetObjectLocal(ctx, filePath, "./fff")
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log(i)
	//f, err := os.Open("./fff")
	//a.PrivateBucket().Upload(ctx, filePath+"/fefe", f, common.OptionKv{})
	//url, err = a.PrivateBucket().SignedGetUrl(ctx, filePath, 1000, common.OptionKv{})
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log(url)
	//url, err = a.PublicBucket().GetUrl(filePath, false)
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Run("123", func(t *testing.T) {
	//	Start()
	//})
}
