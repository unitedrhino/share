package oss

import (
	"context"
	"io"
	"strconv"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/oss/common"
	aliOss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliYunOss struct {
	setting         conf.AliYunConf
	client          *aliOss.Client
	bucket          *aliOss.Bucket
	b4IsObjectExist *aliOss.Bucket
}

func (a *AliYunOss) GetObjectLocal(ctx context.Context, filePath string, localPath string) error {
	//TODO implement me
	panic("implement me")
}

func newAliYunOss(conf conf.AliYunConf) (*AliYunOss, error) {
	client, err := aliOss.New(conf.GetEndPoint(), conf.AccessKeyID, conf.AccessKeySecret, conf.GenClientOption()...)
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(conf.PrivateBucketName)
	if err != nil {
		return nil, err
	}
	b4IsObjectExist := bucket
	return &AliYunOss{client: client, bucket: bucket, setting: conf, b4IsObjectExist: b4IsObjectExist}, nil
}

func (a *AliYunOss) Bucket(name string) Handle {
	a.bucket, _ = a.client.Bucket(name)
	return a
}
func (a *AliYunOss) PrivateBucket() Handle {
	a.bucket, _ = a.client.Bucket(a.setting.PrivateBucketName)
	return a
}
func (a *AliYunOss) PublicBucket() Handle {
	a.bucket, _ = a.client.Bucket(a.setting.PublicBucketName)
	return a
}
func (a *AliYunOss) TemporaryBucket() Handle {
	a.bucket, _ = a.client.Bucket(a.setting.PublicBucketName)
	return a
}

// 获取put上传url
func (a *AliYunOss) SignedPutUrl(ctx context.Context, fileDir string, expiredSec int64, opKv common.OptionKv) (string, error) {
	return a.bucket.SignURL(fileDir, aliOss.HTTPPut, expiredSec, opKv.ToAliYunOptions()...)
}
func (a *AliYunOss) SignedGetUrl(ctx context.Context, fileDir string, expiredSec int64, opKv common.OptionKv) (string, error) {
	return a.bucket.SignURL(fileDir, aliOss.HTTPGet, expiredSec, opKv.ToAliYunOptions()...)
}
func (a *AliYunOss) Delete(ctx context.Context, fileDir string, opKv common.OptionKv) error {
	return a.bucket.DeleteObject(fileDir, opKv.ToAliYunOptions()...)
}
func (a *AliYunOss) GetObjectInfo(ctx context.Context, fileDir string) (*common.StorageObjectInfo, error) {
	metaInfo, err := a.bucket.GetObjectDetailedMeta(fileDir)
	if err != nil {
		return nil, err
	}
	contentLength := metaInfo.Get("Content-Length")
	size, err := strconv.ParseInt(contentLength, 10, 64)

	return &common.StorageObjectInfo{
		Size: size,
		Md5:  "",
	}, err
}

func (a *AliYunOss) ListObjects(ctx context.Context, prefix string) (ret []*common.StorageObjectInfo, err error) {
	//todo
	return
}

func (a *AliYunOss) Upload(ctx context.Context, filePath string, content io.Reader, opKv common.OptionKv) (string, error) {
	err := a.bucket.PutObject(filePath, content, opKv.ToAliYunOptions()...)
	//fmt.Println(uploadInfo)
	return filePath, err
}
func (a *AliYunOss) CopyFromTempBucket(tempPath, dstPath string) (string, error) {
	return "", nil //TODO
}
func (a *AliYunOss) GetUrl(filePath string, withHost bool) (string, error) {
	return "", nil //TODO
}

func (m *AliYunOss) IsFilePath(filePath string) bool {
	return isFilePath(m.setting.OssConf, filePath)
}

func (m *AliYunOss) IsFileUrl(url string) bool {
	return isFileUrl(m.setting.OssConf, url)
}

func (m *AliYunOss) FileUrlToFilePath(url string) (bucket string, filePath string) {
	return fileUrlToFilePath(m.setting.OssConf, url)
}
