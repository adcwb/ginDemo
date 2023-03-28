package test

import (
	"fmt"
	"ginDemo/global"
	"github.com/gin-gonic/gin"
)

func test() {
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println(">>>>>>>>>>>>>>>>Hello Job !!>>>>>>>>>>>>>>>>")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
}
func JobTest(c *gin.Context) {
	global.JobS.Every(1).Seconds().Do(test)
}

func JobStop(c *gin.Context) {
	global.JobS.Stop()
}
