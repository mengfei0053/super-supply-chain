---
sidebar_position: 5
title: API 与 MCP
slug: api-and-mcp
---

# API 与 MCP

JSON API 在站点源站的 `/api` 下，**不是** `/super-supply-chain` 下面。`/super-supply-chain/` 只提供管理后台的静态页面。

生产源站示例：`https://ssc.mengfei.tech`。本地则是后端监听的地址，例如 `http://localhost:8081`。

怎么把这套 API 接到 Cursor，见独立说明 [MCP 接入与使用](/mcp)。本页列出接口范围，以及 `ssc-mcp` 已经包好的 5 个常用工具。

## 认证

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/api/register` | 注册。MCP 未包装此接口。 |
| POST | `/api/login` | 登录，返回 JWT 和用户信息。MCP 不提供登录工具；设置 `SSC_TOKEN`，或设置 `SSC_USERNAME` 与 `SSC_PASSWORD` 让进程自己调用这里。 |

其余业务接口在 `/api/admin`，请求头为 `Authorization: Bearer <token>`。JWT 为 HS256，有效期 24 小时。也可以改用个人访问令牌：在管理后台创建，调用时仍用 `Authorization: Bearer`，也接受 `X-API-Key`，以及不推荐的查询参数 `token`。见 [个人访问令牌](/personal-access-tokens)。

登录后的主要资源包括结算单、Excel 读取规则、字典、动态 Excel 表、导出规则、公司基础表，以及菜单和选项接口。动态表列表沿用 React Admin 的 `range`、`sort`、`filter` 和 `Content-Range`。`range` 为 `[start, end]` 或重复的 `range=start&range=end`，`end` 是开区间端点。`ssc_list_excel` 不传 `range` 时发送 `range=0&range=50`。

完整路径表在仓库的 [`ssc-mcp/API_CATALOG.md`](https://github.com/mengfei0053/super-supply-chain/blob/master/ssc-mcp/API_CATALOG.md)，路由注册在 `backend/main.go`。

## MCP 是什么

`ssc-mcp/` 把下面的常用接口包成 MCP 工具，供 Cursor 等客户端调用。默认走 stdio；设置 `SSC_MCP_TRANSPORT=http` 后，Streamable HTTP 在 `/mcp`，旧版 SSE 在 `/sse`。

配置时注意：

- `SSC_BASE_URL` 填源站，例如 `https://ssc.mengfei.tech` 或 `http://localhost:8081`。不要带 `/super-supply-chain`。
- 认证用 `SSC_TOKEN`，或用 `SSC_USERNAME` 加 `SSC_PASSWORD` 让服务自己登录。
- 把 `.env.example` 复制为本地 `.env`。不要把 token 或密码提交进仓库。

安装、两种传输、Cursor 的 `mcp.json` 示例和冒烟命令都写在 [MCP 接入与使用](/mcp)。

## 工具一览 \{#tools}

共 5 个工具，都是常用接口。说明根据 `ssc-mcp/README.md` 与 `ssc-mcp/src/server.ts` 里的注册信息整理。

| 工具 | 说明 |
| --- | --- |
| `ssc_upload_excel` | 常用。`POST /api/admin/excel/{tableName}`，multipart 字段 `file` 与 `name`。对应 [上传 Excel](/common-apis/upload)。 |
| `ssc_delete_excel_row` | 常用。`DELETE /api/admin/excel/{tableName}/{id}`。对应 [删除一行](/common-apis/delete-row)。 |
| `ssc_export_excel` | 常用。`GET /api/admin/excel-exports/{tableName}?ids={ids}&type={type}`。`type` 只能是 `shortHaulInvoice`、`invoice_unpacking`、`invoice_clearance_only`、`invoice_freight`，分别对应 [短驳发票](/common-apis/export-short-haul)、[拆箱发票](/common-apis/export-unpacking)、[清关发票](/common-apis/export-clearance)、[运费发票](/common-apis/export-freight)。工具把 xlsx 写到临时路径并返回该路径。 |
| `ssc_list_excel` | 常用。`GET /api/admin/excel/{tableName}`，`filter` 为 `{"start","end"}`，可选 `sort`，`range` 重复传 `range=0&range=50`。对应 [查询列表](/common-apis/list)。 |
| `ssc_search_companies` | 常用。`GET /api/admin/companies?keyword=`，在 `name` 与 `alias` 上做包含匹配，返回全部列。对应 [公司查询](/common-apis/companies)。 |

表名示例：`dynamic_settlement_statement_suqian`。登录、菜单、结算单、字典和读取规则没有 MCP 工具。

## 尚未包装的能力 \{#not-wrapped}

下列接口在后端存在，MCP **没有**对应工具：

| 能力 | HTTP | 说明 |
| --- | --- | --- |
| 结算单文件上传 | `POST /api/admin/settlement-form-entries` | multipart 文件 |
| 注册账号 | `POST /api/register` | 未包装 |

管理后台的 Excel 上传与 `ssc_upload_excel` 走同一条 `POST`。服务端按 `UPLOAD_SERVER` 写入本机目录（`file://` 或 `local://`，生产示例为 `file:///data/ssc-uploads`），不再依赖从部署机访问不到的家里 NAS。Cursor 配置示例见 [MCP 接入与使用](/mcp)。
