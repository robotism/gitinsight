# 构建 Go 二进制
FROM golang:1.24-alpine AS builder

WORKDIR /

# 安装依赖工具
RUN apk add --no-cache git ca-certificates

# 拷贝代码并构建
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o gitinsight ./cmd/gitinsight/

# 生产镜像
FROM alpine:3.18
WORKDIR /
RUN apk add --no-cache ca-certificates
COPY --from=builder /gitinsight .

EXPOSE 8080
CMD ["./gitinsight", "serv"]
