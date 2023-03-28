package global

import (
	"context"
	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis/v8" // 注意导入的是新版本
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"sync"
)

var (
	DB              *gorm.DB
	JobS            *gocron.Scheduler
	REDIS           *redis.Client
	REDISCTX        context.Context
	CONFIG          *viper.Viper
	MONGO           *mongo.Client
	RabbitMQConn    *amqp.Connection
	RabbitMQChannel *amqp.Channel
	InflxDBv2       influxdb2.Client
	InflxDBv1       client.Client
	lock            sync.RWMutex // 全局声明一把读写锁
)
