# 使用官方Go镜像作为构建环境
FROM golang:1.22 as builder

# 设置工作目录
WORKDIR /app

# 复制go mod和sum文件
COPY go.mod go.sum ./

# 下载所有依赖
RUN go mod download

# 复制源代码到容器中
COPY pkg .

# 构建应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# 使用scratch作为运行环境
FROM scratch

# 从builder镜像中复制构建的可执行文件
COPY --from=builder /app/main .

# 运行应用程序
ENTRYPOINT ["./main"]