package wechat

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

type GetWorkUserDataStruct struct {
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	CustomerList []struct {
		ExternalUserid      string `json:"external_userid"`
		Nickname            string `json:"nickname"`
		Avatar              string `json:"avatar"`
		Gender              int    `json:"gender"`
		Unionid             string `json:"unionid"`
		EnterSessionContext struct {
			Scene          string `json:"scene"`
			SceneParam     string `json:"scene_param"`
			WechatChannels struct {
				Nickname string `json:"nickname"`
				Scene    int    `json:"scene"`
			} `json:"wechat_channels"`
		} `json:"enter_session_context"`
	} `json:"customer_list"`
	InvalidExternalUserid []string `json:"invalid_external_userid"`
}

type UserDataStruct struct {
	UserID       string `json:"userID"`
	Username     string `json:"username"`
	Mobile       string `json:"mobile"`
	DeviceNumber string `json:"deviceNumber"`
	DeviceModel  string `json:"deviceModel"`
	IccID        string `json:"IccID"`
	Operator     string `json:"operator"`
	Address      string `json:"address"`
	Comment      string `json:"comment"`
}

type MongoDBUserDataStruct struct {
	Status              string `json:"status" bson:"status"`
	UserID              string `json:"userID"`
	Nickname            string `json:"nickname"`
	Avatar              string `json:"avatar"`
	Gender              int    `json:"gender"`
	Unionid             string `json:"unionid"`
	Username            string `json:"username"`
	Mobile              string `json:"mobile"`
	DeviceNumber        string `json:"deviceNumber"`
	DeviceModel         string `json:"deviceModel"`
	IccID               string `json:"IccID"`
	Operator            string `json:"operator"`
	Address             string `json:"address"`
	Comment             string `json:"comment"`
	EnterSessionContext struct {
		Scene          string `json:"scene"`
		SceneParam     string `json:"scene_param"`
		WechatChannels struct {
			Nickname string `json:"nickname"`
			Scene    int    `json:"scene"`
		} `json:"wechat_channels"`
	} `json:"enter_session_context"`
}
