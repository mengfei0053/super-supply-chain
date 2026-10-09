---
sidebar_position: 3
title: MCP 接入与使用
slug: mcp
---

# MCP 接入与使用

`ssc-mcp` 把 Super Supply Chain 的 HTTP API 包成 MCP 工具，给 Cursor 和其他 Agent 客户端调用。它不另写一套业务接口：登录、菜单、结算单、字典、Excel 读取规则、动态表行和导出规则，都转发到源站的 `/api` 与 `/api/admin`。

生产源站是 `https://ssc.mengfei.tech`。管理后台在 `https://ssc.mengfei.tech/super-supply-chain/`，那只是静态页面。JSON API 与后台路径并列，不在 `/super-supply-chain` 下面。

20 个工具的分组说明见 [API 与 MCP](/api-and-mcp#tools)。仓库里的英文说明在 `ssc-mcp/README.md`，路径表在 `ssc-mcp/API_CATALOG.md`。

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
2. **用户名和密码。** 设置 `SSC_USERNAME` 与 `SSC_PASSWORD`。进程在第一次需要登录的调用上请求 `POST /api/login`，把 JWT 放在内存里，不打印出来。也可以稍后调用工具 `ssc_login`。

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

```bash
export SSC_BASE_URL=https://ssc.mengfei.tech
export SSC_TOKEN='<jwt>'
export SSC_MCP_TRANSPORT=http
export SSC_MCP_HOST=0.0.0.0
export SSC_MCP_PORT=3100
npm run start
```

也可以：`npm run dev -- --transport http --port 3100`。

| 传输 | 地址 | 方法 |
| --- | --- | --- |
| Streamable HTTP（用这个） | `http://127.0.0.1:3100/mcp` | `POST`、`GET`、`DELETE` |
| 旧版 SSE | `http://127.0.0.1:3100/sse` | `GET` |
| 旧版 SSE 消息 | `http://127.0.0.1:3100/messages?sessionId=...` | `POST` |

进程默认听在 `0.0.0.0`。客户端要填 `127.0.0.1` 或其他能访问到的地址，不要填 `0.0.0.0`。

会话保存在内存中。客户端发送 `DELETE /mcp`、旧版 SSE 连接关闭，或 30 分钟内没有进行中的请求时，会话会被清掉。最多接受 200 个会话。

桌面版 Cursor 直接请求 URL，不需要 CORS。浏览器里的 MCP 客户端需要。服务会响应 `OPTIONS`，并设置：

- `Access-Control-Allow-Origin`：默认 `*`；请求的 `Origin` 落在 `SSC_MCP_CORS_ORIGIN`（逗号分隔）里时，回显该来源
- `Access-Control-Allow-Methods`：`GET, POST, DELETE, OPTIONS`
- `Access-Control-Allow-Headers`：`Content-Type`、`Accept`、`Authorization`、`Mcp-Session-Id`、`Mcp-Protocol-Version`、`Last-Event-ID`
- `Access-Control-Expose-Headers`：`Mcp-Session-Id`、`Mcp-Protocol-Version`

`*` 不能和带凭证的浏览器请求一起用。要限制到某个站点，设置 `SSC_MCP_CORS_ORIGIN=https://your-app.example`。

监听 `0.0.0.0` 或 `::` 且没有 `SSC_MCP_ALLOWED_HOSTS` 时，SDK 会关闭 DNS 重绑定检查并打印警告。绑定 `127.0.0.1`、`localhost` 或 `::1` 时会自动打开这项检查。若要听全部网卡又限制 `Host`，设置 `SSC_MCP_ALLOWED_HOSTS=localhost,127.0.0.1`。

### Cursor：HTTP URL

先自己把服务跑起来，再把 Cursor 指到 Streamable HTTP：

```json
{
  "mcpServers": {
    "ssc": {
      "url": "http://127.0.0.1:3100/mcp"
    }
  }
}
```

只支持 SSE 的旧客户端：

```json
{
  "mcpServers": {
    "ssc": {
      "url": "http://127.0.0.1:3100/sse"
    }
  }
}
```

`SSC_TOKEN`、`SSC_USERNAME`、`SSC_PASSWORD` 放在启动 MCP 的进程环境里，不要写进这段 URL 配置。

## 工具概览

服务注册了 **20** 个工具，分成六组。每一项的用途见 [工具一览](/api-and-mcp#tools)。

| 分组 | 工具 |
| --- | --- |
| 登录、状态、菜单 | `ssc_login`、`ssc_status`、`ssc_list_menus` |
| 结算单 | `ssc_list_orders`、`ssc_get_order` |
| 字典 | `ssc_list_dicts`、`ssc_get_dict`、`ssc_get_dict_map`、`ssc_create_dict`、`ssc_update_dict`、`ssc_delete_dict` |
| Excel 读取规则 | `ssc_list_excel_read_rules`、`ssc_get_excel_read_rule` |
| 动态 Excel 行 | `ssc_list_excel_rows`、`ssc_get_excel_row`、`ssc_update_excel_row`、`ssc_delete_excel_row` |
| 导出规则 | `ssc_list_export_rules`、`ssc_list_export_template_options`、`ssc_get_export_rule` |

列表类工具的 `range` 沿用 React Admin，形如 `[start, end]`，不传时默认 `[0, 49]`。动态表名来自 `ssc_list_menus` 返回的 `dynamicTableName`。

目录里记录过的动态表包括 `dynamic_Integrity_packaging_invoice`、`dynamic_customs_declaration_form`、`dynamic_settlement_statement_fenchang`、`dynamic_settlement_statement_suqian`、`dynamic_yifan_cost_cal`。以菜单接口的实时结果为准。

### 还没有包进 MCP 的接口

下面这些 HTTP 能力存在，**还没有**对应的 MCP 工具：

- 动态表 Excel 上传（`POST /api/admin/excel/:tableName`，multipart `file`）
- 结算单文件上传（`POST /api/admin/settlement-form-entries`，multipart `file`）
- 批量导出下载（`GET /api/admin/excel-exports/:tableName`，返回二进制文件）
- 注册账号（`POST /api/register`）

管理后台的 Excel 上传已经可用，文件写到服务端本机目录。MCP 上传工具仍未提供，不能用现有工具代替网页上传。

## 冒烟

需要能访问 SSC，并且已经设置 token 或用户名密码。脚本不会打印密钥。

```bash
export SSC_BASE_URL=https://ssc.mengfei.tech
export SSC_TOKEN='...'
npm run smoke
```

只检查传输、不连接业务账号时，用下面的命令。它会在临时端口上握手 Streamable HTTP 和旧版 SSE，列出工具，调用 `ssc_status`，检查 CORS，再握手默认的 stdio 服务：

```bash
npm run smoke:http
```

接入 Cursor 之后，先调用 `ssc_status`。返回里 `authenticated` 为真，并带有菜单数量，就说明源站和凭证是通的。然后用 `ssc_list_menus` 拿到 `dynamicTableName`，再查表。

## 常见问题

:::caution 源站不要带后台路径
`SSC_BASE_URL` 必须是 API 源站，例如 `https://ssc.mengfei.tech` 或本地 `http://localhost:8081`。不要写成 `https://ssc.mengfei.tech/super-supply-chain`。`/super-supply-chain` 只提供管理后台页面，请求会打到静态文件而不是 `/api`。
:::

**HTTP 的 `url` 写成了 `0.0.0.0`。** `0.0.0.0` 是进程的监听地址。Cursor 里填 `http://127.0.0.1:3100/mcp`。

**把 token 写进 HTTP 配置。** HTTP 模式下，`SSC_TOKEN` 或用户名密码属于启动 `ssc-mcp` 的那个进程。Cursor 的配置只需要 `url`。

**列表为空，或接口要求 `range`。** 字典、读取规则、导出模板等列表沿用 React Admin 的 `range`。工具不传时默认 `[0, 49]`。动态表还可以用 `filterStart`、`filterEnd` 按 `created_at` 过滤，格式 `YYYY-MM-DD`。

**网页能上传，MCP 不能传文件。** 生产机访问不到家里局域网的 NAS（WebDAV）。当前部署把 `UPLOAD_SERVER` 设为 `file://` 或 `local://`（例如 `file:///data/ssc-uploads`），管理后台的 Excel 上传写到本机目录或对应的 Docker 卷。改回一个可达的 WebDAV 地址才会回到 NAS 模式。这条存储切换只影响网站上传；MCP 仍然没有上传、结算单传文件或批量导出下载工具。查行、改行、删行可以，传文件请用管理后台。
