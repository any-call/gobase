# 使用 golang 官方镜像
FROM golang:latest

# 安装 Leptonica 和 Tesseract 依赖
RUN apt-get update && apt-get install -y \
    libleptonica-dev \
    tesseract-ocr \
    pkg-config \
    build-essential

# 设置工作目录
WORKDIR /app

# 将当前目录挂载到容器的 /app 目录
VOLUME /app

# 设置 Go 环境
ENV CGO_ENABLED=1

# 默认命令
CMD ["go", "build", "./cmd/client"]