# 使用官方的 Go 镜像作为基础镜像
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制项目文件
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o im-server ./cmd/main.go

# 使用精简的 Alpine 镜像作为运行时镜像
FROM alpine:latest

# 安装必要的工具
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制可执行文件
COPY --from=builder /app/im-server .
COPY --from=builder /app/config/config-docker.yaml /root/config/config-docker.yaml

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./im-server"]
# docker build -t im-server .

# docker run -d --name im-server --network im-network -p 8080:8080 im-server
# docker network inspect im-network
# docker stop im-server mysql-im redis
# docker network connect im-network im-server
# docker network connect im-network mysql-im
# docker network connect im-network redis
# docker start mysql-im redis im-server