package stores

import (
	"fmt"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/utils"
	"gorm.io/gorm"
)

// 特殊字符需要用该函数来包裹
func Col(column string) string {
	switch rlDBType {
	case conf.Pgsql:
		return fmt.Sprintf(`"%s"`, column)
	default:
		return fmt.Sprintf("`%s`", column)
	}
}

func JsonCol(column string, key string) string {
	switch rlDBType {
	case conf.Pgsql:
		// PostgreSQL jsonb提取语法（->>用于获取文本值）
		// 支持嵌套键：config.timeout -> 'config'->>'timeout'
		return fmt.Sprintf(`"%s"->>'%s'`, column, key)
	default:
		return fmt.Sprintf(" JSON_UNQUOTE(JSON_EXTRACT(`%s`, '$.%s'))", column, key)
	}
}

type CastTo = string

const (
	CastToString CastTo = "string"
	CastToInt    CastTo = "int"
	CastToFloat  CastTo = "float"
)

func JsonCol2(column string, key string, castTo CastTo) string {
	switch rlDBType {
	case conf.Pgsql:
		return fmt.Sprintf(`"%s"`, column)
	default:
		return fmt.Sprintf(" JSON_UNQUOTE(JSON_EXTRACT(`%s`, '$.%s'))", column, key)
	}
}

// JsonCol2 复用JsonCol获取基础表达式，再添加类型转换
func Cast(column string, castTo CastTo) string {
	// 根据转换类型和数据库类型处理
	switch castTo {
	case CastToInt:
		switch rlDBType {
		case conf.Pgsql:
			return fmt.Sprintf("(%s)::int", column)
		default: // MySQL/MariaDB
			return fmt.Sprintf("CAST(%s AS UNSIGNED)", column)
		}
	case CastToFloat:
		switch rlDBType {
		case conf.Pgsql:
			return fmt.Sprintf("(%s)::float8", column)
		default: // MySQL/MariaDB
			return fmt.Sprintf("CAST(%s AS DECIMAL(10,2))", column)
		}
	case CastToString:
		fallthrough // 字符串无需转换，直接返回基础表达式
	default:
		return column
	}
}

func ColWithT(column string, tableAlias string) string {
	if tableAlias != "" {
		tableAlias = tableAlias + "."
	}
	switch rlDBType {
	case conf.Pgsql:
		return fmt.Sprintf(`%s"%s"`, tableAlias, column)
	default:
		return fmt.Sprintf("%s`%s`", tableAlias, column)
	}
}

func WithSelect(db *gorm.DB, columns ...string) *gorm.DB {
	if len(columns) == 0 {
		return db
	}
	var newColumns []string
	for _, v := range columns {
		newColumns = append(newColumns, utils.CamelCaseToUdnderscore(v))
	}
	return db.Select(newColumns)
}

type Compare struct {
	CmpType string //"=":相等 "!=":不相等 ">":大于">=":大于等于"<":小于"<=":小于等于 "like":模糊查询
	Value   string //值
	CastTo  string //参数的类型(填写了会进行数据类型的转换): int float string
}
