package utils

import (
	"database/sql"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/jinzhu/copier"
	"time"
)

var converters []copier.TypeConverter

type (
	TypeConverter = copier.TypeConverter
)

func init() {
	var (
		str1    *string
		str2    *wrappers.StringValue
		str3    sql.NullString
		str4    string
		strCopy = []copier.TypeConverter{
			{SrcType: str1, DstType: str2, Fn: func(src interface{}) (dst interface{}, err error) {
				return ToRpcNullString(src), nil
			}},
			{SrcType: str4, DstType: str2, Fn: func(src interface{}) (dst interface{}, err error) {
				return ToRpcNullString(src), nil
			}},
			{SrcType: str3, DstType: str2, Fn: func(src interface{}) (dst interface{}, err error) {
				return ToRpcNullString(src), nil
			}},
			{SrcType: str2, DstType: str1, Fn: func(src interface{}) (dst interface{}, err error) {
				return ToNullString(src.(*wrappers.StringValue)), nil
			}},
			{SrcType: str2, DstType: str4, Fn: func(src interface{}) (dst interface{}, err error) {
				return ToEmptyString(src.(*wrappers.StringValue)), nil
			}},
			{SrcType: str3, DstType: str4, Fn: func(src interface{}) (dst interface{}, err error) {
				return SqlToString(src.(sql.NullString)), nil
			}},
		}
	)
	converters = append(converters, strCopy...)
	var (
		int1 *wrappers.Int64Value
		int2 int64
		int3 *int64
		int4 sql.NullInt64
	)
	converters = append(converters,
		copier.TypeConverter{SrcType: str1, DstType: int1, Fn: func(src interface{}) (dst interface{}, err error) {
			return ToRpcNullInt64(src), nil
		}},
		copier.TypeConverter{SrcType: str4, DstType: int1, Fn: func(src interface{}) (dst interface{}, err error) {
			return ToRpcNullInt64(src), nil
		}},
		copier.TypeConverter{SrcType: str3, DstType: int1, Fn: func(src interface{}) (dst interface{}, err error) {
			return ToRpcNullInt64(src), nil
		}},
		copier.TypeConverter{SrcType: int3, DstType: int1, Fn: func(src interface{}) (dst interface{}, err error) {
			return ToRpcNullInt64(src), nil
		}},
		copier.TypeConverter{SrcType: int4, DstType: int1, Fn: func(src interface{}) (dst interface{}, err error) {
			return ToRpcNullInt64(src), nil
		}},
		copier.TypeConverter{SrcType: int2, DstType: int1, Fn: func(src interface{}) (dst interface{}, err error) {
			return ToRpcNullInt64(src), nil
		}},
		copier.TypeConverter{SrcType: int1, DstType: int3, Fn: func(src interface{}) (dst interface{}, err error) {
			return ToNullInt64(src.(*wrappers.Int64Value)), nil
		}},
		copier.TypeConverter{SrcType: int1, DstType: int2, Fn: func(src interface{}) (dst interface{}, err error) {
			return ToEmptyInt64(src.(*wrappers.Int64Value)), nil
		}})
	var (
		time1 time.Time
		time2 *time.Time
		time3 sql.NullTime
	)
	converters = append(converters,
		copier.TypeConverter{SrcType: time1, DstType: int2, Fn: func(src interface{}) (dst interface{}, err error) {
			t := src.(time.Time)
			if t.IsZero() {
				return int64(0), nil
			}
			return t.Unix(), nil
		}},
		copier.TypeConverter{SrcType: time2, DstType: int2, Fn: func(src interface{}) (dst interface{}, err error) {
			t := src.(*time.Time)
			if t == nil {
				return int64(0), nil
			}
			return t.Unix(), nil
		}},
		copier.TypeConverter{SrcType: time3, DstType: int2, Fn: func(src interface{}) (dst interface{}, err error) {
			t := src.(sql.NullTime)
			if t.Valid == false {
				return int64(0), nil
			}
			return t.Time.Unix(), nil
		}},
		copier.TypeConverter{SrcType: int2, DstType: time1, Fn: func(src interface{}) (dst interface{}, err error) {
			in := src.(int64)
			if in == 0 {
				return time.Time{}, nil
			}
			return time.Unix(in, 0), nil
		}},
		copier.TypeConverter{SrcType: int2, DstType: time2, Fn: func(src interface{}) (dst interface{}, err error) {
			return Int64ToTimex(src.(int64)), nil
		}},
		copier.TypeConverter{SrcType: int1, DstType: time3, Fn: func(src interface{}) (dst interface{}, err error) {
			return ToNullTime2(src.(*wrappers.Int64Value)), nil
		}},
		copier.TypeConverter{SrcType: time3, DstType: int1, Fn: func(src interface{}) (dst interface{}, err error) {
			return TimeToNullInt(src.(sql.NullTime)), nil
		}},
		copier.TypeConverter{SrcType: int2, DstType: time3, Fn: func(src interface{}) (dst interface{}, err error) {
			return Int64ToSqlTime(src.(int64)), nil
		}})

	var (
		float1    *float32
		float2    *wrappers.FloatValue
		float3    sql.NullFloat64
		float4    float32
		floatCopy = []copier.TypeConverter{
			{SrcType: float1, DstType: float2, Fn: func(src interface{}) (dst interface{}, err error) {
				return ToRpcNullFloat32(src), nil
			}},
			{SrcType: float4, DstType: float2, Fn: func(src interface{}) (dst interface{}, err error) {
				return ToRpcNullFloat32(src), nil
			}},
			{SrcType: float3, DstType: float2, Fn: func(src interface{}) (dst interface{}, err error) {
				return ToRpcNullFloat32(src), nil
			}},
			{SrcType: float2, DstType: float1, Fn: func(src interface{}) (dst interface{}, err error) {
				return ToNullFloat32(src.(*wrappers.FloatValue)), nil
			}},
			{SrcType: float2, DstType: float4, Fn: func(src interface{}) (dst interface{}, err error) {
				return ToEmptyFloat32(src.(*wrappers.FloatValue)), nil
			}},
			{SrcType: float3, DstType: float4, Fn: func(src interface{}) (dst interface{}, err error) {
				return SqlToFloat32(src.(sql.NullFloat64)), nil
			}},
		}
	)
	converters = append(converters, floatCopy...)

	var (
		mapString map[string]string
		mapCopy   = []copier.TypeConverter{
			{SrcType: mapString, DstType: mapString, Fn: func(src interface{}) (dst interface{}, err error) {
				return src, nil
			}},
		}
	)
	converters = append(converters, mapCopy...)

}

func AddConverter(in ...copier.TypeConverter) {
	converters = append(converters, in...)
}

func CopyE(toValue interface{}, fromValue interface{}) (err error) {
	return copier.CopyWithOption(toValue, fromValue, copier.Option{
		DeepCopy:   true,
		Converters: converters,
	})
}

func Copy[toT any](fromValue any) *toT {
	var toValue toT
	if fromValue == nil {
		return nil
	}
	err := CopyE(&toValue, fromValue)
	if err != nil {
		return nil
	}
	return &toValue
}

func Copy2[toT any](fromValue any) toT {
	var toValue toT
	if fromValue == nil {
		return toValue
	}
	err := CopyE(&toValue, fromValue)
	if err != nil {
		return toValue
	}
	return toValue
}

func CopySlice[toT any, fromT any](fromValue []*fromT) []*toT {
	if fromValue == nil {
		return nil
	}
	var ret []*toT
	for _, v := range fromValue {
		ret = append(ret, Copy[toT](v))
	}
	return ret
}

func CopyMap[toT any, fromT any, keyT comparable](fromValue map[keyT]*fromT) map[keyT]*toT {
	if len(fromValue) == 0 {
		return nil
	}
	var ret = map[keyT]*toT{}
	for k, v := range fromValue {
		ret[k] = Copy[toT](v)
	}
	return ret
}
