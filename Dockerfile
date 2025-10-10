# =========================
# 构建阶段
# =========================
FROM golang:1.24-alpine AS builder

# 在容器根目录构建
WORKDIR /src

RUN apk add --no-cache git ca-certificates

# 拷贝依赖文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 拷贝整个项目
COPY . .

# 编译产物到 bin/
RUN go build -o bin/gitinsight ./cmd/gitinsight/

# =========================
# 生产阶段
# =========================
FROM alpine:3.18

# ✅ 设置镜像工作目录为 /app
WORKDIR /app

RUN apk add --no-cache ca-certificates

# 从构建阶段复制二进制文件
COPY --from=builder /src/bin/gitinsight /app/gitinsight

CMD ["./gitinsight", "serv"]
