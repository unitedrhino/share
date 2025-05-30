package oss

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"net/url"
	"time"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/oss/common"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	setting           conf.OssConf
	client            *minio.Client
	core              *minio.Core
	currentBucketName string
}

func (m *Minio) IsFilePath(filePath string) bool {
	return isFilePath(m.setting, filePath)
}

func (m *Minio) IsFileUrl(url string) bool {
	return isFileUrl(m.setting, url)
}

func newMinio(conf conf.MinioConf) (*Minio, error) {
	// 初使化 minio client对象。
	minioClient, err := minio.New(conf.GetEndPoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKeyID, conf.AccessKeySecret, ""),
		Secure: conf.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	core, err := minio.NewCore(conf.GetEndPoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKeyID, conf.AccessKeySecret, ""),
		Secure: conf.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	minioC := &Minio{
		setting: conf.OssConf,
		client:  minioClient,
		core:    core,
	}
	logx.Must(minioC.initPrivatePolicy())   //私有桶
	logx.Must(minioC.initPublicPolicy())    //公共读桶
	logx.Must(minioC.initTemporaryPolicy()) //临时桶
	return minioC, nil
}
func (m *Minio) PrivateBucket() Handle {
	m.currentBucketName = m.setting.PrivateBucketName
	return m
}
func (m *Minio) PublicBucket() Handle {
	m.currentBucketName = m.setting.PublicBucketName
	return m
}
func (m *Minio) TemporaryBucket() Handle {
	m.currentBucketName = m.setting.TemporaryBucketName
	return m
}

func (m *Minio) Bucket(name string) Handle {
	m.currentBucketName = name
	return m
}

// 获取put上传url
func (m *Minio) SignedPutUrl(ctx context.Context, fileDir string, expiredSec int64, opKv common.OptionKv) (string, error) {
	if err := m.checkForbidOverwrite(ctx, fileDir, opKv); err != nil {
		return "", err
	}
	url, err := m.client.PresignedPutObject(ctx, m.currentBucketName, fileDir, time.Duration(expiredSec*int64(time.Second)))
	if err != nil {
		return "", err
	}
	return m.setting.CustomPath + url.RequestURI(), err
}

// 获取get下载url
func (m *Minio) SignedGetUrl(ctx context.Context, filePath string, expiredSec int64, opKv common.OptionKv) (string, error) {
	url, err := m.client.PresignedGetObject(ctx, m.currentBucketName, filePath, time.Duration(expiredSec*int64(time.Second)), opKv.ToMinioReqParams())
	if err != nil {
		return "", err
	}
	return m.setting.CustomPath + url.RequestURI(), err
}

// 删除
func (m *Minio) Delete(ctx context.Context, filePath string, opKv common.OptionKv) error {
	return m.client.RemoveObject(ctx, m.currentBucketName, filePath, minio.RemoveObjectOptions{})
}

// 重名文件检查
func (m *Minio) checkForbidOverwrite(ctx context.Context, filePath string, opKv common.OptionKv) error {
	if opKv.IsForbidOverwrite() {
		ok, err := m.IsObjectExist(ctx, filePath, opKv)
		if err != nil {
			return err
		}
		if ok {
			return common.ForbidWriteErr
		}
	}
	return nil
}

func (m *Minio) IsObjectExist(ctx context.Context, filePath string, opKv common.OptionKv) (bool, error) {
	_, err := m.client.StatObject(ctx, m.currentBucketName, filePath, minio.StatObjectOptions{})
	if err == nil {
		return true, nil
	}
	switch err.(type) {
	case minio.ErrorResponse:
		if err.(minio.ErrorResponse).StatusCode == http.StatusNotFound {
			return false, nil
		}
	}
	return false, err
}
func (m *Minio) Upload(ctx context.Context, filePath string, reader io.Reader, opKv common.OptionKv) (string, error) {
	uploadInfo, err := m.client.PutObject(ctx, m.currentBucketName, filePath, reader, -1, minio.PutObjectOptions{ContentType: common.GetFilePathMineType(filePath), PartSize: 1024 * 1024 * 16})
	uri, _ := url.Parse(uploadInfo.Location)
	return m.setting.CustomPath + uri.RequestURI(), err
}

func (m *Minio) GetObjectLocal(ctx context.Context, filePath string, localPath string) error {
	err := m.client.FGetObject(ctx, m.currentBucketName, filePath, localPath, minio.GetObjectOptions{})
	return err
}

func (m *Minio) GetObjectInfo(ctx context.Context, filePath string) (*common.StorageObjectInfo, error) {
	object, err := m.client.GetObject(ctx, m.currentBucketName, filePath, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	objectInfo, err := object.Stat()
	return &common.StorageObjectInfo{
		FilePath: objectInfo.Key,
		Size:     objectInfo.Size,
		Md5:      objectInfo.ETag,
	}, err
}

func (m *Minio) ListObjects(ctx context.Context, prefix string) (ret []*common.StorageObjectInfo, err error) {
	objs := m.client.ListObjects(ctx, m.currentBucketName, minio.ListObjectsOptions{Prefix: prefix})
	for obj := range objs {
		ret = append(ret, &common.StorageObjectInfo{
			FilePath: obj.Key,
			Size:     obj.Size,
			Md5:      obj.ETag,
		})
	}
	return
}

func (m *Minio) initPrivatePolicy() error {
	if exists, err := m.client.BucketExists(context.Background(), m.setting.PrivateBucketName); err != nil {
		return err
	} else if !exists {
		err = m.client.MakeBucket(context.Background(), m.setting.PrivateBucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
func (m *Minio) initPublicPolicy() error {
	if exists, err := m.client.BucketExists(context.Background(), m.setting.PublicBucketName); err != nil {
		return err
	} else if !exists {
		err = m.client.MakeBucket(context.Background(), m.setting.PublicBucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}
	publicProcyString := `{
		"Version":"2012-10-17",
		"Statement":[
		  {
			"Effect":"Allow",
			"Principal":{
			  "AWS":["*"]
			},
			"Action":[
			  "s3:GetBucketLocation"
			],
			"Resource":[
              "arn:aws:s3:::` + m.setting.PublicBucketName + `"
			]
		  },
		  {
			"Effect":"Allow",
			"Principal":{
			  "AWS":["*"]
			},
			"Action":[
			  "s3:GetObject"
			],
			"Resource":[
			  "arn:aws:s3:::` + m.setting.PublicBucketName + `/*"
			]
		  }
		]
	  }`
	err := m.client.SetBucketPolicy(context.Background(), m.setting.PublicBucketName, publicProcyString)
	return err
}
func (m *Minio) initTemporaryPolicy() error {
	// TODO: 临时桶的自动过期删除文件的策略还没有实现
	// rule := `{
	// 	"Rules": [
	// 		{
	// 			"Expiration": {
	// 				"Date": "2020-04-07T02:00:00.000Z"
	// 			},
	// 			"ID": "Delete very old messenger pictures",
	// 			"Filter": {
	// 				"Prefix": "uploads/2015/"
	// 			},
	// 			"Msg": "Enabled"
	// 		},
	// 		{
	// 			"Expiration": {
	// 				"Days": 7
	// 			},
	// 			"ID": "Delete temporary uploads",
	// 			"Filter": {
	// 				"Prefix": "temporary-uploads/"
	// 			},
	// 			"Msg": "Enabled"
	// 		}
	// 	]
	// }`
	if exists, err := m.client.BucketExists(context.Background(), m.setting.TemporaryBucketName); err != nil {
		return err
	} else if !exists {
		err = m.client.MakeBucket(context.Background(), m.setting.TemporaryBucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		publicProcyString := `{
		"Version":"2012-10-17",
		"Statement":[
		  {
			"Effect":"Allow",
			"Principal":{
			  "AWS":["*"]
			},
			"Action":[
			  "s3:GetBucketLocation"
			],
			"Resource":[
              "arn:aws:s3:::` + m.setting.TemporaryBucketName + `"
			]
		  },
		  {
			"Effect":"Allow",
			"Principal":{
			  "AWS":["*"]
			},
			"Action":[
			  "s3:GetObject"
			],
			"Resource":[
			  "arn:aws:s3:::` + m.setting.TemporaryBucketName + `/*"
			]
		  }
		]
	  }`
		err := m.client.SetBucketPolicy(context.Background(), m.setting.TemporaryBucketName, publicProcyString)
		return err
	}
	return nil
}

func (m *Minio) CopyFromTempBucket(tempPath, dstPath string) (string, error) {
	src := minio.CopySrcOptions{
		Bucket: m.setting.TemporaryBucketName,
		Object: tempPath,
	}
	dst := minio.CopyDestOptions{
		Bucket: m.currentBucketName,
		Object: dstPath,
	}

	ui, err := m.client.CopyObject(context.Background(), dst, src)
	if err != nil {
		return "", err
	}
	return ui.Key, err
}

// 获取完整链接
func (m *Minio) GetUrl(path string, withHost bool) (string, error) {
	if path[0] == '/' {
		path = path[1:]
	}
	if withHost {
		return m.setting.CustomHost + "/" + m.currentBucketName + "/" + path, nil
	}
	return m.setting.CustomPath + "/" + m.currentBucketName + "/" + path, nil
}

func (m *Minio) FileUrlToFilePath(url string) (bucket string, filePath string) {
	return fileUrlToFilePath(m.setting, url)
}
