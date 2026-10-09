# 文档站

Super Supply Chain 的说明站，使用 [Docusaurus](https://docusaurus.io/) 3.x。

站点 `baseUrl` 是 `/docs/`。本地和线上地址都带这个前缀，例如 `/docs/` 和 `/docs/intro/`。

管理后台部署在 `/super-supply-chain/`。文档放在主机根路径 `/docs/`，与后台、`/api` 并列。后台顶栏的「文档」链接到 `/docs/`。

## 开发

在 `website/` 目录：

```sh
npm install
npm start
```

开发服务器：<http://localhost:3000/docs/>。

也可以用 pnpm：`pnpm install` 然后 `pnpm start`。本目录提交了 `package-lock.json`，Docker 构建使用 `npm ci`。

前端开发服务器会把 `/docs` 代理到 `127.0.0.1:3000`。两个进程都开着时，后台顶栏的「文档」会打开本站。

## 构建

```sh
npm run build
```

静态文件输出到 `website/build/`。预览：

```sh
npm run serve
```

类型检查：

```sh
npm run typecheck
```

## 生产环境怎么挂上

根目录 `Dockerfile` 会构建本目录，并把 `build/` 复制进镜像的 `website/build`。`ENVIRONMENT=production` 时，Go 进程把该目录挂到 `/docs`。开发模式（`./dev.sh`）不托管这些静态文件。

路径选择和斜杠重定向见文档页「部署」。
