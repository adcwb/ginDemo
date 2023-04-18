package test

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine) {
	test := e.Group("/test")
	{
		test.GET("/ping", Pang)
		test.GET("/test", RabbitTest)
		test.GET("/gocron", JobTest)
		test.GET("/stop", JobStop)
		test.GET("/operator", Operator)
		test.GET("/ThreeCodeMutualCheck", ThreeCodeMutualCheck)

	}
}
