package oss

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/oss/common"
	"github.com/google/uuid"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Local struct {
	setting           conf.OssConf
	currentBucketName string
}

func newLocal(conf conf.OssConf) (*Local, error) {
	conf.CustomPath = "/api/v1/system/common/download-file"
	_, err := os.Stat(conf.StorePath)
	if err != nil {
		err = os.Mkdir(conf.StorePath, os.ModeDir)
		if err != nil {
			return nil, err
		}
		err = os.Mkdir(fmt.Sprintf("%s%c%s", conf.StorePath, os.PathSeparator, conf.PublicBucketName), os.ModeDir)
		if err != nil {
			return nil, err
		}
		err = os.Mkdir(fmt.Sprintf("%s%c%s", conf.StorePath, os.PathSeparator, conf.PrivateBucketName), os.ModeDir)
		if err != nil {
			return nil, err
		}
		err = os.Mkdir(fmt.Sprintf("%s%c%s", conf.StorePath, os.PathSeparator, conf.TemporaryBucketName), os.ModeDir)
		if err != nil {
			return nil, err
		}
	}
	return &Local{setting: conf}, nil
}
func (m *Local) PrivateBucket() Handle {
	m.currentBucketName = m.setting.PrivateBucketName
	return m
}
func (m *Local) PublicBucket() Handle {
	m.currentBucketName = m.setting.PublicBucketName
	return m
}
func (m *Local) TemporaryBucket() Handle {
	m.currentBucketName = m.setting.TemporaryBucketName
	return m
}

func (m *Local) Bucket(name string) Handle {
	m.currentBucketName = name
	return m
}

// 获取put上传url
func (m *Local) SignedPutUrl(ctx context.Context, fileDir string, expiredSec int64, opKv common.OptionKv) (string, error) {
	//不使用
	//if err := m.checkForbidOverwrite(ctx, fileDir, opKv); err != nil {
	//	return "", err
	//}
	//url, err := m.client.PresignedPutObject(ctx, m.currentBucketName, fileDir, time.Duration(expiredSec*int64(time.Second)))
	//if err != nil {
	//	return "", err
	//}
	return m.setting.CustomPath, nil
}

type CacheBody struct {
	FilePath string
}

func (m *Local) fmtKey(uuid string) string {
	return fmt.Sprintf("oss:local:%s:%s", m.currentBucketName, uuid)
}

// 获取get下载url
func (m *Local) SignedGetUrl(ctx context.Context, filePath string, expiredSec int64, opKv common.OptionKv) (string, error) {
	path := fmt.Sprintf("%s/%s", m.currentBucketName, filePath)
	params := url.Values{}
	params.Add("filePath", path)
	id := uuid.NewString()
	params.Add("sign", id)
	err := caches.GetStore().SetexCtx(ctx, m.fmtKey(id), path, int(expiredSec))
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s?%s", m.setting.CustomPath, params.Encode())
	return url, nil
}

// 删除
func (m *Local) Delete(ctx context.Context, filePath string, opKv common.OptionKv) error {
	path := fmt.Sprintf("%s/%s/%s", m.setting.StorePath, m.currentBucketName, filePath)

	if _, err := os.Stat(path); err == nil {
		err := os.Remove(path)
		if err != nil {
			return errors.System.AddMsg("删除失败").AddDetail(err)
		}
	}
	return nil
}

func (m *Local) IsObjectExist(ctx context.Context, filePath string, opKv common.OptionKv) (bool, error) {
	path := fmt.Sprintf("%s/%s/%s", m.setting.StorePath, m.currentBucketName, filePath)
	if _, err := os.Stat(path); err == nil {
		return true, nil
	}
	return false, nil
}

func (m *Local) Upload(ctx context.Context, filePath string, reader io.Reader, opKv common.OptionKv) (string, error) {
	path := fmt.Sprintf("%s/%s/%s", m.setting.StorePath, m.currentBucketName, filePath)
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", errors.System.AddDetailf("无法创建目录: %v", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return "", errors.System.AddDetailf("无法打开文件: %v", err)
	}
	// 确保文件在函数结束时关闭
	defer file.Close()

	// 将 reader 中的数据复制到文件中
	_, err = io.Copy(file, reader)
	if err != nil {
		return "", errors.System.AddDetailf("写入文件时出错: %v", err)
	}
	return m.GetUrl(filePath, false)
}

func (m *Local) DownloadFile(ctx context.Context, filePath string, sign string, w http.ResponseWriter) error {
	if !strings.HasPrefix(filePath, m.setting.TemporaryBucketName) {
		if sign == "" {
			return errors.Permissions
		}
		bucket, _, _ := strings.Cut(filePath, "/")
		path, err := caches.GetStore().Get(fmt.Sprintf("oss:local:%s:%s", bucket, sign))
		if err != nil {
			return errors.Permissions.AddDetail(err)
		}
		if path != filePath {
			return errors.Permissions
		}
	}
	// 打开文件
	file, err := os.Open(fmt.Sprintf(m.setting.StorePath + "/" + filePath))
	if err != nil {
		return errors.System.AddMsgf("文件不存在").AddDetail(err)
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return errors.System.AddMsgf("无法获取文件信息").AddDetail(err)
	}
	// 设置响应头
	// Content-Disposition 用于指定文件名
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.Name()))
	c := mime.TypeByExtension(path.Ext(file.Name()))
	if c != "" {
		w.Header().Set("Content-Type", c)
	} else {
		// 只需要读取前 512 个字节来判断内容类型
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil && err != io.EOF {
			return errors.System.AddMsgf("读取文件内容失败").AddDetail(err)
		}
		// 重置文件读取位置，以便后续操作
		_, err = file.Seek(0, 0)
		if err != nil {
			return errors.System.AddMsgf("重置文件指针失败").AddDetail(err)
		}
		// Content-Type 可以根据文件类型设置，这里使用通用的二进制类型
		w.Header().Set("Content-Type", http.DetectContentType(buffer))
	}

	// Content-Length 用于指定文件大小
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	_, err = io.Copy(w, file)
	if err != nil {
		return errors.System.AddMsgf("文件下载出错").AddDetail(err)
	}
	return nil
}

func (m *Local) GetObjectLocal(ctx context.Context, filePath string, localPath string) error {
	sourceFile, err := os.Open(fmt.Sprintf(m.setting.StorePath + "/" + m.currentBucketName + "/" + filePath))
	if err != nil {
		return errors.System.AddMsgf("文件不存在").AddDetail(err)
	}
	defer sourceFile.Close()
	dir := filepath.Dir(localPath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return errors.System.AddDetailf("无法创建目录: %v", err)
	}
	// 创建目标文件
	destinationFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 使用 io.Copy 复制文件内容
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	// 同步文件内容到磁盘
	err = destinationFile.Sync()
	if err != nil {
		return err
	}

	return err
}

func (m *Local) GetObjectInfo(ctx context.Context, filePath string) (*common.StorageObjectInfo, error) {
	// 打开文件
	file, err := os.Open(fmt.Sprintf(m.setting.StorePath + "/" + m.currentBucketName + "/" + filePath))
	if err != nil {
		return nil, errors.System.AddMsgf("文件不存在").AddDetail(err)
	}
	defer file.Close()

	// 创建一个新的 MD5 哈希对象
	hash := md5.New()

	// 将文件内容复制到哈希对象中
	if _, err := io.Copy(hash, file); err != nil {
		return nil, err
	}

	// 获取哈希值的字节切片
	hashInBytes := hash.Sum(nil)

	// 将字节切片转换为十六进制字符串
	hashString := hex.EncodeToString(hashInBytes)
	fileStat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return &common.StorageObjectInfo{
		FilePath: filePath,
		Size:     fileStat.Size(),
		Md5:      hashString,
	}, err
}

func (m *Local) ListObjects(ctx context.Context, prefix string) (ret []*common.StorageObjectInfo, err error) {
	//objs := m.client.ListObjects(ctx, m.currentBucketName, minio.ListObjectsOptions{Prefix: prefix})
	//for obj := range objs {
	//	ret = append(ret, &common.StorageObjectInfo{
	//		FilePath: obj.Key,
	//		Size:     obj.Size,
	//		Md5:      obj.ETag,
	//	})
	//}
	return
}

func (m *Local) CopyFromTempBucket(tempPath, dstPath string) (string, error) {
	bucket := m.currentBucketName
	err := m.TemporaryBucket().GetObjectLocal(context.Background(), tempPath,
		fmt.Sprintf(m.setting.StorePath+"/"+bucket+"/"+dstPath))
	if err != nil {
		return "", err
	}
	return dstPath, nil
}

// 获取完整链接
func (m *Local) GetUrl(path string, withHost bool) (string, error) {
	if path[0] == '/' {
		path = path[1:]
	}
	path = fmt.Sprintf("%s/%s", m.currentBucketName, path)
	params := url.Values{}
	params.Add("filePath", path)
	if withHost {
		url := fmt.Sprintf("%s?%s", m.setting.CustomPath, params.Encode())
		return m.setting.CustomHost + "/" + url, nil
	}
	url := fmt.Sprintf("%s?%s", m.setting.CustomPath, params.Encode())
	return url, nil
}

func (m *Local) IsFilePath(filePath string) bool {
	return isFilePath(m.setting, filePath)
}

func (m *Local) IsFileUrl(url string) bool {
	return isFileUrl(m.setting, url)
}
func (m *Local) FileUrlToFilePath(url string) (bucket string, filePath string) {
	return fileUrlToFilePath(m.setting, url)
}
