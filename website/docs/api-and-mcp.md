---
sidebar_position: 3
title: API 与 MCP
slug: api-and-mcp
---

# API 与 MCP

JSON API 在站点源站的 `/api` 下，**不是** `/super-supply-chain` 下面。`/super-supply-chain/` 只提供管理后台的静态页面。

生产源站示例：`https://ssc.mengfei.tech`。本地则是后端监听的地址，例如 `http://localhost:8081`。

## 认证

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/api/register` | 注册 |
| POST | `/api/login` | 登录，返回 JWT 和用户信息 |

其余业务接口在 `/api/admin`，请求头为 `Authorization: Bearer <token>`。

登录后的主要资源包括结算单、Excel 读取规则、字典、动态 Excel 表、导出规则，以及菜单和选项接口。列表接口沿用 React Admin 的 `range`、`sort`、`filter` 和 `Content-Range`。

完整路径表在仓库的 [`ssc-mcp/API_CATALOG.md`](https://github.com/mengfei0053/super-supply-chain/blob/master/ssc-mcp/API_CATALOG.md)，路由注册在 `backend/main.go`。这里不重复整张表。

## MCP

`ssc-mcp/` 把上述 HTTP API 包成 MCP 工具（登录、菜单、订单、字典、Excel 行、导出规则等）。安装和工具列表见 [`ssc-mcp/README.md`](https://github.com/mengfei0053/super-supply-chain/blob/master/ssc-mcp/README.md)。

配置时注意：

- `SSC_BASE_URL` 填源站，例如 `https://ssc.mengfei.tech` 或 `http://localhost:8081`。不要带 `/super-supply-chain`。
- 认证用 `SSC_TOKEN`，或用 `SSC_USERNAME` 加 `SSC_PASSWORD` 让服务自己登录。
- 把 `.env.example` 复制为本地 `.env`。不要把 token 或密码提交进仓库。

未包装的能力主要是文件上传、批量导出下载这类二进制接口。
