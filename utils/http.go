package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"sync"
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

func sendHTTPRequest(method string, url string, body []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func ConcurrentHttpClient(urls, method, data string, Token string, concurrency int) [][]byte {
	// 创建通道和等待组
	resChan := make(chan []byte, concurrency)
	wg := sync.WaitGroup{}

	// 分割请求url
	urlArr := strings.Split(urls, ",")

	// 构造请求数据
	payload := strings.NewReader(data)

	// 构造请求客户端
	client := http.Client{}

	// 循环发起请求
	for _, url := range urlArr {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			req, err := http.NewRequest(method, url, payload)
			if err != nil {
				zap.L().Error("构造请求失败", zap.Error(err))
				return
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
				return
			}

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				zap.L().Error("读取body数据错误", zap.Error(err))
				return
			}

			// 将获取到的body数据放入通道
			resChan <- body
		}(url)
	}

	// 等待所有请求处理完成并关闭通道
	go func() {
		wg.Wait()
		close(resChan)
	}()

	// 从通道中读取所有响应数据并返回
	var result [][]byte
	for res := range resChan {
		result = append(result, res)
	}
	return result
}
