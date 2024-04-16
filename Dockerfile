FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

# 移动到工作目录：/build
WORKDIR /build

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文件app
#RUN export GIN_MODE=release && go env -w GOPROXY=https://goproxy.cn,direct && go build -ldflags "-s -w" -o app .
RUN go mod tidy && go build -ldflags "-s -w"  -toolexec="/build/bin/skywalking-go-agent-0.4.0-linux-amd64" -a  -o app .

FROM busybox:latest

COPY --from=builder /build /build

# 声明服务端口
EXPOSE 8000

WORKDIR /build
# 启动容器时运行的命令
CMD ["/build/app"]