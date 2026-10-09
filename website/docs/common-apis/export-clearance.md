---
sidebar_position: 5
title: 清关发票
slug: /common-apis/export-clearance
---

# 清关发票

按行 id 下载清关发票（仅清关，不含掏箱）。这是文件下载，响应体是 xlsx，不是 JSON。

| 项 | 值 |
| --- | --- |
| 方法 | `GET` |
| 路径 | `/api/admin/excel-exports/{tableName}` |
| 认证 | `Authorization: Bearer <JWT>` |
| `type` | `invoice_clearance_only` |

示例（行 id `896`）：

```http
GET https://ssc.mengfei.tech/api/admin/excel-exports/dynamic_settlement_statement_suqian?ids=896&type=invoice_clearance_only
Authorization: Bearer <JWT>
```

## 查询参数

| 参数 | 说明 |
| --- | --- |
| `ids` | 动态表行 id。单行：`ids=896`。多行重复参数：`ids=896&ids=897`。 |
| `type` | 固定 `invoice_clearance_only`。不要写成 `invoice_clearance`，后者是「清关-掏箱」另一种导出。 |

成功时响应头 `Content-Disposition` 为附件。下载文件名形如 `导出发票-清关_2026_10_09_14_30_00.xlsx`。请把响应保存为文件。出错时才返回 JSON，HTTP `500`。

```bash
curl -L -o clearance-invoice.xlsx \
  -H 'Authorization: Bearer <JWT>' \
  'https://ssc.mengfei.tech/api/admin/excel-exports/dynamic_settlement_statement_suqian?ids=896&type=invoice_clearance_only'
```

同一路径上的其他常用类型：[短驳发票](/common-apis/export-short-haul)、[拆箱发票](/common-apis/export-unpacking)、[运费发票](/common-apis/export-freight)。

## MCP

批量导出下载还没有 MCP 工具。以后若包装，应请求本页的 `GET`，`type` 使用 `invoice_clearance_only`，并原样带回 xlsx。接入见 [MCP 接入与使用](/mcp)。
