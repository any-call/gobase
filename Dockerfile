# 使用 golang 官方镜像
#FROM  golang:latest
FROM  --platform=linux/amd64  golang:latest
#FROM golang:1.20-buster

# 安装 Leptonica 和 Tesseract 依赖
RUN apt-get update && apt-get install -y \
    libtesseract-dev \
    libleptonica-dev \
    tesseract-ocr \
    pkg-config \
    gcc \
    g++ \
    make \
    libc-dev \
    g++-x86-64-linux-gnu \
    libc6-dev-amd64-cross \
    build-essential

# 设置工作目录
WORKDIR /app

# 将当前目录挂载到容器的 /app 目录
VOLUME /app

# 设置 Go 环境
ENV CGO_ENABLED=1
#ENV GOOS=linux
#ENV GOARCH=amd64
#ENV GOARCH=arm64
# 默认命令
#CMD ["go", "build", "./cmd/client"]
#CMD GOOS=linux GOARCH=amd64 go build ./cmd/client