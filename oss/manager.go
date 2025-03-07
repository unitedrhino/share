package oss

import (
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/def"
)

func newOssManager(setting conf.OssConf) (sm Handle, err error) {
	OssType := setting.OssType
	switch OssType {
	case def.OssAliyun:
		sm, err = newAliYunOss(conf.AliYunConf{OssConf: setting})
	case def.OssMinio:
		sm, err = newMinio(conf.MinioConf{OssConf: setting})
	case def.OssLocal:
		sm, err = newLocal(setting)
	}
	return sm, err
}
