package pays

import (
	"github.com/gin-gonic/gin"
)

// DeviceSaleMode xxx
func DeviceSaleMode(DeviceNumber, PackageID, Money, PackageMoney, OutTradeNo, PayChannel string, c *gin.Context) (Status bool) {
	return true
}
