package oss

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"gitee.com/unitedrhino/share/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/oss/common"
)

type Aws struct {
	setting           conf.OssConf
	currentBucketName string
	cli               *s3.Client
}

func (m *Aws) IsFilePath(filePath string) bool {
	return isFilePath(m.setting, filePath)
}

func (m *Aws) IsFileUrl(url string) bool {
	return isFileUrl(m.setting, url)
}

func newAws(conf conf.AwsConf) (*Aws, error) {
	creds := credentials.NewStaticCredentialsProvider(conf.AccessKeyID, conf.AccessKeySecret, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(conf.Region),
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               conf.Location,
				HostnameImmutable: true,
				SigningRegion:     conf.Region,
			}, nil
		})),
	)
	if err != nil {
		return nil, err
	}
	svc := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	return &Aws{
		setting: conf.OssConf,
		cli:     svc,
	}, nil
}
func (m *Aws) PrivateBucket() Handle {
	m.currentBucketName = m.setting.PrivateBucketName
	return m
}
func (m *Aws) PublicBucket() Handle {
	m.currentBucketName = m.setting.PublicBucketName
	return m
}
func (m *Aws) TemporaryBucket() Handle {
	m.currentBucketName = m.setting.TemporaryBucketName
	return m
}

func (m *Aws) Bucket(name string) Handle {
	m.currentBucketName = name
	return m
}

// 获取put上传url
func (m *Aws) SignedPutUrl(ctx context.Context, filePath string, expiredSec int64, opKv common.OptionKv) (string, error) {
	presignClient := s3.NewPresignClient(m.cli)
	resp, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(m.currentBucketName),
		Key:    aws.String(filePath),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expiredSec) * time.Second
	})
	if err != nil {
		return "", err
	}
	return resp.URL, nil
}

// 获取get下载url
func (m *Aws) SignedGetUrl(ctx context.Context, filePath string, expiredSec int64, opKv common.OptionKv) (string, error) {
	presignClient := s3.NewPresignClient(m.cli)
	resp, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(m.currentBucketName),
		Key:    aws.String(filePath),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expiredSec) * time.Second
	})
	if err != nil {
		return "", err
	}
	return resp.URL, nil
}

// 删除
func (m *Aws) Delete(ctx context.Context, filePath string, opKv common.OptionKv) error {
	_, err := m.cli.DeleteObject(ctx, &s3.DeleteObjectInput{
		Key:    aws.String(filePath),
		Bucket: aws.String(m.currentBucketName),
	})
	return err
}

func (m *Aws) Upload(ctx context.Context, filePath string, reader io.Reader, opKv common.OptionKv) (string, error) {
	r, err := utils.ReaderToReadSeeker(reader)
	if err != nil {
		return "", err
	}
	_, err = m.cli.PutObject(ctx, &s3.PutObjectInput{
		Body:   r,
		Bucket: aws.String(m.currentBucketName),
		Key:    aws.String(filePath),
	})
	return m.SignedGetUrl(ctx, filePath, 1000, opKv)
}

func (m *Aws) GetObjectLocal(ctx context.Context, filePath string, localPath string) error {
	obj, err := m.cli.GetObject(ctx, &s3.GetObjectInput{
		Key:    aws.String(filePath),
		Bucket: aws.String(m.currentBucketName),
	})
	if err != nil {
		return err
	}
	defer obj.Body.Close()
	f, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, obj.Body)
	return err
}

func (m *Aws) GetObjectInfo(ctx context.Context, filePath string) (*common.StorageObjectInfo, error) {
	obj, err := m.cli.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(m.currentBucketName),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return nil, err
	}
	defer obj.Body.Close()
	return &common.StorageObjectInfo{
		FilePath: filePath,
		Size:     *obj.ContentLength,
		Md5:      *obj.ETag,
	}, nil
}

func (m *Aws) ListObjects(ctx context.Context, prefix string) (ret []*common.StorageObjectInfo, err error) {
	objs, err := m.cli.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(m.currentBucketName),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, err
	}
	for _, obj := range objs.Contents {
		ret = append(ret, &common.StorageObjectInfo{
			FilePath: *obj.Key,
			Size:     *obj.Size,
			Md5:      *obj.ETag,
		})
	}
	return
}

func (m *Aws) CopyFromTempBucket(tempPath, dstPath string) (string, error) {
	_, err := m.cli.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     aws.String(m.currentBucketName),
		CopySource: aws.String(fmt.Sprintf("%s/%s", m.setting.TemporaryBucketName, url.PathEscape(tempPath))),
		Key:        aws.String(dstPath),
	})
	return dstPath, err
}

// 获取完整链接
func (m *Aws) GetUrl(path string, withHost bool) (string, error) {
	if path[0] == '/' {
		path = path[1:]
	} //示例: https://tier0-upload-pub.s3.ap-southeast-1.amazonaws.com/screenshot-20251202-231438.png
	return "https://" + m.currentBucketName + "." + m.setting.Location + "/" + path, nil
}

func (m *Aws) FileUrlToFilePath(url string) (bucket string, filePath string) {
	return fileUrlToFilePath(m.setting, url)
}
