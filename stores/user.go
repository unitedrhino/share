package stores

import (
	"context"
	"database/sql/driver"
	"reflect"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type DeletedBy int64
type CreatedBy int64
type UpdatedBy int64

func (t CreatedBy) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口
	expr = clause.Expr{SQL: "?", Vars: []interface{}{int64(t)}}
	return
}

func (t *CreatedBy) Scan(value interface{}) error {
	ret := utils.ToInt64(value)
	p := CreatedBy(ret)
	*t = p
	return nil
}

// Value implements the driver Valuer interface.
func (t CreatedBy) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t CreatedBy) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{UserByClause[CreatedBy]{Field: f, Opt: Create}}
}

func (t UpdatedBy) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口
	expr = clause.Expr{SQL: "?", Vars: []interface{}{int64(t)}}
	return
}

func (t *UpdatedBy) Scan(value interface{}) error {
	ret := utils.ToInt64(value)
	p := UpdatedBy(ret)
	*t = p
	return nil
}

// Value implements the driver Valuer interface.
func (t UpdatedBy) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t UpdatedBy) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{UserByClause[UpdatedBy]{Field: f, Opt: Update}}
}

func (t DeletedBy) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口
	expr = clause.Expr{SQL: "?", Vars: []interface{}{int64(t)}}
	return
}

func (t *DeletedBy) Scan(value interface{}) error {
	ret := utils.ToInt64(value)
	p := DeletedBy(ret)
	*t = p
	return nil
}

// Value implements the driver Valuer interface.
func (t DeletedBy) Value() (driver.Value, error) {
	return int64(t), nil
}

type UserByClause[keyT ~int64] struct {
	Field *schema.Field
	Opt   Opt
}

func (sd UserByClause[keyT]) Name() string {
	return ""
}

func (sd UserByClause[keyT]) Build(clause.Builder) {
}

func (sd UserByClause[keyT]) MergeClause(*clause.Clause) {
}

func (sd UserByClause[keyT]) ModifyStatement(stmt *gorm.Statement) { //查询的时候会调用此接口
	ctx := stmt.Context
	uc := ctxs.GetUserCtx(ctx)
	if uc == nil {
		return
	}
	var userID = keyT(uc.UserID)
	destV := reflect.ValueOf(stmt.Dest)
	if destV.Kind() == reflect.Array || destV.Kind() == reflect.Slice {
		if destV.Kind() == reflect.Map {
			stmt.SetColumn(sd.Field.DBName, userID, true)
			return
		}
		for i := 0; i < destV.Len(); i++ {
			dest := destV.Index(i)
			if dest.Kind() == reflect.Pointer || dest.Kind() == reflect.Interface {
				dest = dest.Elem()
			}
			field := dest.FieldByName(sd.Field.Name)
			if sd.Opt == Create && !field.IsZero() { //只有root权限的租户可以设置为其他租户
				continue
			}
			field.Set(reflect.ValueOf(userID))
		}
		return
	}
	if destV.Kind() == reflect.Map {
		stmt.SetColumn(sd.Field.DBName, userID, true)
		return
	}
	field := destV.Elem().FieldByName(sd.Field.Name)
	if sd.Opt == Create && !field.IsZero() { //只有root权限的租户可以设置为其他租户
		return
	}
	field.Set(reflect.ValueOf(userID))
}
