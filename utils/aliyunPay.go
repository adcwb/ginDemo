package utils

import (
	"github.com/smartwalle/alipay/v2"
	"go.uber.org/zap"
)

const (
	// AppId 沙箱环境
	AppId = ""
	// PrivateKey 应用私钥
	PrivateKey = ""
	// 支付宝公钥
	aliPublicKey = ""
)

func AliPayClientInitV2() (aliClient *alipay.Client, err error) {
	aliClient, err = alipay.New(AppId, aliPublicKey, PrivateKey, true)

	if err != nil {
		zap.L().Error("初始化支付宝失败！", zap.Error(err))
	}

	return
}
