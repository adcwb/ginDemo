package initialization

import (
	"context"
	"errors"
	"fmt"
	"ginDemo/global"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis/v8" // 注意导入的是新版本
	socketio "github.com/googollee/go-socket.io"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
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

// InitConfigFile 初始化配置文件
func InitConfigFile() {
	global.CONFIG = viper.New()                  // 初始化配置文件
	global.CONFIG.SetConfigFile("./config.json") // 指定配置文件路径
	err := global.CONFIG.ReadInConfig()          // 查找并读取配置文件

	// 处理读取配置文件的错误
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			// 配置文件未找到错误；如果需要可以忽略
			zap.L().Info("配置文件未找到，请核实配置文件路径是否正确！")
		} else {
			// 配置文件被找到了，但产生了另外的错误
			panic(fmt.Errorf("读取配置文件错误: %s \n", err))
		}
	}

	err = global.CONFIG.Unmarshal(&global.ConfigMap)
	if err != nil {
		zap.L().Error("反序列化配置信息失败，请核实原因！", zap.Error(err))
	}

	// 监听配置文件
	global.CONFIG.WatchConfig()
	global.CONFIG.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		zap.L().Info("配置文件发生变化！", zap.String("ConfigName", e.Name))
		err = global.CONFIG.Unmarshal(&global.ConfigMap)
		if err != nil {
			zap.L().Error("配置文件发生变化后，反序列化配置信息失败！", zap.Error(err))
		}
	})
}

// InitMongoDBClient MongoDB初始化连接
func InitMongoDBClient() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI(global.ConfigMap.MongoConfig.Mongo).
		SetMaxPoolSize(uint64(global.ConfigMap.MongoConfig.MaxPoolSize)).
		SetMinPoolSize(uint64(global.ConfigMap.MongoConfig.MinPoolSize)).
		SetMaxConnIdleTime(time.Duration(global.ConfigMap.MongoConfig.MaxIdleTimeMS) * time.Second)

	var err error
	// 连接到MongoDB
	global.MONGO, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		zap.L().Error("连接到MongoDB出现错误！", zap.Error(err))
	}
	// 检查连接
	err = global.MONGO.Ping(ctx, nil)
	if err != nil {
		zap.L().Error("MongoDB初始化失败！", zap.Error(err))
	} else {
		zap.L().Info("MongoDB初始化成功!")
	}
}

// InitRedisClient redis初始化连接
func InitRedisClient() {
	global.REDIS = redis.NewClient(&redis.Options{
		Addr:         global.ConfigMap.RedisConfig.HOST,            //Addr:     "127.0.0.1:6379",
		Password:     global.ConfigMap.RedisConfig.PWD,             // no password set
		DB:           global.ConfigMap.RedisConfig.DB,              // use default DB
		MinIdleConns: global.ConfigMap.RedisConfig.MinIdleConnTime, // 最小连接数
		PoolSize:     global.ConfigMap.RedisConfig.PoolSize,        // 连接池大小
		MaxConnAge:   time.Duration(global.ConfigMap.RedisConfig.MaxConnAge) * time.Second,
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
func InitMySqlClient() (db *gorm.DB, err error) {
	mysqlServer := global.ConfigMap.MySqlConfig.Mysql
	DBName := global.ConfigMap.MySqlConfig.Dbname
	dsn := mysqlServer + "/" + DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("MySQL组件初始化失败！", zap.Error(err))
	}

	// 获取底层的 *sql.DB 以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Error("获取底层数据库连接失败！", zap.Error(err))
		return
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(100)                // 设置最大打开连接数
	sqlDB.SetMaxIdleConns(10)                 // 设置最大空闲连接数
	sqlDB.SetConnMaxLifetime(time.Minute * 5) // 设置连接的最大存活时间

	zap.L().Info("MySQL组件初始化成功！")
	return
}

// InitRpc 初始化RPC连接
func InitRpc() {

}

// InitInfluxDB 初始化influxDB
func InitInfluxDB() {
	// influxDB 1.x 版本使用
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     global.ConfigMap.InfluxDBConfig.InfluxDBV1.Url,
		Username: global.ConfigMap.InfluxDBConfig.InfluxDBV1.User,
		Password: global.ConfigMap.InfluxDBConfig.InfluxDBV1.Pwd,
	})
	if err != nil {
		zap.L().Error("InfluxDB组件初始化失败！", zap.Error(err))
	}
	global.InflxDBv1 = cli
}

// InitInfluxDBV2 初始化influxDBV2
func InitInfluxDBV2() {
	// influxDB 2.x 版本使用
	global.InflxDBv2 = influxdb2.NewClient(
		global.ConfigMap.InfluxDBConfig.InfluxDBV2.Url,
		global.ConfigMap.InfluxDBConfig.InfluxDBV2.Token)

	// 测试连接（可选）
	ping, err := global.InflxDBv2.Ping(context.Background())

	if !ping {
		zap.L().Error("连接到 InfluxDB v2 失败！", zap.Error(err))
		return
	}

	zap.L().Info("InfluxDB v2 客户端初始化成功！")
}

// InitRabbitMQ 初始化消息队列
func InitRabbitMQ() {
	var err error
	global.RabbitMQConn, err = amqp.Dial(global.ConfigMap.RabbitMQConfig.URL)
	if err != nil {
		zap.L().Error("RabbitMQ连接失败！", zap.Error(err))
	} else {
		global.RabbitMQChannel, err = global.RabbitMQConn.Channel()
		if err != nil {
			zap.L().Error("获取RabbitMQ channel失败！", zap.Error(err))
		} else {
			zap.L().Info("RabbitMQ初始化成功！")
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

// InitAliYunOss 初始化阿里云OSS存储
func InitAliYunOss() {
	var err error
	AccessKeyId := global.ConfigMap.AliYunOssConfig.AccessKeyId
	AccessKeySecret := global.ConfigMap.AliYunOssConfig.AccessKeySecret
	Endpoint := global.ConfigMap.AliYunOssConfig.Endpoint
	global.AliStorage, err = oss.New(Endpoint, AccessKeyId, AccessKeySecret)
	if err != nil {
		zap.L().Error("初始化阿里云OSS存储失败，请核实原因！", zap.Error(err))
	} else {
		zap.L().Info("初始化阿里云OSS存储成功！")
	}
}

// InitCasDoorSDK 初始化统一认证平台SDK
func InitCasDoorSDK() {
	AuthConfig := casdoorsdk.AuthConfig{
		Endpoint:         global.ConfigMap.SSOAuthentication.Endpoint,
		ClientId:         global.ConfigMap.SSOAuthentication.ClientID,
		ClientSecret:     global.ConfigMap.SSOAuthentication.ClientSecret,
		Certificate:      "",
		OrganizationName: "built-in",
		ApplicationName:  "golang-gin-demo",
	}
	global.CasDoorClient = casdoorsdk.NewClientWithConf(&AuthConfig)
}
