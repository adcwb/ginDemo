package utils

import (
	"context"
	"fmt"
	"ginDemo/global"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	client "github.com/influxdata/influxdb1-client/v2"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"reflect"
)

// QueryDBApiV1 读取influxDB数据
func QueryDBApiV1(cli client.Client, cmd, db string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: db,
	}
	if response, err := cli.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

// QueryDBApiV2 读取influxDB数据
func QueryDBApiV2(client influxdb2.Client, cmd string) (res []influxdb2.Client, err error) {
	org := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".influxOrg")
	queryAPI := client.QueryAPI(org)
	// get QueryTableResult
	result, err := queryAPI.Query(context.Background(), `from(bucket:"my-bucket")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "stat")`)
	if err == nil {
		// Iterate over query response
		for result.Next() {
			// Notice when group key has changed
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// Access data
			fmt.Printf("value: %v\n", result.Record().Value())
		}
		// check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}
	return res, err
}

// WriteDBApiV2 写入数据
func WriteDBApiV2(client influxdb2.Client, p *write.Point) {
	org := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".influxOrg")
	bucket := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".influxBucket")
	writeAPI := client.WriteAPIBlocking(org, bucket)

	//p := influxdb2.NewPoint("stat",
	//	map[string]string{"unit": "temperature"},
	//	map[string]interface{}{"avg": 24.5, "max": 45},
	//	time.Now())
	// Write point immediately
	ctx := context.Background()
	err := writeAPI.WritePoint(ctx, p)
	if err != nil {
		zap.L().Error("influxDB写入数据错误！", zap.Error(err))
	}
	err = writeAPI.Flush(ctx)
	if err != nil {
		zap.L().Error("influxDB写入数据错误！", zap.Error(err))
	}
	// Ensures background processes finishes
	//client.Close()
}

// CheckRabbitMQChannelClosed 0表示channel未关闭，1表示channel已关闭
func CheckRabbitMQChannelClosed(ch *amqp.Channel) int64 {
	d := reflect.ValueOf(ch)
	i := d.FieldByName("closed").Int()
	return i
}
