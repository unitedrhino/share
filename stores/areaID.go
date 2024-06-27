package stores

import (
	"context"
	"database/sql/driver"
	"fmt"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type AreaID int64

func (t AreaID) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口
	stmt := db.Statement
	uc := ctxs.GetUserCtxOrNil(ctx)
	expr = clause.Expr{SQL: "?", Vars: []interface{}{int64(t)}}

	authType, areas := ctxs.GetAreaIDs(uc.ProjectID, uc.ProjectAuth)
	if t != def.NotClassified && !(uc.IsAdmin || uc.AllArea || authType <= def.AuthReadWrite || utils.SliceIn(int64(t), areas...)) { //如果没有权限
		stmt.Error = errors.Permissions.WithMsg("区域权限不足")
	}
	return
}
func (t *AreaID) Scan(value interface{}) error {
	ret := utils.ToInt64(value)
	p := AreaID(ret)
	*t = p
	return nil
}

// Value implements the driver Valuer interface.
func (t AreaID) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t AreaID) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{AreaClause{Field: f, T: t, Opt: Select}}
}
func (t AreaID) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{AreaClause{Field: f, T: t, Opt: Update}}
}

func (t AreaID) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{AreaClause{Field: f, T: t, Opt: Create}}
}

func (t AreaID) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{AreaClause{Field: f, T: t, Opt: Delete}}
}

type AreaClause struct {
	clauseInterface
	Field *schema.Field
	T     AreaID
	Opt   Opt
}

func (sd AreaClause) GenAuthKey() string { //查询的时候会调用此接口
	return fmt.Sprintf(AuthModify, "areaID")
}

func (sd AreaClause) ModifyStatement(stmt *gorm.Statement) { //查询的时候会调用此接口
	uc := ctxs.GetUserCtxOrNil(stmt.Context)
	if uc == nil {
		return
	}
	authType, areas := ctxs.GetAreaIDs(uc.ProjectID, uc.ProjectAuth)
	if uc.IsAdmin || uc.AllArea || authType <= def.AuthReadWrite {
		return
	}
	switch sd.Opt {
	case Create:
	case Update, Delete, Select:
		if _, ok := stmt.Clauses[sd.GenAuthKey()]; !ok {
			if c, ok := stmt.Clauses["WHERE"]; ok {
				if where, ok := c.Expression.(clause.Where); ok && len(where.Exprs) > 1 {
					for _, expr := range where.Exprs {
						if orCond, ok := expr.(clause.OrConditions); ok && len(orCond.Exprs) == 1 {
							where.Exprs = []clause.Expression{clause.And(where.Exprs...)}
							c.Expression = where
							stmt.Clauses["WHERE"] = c
							break
						}
					}
				}
			}
			if len(areas) == 0 { //如果没有权限
				stmt.Error = errors.Permissions.WithMsg("区域权限不足")
				return
			}
			var values = []any{def.NotClassified}
			for _, v := range areas {
				values = append(values, v)
			}
			stmt.AddClause(clause.Where{Exprs: []clause.Expression{
				clause.IN{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Values: values},
			}})
			stmt.Clauses[sd.GenAuthKey()] = clause.Clause{}
		}
	}
}
func GenAreaAuthScope(ctx context.Context, db *gorm.DB) *gorm.DB {
	uc := ctxs.GetUserCtxOrNil(ctx)
	if uc == nil {
		return db
	}
	authType, areas := ctxs.GetAreaIDs(uc.ProjectID, uc.ProjectAuth)
	if uc.IsAdmin || uc.AllArea || authType <= def.AuthReadWrite {
		return db
	}
	if len(areas) == 0 { //如果没有权限
		db.AddError(errors.Permissions.WithMsg("区域权限不足"))
		return db
	}
	var values = []any{def.NotClassified}
	for _, v := range areas {
		values = append(values, v)
	}
	db = db.Where("area_id in ?", values)
	return db
}
