package oss

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/oss/common"
	"gitee.com/unitedrhino/share/utils"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"mime/multipart"
	"path"
	"strings"
	"time"
)

type (
	SceneInfo struct {
		Business string
		Scene    string
		FilePath string
		FileName string
	}
)

// 产品管理
const (
	BusinessDeviceManage = "deviceManage" //设备管理
	SceneDeviceImg       = "deviceImg"    //产品图片
	SceneFile            = "file"         //产品图片

	BusinessDeviceGroup = "deviceGroup"

	BusinessProductManage = "productManage"   //产品管理
	SceneProductImg       = "productImg"      //产品图片
	SceneProductCustomUi  = "productCustomUi" //产品自定义ui
	SceneCategoryImg      = "categoryImg"     //产品品类图片
)

const (
	BusinessProject    = "project"    //项目
	BusinessArea       = "area"       //区域
	BusinessUserManage = "userManage" //产品管理
	SceneHeadIng       = "headImg"    //头像
	SceneConfigFile    = "configFile"
)
const (
	BusinessTenantManage = "tenantManage"  //租户管理
	SceneBackgroundImg   = "backgroundImg" //
	SceneLogoImg         = "logoImg"
)
const (
	BusinessScene = "scene" //场景
)

const (
	SceneFirmware = "firmware"
	BusinessOta   = "ota"
)

const (
	BusinessApp = "app" //应用
)

func GetSceneInfo(filePath string) (*SceneInfo, error) {
	paths := strings.Split(filePath, "/")
	if len(paths) < 3 {
		return nil, errors.Parameter.WithMsg("路径不对")
	}
	scene := &SceneInfo{
		Business: paths[0],
		Scene:    paths[1],
		FilePath: strings.Join(paths[2:], "/"),
		FileName: paths[len(paths)-1],
	}
	return scene, nil
}

func IsFilePath(c conf.OssConf, filePath string) bool {
	if strings.HasPrefix(filePath, "http") || strings.HasPrefix(filePath, c.CustomPath) || strings.HasPrefix(filePath, c.CustomHost) {
		return false
	}
	return true
}

func GenFilePath(ctx context.Context, svrName, business, scene, filePath string) string {
	uc := ctxs.GetUserCtx(ctx)
	return fmt.Sprintf("%s/%s/%s/%s/%s/%s", svrName, uc.TenantCode, uc.AppCode, business, scene, filePath)
}

func GenCommonFilePath(svrName, business, scene, filePath string) string {
	return fmt.Sprintf("%s/%s/%s/%s", svrName, business, scene, filePath)
}

func IsCommonFile(svrName, business, scene string, filePath string) bool {
	part := fmt.Sprintf("%s/common/%s/%s", svrName, business, scene)
	return strings.Contains(filePath, part)
}

func GetFileNameWithPath(path string) string {
	fs := strings.Split(path, "/")
	return fs[len(fs)-1]
}

func SceneToNewPath(ctx context.Context, ossClient *Client, business, scene, filePath, oldFilePath, newFilePath string) (string, error) {
	si, err := GetSceneInfo(newFilePath)
	if err != nil {
		return "", err
	}
	if !(si.Business == business && si.Scene == scene && strings.HasPrefix(si.FilePath, filePath)) {
		return "", errors.Parameter.WithMsgf("图片的路径不对,路径要为/%s/%s/%s开头", business, scene, si.FilePath)
	}
	si.FilePath = filePath
	newPath, err := GetFilePath(si, false)
	if err != nil {
		return "", err
	}
	path, err := ossClient.PrivateBucket().CopyFromTempBucket(newFilePath, newPath)
	if err != nil {
		return "", errors.System.AddDetail(err)
	}
	err = ossClient.TemporaryBucket().Delete(ctx, newFilePath, common.OptionKv{})
	if err != nil {
		logx.WithContext(ctx).Error(err)
	}
	if oldFilePath != "" {
		err = ossClient.PrivateBucket().Delete(ctx, oldFilePath, common.OptionKv{})
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
	}
	return path, nil
}

func GetFilePath2(ctx context.Context, fh *multipart.FileHeader) (string, error) {
	fileName := fh.Filename
	spcChar := []string{`,`, `?`, `*`, `|`, `{`, `}`, `\`, `$`, `、`, `·`, "`", `'`, `"`}
	if strings.ContainsAny(fileName, strings.Join(spcChar, "")) {
		return "", errors.Parameter.WithMsg("包含特殊字符")
	}
	uc := ctxs.GetUserCtx(ctx)
	if uc == nil {
		return "", errors.Permissions.WithMsg("需要登录")
	}
	return fmt.Sprintf("%s/%s/%s/%d/%s/%s", utils.ToYYMMdd2(time.Now().UnixMilli()), uc.TenantCode, uc.AppCode, uc.UserID,
		utils.ToddHHSS2(time.Now().UnixMilli()), fileName), nil

}

func GetFilePath(scene *SceneInfo, rename bool) (string, error) {
	if rename == true {
		ext := path.Ext(scene.FilePath)
		if ext == "" {
			return "", errors.Parameter.WithMsg("未能获取文件后缀名")
		}
		uuid := uuid.NewString()
		scene.FilePath = uuid + ext
	} else {
		spcChar := []string{`,`, `?`, `*`, `|`, `{`, `}`, `\`, `$`, `、`, `·`, "`", `'`, `"`}
		if strings.ContainsAny(scene.FilePath, strings.Join(spcChar, "")) {
			return "", errors.Parameter.WithMsg("包含特殊字符")
		}
	}
	filePath := fmt.Sprintf("%s/%s/%s", scene.Business, scene.Scene, scene.FilePath)
	return filePath, nil
}

func CheckWithCopy(ctx context.Context, handle Handle, srcPath string, business, scene string) (string, error) {
	//如果第一次就提交了模型文件
	si, err := GetSceneInfo(srcPath)
	if err != nil {
		return "", err
	}
	if !(si.Business == business && si.Scene == scene) {
		return "", errors.Parameter.WithMsg(scene + "的路径不对")
	}
	nwePath, err := GetFilePath(si, false)
	if err != nil {
		return "", err
	}
	path, err := handle.CopyFromTempBucket(srcPath, nwePath)
	if err != nil {
		return "", errors.System.AddDetail(err)
	}
	return path, nil
}
