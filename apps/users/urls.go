package users

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine) {
	users := e.Group("/users")
	{
		// GET请求
		users.GET("/openid", WeChatOpenid)

		// POST请求
		users.POST("/", MiddleWareGetPost(), Test)

		// 发送短信
		users.POST("/SendSms", SendSmsCode)
	}

	e.GET("/version", MiddleWareGetPost(), Version)
	e.Any("/redirect", RedirectObj)

}
