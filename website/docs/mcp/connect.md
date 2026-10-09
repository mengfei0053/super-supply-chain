---
id: mcp
sidebar_position: 1
title: MCP 接入与使用
slug: /mcp
---

# MCP 接入与使用

`ssc-mcp` 把 Super Supply Chain 里常用的 Excel 上传、删除和四类发票导出包成 MCP 工具，给 Cursor 和其他 Agent 客户端调用。它不另写一套业务接口，请求都转到源站的 `/api` 与 `/api/admin`。

生产源站是 `https://ssc.mengfei.tech`。管理后台在 `https://ssc.mengfei.tech/super-supply-chain/`，那只是静态页面。JSON API 与后台路径并列，不在 `/super-supply-chain` 下面。

三个工具的说明见 [API 与 MCP](/api-and-mcp#tools)。仓库里的英文说明在 `ssc-mcp/README.md`，路径表在 `ssc-mcp/API_CATALOG.md`。

## 安装与构建

需要 Node.js 20+。在仓库的 `ssc-mcp/` 目录执行：

```bash
cd ssc-mcp
npm install
npm run build
```

`npm run build` 生成 `dist/index.js`。stdio 模式用这个文件。还没构建时，也可以用 `npx tsx src/index.ts` 直接跑源码。

## 环境变量

把 `ssc-mcp/.env.example` 复制为 `ssc-mcp/.env`，权限建议 `600`，只在本机填写。不要把 token 或密码提交进仓库。

| 变量 | 是否必填 | 说明 |
| --- | --- | --- |
| `SSC_BASE_URL` | 否 | API 源站，默认 `https://ssc.mengfei.tech`。不要带 `/super-supply-chain`，也不要末尾斜杠。 |
| `SSC_TOKEN` | 与账号密码二选一 | 登录得到的 JWT，优先使用 |
| `SSC_USERNAME` | 与密码一起 | 交给服务调用 `POST /api/login` |
| `SSC_PASSWORD` | 与用户名一起 | 登录密码 |
| `SSC_MCP_TRANSPORT` | 否 | `stdio`（默认）或 `http` |
| `SSC_MCP_HOST` | 否 | HTTP 监听地址，默认 `0.0.0.0` |
| `SSC_MCP_PORT` | 否 | HTTP 端口，默认 `3100` |
| `SSC_MCP_CORS_ORIGIN` | 否 | 浏览器 CORS 允许的来源，默认 `*` |
| `SSC_MCP_ALLOWED_HOSTS` | 否 | 可选，逗号分隔的 `Host` 白名单，用来防 DNS 重绑定 |

命令行参数会覆盖同名环境变量：`--transport`、`--host`、`--port`、`--cors-origin`、`--allowed-hosts`。

## 认证

除登录以外的工具都要已有凭证。两种方式：

1. **Token（优先）。** 从浏览器登录态或仓库 `cli/` 的 `ssc login` 取出 JWT，设为 `SSC_TOKEN`。请求头是 `Authorization: Bearer <token>`。JWT 为 HS256，有效期 24 小时。
2. **用户名和密码。** 设置 `SSC_USERNAME` 与 `SSC_PASSWORD`。进程在第一次需要登录的调用上请求 `POST /api/login`，把 JWT 放在内存里，不打印出来。没有单独的登录工具。

同一个操作系统进程里的全部会话共用这一次登录。不同账号请各起一个进程。HTTP 模式不要把端口暴露到不受信任的网络。

仓库不保存密码。`ssc-mcp/README.md` 里提到的现有账号名是 `testuser`。

## stdio（默认）

Cursor 在本机拉起进程时用 stdio。不设置 `SSC_MCP_TRANSPORT`，或显式设为 `stdio`。

```bash
export SSC_BASE_URL=https://ssc.mengfei.tech
export SSC_TOKEN='<jwt>'
npm run start
```

### Cursor：本地命令

路径换成你本机克隆目录里的 `ssc-mcp/dist/index.js`。

```json
{
  "mcpServers": {
    "ssc": {
      "command": "node",
      "args": ["/path/to/super-supply-chain/ssc-mcp/dist/index.js"],
      "env": {
        "SSC_BASE_URL": "https://ssc.mengfei.tech",
        "SSC_TOKEN": "<paste JWT here>"
      }
    }
  }
}
```

用账号密码、且尚未构建时：

```json
{
  "mcpServers": {
    "ssc": {
      "command": "npx",
      "args": ["tsx", "/path/to/super-supply-chain/ssc-mcp/src/index.ts"],
      "env": {
        "SSC_BASE_URL": "https://ssc.mengfei.tech",
        "SSC_USERNAME": "testuser",
        "SSC_PASSWORD": "<password>"
      }
    }
  }
}
```

## HTTP

HTTP 需要显式打开。工具和登录行为与 stdio 相同。当前协议的 Streamable HTTP 在 `/mcp`。旧客户端使用的 HTTP+SSE（协议 2024-11-05）在 `/sse`。

**生产环境已经挂到公网域名**（Nginx Proxy Manager 把 `/mcp`、`/sse`、`/messages` 转到本机 `172.17.0.1:3100` 上的 `ssc-mcp` 容器）。日常接入直接用 HTTPS，不必自己在本机起进程。

| 传输 | 生产地址 | 本机调试 | 方法 |
| --- | --- | --- | --- |
| Streamable HTTP（用这个） | `https://ssc.mengfei.tech/mcp` | `http://127.0.0.1:3100/mcp` | `POST`、`GET`、`DELETE` |
| 旧版 SSE | `https://ssc.mengfei.tech/sse` | `http://127.0.0.1:3100/sse` | `GET` |
| 旧版 SSE 消息 | `https://ssc.mengfei.tech/messages?sessionId=...` | `http://127.0.0.1:3100/messages?sessionId=...` | `POST` |

本机自己跑 HTTP 时：

```bash
export SSC_BASE_URL=https://ssc.mengfei.tech
export SSC_TOKEN='<jwt>'
export SSC_MCP_TRANSPORT=http
export SSC_MCP_HOST=0.0.0.0
export SSC_MCP_PORT=3100
npm run start
```

也可以：`npm run dev -- --transport http --port 3100`。进程默认听在 `0.0.0.0`；客户端不要填 `0.0.0.0`，填 `127.0.0.1` 或生产域名。

会话保存在内存中。客户端发送 `DELETE /mcp`、旧版 SSE 连接关闭，或 30 分钟内没有进行中的请求时，会话会被清掉。最多接受 200 个会话。

桌面版 Cursor 直接请求 URL，不需要 CORS。浏览器里的 MCP 客户端需要。服务会响应 `OPTIONS`，并设置：

- `Access-Control-Allow-Origin`：默认 `*`；请求的 `Origin` 落在 `SSC_MCP_CORS_ORIGIN`（逗号分隔）里时，回显该来源
- `Access-Control-Allow-Methods`：`GET, POST, DELETE, OPTIONS`
- `Access-Control-Allow-Headers`：`Content-Type`、`Accept`、`Authorization`、`Mcp-Session-Id`、`Mcp-Protocol-Version`、`Last-Event-ID`
- `Access-Control-Expose-Headers`：`Mcp-Session-Id`、`Mcp-Protocol-Version`

`*` 不能和带凭证的浏览器请求一起用。要限制到某个站点，设置 `SSC_MCP_CORS_ORIGIN=https://your-app.example`。

监听 `0.0.0.0` 或 `::` 且没有 `SSC_MCP_ALLOWED_HOSTS` 时，SDK 会关闭 DNS 重绑定检查并打印警告。绑定 `127.0.0.1`、`localhost` 或 `::1` 时会自动打开这项检查。若要听全部网卡又限制 `Host`，设置 `SSC_MCP_ALLOWED_HOSTS=localhost,127.0.0.1`。

### Cursor：HTTP URL（生产推荐）

生产上 MCP 已在域名后面跑着，Cursor 直接填 HTTPS：

```json
{
  "mcpServers": {
    "ssc": {
      "url": "https://ssc.mengfei.tech/mcp"
    }
  }
}
```

只支持 SSE 的旧客户端：

```json
{
  "mcpServers": {
    "ssc": {
      "url": "https://ssc.mengfei.tech/sse"
    }
  }
}
```

业务侧的 `SSC_TOKEN` / `SSC_USERNAME`+`SSC_PASSWORD` 配在 **101 上 `ssc-mcp` 容器的环境变量**（`.env.compose`）里，不要写进 Cursor 的 URL 配置。本机临时起 HTTP 调试时，才用 `http://127.0.0.1:3100/mcp`，并把凭证放进启动进程的环境变量。

## 工具概览

服务只注册 **3** 个工具，都是常用接口。说明见 [工具一览](/api-and-mcp#tools)。

| 常用接口 | 工具 |
| --- | --- |
| [上传 Excel](/common-apis/upload) | `ssc_upload_excel` |
| [删除一行](/common-apis/delete-row) | `ssc_delete_excel_row` |
| [短驳发票](/common-apis/export-short-haul)、[拆箱发票](/common-apis/export-unpacking)、[清关发票](/common-apis/export-clearance)、[运费发票](/common-apis/export-freight) | `ssc_export_excel`（`type` 分别为 `shortHaulInvoice`、`invoice_unpacking`、`invoice_clearance_only`、`invoice_freight`） |

导出工具把 xlsx 写到临时文件，并在结果里返回路径。登录、菜单、结算单、字典和读取规则没有对应工具。凭证只放在 `SSC_TOKEN` 或用户名密码环境变量里。

尚未包装的还有：结算单文件上传（`POST /api/admin/settlement-form-entries`）和注册账号（`POST /api/register`）。

## 冒烟

需要能访问 SSC，并且已经设置 token 或用户名密码。脚本不会打印密钥。

```bash
export SSC_BASE_URL=https://ssc.mengfei.tech
export SSC_TOKEN='...'
npm run smoke
```

只检查传输、不连接业务账号时，用下面的命令。它会在临时端口上握手 Streamable HTTP 和旧版 SSE，确认只有这三个工具，在没有凭证时调用删除，检查 CORS，再握手默认的 stdio 服务：

```bash
npm run smoke:http
```

接入 Cursor 之后，用一次只读的 `ssc_export_excel`（例如宿迁表、已有 id、`type` 为 `shortHaulInvoice`）确认源站和凭证是通的。成功时结果里有本地 xlsx 路径。

## 常见问题

:::caution 源站不要带后台路径
`SSC_BASE_URL` 必须是 API 源站，例如 `https://ssc.mengfei.tech` 或本地 `http://localhost:8081`。不要写成 `https://ssc.mengfei.tech/super-supply-chain`。`/super-supply-chain` 只提供管理后台页面，请求会打到静态文件而不是 `/api`。
:::

**HTTP 的 `url` 写成了 `0.0.0.0`。** `0.0.0.0` 是进程的监听地址。生产填 `https://ssc.mengfei.tech/mcp`；本机调试填 `http://127.0.0.1:3100/mcp`。

**把 token 写进 HTTP 配置。** HTTP 模式下，`SSC_TOKEN` 或用户名密码属于跑 `ssc-mcp` 的那个进程（生产是 101 上的容器）。Cursor 的配置只需要 `url`。

**网页上传和 MCP 上传走同一条接口。** 生产机访问不到家里局域网的 NAS（WebDAV）。当前部署把 `UPLOAD_SERVER` 设为 `file://` 或 `local://`（例如 `file:///data/ssc-uploads`），Excel 上传写到本机目录或对应的 Docker 卷。`ssc_upload_excel` 调用的也是这条 `POST /api/admin/excel/{tableName}`。改回一个可达的 WebDAV 地址才会回到 NAS 模式。结算单文件上传（`POST /api/admin/settlement-form-entries`）仍然没有 MCP 工具。
