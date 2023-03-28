package pays

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine) {
	OpenPay := e.Group("/pays")
	{
		//OpenPay.GET("V2/TradePagePay", WebTradeWapPayV3)
		//OpenPay.GET("V2/TradeWapPay", AppTradeWapPayV3)

		OpenPay.Any("/callback", AliCallback)                                // 支付宝支付回调地址
		OpenPay.Any("/WechatPayCallback", WechatPayCallback)                 // 微信支付回调地址
		OpenPay.GET("/TradeRefund", OpenTradeRefund)                         // 支付宝统一收单交易退款接口
		OpenPay.GET("/TradeFastPayRefundQuery", OpenTradeFastPayRefundQuery) // 支付宝统一收单交易退款查询
		OpenPay.GET("/payHistory", PayHistory)                               // 充值历史
		OpenPay.GET("/PayDataAll", PayDataAll)                               // 前端使用的带分页的充值历史，支持搜索

		// APP支付 V2版本
		OpenPay.POST("/TradeAppPay", MobileAppWapPay)     // 支付宝支付
		OpenPay.POST("/WechatPay", WechatPay)             // 微信支付
		OpenPay.POST("/WechatPayH5", WechatPayH5)         // 微信支付H5版
		OpenPay.POST("/WechatPayRefund", WechatPayRefund) // 微信退款

		// 代理商支付宝keys管理
		OpenPay.GET("/keys", GerAllKeys)
		OpenPay.POST("/addKeys", AddKeys)
		OpenPay.PUT("/editKeys", EditKeys)
		OpenPay.DELETE("/delKeys", DelKeys)

		OpenPay.Any("/quitUrl", quitUrl)
		OpenPay.Any("/notify", AliNotify)

	}
}
