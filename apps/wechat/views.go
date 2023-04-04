package wechat

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"ginDemo/global"
	"ginDemo/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// GetWeChatAccessToken 获取企业微信Token
func GetWeChatAccessToken(c *gin.Context) {
	// 判断Redis数据库中是否有记录
	ctx := context.Background()

	result, err := global.REDIS.Get(ctx, "WeChatAccessToken").Result()

	if err == redis.Nil {
		zap.L().Error("RedisKey PhoneNumbers does not exist", zap.Error(err))

	} else if err != nil {
		zap.L().Error("RedisKey PhoneNumbers does not exist", zap.Error(err))
	}
	// 若有缓存直接返回
	if len(result) > 5 && result != "" {
		ReturnData := map[string]interface{}{
			"errcode":      0,
			"errmsg":       "ok",
			"access_token": result,
			"expires_in":   1800,
		}
		c.JSON(http.StatusOK, ReturnData)
		return
	}

	CorpID := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".WorkWechatCorpID")
	CorpSecret := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".WorkWechatCorpSecret")

	data, err := utils.HttpClient(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".WorkWechatUrl")+"/cgi-bin/gettoken?corpid="+CorpID+"&corpsecret="+CorpSecret, "GET", "", "")
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
		ctx := context.Background()
		seconds := 7150
		err = global.REDIS.Set(ctx, "WeChatAccessToken", ReturnData.AccessToken, time.Duration(seconds)*time.Second).Err()

		if err != nil {
			zap.L().Error("Redis Set Key Error", zap.String("keys", "WeChatAccessToken"), zap.String("value", ReturnData.AccessToken), zap.Error(err))
		}
		c.JSON(http.StatusOK, ReturnData)
	} else {
		c.JSON(http.StatusOK, ReturnData)
	}
}

// CallbackWechat 回调地址
func CallbackWechat(c *gin.Context) {
	method := c.Request.Method
	token := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".WorkWechatCorpToken")
	encodingAeskey := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".WorkWechatCorpEncodingAes")
	receiverId := global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".WorkWechatCorpReceiverId")
	wxcpt := utils.NewWXBizMsgCrypt(token, encodingAeskey, receiverId, utils.XmlType)
	if method == "GET" {
		// 解析出url上的参数值如下：
		verifyMsgSign, _ := c.GetQuery("msg_signature")
		verifyTimestamp, _ := c.GetQuery("timestamp")
		verifyNonce, _ := c.GetQuery("nonce")
		verifyEchoStr, _ := c.GetQuery("echoStr")
		echoStr, cryptErr := wxcpt.VerifyURL(verifyMsgSign, verifyTimestamp, verifyNonce, verifyEchoStr)
		if nil != cryptErr {
			zap.L().Error("verifyUrl fail!", zap.String("cryptErrMsg", cryptErr.ErrMsg), zap.Int("cryptErrCode", cryptErr.ErrCode))
		}
		zap.L().Info("verifyUrl success echoStr", zap.String("echoStr", string(echoStr)))
		// 验证URL成功，将sEchoStr返回
		c.JSON(http.StatusOK, string(echoStr))
		return
	} else if method == "POST" {
		reqMsgSign, _ := c.GetQuery("msg_signature")
		reqTimestamp, _ := c.GetQuery("timestamp")
		reqNonce, _ := c.GetQuery("nonce")
		// post请求的密文数据
		reqData, _ := c.GetRawData()

		msg, cryptErr := wxcpt.DecryptMsg(reqMsgSign, reqTimestamp, reqNonce, reqData)
		if nil != cryptErr {
			zap.L().Error("DecryptMsg fail!", zap.String("cryptErrMsg", cryptErr.ErrMsg), zap.Int("cryptErrCode", cryptErr.ErrCode))
		}
		zap.L().Info("after decrypt msg: ", zap.String("echoStr", string(msg)))
		// TODO: 解析出明文xml标签的内容进行处理
		// For example:

		var msgContent MsgContent
		err := xml.Unmarshal(msg, &msgContent)
		if nil != err {
			zap.L().Error("xml Unmarshal失败!", zap.Error(err))
		}

	} else {
		c.String(http.StatusNotFound, "404 page not found!")
		return
	}
}
