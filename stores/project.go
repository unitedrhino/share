package stores

import (
	"context"
	"database/sql/driver"
	"fmt"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"reflect"
)

type ProjectID int64

func (t ProjectID) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口
	stmt := db.Statement
	//authIDs, err := caches.GatherUserAuthProjectIDs(ctx)
	//if err != nil {
	//	stmt.Error = err
	//	return
	//}
	uc := ctxs.GetUserCtx(ctx)
	if t == 0 && uc != nil && uc.ProjectID != 0 {
		t = ProjectID(uc.ProjectID)
	}
	expr = clause.Expr{SQL: "?", Vars: []interface{}{int64(t)}}

	if !(uc == nil || uc.IsSuperAdmin || uc.AllProject) { //如果没有权限
		pa := uc.ProjectAuth[int64(t)]
		if pa == nil { //要有写权限
			stmt.Error = errors.Permissions.WithMsg("项目权限不足")
		}
	}

	return
}
func (t *ProjectID) Scan(value interface{}) error {
	ret := utils.ToInt64(value)
	p := ProjectID(ret)
	*t = p
	return nil
}

// Value implements the driver Valuer interface.
func (t ProjectID) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t ProjectID) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{ProjectClause{Field: f, T: t, Opt: Select}}
}
func (t ProjectID) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{ProjectClause{Field: f, T: t, Opt: Create}}
}

func (t ProjectID) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{ProjectClause{Field: f, T: t, Opt: Delete}}
}
func (t ProjectID) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{ProjectClause{Field: f, T: t, Opt: Update}}
}

type ProjectClause struct {
	clauseInterface
	Field *schema.Field
	T     ProjectID
	Opt   Opt
}

func (sd ProjectClause) GenAuthKey() string { //查询的时候会调用此接口
	return fmt.Sprintf(AuthModify, "projectID")
}

func (sd ProjectClause) ModifyStatement(stmt *gorm.Statement) { //查询的时候会调用此接口
	uc := ctxs.GetUserCtxNoNil(stmt.Context)
	if uc.ProjectID == 0 || uc.ProjectID == def.NotClassified {
		ti, err := caches.GetTenant(stmt.Context, uc.TenantCode)
		if err != nil {
			uc.ProjectID = def.NotClassified
		} else {
			uc.ProjectID = ti.DefaultProjectID
		}
	}
	switch sd.Opt {
	case Create:
		f := stmt.Schema.FieldsByName[sd.Field.Name]
		if f == nil {
			return
		}
		destV := reflect.ValueOf(stmt.Dest)
		if destV.Kind() == reflect.Array || destV.Kind() == reflect.Slice {
			for i := 0; i < destV.Len(); i++ {
				dest := destV.Index(i)
				field := GetField(dest, f.BindNames...)
				if !field.IsValid() || !field.IsZero() { //如果不是零值
					continue
				}
				var v ProjectID
				v = ProjectID(uc.ProjectID)
				field.Set(reflect.ValueOf(v))
			}
			return
		}
		field := GetField(destV, f.BindNames...)
		if !field.IsZero() { //只有root权限的租户可以设置为其他租户
			return
		}
		var v ProjectID
		v = ProjectID(uc.ProjectID)
		field.Set(reflect.ValueOf(v))

	case Update, Delete, Select:
		if uc == nil || uc.AllProject { //root 权限不用管
			return
		}
		if uc.ProjectID > def.NotClassified && !(uc.IsSuperAdmin || uc.AllProject) {
			pa := uc.ProjectAuth[uc.ProjectID]
			if pa == nil {
				stmt.Error = errors.Permissions.WithMsg("项目权限不足")
				return
			}
		}
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
			var values = []any{uc.ProjectID}
			if uc.ProjectID < def.NotClassified { //如果没有传项目ID,那么就是需要获取所有项目的参数
				values = nil
				for k := range uc.ProjectAuth {
					values = append(values, k)
				}
			}
			stmt.AddClause(clause.Where{Exprs: []clause.Expression{
				clause.IN{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Values: values},
			}})
			stmt.Clauses[sd.GenAuthKey()] = clause.Clause{}
		}
	}
}

func GetProjectAuthIDs(ctx context.Context) ([]int64, error) {
	uc := ctxs.GetUserCtxNoNil(ctx)
	if uc == nil || uc.AllProject { //root 权限不用管
		return nil, nil
	}
	if uc.ProjectID > def.NotClassified && !(uc.IsSuperAdmin || uc.AllProject) {
		pa := uc.ProjectAuth[uc.ProjectID]
		if pa == nil {
			return nil, errors.Permissions.WithMsg("项目权限不足")
		}
	}
	var values = []int64{uc.ProjectID}
	if uc.ProjectID <= def.NotClassified { //如果没有传项目ID,那么就是需要获取所有项目的参数
		values = nil
		for k := range uc.ProjectAuth {
			values = append(values, k)
		}
	}
	return values, nil
}

func GenProjectAuthScope(ctx context.Context, db *gorm.DB) *gorm.DB {
	uc := ctxs.GetUserCtxNoNil(ctx)
	if uc == nil || uc.AllProject || (uc.ProjectID <= def.NotClassified && uc.IsAdmin) { //root 权限不用管
		return db
	}
	if uc.ProjectID > def.NotClassified && !(uc.IsSuperAdmin || uc.AllProject) {
		pa := uc.ProjectAuth[uc.ProjectID]
		if pa == nil {
			db.AddError(errors.Permissions.WithMsg("项目权限不足"))
			return db
		}
	}
	var values = []any{uc.ProjectID}
	if uc.ProjectID <= def.NotClassified { //如果没有传项目ID,那么就是需要获取所有项目的参数
		values = nil
		for k := range uc.ProjectAuth {
			values = append(values, k)
		}
	}
	db = db.Where("project_id in ?", values)
	return db

}
