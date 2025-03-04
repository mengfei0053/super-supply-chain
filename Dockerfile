# 第一阶段：构建阶段
FROM node:22 AS builder

WORKDIR /usr/src/app


COPY ./frontend .
RUN yarn
RUN npm run build

FROM golang

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/dist ./frontend/dist
COPY . .

WORKDIR /usr/src/app/backend

RUN go mod download

RUN go build -o app .

EXPOSE 8081

ENV ENVIRONMENT="production"

WORKDIR /usr/src/app

CMD ["ls -al ./", "ls -al ./configs"]