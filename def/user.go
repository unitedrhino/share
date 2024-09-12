package def

type AccountType = string

const (
	AccountUserName         = "userName"
	AccountTypePhone        = "phone"
	AccountTypeEmail        = "email"
	AccountTypeWechatUnion  = "wechatUnion"
	AccountTypeWechatOpen   = "wechatOpen"
	AccountTypeDingTalkUser = "dingTalkUser"
)

type RoleCode = string

const (
	RoleCodeDistributor RoleCode = "distributor" //经销商
	RoleCodeClient      RoleCode = "client"      //c端普通用户
	RoleCodeAdmin       RoleCode = "admin"       //管理员
	RoleCodeSupper      RoleCode = "supper"      //超级管理员(平台管理员)
)
