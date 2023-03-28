package utils

import (
	"ginDemo/global"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"go.uber.org/zap"
)

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */

// CreateClient 初始化阿里云发送短信客户端
func CreateClient(accessKeyId *string, accessKeySecret *string) (result *dysmsapi20170525.Client, err error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	result = &dysmsapi20170525.Client{}
	result, err = dysmsapi20170525.NewClient(config)
	return result, err
}

func SendSmSCode(sendSmsRequest *dysmsapi20170525.SendSmsRequest) (err error) {

	client, err := CreateClient(tea.String(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".aliYunAccessKeyId")), tea.String(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".aliYunAccessKeySecret")))
	if err != nil {
		zap.L().Error("发送短信模块初始化失败，请检查原因！", zap.Error(err))
	}

	ResponseData, err := client.SendSms(sendSmsRequest)
	message := *ResponseData.Body.Message

	zap.L().Info("短信发送返回信息：", zap.String("data", message))
	if err != nil {
		zap.L().Error("短信发送失败！", zap.Error(err))
	}
	return err
}
