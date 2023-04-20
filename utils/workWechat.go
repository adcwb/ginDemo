package utils

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"ginDemo/global"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
)

// GetAccessTokenReturnStruct 企业微信AccessToken
type GetAccessTokenReturnStruct struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// GetTicketReturnStruct 获取企业的jsapi_ticket
type GetTicketReturnStruct struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

type MsgContent struct {
	ToUsername   string `xml:"ToUserName"`
	FromUsername string `xml:"FromUserName"`
	CreateTime   uint32 `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	Msgid        string `xml:"MsgId"`
	Agentid      uint32 `xml:"AgentId"`
}

// GetWechatToken 获取企业微信Token
func GetWechatToken() (result string) {
	ctx := context.Background()

	result, err := global.REDIS.Get(ctx, "WeChatAccessToken").Result()

	if err == redis.Nil {
		zap.L().Error("RedisKey PhoneNumbers does not exist", zap.Error(err))

	} else if err != nil {
		zap.L().Error("RedisKey PhoneNumbers does not exist", zap.Error(err))
	}
	// 若有缓存直接返回
	if len(result) > 5 && result != "" {
		return result
	} else {
		CorpID := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".WorkWechatCorpID")
		CorpSecret := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".WorkWechatCorpSecret")

		data, err := HttpClient(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".WorkWechatUrl")+"/cgi-bin/gettoken?corpid="+CorpID+"&corpsecret="+CorpSecret, "GET", "", "")
		if err != nil {
			zap.L().Error("HttpClient请求发送失败", zap.Error(err))
		}
		var ReturnData GetAccessTokenReturnStruct
		err = json.Unmarshal(data, &ReturnData)
		if err != nil {
			zap.L().Error("WeChatGetAccessToken接口返回数据序列化失败！", zap.Error(err))
		}
		if ReturnData.ErrCode == 0 {
			// 将数据记录到Redis数据库一份
			seconds := 7150
			err = global.REDIS.Set(ctx, "WeChatAccessToken", ReturnData.AccessToken, time.Duration(seconds)*time.Second).Err()

			if err != nil {
				zap.L().Error("Redis Set Key Error", zap.String("keys", "WeChatAccessToken"), zap.String("value", ReturnData.AccessToken), zap.Error(err))
			}
		}
		return ReturnData.AccessToken
	}

}

// GetWechatAgentToken 获取企业微信应用Token
func GetWechatAgentToken() (result string) {
	ctx := context.Background()

	result, err := global.REDIS.Get(ctx, "WeChatAgentAccessToken").Result()

	if err == redis.Nil {
		zap.L().Error("RedisKey WeChatAgentAccessToken does not exist", zap.Error(err))

	} else if err != nil {
		zap.L().Error("RedisKey WeChatAgentAccessToken does not exist", zap.Error(err))
	}
	// 若有缓存直接返回
	if len(result) > 5 && result != "" {
		fmt.Println("GetWechatAgentToken接口有缓存，直接返回数据：" + result)
		return result
	} else {
		CorpID := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".WorkWechatCorpID")
		CorpSecret := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".WorkWechatAgentCorpSecret")
		fmt.Println("GetWechatAgentToken接口无缓存，CorpID：" + CorpID)
		fmt.Println("GetWechatAgentToken接口无缓存，CorpSecret：" + CorpSecret)
		fmt.Println("GetWechatAgentToken接口无缓存，URLS：" + global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".WorkWechatUrl") + "/cgi-bin/gettoken?corpid=" + CorpID + "&corpsecret=" + CorpSecret)
		data, err := HttpClient(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".WorkWechatUrl")+"/cgi-bin/gettoken?corpid="+CorpID+"&corpsecret="+CorpSecret, "GET", "", "")
		if err != nil {
			zap.L().Error("HttpClient请求发送失败", zap.Error(err))
		}
		var ReturnData GetAccessTokenReturnStruct
		fmt.Println("GetWechatAgentToken接口无缓存，ReturnData：" + string(data))
		err = json.Unmarshal(data, &ReturnData)
		if err != nil {
			zap.L().Error("WeChatGetAccessToken接口返回数据序列化失败！", zap.Error(err))
		}
		if ReturnData.ErrCode == 0 {
			// 将数据记录到Redis数据库一份
			seconds := 7150
			err = global.REDIS.Set(ctx, "WeChatAgentAccessToken", ReturnData.AccessToken, time.Duration(seconds)*time.Second).Err()

			if err != nil {
				zap.L().Error("Redis Set Key Error", zap.String("keys", "WeChatAgentAccessToken"), zap.String("value", ReturnData.AccessToken), zap.Error(err))
			}
		}
		return ReturnData.AccessToken
	}

}

// GetWorkJsAPITicket 获取企业的jsapi_ticket
func GetWorkJsAPITicket() (JsAPITicket string) {
	token := GetWechatAgentToken()
	ctx := context.Background()
	result, err := global.REDIS.Get(ctx, "WeChatWorkJsAPITicketToken").Result()

	if err == redis.Nil {
		zap.L().Error("RedisKey WeChatWorkJsAPITicketToken does not exist", zap.Error(err))

	} else if err != nil {
		zap.L().Error("RedisKey WeChatWorkJsAPITicketToken does not exist", zap.Error(err))
	}

	// 若有缓存直接返回
	if len(result) > 5 && result != "" {
		return result
	}

	data, err := HttpClient(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".WorkWechatUrl")+"/cgi-bin/get_jsapi_ticket?access_token="+token, "GET", "", "")
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}

	var ReturnData GetTicketReturnStruct
	err = json.Unmarshal(data, &ReturnData)
	if err != nil {
		zap.L().Error("WeChatGetTicketReturnStruct接口返回数据序列化失败！", zap.Error(err))
	}
	if ReturnData.ErrCode == 0 {
		// 将数据记录到Redis数据库一份

		seconds := 7150
		err = global.REDIS.Set(ctx, "WeChatWorkJsAPITicketToken", ReturnData.Ticket, time.Duration(seconds)*time.Second).Err()

		if err != nil {
			zap.L().Error("Redis Set Key Error", zap.String("keys", "WeChatWorkJsAPITicketToken"), zap.String("value", ReturnData.Ticket), zap.Error(err))
		}
	}
	return ReturnData.Ticket
}

// GetJsAPITicket 获取应用的jsapi_ticket
func GetJsAPITicket() (Ticket string) {
	token := GetWechatAgentToken()
	fmt.Println(token)
	ctx := context.Background()
	result, err := global.REDIS.Get(ctx, "WeChatTicketToken").Result()

	if err == redis.Nil {
		zap.L().Error("RedisKey WeChatTicketToken does not exist", zap.Error(err))

	} else if err != nil {
		zap.L().Error("RedisKey WeChatTicketToken does not exist", zap.Error(err))
	}

	// 若有缓存直接返回
	if len(result) > 5 && result != "" {
		return result
	}
	data, err := HttpClient(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".WorkWechatUrl")+"/cgi-bin/ticket/get?access_token="+token+"&type=agent_config", "GET", "", "")
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}
	var ReturnData GetTicketReturnStruct
	err = json.Unmarshal(data, &ReturnData)
	if err != nil {
		zap.L().Error("WeChatGetTicketReturnStruct接口返回数据序列化失败！", zap.Error(err))
	}

	if ReturnData.ErrCode == 0 {
		// 将数据记录到Redis数据库一份
		seconds := 7150
		err = global.REDIS.Set(ctx, "WeChatTicketToken", ReturnData.Ticket, time.Duration(seconds)*time.Second).Err()

		if err != nil {
			zap.L().Error("Redis Set Key WeChatTicketToken Error", zap.String("keys", "WeChatTicketToken"), zap.String("value", ReturnData.Ticket), zap.Error(err))
		}
	}
	return ReturnData.Ticket
}

// GenerateSignature  企业微信JS-SDK使用权限签名算法
func GenerateSignature(noncestr string, jsapiTicket string, timestamp int64, url string) string {
	string1 := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", jsapiTicket, noncestr, timestamp, url)
	h := sha1.New()
	h.Write([]byte(string1))
	signature := fmt.Sprintf("%x", h.Sum(nil))
	return signature
}
