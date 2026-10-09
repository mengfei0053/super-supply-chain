# 第一阶段：构建前端
FROM registry.cn-hangzhou.aliyuncs.com/mengfei0053/node:22 AS frontend-builder

WORKDIR /usr/src/app

COPY ./frontend .
RUN yarn
RUN npm run build

# 第二阶段：构建文档站，产物由 Go 挂到 /docs/
FROM registry.cn-hangzhou.aliyuncs.com/mengfei0053/node:22 AS docs-builder

WORKDIR /usr/src/app

COPY website/package.json website/package-lock.json ./
RUN npm ci
COPY website/ ./
RUN npm run build

FROM registry.cn-hangzhou.aliyuncs.com/mengfei0053/golang:1.23.6-alpine

WORKDIR /usr/src/app

COPY --from=frontend-builder /usr/src/app/dist ./frontend/dist
COPY --from=docs-builder /usr/src/app/build ./website/build
COPY . .

WORKDIR /usr/src/app/backend

RUN go mod download

RUN go build -o app .

EXPOSE 8081

ENV ENVIRONMENT="production"

CMD ["./app"]
