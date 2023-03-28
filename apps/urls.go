package apps

import (
	"ginDemo/global"
	"ginDemo/middleware"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Option func(*gin.Engine)

var options []Option

// Include 注册app的路由配置
func Include(opts ...Option) {
	options = append(options, opts...)
}

// Init 初始化
func Init(ENV string) *gin.Engine {
	// 生产中启动，关闭DeBug模式
	if global.CONFIG.GetBool(ENV + ".DeBug") {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 启用性能分析工具
	pprof.Register(r)
	r.Use(
		// 跨域配置
		middleware.Cors(),
		middleware.JwtCheck(),

		// 限速中间件，初始token100，每秒增加100
		middleware.RateLimitMiddleware(time.Second, 100, 10),

		// 日志中间件
		middleware.GinLogger(),
		middleware.GinRecovery(true),
	)

	r.LoadHTMLGlob("templates/**/*")
	r.Static("/assets", "./assets")
	r.StaticFS("/logs", http.Dir("./logs"))
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "views/404.html", nil)
	})

	for _, opt := range options {
		opt(r)
	}
	return r
}
