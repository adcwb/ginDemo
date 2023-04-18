package utils

import (
	"encoding/base64"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"time"
)

type TokenStruct struct {
	Sign      string `json:"Sign"`
	Timestamp string `json:"Timestamp"`
	AppKey    string `json:"AppKey"`
}

// HttpClient http请求客户端封装
func HttpClient(urls, method, data string, Token string) (body []byte, err error) {
	zap.L().Debug("interface start", zap.String("data", time.Now().Format("2006-01-02 15:04:05")))
	payload := strings.NewReader(data)
	client := http.Client{}

	req, err := http.NewRequest(method, urls, payload)
	if err != nil {
		return body, err
	}

	if Contains(Token, "AppKey") && Contains(Token, "Sign") {
		var tempData TokenStruct
		err = json.Unmarshal([]byte(Token), &tempData)
		if err != nil {
			zap.L().Error("Json序列化失败", zap.Error(err))
		}
		req.Header.Add("Sign", tempData.Sign)
		req.Header.Add("Timestamp", tempData.Timestamp)
		req.Header.Add("AppKey", tempData.AppKey)

	} else if len(Token) > 1 && Token != "xml" {
		StrBytes := []byte(Token)
		TokenEncoded := base64.StdEncoding.EncodeToString(StrBytes)
		req.Header.Add("X-Token", TokenEncoded)
	} else if Token == "xml" {
		req.Header.Add("Content-Type", "application/xml")
	} else if Token == "form" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Add("Content-Type", "application/json")
	}

	// 发起请求
	res, err := client.Do(req)
	if err != nil {
		zap.L().Error("HttpClient请求发送失败", zap.Error(err))
		return body, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			zap.L().Error("Body Close Error", zap.Error(err))
		}
	}(res.Body)

	body, err = io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("读取body数据错误", zap.Error(err))
		return body, err
	}

	//// 将返回的数据进行反序列化
	//var m map[string]interface{}
	//err = json.Unmarshal([]byte(body), &m)
	//if err != nil {
	//	fmt.Println("Unmarshal failed, ", err)
	//	return body, err
	//}

	// 将获取到的token返回出去
	zap.L().Debug("interface return", zap.String("data: ", time.Now().Format("2006-01-02 15:04:05")))
	return body, err

}
