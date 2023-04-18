package test

import (
	"encoding/json"
	"fmt"
	"ginDemo/global"
	"ginDemo/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func test() {
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println(">>>>>>>>>>>>>>>>Hello Job !!>>>>>>>>>>>>>>>>")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
}
func JobTest(c *gin.Context) {
	global.JobS.Every(1).Seconds().Do(test)
}

func JobStop(c *gin.Context) {
	global.JobS.Stop()
}

func Operator(c *gin.Context) {
	params := make(map[string]string)
	err := c.ShouldBindQuery(params)
	if err != nil {
		zap.L().Error("ShouldBindQuery Error", zap.Error(err))
	}
	b, _ := c.GetRawData()
	ycCTCC := utils.OperatorSign(utils.YanChengCTCC{})
	timestamp := time.Now().Format("20060102150405")

	sign := ycCTCC.Sign(utils.SortMapToURLParams(params), string(b), timestamp)
	tempKeys := map[string]string{
		"Sign":      sign,
		"Timestamp": timestamp,
		"AppKey":    global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".YanChengChinaTelecomAppKey"),
	}

	marshal, err := json.Marshal(tempKeys)
	if err != nil {
		zap.L().Error("Json序列化失败", zap.Error(err))
	}
	ReturnData, err := utils.HttpClient("https://cmp-api.ctwing.cn:20164/openapi/v1/prodinst/realNameQueryIot?"+utils.SortMapToURLParams(params), "GET", "", string(marshal))
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}
	toJson, err := utils.XmlToJson(string(ReturnData))
	if err != nil {
		zap.L().Error("XmlToJson转化失败", zap.Error(err))
	}

	var tempData T
	err = json.Unmarshal([]byte(toJson), &tempData)
	if err != nil {
		zap.L().Error("Json序列化失败", zap.Error(err))
	}

	c.JSON(http.StatusOK, tempData)
}

// ThreeCodeMutualCheck 电信三号互查接口
func ThreeCodeMutualCheck(c *gin.Context) {
	params := make(map[string]string)
	err := c.ShouldBindQuery(params)
	if err != nil {
		zap.L().Error("ShouldBindQuery Error", zap.Error(err))
	}
	b, _ := c.GetRawData()

	ycCTCC := utils.OperatorSign(utils.YanChengCTCC{})
	timestamp := time.Now().Format("20060102150405")
	sign := ycCTCC.Sign(utils.SortMapToURLParams(params), string(b), timestamp)
	tempKeys := map[string]string{
		"Sign":      sign,
		"Timestamp": timestamp,
		"AppKey":    global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".YanChengChinaTelecomAppKey"),
	}

	marshal, err := json.Marshal(tempKeys)
	if err != nil {
		zap.L().Error("Json序列化失败", zap.Error(err))
	}
	ReturnData, err := utils.HttpClient("https://cmp-api.ctwing.cn:20164/openapi/v1/prodinst/getSIMList?"+utils.SortMapToURLParams(params), "GET", "", string(marshal))
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}

	var tempData T3
	err = json.Unmarshal(ReturnData, &tempData)
	if err != nil {
		zap.L().Error("Json序列化失败", zap.Error(err))
	}

	c.JSON(http.StatusOK, tempData)
}
