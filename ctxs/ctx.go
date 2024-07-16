package ctxs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/spf13/cast"
	"google.golang.org/grpc/metadata"
	"net/http"
)

type UserCtx struct {
	IsOpen         bool //是否开放认证用户
	AppCode        string
	Token          string
	TenantCode     string //租户Code
	AcceptLanguage string
	ProjectID      int64 `json:",string"`
	IsAdmin        bool  //是否是超级管理员
	IsSuperAdmin   bool
	UserID         int64   `json:",string"` //用户id（开放认证用户值为0）
	RoleIDs        []int64 //用户使用的角色（开放认证用户值为0）
	RoleCodes      []string
	IsAllData      bool   //是否所有数据权限（开放认证用户值为true）
	IP             string //用户的ip地址
	Os             string //操作系统
	UserName       string
	Account        string
	ProjectAuth    map[int64]*ProjectAuth
	InnerCtx
}

type ProjectAuth struct {
	Area map[int64]int64 //key是区域ID,value是授权类型
	// 1 //管理权限,可以修改别人的权限,及读写权限 管理权限不限制区域权限
	// 2 //读权限,只能读,不能修改
	// 3 //读写权限,可以读写该权限
	AuthType def.AuthType //项目的授权类型
}

func GetAllAreaIDs(in map[int64]*ProjectAuth) (areas []int64) {
	for _, v := range in {
		for area := range v.Area {
			areas = append(areas, area)
		}
	}
	return
}

func GetAreaIDs(projectID int64, in map[int64]*ProjectAuth) (authType def.AuthType, areas []int64) {
	v := in[projectID]
	if v == nil {
		return
	}
	authType = v.AuthType
	for area := range v.Area {
		areas = append(areas, area)
	}
	return
}

func (u *UserCtx) ClearInner() *UserCtx {
	if u == nil {
		return nil
	}
	newCtx := *u
	newCtx.InnerCtx = InnerCtx{}
	return &newCtx
}

func (u *UserCtx) HasRole(roleCodes ...string) bool {
	if u == nil {
		return false
	}
	for _, v := range roleCodes {
		if !utils.SliceIn(v, u.RoleCodes...) {
			return false
		}
	}
	return true
}

type InnerCtx struct {
	AllProject       bool
	AllArea          bool //内部使用,不限制区域
	AllTenant        bool //所有租户的权限
	WithCommonTenant bool //同时获取公共租户
}

func GetHandle(r *http.Request, keys ...string) string {
	var val string
	for _, v := range keys {
		val = r.Header.Get(v)
		if val != "" {
			return val
		}
		val = r.URL.Query().Get(v)
		if val != "" {
			return val
		}
	}
	return val
}

func InitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uc := GetUserCtx(r.Context())
		if uc == nil {
			strIP, _ := utils.GetIP(r)
			appCode := GetHandle(r, UserAppCodeKey)
			tenantCode := GetHandle(r, UserTenantCodeKey)
			uc = &UserCtx{
				AppCode:    appCode,
				TenantCode: tenantCode,
				IP:         strIP,
			}
			c := context.WithValue(r.Context(), UserInfoKey, uc)
			r = r.WithContext(c)
		}
		strProjectID := GetHandle(r, UserProjectID)
		projectID := cast.ToInt64(strProjectID)
		if projectID == 0 {
			projectID = def.NotClassified
		}
		uc.AppCode = GetHandle(r, UserAppCodeKey)
		if uc.AppCode == "" {
			uc.AppCode = def.AppCore
		}
		if uc.TenantCode == "" {
			uc.TenantCode = def.TenantCodeDefault
		}
		uc.ProjectID = projectID
		uc.Os = GetHandle(r, "User-Agent")
		uc.AcceptLanguage = GetHandle(r, "Accept-Language")
		uc.Token = GetHandle(r, UserTokenKey)
		ctx := SetUserCtx(r.Context(), uc)
		r = r.WithContext(ctx)
		next(w, r)
	}
}

func BindTenantCode(ctx context.Context, tenantCode string, projectID int64) context.Context {
	uc := GetUserCtx(ctx)
	if uc == nil {
		if tenantCode == "" {
			tenantCode = def.TenantCodeDefault
		}
		uc = &UserCtx{
			TenantCode: tenantCode,
			ProjectID:  projectID,
		}
		ctx = context.WithValue(ctx, UserInfoKey, uc)

	} else {
		uc.TenantCode = tenantCode
		if projectID != 0 {
			uc.ProjectID = projectID
		}
	}
	return ctx
}

func UpdateUserCtx(ctx context.Context) context.Context {
	uc := GetUserCtx(ctx)
	if uc == nil {
		return ctx
	}
	return SetUserCtx(ctx, uc)
}

func SetUserCtx(ctx context.Context, userCtx *UserCtx) context.Context {
	if userCtx == nil {
		return ctx
	}
	info, _ := json.Marshal(userCtx)
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
		UserInfoKey, base64.StdEncoding.EncodeToString(info),
	))
	return context.WithValue(ctx, UserInfoKey, userCtx)
}
func SetInnerCtx(ctx context.Context, inner InnerCtx) context.Context {
	uc := GetUserCtx(ctx)
	if uc == nil {
		return ctx
	}
	uc.InnerCtx = inner
	return SetUserCtx(ctx, uc)
}

func GetInnerCtx(ctx context.Context) InnerCtx {
	uc := GetUserCtx(ctx)
	if uc == nil {
		return InnerCtx{}
	}
	return uc.InnerCtx
}

// 使用该函数前必须传了UserCtx
func GetUserCtx(ctx context.Context) *UserCtx {
	val, ok := ctx.Value(UserInfoKey).(*UserCtx)
	if !ok { //这里线上不能获取不到
		return nil
	}
	return val
}

func GetUserCtxNoNil(ctx context.Context) *UserCtx {
	val, ok := ctx.Value(UserInfoKey).(*UserCtx)
	if !ok { //这里线上不能获取不到
		return &UserCtx{ProjectID: def.NotClassified, TenantCode: def.TenantCodeDefault}
	}
	return val
}
func WithRoot(ctx context.Context) context.Context {
	uc := *GetUserCtxNoNil(ctx)
	uc.TenantCode = def.TenantCodeDefault //只有default租户有root权限去读其他租户的数据
	uc.AllTenant = true
	uc.AllProject = true
	uc.AllArea = true
	uc.IsSuperAdmin = true
	uc.IsAdmin = true
	return SetUserCtx(ctx, &uc)
}

// 如果是default租户直接给root权限
func WithDefaultRoot(ctx context.Context) context.Context {
	uc := GetUserCtxNoNil(ctx)
	if uc.TenantCode != def.TenantCodeDefault || !uc.IsSuperAdmin || uc.ProjectID > 3 { //传了项目ID则不是root权限
		return ctx
	}
	return WithRoot(ctx)
}

// 如果是管理员,且没有传项目id,则直接给所有项目权限
func WithDefaultAllProject(ctx context.Context) context.Context {
	uc := GetUserCtxNoNil(ctx)
	if !uc.IsAdmin || uc.ProjectID > 3 { //传了项目ID则不是root权限
		return ctx
	}
	return WithAllProject(ctx)
}

func IsAdmin(ctx context.Context) error {
	uc := GetUserCtxNoNil(ctx)
	if uc.IsAdmin || uc.IsSuperAdmin {
		return nil
	}
	return errors.Permissions.AddMsg("只允许管理员操作")
}

func WithAllProject(ctx context.Context) context.Context {
	uc := *GetUserCtxNoNil(ctx)
	uc.AllProject = true
	return SetUserCtx(ctx, &uc)
}

func WithCommonTenant(ctx context.Context) context.Context {
	uc := *GetUserCtxNoNil(ctx)
	uc.WithCommonTenant = true
	return SetUserCtx(ctx, &uc)
}

func NewUserCtx(ctx context.Context) context.Context {
	val, ok := ctx.Value(UserInfoKey).(*UserCtx)
	if !ok { //这里线上不能获取不到
		return ctx
	}
	var newUc UserCtx
	newUc = *val
	return context.WithValue(context.Background(), UserInfoKey, &newUc)
}

func IsRoot(ctx context.Context) error {
	uc := GetUserCtx(ctx)
	if uc == nil || uc.TenantCode != def.TenantCodeDefault || !uc.IsSuperAdmin {
		return errors.Permissions.AddDetailf("需要超管才能操作")
	}
	return nil
}

// 使用该函数前必须传了UserCtx
func GetUserCtxOrNil(ctx context.Context) *UserCtx {
	val, ok := ctx.Value(UserInfoKey).(*UserCtx)
	if !ok { //这里线上不能获取不到
		return nil
	}
	return val
}

type MetadataCtx = map[string][]string

func SetMetaCtx(ctx context.Context, maps MetadataCtx) context.Context {
	return context.WithValue(ctx, MetadataKey, maps)
}
func GetMetaCtx(ctx context.Context) MetadataCtx {
	val, ok := ctx.Value(MetadataKey).(MetadataCtx)
	if !ok {
		return nil
	}
	return val
}

func GetMetaVal(ctx context.Context, field string) []string {
	mdCtx := GetMetaCtx(ctx)
	if val, ok := mdCtx[field]; !ok {
		return nil
	} else {
		return val
	}
}

//// 指定项目id（企业版功能）
//func SetMetaProjectID(ctx context.Context, projectID int64) {
//	mc := GetMetaCtx(ctx)
//	projectIDStr := utils.ToString(projectID)
//	mc[string(MetaFieldProjectID)] = []string{projectIDStr}
//}
//
//// 获取meta里的项目ID（企业版功能）
//func ClearMetaProjectID(ctx context.Context) {
//	mc := GetMetaCtx(ctx)
//	delete(mc, string(MetaFieldProjectID))
//}
