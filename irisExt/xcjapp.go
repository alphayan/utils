package xcjIris

import (
	"github.com/alphayan/iris"
	"github.com/alphayan/iris/context"
	"github.com/alphayan/iris/core/errors"
	"github.com/alphayan/iris/middleware/logger"
	"github.com/thoas/go-funk"
	"utils/config"
	"utils/handlers"
	"utils/xcj/redis"
	"utils/xcjswagger"
)

const (
	ApiVersion = "v1"
	SessionKey = "xcjSession"
)

type XcjApper interface {
	NewApp() (*iris.Application, iris.Party)
}
type xcjApp struct {
	isLoad     bool
	app        *iris.Application
	apiVersion iris.Party
}

var (
	_XcjApp xcjApp
)

//检查是否实现了XcjApper接口
var _ XcjApper = new(xcjApp)

func (app *xcjApp) NewApp() (*iris.Application, iris.Party) {
	if !_XcjApp.isLoad {
		_XcjApp.app = iris.New()
		_XcjApp.isLoad = true
		xcjRedis.LoadClient(config.Config.Redis)
		BuildIrisSession(config.Config)
		if config.Config.Debug {
			_XcjApp.app.Logger().Debug(true)
			//定义客户日志
			customLogger := logger.New(logger.Config{
				// Status displays status code
				Status: true,
				// IP displays request's remote address
				IP: true,
				// Method displays the http method
				Method: true,
				// Path displays the request path
				Path: true,

				//Columns: true,

				// if !empty then its contents derives from `ctx.Values().Get("logger_message")
				// will be added to the logs.
				MessageContextKey: "logger_message",
			})
			_XcjApp.app.Use(customLogger)
			_XcjApp.app.OnAnyErrorCode(customLogger, func(ctx iris.Context) {
				// this should be added to the logs, at the end because of the `logger.Config#MessageContextKey`
				ctx.Values().Set("logger_message",
					"the url not found in app router")
				ctx.Writef("The XCJ API NOT FOUND")
			})
		}
		_XcjApp.apiVersion = _XcjApp.app.Party("/"+ApiVersion, func(ctx context.Context) {
			//不需要登录的地址
			urls := []string{"/v1/user/login", "/v1/user/logout",
				"/v1/xq/school/receiveschooldata/class/info",
				"/v1/xq/school/receiveschooldata/class/resume/teacher/info",
				"/v1/xq/school/receiveschooldata/practice/task/catalog",
				"/v1/xq/school/receiveschooldata/practice/task/type",
				"/v1/xq/school/receiveschooldata/practice/task",
				"/v1/xq/school/receiveschooldata/practice/practise/plan",
				"/v1/xq/school/receiveschooldata/practice/practise/plan/task",
				"/v1/xq/school/receiveschooldata/practice/practise/plan/class/user",
				"/v1/xcj/practice/all/school",
				"/v1/xq/school/receiveschooldata/class/update",
				"/v1/xq/school/receiveschooldata/class/import"}
			if funk.ContainsString(urls, ctx.Request().URL.Path) {
				ctx.Next()
				return
			}
			sess, _ := NewSessionFromIris(ctx, SessionKey)
			if !sess.IsValued() {
				handlers.AfterHandlersJSON(ctx, nil, errors.New("您尚未登录！"))
				return
			} else {
				//已经登录过后将session更新一下
				sess, _ := NewSessionFromIris(ctx, SessionKey)
				sess.SaveIris(ctx, SessionKey)
				ctx.Next()
			}
		})
		_XcjApp.app.StaticWeb("/static", config.Config.StaticPath)
		if config.Config.Swagger.Enabledoc {
			_XcjApp.app.StaticWeb("/swagger", config.Config.Swagger.Swaggerpath)
			xcjswag.InitSwagger(_XcjApp.app)
		}
		_XcjApp.app.StaticWeb("/", config.Config.WebApps)
	}
	return _XcjApp.app, _XcjApp.apiVersion
}

func NewXcjApp() XcjApper {
	if _XcjApp.isLoad {
		return &_XcjApp
	}
	_XcjApp := new(xcjApp)
	return _XcjApp
}
