package utils

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"ginDemo/global"
	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Exist 判断所给路径文件/文件夹是否存在
func Exist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// RmFiles 删除文件
func RmFiles(path string) bool {
	err := os.Remove(path)
	if err != nil {
		return false
	}
	return true
}

// FileWriteString 将字符串写入指定文件中
func FileWriteString(filename, content string) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("open file failed, err:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			zap.L().Error("文件关闭失败！", zap.Error(err))
		}
	}(file)

	writeString, err := file.WriteString(content)
	if err != nil {
		zap.L().Error("文件写入失败！", zap.Error(err))
	} //直接写入字符串数据

	zap.L().Info("文件写入成功！", zap.String("文件名称", filename), zap.Int("写入数据量", writeString))
}

// UUID5 生成UUID5
func UUID5(temp string) string {
	namespace := [16]byte{128}
	data := uuid.NewV5(namespace, temp).String()
	return data
}

// IsMobile 校验手机号格式是否正确
func IsMobile(mobile string) (status bool) {
	result, _ := regexp.MatchString(`^(1[3|4|5|6|7|8|9][0-9]\d{4,8})$`, mobile)
	if result {
		return true
	} else {
		return false
	}
}

// IsEmail 校验邮箱格式是否正确
func IsEmail(email string) (status bool) {
	result, _ := regexp.MatchString("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", email)
	if result {
		return true
	} else {
		return false
	}
}

// NumberCode 生成数字验证码
func NumberCode() (data string) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	data = fmt.Sprintf("%06v", rnd.Int31n(1000000))
	return data
}

// MD5 MD5加密
func MD5(str string) string {
	data := []byte(str) //切片
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

func ToHexStr(str string) string {
	var newStr string
	for _, r := range str {
		hexStr := strconv.FormatInt(int64(r&0xFF), 16)
		if len(hexStr) == 1 {
			hexStr = "0" + hexStr
		}
		newStr += hexStr
	}
	return newStr
}

// GetRandomString 获取一个随机字符串
func GetRandomString(n int) string {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	bytes := []byte(str)
	var result []byte
	for i := 0; i < n; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

// GetLocalIP 获取本地IP
func GetLocalIP() []string {
	var ipStr []string
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces error:", err.Error())
		return ipStr
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					//获取IPv6
					/*if ipnet.IP.To16() != nil {
					    fmt.Println(ipnet.IP.String())
					    ipStr = append(ipStr, ipnet.IP.String())

					}*/
					//获取IPv4
					if ipnet.IP.To4() != nil {
						ipStr = append(ipStr, ipnet.IP.String())
					}
				}
			}
		}
	}
	return ipStr
}

// IsNumber 判断给定字符串是否是纯数字组成
func IsNumber(data string) bool {
	matchString, err := regexp.MatchString("\\d+", data)
	if err != nil {
		return false
	} else {
		return matchString
	}
}

// Paginate 分页器
func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(q.Get("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// FindList 判断元素是否在切片中
func FindList(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func ReturnEveryDays(startData, endData string, status bool) []string {
	layout := ""
	if status {
		layout = "2006-01-02 15:04:05"
	} else {
		layout = "2006-01-02"
	}

	// 加载时区
	locTime, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		zap.L().Error("加载时区失败！", zap.Error(err))
	}

	// 按照指定时区和指定格式解析字符串时间
	timeObj1, err := time.ParseInLocation(layout, startData, locTime)
	if err != nil {
		zap.L().Error("解析时间字符串失败！", zap.Error(err))
	}

	timeObj2, err2 := time.ParseInLocation(layout, endData, locTime)
	if err2 != nil {
		zap.L().Error("解析时间字符串失败！", zap.Error(err))
	}
	start := timeObj1.Unix()
	end := timeObj2.Unix()

	allDateArray := make([]string, 0)
	startTime := time.Unix(start, 0)
	endTime := time.Unix(end, 0)
	//After方法 a.After(b) a,b Time类型 如果a时间在b时间之后，则返回true
	for endTime.After(startTime) {
		allDateArray = append(allDateArray, startTime.Format(layout))
		startTime = startTime.AddDate(0, 0, 1)
	}
	allDateArray = append(allDateArray, endTime.Format(layout))
	return allDateArray
}

// CheckNetWorkStatus 网络检测函数
func CheckNetWorkStatus() bool {
	// 判定网络连接的变量，当以后有更好的判定方法时，可以在此函数继续拓展
	var temp1, temp2, temp3 bool

	// 检测设备到百度之前的网络连接
	var cmd1, cmd2 *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd1 = exec.Command("ping", "www.baidu.com")
	} else {
		cmd1 = exec.Command("ping", "www.baidu.com", "-c", "4", "-W", "5")
	}

	err1 := cmd1.Run()
	if err1 != nil {
		temp1 = false
	} else {
		temp1 = true
	}

	// 检测设备到服务端之间的连接
	if runtime.GOOS == "windows" {
		cmd2 = exec.Command("ping", "pays.tianchao.pro")
	} else {
		cmd2 = exec.Command("ping", "pays.tianchao.pro", "-c", "4", "-W", "5")
	}
	err2 := cmd2.Run()

	// 检测服务端是否在线
	urlTempData := "https://pays.tianchao.pro/api/test/ping"
	method := "GET"

	client := &http.Client{}
	req, _ := http.NewRequest(method, urlTempData, nil)
	res, _ := client.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	body, _ := io.ReadAll(res.Body)

	if strings.ToLower(string(body)) == "Pang" {
		temp3 = true
	} else {
		temp3 = false
	}

	if err2 != nil {
		temp2 = false
	} else {
		temp2 = true
	}
	fmt.Println(temp1, temp2, temp3)
	if temp1 || temp2 || temp3 {
		zap.L().Debug("网络检测：CheckNetWorkStatus函数运行，网络检测通过......")
		return true
	} else if temp1 == false && temp2 == false {
		zap.L().Debug("网络检测：CheckNetWorkStatus函数运行，请检查服务器，网络检测不通过.....")
		return false
	} else {
		zap.L().Debug("网络检测：CheckNetWorkStatus函数运行，网络检测通过.....")
		return true
	}
}

// SortMapToURLParams 将接收到的字典转化为URL查询字符串
func SortMapToURLParams(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var params []string
	for _, k := range keys {
		v := url.QueryEscape(m[k])
		params = append(params, k+"="+v)
	}
	return strings.Join(params, "&")
}

// Contains 判断A字符串是否包含B字符串
func Contains(a, b string) bool {
	return strings.Contains(a, b)
}

type Result struct {
	ActiveTime         string `xml:"activeTime"`
	ProdStatusName     string `xml:"prodStatusName"`
	ProdMainStatusName string `xml:"prodMainStatusName"`
	CertNumber         string `xml:"certNumber"`
	Number             string `xml:"number"`
}

type SvcCont struct {
	Result        Result `xml:"RESULT"`
	ResultCode    int    `xml:"resultCode"`
	ResultMsg     string `xml:"resultMsg"`
	TransactionID string `xml:"GROUP_TRANSACTIONID"`
}

// XmlToJson xml数据转化为json
func XmlToJson(xmlStr string) (string, error) {
	var svcCont SvcCont
	err := xml.Unmarshal([]byte(xmlStr), &svcCont)
	if err != nil {
		return "", err
	}
	jsonBytes, err := json.Marshal(svcCont)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// GetRandomElement 返回数组的随机一个元素
func GetRandomElement(arr interface{}) interface{} {
	rand.Seed(time.Now().UnixNano())

	switch a := arr.(type) {
	case []int:
		return a[rand.Intn(len(a))]
	case []string:
		return a[rand.Intn(len(a))]
	default:
		return nil
	}
}

func GetRedisKey(ctx context.Context, key string) (result string) {
	result, err := global.REDIS.Get(ctx, key).Result()

	if err == redis.Nil {
		zap.L().Error("RedisKey "+key+"does not exist", zap.Error(err))

	} else if err != nil {
		zap.L().Error("RedisKey "+key+"does not exist", zap.Error(err))
	}
	// 若有缓存直接返回
	if len(result) > 5 && result != "" {
		return result
	} else {
		return ""
	}
}

// GenerateRandomNumber 返回随机数字
func GenerateRandomNumber(digits int) string {
	rand.Seed(time.Now().UnixNano())

	// 生成指定位数的随机数字
	randNum := rand.Intn(int(math.Pow10(digits)))

	// 将结果转化为字符串并返回
	return strconv.Itoa(randNum)
}

func IsFutureTime(timeStr string) bool {
	inputTime, err := time.Parse("2006-01-02", timeStr)
	if err != nil {
		fmt.Println("Invalid time format")
		return false
	}

	currentTime := time.Now()

	return inputTime.After(currentTime)
}
