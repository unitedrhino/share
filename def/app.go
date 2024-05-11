package def

const (
	AppCore = "core"
	AppAll  = "all" //绑定所有模块
)

const (
	TenantCodeDefault = "default" //默认租户
	TenantCodeCommon  = "common"  //特殊租户:所有租户都能获取该租户的信息
)
const (
	ModuleSystemManage = "systemManage" //租户内的系统管理
	ModuleTenantManage = "tenantManage" //租户系统管理
	ModuleThings       = "things"       //物联网模块
	ModuleView         = "view"         //大屏模块
	ModuleVideo        = "video"        //音视频
)

type AppType = string

const (
	AppTypeWeb  = "web"
	AppTypeApp  = "app"
	AppTypeMini = "mini"
)

type AppSubType = string

const (
	AppSubTypeWx   = "wx"
	AppSubTypeDing = "ding"

	AppSubTypeAndroid = "android"
	AppSubTypeIos     = "ios"
)
