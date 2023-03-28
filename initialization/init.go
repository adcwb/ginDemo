package initialization

import (
	"context"
	"fmt"
	"ginDemo/global"
	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis/v8" // 注意导入的是新版本
	socketio "github.com/googollee/go-socket.io"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func InitStorage() {

}

// InitConfigFile 初始化配置文件
func InitConfigFile() {
	// 初始化配置文件
	global.CONFIG = viper.New()
	global.CONFIG.SetConfigFile("./config.json")
	err := global.CONFIG.ReadInConfig() // 查找并读取配置文件
	if err != nil {                     // 处理读取配置文件的错误
		panic(fmt.Errorf("读取配置文件错误: %s \n", err))
	}
	global.CONFIG.WatchConfig()
	global.CONFIG.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		zap.L().Info("配置文件发生变化！", zap.String("ConfigName", e.Name))
	})
}

// InitMongoDBClient MongoDB初始化连接
func InitMongoDBClient(env string) {
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI(global.CONFIG.GetString(env + ".mongo"))
	var err error
	// 连接到MongoDB
	global.MONGO, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		zap.L().Error("连接到MongoDB出现错误！", zap.Error(err))
	}
	// 检查连接
	err = global.MONGO.Ping(context.TODO(), nil)
	if err != nil {
		zap.L().Error("MongoDB初始化失败！", zap.Error(err))
	} else {
		zap.L().Info("MongoDB初始化成功!")
	}
}

// InitRedisClient redis初始化连接
func InitRedisClient(env string) {
	global.REDIS = redis.NewClient(&redis.Options{
		Addr: global.CONFIG.GetString(env + ".redis"),
		//Addr:     "127.0.0.1:6379",
		Password: "",  // no password set
		DB:       0,   // use default DB
		PoolSize: 100, // 连接池大小
	})

	var cancel context.CancelFunc
	global.REDISCTX, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := global.REDIS.Ping(global.REDISCTX).Result()
	if err != nil {
		zap.L().Error("Redis初始化失败！", zap.Error(err))
	}

	zap.L().Info("Redis初始化成功")
}

// InitMySqlClient 初始化数据库连接
func InitMySqlClient(env string) (db *gorm.DB, err error) {
	mysqlServer := global.CONFIG.GetString(env + ".mysql")
	DBName := global.CONFIG.GetString(env + ".dbname")
	dsn := mysqlServer + "/" + DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("MySQL组件初始化失败！", zap.Error(err))
	} else {
		zap.L().Info("MySQL组件初始化成功！")
	}
	return
}

// InitRpc 初始化RPC连接
func InitRpc() {

}

// InitInfluxDB 初始化influxDB
func InitInfluxDB(env string) {
	// influxDB 2.x 版本使用
	//global.InflxDB = influxdb2.NewClient(global.CONFIG.GetString(env+".influxDBv2"), global.CONFIG.GetString(env+".influxDBAuthToken"))

	// influxDB 1.x 版本使用
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     global.CONFIG.GetString(env + ".influxDBv1"),
		Username: "admin",
		Password: "",
	})
	if err != nil {
		zap.L().Error("InfluxDB组件初始化失败！", zap.Error(err))
	}
	global.InflxDBv1 = cli
}

// InitRabbitMQ 初始化消息队列
func InitRabbitMQ(env string) {
	var err error
	global.RabbitMQConn, err = amqp.Dial(global.CONFIG.GetString(env + ".RabbitMQUrl"))
	if err != nil {
		zap.L().Error("RabbitMQ连接失败！", zap.Error(err))
	} else {
		global.RabbitMQChannel, err = global.RabbitMQConn.Channel()
		if err != nil {
			zap.L().Error("获取RabbitMQ channel失败！", zap.Error(err))
		}
	}
}

// InitSocketIO 初始化WebSocket
func InitSocketIO() {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
}
