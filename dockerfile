# 使用 Golang 官方镜像作为基础镜像
FROM golang

# 设置代理
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 将 go.mod 和 go.sum 文件复制到容器内
COPY go.mod go.sum ./

# 下载 Go 依赖
RUN go mod download

# 将整个项目复制到容器中
COPY . .

# 编译 Go 程序
RUN go build -o main .

# 启动程序
CMD ["/app/main"]
