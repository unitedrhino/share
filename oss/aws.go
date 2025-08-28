package oss

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"gitee.com/unitedrhino/share/utils"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/oss/common"
)

type Aws struct {
	setting           conf.OssConf
	currentBucketName string
	cli               *s3.S3
}

func (m *Aws) IsFilePath(filePath string) bool {
	return isFilePath(m.setting, filePath)
}

func (m *Aws) IsFileUrl(url string) bool {
	return isFileUrl(m.setting, url)
}

func newAws(conf conf.AwsConf) (*Aws, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(conf.AccessKeyID, conf.AccessKeySecret, ""),
		Endpoint:         &conf.Location,
		Region:           &conf.Region,
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	svc := s3.New(sess)
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
	req, _ := m.cli.PutObjectRequest(&s3.PutObjectInput{
		Bucket: &m.currentBucketName,
		Key:    &filePath,
	})
	url, err := req.Presign(time.Duration(expiredSec) * time.Second)
	if err != nil {
		return "", err
	}
	return url, nil
}

// 获取get下载url
func (m *Aws) SignedGetUrl(ctx context.Context, filePath string, expiredSec int64, opKv common.OptionKv) (string, error) {
	req, _ := m.cli.GetObjectRequest(&s3.GetObjectInput{Key: &filePath, Bucket: &m.currentBucketName})
	url, err := req.Presign(time.Duration(expiredSec) * time.Second)
	return url, err
}

// 删除
func (m *Aws) Delete(ctx context.Context, filePath string, opKv common.OptionKv) error {
	_, err := m.cli.DeleteObject(&s3.DeleteObjectInput{Key: &filePath, Bucket: &m.currentBucketName})
	return err
}

func (m *Aws) Upload(ctx context.Context, filePath string, reader io.Reader, opKv common.OptionKv) (string, error) {
	r, err := utils.ReaderToReadSeeker(reader)
	if err != nil {
		return "", err
	}
	_, err = m.cli.PutObject(&s3.PutObjectInput{
		Body:   r,
		Bucket: &m.currentBucketName,
		Key:    &filePath,
	})
	return m.SignedGetUrl(ctx, filePath, 1000, opKv)
}

func (m *Aws) GetObjectLocal(ctx context.Context, filePath string, localPath string) error {
	obj, err := m.cli.GetObjectWithContext(ctx, &s3.GetObjectInput{Key: &filePath, Bucket: &m.currentBucketName})
	if err != nil {
		return err
	}
	f, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, obj.Body)
	return err
}

func (m *Aws) GetObjectInfo(ctx context.Context, filePath string) (*common.StorageObjectInfo, error) {
	obj, err := m.cli.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: &m.currentBucketName,
		Key:    &filePath,
	})
	if err != nil {
		return nil, err
	}
	return &common.StorageObjectInfo{
		FilePath: filePath,
		Size:     *obj.ContentLength,
		Md5:      *obj.ETag,
	}, nil
}

func (m *Aws) ListObjects(ctx context.Context, prefix string) (ret []*common.StorageObjectInfo, err error) {
	objs, err := m.cli.ListObjectsWithContext(ctx, &s3.ListObjectsInput{
		Bucket: &m.currentBucketName,
		Prefix: tea.String(prefix),
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
	_, err := m.cli.CopyObject(&s3.CopyObjectInput{
		Bucket:     &m.currentBucketName,
		CopySource: aws.String(fmt.Sprintf("%s/%s", m.setting.TemporaryBucketName, url.PathEscape(tempPath))),
		Key:        &dstPath,
	})
	return dstPath, err
}

// 获取完整链接
func (m *Aws) GetUrl(path string, withHost bool) (string, error) {
	if path[0] == '/' {
		path = path[1:]
	}
	if withHost {
		return m.setting.CustomHost + "/" + m.currentBucketName + "/" + path, nil
	}
	return m.setting.CustomPath + "/" + m.currentBucketName + "/" + path, nil
}

func (m *Aws) FileUrlToFilePath(url string) (bucket string, filePath string) {
	return fileUrlToFilePath(m.setting, url)
}
