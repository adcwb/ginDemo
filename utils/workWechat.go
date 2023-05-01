package utils

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"ginDemo/global"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
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

type GetKfIDListStruct struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccountList []struct {
		OpenKfId        string `json:"open_kfid"`
		Name            string `json:"name"`
		Avatar          string `json:"avatar"`
		ManagePrivilege bool   `json:"manage_privilege"`
	} `json:"account_list"`
}

type GetServiceListStruct struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	ServiceList []struct {
		Userid       string `json:"userid,omitempty"`
		Status       int    `json:"status,omitempty"`
		DepartmentId int    `json:"department_id,omitempty"`
	} `json:"servicer_list"`
}

type SetServiceStruct struct {
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	ResultList []struct {
		Userid       string `json:"userid,omitempty"`
		ErrCode      int    `json:"errcode"`
		ErrMsg       string `json:"errmsg"`
		DepartmentId int    `json:"department_id,omitempty"`
	} `json:"result_list"`
}

type GetServiceStateStruct struct {
	ErrCode        int    `json:"errcode"`
	ErrMsg         string `json:"errmsg"`
	ServiceState   int    `json:"service_state"`
	ServicerUserid string `json:"servicer_userid"`
}

type GetMessageStruct struct {
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	NextCursor string `json:"next_cursor"`
	HasMore    int    `json:"has_more"`
	MsgList    []struct {
		MsgID          string `json:"msgid"`
		OpenKfID       string `json:"open_kfid"`
		ExternalUserid string `json:"external_userid"`
		SendTime       int    `json:"send_time"`
		Origin         int    `json:"origin"`
		ServiceUserid  string `json:"servicer_userid"`
		MsgType        string `json:"msgtype"`
		Event          struct {
			EventType      string `json:"event_type"`
			Scene          int    `json:"scene"`
			OpenKfID       string `json:"open_kfid"`
			ExternalUserid string `json:"external_userid"`
			WelcomeCode    string `json:"welcome_code"`
		} `json:"event"`
	} `json:"msg_list"`
}

type MongoDBUserScheduleStruct struct {
	userid         string
	serviceState   int
	sendTime       int
	msgType        string
	openKfId       string
	externalUserid string
}

type SendMsgOnEventStruct struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	MsgId   string `json:"msgid"`
}

// Queue 创建一个队列
type Queue[T string | int64] struct {
	items []T
}

func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

func (q *Queue[T]) Dequeue() T {
	if len(q.items) == 0 {
		var temp T
		return temp
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

type WorkSendData struct{}

// Text 发送文本
func (send *WorkSendData) Text(message string) map[string]string {
	data := map[string]string{
		"content": message,
	}

	return data
}

// Images 发送图片
func (send *WorkSendData) Images(message string) map[string]string {
	data := map[string]string{
		"content": message,
	}

	return data
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

// GetKfIDList 获取客服账号列表
func GetKfIDList() (data GetKfIDListStruct) {
	token := GetWechatToken()
	bodyData := map[string]int64{
		"offset": 0,
		"limit":  100,
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		zap.L().Error("Json序列化失败", zap.Error(err))
	}
	tempData, err := HttpClient("https://qyapi.weixin.qq.com/cgi-bin/kf/account/list?access_token="+token, "POST", string(body), "")
	var ReturnData GetKfIDListStruct
	err = json.Unmarshal(tempData, &ReturnData)
	if err != nil {
		zap.L().Error("Json序列化失败", zap.Error(err))
	}
	/*
		返回数据示例
		{wk5aTbYAAAfXatwdbTTfQB9JY_TI96vQ 点击联系在线客服 https://wework.qpic.cn/wwpic/329236_JjHmAOdGQGeRr5e_1678626977/0 false}
		{wk5aTbYAAAHF-IN-YZtmyn8s4369j47w 测试API https://wework.qpic.cn/wwpic/236392_Ic0lpXNlQIOpQMB_1678541469/0 true}
		{wk5aTbYAAAjGJ5KYoUiOB3wige3YTvHA 深圳优度通信科技有限公司客服 https://wwcdn.weixin.qq.com/node/wework/images/kf_head_image_url_1.png false}
		{wk5aTbYAAAzssxTaC6LK2NOhFPMKLkjw 优度无线宽带客服 http://wx.qlogo.cn/finderhead/Q3auHgzwzM78Qvu8sxQrQz4nrX7iac2BzqUnVOpbbjbiaBhTeUyKQicDg/0/0 false}
	*/
	if ReturnData.ErrMsg == "ok" {
		for i, v := range ReturnData.AccountList {
			fmt.Println(i, v)
		}
	}

	return ReturnData
}

// GetServiceList 获取客服账号接待人员列表
func GetServiceList(kfID string) (ReturnData GetServiceListStruct) {
	token := GetWechatToken()
	tempData, err := HttpClient("https://qyapi.weixin.qq.com/cgi-bin/kf/servicer/list?access_token="+token+"&open_kfid="+kfID, "GET", "", "")
	if err != nil {
		zap.L().Error("HTTP请求发送失败！", zap.Error(err))
	}

	err = json.Unmarshal(tempData, &ReturnData)
	if err != nil {
		zap.L().Error("Json序列化失败！", zap.Error(err))
	}
	return ReturnData
}

// SetService 添加客服账号接待人员
func SetService(kfID string, userData []string) {
	token := GetWechatToken()
	bodyData := map[string]interface{}{
		"open_kfid":   kfID,
		"userid_list": userData,
	}

	body, err := json.Marshal(bodyData)
	if err != nil {
		zap.L().Error("Json序列化失败！", zap.Error(err))
	}

	tempData, err := HttpClient("https://qyapi.weixin.qq.com/cgi-bin/kf/servicer/add?access_token="+token, "POST", string(body), "")
	if err != nil {
		zap.L().Error("HTTP请求发送失败！", zap.Error(err))
	}
	var ReturnData SetServiceStruct
	err = json.Unmarshal(tempData, &ReturnData)
	if err != nil {
		zap.L().Error("Json序列化失败！", zap.Error(err))
	}
}

// DelService 删除客服账号接待人员
func DelService(kfID string, userData []string) {
	token := GetWechatToken()
	bodyData := map[string]interface{}{
		"open_kfid":   kfID,
		"userid_list": userData,
	}

	body, err := json.Marshal(bodyData)
	if err != nil {
		zap.L().Error("Json序列化失败！", zap.Error(err))
	}

	tempData, err := HttpClient("https://qyapi.weixin.qq.com/cgi-bin/kf/servicer/del?access_token="+token, "POST", string(body), "")
	if err != nil {
		zap.L().Error("HTTP请求发送失败！", zap.Error(err))
	}
	var ReturnData SetServiceStruct
	err = json.Unmarshal(tempData, &ReturnData)
	if err != nil {
		zap.L().Error("Json序列化失败！", zap.Error(err))
	}
}

// GetServiceState 获取会话状态
func GetServiceState(kfID, userID string) (ReturnData GetServiceStateStruct) {
	token := GetWechatToken()
	bodyData := map[string]interface{}{
		"open_kfid":       kfID,
		"external_userid": userID,
	}

	body, err := json.Marshal(bodyData)
	if err != nil {
		zap.L().Error("Json序列化失败！", zap.Error(err))
	}

	tempData, err := HttpClient("https://qyapi.weixin.qq.com/cgi-bin/kf/service_state/get?access_token="+token, "POST", string(body), "")
	if err != nil {
		zap.L().Error("HTTP请求发送失败！", zap.Error(err))
	}

	err = json.Unmarshal(tempData, &ReturnData)
	if err != nil {
		zap.L().Error("Json序列化失败！", zap.Error(err))
	}
	return ReturnData
}

// TransServiceState 变更会话状态
func TransServiceState(kfID, userID, jdKf string, state int) (ReturnData SetServiceStruct) {
	token := GetWechatToken()
	bodyData := map[string]interface{}{
		"open_kfid":       kfID,
		"external_userid": userID,
		"service_state":   state,
		"servicer_userid": jdKf,
	}

	body, err := json.Marshal(bodyData)
	if err != nil {
		zap.L().Error("Json序列化失败！", zap.Error(err))
	}

	tempData, err := HttpClient("https://qyapi.weixin.qq.com/cgi-bin/kf/service_state/trans?access_token="+token, "POST", string(body), "")
	if err != nil {
		zap.L().Error("HTTP请求发送失败！", zap.Error(err))
	}
	err = json.Unmarshal(tempData, &ReturnData)
	if err != nil {
		zap.L().Error("Json序列化失败！", zap.Error(err))
	}

	return ReturnData
}

// GetMessage 同步消息
func GetMessage(kfID, token string) (ReturnData GetMessageStruct) {
	worktoken := GetWechatToken()

	// Redis加缓存
	ctx := context.Background()
	redisData := GetRedisKey(ctx, "syncMsgNextCursor")
	bodyData := make(map[string]interface{})
	if redisData == "" {
		bodyData = map[string]interface{}{
			"token":        token,
			"limit":        1000,
			"voice_format": 0,
			"open_kfid":    kfID,
		}
	} else {
		bodyData = map[string]interface{}{
			"cursor":       redisData,
			"token":        token,
			"limit":        1000,
			"voice_format": 0,
			"open_kfid":    kfID,
		}
	}

	body, err := json.Marshal(bodyData)
	if err != nil {
		zap.L().Error("Json序列化失败！", zap.Error(err))
	}

	tempData, err := HttpClient("https://qyapi.weixin.qq.com/cgi-bin/kf/sync_msg?access_token="+worktoken, "POST", string(body), "")
	if err != nil {
		zap.L().Error("HTTP请求发送失败！", zap.Error(err))
	}

	err = json.Unmarshal(tempData, &ReturnData)
	if err != nil {
		zap.L().Error("Json序列化失败！", zap.Error(err))
	}
	if redisData != "" {
		err = global.REDIS.Set(ctx, "syncMsgNextCursor", ReturnData.NextCursor, 0).Err()

		if err != nil {
			zap.L().Error("Redis Set Key Error", zap.String("keys", "syncMsgNextCursor"), zap.String("value", ReturnData.NextCursor), zap.Error(err))
		}
	}

	// MongoDB 调度用户数据
	collection := global.MONGO.Database("workWechat").Collection("userSchedule")

	// 将同步到的数据全部缓存到MongoDB中
	for _, v := range ReturnData.MsgList {
		var result MongoDBUserScheduleStruct
		filter := bson.M{"userid": v.ExternalUserid}
		err = collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			zap.L().Error("查询数据失败", zap.Error(err))
		}
		if result.userid != "" {
			update := bson.M{
				"$set": bson.M{
					"serviceState":   v.Event.Scene,
					"sendTime":       v.SendTime,
					"msgType":        v.MsgType,
					"openKfId":       v.OpenKfID,
					"externalUserid": v.ExternalUserid,
				},
			}

			updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				zap.L().Error("更新数据失败", zap.Error(err))
			}

			zap.L().Info("更新数据成功", zap.Any("Matched documents and updated  documents", updateResult.ModifiedCount))

		} else {
			s1 := MongoDBUserScheduleStruct{
				userid:         v.ExternalUserid,
				serviceState:   v.Event.Scene,
				sendTime:       v.SendTime,
				msgType:        v.MsgType,
				openKfId:       v.OpenKfID,
				externalUserid: v.ExternalUserid,
			}
			insertResult, err := collection.InsertOne(context.TODO(), s1)
			if err != nil {
				zap.L().Error("插入数据失败", zap.Error(err))
			}
			zap.L().Info("插入数据成功", zap.Any("data", insertResult.InsertedID))
		}
	}

	return ReturnData
}

// SendMsgOnEvent 发送事件响应消息
func SendMsgOnEvent(code string) (ReturnData SendMsgOnEventStruct) {
	worktoken := GetWechatToken()
	temp := new(WorkSendData)

	bodyData := map[string]interface{}{
		"code":    code,
		"msgtype": "text",
		"text":    temp.Text("欢迎咨询，你已进入人工会话!"),
	}

	body, err := json.Marshal(bodyData)
	if err != nil {
		zap.L().Error("Json序列化失败！", zap.Error(err))
	}

	tempData, err := HttpClient("https://qyapi.weixin.qq.com/cgi-bin/kf/send_msg_on_event?access_token="+worktoken, "POST", string(body), "")
	if err != nil {
		zap.L().Error("HTTP请求发送失败！", zap.Error(err))
	}
	err = json.Unmarshal(tempData, &ReturnData)
	if err != nil {
		zap.L().Error("Json序列化失败！", zap.Error(err))
	}

	return ReturnData
}

// SendMsgData 发送普通消息
func SendMsgData(userID, kfID, message string) (ReturnData SendMsgOnEventStruct) {
	worktoken := GetWechatToken()
	temp := new(WorkSendData)

	bodyData := map[string]interface{}{
		"touser":    userID,
		"open_kfid": kfID,
		"msgtype":   "text",
		"text":      temp.Text(message),
	}

	body, err := json.Marshal(bodyData)
	if err != nil {
		zap.L().Error("Json序列化失败！", zap.Error(err))
	}

	tempData, err := HttpClient("https://qyapi.weixin.qq.com/cgi-bin/kf/send_msg?access_token="+worktoken, "POST", string(body), "")
	if err != nil {
		zap.L().Error("HTTP请求发送失败！", zap.Error(err))
	}
	err = json.Unmarshal(tempData, &ReturnData)
	if err != nil {
		zap.L().Error("Json序列化失败！", zap.Error(err))
	}

	return ReturnData
}
