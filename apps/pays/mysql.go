package pays

import "gorm.io/gorm"

// PayConfigData 存放支付的相关配置信息
type PayConfigData struct {
	gorm.Model
	Name         string `gorm:"unique"`
	Mobile       string
	Email        string
	PrivateKey   string
	PublicKey    string
	AppID        string
	MchID        string
	PayNotify    string
	RefundNotify string
	Secret       string
	Status       bool
	IsDelete     bool
	AgentID      string
}

type AliPayData struct {
	gorm.Model
	Device      string // 设备号
	PackageID   int    // 套餐ID
	TotalAmount string // 支付金额
	TimeStamp   string // 支付时间
	OutTradeNo  string
	TradeNo     string
	SellerId    string
}

type PayData struct {
	gorm.Model
	DeviceId            string // 设备ID
	Mobile              string // 手机号
	PackageId           string // 套餐ID
	PackageName         string // 套餐名称name
	PackageMoney        string // 套餐真实价格price
	PackageVirtualMoney string // 套餐虚拟价格fictitiousPrice
	PackageType         string // 套餐类型SetMealType_id
	PackageDay          string // 套餐有效天数day
	PackageTotalFlow    string // 套餐总流量totalFlow
	AgentID             string // 代理商ID
	GroupID             string // 设备组ID
	WeChatAppID         string // 商户号
	WeChatMchID         string // 商户号
	Money               string // 支付金额
	OpenId              string // 微信支付客户ID
	OutTradeNo          string // 订单号
	Body                string // 订单描述
	Status              int    // 是否支付成功(回调地址更新， 1：未支付，2：已支付，3：已退款)
	PayType             string // 支付类型(支付宝，微信)
	TimeEnd             string // 支付完成时间
	IsDelete            bool
}
