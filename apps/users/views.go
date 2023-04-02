package users

import (
	"context"
	"encoding/json"
	"fmt"
	"ginDemo/global"
	"ginDemo/utils"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func Test(c *gin.Context) {
	method := c.Request.Method
	if method == "GET" {
		//
		data, err := utils.HttpClient(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".UserServer")+"/users/", "GET", "", c.GetString("X-Token"))
		if err != nil {
			zap.L().Error("HttpClient请求发送失败", zap.Error(err))
		}
		// 将返回的数据进行反序列化
		var m testStruct
		err = json.Unmarshal(data, &m)
		if err != nil {
			zap.L().Error("json序列化失败", zap.Error(err))
		}

		c.JSON(http.StatusOK, m)

	} else if method == "POST" {
		c.String(http.StatusNotFound, "404 page not found!")
	} else {
		c.String(http.StatusNotFound, "404 page not found!")
	}
}

// WeChatOpenid 根据code获取openid
func WeChatOpenid(c *gin.Context) {

	code := c.Query("code")   // 获取code
	brand := c.Query("brand") // 获取平台，
	fmt.Println(brand)
	// TODO 是否要记录OpenID和设备号，平台的对应关系
	if code == "" {
		ReturnData := map[string]interface{}{
			"code": 20001,
			"msg":  "code为空",
		}

		c.JSON(http.StatusOK, ReturnData)
		return
	}
	appid := "wxb4276ff58b1587a6"
	secret := "7d3d63bc6026d61e4a6ad5b4c0fc670a"
	/* 测试使用，不可调用支付
	PLATFORM = {
	    "CH": {
	        "app_id": "wxea09d930a8331bc5",
	        "secret": "b99864da5cfaa31f1f89892bf348e7be",
	        "pay_app_id": "81373949190384927734",
	        "pay_key": "3lzfexKlEyNFv90i0jZY"
	    },
	    "IWANT": {
	        "app_id": "wxea09d930a8331bc5",
	        "secret": "b99864da5cfaa31f1f89892bf348e7be",
	        "pay_app_id": "81373949190384927734",
	        "pay_key": "3lzfexKlEyNFv90i0jZY"
	    }
	}

	*/

	data, err := utils.HttpClient(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".WeChatOpenIdServer")+"/sns/oauth2/access_token?appid="+appid+"&secret="+secret+"&code="+code+"&grant_type=authorization_code", "GET", "", c.GetString(""))
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
	}
	// 将返回的数据进行反序列化
	var m WeChatOpenIDRequestStruct
	err = json.Unmarshal(data, &m)
	if err != nil {
		zap.L().Error("获取用户OpenID失败，序列化失败", zap.Error(err))
	}
	returnData := map[string]interface{}{
		"code": 20000,
		"data": m,
	}
	c.JSON(http.StatusOK, returnData)
}

// authHandler JWT测试函数
func authHandler(c *gin.Context) {
	// 用户发送用户名和密码过来
	data, _ := utils.GenToken("admin")
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "success",
		"data": gin.H{"token": data},
	})
	return
}

// Version 返回当前软件版本号
func Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "success",
		"version": "0.0.1",
	})
}

// RedirectObj 重定向测试函数，重定向所有的请求到baidu
func RedirectObj(c *gin.Context) {
	// HTTP重定向
	c.Redirect(http.StatusMovedPermanently, "https://www.baidu.com/")

	// c.Request.URL.Path = "/test2"
	// r.HandleContext(c) 路由重定向
}

// SendSmsCode 发送短信接口
//
//	@BasePath		/users/SendSmsCode
//	@Summary		发送短信接口
//	@Description	接收前端传递过来的手机号，并生成随机验证码发送给客户
//	@Tags			Users
//	@Accept			application/json
//	@Produce		application/json
//	@Param			Authorization	header	string			false	"Bearer 用户令牌"
//	@Param			object			body	SendDataStruct	true	"查询参数"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	testStruct "{"code":200,"data":"ok","msg":"ok"}"
//	@Router			/SendSmsCode [post]
func SendSmsCode(c *gin.Context) {
	b, _ := c.GetRawData()
	var SendData SendDataStruct
	err := json.Unmarshal(b, &SendData)
	if err != nil {
		zap.L().Error("Json序列化失败！！", zap.Error(err))
	}

	// 判断Redis数据库中是否有记录
	ctx := context.Background()

	result, err := global.REDIS.Get(ctx, SendData.PhoneNumbers).Result()

	if err == redis.Nil {
		zap.L().Error("RedisKey PhoneNumbers does not exist", zap.Error(err))

	} else if err != nil {
		zap.L().Error("RedisKey PhoneNumbers does not exist", zap.Error(err))
	}

	if len(result) > 0 {
		zap.L().Info("短信验证码安全机制触发！")
		ReturnData := map[string]interface{}{
			"code": 20010,
			"data": "短信发送失败，三分钟内只可以获取一次！",
		}
		c.JSON(http.StatusOK, ReturnData)
		return
	}

	code := utils.NumberCode()
	zap.L().Info("短信验证码", zap.String("mobile", SendData.PhoneNumbers), zap.String("value", code), zap.Error(err))

	if utils.IsMobile(SendData.PhoneNumbers) {
		sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
			PhoneNumbers:  tea.String(SendData.PhoneNumbers),
			SignName:      tea.String(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".aliYunSignName")),
			TemplateCode:  tea.String(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig") + ".aliYunTemplateCode")),
			TemplateParam: tea.String("{\"code\":" + "\"" + code + "\"" + "}"),
		}
		// 判断OpenID是否存在
		if SendData.OpenID != "" {
			var UserData User
			err := global.DB.Where("mobile = ? AND open_id = ?", SendData.PhoneNumbers, SendData.OpenID).First(&UserData).Error
			if err == nil {
				if UserData.ID > 0 {
					temp := map[string]string{
						"username":    SendData.PhoneNumbers,
						"device_type": "app",
					}
					marshal, err := json.Marshal(temp)
					if err != nil {
						zap.L().Error("序列化失败", zap.Error(err))
					}
					data, err := utils.HttpClient(global.CONFIG.GetString(global.CONFIG.GetString("RunConfig")+".UserServer")+"/users/getToken", "POST", string(marshal), c.GetString("X-Token"))
					if err != nil {
						zap.L().Error("HttpClient请求发送失败", zap.Error(err))
					}
					// 将返回的数据进行反序列化
					var m getTokenStruct
					err = json.Unmarshal(data, &m)
					if err != nil {
						zap.L().Error("用户获取Token失败，序列化失败", zap.Error(err))
					} else {
						zap.L().Info("用户" + SendData.PhoneNumbers + "自动登录成功！")
					}
					c.JSON(http.StatusOK, m)
					return
				}
			}
		}
		err := utils.SendSmSCode(sendSmsRequest)
		if err != nil {
			ReturnData := map[string]interface{}{
				"code": 20010,
				"data": "短信发送失败！",
			}
			c.JSON(http.StatusOK, ReturnData)
			return
		}

		ReturnData := map[string]interface{}{
			"code": 20000,
			"data": "短信发送成功！",
		}
		// 将数据记录到Redis数据库一份
		ctx := context.Background()
		err = global.REDIS.Set(ctx, SendData.PhoneNumbers, code, 180*time.Second).Err()

		if err != nil {
			zap.L().Error("Redis Set Key Error", zap.String("keys", SendData.PhoneNumbers), zap.String("value", code), zap.Error(err))
		}

		// 将数据记录到MySQL一份
		// 保存数据
		global.DB.Create(&SendSmsCodeData{
			Model:  gorm.Model{},
			Mobile: SendData.PhoneNumbers,
			Code:   code,
		})

		c.JSON(http.StatusOK, ReturnData)
	} else {
		ReturnData := map[string]interface{}{
			"code": 20012,
			"data": "手机号格式校验失败！",
		}
		c.JSON(http.StatusOK, ReturnData)
	}
}
