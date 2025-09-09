package ctxs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

type UserCtx struct {
	IsOpen         bool `json:",omitempty"` //是否开放认证用户
	AppCode        string
	DeviceID       string
	Token          string   `json:",omitempty"`
	TenantCode     string   //租户Code
	AcceptLanguage string   `json:",omitempty"`
	ProjectID      int64    `json:",string"`
	IsAdmin        bool     `json:",omitempty"` //是否是超级管理员
	IsSuperAdmin   bool     `json:",omitempty"`
	UserID         int64    `json:",string"`    //用户id（开放认证用户值为0）
	RoleIDs        []int64  `json:",omitempty"` //用户使用的角色（开放认证用户值为0）
	RoleCodes      []string `json:",omitempty"`
	IsAllData      bool     `json:",omitempty"` //是否所有数据权限（开放认证用户值为true）
	IP             string   `json:",omitempty"` //用户的ip地址
	Os             string   `json:",omitempty"` //操作系统
	UserName       string
	Account        string
	ProjectAuth    map[int64]*ProjectAuth  `json:",omitempty"`
	Dept           map[int64]def.AuthType  `json:",omitempty"` //key是区域ID,value是授权类型
	DeptPath       map[string]def.AuthType `json:",omitempty"` //key是区域ID路径,value是授权类型
	InnerCtx
}

type ProjectAuth struct {
	Area     map[int64]def.AuthType  `json:",omitempty"` //key是区域ID,value是授权类型
	AreaPath map[string]def.AuthType `json:",omitempty"` //key是区域ID路径,value是授权类型
	// 1 //管理权限,可以修改别人的权限,及读写权限 管理权限不限制区域权限
	// 2 //读权限,只能读,不能修改
	// 3 //读写权限,可以读写该权限
	AuthType def.AuthType `json:",omitempty"` //项目的授权类型
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

func GetAreaIDPaths(projectID int64, in map[int64]*ProjectAuth) (authType def.AuthType, areas []string) {
	v := in[projectID]
	if v == nil {
		return
	}
	authType = v.AuthType
	for area := range v.AreaPath {
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

func (u *UserCtx) IsRoot() bool {
	if u == nil || u.TenantCode != def.TenantCodeDefault || !u.IsAdmin {
		return false
	}
	return true
}

type InnerCtx struct {
	AllProject       bool `json:",omitempty"`
	AllArea          bool `json:",omitempty"` //内部使用,不限制区域
	AllTenant        bool `json:",omitempty"` //所有租户的权限
	WithCommonTenant bool `json:",omitempty"` //同时获取公共租户
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

func InitCtxWithReq(r *http.Request) *http.Request {
	uc := GetUserCtx(r.Context())
	if uc == nil {
		strIP, _ := utils.GetIP(r)
		appCode := GetHandle(r, UserAppCodeKey, UserAppCodeKey2)
		tenantCode := GetHandle(r, UserTenantCodeKey, UserTenantCodeKey2)
		uc = &UserCtx{
			AppCode:    appCode,
			TenantCode: tenantCode,
			IP:         strIP,
		}
		c := context.WithValue(r.Context(), UserInfoKey, uc)
		r = r.WithContext(c)
	}
	strProjectID := GetHandle(r, UserProjectID, UserProjectID2)
	projectID := cast.ToInt64(strProjectID)
	if projectID == 0 {
		projectID = def.NotClassified
	}
	uc.DeviceID = GetHandle(r, UserDeviceIDKey)
	uc.AppCode = GetHandle(r, UserAppCodeKey, UserAppCodeKey2)
	if uc.AppCode == "" {
		uc.AppCode = def.AppCore
	}
	if uc.TenantCode == "" {
		uc.TenantCode = def.TenantCodeCommon
	}
	uc.ProjectID = projectID
	uc.Os = GetHandle(r, "User-Agent")
	uc.AcceptLanguage = GetHandle(r, "Accept-Language")
	uc.Token = GetHandle(r, UserTokenKey, UserToken2Key)
	if uc.IP == "" {
		strIP, _ := utils.GetIP(r)
		uc.IP = strIP
	}
	ctx := SetUserCtx(r.Context(), uc)
	r = r.WithContext(ctx)
	return r
}

func BindUser(ctx context.Context, userID int64, account string) context.Context {
	uc := GetUserCtxNoNil(ctx)
	uc.UserID = userID
	uc.Account = account
	return SetUserCtx(ctx, uc)
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
		return SetUserCtx(ctx, uc)
	} else {
		ucc := *uc
		ucc.TenantCode = tenantCode
		ucc.AllTenant = false
		if projectID != 0 {
			ucc.ProjectID = projectID
			ucc.AllProject = true
		}
		return SetUserCtx(ctx, &ucc)
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

type ctxStr struct {
	UserCtx *UserCtx
	Trace   string
}

func ToString(ctx context.Context) string {
	uc := GetUserCtx(ctx)
	span := trace.SpanFromContext(ctx)
	traceinfo, _ := span.SpanContext().MarshalJSON()
	ctxstr := ctxStr{
		UserCtx: uc,
		Trace:   string(traceinfo),
	}
	return utils.MarshalNoErr(ctxstr)
}

type mySpanContextConfig struct {
	TraceID string
	SpanID  string
}

func StringParse(ctx context.Context, str string) (context.Context, bool) {
	var cs ctxStr
	err := json.Unmarshal([]byte(str), &cs)
	if err != nil {
		return ctx, false
	}
	var msg mySpanContextConfig
	err = json.Unmarshal([]byte(cs.Trace), &msg)
	if err != nil {
		logx.Errorf("[GetCtx]|json Unmarshal trace.SpanContextConfig err:%v", err)
		return ctx, false
	}
	//将MsgHead 中的msg链路信息 重新注入ctx中并返回
	t, err := trace.TraceIDFromHex(msg.TraceID)
	if err != nil {
		logx.Errorf("[GetCtx]|TraceIDFromHex err:%v", err)
		return ctx, false
	}
	s, err := trace.SpanIDFromHex(msg.SpanID)
	if err != nil {
		logx.Errorf("[GetCtx]|SpanIDFromHex err:%v", err)
		return ctx, false
	}
	parent := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    t,
		SpanID:     s,
		TraceFlags: 0x1,
	})
	ctx2 := trace.ContextWithRemoteSpanContext(ctx, parent)
	return SetUserCtx(ctx2, cs.UserCtx), true

}

func SetUserCtx(ctx context.Context, userCtx *UserCtx) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
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

func GetAppCode(ctx context.Context) string {
	uc := GetUserCtx(ctx)
	if uc == nil {
		return ""
	}
	return uc.AppCode
}

func GenGrpcurlHandle(ctx context.Context) string {
	uc := GetUserCtx(ctx)
	if uc == nil {
		return ""
	}
	info, _ := json.Marshal(uc)
	return UserInfoKey + ":" + base64.StdEncoding.EncodeToString(info)
}

func GetUserCtxNoNil(ctx context.Context) *UserCtx {
	if ctx == nil {
		ctx = context.Background()
	}
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

func WithAdmin(ctx context.Context) context.Context {
	uc := *GetUserCtxNoNil(ctx)
	uc.AllProject = true
	uc.AllArea = true
	uc.IsAdmin = true
	return SetUserCtx(ctx, &uc)
}

func WithProjectID(ctx context.Context, projectID int64) context.Context {
	uc := *GetUserCtxNoNil(ctx)
	uc.AllProject = false
	uc.ProjectID = projectID
	return SetUserCtx(ctx, &uc)
}

// 如果是default租户直接给root权限
func WithDefaultRoot(ctx context.Context) context.Context {
	uc := GetUserCtxNoNil(ctx)
	if uc.TenantCode != def.TenantCodeDefault || !uc.IsAdmin || uc.ProjectID > 3 { //传了项目ID则不是root权限
		return ctx
	}
	return WithRoot(ctx)
}

func CommonWithDefault(ctx context.Context) context.Context {
	uc := GetUserCtxNoNil(ctx)
	if uc.TenantCode != def.TenantCodeCommon { //传了项目ID则不是root权限
		return ctx
	}
	return BindTenantCode(ctx, def.TenantCodeDefault, 0)
}

func CommonWithRoot(ctx context.Context) context.Context {
	uc := GetUserCtxNoNil(ctx)
	if uc.TenantCode != def.TenantCodeCommon { //传了项目ID则不是root权限
		return ctx
	}
	return WithRoot(ctx)
}

func IsTenantDefault(ctx context.Context) bool {
	uc := GetUserCtxNoNil(ctx)
	return uc.TenantCode == def.TenantCodeDefault
}

func CanHandleTenantCommon[t ~string](ctx context.Context, tenantCode t) bool {
	uc := GetUserCtxNoNil(ctx)
	if !(uc.TenantCode == def.TenantCodeDefault || string(tenantCode) == uc.TenantCode) {
		return false
	}
	return true
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

func WithAllArea(ctx context.Context) context.Context {
	uc := *GetUserCtxNoNil(ctx)
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
	newCtx := context.WithValue(context.Background(), UserInfoKey, &newUc)
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		newCtx = metadata.NewOutgoingContext(newCtx, md)
	}
	return newCtx
}

func IsRoot(ctx context.Context) error {
	uc := GetUserCtx(ctx)
	if !uc.IsRoot() {
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
