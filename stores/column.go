package stores

import (
	"fmt"
	"gitee.com/i-Things/share/conf"
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
