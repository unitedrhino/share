package oss

import (
	"context"
	"io"

	"gitee.com/unitedrhino/share/oss/common"
)

type Handle interface {
	SignedPutUrl(ctx context.Context, filePath string, expiredSec int64, opKv common.OptionKv) (string, error)
	SignedGetUrl(ctx context.Context, filePath string, expiredSec int64, opKv common.OptionKv) (string, error)
	Delete(ctx context.Context, filePath string, opKv common.OptionKv) error
	Upload(ctx context.Context, filePath string, reader io.Reader, opKv common.OptionKv) (string, error)
	GetObjectInfo(ctx context.Context, filePath string) (*common.StorageObjectInfo, error)
	GetObjectLocal(ctx context.Context, filePath string, localPath string) error
	ListObjects(ctx context.Context, prefix string) (ret []*common.StorageObjectInfo, err error)
	PrivateBucket() Handle
	PublicBucket() Handle
	TemporaryBucket() Handle
	CopyFromTempBucket(tempPath, dstPath string) (string, error)
	GetUrl(path string, withHost bool) (string, error)
	//List(ctx context.Context)
}
