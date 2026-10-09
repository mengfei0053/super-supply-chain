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

成功时响应头 `Content-Disposition` 为附件。下载文件名形如 `导出发票-清关_2026_10_09_14_30_00.xlsx`。请按这个文件名把响应保存为 `.xlsx`。`Content-Type` 可能是 `application/zip`（Office Open XML 本身是 zip），不要因此改成 `.zip`。出错时才返回 JSON，HTTP `500`。

```bash
curl -L -o clearance-invoice.xlsx \
  -H 'Authorization: Bearer <JWT>' \
  'https://ssc.mengfei.tech/api/admin/excel-exports/dynamic_settlement_statement_suqian?ids=896&type=invoice_clearance_only'
```

同一路径上的其他常用类型：[短驳发票](/common-apis/export-short-haul)、[拆箱发票](/common-apis/export-unpacking)、[运费发票](/common-apis/export-freight)。

## MCP

对应工具 `ssc_export_excel`，`type` 为 `invoice_clearance_only`。工具请求本页的 `GET`，把 xlsx 写到临时文件并返回路径。接入见 [MCP 接入与使用](/mcp)。
