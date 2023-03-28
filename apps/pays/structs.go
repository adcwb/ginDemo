package pays

import xmljson "encoding/xml"

type testStruct struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data string `json:"data"`
}

type test2Struct struct {
	Data []string `json:"data"`
	Code int      `json:"code"`
}

type PackageInfoStruct struct {
	Data struct {
		TotalFlow    string `json:"total_flow"`
		PackageName  string `json:"package_name"`
		PackageDays  string `json:"package_days"`
		PackagePrice string `json:"package_price"`
		EndTime      string `json:"end_time"`
	} `json:"data"`
	Code int `json:"code"`
}

type KeysData struct {
	Name         string `json:"name"`        // 客户名称
	Mobile       string `json:"mobile"`      // 客户手机号
	Email        string `json:"email"`       // 客户邮箱
	PrivateKey   string `json:"private_key"` // PrivateKey
	PublicKey    string `json:"public_key"`  // PublicKey
	Status       bool   `json:"status"`      // 状态信息
	AppID        string `json:"app_id"`
	MchID        string `json:"mch_id"`
	PayNotify    string `json:"pay_notify"`
	RefundNotify string `json:"refund_notify"`
	Secret       string `json:"secret"`
}

type DelKeysStruct struct {
	ID uint `json:"id"`
}

type PayCallback struct {
	Data struct {
		AlipayTradeAppPayResponse struct {
			Code        string `json:"code"`
			Msg         string `json:"msg"`
			AppId       string `json:"app_id"`
			AuthAppId   string `json:"auth_app_id"`
			Charset     string `json:"charset"`
			Timestamp   string `json:"timestamp"`
			OutTradeNo  string `json:"out_trade_no"`
			TotalAmount string `json:"total_amount"`
			TradeNo     string `json:"trade_no"`
			SellerId    string `json:"seller_id"`
		} `json:"alipay_trade_app_pay_response"`
		Sign     string `json:"sign"`
		SignType string `json:"sign_type"`
	} `json:"data"`
	Imei      string `json:"imei"`
	PackageId string `json:"package_id"`
}

type WechatStructData struct {
	Appid          string `xml:"appid"`            // 应用ID
	MchId          string `xml:"mch_id"`           // 商户号
	NonceStr       string `xml:"nonce_str"`        // 随机字符串
	Sign           string `xml:"sign"`             // 签名
	Body           string `xml:"body"`             // 商品描述 应用市场上的APP名字-实际商品名称
	OutTradeNo     string `xml:"out_trade_no"`     // 商户订单号
	TotalFee       string `xml:"total_fee"`        // 总金额
	SpbillCreateIp string `xml:"spbill_create_ip"` // 终端IP
	NotifyUrl      string `xml:"notify_url"`       // 通知地址
	TradeType      string `xml:"trade_type"`       // 交易类型 APP
}

type XmlData struct {
	XMLName    xmljson.Name `xml:"xml"`
	ReturnCode string       `xml:"return_code"`
	ReturnMsg  string       `xml:"return_msg"`
	ResultCode string       `xml:"result_code"`
	MchId      string       `xml:"mch_id"`
	Appid      string       `xml:"appid"`
	NonceStr   string       `xml:"nonce_str"`
	Sign       string       `xml:"sign"`
	PrepayId   string       `xml:"prepay_id"`
	TradeType  string       `xml:"trade_type"`
}

type UserInfo struct {
	Money     float64 `json:"money"`
	Body      string  `json:"body"`
	Imei      string  `json:"imei"`
	Mobile    string  `json:"mobile"`
	PackageID string  `json:"package_id"`
}

type UserInfoH5 struct {
	Terrace   string  `json:"terrace"`
	Money     float64 `json:"money"`
	Body      string  `json:"body"`
	Imei      string  `json:"imei"`
	Mobile    string  `json:"mobile"`
	PackageId int     `json:"package_id"`
	Openid    string  `json:"openid"`
}

type PackageIdData struct {
	Data struct {
		SetMealBaseName       interface{} `json:"setMealBaseName"`
		Name                  string      `json:"name"`
		Price                 string      `json:"price"`
		FictitiousPrice       string      `json:"fictitiousPrice"`
		Day                   string      `json:"day"`
		IsWholeMonth          bool        `json:"is_wholeMonth"`
		TotalFlow             string      `json:"totalFlow"`
		FictitiousFlowIndexID interface{} `json:"fictitiousFlowIndexID"`
		SpeedFlowIndexID      interface{} `json:"SpeedFlowIndexID"`
		Integral              string      `json:"integral"`
		PurchaseType          string      `json:"purchaseType"`
		ExcessAction          string      `json:"excess_action"`
		ExcessPrice           interface{} `json:"excess_price"`
		CostPrice             string      `json:"cost_price"`
		SetMealTypeId         string      `json:"SetMealType_id"`
		AgentID               string      `json:"agent_id"`
		GroupID               string      `json:"group_id"`
	} `json:"data"`
	Code int `json:"code"`
}

type WechatPayCallbackStruct struct {
	XMLName       xmljson.Name `xml:"xml"`
	Appid         string       `xml:"appid"`
	BankType      string       `xml:"bank_type"`
	CashFee       string       `xml:"cash_fee"`
	FeeType       string       `xml:"fee_type"`
	IsSubscribe   string       `xml:"is_subscribe"`
	MchId         string       `xml:"mch_id"`
	NonceStr      string       `xml:"nonce_str"`
	Openid        string       `xml:"openid"`
	OutTradeNo    string       `xml:"out_trade_no"`
	ResultCode    string       `xml:"result_code"`
	ReturnCode    string       `xml:"return_code"`
	Sign          string       `xml:"sign"`
	TimeEnd       string       `xml:"time_end"`
	TotalFee      string       `xml:"total_fee"`
	TradeType     string       `xml:"trade_type"`
	TransactionId string       `xml:"transaction_id"`
}

type WechatPayRefundStruct struct {
	XMLName     xmljson.Name `xml:"xml"`
	Appid       string       `xml:"appid"`         // 公众账号ID
	MchId       string       `xml:"mch_id"`        // 商户号
	NonceStr    string       `xml:"nonce_str"`     // 随机字符串
	OutTradeNo  string       `xml:"out_trade_no"`  // 商户订单号
	OutRefundNo string       `xml:"out_refund_no"` // 商户退款单号
	TotalFee    int          `xml:"total_fee"`     // 订单金额
	RefundFee   int          `xml:"refund_fee"`    // 退款金额
	Sign        string       `xml:"sign"`          // 签名
}

type RefundMoneyStruct struct {
	OutTradeNo  string  `json:"out_trade_no"`
	RefundMoney float64 `json:"refund_money"`
}

type OpenTradeRefundMoneyStruct struct {
	OutTradeNo  string `json:"out_trade_no"`
	RefundMoney string `json:"refund_money"`
}

// WechatPayReturnDataRefundStruct 微信支付退款返回数据
type WechatPayReturnDataRefundStruct struct {
	XMLName           xmljson.Name `xml:"xml"`
	ReturnCode        string       `xml:"return_code"`         // 返回状态码
	ReturnMsg         string       `xml:"return_msg"`          // 返回信息
	ResultCode        string       `xml:"result_code"`         // 业务结果
	Appid             string       `xml:"appid"`               // 公众账号ID
	MchId             string       `xml:"mch_id"`              // 商户号
	NonceStr          string       `xml:"nonce_str"`           // 随机字符串
	TransactionId     string       `xml:"transaction_id"`      // 微信支付订单号
	OutTradeNo        string       `xml:"out_trade_no"`        // 商户订单号
	OutRefundNo       string       `xml:"out_refund_no"`       // 商户退款单号
	RefundId          string       `xml:"refund_id"`           // 微信退款单号
	TotalFee          int          `xml:"total_fee"`           // 订单金额
	RefundFee         int          `xml:"refund_fee"`          // 退款金额
	CashFee           int          `xml:"cash_fee"`            // 现金支付金额
	RefundChannel     string       `xml:"refund_channel"`      //
	CouponRefundFee   string       `xml:"coupon_refund_fee"`   // 代金券退款总金额
	CouponRefundCount string       `xml:"coupon_refund_count"` // 退款代金券使用数量
	CashRefundFee     string       `xml:"cash_refund_fee"`     // 现金退款金额
	Sign              string       `xml:"sign"`                // 签名
}

type WechatPayCallbackReturnStruct struct {
	XMLName    xmljson.Name `xml:"xml"`
	ReturnCode string       `xml:"return_code"`
	ReturnMsg  string       `xml:"return_msg"`
}

type TokenStruct struct {
	UserAgentCard string `json:"user_agent_card"`
	UserAgentDev  string `json:"user_agent_dev"`
	LoginTime     string `json:"login_time"`
	UserId        int    `json:"user_id"`
	IsRoot        bool   `json:"is_root"`
}

type SetMealSurplusDaysStruct struct {
	Code int `json:"code"`
	Data struct {
		EndDays        string  `json:"end_days"`
		CreatedTime    string  `json:"created_time"`
		FlowTotal      float64 `json:"flow_total"`
		TotalFlowCount float64 `json:"TotalFlowCount"`
		ProductName    string  `json:"product_name"`
		PkgName        string  `json:"pkg_name"`
		RestDay        int     `json:"rest_day"`
		RestData       int     `json:"rest_data"`
		TotalDays      int     `json:"total_days"`
	} `json:"data"`
}
