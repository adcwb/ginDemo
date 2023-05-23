package main

import (
	"ginDemo/apps"
	"ginDemo/apps/count"
	"ginDemo/apps/dingtalk"
	"ginDemo/apps/pays"
	"ginDemo/apps/test"
	"ginDemo/apps/users"
	"ginDemo/apps/wechat"
	"ginDemo/global"
	"ginDemo/initialization"
	"ginDemo/middleware"
	"ginDemo/utils"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

//	@title			ginDemo
//	@version		1.0
//	@description	此项目用于学习Golang的gin框架
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	root
//	@contact.email	root@adcwb.com
//	@contact.url	http://www.swagger.io/support

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8000
//	@BasePath	/

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	// 加载多个APP的路由配置
	apps.Include(
		users.Routers,
		dingtalk.Routers,
		pays.Routers,
		count.Routers,
		test.Routers,
		wechat.Routers,
	)

	// 初始化配置文件
	initialization.InitConfigFile()

	// 获取配置文件中需要使用的环境
	ENV := global.CONFIG.GetString("RunConfig")

	// 初始化日志
	if err := middleware.InitLogger(); err != nil {
		zap.L().Error("日志模块加载失败！")
	}
	var err error
	// 初始化数据库
	//global.DB, _ = initialization.InitMySqlClient(global.CONFIG.GetString("RunConfig"))

	// 数据迁移
	//err := global.DB.AutoMigrate(
	//	&users.User{},
	//	&pays.PayConfigData{},
	//	&pays.PayData{},
	//)
	//if err != nil {
	//	zap.L().Error("数据库自动迁移失败！", zap.Error(err))
	//}
	// 初始化MongoDB数据库
	initialization.InitMongoDBClient(ENV)

	// 初始化Redis数据库
	initialization.InitRedisClient(ENV)

	// 初始化InfluxDB数据库
	//initialization.InitInfluxDB(global.CONFIG.GetString("RunConfig"))

	// 初始化RabbitMQ消息队列
	initialization.InitRabbitMQ(ENV)

	// 初始化企业微信所需要的参数
	utils.InitWorkWechatData(global.CONFIG.GetString(ENV + ".WechatOpenKfId"))

	// 初始化阿里云OSS存储
	//initialization.InitAliYunOss(global.CONFIG.GetString("RunConfig"))

	// 初始化定时任务
	//global.JobS = gocron.NewScheduler(time.UTC)

	// 运行调度任务，共有两种方式
	//global.JobS.StartAsync() // 异步启动调度器
	//global.JobS.StartBlocking() // 启动调度器并阻塞当前执行路径

	// 初始化路由
	r := apps.Init(ENV)
	//initialization.InitSocketIO()

	if global.CONFIG.GetBool(ENV + ".RunningTLS") {
		// 启用https
		if err = r.RunTLS(
			global.CONFIG.GetString(ENV+".host")+":"+strconv.Itoa(global.CONFIG.GetInt(ENV+".port")), global.CONFIG.GetString(ENV+".RunningCertFile"), global.CONFIG.GetString(ENV+".RunningKeyFile"),
		); err != nil {
			zap.L().Error("项目启动失败！", zap.Error(err))
		}

	} else {
		if err = r.Run(global.CONFIG.GetString(ENV+".host") + ":" + strconv.Itoa(global.CONFIG.GetInt(ENV+".port"))); err != nil {
			zap.L().Error("项目启动失败！", zap.Error(err))
		}
	}

	// 接收退出的信号，并做处理
	q := make(chan os.Signal)
	//接收ctrl + c ，kill(排除 kill -9)
	signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)
	<-q

	//后续操作处理，比如主动从服务中心中移除当前节点
	zap.L().Error("项目停止运行了......")

}
