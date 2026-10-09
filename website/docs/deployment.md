---
sidebar_position: 6
title: 部署
slug: deployment
---

# 部署与路径

管理后台已经部署在路径前缀 `/super-supply-chain/`。文档站选择挂在**主机根路径** `/docs/`，而不是 `/super-supply-chain/docs/`。

这样生产地址是：

- `https://<主机>/docs/`
- `https://<主机>/docs/intro/`

选择根路径 `/docs/` 有三个原因：

1. 前端 Vite `base` 固定为 `/super-supply-chain`。文档若再套一层此前缀，静态资源和路由都要跟随后台一起改。
2. API 已经在 `/api`。文档、后台、接口三者并列，互不嵌套。
3. 管理后台顶栏使用同源绝对地址 `/docs/`。无论当前页面在 `/super-supply-chain/#/...` 的哪个 hash 路由，点击都会离开单页应用，打开文档站。

## 镜像怎么打

根目录 `Dockerfile` 分三段：

1. 用 Node 构建 `frontend/`，得到 `frontend/dist`。
2. 用 Node 在 `website/` 执行 `npm ci` 和 `npm run build`，得到 `website/build`。
3. 用 Go 镜像编译 `backend`，并把上面两份静态产物复制进最终镜像。

生产进程 `ENVIRONMENT=production` 时，`backend/controllers/static.go` 会：

- 把 `frontend/dist` 挂到 `/super-supply-chain`
- 把 `website/build` 挂到 `/docs`

Docusaurus 打开了 `trailingSlash: true`，每个页面是目录下的 `index.html`。Go 对不带斜杠的目录路径做 301，例如 `/docs/intro` 转到 `/docs/intro/`。

本地开发（`ENVIRONMENT` 不是 `production`）不走这段静态托管。前端用 Vite，文档用 `npm start`。

## Kubernetes

`ssc-deployment.yaml` 把容器发布到命名空间 `ssc`。运行时环境来自 `PORT`、`ENVIRONMENT=production` 和 Secret `ssc-secret`。更新镜像标签后，同一 Service 会同时提供后台、API 和 `/docs/`。

GitHub Actions 工作流 `.github/workflows/docker-build-push.yml` 只在 `v*` tag 上构建并推送镜像。

## MCP HTTP（生产）

Compose 服务 `ssc-mcp` 监听容器内 `3100`，只发布到 `172.17.0.1:3100`。Nginx Proxy Manager 在 `ssc.mengfei.tech` 上把 `/mcp`、`/sse`、`/messages` 转到该地址；`/api`、`/docs`、`/super-supply-chain` 仍走 `172.17.0.1:8088`。

Cursor 填完整远程配置：`url` 加上个人访问令牌请求头。令牌在管理后台 **工具 → 个人访问令牌** 创建，下面只是占位，不要写真实令牌。`/sse` 与可选的 `X-API-Key` 见 [MCP 接入与使用](/mcp)。

```json
{
  "mcpServers": {
    "ssc": {
      "url": "https://ssc.mengfei.tech/mcp",
      "headers": {
        "Authorization": "Bearer ssc_pat_你的令牌"
      }
    }
  }
}
```

自己在服务器上跑 `ssc-mcp` 时，进程仍用 `.env.compose` 里的 `SSC_TOKEN`（或 `SSC_USERNAME` + `SSC_PASSWORD`）去调用源站。那不能代替客户端的 `headers`。

