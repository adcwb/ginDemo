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
	"sync"
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

type TransServiceStateStruct struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	MsgCode string `json:"msg_code"`
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
	MsgList    []struct {
		MsgId    string `json:"msgid"`
		SendTime int    `json:"send_time"`
		Origin   int    `json:"origin"`
		MsgType  string `json:"msgtype"`
		Event    struct {
			EventType      string `json:"event_type"`
			OpenKfid       string `json:"open_kfid"`
			ExternalUserid string `json:"external_userid"`
			Scene          string `json:"scene"`
			SceneParam     string `json:"scene_param"`
			WelcomeCode    string `json:"welcome_code"`
			WechatChannels struct {
				Nickname string `json:"nickname"`
				Scene    int    `json:"scene"`
			} `json:"wechat_channels"`
			FailMsgId        string `json:"fail_msgid"`
			FailType         int    `json:"fail_type"`
			ServiceUserid    string `json:"servicer_userid"`
			Status           int    `json:"status"`
			ChangeType       int    `json:"change_type"`
			OldServiceUserid string `json:"old_servicer_userid"`
			NewServiceUserid string `json:"new_servicer_userid"`
			MsgCode          string `json:"msg_code"`
			RecallMsgId      string `json:"recall_msgid"`
			RejectSwitch     int    `json:"reject_switch"`
		} `json:"event"`
		OpenKfId       string `json:"open_kfid,omitempty"`
		ExternalUserid string `json:"external_userid,omitempty"`
		Text           struct {
			Content string `json:"content"`
			MenuId  string `json:"menu_id,omitempty"`
		} `json:"text,omitempty"`

		Image struct {
			MediaId string `json:"media_id,omitempty"`
		} `json:"image,omitempty"`

		Voice struct {
			MediaId string `json:"media_id"`
		} `json:"voice"`

		Video struct {
			MediaId string `json:"media_id"`
		} `json:"video"`

		File struct {
			MediaId string `json:"media_id"`
		} `json:"file"`

		Location struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Name      string  `json:"name"`
			Address   string  `json:"address"`
		} `json:"location"`
	} `json:"msg_list"`
	HasMore int `json:"has_more"`
}

type MongoDBUserScheduleStruct struct {
	Userid         string `bson:"userid"`
	ServiceState   int    `bson:"serviceState"`
	SendTime       int    `bson:"sendTime"`
	MsgType        string `bson:"msgType"`
	OpenKfId       string `bson:"openKfId"`
	ExternalUserid string `bson:"externalUserid"`
}

type SendMsgOnEventStruct struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	MsgId   string `json:"msgid"`
}

type ChatRecord struct {
	MsgId       string `json:"msgid"`
	Sender      string `bson:"sender"`
	Receiver    string `bson:"receiver"`
	Message     string `bson:"message"`
	MessageType string `bson:"messageType"`
	Timestamp   int    `bson:"timestamp"`
}
type MessageSession struct {
	Userid string       `bson:"userid"`
	Data   []ChatRecord `bson:"data"`
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

func (q *Queue[T]) IndexOf(item T) int {
	for i, val := range q.items {
		if val == item {
			return i + 1
		}
	}
	return -1
}

func (q *Queue[T]) Length() int {
	return len(q.items)
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
func GetServiceState(kfID, userID string) {
	token := GetWechatToken()
	var ReturnData GetServiceStateStruct
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
	if ReturnData.ErrCode == 0 {
		// 更新MongoDB会话状态
		zap.L().Debug("更新MongoDB数据库：", zap.String("external_userid", userID), zap.Int("ServiceState", ReturnData.ServiceState))
		if ReturnData.ServiceState == 0 {
			TransServiceState(kfID, userID, "", 1)
			go UpdateMongoMessage("workWechat", "userSchedule", userID, ReturnData.ServiceState)

		} else {
			go UpdateMongoMessage("workWechat", "userSchedule", userID, ReturnData.ServiceState)
		}
	}
	zap.L().Info("获取会话状态GetServiceState返回数据：", zap.Any("data", ReturnData))
}

// UpdateMongoMessage 更新MongoDB数据函数
func UpdateMongoMessage(dbName, tableName, userID string, State int) {
	var result MongoDBUserScheduleStruct
	//collection := global.MONGO.Database("workWechat").Collection("userSchedule")
	collection := global.MONGO.Database(dbName).Collection(tableName)
	filter := bson.D{{"userid", userID}}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		zap.L().Error("查询数据失败", zap.Error(err))
	}
	if result.Userid != "" && err == nil {
		update := bson.M{
			"$set": bson.M{
				"serviceState": State,
			},
		}
		updateResult, updateErr := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			zap.L().Error("更新数据失败", zap.Error(updateErr))
		}
		zap.L().Info("更新数据成功", zap.Any("Matched documents and updated  documents", updateResult.ModifiedCount))
	}
}

// SaveSession 将收到的用户信息保存到MongoDB中
func SaveSession(dbName, tableName, Sender, Receiver, Message, MessageType, UserId, MsgId string, Timestamp int) {
	result := MessageSession{}
	collection := global.MONGO.Database(dbName).Collection(tableName)
	filter := bson.D{{"userid", UserId}}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		zap.L().Error("查询数据失败", zap.Error(err))
	}
	if result.Userid != "" && err == nil {
		update := bson.M{
			"$set": bson.M{
				"data": append(result.Data, ChatRecord{
					MsgId:       MsgId,
					Sender:      Sender,
					Receiver:    Receiver,
					Message:     Message,
					MessageType: MessageType,
					Timestamp:   Timestamp,
				}),
			},
		}
		updateResult, updateErr := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			zap.L().Error("更新数据失败", zap.Error(updateErr))
		}
		zap.L().Info("更新数据成功", zap.Any("Matched documents and updated  documents", updateResult.ModifiedCount))
	} else {
		temp := make([]ChatRecord, 0, 16)
		temp = append(temp, ChatRecord{
			MsgId:       MsgId,
			Sender:      Sender,
			Receiver:    Receiver,
			Message:     Message,
			MessageType: MessageType,
			Timestamp:   Timestamp,
		})

		data := MessageSession{
			Userid: UserId,
			Data:   temp,
		}

		insertResult, err := collection.InsertOne(context.TODO(), data)
		if err != nil {
			zap.L().Error("插入数据失败", zap.Error(err))
		}
		zap.L().Info("插入数据成功", zap.Any("data", insertResult.InsertedID))
	}
}

// UpdateMongoMsgSendFail 记录消息发送失败事件
func UpdateMongoMsgSendFail(dbName, tableName, userID, kfID, FID, FType string) {
	type MongoMsgSendFaiStruct struct {
		OpenKfId       string `bson:"open_kfid"`
		ExternalUserid string `bson:"external_userid"`
		FailMsgId      string `bson:"fail_msgid"`
		FailType       string `bson:"fail_type"`
	}
	collection := global.MONGO.Database(dbName).Collection(tableName)

	data := MongoMsgSendFaiStruct{
		OpenKfId:       kfID,
		ExternalUserid: userID,
		FailMsgId:      FID,
		FailType:       FType,
	}

	insertResult, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		zap.L().Error("插入数据失败", zap.Error(err))
	}
	zap.L().Info("插入数据成功", zap.Any("data", insertResult.InsertedID))

}

// UpdateMongoServiceStatusChange 记录接待人员接待状态变更
func UpdateMongoServiceStatusChange(dbName, tableName, sID, kfID string, status int) {
	type MongoServiceStatusChangeStruct struct {
		ServiceUserid string `json:"servicer_userid"`
		OldStatus     int    `json:"old_status"`
		NewStatus     int    `json:"new_status"`
		OpenKfId      string `json:"open_kfid"`
	}
	collection := global.MONGO.Database(dbName).Collection(tableName)
	oldStatus := 0
	if status == 1 {
		oldStatus = 2
	} else {
		oldStatus = 1
	}
	data := MongoServiceStatusChangeStruct{
		ServiceUserid: sID,
		OldStatus:     oldStatus,
		NewStatus:     status,
		OpenKfId:      kfID,
	}

	insertResult, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		zap.L().Error("插入数据失败", zap.Error(err))
	}
	zap.L().Info("插入数据成功", zap.Any("data", insertResult.InsertedID))

}

func UpdateMongoRejectCustomerMsgSwitchChange(dbName, tableName, userID, kfID, sID string, status int) {
	type MongoRejectCustomerMsgSwitchChangeStruct struct {
		ServiceUserid  string `json:"servicer_userid"`
		OpenKfId       string `json:"open_kfid"`
		ExternalUserid string `json:"external_userid"`
		RejectSwitch   string `json:"reject_switch"`
	}
	collection := global.MONGO.Database(dbName).Collection(tableName)
	oldStatus := ""
	if status == 1 {
		oldStatus = "接待人员拒收了客户消息"
	} else if status == 0 {
		oldStatus = "接待人员取消拒收客户消息"
	}
	data := MongoRejectCustomerMsgSwitchChangeStruct{
		ServiceUserid:  sID,
		OpenKfId:       kfID,
		ExternalUserid: userID,
		RejectSwitch:   oldStatus,
	}

	insertResult, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		zap.L().Error("插入数据失败", zap.Error(err))
	}
	zap.L().Info("插入数据成功", zap.Any("data", insertResult.InsertedID))

}

// TransServiceState 变更会话状态
func TransServiceState(kfID, userID, jdKf string, state int) (ReturnData TransServiceStateStruct) {
	token := GetWechatToken()
	bodyData := make(map[string]interface{})
	if state == 4 || state == 2 || state == 1 {
		bodyData = map[string]interface{}{
			// 结束会话不用分配客服ID
			"open_kfid":       kfID,
			"external_userid": userID,
			"service_state":   state,
		}
	} else {
		bodyData = map[string]interface{}{
			"open_kfid":       kfID,
			"external_userid": userID,
			"service_state":   state,
			"servicer_userid": jdKf,
		}
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

	if state == 1 || state == 2 {
		if ReturnData.MsgCode != "" {
			go SendMsgOnEvent(ReturnData.MsgCode, "正在排队中，预计还需要5分钟~")
		}

		go SendMsgData(userID, kfID, "这是普通消息，您正在排队中，预计还需要5分钟~")

	} else if state == 3 {
		if ReturnData.MsgCode != "" {
			go SendMsgOnEvent(ReturnData.MsgCode, "你已进入人工会话")
		}
		go SendMsgData(userID, kfID, "这是普通消息2，您正在排队中，预计还需要5分钟~")
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
	err = global.REDIS.Set(ctx, "syncMsgNextCursor", ReturnData.NextCursor, 0).Err()

	if err != nil {
		zap.L().Error("Redis Set Key Error", zap.String("keys", "syncMsgNextCursor"), zap.String("value", ReturnData.NextCursor), zap.Error(err))
	}

	// MongoDB 调度用户数据
	collection := global.MONGO.Database("workWechat").Collection("userSchedule")

	// 将同步到的数据全部缓存到MongoDB中
	for _, v := range ReturnData.MsgList {
		var result MongoDBUserScheduleStruct
		filter := bson.D{{"userid", v.ExternalUserid}}
		err = collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			zap.L().Error("查询数据失败", zap.Error(err))
		}
		if result.Userid != "" && err == nil {
			update := bson.M{
				"$set": bson.M{
					"sendTime":       v.SendTime,
					"msgType":        v.MsgType,
					"openKfId":       v.OpenKfId,
					"externalUserid": v.ExternalUserid,
				},
			}

			updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				zap.L().Error("更新数据失败", zap.Error(err))
			}

			zap.L().Info("更新数据成功", zap.Any("Matched documents and updated  documents", updateResult.ModifiedCount))

		} else {
			if v.ExternalUserid != "" {
				s1 := MongoDBUserScheduleStruct{
					v.ExternalUserid,
					0,
					v.SendTime,
					v.MsgType,
					v.OpenKfId,
					v.ExternalUserid,
				}
				insertResult, err := collection.InsertOne(context.TODO(), s1)
				if err != nil {
					zap.L().Error("插入数据失败", zap.Error(err))
				}

				zap.L().Info("插入数据成功", zap.Any("data", insertResult.InsertedID))
			}
		}

		// 将解析到的聊天记录保存一份
		if v.Origin == 3 {
			// 微信客户发送的消息
			go SaveSession("workWechat", "userMessageSession", v.ExternalUserid, v.OpenKfId, v.Text.Content, v.MsgType, v.ExternalUserid, v.MsgId, v.SendTime)
		} else if v.Origin == 5 {
			// 接待人员在企业微信客户端发送的消息
			go SaveSession("workWechat", "userMessageSession", v.OpenKfId, v.ExternalUserid, v.Text.Content, v.MsgType, v.ExternalUserid, v.MsgId, v.SendTime)
		} else if v.Origin == 4 && v.MsgType == "event" {
			zap.L().Error("v", zap.Any("v", v))
			zap.L().Error("v.Event", zap.Any("v.event", v.Event))
			// 系统推送的事件消息
			switch v.Event.EventType {
			case "enter_session":
				// 用户进入会话事件, ，并将用户信息放入调度队列中
				if v.Event.WelcomeCode != "" {
					// 将数据写入MongoDB
					if v.ExternalUserid != "" {
						s1 := MongoDBUserScheduleStruct{
							v.ExternalUserid,
							0,
							v.SendTime,
							v.MsgType,
							v.OpenKfId,
							v.ExternalUserid,
						}
						insertResult, err := collection.InsertOne(context.TODO(), s1)
						if err != nil {
							zap.L().Error("插入数据失败", zap.Error(err))
						}
						zap.L().Info("插入数据成功", zap.Any("data", insertResult.InsertedID))
					}
					// 发送欢迎语以及排队人数
					fmt.Println("enter_session, 用户进入会话")
					go SendMsgOnEvent(v.Event.WelcomeCode, "哈罗，尊敬的客户您好，很高兴为您服务，请问有什么可以帮您？")
					go SendMsgData(v.ExternalUserid, v.OpenKfId, "当前排队人数为10人, 请耐心等待~")
				}
			case "msg_send_fail":
				// 消息发送失败事件, 记录MongoDB
				FailTypeMessage := ""
				switch v.Event.FailType {
				case 0:
					FailTypeMessage = "未知原因"
				case 1:
					FailTypeMessage = "客服帐号已删除"
				case 2:
					FailTypeMessage = "应用已关闭"
				case 4:
					FailTypeMessage = "会话已过期，超过48小时"
				case 5:
					FailTypeMessage = "会话已关闭"
				case 6:
					FailTypeMessage = "超过5条限制"
				case 7:
					FailTypeMessage = "未绑定视频号"
				case 8:
					FailTypeMessage = "主体未验证"
				case 9:
					FailTypeMessage = "未绑定视频号且主体未验证"
				case 10:
					FailTypeMessage = "用户拒收"
				default:
					FailTypeMessage = ""
				}
				zap.L().Error("接收到微信事件发送消息失败", zap.String("fail_type", FailTypeMessage))
				go UpdateMongoMsgSendFail("workWechat", "userMsgSendFail", v.Event.ExternalUserid, v.Event.OpenKfid, v.Event.FailMsgId, FailTypeMessage)

			case "servicer_status_change":
				// 接待人员接待状态变更事件 如从接待中变更为停止接待
				go UpdateMongoServiceStatusChange("workWechat", "userServiceStatusChange", v.Event.ServiceUserid, v.Event.OpenKfid, v.Event.Status)

			case "session_status_change":
				/*
					会话状态变更事件:
						1-从接待池接入会话
						2-转接会话
						3-结束会话
						4-重新接入已结束/已转接会话
				*/
				fmt.Println("session_status_change, 会话状态变更")
				switch v.Event.ChangeType {
				case 1:

					if v.Event.MsgCode != "" {
						go SendMsgOnEvent(v.Event.WelcomeCode, "提供设备号，姓名，手机号")
					}

				case 2:
				case 3:
					fmt.Println("本次会话已结束")
					// TODO 客服若主动结束会话，此处将减少客服消息队列中的人数，客户进线时记录客户进线的时间，若超时则自动结束会话

					if v.Event.MsgCode != "" {
						go SendMsgOnEvent(v.Event.WelcomeCode, "评价发送")
					}
				case 4:

				default:

				}

			case "user_recall_msg":
				// 用户撤回消息事件

			case "servicer_recall_msg":
				// 接待人员撤回消息事件

			case "reject_customer_msg_switch_change":
				// 拒收客户消息变更事件
				go UpdateMongoRejectCustomerMsgSwitchChange("workWechat", "userRejectCustomerMsgSwitchChange", v.Event.ExternalUserid, v.Event.OpenKfid, v.Event.ServiceUserid, v.Event.RejectSwitch)

			default:
				fmt.Println("无效的输入！")
			}
		}
	}
	return ReturnData
}

// SendMsgOnEvent 发送事件响应消息
func SendMsgOnEvent(code, message string) (ReturnData SendMsgOnEventStruct) {
	worktoken := GetWechatToken()
	temp := new(WorkSendData)
	if message == "" {
		message = "欢迎咨询，你已进入人工会话!"
	}
	bodyData := map[string]interface{}{
		"code":    code,
		"msgtype": "text",
		"text":    temp.Text(message),
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
	zap.L().Info("发送事件响应消息返回信息", zap.Any("ReturnData", ReturnData))
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
	zap.L().Info("发送普通消息返回信息", zap.Any("ReturnData", ReturnData))
	return ReturnData
}

// 客服结构体
type CustomerService struct {
	ID       int        // 客服ID
	IsPaused bool       // 是否暂停接入客户
	IsBusy   bool       // 是否正在处理客户
	Lock     sync.Mutex // 锁，控制并发处理客户
}

// 客户结构体
type Customer struct {
	ID int // 客户ID
}

// 客服队列
type CustomerServiceQueue struct {
	Services []*CustomerService // 客服队列
}

// 客户咨询队列
type CustomerQueue struct {
	Customers []*Customer // 客户队列
}

// Scheduler 调度器结构体
type Scheduler struct {
	ServiceQueue  *CustomerServiceQueue // 客服队列
	CustomerQueue *CustomerQueue        // 客户咨询队列
	Lock          sync.Mutex            // 锁，控制并发访问调度器
}

// NewCustomerServiceQueue 初始化客服队列
func NewCustomerServiceQueue() *CustomerServiceQueue {
	return &CustomerServiceQueue{
		Services: make([]*CustomerService, 0),
	}
}

// NewCustomerQueue 初始化客户咨询队列
func NewCustomerQueue() *CustomerQueue {
	return &CustomerQueue{
		Customers: make([]*Customer, 0),
	}
}

// NewScheduler 初始化调度器
func NewScheduler(serviceQueue *CustomerServiceQueue, customerQueue *CustomerQueue) *Scheduler {
	return &Scheduler{
		ServiceQueue:  serviceQueue,
		CustomerQueue: customerQueue,
	}
}

// AddService 将客服加入队列
func (q *CustomerServiceQueue) AddService(service *CustomerService) {
	q.Services = append(q.Services, service)
}

// AddCustomer 将客户加入队列
func (q *CustomerQueue) AddCustomer(customer *Customer) {
	q.Customers = append(q.Customers, customer)
}

// AllocateCustomerToService 分配客户给客服
func (s *Scheduler) AllocateCustomerToService(service *CustomerService, customer *Customer) {
	// 尝试获取锁
	service.Lock.Lock()
	defer service.Lock.Unlock()

	// 标记客服正在处理客户
	service.IsBusy = true

	// 处理客户
	// ...

	// 标记客服处理完客户
	service.IsBusy = false
}

// Run 调度器主循环
func (s *Scheduler) Run() {
	for {
		// 获取锁，避免并发访问客服队列和客户咨询队列
		s.Lock.Lock()

		// 检查客服队列是否有客服
		if len(s.ServiceQueue.Services) > 0 {
			// 遍历客服队列
			for _, service := range s.ServiceQueue.Services {
				// 如果客服正在处理客户，则跳过
				if service.IsBusy {
					continue
				}

				// 如果客服设置了暂停接入客户，则跳过
				if service.IsPaused {
					continue
				}

				// 检查客户咨询队列是否有客户
				if len(s.CustomerQueue.Customers) > 0 {
					// 取出客户
					customer := s.CustomerQueue.Customers[0]
					s.CustomerQueue.Customers = s.CustomerQueue.Customers[1:]

					// 分配客户给客服
					go s.AllocateCustomerToService(service, customer)
				}
			}
		}

		// 释放锁
		s.Lock.Unlock()

		// 等待一段时间后再次执行
		time.Sleep(1 * time.Second)
	}
}

func RunningQuery() {
	// 初始化客服队列和客户咨询队列
	serviceQueue := NewCustomerServiceQueue()
	customerQueue := NewCustomerQueue()

	// 创建客服和客户，并加入队列
	service1 := &CustomerService{ID: 1}
	service2 := &CustomerService{ID: 2}
	customer1 := &Customer{ID: 1}
	customer2 := &Customer{ID: 2}
	serviceQueue.AddService(service1)
	serviceQueue.AddService(service2)
	customerQueue.AddCustomer(customer1)
	customerQueue.AddCustomer(customer2)

	// 创建调度器并运行
	scheduler := NewScheduler(serviceQueue, customerQueue)
	scheduler.Run()
}
