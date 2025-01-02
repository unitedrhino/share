package stores

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"reflect"
	"time"
)

// ColumnDef 定义列的结构，包含列名、GORM标签和类型描述
type ColumnDef struct {
	Tag  string
	Type string
}

/*
// 列定义示例
	columns := map[string]stores.ColumnDef{
		"ID":      {Tag: `gorm:"primary_key"`, Type: "int64"},
		"Name":    {Tag: `gorm:"size:255"`, Type: "string"},
		"Age":     {Tag: ``, Type: "int64"},
		"Active":  {Tag: `gorm:"default:true"`, Type: "bool"},
		"Created": {Tag: ``, Type: "time"},
	}
*/

func AutoMigrateDynamicTable(ctx context.Context, db *DB, tableName string, columnDefs map[string]ColumnDef) error {
	i, err := GenerateDynamicTable(columnDefs)
	if err != nil {
		return err
	}
	err = db.Table(tableName).AutoMigrate(ctx, i)
	return err
}

// GenerateDynamicTable 根据给定的列定义动态生成表结构
func GenerateDynamicTable(columnDefs map[string]ColumnDef) (interface{}, error) {
	// 定义一个结构体，用于存储列定义
	var fields []reflect.StructField

	// 遍历列定义，创建结构体字段
	for columnName, columnDef := range columnDefs {
		var fieldType reflect.Type
		switch columnDef.Type {
		case "string":
			fieldType = reflect.TypeOf("")
		case "int64":
			fieldType = reflect.TypeOf(int64(0))
		case "float64":
			fieldType = reflect.TypeOf(float64(0.0))
		case "bool":
			fieldType = reflect.TypeOf(false)
		case "time":
			fieldType = reflect.TypeOf(time.Time{})
		default:
			// 如果类型未知，使用string类型作为默认值
			return nil, errors.NotRealize.AddMsg(columnDef.Type)
		}

		// 创建结构体字段
		field := reflect.StructField{
			Name: columnName,
			Type: fieldType,
			Tag:  reflect.StructTag(columnDef.Tag),
		}
		fields = append(fields, field)
	}

	// 创建一个新的结构体类型
	newType := reflect.StructOf(fields)
	// 创建一个新的实例
	newInstance := reflect.New(newType).Interface()

	// 如果类型断言失败，返回nil
	return newInstance, nil
}

//func main() {
//	// 列定义示例
//	columns := map[string]ColumnDef{
//		"ID":     {Tag: `gorm:"primary_key"`, Type: "int64"},
//		"Name":   {Tag: `gorm:"size:255"`, Type: "string"},
//		"Age":    {Tag: ``, Type: "int64"},
//		"Active": {Tag: `gorm:"default:true"`, Type: "bool"},
//		"Created": {Tag: ``, Type: "time"},
//	}
//
//	// 表名
//	tableName := "users"
//
//	// 生成动态表结构
//	dynamicTable := GenerateDynamicTable(columns, tableName)
//
//	// 假设我们有一个GORM的DB实例
//	db, err := gorm.Open(gorm.Dialect{}, ":memory:")
//	if err != nil {
//		panic("failed to connect database")
//	}
//
//	// 使用动态表结构创建表
//	db.AutoMigrate(dynamicTable)
//
//	// 这里可以继续进行数据库操作，例如插入、查询等
//}
