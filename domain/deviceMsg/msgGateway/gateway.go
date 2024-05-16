package msgGateway

import (
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/utils"
)

type (
	//Msg 请求和回复结构体
	Msg struct {
		*deviceMsg.CommonMsg
		Payload *GatewayPayload `json:"payload,omitempty"`
	}
	Devices []*Device
	Device  struct {
		ProductID    string `json:"productID"`              //产品id
		DeviceName   string `json:"deviceName"`             //设备名称
		DeviceSecret string `json:"deviceSecret,omitempty"` //设备秘钥
		Register
		Code int64  `json:"code,omitempty"` //子设备绑定结果
		Msg  string `json:"msg,omitempty"`  //错误原因
	}
	Register struct {
		/*
			子设备绑定签名串。 签名算法：
			1. 签名原串，将产品 GroupIDs 设备名称，随机数，时间戳拼接：text=${product_id};${device_name};${random};${expiration_time}
			2. 使用设备 Psk 密钥，或者证书的 Sha1 摘要，进行签名：hmac_sha1(device_secret, text)
		*/
		Signature  string `json:"signature,omitempty"`
		Random     int64  `json:"random,omitempty"`     //随机数。
		Timestamp  int64  `json:"timestamp,omitempty"`  //时间戳，单位：秒。
		SignMethod string `json:"signMethod,omitempty"` //签名算法。支持 hmacsha1、hmacsha256
	}
	GatewayPayload struct {
		Status  def.GatewayStatus `json:"status,omitempty"`
		Devices Devices           `json:"devices"`
	}
)
type PackReport struct {
	Id      string `json:"id"`
	Version string `json:"version"`
	Sys     struct {
		Ack int `json:"ack"`
	} `json:"sys"`
	Params struct {
		Properties struct {
			Power struct {
				Value string `json:"value"`
				Time  int64  `json:"time"`
			} `json:"Power"`
			WF struct {
				Value struct {
				} `json:"value"`
				Time int64 `json:"time"`
			} `json:"WF"`
		} `json:"properties"`
		Events struct {
			AlarmEvent1 struct {
				Value struct {
					Param1 string `json:"param1"`
					Param2 string `json:"param2"`
				} `json:"value"`
				Time int64 `json:"time"`
			} `json:"alarmEvent1"`
			AlertEvent2 struct {
				Value struct {
					Param1 string `json:"param1"`
					Param2 string `json:"param2"`
				} `json:"value"`
				Time int64 `json:"time"`
			} `json:"alertEvent2"`
		} `json:"events"`
		SubDevices []struct {
			Identity struct {
				ProductKey string `json:"productKey"`
				DeviceName string `json:"deviceName"`
			} `json:"identity"`
			Properties struct {
				Power struct {
					Value string `json:"value"`
					Time  int64  `json:"time"`
				} `json:"Power"`
				WF struct {
					Value struct {
					} `json:"value"`
					Time int64 `json:"time"`
				} `json:"WF"`
			} `json:"properties"`
			Events struct {
				AlarmEvent1 struct {
					Value struct {
						Param1 string `json:"param1"`
						Param2 string `json:"param2"`
					} `json:"value"`
					Time int64 `json:"time"`
				} `json:"alarmEvent1"`
				AlertEvent2 struct {
					Value struct {
						Param1 string `json:"param1"`
						Param2 string `json:"param2"`
					} `json:"value"`
					Time int64 `json:"time"`
				} `json:"alertEvent2"`
			} `json:"events"`
		} `json:"subDevices"`
	} `json:"params"`
	Method string `json:"method"`
}

const (
	TypeTopo   = "topo"   //拓扑关系管理
	TypeStatus = "status" //代理子设备上下线
)

// 获取产品id列表(不重复的)
func (d Devices) GetProductIDs() []string {
	var (
		set = map[string]struct{}{}
	)
	for _, v := range d {
		set[v.ProductID] = struct{}{}
	}
	return utils.SetToSlice(set)
}
func (d Devices) GetCore() Devices {
	if d == nil {
		return nil
	}
	var ret Devices
	for _, v := range d {
		ret = append(ret, &Device{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	return ret
}
func (d Devices) GetDevCore() []*devices.Core {
	if d == nil {
		return nil
	}
	var ret []*devices.Core
	for _, v := range d {
		ret = append(ret, &devices.Core{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	return ret
}
