package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"ginDemo/global"
	"ginDemo/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type TokenData struct {
	Header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	} `json:"Header"`
	Payload struct {
		UserAgentCard string `json:"user_agent_card"`
		UserAgentDev  string `json:"user_agent_dev"`
		LoginTime     string `json:"login_time"`
		UserId        int    `json:"user_id"`
		IsRoot        bool   `json:"is_root"`
	} `json:"Payload"`
}

// JwtCheck 校验JwtCheck
func JwtCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		Token := c.Request.Header.Get("X-Token")
		zap.L().Debug("当前登录用户的Token为：", zap.String("Token: ", Token))
		uriTemp := make([]string, 0, 8)
		uriTemp = append(uriTemp,
			"/devices/DownloadTemplateFile",
			"/devices/DownloadDevicesTemplateFile",
			"/devices/DownloadDevicesGroupTemplateFile",
			"/users/getToken", "/app/login", "/users/SendSms",
			"/logs/gin.log", "/users/openid", "/count/ExpressDeliveryCallback",
		)
		if utils.FindList(uriTemp, c.Request.URL.Path) {
			c.Set("X-Token", "")
		} else if strings.Contains(c.Request.URL.Path, "assets") || strings.Contains(c.Request.URL.Path, "templates") || strings.Contains(c.Request.URL.Path, "views") {
			c.Set("X-Token", "")
		} else {
			if Token != "" {
				decoded, err := base64.StdEncoding.DecodeString(Token)
				if err != nil {
					zap.L().Error("解析Token失败", zap.Error(err))
				}

				decodeStrToken := strings.Split(string(decoded), ",")
				key := decodeStrToken[0]
				ctx := context.Background()

				result, err := global.REDIS.Get(ctx, "USER_TOKEN_KEY").Result()

				if err == redis.Nil {
					zap.L().Error("RedisKey USER_TOKEN_KEY does not exist", zap.Error(err))
				} else if err != nil {
					zap.L().Error("RedisKey USER_TOKEN_KEY does not exist", zap.Error(err))
				}

				if key != result {
					zap.L().Error("RedisKey USER_TOKEN_KEY 校验失败")

					returnData := map[string]interface{}{
						"data": "TOKEN_KEY 校验失败",
						"code": 9999,
					}

					c.JSON(http.StatusOK, returnData)
					c.Abort()
					return
				}

				decoded1, err := base64.StdEncoding.DecodeString(decodeStrToken[1])
				if err != nil {
					zap.L().Error("解析Token失败", zap.Error(err))
				}

				var temp1 TokenData
				err = json.Unmarshal(decoded1, &temp1)
				if err != nil {
					zap.L().Error("序列化Token失败", zap.Error(err))
				}

				tempData := temp1.Payload.LoginTime                    // JWT时间
				timeFormat := time.Now().Format("2006-01-02 15:04:05") // 当前时间

				t1, err := time.Parse("2006-01-02 15:04:05", tempData)
				t2, err := time.Parse("2006-01-02 15:04:05", timeFormat)

				if err == nil && t1.Before(t2) {
					zap.L().Error("当前Token已过期", zap.Error(err))
					returnData := map[string]interface{}{
						"data": "当前Token已过期",
						"code": 9999,
					}
					c.JSON(http.StatusOK, returnData)
					c.Abort()
					return
				}

				// 将数据json化，装入请求头中
				c.Set("UserID", temp1.Payload.UserId)
				token, err := json.Marshal(temp1.Payload)

				if err != nil {
					zap.L().Error("序列化Token失败", zap.Error(err))
				}
				c.Set("X-Token", string(token))

			} else {
				//returnData := map[string]interface{}{
				//	"data": "未检测到X-Token，验证失败！",
				//	"code": 9999,
				//}
				//c.JSON(http.StatusOK, returnData)
				//c.Abort()
				//return
			}
		}
	}
}
