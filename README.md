# 项目名称

<!--一句话简短的描述项目-->
仿Django框架Demo


## 功能特性

<!--描述项目的核心功能-->

## 项目目录
```go
├─apps  // 项目代码
│  ├─count
│  ├─dingtalk
│  ├─pays
│  │  └─keys
│  ├─test
│  └─users
├─assets            // 静态资源
│  └─video
├─CA                // 证书
│  └─TLS
├─Docs              // 文档
├─downloadFile      // 下载文件
├─global            // 全局变量定义  
├─initialization    // 初始化组件
├─logs              // 日志
├─middleware        // 中间件
├─templates         // 模板文件
│  └─views          
├─uploadFile        // 上传文件存放目录
└─utils             // 第三方工具

```

## 软件架构

<!--可选，简单描述一下项目的架构-->



## 快速开始

<!--如何快速的部署项目-->
```bash
go build main.go
./main
```


### 依赖检查

<!--描述该项目的依赖，比如依赖的包，工具或者其他任何依赖项-->



### 运行项目

<!--描述如何运行该项目-->
```dockerfile
FROM golang:alpine

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文件app
RUN export GIN_MODE=release && go env -w GOPROXY=https://goproxy.cn,direct && go build -o app .
# RUN go build -o app .

# 移动到用于存放生成的二进制文件的 /dist 目录
# WORKDIR /dist

# 将二进制文件从 /build 目录复制到这里
# RUN cp /build/app .

# 声明服务端口
EXPOSE 20000

# 启动容器时运行的命令
CMD ["/build/app"]
```


## 使用指南

<!--描述如何使用该项目-->
#### 接口文档编写
- 官方文档 https://github.com/swaggo/gin-swagger
```GO
// @title 这里写标题
// @version 1.0
// @description 这里写描述信息
// @termsOfService http://swagger.io/terms/

// @contact.name 这里写联系人信息
// @contact.email 这里写联系人邮箱
// @contact.url http://www.swagger.io/support
// 开源协议
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 这里写接口服务的host
// @BasePath 这里写base path
func main() {
	r := gin.New()

	// liwenzhou.com ...

	r.Run()
}

```

#### 生成接口文档
```bash
# 编写完注释后，使用以下命令安装swag工具：
go get -u github.com/swaggo/swag/cmd/swag

# 在项目根目录执行以下命令，使用swag工具生成接口文档数据。
swag fmt && swag init
```

#### 引入gin-swagger渲染文档数据
```go
import (
  _ "ginDemo/Docs"
  "github.com/gin-gonic/gin"
  swaggerFiles "github.com/swaggo/files"
  gs "github.com/swaggo/gin-swagger"
)
// 路由注册
r.GET("/swagger/*any", gs.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE"))
// 打开浏览器访问http://localhost:8000/swagger/index.html就能看到Swagger 2.0 Api文档了。
```

## 如何贡献

<!--告诉其他开发者如何给该项目贡献源码-->



## 社区论坛

<!--如果项目有社区或者论坛反馈，可以在此处加上-->



## 关于作者

<!--此处可以写上项目作者的联系方式-->



## 客户案例

<!--描述一下谁在使用该项目，展示一下项目的影响力，一般没啥用-->



## 许可证

<!--这里链接上该项目的开源许可证-->



### Golang 在 三大平台下如何交叉编译

- GOOS：目标平台的操作系统（部分不常见的没有列出）
    - darwin
    - freebsd
    - openbsd
    - linux
    - windows
    - ios
    - android
    - aix
- GOARCH：目标平台的体系架构（）
    - 386
    - amd64
    - arm
- 交叉编译不支持 CGO 所以要禁用它



```bash
# 查看支持的平台

go tool dist list

# 终极打包工具， 经测试目前无法使用
https://github.com/upx/upx
```


#### Windows平台编译设置
```bash
# 设置为Mac平台
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build main.go

# 设置为Linux系统
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build main.go

# 编译为Mips
SET GOOS=linux
SET GOARCH=mipsle
SET GOMIPS=softfloat
SET CGO_ENABLED=0
go build -trimpath -ldflags="-s -w"  main.go
upx -9 main # 未测试通过


# 调整至windows
SET GOOS=windows
SET GOARCH=amd64

```



#### Linxu平台编译设置
```bash
# 编译为Mac平台
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go

# 编译为windows系统
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go

#编译为Mips系统
GOOS=linux GOARCH=mipsle GOMIPS=softfloat CGO_ENABLED=0 go build -trimpath -ldflags="-s -w"  main.go

```


#### Mac平台编译设置
```bash
# 设置为Linux平台
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

# 设置为windows系统
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
```