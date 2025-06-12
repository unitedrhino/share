package stores

import (
	"fmt"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/utils"
	"gorm.io/gorm"
	"strings"
)

var Expr = gorm.Expr

type CmpType = string

const (
	CmpTypeEq    CmpType = "="      //相等
	CmpTypeNot   CmpType = "!="     //不相等
	CmpTypeGt    CmpType = ">"      //大于
	CmpTypeGte   CmpType = ">="     //大于等于
	CmpTypeLt    CmpType = "<"      //小于
	CmpTypeLte   CmpType = "<="     //小于等于
	CmpTypeIn    CmpType = "in"     //在xx值之中,可以有n个参数
	CmpTypeNotIn CmpType = "not in" //在xx值之中,可以有n个参数
	CmpTypeLike  CmpType = "like"   //模糊查询
)

type toSqlFunc func(db *DB, column string) string

type Cmp struct {
	Value     any
	toSqlFunc toSqlFunc
}

func defaultToSql(c CmpType) toSqlFunc {
	return func(db *DB, column string) string {
		return fmt.Sprintf("%s %s ?", column, string(c))
	}
}

func GetCmp(cmpType CmpType, value any) *Cmp {
	switch cmpType {
	case CmpTypeEq:
		return CmpEq(value)
	case CmpTypeNot:
		return CmpNot(value)
	case CmpTypeGt:
		return CmpGt(value)
	case CmpTypeGte:
		return CmpGte(value)
	case CmpTypeLt:
		return CmpLt(value)
	case CmpTypeLte:
		return CmpLte(value)
	}
	return nil
}

func CmpEq(value any) *Cmp {
	return &Cmp{toSqlFunc: defaultToSql(CmpTypeEq), Value: value}
}
func CmpNot(value any) *Cmp {
	return &Cmp{toSqlFunc: defaultToSql(CmpTypeNot), Value: value}
}
func CmpGt(value any) *Cmp {
	return &Cmp{toSqlFunc: defaultToSql(CmpTypeGt), Value: value}
}
func CmpGte(value any) *Cmp {
	return &Cmp{toSqlFunc: defaultToSql(CmpTypeGte), Value: value}
}
func CmpLt(value any) *Cmp {
	return &Cmp{toSqlFunc: defaultToSql(CmpTypeLt), Value: value}
}

// 是否是null
func CmpIsNull(isNull bool) *Cmp {
	return &Cmp{toSqlFunc: func(db *DB, column string) string {
		var isNullStr = "null"
		if isNull == false {
			isNullStr = "not null"
		}
		return fmt.Sprintf("%s is %s", column, isNullStr)
	}}
}

func CmpEqColumn(c CmpType, columnLeft string) *Cmp {
	return &Cmp{toSqlFunc: func(db *DB, column string) string {
		return fmt.Sprintf("%s %s %s", columnLeft, c, column)
	}}
}

func CmpLike(value string) *Cmp {
	return &Cmp{toSqlFunc: func(db *DB, column string) string {
		return fmt.Sprintf("%s like %s", column, "%"+value+"%")
	}}
}

func CmpLte(value any) *Cmp {
	return &Cmp{toSqlFunc: defaultToSql(CmpTypeLte), Value: value}
}

func CmpIn[t any](values ...t) *Cmp {
	if len(values) == 0 {
		return nil
	}
	return &Cmp{toSqlFunc: defaultToSql(CmpTypeIn), Value: values}
}

func CmpNotIn[t any](values ...t) *Cmp {
	if len(values) == 0 {
		return nil
	}
	return &Cmp{toSqlFunc: defaultToSql(CmpTypeNotIn), Value: values}
}

// json对象中 obj.key = v
func CmpJsonObjEq(k string, v any) *Cmp {
	switch GetDBType() {
	case conf.Pgsql:
		return &Cmp{toSqlFunc: func(db *DB, column string) string {
			return fmt.Sprintf(` %s::jsonb @> '{"%v": "%v"}'::jsonb  `, column, k, v)
		}}
	default:
		return &Cmp{toSqlFunc: func(db *DB, column string) string {
			return fmt.Sprintf("JSON_CONTAINS(%s, JSON_OBJECT(?,?))", column)
		}, Value: []any{k, v}}
	}

}

// json对象中 obj.key like '%v%'
func CmpJsonObjLike(k string, v any) *Cmp {
	return &Cmp{toSqlFunc: func(db *DB, column string) string {
		switch GetDBType() {
		case conf.Pgsql:
			return fmt.Sprintf("jsonb_path_exists(%s, '$.\"%s\" ? (@ like_regex $pattern flags \"i\")', jsonb_build_object('pattern', ?))",
				column, k)
		default:
			return fmt.Sprintf("JSON_SEARCH(%s, 'all', ?, NULL, '$.\"%s\"') IS NOT NULL", column, k)
		}
	}, Value: []any{"%" + utils.ToString(v) + "%"}}

}

// json 数组中包含 v的值 或方式
func CmpJsonArrayOrHas[t any](v ...t) *Cmp {
	switch GetDBType() {
	case conf.Pgsql:
		return &Cmp{toSqlFunc: func(db *DB, column string) string {

			// PostgreSQL 实现：使用 jsonb_path_exists 函数和路径表达式
			var conditions []string
			for _, val := range v {
				// 将值转换为 JSON 字符串并处理特殊字符
				// 使用路径表达式检查数组中是否包含该值
				conditions = append(conditions,
					fmt.Sprintf("jsonb_path_exists(%s, '$[*] ? (@ == \"%s\")')",
						column, val))
			}
			return strings.Join(conditions, " OR ")
		}}

	default:
		return &Cmp{toSqlFunc: func(db *DB, column string) string {

			// MySQL 实现：使用 JSON_SEARCH 函数
			var conditions []string
			for range v {
				conditions = append(conditions,
					fmt.Sprintf("JSON_SEARCH(%s, 'all', ?) IS NOT NULL", column))
			}
			return strings.Join(conditions, " OR ")
		}, Value: utils.ToAnySlice(v)}
	}
}

// 过滤二进制比特位
func CmpBinEq(bit int64, isHigh int64) *Cmp {
	return &Cmp{toSqlFunc: func(db *DB, column string) string {
		return fmt.Sprintf("(%s >> %d) & 1 = ?", column, bit)
	}, Value: isHigh}
}

func CmpAnd(cs ...*Cmp) *Cmp {
	var values []any
	for _, v := range cs {
		values = append(values, v.ToValues()...)
	}
	return &Cmp{Value: values, toSqlFunc: func(db *DB, column string) string {
		var sqls []string
		for _, v := range cs {
			sqls = append(sqls, v.ToSql(db, column))
		}
		return strings.Join(sqls, " and ")
	}}
}

func CmpOr(cs ...*Cmp) *Cmp {
	var values []any
	for _, v := range cs {
		values = append(values, v.ToValues()...)
	}
	return &Cmp{Value: values, toSqlFunc: func(db *DB, column string) string {
		var sqls []string
		for _, v := range cs {
			sqls = append(sqls, v.ToSql(db, column))
		}
		return strings.Join(sqls, " or ")
	}}
}

// 大于=? 小于等于? 在xx之间
func CmpBtw(max any, min any) *Cmp {
	return &Cmp{
		Value: []any{max, min},
		toSqlFunc: func(db *DB, column string) string {
			return fmt.Sprintf("%s <= ? && %s >= ?", column, column)
		},
	}
}

func (g *Cmp) ToValues() []any {
	if g == nil {
		return nil
	}
	switch g.Value.(type) {
	case []any:
		return g.Value.([]any)
	case nil:
		return nil
	default:
		return []any{g.Value}
	}
}

func (g *Cmp) ToSql(db *DB, column string) string {
	if g == nil {
		return ""
	}
	return g.toSqlFunc(db, column)
}

func (g *Cmp) Where(db *DB, column string) *gorm.DB {
	if g != nil {
		db = db.Where(g.ToSql(db, column), g.ToValues()...)
	}
	return db
}
