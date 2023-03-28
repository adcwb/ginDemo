package dingtalk

import (
	"context"
	"encoding/json"
	"ginDemo/global"
	"ginDemo/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

func Test(c *gin.Context) {
	method := c.Request.Method
	if method == "GET" {
		//
		data, err := utils.HttpClient(global.CONFIG.GetString("UserServer")+"/users/", "GET", "", c.GetString("X-Token"))
		if err != nil {
			zap.L().Error("HttpClient请求发送失败", zap.Error(err))
		}
		// 将返回的数据进行反序列化
		var m testStruct
		err = json.Unmarshal(data, &m)
		if err != nil {
			zap.L().Error("json序列化失败", zap.Error(err))
		}

		c.JSON(http.StatusOK, m)

	} else if method == "POST" {
		c.String(http.StatusNotFound, "404 page not found!")
	} else {
		c.String(http.StatusNotFound, "404 page not found!")
	}
}

// GetToken 获取钉钉开放平台-企业内部应用access_token
func GetToken(c *gin.Context) {
	// 先从Redis数据库中查询看是否有token，若有token且token并未过期则直接获取数据返回，否则重新请求

	ctx := context.Background()
	result, err := global.REDIS.Get(ctx, "DING_TALK_TOKEN_KEY").Result()
	if err == redis.Nil {
		zap.L().Error("RedisKey DING_TALK_TOKEN_KEY does not exist", zap.Error(err))
		urls := global.CONFIG.GetString("dingTalkServer") + "/gettoken" + "?appkey=???"
		data, err := utils.HttpClient(urls, "GET", "", "")
		if err != nil {
			zap.L().Error("HttpClient请求发送失败", zap.Error(err))
		}

		var ReturnData GetTokenStruct
		err = json.Unmarshal(data, &ReturnData)
		if err != nil {
			zap.L().Error("GetTokenStruct序列化失败", zap.Error(err))
		}

		err = global.REDIS.Set(ctx, "DING_TALK_TOKEN_KEY", ReturnData.AccessToken, 7100*time.Second).Err()
		if err != nil {
			zap.L().Error("RedisKey Set DING_TALK_TOKEN_KEY Error", zap.Error(err))
		}
		zap.L().Info("Redis Set DING_TALK_TOKEN_KEY Success!")
		c.JSON(http.StatusOK, ReturnData)

	} else if err != nil {
		zap.L().Error("RedisKey Get DING_TALK_TOKEN_KEY Error", zap.Error(err))
	}

	if result != "" {
		ReturnData := map[string]interface{}{
			"errcode":      0,
			"access_token": result,
			"errmsg":       "ok",
			"expires_in":   3600,
		}
		c.JSON(http.StatusOK, ReturnData)
	}
}

// GetAllForms 获取当前企业所有可管理的表单
func GetAllForms(c *gin.Context) {
	urlData := c.Request.URL.Query().Encode()
	urls := global.CONFIG.GetString("dingTalkServer") + "/topapi/process/template/manage/get?" + urlData
	data, err := utils.HttpClient(urls, "POST", "{\"userid\": \"\"}", "")
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}

	var ReturnData GetAllFormsStruct
	err = json.Unmarshal(data, &ReturnData)
	if err != nil {
		zap.L().Error("GetTokenStruct序列化失败", zap.Error(err))
	}
	c.JSON(http.StatusOK, ReturnData)
}

// GetFromProcessCode 获取表单schema
func GetFromProcessCode(c *gin.Context) {
	urlData := c.Request.URL.Query().Encode()
	url := "https://api.dingtalk.com/v1.0/workflow/forms/schemas/processCodes?" + urlData

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}
	ctx := context.Background()
	result, err := global.REDIS.Get(ctx, "DING_TALK_TOKEN_KEY").Result()
	if err == redis.Nil {
		zap.L().Error("RedisKey DING_TALK_TOKEN_KEY does not exist", zap.Error(err))
		urls := global.CONFIG.GetString("dingTalkServer") + "/gettoken" + "?appkey=???"
		data, err := utils.HttpClient(urls, "GET", "", "")
		if err != nil {
			zap.L().Error("HttpClient请求发送失败", zap.Error(err))
		}

		var ReturnData GetTokenStruct
		err = json.Unmarshal(data, &ReturnData)
		if err != nil {
			zap.L().Error("GetTokenStruct序列化失败", zap.Error(err))
		}

		err = global.REDIS.Set(ctx, "DING_TALK_TOKEN_KEY", ReturnData.AccessToken, 7100*time.Second).Err()
		if err != nil {
			zap.L().Error("RedisKey Set DING_TALK_TOKEN_KEY Error", zap.Error(err))
		}
		zap.L().Info("Redis Set DING_TALK_TOKEN_KEY Success!")
		result = ReturnData.AccessToken
	}

	req.Header.Add("x-acs-dingtalk-access-token", result)

	res, err := client.Do(req)
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}

	var ReturnData GetFromProcessCodeStruct
	err = json.Unmarshal(body, &ReturnData)
	if err != nil {
		zap.L().Error("GetTokenStruct序列化失败", zap.Error(err))
	}
	c.JSON(http.StatusOK, ReturnData)

}

func GetUUID5(c *gin.Context) {
	data := utils.UUID5("test")
	ReturnData := map[string]string{
		"UUID5": data,
	}
	c.JSON(http.StatusOK, ReturnData)
}

func DingToken() (result string, err error) {
	ctx := context.Background()
	result, err = global.REDIS.Get(ctx, "DING_TALK_TOKEN_KEY").Result()
	if err == redis.Nil {
		zap.L().Error("RedisKey DING_TALK_TOKEN_KEY does not exist", zap.Error(err))
		urls := global.CONFIG.GetString("dingTalkServer") + "/gettoken" + "?appkey=???"
		data, err := utils.HttpClient(urls, "GET", "", "")
		if err != nil {
			zap.L().Error("HttpClient请求发送失败", zap.Error(err))
		}

		var ReturnData GetTokenStruct
		err = json.Unmarshal(data, &ReturnData)
		if err != nil {
			zap.L().Error("GetTokenStruct序列化失败", zap.Error(err))
		}

		err = global.REDIS.Set(ctx, "DING_TALK_TOKEN_KEY", ReturnData.AccessToken, 7100*time.Second).Err()
		if err != nil {
			zap.L().Error("RedisKey Set DING_TALK_TOKEN_KEY Error", zap.Error(err))
		}
		zap.L().Info("Redis Set DING_TALK_TOKEN_KEY Success!")
		result = ReturnData.AccessToken
		return result, err
	} else {
		return result, err
	}
}

// 获取部门列表
func test(c *gin.Context) {
	token, err := DingToken()
	if err != nil {
		zap.L().Error("RedisKey Get DING_TALK_TOKEN_KEY Error", zap.Error(err))
	}
	c.JSON(http.StatusOK, token)
}
