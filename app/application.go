package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	// PRODUCTION 生产环境
	PRODUCTION = "production"
	// DEVELOPMENT 开发环境
	DEVELOPMENT = "development"
)

// Application 应用
type Application struct {
	Config      *viper.Viper
	Logger      *logrus.Logger
	Environment string

	httpServer     *gin.Engine
	initHooks      []ApplicationHook
	beforeRunHooks []ApplicationHook
	destroyHooks   []ApplicationHook
}

// ApplicationHook 生命周期中调用的钩子函数
type ApplicationHook func(Application) error

// Default 返回默认的应用
func Default() *Application {
	o := &Application{
		Config: viper.GetViper(),
		Logger: logrus.New(),

		initHooks:      make([]ApplicationHook, 0),
		beforeRunHooks: make([]ApplicationHook, 0),
		destroyHooks:   make([]ApplicationHook, 0),
	}

	env := o.Config.GetString("environment")
	switch env {
	case PRODUCTION:
	case DEVELOPMENT:
		o.Environment = env
		break
	default:
		// FIXME: 添加环境参数不支持警告
		o.Environment = PRODUCTION
		break
	}

	if o.Environment == PRODUCTION {
		gin.SetMode(gin.ReleaseMode)

		// Logger
		o.Logger.SetFormatter(&logrus.JSONFormatter{})
		o.Logger.SetLevel(logrus.WarnLevel)
	} else {
		o.Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		o.Logger.SetLevel(logrus.DebugLevel)
	}
	o.Logger.SetNoLock()

	o.httpServer = gin.Default()

	return o
}

// Group 代理 gin.Engine 的 Group 函数
func (o *Application) Group(path string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return o.httpServer.Group(path, handlers...)
}

// Run 运行服务
func (o *Application) Run() {
	errc := make(chan error)

	go func() { errc <- o.runHTTPServer() }()

	o.Logger.Panic(<-errc)
}

func (o *Application) runHTTPServer() error {
	addr := o.Config.GetString("http_addr")
	if len(addr) == 0 {
		addr = ":8080"
	}
	o.Logger.WithField("addr", addr).Info("Running HTTP server...")
	return o.httpServer.Run(addr)
}
