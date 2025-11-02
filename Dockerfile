# 使用官方的Go运行时作为基础镜像
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 配置 Go 模块代理为国内源
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o stable-diffusion-webui-mcp .

# 使用轻量级的运行时镜像
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/stable-diffusion-webui-mcp .

# 暴露端口
EXPOSE 18080

# 环境变量
ENV SDWEBUI_URL="http://127.0.0.1:7860"

# 运行应用，使用固定端口
CMD ["./stable-diffusion-webui-mcp", "-port", ":18080", "-sdwebui-url", "${SDWEBUI_URL}"]