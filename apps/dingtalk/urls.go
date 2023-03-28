package dingtalk

import "github.com/gin-gonic/gin"

func Routers(e *gin.Engine) {
	dingTalk := e.Group("/dingTalk")
	{
		// GET请求
		dingTalk.GET("/", Test)
		dingTalk.GET("/getToken", GetToken)
		dingTalk.GET("/GetAllForms", GetAllForms)
		dingTalk.GET("/GetFromProcessCode", GetFromProcessCode)
		dingTalk.GET("/test", test)

		dingTalk.GET("/GetUUID5", GetUUID5)

	}
}
