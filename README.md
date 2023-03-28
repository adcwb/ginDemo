# 项目名称

<!--一句话简短的描述项目-->



## 功能特性

<!--描述项目的核心功能-->



## 软件架构

<!--可选，简单描述一下项目的架构-->



## 快速开始

<!--如何快速的部署项目-->



### 依赖检查

<!--描述该项目的依赖，比如依赖的包，工具或者其他任何依赖项-->



### 运行项目

<!--描述如何运行该项目-->



## 使用指南

<!--描述如何使用该项目-->



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
