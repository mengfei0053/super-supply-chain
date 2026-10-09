---
id: mcp
sidebar_position: 1
title: MCP 接入与使用
slug: /mcp
---

# MCP 接入与使用

`ssc-mcp` 把 Super Supply Chain 里常用的 Excel 上传、列表、删除、四类发票导出，以及公司查询包成 MCP 工具，给 Cursor 和其他 Agent 客户端调用。它不另写一套业务接口，请求都转到源站的 `/api` 与 `/api/admin`。

管理后台在 `https://ssc.mengfei.tech/super-supply-chain/`，那只是静态页面。JSON API 与后台路径并列，不在 `/super-supply-chain` 下面。

五个工具的说明见 [API 与 MCP](/api-and-mcp#tools)。仓库里的英文说明在 `ssc-mcp/README.md`，路径表在 `ssc-mcp/API_CATALOG.md`。

## 连接

接入只使用下面两个公网地址。

| 传输 | 地址 | 方法 |
| --- | --- | --- |
| Streamable HTTP（用这个） | `https://ssc.mengfei.tech/mcp` | `POST`、`GET`、`DELETE` |
| 旧版 SSE（协议 2024-11-05） | `https://ssc.mengfei.tech/sse` | `GET` |

只支持 SSE 的客户端在拿到 `sessionId` 之后，把后续消息 `POST` 到 `https://ssc.mengfei.tech/messages?sessionId=...`。Cursor 这类只填一个 URL 的客户端填 `https://ssc.mengfei.tech/sse` 即可，由客户端按协议访问 `/messages`。

先在管理后台侧栏 **工具 → 个人访问令牌** 创建一把令牌（明文只显示一次），再替换下面的 `ssc_pat_你的令牌`。不要把真实令牌提交进仓库、聊天或本文。创建、权限范围和撤销见 [个人访问令牌](/personal-access-tokens)。

### Cursor：Streamable HTTP

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

### Cursor：旧版 SSE

```json
{
  "mcpServers": {
    "ssc": {
      "url": "https://ssc.mengfei.tech/sse",
      "headers": {
        "Authorization": "Bearer ssc_pat_你的令牌"
      }
    }
  }
}
```

`Authorization` 的值要带 `Bearer ` 前缀（含空格）。也可以改用 `X-API-Key`，这时不要加 `Bearer `，值仍是同一串令牌。`/mcp` 和 `/sse` 只换 `url`：

```json
{
  "mcpServers": {
    "ssc": {
      "url": "https://ssc.mengfei.tech/mcp",
      "headers": {
        "X-API-Key": "ssc_pat_你的令牌"
      }
    }
  }
}
```

## 认证

凭证仍是个人访问令牌（PAT）或登录 JWT。不要把真实令牌写进仓库、聊天或本文。

- **写进 Cursor。** 上面的配置要同时有 `url` 和 `headers`。优先 `"Authorization": "Bearer ssc_pat_你的令牌"`。也可以改用 `"X-API-Key": "ssc_pat_你的令牌"`（不要 `Bearer ` 前缀）。
- **配在服务端。** 生产上的 MCP 服务读取 `SSC_TOKEN`（PAT 或 JWT），或 `SSC_USERNAME` 与 `SSC_PASSWORD`。它访问源站时带请求头 `Authorization: Bearer`。服务端环境变量不能代替客户端的 `headers`。

JWT 为 HS256，有效期 24 小时。个人访问令牌在撤销或到期前一直有效，创建方式见 [个人访问令牌](/personal-access-tokens)。没有单独的登录工具。

## 工具概览

服务只注册 **5** 个工具，都是常用接口。说明见 [工具一览](/api-and-mcp#tools)。

| 常用接口 | 工具 |
| --- | --- |
| [上传 Excel](/common-apis/upload) | `ssc_upload_excel` |
| [查询列表](/common-apis/list) | `ssc_list_excel` |
| [删除一行](/common-apis/delete-row) | `ssc_delete_excel_row` |
| [短驳发票](/common-apis/export-short-haul)、[拆箱发票](/common-apis/export-unpacking)、[清关发票](/common-apis/export-clearance)、[运费发票](/common-apis/export-freight) | `ssc_export_excel`（`type` 分别为 `shortHaulInvoice`、`invoice_unpacking`、`invoice_clearance_only`、`invoice_freight`） |
| [公司查询](/common-apis/companies) | `ssc_search_companies` |

导出工具把 xlsx 写到临时文件，并在结果里返回路径。列表和公司查询返回 JSON 行。登录、菜单、结算单、字典和读取规则没有对应工具。

尚未包装的还有：结算单文件上传（`POST /api/admin/settlement-form-entries`）和注册账号（`POST /api/register`）。

接上之后，用一次只读的 `ssc_export_excel`（例如宿迁表、已有 id、`type` 为 `shortHaulInvoice`）确认源站和凭证是通的。成功时结果里有 xlsx 路径。

## 常见问题

:::caution 不要连到后台页面路径
MCP 地址是 `https://ssc.mengfei.tech/mcp` 或 `https://ssc.mengfei.tech/sse`。`/super-supply-chain` 只提供管理后台页面，不是 MCP，也不是 JSON API。
:::

**不要把令牌写进 URL。** 令牌放在 `headers` 里：`"Authorization": "Bearer ssc_pat_你的令牌"`，或 `"X-API-Key": "ssc_pat_你的令牌"`（不要 `Bearer ` 前缀）。在管理后台 **工具 → 个人访问令牌** 创建，见 [个人访问令牌](/personal-access-tokens)。

**网页上传和 MCP 上传走同一条接口。** 生产机访问不到家里局域网的 NAS（WebDAV）。当前部署把 `UPLOAD_SERVER` 设为 `file://` 或 `local://`（例如 `file:///data/ssc-uploads`），Excel 上传写到本机目录或对应的 Docker 卷。`ssc_upload_excel` 调用的也是这条 `POST /api/admin/excel/{tableName}`。结算单文件上传（`POST /api/admin/settlement-form-entries`）仍然没有 MCP 工具。
