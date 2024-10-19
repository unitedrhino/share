package schema

type (
	// Model 物模型协议-数据模板定义
	ModelSimple struct {
		Properties PropertiesSimple `json:"properties,omitempty"` //属性
		Events     EventsSimple     `json:"events,omitempty"`     //事件
		Actions    ActionsSimple    `json:"actions,omitempty"`    //行为
	}
	CommonParamSimple struct {
		Identifier   string `json:"identifier"`             //标识符 (统一)
		Name         string `json:"name"`                   //功能名称
		Desc         string `json:"desc,omitempty"`         //描述
		ExtendConfig string `json:"extendConfig,omitempty"` //拓展参数,json格式
	}

	/*事件*/
	EventSimple struct {
		CommonParamSimple
		Type   EventType `json:"type"`   //事件类型: 1:信息:info  2:告警alert  3:故障:fault
		Params Params    `json:"params"` //事件参数
	}
	EventsSimple []EventSimple

	/*行为*/
	ActionSimple struct {
		CommonParamSimple
		Dir    ActionDir `json:"dir"`    //调用方向
		Input  Params    `json:"input"`  //调用参数
		Output Params    `json:"output"` //返回参数
	}
	ActionsSimple []ActionSimple

	/*属性*/
	PropertySimple struct {
		CommonParamSimple
		Mode   PropertyMode `json:"mode"`   //读写类型:rw(可读可写) r(只读)
		Define Define       `json:"define"` //数据定义
	}
	PropertiesSimple []PropertySimple
)
