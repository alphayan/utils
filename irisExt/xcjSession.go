package xcjIris

import (
	"github.com/alphayan/iris/context"
	"github.com/alphayan/iris/sessions"
	"time"
)

const (
	XcjCookieName = "xcjsessionid"
	XcjUserToken  = "UserToken"
)

var (
	sessionsXcj *sessions.Sessions
	xcjSee      *sessions.Session
)

//初始化小车匠session
func NewSessionXcj() {
	if sessionsXcj == nil {
		sessionsXcj = sessions.New(sessions.Config{
			Cookie:  XcjCookieName,
			Expires: 24 * time.Hour,
		})
	}
}

//开始session
func XcjSessionStart(ctx context.Context) {
	if xcjSee == nil {
		xcjSee = sessionsXcj.Start(ctx)
	}
}

//保存登录信息
func LoginSet(ctx context.Context, val interface{}) {
	if xcjSee == nil {
		xcjSee = sessionsXcj.Start(ctx)
	}
	xcjSee.Set(XcjCookieName, val)
}

//获取当前登录信息
func GetCurrentLoginInfo() interface{} {
	return xcjSee.Get(XcjCookieName)
}

//是否登录
func IsLoggedIn() bool {
	if GetCurrentLoginInfo() != nil {
		return true
	} else {
		return false
	}
}
