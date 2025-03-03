FROM alpine:latest

# 安装 Node.js
RUN apk add --no-cache curl tar && \
    curl -fsSL https://unofficial-builds.nodejs.org/download/release/v18.17.1/node-v18.17.1-linux-x64-musl.tar.gz | tar -xz -C /usr/local --strip-components=1 && \
    apk del curl tar

# 安装 Go
RUN apk add --no-cache bash git gcc musl-dev && \
    wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz && \
    rm go1.24.0.linux-amd64.tar.gz && \
    export PATH=$PATH:/usr/local/go/bin

# 设置工作目录
WORKDIR /app

# 复制应用程序代码
COPY . .
