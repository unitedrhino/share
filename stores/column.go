package stores

import (
	"fmt"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/utils"
	"gorm.io/gorm"
)

// 特殊字符需要用该函数来包裹
func Col(column string) string {
	switch dbType {
	case conf.Pgsql:
		return fmt.Sprintf(`"%s"`, column)
	default:
		return fmt.Sprintf("`%s`", column)
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
