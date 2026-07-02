# 多阶段构建：第一阶段编译，第二阶段产出精简镜像
FROM golang:1.22-alpine AS builder

WORKDIR /build

# 复制 go.mod / go.sum 并下载依赖
COPY go.mod go.sum* ./
RUN go mod download

# 复制源码
COPY . .

# CGO_ENABLED=0 产出纯静态二进制，可在 scratch 镜像运行
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o smart-learning-server ./cmd/server

# 运行阶段：scratch 镜像（最小攻击面）
FROM scratch AS runtime

# 时区与证书
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 拷贝二进制
COPY --from=builder /build/smart-learning-server /smart-learning-server

EXPOSE 8080
ENTRYPOINT ["/smart-learning-server"]