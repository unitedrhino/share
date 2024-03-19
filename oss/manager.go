package oss

import (
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/def"
)

func newOssManager(setting conf.OssConf) (sm Handle, err error) {
	OssType := setting.OssType
	switch OssType {
	case def.OssAliyun:
		sm, err = newAliYunOss(conf.AliYunConf{OssConf: setting})
	case def.OssMinio:
		sm, err = newMinio(conf.MinioConf{OssConf: setting})
	}
	return sm, err
}
