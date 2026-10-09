---
sidebar_position: 2
title: 本地运行
slug: local-development
---

# 本地运行

先准备 Node.js 20+、Go 1.23 和 MySQL。密钥只放在本地环境变量或 `configs/.env` 里，不要提交到仓库。

## 后端

在仓库根目录执行：

```sh
./dev.sh
```

脚本会设置 `ENVIRONMENT=development`，进入 `backend/`，并用 [air](https://github.com/air-verse/air) 热重载。需要本机已安装 `air`。

开发模式下后端**不**托管前端和文档静态文件。接口默认在 `http://localhost:8081`。

后端会读取这些环境变量：

- `ENVIRONMENT`
- `PORT`
- `MYSQL_USER`
- `MYSQL_PASSWORD`
- `MYSQL_SERVER`
- `UPLOAD_USER`
- `UPLOAD_PASSWORD`
- `UPLOAD_SERVER`

开发环境从 `configs/.env` 加载（相对 `backend/` 工作目录的 `../configs/.env`）。生产环境从进程或 Kubernetes Secret 读取。

不用热重载时，可以在 `backend/` 里执行 `go build ./...`。

## 前端

```sh
cd frontend
yarn
yarn dev
```

仓库里有 `frontend/yarn.lock`，请用 yarn。Vite 的 `base` 是 `/super-supply-chain`，开发地址类似 `http://localhost:5173/super-supply-chain/`。

`vite.config.ts` 会把 `/api` 代理到 `http://localhost:8081`，并把 `/docs` 代理到本机文档站 `http://127.0.0.1:3000`。文档站没启动时，顶栏「文档」会打不开，先按下一节把文档站跑起来。

常用检查：

```sh
yarn type-check
yarn build
```

## 文档站

```sh
cd website
npm install
npm start
```

开发服务器地址是 [http://localhost:3000/docs/](http://localhost:3000/docs/)。站点配置了 `baseUrl: '/docs/'`，所以本地也带这个前缀。

生产构建：

```sh
npm run build
npm run serve
```

`npm run build` 把静态文件生成到 `website/build/`。`npm run serve` 用来预览这次构建。

## 命令行

```sh
cd cli
go run . login --base-url http://localhost:8081 --username <账号> --password <密码>
go run . status
```

也可以用环境变量 `SSC_BASE_URL`、`SSC_USERNAME`、`SSC_PASSWORD`。说明见仓库 `cli/README.md`。
