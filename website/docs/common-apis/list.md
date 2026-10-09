---
sidebar_position: 8
title: 查询列表
slug: /common-apis/list
---

# 查询列表

按表名列出动态 Excel 行。管理后台打开一张表时走的就是这条接口。结果按 `created_at` 落在 `filter` 的日期窗口里筛选，再用 `range` 分页。

| 项 | 值 |
| --- | --- |
| 方法 | `GET` |
| 路径 | `/api/admin/excel/{tableName}` |
| 认证 | `Authorization: Bearer <token>` |

`<token>` 与本分组其他接口相同：登录 JWT，或个人访问令牌。

示例（宿迁结算表，当天创建的前 50 行）：

```http
GET https://ssc.mengfei.tech/api/admin/excel/dynamic_settlement_statement_suqian?filter={"start":"2026-10-09","end":"2026-10-09"}&sort={"field":"id","order":"ASC"}&range=0&range=50
Authorization: Bearer <token>
```

## 查询参数

| 参数 | 说明 |
| --- | --- |
| `filter` | JSON 对象，例如 `{"start":"2026-10-09","end":"2026-10-09"}`。服务端按 `created_at` 在 `start 00:00:00` 与 `end 23:59:59` 之间过滤。宿迁结算表不带有效起止日期时，窗口对不上真实时间，结果是空的。也可以写成 `filter.start` 与 `filter.end` 两个参数。 |
| `sort` | 可选。管理后台发送 JSON 对象 `{"field":"id","order":"ASC"}`。部分客户端发送数组 `["id","ASC"]`。两种形式接口都接受。当前列表 SQL 只做日期过滤和分页，不会按 `sort` 做 `ORDER BY`。 |
| `range` | 分页。浏览器实际发送的是重复键 `range=0&range=50`：起点是偏移，终点是开区间，条数 = 终点 − 起点。`range=[0,50]`（URL 编码后的 JSON）同样可用。缺 `range` 时接口返回 400。 |

```bash
curl -G \
  -H 'Authorization: Bearer <token>' \
  --data-urlencode 'filter={"start":"2026-10-09","end":"2026-10-09"}' \
  --data-urlencode 'sort={"field":"id","order":"ASC"}' \
  --data-urlencode 'range=0' \
  --data-urlencode 'range=50' \
  'https://ssc.mengfei.tech/api/admin/excel/dynamic_settlement_statement_suqian'
```

成功时 HTTP `200`，响应体是 JSON 数组。每一行包含 `id`、`fileName`、`datas`（其中有 `baseData` 和 `list`），以及 `uploadFilePath`、`nasFileName` 等字段。响应头 `Content-Range` 可能出现，值是匹配总数的十进制字符串，不是字节范围。

```json
[
  {
    "id": 896,
    "fileName": "statement.xlsx",
    "datas": {
      "baseData": {},
      "list": []
    }
  }
]
```

## MCP

对应工具 `ssc_list_excel`。参数是 `tableName`、`filterStart`、`filterEnd`，可选 `sort` 和 `range`。工具请求本页的 `GET`：`filter` 用 JSON 对象，`range` 用重复的 `range=0&range=50`（不传时默认这一页）。工具说明见 [API 与 MCP](/api-and-mcp#tools)。接入见 [MCP 接入与使用](/mcp)。
