package xcjIris

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/json-iterator/go"
	"github.com/alphayan/iris"
	"github.com/alphayan/iris/context"
	"github.com/alphayan/iris/sessions"
	"github.com/alphayan/iris/sessions/sessiondb/redis"
	"github.com/alphayan/iris/sessions/sessiondb/redis/service"
	"net/http"
	"net/http/httptest"
	"time"
	"utils/config"
	"utils/xcj/redis"
)

const (
	CookieName = "xcjsessionid"
	UserToken  = "UserToken"
)

type XcjSession struct {
	ctx      context.Context
	Iris     *sessions.Session `json:"-"`
	key      string //保存在iris中的键
	isClear  bool
	Sid      string
	Uid      uint32
	RoleType uint32
	Username string
	LastTime time.Time
	SchoolId uint32 //班级id
}

type IXcjSession interface {
	IsClear() bool
	Clear()
	IsValued() bool
	Refresh()
	RefreshAuto()
}

var (
	_session          *sessions.Sessions
	_errNotSet        = errors.New("utils.session not set")
	_rdb              *redis.Database
	_sessionKeyPrefix = "xcj_"
)

func BuildIrisSession(conf config.DLZConfig) {
	rconf := conf.Redis
	_rdb = redis.New(service.Config{
		Network:     service.DefaultRedisNetwork,
		Addr:        rconf.Addr,
		Password:    rconf.Password,
		Database:    rconf.Database,
		MaxIdle:     0,
		MaxActive:   0,
		IdleTimeout: service.DefaultRedisIdleTimeout,
		Prefix:      _sessionKeyPrefix})

	_rdb.Async(true)
	iris.RegisterOnInterrupt(func() {
		_rdb.Close()
	})

	exp := time.Duration(conf.CookieExpires)
	if exp > 0 {
		exp = exp * time.Second
	}
	mySessions := sessions.New(sessions.Config{
		Cookie:       CookieName,
		Expires:      exp,
		AllowReclaim: true,
	})
	mySessions.UseDatabase(_rdb)

	SetIrisSessionFactory(mySessions)
	return
}

func SetIrisSessionFactory(s *sessions.Sessions) {
	if _session != nil {
		return
	}
	_session = s
}

func checkSession() bool {
	return _session != nil
}

func GetIrisSession(ctx context.Context) *sessions.Session {
	sess := _session.Start(ctx)
	return sess
}
func GetSessionCtx(key string) context.Context {
	ctx := context.NewContext(nil)
	req := httptest.NewRequest("GET", "http://localhost", nil)
	req.AddCookie(&http.Cookie{Name: "smsid", Value: key})
	w := httptest.NewRecorder()
	ctx.BeginRequest(w, req)
	return ctx
}

func GetIrisSessionByKey(key string) *sessions.Session {
	ctx := GetSessionCtx(key)
	defer func() {
		ctx.EndRequest()
	}()
	return GetIrisSession(ctx)
}

func NewSessionByKey(key string) (*XcjSession, error) {
	ctx := GetSessionCtx(key)
	defer func() {
		ctx.EndRequest()
	}()
	return NewSessionFromIris(ctx, UserToken)
}

func NewSessionFromIris(ctx context.Context, key string) (*XcjSession, error) {
	if !checkSession() {
		return nil, _errNotSet
	}
	sess := _session.Start(ctx)
	val := sess.Get(key)
	if val == nil {
		return &XcjSession{ctx: ctx, key: key, Iris: sess}, nil
	}
	xsess, err := NewSessionFromGob(val.([]byte))
	xsess.RefreshAuto()
	xsess.ctx = ctx
	xsess.key = key
	xsess.Iris = sess
	xsess.Sid = sess.ID()
	return xsess, err
}

func NewSessionFromGob(bs []byte) (*XcjSession, error) {
	//buf := bytes.NewBuffer(bs)
	//dec := gob.NewDecoder(buf)
	sess := &XcjSession{}
	//err := dec.Decode(&sess)
	err := jsoniter.Unmarshal(bs, sess)
	return sess, err
}

func (x *XcjSession) ToGob() ([]byte, error) {
	//var buf bytes.Buffer
	//enc := gob.NewEncoder(&buf)
	////make a copy
	//dx := &XSession{}
	//copier.Copy(dx, x)
	//dx.ctx = nil
	//dx.Iris = nil
	//err := enc.Encode(dx)
	//return buf.Bytes(), err
	return json.Marshal(x)
}

func (x *XcjSession) UpdateExpiration(expires time.Duration) {
	_session.UpdateExpiration(x.ctx, expires)
}

func (x *XcjSession) IsClear() bool {
	return x.isClear
}

func (x *XcjSession) Clear() {
	newSess := XcjSession{ctx: x.ctx, key: x.key}
	copier.Copy(x, newSess)
	x.isClear = true
}

func (x *XcjSession) SaveIris(ctx context.Context, key string) error {
	if !checkSession() {
		return _errNotSet
	}
	//被清空的不需要保存
	if x.isClear {
		_session.Destroy(ctx)
		x.isClear = false
	}

	g, err := x.ToGob()
	sess := x.Iris
	if sess == nil {
		sess = _session.Start(ctx)
		x.Iris = sess
	}
	sess.Set(key, g)
	return err
}

// 直接保存
func (x *XcjSession) SaveIrisD() error {
	if x.ctx == nil || x.key == "" {
		return errors.New("XSession ctx or key is not set")
	}

	return x.SaveIris(x.ctx, x.key)
}

func (x *XcjSession) IsValued() bool {
	return x.Uid > 0
}

func (x *XcjSession) Refresh() {
	//acc := x.Account()
	//if acc == nil {
	//	//清空
	//	x.Clear()
	//} else {
	//	x.Group = acc.Group
	//	x.LastTime = time.Now()
	//}
}

//自动session检查
func (x *XcjSession) RefreshAuto() {
	//5分钟检查一次session
	//用户可能被删除、更新的情况
	if x.Uid > 0 && time.Now().After(x.LastTime.Add(5*time.Minute)) {
		x.Refresh()
	}
}

//删除所有的用户登录session
func ClearUserAllSessions(uid uint32) (err error) {
	userKey := fmt.Sprintf("core_user_%d_sids", uid)
	sids := []string{}
	err = xcjRedis.MainClient.Csmembers(userKey, &sids)
	if err != nil {
		return
	}
	for i, v := range sids {
		_session.DestroyByID(v)
		if i == 0 { //初始化
			_rdb.Load("")
		}
		//手动调用，可能不在内存里
		_rdb.Sync(sessions.SyncPayload{
			SessionID: v,
			Action:    sessions.ActionDestroy,
		})
	}
	//删除自己
	xcjRedis.MainClient.Del(userKey)
	return
}

//获取session
func GetXcjSession(ctx context.Context) (*XcjSession, error) {
	return NewSessionFromIris(ctx, SessionKey)
}

func (x *XcjSession) Destroy(uid uint32, ctx context.Context) {
	_session.Destroy(ctx)
}
