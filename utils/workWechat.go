package utils

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"ginDemo/global"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"strconv"
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
		zap.L().Error("RedisKey WeChatAccessToken does not exist", zap.Error(err))

	} else if err != nil {
		zap.L().Error("RedisKey WeChatAccessToken does not exist", zap.Error(err))
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
func GetServiceState(kfID, userID string) GetServiceStateStruct {
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
		fmt.Println("将客户信息写入MQ队列ExternalUseridQueue")
		if ReturnData.ServiceState != 3 && ReturnData.ServiceState != 4 {
			err = SendQueue("ExternalUseridQueue", userID)
			if err != nil {
				zap.L().Error("客户数据写入MQ失败！", zap.Error(err), zap.String("data", userID))
			}
		}
		// 更新MongoDB会话状态
		//zap.L().Debug("更新MongoDB数据库：", zap.String("external_userid", userID), zap.Int("ServiceState", ReturnData.ServiceState))
		//if ReturnData.ServiceState == 0 {
		//	TransServiceState(kfID, userID, "", 2)
		//	go UpdateMongoMessage("workWechat", "userSchedule", userID, ReturnData.ServiceState)
		//
		//} else {
		//	go UpdateMongoMessage("workWechat", "userSchedule", userID, ReturnData.ServiceState)
		//}
	}
	zap.L().Info("获取会话状态GetServiceState返回数据：", zap.Any("data", ReturnData))
	return ReturnData
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
func TransServiceState(kfID, userID string) (ReturnData TransServiceStateStruct) {
	token := GetWechatToken()
	ctx := context.Background()
	bodyData := make(map[string]interface{})
	// 取出可用的客服ID
	ServiceUserid, err1 := global.REDIS.LPop(ctx, "ServiceUseridUpQueue").Result()
	if err1 != nil {
		zap.L().Error("取出可用的客服ID失败！", zap.Error(err1))
	}

	if ServiceUserid == "" {
		zap.L().Error("当前客服在线队列ServiceUseridUpQueue为空，会话分配暂停中.......")
		return ReturnData
	}

	// 分配会话时候，先将客户会话状态改为2并且发送事件消息，并且判断当前是否有可接待的客服，如果有直接分配，若没有则将客户加入新的队列中，下次优先调度
	// 判断本次取出的客服是否可接待
	fmt.Println("=====>> ServiceUserid", ServiceUserid)
	member, err := global.REDIS.ZScore(ctx, "ServiceUseridData", ServiceUserid).Result()
	if err != nil {
		zap.L().Error("Redis ZScore ServiceUseridData Error! ", zap.Error(err))
	}
	fmt.Println("member 本次取出的客服接待的人数", member)
	// 最大接待十个人
	if member < 1 {
		// 将取到的接待人员重新放入队列中
		global.REDIS.RPush(ctx, "ServiceUseridUpQueue", ServiceUserid)
		bodyData = map[string]interface{}{
			"open_kfid":       kfID,
			"external_userid": userID,
			"service_state":   3,
			"servicer_userid": ServiceUserid,
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
		if ReturnData.MsgCode != "" {
			go SendMsgOnEvent(ReturnData.MsgCode, "尊敬的客户您好，系统提示您已进入人工会话，欢迎您的来访~")
		}
		// 会话分配成功，记录哪个客户分配给哪个客服，后期做会话超时
		if ReturnData.Errcode == 0 {
			// 会话成功，将当前正在接待的客服的接待人员数加1 ServiceUseridData
			err = global.REDIS.ZIncrBy(ctx, "ServiceUseridData", 1, ServiceUserid).Err()
			if err != nil {
				zap.L().Error("增加客服接待人数失败，请排查原因~", zap.Error(err))
			}

			ServiceUserHistoryMongo := global.MONGO.Database("workWechat").Collection("ServiceUserHistory")
			data := ServiceUserHistoryStruct{
				ServiceUserid:  ServiceUserid,
				OpenKfId:       kfID,
				ExternalUserid: userID,
				ServiceStatus:  1,
				ServiceData:    time.Now().Format("2006-01-02 15:04:05"),
			}

			insertResult, err := ServiceUserHistoryMongo.InsertOne(context.TODO(), data)
			if err != nil {
				zap.L().Error("插入数据失败", zap.Error(err))
			}
			zap.L().Info("插入数据成功", zap.Any("data", insertResult.InsertedID))
		}
		return ReturnData
	} else {
		// 循环指定次数 获取有序集合的成员数
		numberTemp := global.REDIS.ZCard(ctx, "ServiceUseridData").String() // 获取集合元素个数
		number, _ := strconv.Atoi(numberTemp)
		for i := 0; i <= number; i++ {
			// 循环时将上次的拿出的客服重新放入队列中，避免队列为空
			global.REDIS.RPush(ctx, "ServiceUseridUpQueue", ServiceUserid)
			var ServiceUseridErr error
			ServiceUserid, ServiceUseridErr = global.REDIS.LPop(ctx, "ServiceUseridUpQueue").Result()
			if ServiceUseridErr != nil {
				zap.L().Error(".REDIS LPop ServiceUseridUpQueue Error! ", zap.Error(ServiceUseridErr))
			}

			if ServiceUserid == "" {
				zap.L().Error("当前客服在线队列ServiceUseridUpQueue为空，会话分配暂停中.......")
				return ReturnData
			}

			member, err = global.REDIS.ZScore(context.Background(), "ServiceUseridData", ServiceUserid).Result()
			if err != nil {
				zap.L().Error("Redis ZScore ServiceUseridData Error! ", zap.Error(err))
			}

			if member < 1 {
				global.REDIS.RPush(ctx, "ServiceUseridUpQueue", ServiceUserid)
				bodyData = map[string]interface{}{
					"open_kfid":       kfID,
					"external_userid": userID,
					"service_state":   3,
					"servicer_userid": ServiceUserid,
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
				if ReturnData.Errcode == 0 {
					// 会话成功，将当前正在接待的客服的接待人员数加1 ServiceUseridData
					err = global.REDIS.ZIncrBy(ctx, "ServiceUseridData", 1, ServiceUserid).Err()
					if err != nil {
						zap.L().Error("增加客服接待人数失败，请排查原因~", zap.Error(err))
					}
				}
				if ReturnData.MsgCode != "" {
					go SendMsgOnEvent(ReturnData.MsgCode, "尊敬的客户您好，系统提示您已进入人工会话，欢迎您的来访~")
				} else {
					go SendMsgData(userID, kfID, "尊敬的客户您好，系统提示您已进入人工会话，欢迎您的来访~")
				}

				return ReturnData
			}
		}
		// 若五次循环都没有成功分配，则证明当前客服繁忙，将客户放入接待池中
		bodyData = map[string]interface{}{
			// 结束会话不用分配客服ID
			"open_kfid":       kfID,
			"external_userid": userID,
			"service_state":   2,
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
		// 将客户放入优先接待队列
		global.REDIS.RPush(ctx, "ExternalUseridUpQueue", userID)
		// 获取队列长度，并且将排队信息返回给客户
		tempNumber, err := global.REDIS.LLen(ctx, "ExternalUseridUpQueue").Result()
		if err != nil {
			zap.L().Error("Redis LLen ExternalUseridUpQueue Error!", zap.Error(err))
		}
		messageData := ""
		if tempNumber < 0 {
			messageData = "尊敬的客户您好, 当前客服坐席繁忙, 您正在排队中, 预计还有1人, 请您稍后......"
		} else if tempNumber == 0 {
			messageData = "尊敬的客户您好, 当前客服坐席繁忙, 您正在排队中, 预计还有1人, 请您稍后......"
		} else {
			// 索引，需要加1
			tempNumber = tempNumber + 1
			messageData = "尊敬的客户您好, 当前客服坐席繁忙, 您正在排队中, 预计还有" + strconv.FormatInt(tempNumber, 10) + "人, 请您稍后......"
		}
		if ReturnData.MsgCode != "" {
			go SendMsgOnEvent(ReturnData.MsgCode, messageData)
		} else {
			go SendMsgData(userID, kfID, messageData)
		}
	}
	return ReturnData
}

// TimeOutCheck 超时检测
func TimeOutCheck() {
	// 检测MongoDB中的数据，判断是否超时，超时则把会话状态改为4，并从数据库中移除数据 ServiceUserHistoryStruct
	ServiceUserHistoryMongo := global.MONGO.Database("workWechat").Collection("ServiceUserHistory")
	// 取出所有的数据，判断是否会话结束，若未结束则查看时间是否超过五分钟，超时直接发送信息提醒，并结束会话
	filter := bson.M{}
	findOptions := options.Find().SetSort(bson.M{"serviceData": 1})
	cur, err := ServiceUserHistoryMongo.Find(context.Background(), filter, findOptions)
	if err != nil {
		zap.L().Error("读取数据错误！", zap.Error(err))
	}

	for cur.Next(context.Background()) {
		var result ServiceUserHistoryStruct
		err1 := cur.Decode(&result)
		if err1 != nil {
			zap.L().Error("读取数据中发生错误，请检查！", zap.Error(err1))
		}

		// 判断会话开始时间师傅超过五分钟，超过则直接结束会话
		loc, _ := time.LoadLocation("Asia/Shanghai")
		now := time.Now().In(loc)
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", result.ServiceData, loc)
		// 判断会话开始时间是否超过十分钟
		if now.After(start.Add(10 * time.Minute)) {
			if result.ServiceUserid == "e0c48f2ff5b73bc696c2aa1a0f666f3e" {
				continue
			}
			// 如果超过十分钟则查询会话是否还是在线，若在线则中断会话，并且删除数据，减少对应客服的接待人数
			token := GetWechatToken()
			var ReturnData GetServiceStateStruct
			bodyData := map[string]interface{}{
				"open_kfid":       result.OpenKfId,
				"external_userid": result.ExternalUserid,
			}

			body, err2 := json.Marshal(bodyData)
			if err2 != nil {
				zap.L().Error("Json序列化失败！", zap.Error(err2))
			}

			tempData, err3 := HttpClient("https://qyapi.weixin.qq.com/cgi-bin/kf/service_state/get?access_token="+token, "POST", string(body), "")
			if err != nil {
				zap.L().Error("HTTP请求发送失败！", zap.Error(err3))
			}

			err = json.Unmarshal(tempData, &ReturnData)
			if err != nil {
				zap.L().Error("Json序列化失败！", zap.Error(err))
			}

			if ReturnData.ServiceState == 3 {
				// 设备在线的时候, 将会话状态变成4
				bodyData = map[string]interface{}{
					"open_kfid":       result.OpenKfId,
					"external_userid": result.ExternalUserid,
					"service_state":   4,
					"servicer_userid": result.ServiceUserid,
				}
				body1, err5 := json.Marshal(bodyData)
				if err5 != nil {
					zap.L().Error("Json序列化失败！", zap.Error(err5))
				}
				var ReturnDataTrans TransServiceStateStruct
				tempDataTrans, errTrans := HttpClient("https://qyapi.weixin.qq.com/cgi-bin/kf/service_state/trans?access_token="+token, "POST", string(body1), "")
				if errTrans != nil {
					zap.L().Error("HTTP请求发送失败！", zap.Error(errTrans))
				}
				err = json.Unmarshal(tempDataTrans, &ReturnDataTrans)
				if err != nil {
					zap.L().Error("Json序列化失败！", zap.Error(err))
				}
			}
			// 会话不在线的时候直接删除MongoDB信息
			one, errDel := ServiceUserHistoryMongo.DeleteOne(context.TODO(), bson.D{{"externalUserid", result.ExternalUserid}})
			if errDel != nil {
				zap.L().Error("删除MongoDB数据失败！", zap.Error(errDel), zap.String("data", strconv.FormatInt(one.DeletedCount, 10)))
			}
		}

	}

	if err1 := cur.Err(); err1 != nil {
		zap.L().Error("读取数据中发生错误，请检查！", zap.Error(err1))
	}

	err = cur.Close(context.Background())
	if err != nil {
		zap.L().Error("关闭会话失败！", zap.Error(err))
	}

	// TODO 有一些异常会话的客户，怎么处理
}

// InitWorkWechatData 初始化企业微信--微信客服需要的数据, 仅在项目运行的时候执行一次
func InitWorkWechatData(kfID string) {
	ReturnData := GetServiceList(kfID)
	// TODO 项目刚启动时候检测需要的队列是否存在，若不存在则新建，存在则清空，避免因为项目重启导致数据异常
	if ReturnData.ErrCode == 0 || ReturnData.ErrMsg == "ok" {
		for _, v := range ReturnData.ServiceList {
			// 保存客服接待人员数量
			global.REDIS.ZAdd(context.Background(), "ServiceUseridData", &redis.Z{
				Score:  0,
				Member: v.Userid,
			})
			if v.Status == 0 {
				// 将客服信息放入接待中队列
				global.REDIS.RPush(context.Background(), "ServiceUseridUpQueue", v.Userid)
			}
		}
	}

	go func() {
		for {
			time.Sleep(time.Second * 3)
			fmt.Println("将客户信息调度给客服进程正在执行中......")
			zap.L().Debug("将客户信息调度给客服进程正在执行中......")
			ControlMessage(kfID)
		}
	}()
}

// GetMessage 同步消息
func GetMessage(kfID, token string) (ReturnData GetMessageStruct) {
	worktoken := GetWechatToken()
	// Redis加缓存
	ctx := context.Background()
	redisData := GetRedisKey(ctx, "syncMsgNextCursor")
	bodyData := make(map[string]interface{})
	if redisData != "" {
		bodyData["cursor"] = redisData
	}
	bodyData["token"] = token
	bodyData["open_kfid"] = kfID
	bodyData["voice_format"] = 0
	bodyData["limit"] = 1000

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
	//collection := global.MONGO.Database("workWechat").Collection("userSchedule")

	// 将同步到的数据全部缓存到MongoDB中
	for _, v := range ReturnData.MsgList {

		// Redis中缓存正在会话中的客户ID KEY InConversationExternalUserid []string Redis Type List
		// 此处做判定，如果是正在会话中的则不再往队列中丢，筛选完成以后再次判断会话状态是否为已结束或未开始，状态等于4的剔除，其余的丢进MQ队列中
		// MQ中保存未调度会话ID的队列名称为ExternalUseridQueue
		result, err1 := global.REDIS.LRange(ctx, "InConversationExternalUserid", 0, -1).Result()
		if err != nil {
			zap.L().Error("Redis LRange InConversationExternalUserid Error！", zap.Error(err1))
		}
		fmt.Println(FindList(result, v.ExternalUserid))
		if FindList(result, v.ExternalUserid) == false {
			fmt.Println("开始获取会话状态", v.OpenKfId, v.ExternalUserid)
			GetServiceState(v.OpenKfId, v.ExternalUserid)
		}

		//将解析到的聊天记录保存一份
		if v.Origin == 3 {
			// 保存微信客户发送的消息
			go SaveSession("workWechat", "userMessageSession", v.ExternalUserid, v.OpenKfId, v.Text.Content, v.MsgType, v.ExternalUserid, v.MsgId, v.SendTime)
		} else if v.Origin == 5 {
			// 接待人员在企业微信客户端发送的消息
			go SaveSession("workWechat", "userMessageSession", v.OpenKfId, v.ExternalUserid, v.Text.Content, v.MsgType, v.ExternalUserid, v.MsgId, v.SendTime)
		} else if v.Origin == 4 && v.MsgType == "event" {
			zap.L().Error("接到企业微信推送事件消息，v.Event", zap.Any("v.event", v.Event))
			// 系统推送的事件消息
			switch v.Event.EventType {
			case "enter_session":
				// 用户进入会话事件，将用户放入正在接待的客户队列，队列名称 InConversationExternalUserid
				// 当用户点击链接跳转进来时候，会优先进入此接口，此时会话状态为4
				result2, err2 := global.REDIS.LRange(ctx, "InConversationExternalUserid", 0, -1).Result()
				if err2 != nil {
					zap.L().Error("Redis LRange InConversationExternalUserid Error！", zap.Error(err1))
				}
				// 判断进入的用户是否已经在会话队列中
				if FindList(result2, v.Event.ExternalUserid) == false {
					// 判断会话状态，不为3的不准添加
					bodyData1 := map[string]interface{}{
						"open_kfid":       kfID,
						"external_userid": v.Event.ExternalUserid,
					}

					body1, err3 := json.Marshal(bodyData1)
					if err3 != nil {
						zap.L().Error("Json序列化失败！", zap.Error(err3))
					}

					tempData1, err4 := HttpClient("https://qyapi.weixin.qq.com/cgi-bin/kf/service_state/get?access_token="+token, "POST", string(body1), "")
					if err4 != nil {
						zap.L().Error("HTTP请求发送失败！", zap.Error(err4))
					}
					var ReturnDataTemp GetServiceStateStruct
					err = json.Unmarshal(tempData1, &ReturnDataTemp)
					if err != nil {
						zap.L().Error("Json序列化失败！", zap.Error(err))
					}
					if ReturnDataTemp.ServiceState == 3 {
						global.REDIS.RPush(ctx, "InConversationExternalUserid", v.Event.ExternalUserid)
					}
				}

				if v.Event.WelcomeCode != "" {
					// 条件为：用户在过去48小时里未收过欢迎语，且未向客服发过消息），会返回该字段。
					go SendMsgOnEvent(v.Event.WelcomeCode, "哈罗，尊敬的客户您好，很高兴为您服务，请问有什么可以帮您？")
				}

			case "msg_send_fail":
				// 消息发送失败事件, 记录MongoDB
				FailTypeMessageMap := map[int]string{
					0:  "未知原因",
					1:  "客服帐号已删除",
					2:  "应用已关闭",
					4:  "会话已过期，超过48小时",
					5:  "会话已关闭",
					6:  "超过5条限制",
					7:  "未绑定视频号",
					8:  "主体未验证",
					9:  "未绑定视频号且主体未验证",
					10: "用户拒收",
				}
				zap.L().Error("接收到微信事件发送消息失败", zap.String("fail_type", FailTypeMessageMap[v.Event.FailType]))
				go UpdateMongoMsgSendFail("workWechat", "userMsgSendFail", v.Event.ExternalUserid, v.Event.OpenKfid, v.Event.FailMsgId, FailTypeMessageMap[v.Event.FailType])

			case "servicer_status_change":
				// 1-接待中 2-停止接待
				// 接待人员接待状态变更事件 如从接待中变更为停止接待, 此处判定，接待人员是否合法，若状态为在线则将接待人员放入Redis队列中，若离线则从队列中删除
				if v.Event.Status == 1 {
					// 将客服加入接待队列尾部 ServiceUseridUpQueue
					if v.Event.ServiceUserid == "e0c48f2ff5b73bc696c2aa1a0f666f3e" {
						global.REDIS.RPush(ctx, "ServiceUseridUpQueue", v.Event.ServiceUserid) // 正常接待
					} else {
						global.REDIS.LPushX(ctx, "ServiceUseridUpQueue", v.Event.ServiceUserid) // 谁去洗手间回来优先接待
					}

				} else if v.Event.Status == 2 {
					// 先判断客服是否存在于队列中，若存在则将客服移除接待队列，不再进行会话分配
					strings, err2 := global.REDIS.LRange(ctx, "ServiceUseridUpQueue", 0, -1).Result()
					if err2 != nil {
						zap.L().Error("Redis LRange ServiceUseridUpQueue Error！", zap.Error(err2))
					}
					if FindList(strings, v.Event.ServiceUserid) {
						global.REDIS.LRem(ctx, "ServiceUseridUpQueue", 0, v.Event.ServiceUserid)
					}
				}
				// 记录接待人员接待状态变更
				go UpdateMongoServiceStatusChange("workWechat", "userServiceStatusChange", v.Event.ServiceUserid, v.Event.OpenKfid, v.Event.Status)

			case "session_status_change":
				/*
					会话状态变更事件:
						1-从接待池接入会话
						2-转接会话
						3-结束会话
						4-重新接入已结束/已转接会话
				*/

				switch v.Event.ChangeType {
				case 1:
					// 从接待池接入会话时，将客户ID保存一份到正在会话的客户队列 队列名称 InConversationExternalUserid
					result2, err2 := global.REDIS.LRange(ctx, "InConversationExternalUserid", 0, -1).Result()
					if err2 != nil {
						zap.L().Error("Redis LRange InConversationExternalUserid Error！", zap.Error(err1))
					}

					if FindList(result2, v.Event.ExternalUserid) == false {
						global.REDIS.RPush(ctx, "InConversationExternalUserid", v.Event.ExternalUserid)
					}

					err = global.REDIS.ZIncrBy(ctx, "ServiceUseridData", 1, v.Event.NewServiceUserid).Err()
					if err != nil {
						zap.L().Error("增加客服接待人数失败，请排查原因~", zap.Error(err))
					}
					// 记录会话状态到MongoDB中，方便后期做会话超时处理
					ServiceUserHistoryMongo := global.MONGO.Database("workWechat").Collection("ServiceUserHistory")
					data := ServiceUserHistoryStruct{
						ServiceUserid:  v.Event.NewServiceUserid,
						OpenKfId:       kfID,
						ExternalUserid: v.Event.ExternalUserid,
						ServiceStatus:  1,
						ServiceData:    time.Now().Format("2006-01-02 15:04:05"),
					}

					insertResult, err := ServiceUserHistoryMongo.InsertOne(context.TODO(), data)
					if err != nil {
						zap.L().Error("插入数据失败", zap.Error(err))
					}
					zap.L().Info("插入数据成功", zap.Any("data", insertResult.InsertedID))

					if v.Event.MsgCode != "" {
						go SendMsgOnEvent(v.Event.WelcomeCode, "请提供设备号，姓名，手机号")
					}

				case 2:
					// 记录转接会话操作
					ServiceUserHistoryMongo := global.MONGO.Database("workWechat").Collection("ServiceUserTransferHistory")
					data := ServiceUserTransferHistory{
						OldServiceUserid: v.Event.OldServiceUserid,
						OpenKfId:         kfID,
						ExternalUserid:   v.Event.ExternalUserid,
						NewServiceUserid: v.Event.NewServiceUserid,
						ServiceData:      time.Now().Format("2006-01-02 15:04:05"),
					}

					insertResult, err := ServiceUserHistoryMongo.InsertOne(context.TODO(), data)
					if err != nil {
						zap.L().Error("记录转接会话插入数据失败", zap.Error(err))
					}
					zap.L().Info("记录转接会话插入数据成功", zap.Any("data", insertResult.InsertedID))

				case 3:
					// 当客户会话结束时，从正在会话的客户队列中删除记录 队列名称 InConversationExternalUserid
					result2, err2 := global.REDIS.LRange(ctx, "InConversationExternalUserid", 0, -1).Result()
					if err2 != nil {
						zap.L().Error("Redis LRange InConversationExternalUserid Error！", zap.Error(err1))
					}
					if FindList(result2, v.Event.ExternalUserid) {
						global.REDIS.LRem(ctx, "InConversationExternalUserid", 0, v.Event.ExternalUserid)
					}

					zap.L().Debug("本次会话已结束", zap.String("ExternalUserid", v.Event.ExternalUserid))

					err = global.REDIS.ZIncrBy(ctx, "ServiceUseridData", -1, v.Event.OldServiceUserid).Err()
					if err != nil {
						zap.L().Error("减少客服接待人数失败，请排查原因~", zap.Error(err))
					}

					// 从MongoDB中删除正在会话的信息
					ServiceUserHistoryMongo := global.MONGO.Database("workWechat").Collection("ServiceUserHistory")
					one, errDel := ServiceUserHistoryMongo.DeleteOne(context.TODO(), bson.D{{"externalUserid", v.Event.ExternalUserid}})
					if errDel != nil {
						zap.L().Error("删除MongoDB数据失败！", zap.Error(errDel), zap.String("data", strconv.FormatInt(one.DeletedCount, 10)))
					}
					if v.Event.MsgCode != "" {
						go SendMsgOnEvent(v.Event.WelcomeCode, "本次会话已结束，感谢您的来访~")
					}
				case 4:
					zap.L().Info("重新接入已结束/已转接会话", zap.Any("data", v))
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

// ControlMessage 将客户信息调度给客服
func ControlMessage(OpenKfId string) {
	/*
		OpenKfId： 客服账号ID
		kfID：     客服ID
		userID:    客户ID
	*/

	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	start, _ := time.ParseInLocation("2006-01-02 15:04:05", now.Format("2006-01-02")+" 09:00:00", loc)
	end, _ := time.ParseInLocation("2006-01-02 15:04:05", now.Format("2006-01-02")+" 21:30:00", loc)
	ExternalUserid := ""
	if now.After(start) && now.Before(end) {
		/*
			将当前客服正在接待的客户放入MongoDB
			{
				ServiceUserid string `json:"servicer_userid"`			        // 客服ID
				OpenKfId       string `json:"open_kfid,omitempty"`			    // 客服账号
				ExternalUseridList string `json:"external_userid,omitempty"`    // 会话客户
				ServiceStatus  int `json:"servicer_userid"`
			}
		*/

		// 将客户分配给客服 判断接待池中是否有客户，若有则优先分配 ExternalUseridUpQueue
		ExternalUserid, _ = global.REDIS.LPop(context.Background(), "ExternalUseridUpQueue").Result()
		fmt.Println("ExternalUserid", ExternalUserid)
		if ExternalUserid == "" {
			// 优先队列没有人的时候，去MQ中取一个 ExternalUseridQueue
			queue, err := PullQueue("ExternalUseridQueue")
			if err != nil {
				zap.L().Error("从ExternalUseridQueue队列中取值失败！", zap.Error(err))
			}
			ExternalUserid = string(queue)
			fmt.Println("从ExternalUseridQueue队列中取值 ExternalUserid", ExternalUserid)
			if ExternalUserid != "" {
				TransServiceStateData := TransServiceState(OpenKfId, ExternalUserid)
				zap.L().Info("分配会话状态返回数据：", zap.Any("data", TransServiceStateData))
			}
		}
	} else {
		// 推送客服不在线的消息
		sendDataReturn := SendMsgData(ExternalUserid, OpenKfId, "您好，非常感谢您的咨询~ \n"+
			"很抱歉告知您，此时是我们的非工作时间，我们的客服人员已下班。如果您有任何问题或需求，请您在工作时间(09:00~21:30)再次联系我们，我们将竭诚为您提供满意的服务。感谢您的理解和支持，祝您生活愉快！")
		zap.L().Info("发送普通消息返回数据：", zap.Any("data", sendDataReturn))
		// 将会话状态重置为4
		TransServiceStateData := TransServiceState(OpenKfId, ExternalUserid)
		zap.L().Info("分配会话状态返回数据：", zap.Any("data", TransServiceStateData))
	}
}
