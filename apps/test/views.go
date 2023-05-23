package test

import (
	"context"
	"encoding/json"
	"fmt"
	"ginDemo/global"
	"ginDemo/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

func Pang(c *gin.Context) {
	c.JSON(http.StatusOK, "Pang")
}

func MongoTest(c *gin.Context) {
	s1 := Student{"小红", 12}
	s2 := Student{"小兰", 10}
	s3 := Student{"小黄", 11}
	s4 := Student{"小瑶", 11}

	// 指定获取要操作的数据集
	collection := global.MONGO.Database("server").Collection("student")

	// 插入一条数据
	insertResult, err := collection.InsertOne(context.TODO(), s1)
	if err != nil {
		zap.L().Error("插入数据失败", zap.Error(err))
	} else {
		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	}

	// 插入多条数据
	students := []interface{}{s2, s3, s4}
	insertManyResult, err := collection.InsertMany(context.TODO(), students)
	if err != nil {
		zap.L().Error("插入数据失败", zap.Error(err))
	} else {
		fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
	}

	// 更新文档
	/*
		D：一个BSON文档, 这种类型应该在顺序重要的情况下使用, 比如MongoDB命令.
		M：一张无序的map, 它和D是一样的, 只是它不保持顺序.
		A：一个BSON数组.
		E：D里面的一个元素.
	*/
	filter := bson.D{{"name", "小兰"}}

	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		zap.L().Error("更新数据失败", zap.Error(err))
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// 创建一个Student变量用来接收查询的结果
	var result Student
	selectData := bson.D{{"name", "小兰"}}
	err = collection.FindOne(context.TODO(), selectData).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found a single document: %+v\n", result)

	// 查询多个
	// 将选项传递给Find()
	findOptions := options.Find()
	findOptions.SetLimit(2)

	// 定义一个切片用来存储查询结果
	var results []*Student

	// 把bson.D{{}}作为一个filter来匹配所有文档
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 查找多个文档返回一个光标
	// 遍历游标允许我们一次解码一个文档
	for cur.Next(context.TODO()) {
		// 创建一个值，将单个文档解码为该值
		var elem Student
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// 完成后关闭游标
	cur.Close(context.TODO())
	fmt.Printf("Found multiple documents (array of pointers): %#v\n", results)

	// 删除名字是小黄的那个
	deleteResult1, err := collection.DeleteOne(context.TODO(), bson.D{{"name", "小黄"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult1.DeletedCount)

	// 删除所有
	deleteResult2, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult2.DeletedCount)

}

func RabbitTest(c *gin.Context) {
	//var ReturnData WechatPayCallbackReturnStruct
	////ReturnData.ReturnCode = "<![CDATA[OK]]>"
	////ReturnData.ReturnMsg = "<![CDATA[SUCCESS]]>"
	//
	//ReturnData.ReturnCode = "OK"
	//ReturnData.ReturnMsg = "SUCCESS"
	////
	////marshal := xml.Unmarshal()
	////if err != nil {
	////	fmt.Println(err)
	////	return
	////}
	////fmt.Println(string(marshal))
	//marshal, err := json.Marshal(ReturnData)
	//if err != nil {
	//	return
	//}
	//fmt.Println(string(marshal))
	//c.XML(http.StatusOK, ReturnData)
}

func QueryDB(c *gin.Context) {
	deviceNumber, _ := c.GetQuery("device_number")
	startMonth, _ := c.GetQuery("startMonth")
	endMonth, _ := c.GetQuery("endMonth")
	days, _ := c.GetQuery("days")

	cmd := ""
	db := "f_c_work"

	if deviceNumber == "" {
		returnData := map[string]interface{}{
			"code": 20001,
			"msg":  "设备号不可为空",
		}
		c.JSON(http.StatusOK, returnData)
		return
	} else {
		locTime, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			zap.L().Error("加载时区失败！", zap.Error(err))
		}

		if startMonth != "" && endMonth != "" {
			data := utils.ReturnEveryDays(startMonth, endMonth, false)

			timeObj1, err := time.ParseInLocation("2006-01-02 15:04:05", data[0]+" 00:00:00", locTime)
			if err != nil {
				zap.L().Error("解析时间字符串失败！", zap.Error(err))
			}
			timeObj2, err2 := time.ParseInLocation("2006-01-02 15:04:05", data[len(data)-1]+" 23:59:59", locTime)
			if err2 != nil {
				zap.L().Error("解析时间字符串失败！", zap.Error(err))
			}
			start := timeObj1.Format(time.RFC3339)
			end := timeObj2.Format(time.RFC3339)
			cmd = fmt.Sprintf("SELECT * FROM %s where \"DeviceNumber\" = '%s' and time >= '%s'  and time <= '%s'", "dev_local_heart", deviceNumber, start, end)
		} else if days != "" {
			// 加载时区
			timeObj1, err := time.ParseInLocation("2006-01-02 15:04:05", days+" 00:00:00", locTime)
			if err != nil {
				zap.L().Error("解析时间字符串失败！", zap.Error(err))
			}
			timeObj2, err2 := time.ParseInLocation("2006-01-02 15:04:05", days+" 23:59:59", locTime)
			if err2 != nil {
				zap.L().Error("解析时间字符串失败！", zap.Error(err))
			}
			//start := timeObj1.Unix()
			//end := timeObj2.Unix()
			start := timeObj1.Format(time.RFC3339)
			end := timeObj2.Format(time.RFC3339)
			cmd = fmt.Sprintf("SELECT * FROM %s where \"DeviceNumber\" = '%s' and time >= '%s'  and time <= '%s'", "dev_local_heart", deviceNumber, start, end)
		} else {
			returnData := map[string]interface{}{
				"code": 20001,
				"msg":  "日期或天数必须传递一个，不可同时为空",
			}
			c.JSON(http.StatusOK, returnData)
			return
		}
		fmt.Println("influxDB查询语句：", cmd)
		res, err := utils.QueryDBApiV1(global.InflxDBv1, cmd, db)
		if err != nil {
			zap.L().Error("数据库查询失败！", zap.Error(err))
			c.JSON(http.StatusOK, "数据库查询失败")
			return
		}
		//fmt.Printf("%v --- %T", res, res)
		// 开始解析查询到的数据，注意此处返回的数据类型是json.Number，需要手动断言才可用
		a := 0.0

		if len(res[0].Series) == 0 {
			ReturnData := map[string]interface{}{
				"code":  20000,
				"data":  "暂未查询到数据！",
				"count": a,
			}
			c.JSON(http.StatusOK, ReturnData)
			return
		}
		tempData := make([]map[string]interface{}, 0, 8)
		for _, row := range res[0].Series[0].Values {
			dataTime := ""
			dataTotal := float64(0)

			for k, v := range row {
				if k == 0 {
					value, ok := v.(string)
					if ok {
						timeObj1, err3 := time.ParseInLocation(time.RFC3339, value, locTime)
						if err3 != nil {
							zap.L().Error("解析时间字符串失败！", zap.Error(err))
						}
						dataTime = timeObj1.Format("2006-01-02 15:04:05")
					}
				}
				if k == 4 {
					value, ok := v.(json.Number)
					if ok {
						dataTotal, _ = value.Float64()
						a = a + dataTotal
					}
				}
			}

			temp1 := map[string]interface{}{
				"dataTime": dataTime,
				"total":    dataTotal,
			}
			tempData = append(tempData, temp1)
		}
		ReturnData := map[string]interface{}{
			"code":  20000,
			"data":  tempData,
			"count": a,
		}
		c.JSON(http.StatusOK, ReturnData)
	}
}

type SendMailStruct struct {
	Mail    string `json:"mail"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func SendMail(c *gin.Context) {
	b, _ := c.GetRawData()
	var tempData SendMailStruct
	err := json.Unmarshal(b, &tempData)
	if err != nil {
		zap.L().Error("Json序列化失败，请核对！", zap.Error(err))
	}

	utils.SendEmail(tempData.Mail, tempData.Subject, tempData.Message)
	c.JSON(http.StatusOK, "OK")
}
