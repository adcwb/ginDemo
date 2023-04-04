package wechat

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine) {
	chat := e.Group("/chat")
	{
		// GET请求
		chat.GET("/GetAccessToken", GetWeChatAccessToken)
		chat.Any("/Callback", CallbackWechat)

	}
}
