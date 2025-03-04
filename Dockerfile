# 第一阶段：构建阶段
FROM registry.cn-hangzhou.aliyuncs.com/mengfei0053/node:22-alpine AS builder

WORKDIR /usr/src/app


COPY ./frontend .
RUN yarn
RUN npm run build

FROM golang:1.23.6-alpine

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/dist ./frontend/dist
COPY . .

WORKDIR /usr/src/app/backend

RUN go mod download

RUN go build -o app .

EXPOSE 8081

ENV ENVIRONMENT="production"

CMD ["./app"]