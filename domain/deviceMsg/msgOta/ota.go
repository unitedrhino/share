package msgOta

import (
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/errors"
)

const (
	TypeReport   = "report"
	TypeUpgrade  = "upgrade" //固件升级消息下行  返回升级信息，版本、固件地址
	TypeProgress = "progress"
)

type (
	Req struct {
		deviceMsg.CommonMsg
		Params params `json:"params,optional"`
	}
	Process struct {
		deviceMsg.CommonMsg
		Params processParams `json:"params,optional"`
	}
	params struct {
		ID      int64  `json:"id"`
		Version string `json:"version"`
	}
	processParams struct {
		ID   int64  `json:"id"`
		Step int64  `json:"step"`
		Desc string `json:"desc"`
	}

	//ota下行消息
	Upgrade struct {
		deviceMsg.CommonMsg
		Params UpgradeParams
	}
	UpgradeParams struct {
		Version          string    `json:"version"`
		IsDiff           int64     `json:"is_diff"`
		SignMethod       string    `json:"sign_method"`
		Files            []File    `json:"files"`
		DownloadProtocol string    `json:"download_protocol"`
		ExtData          []ExtData `json:"ext_data"`
	}
	ExtData struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	File struct {
		Size      int64  `json:"size"`
		Name      string `json:"name"`
		FilePath  string `json:"file_path"`
		FileMd5   string `json:"file_md5"`
		Signature string `json:"signature"`
	}
)

func (d *Req) VerifyReqParam() error {

	return nil
}
func (d *Req) GetVersion() string {
	return d.Params.Version
}
func (d *Process) VerifyReqParam() error {
	if d.Params.Step == 0 {
		return errors.Parameter.AddDetail("need add Step")
	}
	return nil
}

// 定义升级包状态常量
const (
	OtaFirmwareStatusNotRequired        = 1
	OtaFirmwareStatusNotVerified        = 2
	OtaFirmwareStatusVerified           = 3
	OtaFirmwareStatusVerifying          = 4
	OtaFirmwareStatusVerificationFailed = 5
)

// 定义升级包状态映射
var OtaFirmwareStatusMap = map[int]string{
	OtaFirmwareStatusNotRequired:        "不需要验证",
	OtaFirmwareStatusNotVerified:        "未验证",
	OtaFirmwareStatusVerified:           "已验证",
	OtaFirmwareStatusVerifying:          "验证中",
	OtaFirmwareStatusVerificationFailed: "验证失败",
}

// 根据状态值返回中文字符串
func GetOtaFirmwareStatusString(status int) string {
	if statusString, ok := OtaFirmwareStatusMap[status]; ok {
		return statusString
	}
	return "未知状态"
}

// 定义升级批次常量
const (
	ValidateUpgrade = iota + 1
	BatchUpgrade
)

var JobTypeMap = map[int]string{
	ValidateUpgrade: "验证升级包",
	BatchUpgrade:    "批量升级",
}

// 定义升级任务常量
const (
	UpgradeStatusConfirm = iota + 1
	UpgradeStatusQueued
	UpgradeStatusNotified
	UpgradeStatusInProgress
	UpgradeStatusSucceeded
	UpgradeStatusFailed
	UpgradeStatusCanceled
)

var TaskStatusMap = map[int]string{
	UpgradeStatusConfirm:    "待确认",
	UpgradeStatusQueued:     "待推送",
	UpgradeStatusNotified:   "已推送",
	UpgradeStatusInProgress: "升级中",
	UpgradeStatusSucceeded:  "升级成功",
	UpgradeStatusFailed:     "升级失败",
	UpgradeStatusCanceled:   "已取消",
}

// 定义升级批次常量

/*
静态升级：对于选定的升级范围，仅升级当前满足升级条件的设备。
动态升级：对于选定的升级范围，升级当前满足升级条件的设备，并且持续监测该范围内的设备。只要符合升级条件，物联网平台就会自动推送升级信息。包括但不限于以下设备：
满足升级条件的后续新激活设备。
当前上报的OTA模块版本号不满足升级条件，后续满足升级条件的设备。
*/
const (
	StaticUpgrade = iota + 1
	DynamicUpgrade
)

var UpgradeTypeMap = map[int]string{
	StaticUpgrade:  "静态升级",
	DynamicUpgrade: "动态升级",
}

const (
	AllUpgrade = iota + 1
	SpecificUpgrade
	GrayUpgrade
	GroupUpgrade
	AreaUpgrade
)

var UpgradeModeMap = map[int]string{
	AreaUpgrade:     "区域升级",
	AllUpgrade:      "全量升级",
	SpecificUpgrade: "定向升级",
	GrayUpgrade:     "灰度升级",
	GroupUpgrade:    "分组升级",
}

const (
	DiffPackage = iota
	FullPackage
)

var PackageTypeMap = map[int]string{
	FullPackage: "整包",
	DiffPackage: "差包",
}

const (
	JobStatusPlanned = iota + 1
	JobStatusInProgress
	JobStatusCompleted
	JobStatusCanceled
)

var JobStatusMap = map[int]string{
	JobStatusPlanned:    "计划中",
	JobStatusInProgress: "执行中",
	JobStatusCompleted:  "已完成",
	JobStatusCanceled:   "已取消",
}
