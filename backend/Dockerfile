# 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /build

# 安装依赖
RUN apk add --no-cache git

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o x-store-backend ./cmd/main.go

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /build/x-store-backend .
COPY --from=builder /build/config.yaml .

# 设置时区
ENV TZ=Asia/Shanghai

EXPOSE 8082

CMD ["./x-store-backend"]
