package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
)

func Unmarshal(data []byte, v any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	return decoder.Decode(v)
}

func UnmarshalNoErr[inT any](data string) inT {
	var ret inT
	json.Unmarshal([]byte(data), &ret)
	return ret
}

func UnmarshalSlices[inT any](datas []string) (ret []*inT, err error) {
	for _, data := range datas {
		var one inT
		err := json.Unmarshal([]byte(data), &one)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &one)
	}
	return
}

func MarshalSlices[inT any](datas []*inT) (ret []string, err error) {
	for _, data := range datas {
		v, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		ret = append(ret, string(v))
	}
	return
}

func MarshalNoErr(v any) string {
	ret, _ := json.Marshal(v)
	return string(ret)
}

// Fmt 将结构以更方便看的方式打印出来
func Fmt(v any) string {
	switch v.(type) {
	case string:
		return v.(string)
	case []byte:
		return string(v.([]byte))
	case error:
		return v.(error).Error()
	case interface{ String() string }:
		return v.(interface{ String() string }).String()
	default:
		val, err := cast.ToStringE(v)
		if err == nil {
			return val
		}
		js, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%#v", js)
		}
		return string(js)
	}
}

func Fmt2(v any) string {
	switch v.(type) {
	case string:
		return v.(string)
	case []byte:
		return string(v.([]byte))
	case error:
		return v.(error).Error()
	default:
		val, err := cast.ToStringE(v)
		if err == nil {
			return val
		}
		js, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%#v", js)
		}
		return string(js)
	}
}
