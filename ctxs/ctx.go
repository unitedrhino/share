package ctxs

import (
	"context"
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
	ProjectID      int64  `json:",string"`
	IsAdmin        bool   //是否是超级管理员
	UserID         int64  `json:",string"` //用户id（开放认证用户值为0）
	RoleID         int64  //用户使用的角色（开放认证用户值为0）
	IsAllData      bool   //是否所有数据权限（开放认证用户值为true）
	IP             string //用户的ip地址
	Os             string //操作系统
	UserName       string
	Account        string
	InnerCtx
}

func (u *UserCtx) ClearInner() *UserCtx {
	if u == nil {
		return nil
	}
	newCtx := *u
	newCtx.InnerCtx = InnerCtx{}
	return &newCtx
}

type InnerCtx struct {
	AllProject       bool
	AllArea          bool //内部使用,不限制区域
	AllTenant        bool //所有租户的权限
	WithCommonTenant bool //同时获取公共租户
}

func InitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uc := GetUserCtx(r.Context())
		if uc == nil {
			strIP, _ := utils.GetIP(r)
			appCode := r.Header.Get(UserAppCodeKey)
			tenantCode := r.Header.Get(UserTenantCodeKey)
			uc = &UserCtx{
				AppCode:    appCode,
				TenantCode: tenantCode,
				IP:         strIP,
			}
			c := context.WithValue(r.Context(), UserInfoKey, uc)
			r = r.WithContext(c)
		}
		strProjectID := r.Header.Get(UserProjectID)
		projectID := cast.ToInt64(strProjectID)
		if projectID == 0 {
			projectID = def.NotClassified
		}
		uc.AppCode = r.Header.Get(UserAppCodeKey)
		if uc.AppCode == "" {
			uc.AppCode = def.AppCore
		}
		if uc.TenantCode == "" {
			uc.TenantCode = def.TenantCodeDefault
		}
		uc.ProjectID = projectID
		uc.Os = r.Header.Get("User-Agent")
		uc.AcceptLanguage = r.Header.Get("Accept-Language")
		uc.Token = r.Header.Get(UserTokenKey)
		next(w, r)
	}
}

func BindTenantCode(ctx context.Context, tenantCode string) context.Context {
	uc := GetUserCtx(ctx)
	if uc == nil {
		if tenantCode == "" {
			tenantCode = def.TenantCodeDefault
		}
		uc = &UserCtx{
			TenantCode: tenantCode,
		}
		ctx = context.WithValue(ctx, UserInfoKey, uc)

	} else {
		uc.TenantCode = tenantCode
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
		UserInfoKey, string(info),
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
	return SetUserCtx(ctx, &uc)
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
	if uc == nil || uc.TenantCode != def.TenantCodeDefault {
		return errors.Permissions.AddDetailf("需要主租户才能操作")
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
