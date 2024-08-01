package stores

import (
	"gorm.io/gorm"
	"reflect"
)

func GetField(val reflect.Value, bindNames ...string) (ret reflect.Value) {
	var field interface{}
	for _, bindName := range bindNames {
		ele := val
		if !(val.Kind() == reflect.Struct) {
			ele = val.Elem()
			if !ele.IsValid() {
				return ele
			}
		}

		ret = ele.FieldByName(bindName)
		if !ret.IsValid() {
			return
		}
		field = ret.Interface()
		val = reflect.ValueOf(field)
	}
	return
}

func SetValue[valT any](stmt *gorm.Statement, fieldName string, v valT) {
	f := stmt.Schema.FieldsByName[fieldName]
	if f == nil {
		return
	}
	destV := reflect.ValueOf(stmt.Dest)
	if destV.Kind() == reflect.Array || destV.Kind() == reflect.Slice {
		for i := 0; i < destV.Len(); i++ {
			dest := destV.Index(i)
			field := GetField(dest, f.BindNames...)
			if !field.IsZero() { //只有root权限的租户可以设置为其他租户
				continue
			}
			field.Set(reflect.ValueOf(v))
		}
		return
	}
	field := GetField(destV, f.BindNames...)
	if !field.IsZero() { //只有root权限的租户可以设置为其他租户
		return
	}
	field.Set(reflect.ValueOf(v))
}
