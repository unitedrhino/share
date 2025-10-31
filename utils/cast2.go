package utils

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"github.com/spf13/cast"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func Int8ToBool(i int8) bool {
	if i == 0 {
		return false
	}
	return true
}

func ToInt64E(i any) (int64, error) {
	switch i.(type) {
	case string:
		i = strings.TrimLeft(i.(string), "0")
		return cast.ToInt64E(i)
	case []uint8:
		t := i.([]uint8)
		return cast.ToInt64E(string(t))
	case sql.NullTime:
		t := i.(sql.NullTime)
		if t.Valid == false {
			return 0, nil
		}
		return i.(sql.NullTime).Time.Unix(), nil
	case *time.Time:
		t := i.(*time.Time)
		if t == nil {
			return 0, nil
		}
		return t.Unix(), nil
	case time.Time:
		return i.(time.Time).Unix(), nil
	case *wrapperspb.Int64Value:
		return i.(*wrapperspb.Int64Value).GetValue(), nil
	default:
		return cast.ToInt64E(i)
	}
}

func ToInt64(i any) int64 {
	v, _ := ToInt64E(i)
	return v
}

func ToBool(i any) bool {
	switch i.(type) {
	case int8:
		if i.(int8) == 0 {
			return false
		}
		return true
	default:
		return cast.ToBool(i)
	}
}

// ToTime casts an interface to a time.Time type.
func ToTime(i any) time.Time {
	return cast.ToTime(i)
}

// ToDuration casts an interface to a time.Duration type.
func ToDuration(i any) time.Duration {
	return cast.ToDuration(i)
}

// ToFloat64 casts an interface to a float64 type.
func ToFloat64(i any) float64 {
	return cast.ToFloat64(i)
}

// ToFloat32 casts an interface to a float32 type.
func ToFloat32(i any) float32 {
	switch i.(type) {
	case *wrapperspb.FloatValue:
		v := i.(*wrapperspb.FloatValue)
		if v == nil {
			return 0
		}
		return v.GetValue()
	}
	return cast.ToFloat32(i)
}

// ToInt32 casts an interface to an int32 type.
func ToInt32(i any) int32 {
	return cast.ToInt32(i)
}

// ToInt16 casts an interface to an int16 type.
func ToInt16(i any) int16 {
	return cast.ToInt16(i)
}

// ToInt8 casts an interface to an int8 type.
func ToInt8(i any) int8 {
	return cast.ToInt8(i)
}

// ToInt casts an interface to an int type.
func ToInt(i any) int {
	switch i.(type) {
	case string:
		i = strings.TrimLeft(i.(string), "0")
	}
	return cast.ToInt(i)
}

func ToIntE(i any) (int, error) {
	switch i.(type) {
	case string:
		i = strings.TrimLeft(i.(string), "0")
	}
	return cast.ToIntE(i)
}

// ToUint casts an interface to a uint type.
func ToUint(i any) uint {
	return cast.ToUint(i)
}

// ToUint64 casts an interface to a uint64 type.
func ToUint64(i any) uint64 {
	return cast.ToUint64(i)
}

// ToUint32 casts an interface to a uint32 type.
func ToUint32(i any) uint32 {
	return cast.ToUint32(i)
}

// ToUint16 casts an interface to a uint16 type.
func ToUint16(i any) uint16 {
	return cast.ToUint16(i)
}

// ToUint8 casts an interface to a uint8 type.
func ToUint8(i any) uint8 {
	return cast.ToUint8(i)
}

// ToString casts an interface to a string type.
func ToString(i any) string {
	ret, err := cast.ToStringE(i)
	if err != nil {
		ret, _ := json.Marshal(i)
		return string(ret)
	}
	return ret
}

func BoolToInt(in any) any {
	if v, ok := in.(bool); ok {
		if v {
			return 1
		}
		return 0
	}
	return in
}

// ToStringMapStringSlice casts an interface to a map[string][]string type.
func ToStringMapStringSlice(i any) map[string][]string {
	return cast.ToStringMapStringSlice(i)
}

// ToStringMapBool casts an interface to a map[string]bool type.
func ToStringMapBool(i any) map[string]bool {
	return cast.ToStringMapBool(i)
}

// ToStringMapInt casts an interface to a map[string]int type.
func ToStringMapInt(i any) map[string]int {
	return cast.ToStringMapInt(i)
}

// ToStringMapInt64 casts an interface to a map[string]int64 type.
func ToStringMapInt64(i any) map[string]int64 {
	return cast.ToStringMapInt64(i)
}

// ToSlice casts an interface to a []interface{} type.
func ToSlice(i any) []any {
	return cast.ToSlice(i)
}

// ToBoolSlice casts an interface to a []bool type.
func ToBoolSlice(i any) []bool {
	return cast.ToBoolSlice(i)
}

// ToStringSlice casts an interface to a []string type.
func ToStringSlice(i any) []string {
	return cast.ToStringSlice(i)
}

// ToIntSlice casts an interface to a []int type.
func ToIntSlice(i any) []int {
	return cast.ToIntSlice(i)
}

// ToDurationSlice casts an interface to a []time.Duration type.
func ToDurationSlice(i any) []time.Duration {
	return cast.ToDurationSlice(i)
}
