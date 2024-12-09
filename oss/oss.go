package oss

import (
	"context"
	"io"

	"gitee.com/unitedrhino/share/oss/common"
)

type Handle interface {
	//提供带签名的上传url
	SignedPutUrl(ctx context.Context, filePath string, expiredSec int64, opKv common.OptionKv) (string, error)
	//提供带签名的下载url
	SignedGetUrl(ctx context.Context, filePath string, expiredSec int64, opKv common.OptionKv) (string, error)
	//删除文件
	Delete(ctx context.Context, filePath string, opKv common.OptionKv) error
	//上传本地文件
	Upload(ctx context.Context, filePath string, reader io.Reader, opKv common.OptionKv) (string, error)
	//获取文件信息
	GetObjectInfo(ctx context.Context, filePath string) (*common.StorageObjectInfo, error)
	//将文件下载到本地
	GetObjectLocal(ctx context.Context, filePath string, localPath string) error
	//列出符合前缀的文件列表
	ListObjects(ctx context.Context, prefix string) (ret []*common.StorageObjectInfo, err error)
	//获取私有桶
	PrivateBucket() Handle
	//获取公开桶
	PublicBucket() Handle
	//获取临时桶
	TemporaryBucket() Handle
	//从临时桶中将文件复制到选定桶的路径
	CopyFromTempBucket(tempPath, dstPath string) (string, error)
	//不带过期时间的获取文件url
	GetUrl(path string, withHost bool) (string, error)
}
