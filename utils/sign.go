package utils

import "ginDemo/global"

/*
	此组件保存各类运营商签名接口
*/

type OperatorKey struct {
	AppKey    string `json:"AppKey"`
	SecretKey string `json:"secretKey"`
}

// OperatorSign 运营商获取签名接口
type OperatorSign interface {
	// Sign 获取签名
	Sign(params, body, timestamp string) (sign string)
}

// YanChengCTCC 江苏盐城电信
type YanChengCTCC struct{}

// YangZhouCTCC 江苏扬州电信
type YangZhouCTCC struct{}

// Sign YanChengCTCC 接口实现对应的方法
func (yc YanChengCTCC) Sign(params, body, timestamp string) (sign string) {
	var tempData OperatorKey
	tempData.AppKey = global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".YanChengChinaTelecomAppKey")
	tempData.SecretKey = global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".YanChengChinaTelecomSecretKey")
	data := params + body + tempData.SecretKey + timestamp
	sign = ToHexStr(MD5(data))
	return sign
}

// Sign YangZhouCTCC 接口实现对应的方法
func (yc YangZhouCTCC) Sign(params, body, timestamp string) (sign string) {
	var tempData OperatorKey
	tempData.AppKey = global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".YangZhouChinaTelecomAppKey")
	tempData.SecretKey = global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".YangZhouChinaTelecomSecretKey")
	data := params + tempData.SecretKey + timestamp
	sign = ToHexStr(MD5(data))
	return sign
}
