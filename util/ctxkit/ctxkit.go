// Package ctxkit 操作请求 ctx 信息
package ctxkit

import (
	"context"
	"net/http"
	"strings"
)

type key int

const (
	// TraceIDKey 请求唯一标识，类型：string
	TraceIDKey key = iota
	// StartTimeKey 请求开始时间，类型：time.Time
	StartTimeKey
	// UserTokenKey 用户登陆身份，未登录则为 ""，类型：string
	UserTokenKey
	// UserIDKey 用户 ID，未登录则为 0，类型：int64
	UserIDKey
	// UserIPKey 用户 IP，类型：string
	UserIPKey
	// PlatformKey 用户使用平台，ios, android, pc
	PlatformKey
	// ProjectKey
	ProjectKey
	// ProjectDBKey
	ProjectDBNameKey
	// BuildKey 客户端构建版本号
	BuildKey
	// VersionKey 客户端版本号
	VersionKey
	// AccessKeyKey 移动端支付令牌
	AccessKeyKey
	// DeviceKey 移动 app 设备标识，android, phone, pad
	DeviceKey
	// MobiAppKey 移动 app 标识，android, phone, pad
	MobiAppKey
	// UserPortKey 用户端口
	UserPortKey
	// ManageUserKey 管理后台用户名
	ManageUserKey
	// BuvidKey 非登录用户标识
	BuvidKey
	// CookieKey web 用户登录令牌
	CookieKey
	// CompanyAppKeyKey
	CompanyAppKeyKey
	// AppKeyKey 接口签名标识
	AppKeyKey
	// TsKey 时间戳
	TSKey
	// SignKey 签名
	SignKey
	// IsValidSignKeyKey 签名正确则置为 true
	IsValidSignKeyKey
	// http 请求对象
	HttpRawReqKey
	// http 请求响应对象
	HttpRawRespKey
	// User-Agent
	UserAgentKey
	// HttpRefer
	HttpReferKey
	// WxSessionKey
	WxSessionKey
	// WxOpenIdKey
	WxOpenIdKey
	// WxPhoneBindKey
	WxPhoneBindKey
	// MpPlatformIdKey
	MpPlatformIdKey
	// WXAPPIdKey
	WXAPPIdKey
	// QQAPPIdKey
	QQAPPIdKey
)

// GetUserID 获取当前登录用户 ID
func GetUserID(ctx context.Context) int64 {
	uid, _ := ctx.Value(UserIDKey).(int64)
	return uid
}

func WithUserID(ctx context.Context, uid int64) context.Context {
	return context.WithValue(ctx, UserIDKey, uid)
}

func WithMPPlatformId(ctx context.Context, pid int32) context.Context {
	return context.WithValue(ctx, MpPlatformIdKey, pid)
}

func WithDevice(ctx context.Context, device string) context.Context {
	return context.WithValue(ctx, DeviceKey, device)
}

func WithWXAppId(ctx context.Context, appid string) context.Context {
	return context.WithValue(ctx, WXAPPIdKey, appid)
}

func WithQQAppId(ctx context.Context, appid string) context.Context {
	return context.WithValue(ctx, QQAPPIdKey, appid)
}

//func WithWXAPI(ctx context.Context, api *wechat.WechatAPI) context.Context {
//	return context.WithValue(ctx, WXAPIKEY, api)
//}

//func WithQQAPI(ctx context.Context, api *qq.QQAPI) context.Context {
//	return context.WithValue(ctx, QQAPIKEY, api)
//}

func WithMobiApp(ctx context.Context, os string) context.Context {
	return context.WithValue(ctx, MobiAppKey, os)
}

func WithVersion(ctx context.Context, version string) context.Context {
	return context.WithValue(ctx, VersionKey, version)
}

func WithPlatform(ctx context.Context, platform string) context.Context {
	return context.WithValue(ctx, PlatformKey, platform)
}

func WithProject(ctx context.Context, project string) context.Context {
	return context.WithValue(ctx, ProjectKey, project)
}

func WithProjectDBNameKey(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, ProjectDBNameKey, name)
}

func WithUserAgent(ctx context.Context, ua string) context.Context {
	return context.WithValue(ctx, UserAgentKey, ua)
}

func WithHttpRefer(ctx context.Context, refer string) context.Context {
	return context.WithValue(ctx, HttpReferKey, refer)
}

func WithWxSessionKey(ctx context.Context, sessionKey string) context.Context {
	return context.WithValue(ctx, WxSessionKey, sessionKey)
}

func WithWxOpenId(ctx context.Context, openid string) context.Context {
	return context.WithValue(ctx, WxOpenIdKey, openid)
}

func WithWxPhoneBind(ctx context.Context, bind int) context.Context {
	return context.WithValue(ctx, WxPhoneBindKey, bind)
}

// GetWxSession 获取用户 session key
func IsWxPhoneBind(ctx context.Context) bool {
	bind, _ := ctx.Value(WxPhoneBindKey).(int)
	return bind == 1
}

// GetWxSession 获取用户 session key
func GetWxOpenId(ctx context.Context) string {
	openid, _ := ctx.Value(WxOpenIdKey).(string)
	return openid
}

func GetWXAppId(ctx context.Context) string {
	appid, _ := ctx.Value(WXAPPIdKey).(string)
	return appid
}

func GetQQAppId(ctx context.Context) string {
	appid, _ := ctx.Value(QQAPPIdKey).(string)
	return appid
}


// GetWxSession 获取用户 session key
func GetWxSessionKey(ctx context.Context) string {
	sessionKey, _ := ctx.Value(WxSessionKey).(string)
	return sessionKey
}

// GetUserAgent 获取用户 refer
func GetHttpRefer(ctx context.Context) string {
	ua, _ := ctx.Value(HttpReferKey).(string)
	return ua
}

// GetUserAgent 获取用户 ua
func GetUserAgent(ctx context.Context) string {
	ua, _ := ctx.Value(UserAgentKey).(string)
	return ua
}

// GetUserIP 获取用户 IP
func GetUserIP(ctx context.Context) string {
	ip, _ := ctx.Value(UserIPKey).(string)
	return ip
}

func WithUserIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, UserIPKey, ip)
}

// GetUserToken 获取用户 token
func GetUserToken(ctx context.Context) string {
	token, _ := ctx.Value(UserTokenKey).(string)
	return token
}

// GetUserToken 获取用户 token
func GetTokenPlatformId(ctx context.Context) int32 {
	pid, _ := ctx.Value(MpPlatformIdKey).(int32)
	return pid
}

func WithUserToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, UserTokenKey, token)
}

// GetUserPort 获取用户端口
func GetUserPort(ctx context.Context) string {
	port, _ := ctx.Value(UserPortKey).(string)
	return port
}

// GetPlatform 获取用户平台
func GetPlatform(ctx context.Context) string {
	platform, _ := ctx.Value(PlatformKey).(string)
	return platform
}

// GetProject 获取配置项目
func GetProject(ctx context.Context) string {
	project, _ := ctx.Value(ProjectKey).(string)
	return project
}

// GetProjectDBName 获取配置项目
func GetProjectDBName(ctx context.Context) string {
	name, _ := ctx.Value(ProjectDBNameKey).(string)
	if len(name) == 0 {
		return "default"
	}
	return strings.ToUpper(name)
}

// GetPlatform 获取用户平台
func GetCompanyAppKey(ctx context.Context) string {
	companyKey, _ := ctx.Value(CompanyAppKeyKey).(string)
	return companyKey
}

// IsIOSPlatform 判断是否为 IOS 平台
func IsIOSPlatform(ctx context.Context) bool {
	return GetPlatform(ctx) == "ios"
}

// GetTraceID 获取用户请求标识
func GetTraceID(ctx context.Context) string {
	id, _ := ctx.Value(TraceIDKey).(string)
	return id
}

// GetBuild 获取客户端构建版本号
func GetBuild(ctx context.Context) string {
	build, _ := ctx.Value(BuildKey).(string)
	return build
}

// GetDevice 获取用户设备，配合 GetPlatform 使用
func GetDevice(ctx context.Context) string {
	device, _ := ctx.Value(DeviceKey).(string)
	return device
}

// GetMobiApp 获取 APP 标识
func GetMobiApp(ctx context.Context) string {
	app, _ := ctx.Value(MobiAppKey).(string)
	return app
}

// GetVersion 获取客户端版本
func GetVersion(ctx context.Context) string {
	version, _ := ctx.Value(VersionKey).(string)
	return version
}

// GetAccessKey 获取客户端认证令牌
func GetAccessKey(ctx context.Context) string {
	key, _ := ctx.Value(AccessKeyKey).(string)
	return key
}

// WithTraceID 注入 trace_id
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// WithHttpReq 注入 HttpRawReqKey
func WithHttpReq(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, HttpRawReqKey, req)
}

// WithHttpResp 注入 TraceIDKey
func WithHttpResp(ctx context.Context, resp http.ResponseWriter) context.Context {
	return context.WithValue(ctx, HttpRawRespKey, resp)
}

// GetManageUser 获取管理后台用户名
func GetManageUser(ctx context.Context) string {
	user, _ := ctx.Value(ManageUserKey).(string)
	return user
}

// GetBuvid 获取用户 buvid
func GetBuvid(ctx context.Context) string {
	buvid, _ := ctx.Value(BuvidKey).(string)
	return buvid
}

// GetCookie 获取 web cookie
func GetCookie(ctx context.Context) string {
	key, _ := ctx.Value(CookieKey).(string)
	return key
}

func GetSign(ctx context.Context) (appkey, ts, sign string) {
	appkey, _ = ctx.Value(AppKeyKey).(string)
	ts, _ = ctx.Value(TSKey).(string)
	sign, _ = ctx.Value(SignKey).(string)
	return
}

// IsValidSignKey 判断业务签名是否正确
func IsValidSignKey(ctx context.Context) bool {
	valid, _ := ctx.Value(IsValidSignKeyKey).(bool)
	return valid
}

// GetHttpReq 获取 web cookie
func GetHttpReq(ctx context.Context) *http.Request {
	key, _ := ctx.Value(HttpRawReqKey).(*http.Request)
	return key
}

// GetCookie 获取 web cookie
func GetHttpResp(ctx context.Context) http.ResponseWriter {
	key, _ := ctx.Value(HttpRawRespKey).(http.ResponseWriter)
	return key
}
