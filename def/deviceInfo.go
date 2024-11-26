package def

type LogLevel = int64

const (
	LogClose LogLevel = 1 //关闭
	LogError LogLevel = 2 //错误
	LogWarn  LogLevel = 3 //告警
	LogInfo  LogLevel = 4 //信息
	LogDebug LogLevel = 5 //调试
)

var LogLevelTextToIntMap = map[string]LogLevel{
	"关闭": LogClose,
	"错误": LogError,
	"告警": LogWarn,
	"信息": LogInfo,
	"调试": LogDebug,
}

type GatewayStatus = int64

const (
	GatewayBind   GatewayStatus = 1 //绑定
	GatewayUnbind GatewayStatus = 2 //解绑
)

type DeviceStatus = int64

const (
	DeviceStatusInactive  DeviceStatus = 1 // 未激活
	DeviceStatusOnline    DeviceStatus = 2 //在线
	DeviceStatusOffline   DeviceStatus = 3 //离线
	DeviceStatusAbnormal  DeviceStatus = 4 //异常
	DeviceStatusArrearage DeviceStatus = 5 //欠费
	DeviceStatusWarming   DeviceStatus = 6 //告警中
)

type ConnStatus = int64

const (
	ConnectedStatus    ConnStatus = 1
	DisConnectedStatus ConnStatus = 2
)

type MobileOperator = int64 //移动运营商

const (
	MobileOperatorYD   = 1  //移动
	MobileOperatorLT   = 2  //联通
	MobileOperatorDX   = 3  //电信
	MobileOperatorGD   = 4  //移动
	MobileOperatorNone = 10 //无
)
