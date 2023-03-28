package pays

import (
	"crypto/tls"
	"encoding/json"
	"encoding/pem"
	xmljson "encoding/xml"
	"fmt"
	"ginDemo/apps/users"
	"ginDemo/global"
	"ginDemo/utils"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v2"
	"github.com/smartwalle/xid"
	WXPay "github.com/wleven/wxpay"
	"github.com/wleven/wxpay/src/V2"
	"github.com/wleven/wxpay/src/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// GetPackageData 支付时获取套餐ID对应的详细信息，并保存到数据库
func GetPackageData(ID, deviceId, tradeNo, channel, mobile, money, body string) {

	data, err := utils.HttpClient(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".DeviceServer")+"/devices/PayOkGetPackage?package_id="+ID+"&device_id="+deviceId, "GET", "", "")
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}

	var ReturnData PackageIdData
	err = json.Unmarshal(data, &ReturnData)
	if err != nil {
		zap.L().Error("返回值序列化失败！", zap.Error(err))
	}
	// 根据设备号去后台查询出该设备属于哪个代理，并将信息保存

	global.DB.Table("pay_data").Where("device_id = ?", deviceId).Find(&PayData{})

	// 批量插入数据
	var payDataTemp = []PayData{
		{
			PackageId:           ID,
			PackageName:         ReturnData.Data.Name,
			PackageMoney:        ReturnData.Data.Price,
			PackageVirtualMoney: ReturnData.Data.FictitiousPrice,
			PackageType:         ReturnData.Data.PurchaseType,
			PackageDay:          ReturnData.Data.Day,
			PackageTotalFlow:    ReturnData.Data.TotalFlow,
			DeviceId:            deviceId,
			AgentID:             ReturnData.Data.AgentID,
			GroupID:             ReturnData.Data.GroupID,
			Mobile:              mobile,
			Money:               money,
			OpenId:              "",
			OutTradeNo:          tradeNo,
			Body:                body,
			Status:              1,
			PayType:             channel,
			TimeEnd:             "",
		},
	}
	global.DB.Create(&payDataTemp)

}

// SetCallbackData 回调成功时补充数据库参数
func SetCallbackData(tradeNo string) {
	result := global.DB.Table("pay_data").Where("out_trade_no = ?", tradeNo).First(&PayData{}).Updates(PayData{
		Status: 2,
	})
	if result.Error != nil {
		zap.L().Error("数据库查询失败！", zap.Error(result.Error))
	}
}

// MobileAppWapPay 支付宝移动端APP支付
func MobileAppWapPay(c *gin.Context) {
	b, _ := c.GetRawData() // 从c.Request.Body读取请求数据
	fmt.Println(string(b))
	var payConfig UserInfo
	err := json.Unmarshal(b, &payConfig)
	if err != nil {
		zap.L().Error("json序列化失败", zap.Error(err))
	}

	// 生成商家订单号
	var tradeNo = fmt.Sprintf("%d", xid.Next())
	// 初始化订单
	var p = alipay.TradeAppPay{}
	// 异步通知
	p.NotifyURL = "https://pays.example.com/pays" + "/notify"
	// 回调地址
	p.ReturnURL = "https://pays.example.com/pays" + "/callback"
	// 商品说明
	p.Subject = "流量充值:" + payConfig.Imei
	// 商家订单号
	p.OutTradeNo = tradeNo
	// 订单金额
	p.TotalAmount = fmt.Sprintf("%.2f", payConfig.Money)

	// 订单描述
	p.Body = payConfig.Body
	p.TimeoutExpress = "90m"

	// 初始化Alipay
	aliClient, err := utils.AliPayClientInitV2()
	if err != nil {
		zap.L().Error("AliPayClientInit初始化失败！", zap.Error(err))
	}
	url, _ := aliClient.TradeAppPay(p)

	// 将数据写入到数据库中
	GetPackageData(payConfig.PackageID, payConfig.Imei, tradeNo, "支付宝", payConfig.Mobile, fmt.Sprintf("%.2f", payConfig.Money), payConfig.Body)

	c.String(http.StatusOK, url)

}

// WechatPay 微信支付接口
func WechatPay(c *gin.Context) {
	b, _ := c.GetRawData() // 从c.Request.Body读取请求数据

	var UserConfig UserInfo
	err := json.Unmarshal(b, &UserConfig)
	if err != nil {
		zap.L().Error("json序列化失败", zap.Error(err))
	}

	config := entity.PayConfig{
		AppID:        "",
		MchID:        "1610519787",
		PayNotify:    "https://pays.example.com/pays/WechatPayCallback",
		RefundNotify: "https://pays.example.com/pays/WechatPayCallback",
		Secret:       "", // key
	}

	// 获取用户请求的地址，若获取不到，则使用本地的第一个地址
	clientIP := c.ClientIP()
	if clientIP == "" {
		clientIP = utils.GetLocalIP()[0]
	}

	wxpay := WXPay.Init(config)
	var tempData V2.UnifiedOrder
	tradeNo := fmt.Sprintf("%d", xid.Next())
	tempData.Body = "天朝优度-流量充值-" + UserConfig.Imei
	tempData.OutTradeNo = tradeNo
	tempData.NotifyURL = "https://pays.example.com/pays/WechatPayCallback"
	tempData.SpbillCreateIP = clientIP
	tempData.TotalFee = int(UserConfig.Money * 100)
	tempData.TradeType = "APP"

	// 统一下单地址
	data, err := wxpay.V2.WxAppAppPay(tempData)
	if err != nil {
		zap.L().Error("下单失败", zap.Error(err))
	}
	tempSign := data["paySign"]
	delete(data, "paySign")
	data["sign"] = tempSign

	// 将数据写入到数据库中
	GetPackageData(UserConfig.PackageID, UserConfig.Imei, tradeNo, "微信", UserConfig.Mobile, fmt.Sprintf("%.2f", UserConfig.Money), UserConfig.Body)

	c.JSON(http.StatusOK, data)
}

// WechatPayH5 微信支付接口
func WechatPayH5(c *gin.Context) {
	b, _ := c.GetRawData() // 从c.Request.Body读取请求数据
	//fmt.Println(string(b))
	var UserConfig UserInfoH5
	err := json.Unmarshal(b, &UserConfig)
	if err != nil {
		zap.L().Error("json序列化失败------------------", zap.Error(err))
	}
	zap.L().Info(fmt.Sprintf("发起支付请求---------------- %s", string(b)))

	// 通过前端传入的平台信息，去数据库查询对应的支付配置信息，暂时不启用
	var configTemp PayConfigData
	global.DB.Table("pay_config_data").Where("name = ?", UserConfig.Terrace).First(&configTemp)

	config := entity.PayConfig{
		AppID:        "",
		MchID:        "",
		PayNotify:    "https://pays.example.com/pays/WechatPayCallback",
		RefundNotify: "https://pays.example.com/pays/WechatPayCallback",
		Secret:       "", // key
	}

	// 获取用户请求的地址，若获取不到，则使用本地的第一个地址
	clientIP := c.ClientIP()
	if clientIP == "" {
		clientIP = utils.GetLocalIP()[0]
	}

	wxpay := WXPay.Init(config)
	tradeNo := fmt.Sprintf("%d", xid.Next())

	// 微信JSAPI下单地址
	data, err := wxpay.V2.WxAppPay(V2.UnifiedOrder{
		Body:           "-" + UserConfig.Body + "-" + UserConfig.Imei,
		OutTradeNo:     tradeNo,
		NotifyURL:      "https://pays.example.com/pays/WechatPayCallback",
		SpbillCreateIP: clientIP,
		TotalFee:       int(UserConfig.Money * 100),
		OpenID:         UserConfig.Openid,
	})
	if err != nil {
		zap.L().Error("下单失败", zap.Error(err))
	}
	tempSign := data["paySign"]
	delete(data, "paySign")
	data["sign"] = tempSign

	// 将数据写入到数据库中
	GetPackageData(strconv.Itoa(UserConfig.PackageId), UserConfig.Imei, tradeNo, "微信", UserConfig.Mobile, fmt.Sprintf("%.2f", UserConfig.Money), UserConfig.Body)

	c.JSON(http.StatusOK, data)
}

// WechatPayRefund 微信退款接口
func WechatPayRefund(c *gin.Context) {
	b, _ := c.GetRawData() // 从c.Request.Body读取请求数据
	// 调用处传入退款商家订单号，退款金额

	var tempRefundMoneyStruct RefundMoneyStruct
	err := json.Unmarshal(b, &tempRefundMoneyStruct)
	if err != nil {
		zap.L().Error("json反序列化失败！", zap.Error(err))
	}
	var tempData PayData
	global.DB.Table("pay_data").Where("out_trade_no = ?", tempRefundMoneyStruct.OutTradeNo).First(&tempData)

	// 根据商家订单号，查询出设备的充值记录
	if tempData.Money == "" {
		ReturnData := map[string]interface{}{
			"data": "该订单不存在",
			"msg":  20001,
		}

		c.JSON(http.StatusOK, ReturnData)
		return
	}

	// 读取证书文件
	dateKey, _ := os.ReadFile("apps/pays/keys/ap.pem")
	pem.Decode(dateKey)
	var pemBlocks []*pem.Block
	var v *pem.Block
	var pkey []byte
	for {
		v, dateKey = pem.Decode(dateKey)
		if v == nil {
			break
		}
		if v.Type == "PRIVATE KEY" {
			pkey = pem.EncodeToMemory(v)
		} else {
			pemBlocks = append(pemBlocks, v)
		}
	}

	bytes := pem.EncodeToMemory(pemBlocks[0])
	//keyString := string(pkey)
	//CertString := string(bytes)
	//fmt.Printf("Cert :\n %s \n Key:\n %s \n ", CertString, keyString)
	//pool := x509.NewCertPool()
	dataTls, _ := tls.X509KeyPair(bytes, pkey)
	//pool.AppendCertsFromPEM(b)

	cfg := &tls.Config{
		Certificates: []tls.Certificate{dataTls},
	}

	tr := &http.Transport{
		TLSClientConfig: cfg,
	}
	client := &http.Client{Transport: tr}

	// 处理请求体参数
	var refund WechatPayRefundStruct

	refund.Appid = ""
	refund.MchId = ""
	refund.NonceStr = utils.GetRandomString(16)
	refund.OutTradeNo = tempRefundMoneyStruct.OutTradeNo
	refund.OutRefundNo = fmt.Sprintf("%d", xid.Next())
	refund.TotalFee = 1
	refund.RefundFee = int(tempRefundMoneyStruct.RefundMoney * 100)

	// 微信签名数据，手动拼接的
	tempSign := "appid=" + refund.Appid + "&mch_id=" + refund.MchId + "&nonce_str=" + refund.NonceStr + "&out_refund_no=" + refund.OutRefundNo + "&out_trade_no=" + refund.OutTradeNo + "&refund_fee=" + strconv.Itoa(refund.RefundFee) + "&total_fee=" + strconv.Itoa(refund.TotalFee) + "&key=IwAnt86M54idEToLef00kA29yGeintHe"

	refund.Sign = utils.MD5(tempSign)
	// 将数据xml化，然后装入body中准备发送
	marshal, err := xmljson.Marshal(refund)
	if err != nil {
		zap.L().Error("xml序列化失败！", zap.Error(err))
	}

	payload := strings.NewReader(string(marshal))
	request, _ := http.NewRequest("POST", "https://api.mch.weixin.qq.com/secapi/pay/refund", payload)
	request.Header.Set("Content-Type", "application/xml")

	// 接收返回数据
	resp, err := client.Do(request)
	// require.Nil(t, err)
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))

		ReturnData := map[string]interface{}{
			"data": "HttpClient请求发送失败, 请联系管理员",
			"msg":  20001,
		}

		c.JSON(http.StatusOK, ReturnData)

	} else {
		data, _ := io.ReadAll(resp.Body)
		zap.L().Info("退款成功", zap.String("ReturnData", string(data)))
		// 退款成功以后，更新该笔订单的状态，修改为已退款
		global.DB.Table("pay_data").Where("out_trade_no = ?", tempRefundMoneyStruct.OutTradeNo).Update("status", 3)

		// 返回数据给前端
		var ReturnData WechatPayReturnDataRefundStruct

		err = xmljson.Unmarshal(data, &ReturnData)
		if err != nil {
			zap.L().Error("返回值序列化失败！", zap.Error(err))
		}

		c.JSON(http.StatusOK, ReturnData)
	}

}

// AliCallback 支付宝支付回调处理
func AliCallback(c *gin.Context) {

	b, _ := c.GetRawData()
	fmt.Println(string(b))
	var data PayCallback

	err := json.Unmarshal(b, &data)
	if err != nil {
		zap.L().Error("返回值序列化失败！", zap.Error(err))
	}

	var p = alipay.TradeQuery{}
	p.OutTradeNo = data.Data.AlipayTradeAppPayResponse.OutTradeNo
	p.OutTradeNo = data.Data.AlipayTradeAppPayResponse.TotalAmount

	// 初始化Alipay
	aliClient, err := utils.AliPayClientInitV2()
	if err != nil {
		zap.L().Error("AliPayClientInit初始化失败！", zap.Error(err))
	}

	// 统一收单线下交易查询
	rsp, err := aliClient.TradeQuery(p)
	if err != nil {
		ReturnData := map[string]interface{}{
			"code": 20011,
			"msg":  "支付失败，请联系卖家进行核实！",
		}

		zap.L().Error("验证订单信息发生错误", zap.String("订单编号", data.Data.AlipayTradeAppPayResponse.OutTradeNo), zap.Error(err))
		c.JSON(http.StatusOK, ReturnData)
		return
	}
	if rsp.IsSuccess() == false {
		c.String(http.StatusBadRequest, "验证订单 %s 信息发生错误: %s-%s", data.Data.AlipayTradeAppPayResponse.OutTradeNo, rsp.Content.Msg, rsp.Content.SubMsg)
		ReturnData := map[string]interface{}{
			"code": 20011,
			"msg":  "支付失败，请联系卖家进行核实！",
		}

		zap.L().Error("验证订单信息发生错误", zap.String("订单编号", data.Data.AlipayTradeAppPayResponse.OutTradeNo), zap.String("msg", rsp.Content.Msg), zap.String("SubMsg", rsp.Content.SubMsg), zap.Error(err))
		c.JSON(http.StatusOK, ReturnData)
		return
	}

	SetCallbackData(data.Data.AlipayTradeAppPayResponse.OutTradeNo)
	// 将支付信息发送到后端

	RequestPayData := map[string]interface{}{
		"imei":         data.Imei,
		"package_id":   data.PackageId,
		"out_trade_no": data.Data.AlipayTradeAppPayResponse.OutTradeNo,
		"total_amount": data.Data.AlipayTradeAppPayResponse.TotalAmount,
		"trade_no":     data.Data.AlipayTradeAppPayResponse.TradeNo,
		"status":       "Immediate",
	}

	marshal, err := json.Marshal(RequestPayData)
	if err != nil {
		zap.L().Error("支付信息发送到后端序列化失败", zap.String("订单编号", data.Data.AlipayTradeAppPayResponse.OutTradeNo), zap.Error(err))
	}
	// 此处判定，若设备有有效期，则不再继续累加套餐
	data1, err1 := utils.HttpClient(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".UserServer")+"/device/SetMealSurplusDays?device_id="+data.Imei, "GET", "", c.GetString("X-Token"))
	if err1 != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}

	var ReturnDataMsg SetMealSurplusDaysStruct
	err = json.Unmarshal(data1, &ReturnDataMsg)
	if err != nil {
		zap.L().Error("返回值序列化失败！", zap.Error(err))
	}
	t1, err := time.Parse("2006-01-02 15:04:05", ReturnDataMsg.Data.EndDays)
	// time.Now().Before(t1) 当前时间小于t1, 代表设备有有效期，可以进行预存
	if ReturnDataMsg.Data.FlowTotal > 0 && time.Now().Before(t1) {
		// 根据package_id获取套餐详情，并将套餐信息缓存到数据库中
		data2, err2 := utils.HttpClient(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".DeviceServer")+"/devices/PackageInfo?device="+data.Imei+"&id="+data.PackageId, "GET", "", c.GetString("X-Token"))
		if err2 != nil {
			zap.L().Error("HttpClient请求发送失败", zap.Error(err))
		}

		var ReturnDat2 PackageInfoStruct
		err = json.Unmarshal(data2, &ReturnDat2)
		if err != nil {
			zap.L().Error("返回值序列化失败！", zap.Error(err))
		}
		// 将返回的数据写入数据库
		// 保存数据
		global.DB.Create(&users.DevicePreCharge{
			Model:        gorm.Model{},
			UserID:       c.GetInt("UserID"),
			DeviceID:     data.Imei,
			TotalFlow:    ReturnDat2.Data.TotalFlow,
			PackageID:    data.PackageId,
			PackageName:  ReturnDat2.Data.PackageName,
			PackageDays:  ReturnDat2.Data.PackageDays,
			PackagePrice: ReturnDat2.Data.PackagePrice,
			Operator:     c.GetString("UserID"),
			Status:       0,
		})

		PackagePreChargeData := map[string]interface{}{
			"UserID":        c.GetString("UserID"),
			"UserInfo":      data.Data.AlipayTradeAppPayResponse.AuthAppId,
			"DeviceID":      data.Imei,
			"TotalFlow":     ReturnDat2.Data.TotalFlow,
			"PackageID":     data.PackageId,
			"PackageName":   ReturnDat2.Data.PackageName,
			"PackageDays":   ReturnDat2.Data.PackageDays,
			"PackagePrice":  ReturnDat2.Data.PackagePrice,
			"Operator":      c.GetString("UserID"),
			"PackageStatus": 0,
			"OutTradeNo":    data.Data.AlipayTradeAppPayResponse.OutTradeNo,
		}

		PackagePreChargeJson, err := json.Marshal(PackagePreChargeData)
		if err != nil {
			zap.L().Error("返回值序列化失败！", zap.Error(err))
		}
		data3, err3 := utils.HttpClient(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".DeviceServer")+"/devices/PackagePreCharge", "POST", string(PackagePreChargeJson), c.GetString("X-Token"))

		if err3 != nil {
			zap.L().Error("预存信息发送到后端失败", zap.String("订单编号", data.Data.AlipayTradeAppPayResponse.OutTradeNo), zap.Error(err))
		} else {
			zap.L().Info("预存信息发送到后端成功！", zap.String("订单编号", data.Data.AlipayTradeAppPayResponse.OutTradeNo), zap.String("data", string(data3)))
		}

		// 数据缓存完后返回成功
		c.String(http.StatusOK, "SUCCESS")
	}
	zap.L().Info("发送数据到后端PayStatus", zap.String("data", string(marshal)))
	ReturnDataTemp, err := utils.HttpClient(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".DeviceServer")+"/devices/PayStatus", "POST", string(marshal), c.GetString("X-Token"))
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}

	var temp testStruct
	err = json.Unmarshal(ReturnDataTemp, &temp)
	if err != nil {
		zap.L().Error("返回值序列化失败！", zap.Error(err))
	}
	SaleStatus := DeviceSaleMode(data.Imei, data.PackageId, data.Data.AlipayTradeAppPayResponse.TotalAmount, data.Data.AlipayTradeAppPayResponse.TotalAmount, data.Data.AlipayTradeAppPayResponse.OutTradeNo, "AliPay", c)
	if SaleStatus == false {
		zap.L().Error("计算分润数据发送到后端失败，请核实！", zap.String("device_number", data.Imei), zap.String("money", data.Data.AlipayTradeAppPayResponse.TotalAmount))
	}
	if temp.Code == 20000 {
		ReturnData := map[string]interface{}{
			"code":       20000,
			"msg":        "支付成功！",
			"OutTradeNo": data.Data.AlipayTradeAppPayResponse.OutTradeNo,
		}
		c.JSON(http.StatusOK, ReturnData)
	} else {
		ReturnData := map[string]interface{}{
			"code":       20000,
			"msg":        "订单支付成功，发送到后端异常！",
			"OutTradeNo": data.Data.AlipayTradeAppPayResponse.OutTradeNo,
		}
		c.JSON(http.StatusOK, ReturnData)
	}
}

// WechatPayCallback 微信支付回调地址
func WechatPayCallback(c *gin.Context) {
	b, _ := c.GetRawData()

	zap.L().Info("支付回调参数验证！", zap.String("callback", string(b)))

	var data WechatPayCallbackStruct
	err := xmljson.Unmarshal(b, &data)
	if err != nil {
		zap.L().Error("xml反序列化失败", zap.Error(err))
	}

	// 返回数据给微信官方
	var ReturnData WechatPayCallbackReturnStruct
	ReturnData.ReturnCode = "SUCCESS"
	ReturnData.ReturnMsg = "OK"

	result := global.DB.Table("pay_data").Where("out_trade_no = ?", data.OutTradeNo).First(&PayData{}).Updates(PayData{
		Status:      2,
		OpenId:      data.Openid,
		TimeEnd:     data.TimeEnd,
		WeChatAppID: data.Appid,
		WeChatMchID: data.MchId,
	})

	if result.Error != nil {
		zap.L().Error("数据库更新失败失败！", zap.Error(result.Error))
	}
	temp1 := PayData{}

	global.DB.Table("pay_data").Where("out_trade_no = ?", data.OutTradeNo).First(&temp1)
	// 将支付信息发送到后端
	PayData1 := map[string]interface{}{
		"imei":       temp1.DeviceId,
		"package_id": temp1.PackageId,
		"status":     "Immediate",
	}

	// 此处判定，若设备有有效期，则不再继续累加套餐
	data1, err1 := utils.HttpClient(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".DeviceServer")+"/devices/SetMealSurplusDays?device_id="+temp1.DeviceId, "GET", "", c.GetString("X-Token"))
	if err1 != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}
	zap.L().Info("/device/SetMealSurplusDays接口返回信息：", zap.String("data", string(data1)))
	var ReturnDataMsg SetMealSurplusDaysStruct
	err = json.Unmarshal(data1, &ReturnDataMsg)
	if err != nil {
		zap.L().Error("返回值序列化失败！", zap.Error(err))
	}
	t1, err := time.Parse("2006-01-02 15:04:05", ReturnDataMsg.Data.EndDays)

	// 预声明变量，goto跳转的时候不允许声明变量
	var ReturnDataTemp, marshal []byte
	var temp testStruct

	// time.Now().Before(t1) 当前时间小于t1, 代表设备有有效期，可以进行预存
	if ReturnDataMsg.Data.FlowTotal > 0 && time.Now().Before(t1) {
		// 根据package_id获取套餐详情，并将套餐信息缓存到数据库中
		data2, err2 := utils.HttpClient(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".DeviceServer")+"/devices/PackageInfo?device="+temp1.DeviceId+"&id="+temp1.PackageId, "GET", "", c.GetString("X-Token"))
		zap.L().Info("/devices/PackageInfo接口返回信息：", zap.String("data", string(data2)))
		if err2 != nil {
			zap.L().Error("HttpClient请求发送失败", zap.Error(err))
		}

		var ReturnDat2 PackageInfoStruct
		err = json.Unmarshal(data2, &ReturnDat2)
		if err != nil {
			zap.L().Error("返回值序列化失败！", zap.Error(err))
		}
		// 将返回的数据写入数据库
		// 保存数据
		global.DB.Create(&users.DevicePreCharge{
			Model:        gorm.Model{},
			UserID:       c.GetInt("UserID"),
			DeviceID:     temp1.DeviceId,
			TotalFlow:    ReturnDat2.Data.TotalFlow,
			PackageID:    temp1.PackageId,
			PackageName:  ReturnDat2.Data.PackageName,
			PackageDays:  ReturnDat2.Data.PackageDays,
			PackagePrice: ReturnDat2.Data.PackagePrice,
			Operator:     c.GetString("UserID"),
			Status:       0,
		})
		// 将预存的数据保存到后端
		PackagePreChargeData := map[string]interface{}{
			"UserID":        c.GetString("UserID"),
			"UserInfo":      data.Openid,
			"DeviceID":      temp1.DeviceId,
			"TotalFlow":     ReturnDat2.Data.TotalFlow,
			"PackageID":     temp1.PackageId,
			"PackageName":   ReturnDat2.Data.PackageName,
			"PackageDays":   ReturnDat2.Data.PackageDays,
			"PackagePrice":  ReturnDat2.Data.PackagePrice,
			"Operator":      c.GetString("UserID"),
			"PackageStatus": 0,
			"OutTradeNo":    data.OutTradeNo,
		}

		PackagePreChargeJson, err := json.Marshal(PackagePreChargeData)
		if err != nil {
			zap.L().Error("返回值序列化失败！", zap.Error(err))
		}
		data3, err3 := utils.HttpClient(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".DeviceServer")+"/devices/PackagePreCharge", "POST", string(PackagePreChargeJson), c.GetString("X-Token"))
		zap.L().Info("/devices/PackagePreCharge接口返回信息：", zap.String("data", string(data3)))
		if err3 != nil {
			zap.L().Error("预存信息发送到后端失败", zap.String("订单编号", data.OutTradeNo), zap.Error(err))
		} else {
			zap.L().Info("预存信息发送到后端成功！", zap.String("订单编号", data.OutTradeNo), zap.String("data", string(data3)))
		}

		// 数据缓存完后返回成功
		goto LabelDeviceSaleModeFunc
	}

	marshal, err = json.Marshal(PayData1)
	if err != nil {
		zap.L().Error("支付信息发送到后端序列化失败", zap.String("订单编号", data.OutTradeNo), zap.Error(err))
	}

	ReturnDataTemp, err = utils.HttpClient(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".DeviceServer")+"/devices/PayStatus", "POST", string(marshal), c.GetString("X-Token"))
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}
	zap.L().Info("/devices/PayStatus接口调用信息", zap.String("data", string(marshal)))
	zap.L().Info("/devices/PayStatus接口返回信息", zap.String("data", string(ReturnDataTemp)))
	err = json.Unmarshal(ReturnDataTemp, &temp)
	if err != nil {
		zap.L().Error("返回值序列化失败！", zap.Error(err))
	}

	if temp.Code != 20000 {
		zap.L().Error("订单支付成功，发送到后端异常！", zap.String("OutTradeNo", data.OutTradeNo))
	}

	// 将设备号和充值的金额返回给后端，让后端计算分润
LabelDeviceSaleModeFunc:
	SaleStatus := DeviceSaleMode(temp1.DeviceId, temp1.PackageId, temp1.Money, temp1.PackageMoney, temp1.OutTradeNo, "WeChat", c)
	if SaleStatus == false {
		zap.L().Error("计算分润数据发送到后端失败，请核实！", zap.String("device_number", temp1.DeviceId), zap.String("money", temp1.Money))
	}

	c.XML(http.StatusOK, ReturnData)
}

func AliNotify(c *gin.Context) {
	// 初始化Alipay

	aliClient, err := utils.AliPayClientInitV2()
	if err != nil {
		zap.L().Error("AliPayClientInit初始化失败！", zap.Error(err))
	}

	ok, err := aliClient.VerifySign(c.Request.Form)
	if err != nil {
		zap.L().Error("异步通知验证签名发生错误！", zap.Error(err))
	}

	if ok == false {
		zap.L().Error("异步通知验证签名未通过！", zap.Error(err))
	}

	zap.L().Info("异步通知验证签名通过！")

	var outTradeNo = c.Request.Form.Get("out_trade_no")
	var p = alipay.TradeQuery{}
	p.OutTradeNo = outTradeNo
	rsp, err := aliClient.TradeQuery(p)
	if err != nil {
		fmt.Printf("异步通知验证订单 %s 信息发生错误: %s \n", outTradeNo, err.Error())
	}
	if rsp.IsSuccess() == false {
		fmt.Printf("异步通知验证订单 %s 信息发生错误: %s-%s \n", outTradeNo, rsp.Content.Msg, rsp.Content.SubMsg)
	}

	zap.L().Info("订单支付成功！", zap.String("outTradeNo", outTradeNo))
}

func quitUrl(c *gin.Context) {
	c.String(http.StatusOK, "支付中断")
}

// OpenTradeRefund 支付宝统一收单交易退款接口
func OpenTradeRefund(c *gin.Context) {
	b, _ := c.GetRawData() // 从c.Request.Body读取请求数据

	// 调用处传入退款商家订单号，退款金额
	var tempRefundMoneyStruct OpenTradeRefundMoneyStruct
	err := json.Unmarshal(b, &tempRefundMoneyStruct)
	if err != nil {
		zap.L().Error("json反序列化失败！", zap.Error(err))
	}

	// 初始化订单
	var p = alipay.TradeRefund{}
	// 商户订单号
	p.OutTradeNo = tempRefundMoneyStruct.OutTradeNo

	// 退款金额
	p.RefundAmount = tempRefundMoneyStruct.RefundMoney

	// 初始化Alipay
	aliClient, err := utils.AliPayClientInitV2()
	if err != nil {
		zap.L().Error("AliPayClientInit初始化失败！", zap.Error(err))
	}
	result, _ := aliClient.TradeRefund(p)
	if result.Content.Code == "10000" && result.Content.Msg == "Success" {
		marshal, err := json.Marshal(result)
		if err != nil {
			zap.L().Error("JSON序列化失败！", zap.Error(err))
		}
		zap.L().Info("退款成功", zap.String("ReturnData", string(marshal)))
		// 退款成功以后，更新该笔订单的状态，修改为已退款
		global.DB.Table("pay_data").Where("out_trade_no = ?", p.OutTradeNo).Update("status", 3)
	}
	c.JSON(http.StatusOK, result)
}

// OpenTradeFastPayRefundQuery 支付宝统一收单交易退款查询
func OpenTradeFastPayRefundQuery(c *gin.Context) {
	// 初始化订单

	var p = alipay.TradeFastPayRefundQuery{}
	// 商户订单号
	p.OutTradeNo = "3549699361002749952"
	p.OutRequestNo = "" // 必须 请求退款接口时，传入的退款请求号，如果在退款请求时未传入，则该值为创建交易时的外部交易号
	// 初始化Alipay
	aliClient, err := utils.AliPayClientInitV2()
	if err != nil {
		zap.L().Error("AliPayClientInit初始化失败！", zap.Error(err))
	}
	result, _ := aliClient.TradeFastPayRefundQuery(p)

	if result.Content.Code == "10000" && result.Content.Msg == "Success" {

	}
	c.JSON(http.StatusOK, result)
}

// PayHistory 查询设备充值历史
func PayHistory(c *gin.Context) {
	urlData := c.Request.URL.Query().Encode()
	data := strings.Split(urlData, "=")

	var foods []PayData
	result := global.DB.Table("pay_data").Where("device_id = ?", data[1]).Find(&foods)
	if result.Error != nil {
		zap.L().Error("数据库查询失败！", zap.Error(result.Error))
	}

	var ReturnDataTemp []map[string]string
	if len(foods) == 0 {
		ReturnData := map[string]interface{}{
			"code": 20000,
			"data": "",
		}
		c.JSON(http.StatusOK, ReturnData)
		return
	}
	var v PayData
	for _, v = range foods {

		temp2 := ""
		if v.Status == 1 {
			temp2 = "未支付"
		} else if v.Status == 2 {
			temp2 = "已支付"
		} else if v.Status == 3 {
			temp2 = "已退款"
		}
		if v.Status == 2 {
			temp1 := map[string]string{
				"imei":           v.DeviceId,
				"memo":           v.PackageName,
				"total_flow":     v.PackageTotalFlow,
				"out_trade_no":   v.OutTradeNo,
				"recharge_price": v.Money,
				"date_create":    v.CreatedAt.Format("2006-01-02 15:04:05"),
				"status":         temp2,
			}

			ReturnDataTemp = append(ReturnDataTemp, temp1)
		}

	}
	ReturnData := map[string]interface{}{
		"code": 20000,
		"data": ReturnDataTemp,
	}

	c.JSON(http.StatusOK, ReturnData)
}

// PayDataAll 后台页面展示支付信息使用
func PayDataAll(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("pageSize")
	keyword := c.Query("keyword")
	//AgentID := c.DefaultQuery("agent_id", "1")
	//pageTemp, _ := strconv.ParseInt(page, 10, 64)
	pageSizeTemp, _ := strconv.ParseInt(pageSize, 10, 64)
	pageSizeData, _ := strconv.Atoi(pageSize)

	var foods []PayData
	var count int64
	pageTemp, _ := strconv.Atoi(page)
	offset := (pageTemp - 1) * pageSizeData

	var tempTokenData TokenStruct
	err := json.Unmarshal([]byte(c.GetString("X-Token")), &tempTokenData)
	if err != nil {
		zap.L().Error("Json序列化失败", zap.Error(err))
	}

	AgentID := tempTokenData.UserAgentDev
	RequestData, err := utils.HttpClient(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".DeviceServer")+"/devices/get_agent_id?agent_id="+AgentID, "GET", "", c.GetString("X-Token"))
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}
	var tempData test2Struct
	err = json.Unmarshal(RequestData, &tempData)
	if err != nil {
		zap.L().Error("Json序列化失败", zap.Error(err))
	}

	var data *gorm.DB
	if tempTokenData.IsRoot == false {
		data = global.DB.Table("pay_data").Offset(offset).Limit(pageSizeData).Where("agent_id IN ?", tempData.Data).Order("created_at desc").Find(&foods).Offset(-1).Limit(-1).Count(&count)
	} else if keyword != "" {
		data = global.DB.Table("pay_data").Offset(offset).Limit(pageSizeData).Find(&foods, "agent_id IN ? AND mobile = ? OR out_trade_no = ? OR device_id = ?", tempData.Data, keyword, keyword, keyword).Offset(-1).Limit(-1).Count(&count)
	} else {
		data = global.DB.Table("pay_data").Offset(offset).Limit(pageSizeData).Order("created_at desc").Find(&foods).Offset(-1).Limit(-1).Count(&count)
	}

	pageNum := count / pageSizeTemp
	if count%pageSizeTemp != 0 {
		pageNum++
	}

	if data.Error != nil {
		zap.L().Error("数据库查询失败！", zap.Error(data.Error))
	}
	ReturnData := map[string]interface{}{
		"msg":    "OK",
		"code":   20000,
		"data":   foods,
		"pagNum": pageNum,
		"count":  count,
		"page":   pageTemp,
	}
	c.JSON(http.StatusOK, ReturnData)
}

// GerAllKeys 查询所有数据
func GerAllKeys(c *gin.Context) {

	page := c.Query("page")
	pageSize := c.Query("pageSize")
	keyword := c.Query("keyword")
	//AgentID := c.DefaultQuery("agent_id", "1")
	//pageTemp, _ := strconv.ParseInt(page, 10, 64)
	pageSizeTemp, _ := strconv.ParseInt(pageSize, 10, 64)
	pageSizeData, _ := strconv.Atoi(pageSize)

	var foods []PayData
	var count int64
	pageTemp, _ := strconv.Atoi(page)
	offset := (pageTemp - 1) * pageSizeData

	var tempTokenData TokenStruct
	err := json.Unmarshal([]byte(c.GetString("X-Token")), &tempTokenData)
	if err != nil {
		zap.L().Error("Json序列化失败", zap.Error(err))
	}

	var data *gorm.DB
	if tempTokenData.IsRoot == false {
		data = global.DB.Table("pay_config_data").Offset(offset).Limit(pageSizeData).Where("status = ?", 1).Order("created_at desc").Find(&foods).Offset(-1).Limit(-1).Count(&count)
	} else if keyword != "" {
		data = global.DB.Table("pay_config_data").Offset(offset).Limit(pageSizeData).Find(&foods, "status = ? AND app_id = ? OR mch_id = ? OR secret = ?", 1, keyword, keyword, keyword).Offset(-1).Limit(-1).Count(&count)
	} else {
		data = global.DB.Table("pay_config_data").Offset(offset).Limit(pageSizeData).Where("status = ? AND agent_id = ?", 1, tempTokenData.UserAgentDev).Order("created_at desc").Find(&foods).Offset(-1).Limit(-1).Count(&count)
	}

	pageNum := count / pageSizeTemp
	if count%pageSizeTemp != 0 {
		pageNum++
	}

	if data.Error != nil {
		zap.L().Error("数据库查询失败！", zap.Error(data.Error))
	}
	ReturnData := map[string]interface{}{
		"msg":    "OK",
		"code":   20000,
		"data":   foods,
		"pagNum": pageNum,
		"count":  count,
		"page":   pageTemp,
	}
	c.JSON(http.StatusOK, ReturnData)
}

// AddKeys 添加key
func AddKeys(c *gin.Context) {
	b, _ := c.GetRawData() // 从c.Request.Body读取请求数据
	var keyData KeysData
	err := json.Unmarshal(b, &keyData)
	if err != nil {
		zap.L().Error("JSON序列化失败！", zap.Error(err))
		ReturnData := map[string]interface{}{
			"msg":  "参数传入错误",
			"code": 20000,
		}
		c.JSON(http.StatusOK, ReturnData)
		return
	}
	// 校验前端传入的数据是否正确
	if utils.IsMobile(keyData.Mobile) == false {
		ReturnData := map[string]interface{}{
			"msg":  "手机号校验失败！",
			"code": 20000,
		}
		c.JSON(http.StatusOK, ReturnData)
		return
	} else if utils.IsEmail(keyData.Email) == false {
		ReturnData := map[string]interface{}{
			"msg":  "邮箱地址校验失败！",
			"code": 20000,
		}
		c.JSON(http.StatusOK, ReturnData)
		return
	}
	// 判断数据库中是否有记录
	var tempData PayConfigData
	global.DB.Where("app_id = ?", keyData.AppID).First(&tempData)
	if tempData.ID > 0 && tempData.AppID != "" {
		ReturnData := map[string]interface{}{
			"msg":  "数据已存在，不可重复插入！",
			"code": 20000,
		}
		c.JSON(http.StatusOK, ReturnData)
		return
	}

	// 保存数据
	global.DB.Create(&PayConfigData{
		Model:      gorm.Model{},
		Name:       keyData.Name,
		Mobile:     keyData.Mobile,
		Email:      keyData.Email,
		AppID:      keyData.AppID,
		PrivateKey: keyData.PrivateKey,
		PublicKey:  keyData.PublicKey,
		Status:     keyData.Status,
		IsDelete:   false,
	})

	ReturnData := map[string]interface{}{
		"msg":  "ok",
		"code": 20000,
	}

	c.JSON(http.StatusOK, ReturnData)
}

// EditKeys 修改key
func EditKeys(c *gin.Context) {
	b, _ := c.GetRawData() // 从c.Request.Body读取请求数据
	var keyData KeysData
	err := json.Unmarshal(b, &keyData)
	if err != nil {
		zap.L().Error("JSON序列化失败！", zap.Error(err))
		ReturnData := map[string]interface{}{
			"msg":  "参数传入错误",
			"code": 20000,
		}
		c.JSON(http.StatusOK, ReturnData)
		return
	}
	// 校验前端传入的数据是否正确
	if utils.IsMobile(keyData.Mobile) == false {
		ReturnData := map[string]interface{}{
			"msg":  "手机号校验失败！",
			"code": 20000,
		}
		c.JSON(http.StatusOK, ReturnData)
		return
	} else if utils.IsEmail(keyData.Email) == false {
		ReturnData := map[string]interface{}{
			"msg":  "邮箱地址校验失败！",
			"code": 20000,
		}
		c.JSON(http.StatusOK, ReturnData)
		return
	}
	// 修改数据
	var tempData PayConfigData
	global.DB.Where("name = ?", keyData.Name).First(&tempData)
	tempData.Name = keyData.Name
	tempData.Mobile = keyData.Mobile
	tempData.Email = keyData.Email
	tempData.AppID = keyData.AppID
	tempData.PrivateKey = keyData.PrivateKey
	tempData.PublicKey = keyData.PublicKey
	tempData.Status = keyData.Status
	tempData.AppID = keyData.AppID
	tempData.MchID = keyData.MchID
	tempData.PayNotify = keyData.PayNotify
	tempData.RefundNotify = keyData.RefundNotify
	tempData.Secret = keyData.Secret

	global.DB.Save(&tempData)

	ReturnData := map[string]interface{}{
		"msg":  "ok",
		"code": 20000,
	}

	c.JSON(http.StatusOK, ReturnData)

}

// DelKeys 删除key
func DelKeys(c *gin.Context) {
	// 逻辑删除
	b, _ := c.GetRawData() // 从c.Request.Body读取请求数据
	var keyData DelKeysStruct
	err := json.Unmarshal(b, &keyData)
	if err != nil {
		zap.L().Error("JSON序列化失败！", zap.Error(err))
		ReturnData := map[string]interface{}{
			"msg":  "参数传入错误",
			"code": 20000,
		}
		c.JSON(http.StatusOK, ReturnData)
		return
	}
	var temp PayConfigData
	result := global.DB.Model(&temp).Where("id = ?", keyData.ID).Update("status", false)

	if result.Error != nil {
		zap.L().Error("数据库查询失败！", zap.Error(result.Error), zap.String("len", strconv.FormatInt(result.RowsAffected, 10)))
		ReturnData := map[string]interface{}{
			"msg":  "数据库查询错误！",
			"code": 20000,
		}
		c.JSON(http.StatusOK, ReturnData)
		return
	} else {
		ReturnData := map[string]interface{}{
			"msg":  "OK",
			"code": 20000,
		}
		c.JSON(http.StatusOK, ReturnData)
	}
}
