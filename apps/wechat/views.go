package wechat

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"ginDemo/global"
	"ginDemo/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
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

// GetWorkJsAPITicketToken 获取企业的jsapi_ticket
func GetWorkJsAPITicketToken(c *gin.Context) {
	token := utils.GetWorkJsAPITicket()
	if token != "" {
		returnData := map[string]interface{}{
			"jsapi_token": token,
			"errcode":     0,
			"errmsg":      "ok",
			"expires_in":  1800,
		}
		c.JSON(http.StatusOK, returnData)
		return
	} else {
		returnData := map[string]interface{}{
			"jsapi_token": "",
			"errcode":     0,
			"errmsg":      "ok",
			"expires_in":  1800,
		}
		c.JSON(http.StatusOK, returnData)
		return
	}
}

// GetAgentTicketToken 获取应用的jsapi_ticket
func GetAgentTicketToken(c *gin.Context) {
	token := utils.GetJsAPITicket()
	if token != "" {
		returnData := map[string]interface{}{
			"jsapi_token": token,
			"errcode":     0,
			"errmsg":      "ok",
			"expires_in":  1800,
		}
		c.JSON(http.StatusOK, returnData)
		return
	} else {
		returnData := map[string]interface{}{
			"jsapi_token": "",
			"errcode":     0,
			"errmsg":      "ok",
			"expires_in":  1800,
		}
		c.JSON(http.StatusOK, returnData)
		return
	}
}

// GetWorkConfig 获取企业微信config
func GetWorkConfig(c *gin.Context) {
	url := c.DefaultQuery("url", "https://qyapi.ud0.com.cn")
	noncestr := "Wm3WZYTPz0wzccnW"
	jsapi_ticket := utils.GetWorkJsAPITicket()
	timestamp := time.Now().Unix()

	sign := utils.GenerateSignature(noncestr, jsapi_ticket, timestamp, url)
	returnData := map[string]interface{}{
		"noncestr":     noncestr,
		"corpId":       "wwab4a127c8713c62b",
		"jsapi_ticket": jsapi_ticket,
		"timestamp":    timestamp,
		"url":          url,
		"sign":         sign,
	}

	c.JSON(http.StatusOK, returnData)
}

// GetWorkAgentConfig 获取企业微信agent_config
func GetWorkAgentConfig(c *gin.Context) {
	url := c.DefaultQuery("url", "https://qyapi.ud0.com.cn")
	nonceStr := "Wm3WZYTPz0wzccnW"
	jsapiTicket := utils.GetJsAPITicket()
	timestamp := time.Now().Unix()
	corpId := "wwab4a127c8713c62b" // 必填，企业微信的corpid，必须与当前登录的企业一致
	agentId := "1000008"           // 必填，企业微信的应用id （e.g. 1000247）
	sign := utils.GenerateSignature(nonceStr, jsapiTicket, timestamp, url)
	returnData := map[string]interface{}{
		"corpid":       corpId,
		"agentId":      agentId,
		"noncestr":     nonceStr,
		"jsapi_ticket": jsapiTicket,
		"timestamp":    timestamp,
		"sign":         sign,
	}

	c.JSON(http.StatusOK, returnData)
}

// GetWorkUserData 获取用户信息
func GetWorkUserData(c *gin.Context) {
	userID := c.DefaultQuery("userid", "")
	if userID == "" {
		returnData := map[string]interface{}{
			"code": 201,
			"msg":  "查询失败，userid为空",
		}
		c.JSON(http.StatusOK, returnData)
		return
	}

	token := utils.GetWechatToken()
	tempData := make([]string, 0, 4)
	tempData = append(tempData, userID)
	requestData := map[string]interface{}{
		"external_userid_list":       tempData,
		"need_enter_session_context": 1,
	}
	marshal, err := json.Marshal(requestData)
	if err != nil {
		zap.L().Error("Json序列化失败", zap.Error(err))
	}

	data, err := utils.HttpClient(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".WorkWechatUrl")+"/cgi-bin/kf/customer/batchget?access_token="+token, "POST", string(marshal), "")
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}
	var ReturnData GetWorkUserDataStruct
	err = json.Unmarshal(data, &ReturnData)
	if err != nil {
		zap.L().Error("GetWorkUserData接口返回数据序列化失败！", zap.Error(err))
	}
	// 将查询到的数据保存到MongoDB中一份，库名workWechat，表名userdata

	collection := global.MONGO.Database("workWechat").Collection("userdata")
	// 判断文档是否存在，存在则更新，不存在则新增
	var result MongoDBUserDataStruct
	filter := bson.M{"userid": userID}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		zap.L().Error("查询数据失败", zap.Error(err))
	}
	if result.Status == "ok" {
		// 更新的内容
		update := bson.M{
			"$set": bson.M{
				"status":                ReturnData.ErrMsg,
				"userid":                ReturnData.CustomerList[0].ExternalUserid,
				"nickname":              ReturnData.CustomerList[0].Nickname,
				"avatar":                ReturnData.CustomerList[0].Avatar,
				"gender":                ReturnData.CustomerList[0].Gender,
				"unionid":               ReturnData.CustomerList[0].Unionid,
				"enter_session_context": ReturnData.CustomerList[0].EnterSessionContext,
			},
		}

		updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			zap.L().Error("更新数据失败", zap.Error(err))
		}

		zap.L().Info("更新数据成功", zap.Any("Matched documents and updated  documents", updateResult.ModifiedCount))

		err = collection.FindOne(context.TODO(), filter).Decode(&result)
		returnData := map[string]interface{}{
			"code": 200,
			"msg":  "查询成功",
			"data": result,
		}
		c.JSON(http.StatusOK, returnData)
		return
	} else {
		s1 := MongoDBUserDataStruct{
			ReturnData.ErrMsg,
			userID,
			ReturnData.CustomerList[0].Nickname,
			ReturnData.CustomerList[0].Avatar,
			ReturnData.CustomerList[0].Gender,
			ReturnData.CustomerList[0].Unionid,
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			ReturnData.CustomerList[0].EnterSessionContext,
		}
		insertResult, err := collection.InsertOne(context.TODO(), s1)
		if err != nil {
			zap.L().Error("插入数据失败", zap.Error(err))
		}
		zap.L().Info("插入数据成功", zap.Any("data", insertResult.InsertedID))
		err = collection.FindOne(context.TODO(), filter).Decode(&result)
		returnData := map[string]interface{}{
			"code": 200,
			"msg":  "查询成功",
			"data": result,
		}

		c.JSON(http.StatusOK, returnData)
		return
	}
}

// SaveWorkUserData 保存用户信息
func SaveWorkUserData(c *gin.Context) {
	b, _ := c.GetRawData()
	var tempData UserDataStruct
	err := json.Unmarshal(b, &tempData)
	if err != nil {
		zap.L().Error("Json序列化UserDataStruct失败", zap.Error(err))
	}

	// 将接收到的数据保存到MongoDB中，库名workWechat，表名userdata
	collection := global.MONGO.Database("workWechat").Collection("userdata")
	// 更新的条件
	filter := bson.M{"userid": tempData.UserID}

	// 更新的内容
	update := bson.M{
		"$set": bson.M{
			"username":     tempData.Username,
			"mobile":       tempData.Mobile,
			"deviceNumber": tempData.DeviceNumber,
			"deviceModel":  tempData.DeviceModel,
			"IccID":        tempData.IccID,
			"operator":     tempData.Operator,
			"address":      tempData.Address,
			"comment":      tempData.Comment,
		},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		zap.L().Error("更新数据失败", zap.Error(err))
		returnData := map[string]interface{}{
			"code": 201,
			"msg":  "更新失败",
			"data": err,
		}

		c.JSON(http.StatusOK, returnData)
	}
	zap.L().Info("更新数据成功", zap.Any("Matched documents and updated  documents", updateResult.ModifiedCount))
	returnData := map[string]interface{}{
		"code": 200,
		"msg":  "更新成功",
		"data": "",
	}

	c.JSON(http.StatusOK, returnData)

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
