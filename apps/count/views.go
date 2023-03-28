package count

import (
	"context"
	"encoding/json"
	"fmt"
	"ginDemo/global"
	"ginDemo/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetQueryDB 获取历史流量数据
func GetQueryDB(c *gin.Context) {
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

// ExpressDelivery 快递实时信息查询
func ExpressDelivery(c *gin.Context) {
	BackExpressOperator, _ := c.GetQuery("balk_express_operator")
	AfterSalesID, _ := c.GetQuery("after_sale_id")
	K100Key := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".K100Key")
	K100Customer := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".K100Customer")

	// 加缓存
	ctx := context.Background()
	redisKey := "K100_" + BackExpressOperator + "_" + AfterSalesID
	result, err := global.REDIS.Get(ctx, redisKey).Result()

	if err == redis.Nil {
		zap.L().Error(redisKey+"is null", zap.Error(err))
	} else if err != nil {
		zap.L().Error("RedisKey USER_TOKEN_KEY does not exist", zap.Error(err))
	}
	if result != "" {
		var ReturnData ExpressDeliveryStruct
		err = json.Unmarshal([]byte(result), &ReturnData)
		if err != nil {
			zap.L().Error("Json反序列化失败", zap.Error(err))
		}

		c.JSON(http.StatusOK, ReturnData)
		return
	}

	tempData1 := map[string]string{
		"com":      BackExpressOperator,
		"num":      AfterSalesID,
		"resultv2": "4",
		"show":     "0",
	}

	marshal, err := json.Marshal(tempData1)
	if err != nil {
		zap.L().Error("json序列化失败！", zap.Error(err))
	}

	sign := string(marshal) + K100Key + K100Customer
	Md5 := utils.StringToMD5(sign)

	url := "https://poll.kuaidi100.com/poll/query.do"
	method := "POST"

	payload := strings.NewReader("customer=" + K100Customer + "&sign=" + Md5 + "&param=" + string(marshal))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))

	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			zap.L().Error("关闭连接失败", zap.Error(err))
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("读取body数据错误", zap.Error(err))
	}
	var ReturnData ExpressDeliveryStruct
	err = json.Unmarshal(body, &ReturnData)
	if err != nil {
		zap.L().Error("Json反序列化失败", zap.Error(err))
	}

	global.REDIS.Set(ctx, redisKey, string(body), 1*time.Hour)
	c.JSON(http.StatusOK, ReturnData)
}

// ExpressDeliveryMap 快递实时信息查询带地图版
func ExpressDeliveryMap(c *gin.Context) {
	BackExpressOperator, _ := c.GetQuery("balk_express_operator")
	AfterSalesID, _ := c.GetQuery("after_sale_id")
	FromCity, _ := c.GetQuery("from_city")
	ToCity, _ := c.GetQuery("to_city")
	PhoneNumber := c.DefaultQuery("phone_number", "")

	// 加缓存
	ctx := context.Background()
	redisKey := "K100_Map" + BackExpressOperator + "_" + AfterSalesID
	result, err := global.REDIS.Get(ctx, redisKey).Result()

	if err == redis.Nil {
		zap.L().Error(redisKey+"is null", zap.Error(err))
	} else if err != nil {
		zap.L().Error("RedisKey K100_Map does not exist", zap.Error(err))
	}
	if result != "" {
		var ReturnData ExpressDeliveryMapStruct
		err = json.Unmarshal([]byte(result), &ReturnData)
		if err != nil {
			zap.L().Error("Json反序列化失败", zap.Error(err))
		}
		if ReturnData.Message != "ok" {
			err = global.REDIS.Del(ctx, redisKey).Err()
			if err != nil {
				zap.L().Error("RedisKey K100_Map DELETE Error", zap.Error(err))
			}
		}

		c.JSON(http.StatusOK, ReturnData)
		return
	}

	tempData1 := map[string]string{}

	// 顺丰快递需要额外提供手机号参数
	if BackExpressOperator == "shunfeng" || BackExpressOperator == "shunfengkuaiyun" || BackExpressOperator == "shunfenglengyun" {
		if PhoneNumber == "" {
			var ReturnData ExpressDeliveryMapStruct
			ReturnData.Message = "顺丰快递查询必须传递手机号，本次查询失败！"
			c.JSON(http.StatusOK, ReturnData)
			return
		}
		if utils.IsMobile(PhoneNumber) == false {
			var ReturnData ExpressDeliveryMapStruct
			ReturnData.Message = "手机号格式校验失败，请重试！"
			c.JSON(http.StatusOK, ReturnData)
			return
		}

		tempData1 = map[string]string{
			"com":      BackExpressOperator,
			"num":      AfterSalesID,
			"from":     FromCity,
			"to":       ToCity,
			"resultv2": "5",
			"show":     "0",
			"phone":    PhoneNumber,
		}
	} else {
		tempData1 = map[string]string{
			"com":      BackExpressOperator,
			"num":      AfterSalesID,
			"from":     FromCity,
			"to":       ToCity,
			"resultv2": "5",
			"show":     "0",
		}
	}

	marshal, err := json.Marshal(tempData1)
	if err != nil {
		zap.L().Error("json序列化失败！", zap.Error(err))
	}
	K100Key := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".K100Key")
	K100Customer := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".K100Customer")

	sign := string(marshal) + K100Key + K100Customer
	Md5 := utils.StringToMD5(sign)

	url := "https://poll.kuaidi100.com/poll/maptrack.do"
	method := "POST"

	payload := strings.NewReader("customer=" + K100Customer + "&sign=" + Md5 + "&param=" + string(marshal))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))

	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			zap.L().Error("关闭连接失败", zap.Error(err))
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("读取body数据错误", zap.Error(err))
	}
	var ReturnData ExpressDeliveryMapStruct
	err = json.Unmarshal(body, &ReturnData)
	if err != nil {
		zap.L().Error("Json反序列化失败", zap.Error(err))
	}

	global.REDIS.Set(ctx, redisKey, string(body), 1*time.Hour)
	c.JSON(http.StatusOK, ReturnData)
}

// GetAutonumber 智能单号识别
func GetAutonumber(c *gin.Context) {
	AfterSalesID, _ := c.GetQuery("after_sale_id")
	K100Key := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".K100Key")
	url := "http://www.kuaidi100.com/autonumber/auto?num=" + AfterSalesID + "&key=" + K100Key
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		zap.L().Error("HttpClient初始化失败", zap.Error(err))
	}

	res, err := client.Do(req)
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			zap.L().Error("关闭连接失败", zap.Error(err))
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("读取body数据错误", zap.Error(err))
	}

	var ReturnData []GetAutonumberStruct
	err = json.Unmarshal(body, &ReturnData)
	if err != nil {
		zap.L().Error("Json反序列化失败", zap.Error(err))
	}

	c.JSON(http.StatusOK, ReturnData)
}

// ExpressDeliveryPool 快递信息订阅，不带地图版
func ExpressDeliveryPool(c *gin.Context) {
	b, _ := c.GetRawData()
	var tempData ExpressDeliveryPoolStruct

	err := json.Unmarshal(b, &tempData)
	if err != nil {
		zap.L().Error("Json序列化失败！", zap.Error(err))
	}

	// 配置文件中获取key
	K100Key := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".K100Key")
	K100Customer := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".K100Customer")
	K100URLCallback := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".K100URLCallback")
	zap.L().Info("快递信息订阅接口配置文件参数获取", zap.String("K100Key", K100Key), zap.String("K100Customer", K100Customer), zap.String("K100URLCallback", K100URLCallback))

	//构造请求参数
	k100parameters := map[string]string{
		"callbackurl": K100URLCallback,
		"salt":        "tianchao",
		"resultv2":    "4",
	}

	k100param := map[string]interface{}{
		"company":    tempData.BalkExpressOperator,
		"number":     tempData.AfterSaleId,
		"key":        K100Key,
		"parameters": k100parameters,
		"from":       "",
		"to":         "",
	}

	marshal, err := json.Marshal(k100param)
	if err != nil {
		zap.L().Error("json序列化失败！", zap.Error(err))
	}

	RequestData, err := utils.HttpClient("https://poll.kuaidi100.com/poll?schema=json&param="+string(marshal), "POST", "", "form")
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}

	var ReturnData GetAutonumberStruct
	err = json.Unmarshal(RequestData, &ReturnData)
	if err != nil {
		zap.L().Error("返回值序列化失败！", zap.Error(err))
	}

	zap.L().Info("快递消息订阅成功！", zap.String("快递单号", tempData.AfterSaleId), zap.String("快递公司", tempData.BalkExpressOperator))
	c.JSON(http.StatusOK, ReturnData)
}

// ExpressDeliveryPoolMap 快递信息订阅，带地图版
func ExpressDeliveryPoolMap(c *gin.Context) {
	b, _ := c.GetRawData()
	var tempData ExpressDeliveryPoolStruct

	err := json.Unmarshal(b, &tempData)
	if err != nil {
		zap.L().Error("Json序列化失败！", zap.Error(err))
	}

	// 配置文件中获取key
	K100Key := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".K100Key")
	K100Customer := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".K100Customer")
	K100URLCallback := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".K100URLCallback")
	zap.L().Info("快递信息订阅接口配置文件参数获取", zap.String("K100Key", K100Key), zap.String("K100Customer", K100Customer), zap.String("K100URLCallback", K100URLCallback))

	//构造请求参数
	// 顺丰快递需要额外提供手机号参数
	k100parameters := map[string]string{}
	if tempData.BalkExpressOperator == "shunfeng" || tempData.BalkExpressOperator == "shunfengkuaiyun" || tempData.BalkExpressOperator == "shunfenglengyun" {
		if tempData.PhoneNumber == "" {
			var ReturnData ExpressDeliveryMapStruct
			ReturnData.Message = "顺丰快递查询必须传递手机号，本次查询失败！"
			c.JSON(http.StatusOK, ReturnData)
			return
		}
		if utils.IsMobile(tempData.PhoneNumber) == false {
			var ReturnData ExpressDeliveryMapStruct
			ReturnData.Message = "手机号格式校验失败，请重试！"
			c.JSON(http.StatusOK, ReturnData)
			return
		}

		k100parameters = map[string]string{
			"callbackurl": K100URLCallback,
			"salt":        "tianchao",
			"resultv2":    "4",
			"phone":       tempData.PhoneNumber,
		}
	} else {
		k100parameters = map[string]string{
			"callbackurl": K100URLCallback,
			"salt":        "tianchao",
			"resultv2":    "5",
		}
	}

	k100param := map[string]interface{}{
		"company":    tempData.BalkExpressOperator,
		"number":     tempData.AfterSaleId,
		"key":        K100Key,
		"parameters": k100parameters,
		"from":       tempData.FromCity,
		"to":         tempData.ToCity,
	}

	marshal, err := json.Marshal(k100param)
	if err != nil {
		zap.L().Error("json序列化失败！", zap.Error(err))
	}

	RequestData, err := utils.HttpClient("https://poll.kuaidi100.com/pollmap?schema=json&param="+string(marshal), "POST", "", "form")
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}

	var ReturnData GetAutonumberStruct
	err = json.Unmarshal(RequestData, &ReturnData)
	if err != nil {
		zap.L().Error("返回值序列化失败！", zap.Error(err))
	}

	zap.L().Info("快递消息订阅成功！", zap.String("快递单号", tempData.AfterSaleId), zap.String("快递公司", tempData.BalkExpressOperator))
	c.JSON(http.StatusOK, ReturnData)
}

// GetExpressDeliveryPoolMap 快递订阅消息查询接口
func GetExpressDeliveryPoolMap(c *gin.Context) {
	afterSaleId, _ := c.GetQuery("after_sale_id")
	balkExpressOperator, _ := c.GetQuery("balk_express_operator")
	// 创建一个Student变量用来接收查询的结果
	var result MongoGetK100
	selectData := bson.D{{"courier_number", afterSaleId}}
	// 连接到数据库k100Data, 快递公司名称为表
	collection := global.MONGO.Database("k100Data").Collection(balkExpressOperator)

	err := collection.FindOne(context.TODO(), selectData).Decode(&result)
	if err != nil {
		zap.L().Error("MongoDB查询数据失败！", zap.Error(err))
	}

	c.JSON(http.StatusOK, result.Data)
}

// ExpressDeliveryCallback 快递订阅查询回调处理函数
func ExpressDeliveryCallback(c *gin.Context) {
	Sign := c.PostForm("sign")
	// 此处不对数据完整性校验，若以后需要校验md5(param+salt) salt为秘钥，值为tianchao MD5一定要转大写
	Param := c.PostForm("param")
	var tempData ExpressDeliveryCallbackStruct
	err := json.Unmarshal([]byte(Param), &tempData)
	if err != nil {
		zap.L().Error("Json序列化失败！", zap.Error(err))
	}

	zap.L().Info("快递推送消息_PostForm_sign", zap.String("data", Sign))
	zap.L().Info("快递推送消息_PostForm_param", zap.String("data", Param))

	// 将数据保存到MongoDB中，方便以后查询
	tempMongoData := map[string]interface{}{
		"courier_number": tempData.LastResult.Nu,
		"data":           tempData,
	}
	// 连接到数据库k100Data, 快递公司名称为表
	collection := global.MONGO.Database("k100Data").Collection(tempData.LastResult.Com)

	// 创建一个Student变量用来接收查询的结果
	var result MongoGetK100
	selectData := bson.D{{"courier_number", tempData.LastResult.Nu}}
	err = collection.FindOne(context.TODO(), selectData).Decode(&result)
	if err != nil {
		zap.L().Error("MongoDB查询数据失败！", zap.Error(err))
	}

	if result.CourierNumber != "" {
		// 若查询到了直接删除
		deleteResult1, err := collection.DeleteOne(context.TODO(), bson.D{{"courier_number", tempData.LastResult.Nu}})
		if err != nil {
			zap.L().Error("删除原数据失败", zap.Error(err))
		}
		zap.L().Info("Deleted documents in the trainers collection", zap.String("data", strconv.FormatInt(deleteResult1.DeletedCount, 10)))
	}
	// 插入一条数据
	insertResult, err := collection.InsertOne(context.TODO(), tempMongoData)
	if err != nil {
		zap.L().Error("MongoDB插入数据失败", zap.Error(err))
	}
	zap.L().Error("MongoDB插入数据成功", zap.Any("data", insertResult.InsertedID))
	ReturnData := map[string]interface{}{
		"result":     true,
		"returnCode": "200",
		"message":    "成功",
		"data":       "",
	}

	c.JSON(http.StatusOK, ReturnData)

}
