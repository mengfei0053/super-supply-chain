---
sidebar_position: 3
title: 短驳发票
slug: /common-apis/export-short-haul
---

# 短驳发票

按行 id 下载短驳发票。这是文件下载，响应体是 xlsx，不是 JSON。

| 项 | 值 |
| --- | --- |
| 方法 | `GET` |
| 路径 | `/api/admin/excel-exports/{tableName}` |
| 认证 | `Authorization: Bearer <JWT>` |
| `type` | `shortHaulInvoice` |

示例（行 id `896`）：

```http
GET https://ssc.mengfei.tech/api/admin/excel-exports/dynamic_settlement_statement_suqian?ids=896&type=shortHaulInvoice
Authorization: Bearer <JWT>
```

## 查询参数

| 参数 | 说明 |
| --- | --- |
| `ids` | 动态表行 id。单行：`ids=896`。多行重复参数：`ids=896&ids=897`。服务端按 `id in (?)` 查询。 |
| `type` | 固定 `shortHaulInvoice`。 |

成功时响应头 `Content-Disposition` 为附件。下载文件名形如 `导出-短驳费-发票_2026_10_09_14_30_00.xlsx`（前缀「导出-短驳费-发票」加导出时间）。请把响应保存为文件。出错时才返回 JSON，HTTP `500`。

```bash
curl -L -o short-haul-invoice.xlsx \
  -H 'Authorization: Bearer <JWT>' \
  'https://ssc.mengfei.tech/api/admin/excel-exports/dynamic_settlement_statement_suqian?ids=896&type=shortHaulInvoice'
```

同一路径上的其他常用类型：[拆箱发票](/common-apis/export-unpacking)、[清关发票](/common-apis/export-clearance)、[运费发票](/common-apis/export-freight)。

## MCP

对应工具 `ssc_export_excel`，`type` 为 `shortHaulInvoice`。工具请求本页的 `GET`，把 xlsx 写到临时文件并返回路径。接入见 [MCP 接入与使用](/mcp)。
