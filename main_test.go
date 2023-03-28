package main

import (
	"fmt"
	"ginDemo/utils"
	"testing"
)

func TestName(t *testing.T) {
	// 网络连接测试
	ok := utils.CheckNetWorkStatus()
	if ok {
		fmt.Println("网络测试通过！", ok)
	} else {
		fmt.Println("网络测试失败！", ok)
	}
}
