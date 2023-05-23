package wechat

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine) {
	chat := e.Group("/chat")
	{
		// GET请求
		chat.GET("/GetAccessToken", GetWeChatAccessToken)
		//chat.GET("/GetJsAPITicket", GetWorkJsAPITicketToken)
		//chat.GET("/GetAgentTicket", GetAgentTicketToken)
		chat.GET("/GetWorkConfig", GetWorkConfig)
		chat.GET("/GetWorkAgentConfig", GetWorkAgentConfig)
		chat.GET("/GetWorkUserData", GetWorkUserData)
		chat.POST("/SaveWorkUserData", SaveWorkUserData)
		chat.GET("/GetKfList", GetKfList)
		chat.POST("/SendMessageQueue", SendMessageQueue)
		chat.GET("/GetQueueMessage", GetQueueMessage)

		chat.GET("/Callback", GetCallbackWechat)
		chat.POST("/Callback", PostCallbackWechat)
	}
}
