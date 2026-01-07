package stores

import (
	"fmt"
	"reflect"
	"sync"

	"gitee.com/unitedrhino/share/conf"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

var columnCache sync.Map

func SetColumns(db *DB, value any, index string) (ret []clause.Column) {
	s, err := schema.Parse(value, &columnCache, db.NamingStrategy)
	if err != nil {
		panic(err) //生产阶段不应该报错
	}
	for _, field := range s.Fields {
		t := field.TagSettings["UNIQUEINDEX"]
		if t == "" || t != index {
			continue
		}
		ret = append(ret, clause.Column{
			Name: field.DBName,
		})
	}
	return
}

func SetColumnsWithPg(db *DB, value any, index string) (ret []clause.Column) {
	if rlDBType != conf.Pgsql {
		return []clause.Column{}
	}
	return SetColumns(db, value, index)
}

func SetWithPg[v any](value v, def v) v {
	if rlDBType != conf.Pgsql {
		return def
	}
	return value
}

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

// 待完善,pg批量导入之后序列未更新
func CreateInBatches(db *DB, value any, batchSize int) error {
	db = db.Session(&gorm.Session{})//如果不加,pg会报错,模型不更新
	db = db.CreateInBatches(value, batchSize)
	if db.Error != nil {
		return db.Error
	}
	if rlDBType == conf.Pgsql {
		s, err := schema.Parse(value, &columnCache, db.NamingStrategy)
		if err != nil {
			return err
		}
		if s.PrioritizedPrimaryField != nil {
			db.Model(value).Exec(fmt.Sprintf(`SELECT setval( '%s_%s_seq',(SELECT COALESCE(MAX(%s), 0) + 1 FROM %s));`,
				s.Table, s.PrioritizedPrimaryField.DBName, s.PrioritizedPrimaryField.DBName, s.Table))
		}

	}
	return nil
}
