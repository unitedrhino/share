package stores

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type CmpType string

const (
	CmpTypeEq  CmpType = "="  //相等
	CmpTypeNot CmpType = "!=" //不相等
	CmpTypeGt  CmpType = ">"  //大于
	CmpTypeGte CmpType = ">=" //大于等于
	CmpTypeLt  CmpType = "<"  //小于
	CmpTypeLte CmpType = "<=" //小于等于
	CmpTypeIn  CmpType = "in" //在xx值之中,可以有n个参数
)

type toSqlFunc func(column string) string

type Cmp struct {
	Value     any
	toSqlFunc toSqlFunc
}

func defaultToSql(c CmpType) toSqlFunc {
	return func(column string) string {
		return fmt.Sprintf("%s %s ?", column, string(c))
	}
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
	return &Cmp{toSqlFunc: func(column string) string {
		var isNullStr = "null"
		if isNull == false {
			isNullStr = "not null"
		}
		return fmt.Sprintf("%s is %s", column, isNullStr)
	}}
}
func CmpLte(value any) *Cmp {
	return &Cmp{toSqlFunc: defaultToSql(CmpTypeLte), Value: value}
}
func CmpIn(values ...any) *Cmp {
	if len(values) == 0 {
		return nil
	}
	return &Cmp{toSqlFunc: defaultToSql(CmpTypeIn), Value: values}
}

// 过滤二进制比特位
func CmpBinEq(bit int64, isHigh int64) *Cmp {
	return &Cmp{toSqlFunc: func(column string) string {
		return fmt.Sprintf("(%s >> %d) & 1 = ?", column, bit)
	}, Value: isHigh}
}
func CmpAnd(cs ...*Cmp) *Cmp {
	var values []any
	for _, v := range cs {
		values = append(values, v.ToValues()...)
	}
	return &Cmp{Value: values, toSqlFunc: func(column string) string {
		var sqls []string
		for _, v := range cs {
			sqls = append(sqls, v.ToSql(column))
		}
		return strings.Join(sqls, " and ")
	}}
}

func CmpOr(cs ...*Cmp) *Cmp {
	var values []any
	for _, v := range cs {
		values = append(values, v.ToValues()...)
	}
	return &Cmp{Value: values, toSqlFunc: func(column string) string {
		var sqls []string
		for _, v := range cs {
			sqls = append(sqls, v.ToSql(column))
		}
		return strings.Join(sqls, " or ")
	}}
}

// 大于=? 小于等于? 在xx之间
func CmpBtw(max any, min any) *Cmp {
	return &Cmp{
		Value: []any{max, min},
		toSqlFunc: func(column string) string {
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

func (g *Cmp) ToSql(column string) string {
	if g == nil {
		return ""
	}
	return g.toSqlFunc(column)
}

func (g *Cmp) Where(db *gorm.DB, column string) *gorm.DB {
	if g != nil {
		db = db.Where(g.ToSql(column), g.ToValues()...)
	}
	return db
}
