package utils

// ReturnCodeStruct 通用返回信息
type ReturnCodeStruct []struct {
	Msg  string                   `json:"msg"`  // 描述信息
	Code int                      `json:"code"` // 状态码
	Data []map[string]interface{} `json:"data"` // 返回数据
}

type ServiceUserHistoryStruct struct {
	ServiceUserid  string `bson:"serviceUserid"`  // 接待人员ID
	OpenKfId       string `bson:"openKfId"`       // 客服账号
	ExternalUserid string `bson:"externalUserid"` // 用户ID
	ServiceStatus  int    `bson:"serviceStatus"`  // 会话状态 1 接入中  0 会话结束 2 会话超时结束
	ServiceData    string `bson:"serviceData"`    // 会话接入时间
}

type ServiceUserTransferHistory struct {
	OldServiceUserid string `bson:"oldServiceUserid"` // 旧接待人员ID
	NewServiceUserid string `bson:"newServiceUserid"` // 新接待人员ID
	OpenKfId         string `bson:"openKfId"`         // 客服账号
	ExternalUserid   string `bson:"externalUserid"`   // 用户ID
	ServiceData      string `bson:"serviceData"`      // 会话接入时间
}
