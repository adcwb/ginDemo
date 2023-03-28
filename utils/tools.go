package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
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
	url := "https://pays.tianchao.pro/api/test/ping"
	method := "GET"

	client := &http.Client{}
	req, _ := http.NewRequest(method, url, nil)
	res, _ := client.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	body, _ := ioutil.ReadAll(res.Body)

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
