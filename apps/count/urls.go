package count

import "github.com/gin-gonic/gin"

func Routers(e *gin.Engine) {
	count := e.Group("/count")
	{
		count.GET("/GetQueryDB", GetQueryDB)
		count.GET("/ExpressDelivery", ExpressDelivery)
		count.GET("/ExpressDeliveryMap", ExpressDeliveryMap)
		count.GET("/GetAutonumber", GetAutonumber)
		count.GET("/GetExpressDeliveryPoolMap", GetExpressDeliveryPoolMap)

		count.POST("/ExpressDeliveryPool", ExpressDeliveryPool)
		count.POST("/ExpressDeliveryPoolMap", ExpressDeliveryPoolMap)
		count.POST("/ExpressDeliveryCallback", ExpressDeliveryCallback)
	}

}
