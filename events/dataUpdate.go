package events

type DeviceUpdateInfo struct {
	ProductID  string
	DeviceName string
	Data       any
}

type DeviceTransferInfo struct {
	ProductID     string `json:"productID"`
	DeviceName    string `json:"deviceName"`
	OldTenantCode string `json:"oldTenantCode"`
	OldProjectID  int64  `json:"oldProjectID"`
	OldAreaID     int64  `json:"oldAreaID"`
	NewTenantCode string `json:"newTenantCode"`
	NewProjectID  int64  `json:"newProjectID"`
	NewAreaID     int64  `json:"newAreaID"`
	NewAreaIDPath string `json:"newAreaIDPath"`
}

type ChangeInfo struct {
	ID int64
}

type GatewayUpdateInfo struct {
	GatewayProductID  string
	GatewayDeviceName string
	Status            int32         //拓扑关系变化状态。* 2：解绑* 1：绑定
	Devices           []*DeviceCore //子设备列表
}

type DeviceCore struct {
	ProductID  string `json:"productID"`  //产品id
	DeviceName string `json:"deviceName"` //设备名称
}
