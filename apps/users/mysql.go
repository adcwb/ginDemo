package users

import (
	"gorm.io/gorm"
)

/*
权限控制
`gorm:"<-:create"`		允许读取和创建
`gorm:"<-:update"`		允许读取和更新
`gorm:"<-"`				允许读写（创建和更新）
`gorm:"<-:false"`		允许读，禁用写权限
`gorm:"->"`				只读（禁用写权限，除非它已配置）
`gorm:"->;<-:create"`	允许读取和创建
`gorm:"->:false;<-:create"`
`gorm:"-"`				使用 struct 读写时忽略此字段
`gorm:"-:all"`			读写时忽略此字段migrate with struct
 `gorm:"-:migration"`	使用 struct 迁移时忽略该字段
*/

// User 用户信息
type User struct {
	gorm.Model
	Name   string
	Email  *string
	Age    uint8
	Mobile string `gorm:"index"`
	Status bool
	OpenID string `gorm:"index"`
}

// SendSmsCodeData 发送短信记录
type SendSmsCodeData struct {
	gorm.Model
	Mobile string
	Code   string
}

// UserRechargeHistory 用户金额修改记录
type UserRechargeHistory struct {
	gorm.Model
	UserID    int     // 用户ID
	Mobile    string  // 用户手机号
	UserMoney float64 // 未修改前的金额
	Money     float64 // 操作的金额
	Mode      int     // 1 充值  2  扣费
	DeviceID  string  // 当前时间用户手机号绑定的设备号，若没有则不填
	AgentID   int     // 当前手机号是否是代理
	Operator  string  // 操作者
}

// DevicePreCharge 套餐预存表
type DevicePreCharge struct {
	gorm.Model
	UserID       int    // 用户ID
	DeviceID     string // 当前时间用户手机号绑定的设备号，若没有则不填
	TotalFlow    string // 套餐总流量，单位M
	PackageID    string // 套餐ID
	PackageName  string // 套餐名称
	PackageDays  string // 套餐可用天数
	PackagePrice string // 套餐价格
	Operator     string // 操作者
	Status       int    // 是否生效，  0：未生效
}
