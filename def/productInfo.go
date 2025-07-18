package def

type AuthMode = int64

const (
	AuthModePwd  AuthMode = 1 //账密认证
	AuthModeCert AuthMode = 2 //证书认证
)

type Net = int64

const (
	NetOther   Net = 1 //其他
	NetWifi    Net = 2 //wi-fi
	NetG234    Net = 3 //2G/3G/4G
	NetG5      Net = 4 //5G
	NetBle     Net = 5 //蓝牙
	NetLora    Net = 6 //LoRaWAN
	NetWifiBle Net = 6 //Wifi+ble
	NetWired   Net = 7 //有线网
)

type DeviceType = int64

const (
	DeviceTypeDevice  DeviceType = 1 //设备
	DeviceTypeGateway DeviceType = 2 //网关
	DeviceTypeSubset  DeviceType = 3 //子设备
	DeviceTypeMedia   DeviceType = 4 //监控设备
)

type DataProto = int64

const (
	DataProtoCustom   DataProto = 1 //自定义
	DataProtoTemplate DataProto = 2 //数据模板
)

type AutoReg = int64

const (
	AutoRegClose AutoReg = 1 //关闭
	AutoRegOpen  AutoReg = 2 //打开
	AutoRegAuto  AutoReg = 3 //打开并自动创建设备
	AutoRegBind  AutoReg = 4 //在前面的基础上绑定没有也自动创建

)
