package conf

import (
	"fmt"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssConf struct {
	OssType             string `json:",default=minio,options=minio|aliyun|local"` //oss的类型
	AccessKeyID         string `json:",default=root,optional"`                    //账号
	AccessKeySecret     string `json:",default=password,optional"`                //密码
	PublicBucketName    string `json:",default=ithings-public,optional"`          //公开桶的名称
	TemporaryBucketName string `json:",default=ithings-temporary,optional"`       //临时桶,30分钟有效期
	PrivateBucketName   string `json:",default=ithings-private,optional"`         //私有桶的名称
	Location            string `json:",default=localhost:9000,optional"`          // oss的地址
	UseSSL              bool   `json:",optional"`                                 //是否使用ssl
	CustomHost          string `json:",default=/oss,env=OssCustomHost"`           //带上host的返回前缀,支持环境变量,如:http://127.0.0.1:7777/oss
	CustomPath          string `json:",default=/oss,optional"`                    //相对路径返回的前缀
	ConnectTimeout      int64  //连接超时
	ReadWriteTimeout    int64  //读写超时
	StorePath           string `json:",optional,default=../oss"`
}

type LocalConf struct {
	OssConf
}

// minio本地存储配置
type MinioConf struct {
	OssConf
}

func (m MinioConf) GetEndPoint() string {
	return m.Location
}

// 阿里云oss配置
type AliYunConf struct {
	OssConf
}

func (a AliYunConf) GenClientOption() []oss.ClientOption {
	options := make([]oss.ClientOption, 0)
	options = append(options, oss.Timeout(a.ConnectTimeout, a.ReadWriteTimeout))
	if a.CustomHost != "" {
		options = append(options, oss.UseCname(true))
	}
	return options
}

func (a AliYunConf) GetEndPoint() string {
	scheme := "https"
	if !a.UseSSL {
		scheme = "http"
	}
	if a.CustomHost == "" {
		return fmt.Sprintf("%s://%s.aliyuncs.com", scheme, a.Location)
	}
	return fmt.Sprintf("%s://%s", scheme, a.CustomHost)
}
