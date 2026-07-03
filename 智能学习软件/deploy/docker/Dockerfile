# === 调试版 Dockerfile — 智能学习软件后端 ===
# 用途：调试环境，保留调试符号，允许 delve/gdb 附加
# 构建：docker build -f deploy/docker/Dockerfile -t smart-learning-backend:debug ../../backend

FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -gcflags="all=-N -l" -o /app/server ./cmd/server

# Runtime stage — 保留完整工具链便于调试
FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata curl
COPY --from=builder /app/server /server
EXPOSE 8080 2345
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD curl -sf http://localhost:8080/health || exit 1
ENTRYPOINT ["/server"]