---
sidebar_position: 1
title: 一览
slug: /common-apis
---

# 常用接口

宿迁结算表上最常调用的是：上传 Excel、按类型下载四类发票，以及按 id 删除一行。下面每条都是独立说明，侧边栏「常用接口」可以直接点开。

源站是 `https://ssc.mengfei.tech`。路径在 `/api/admin` 下，不要写成 `/super-supply-chain/...`。都要登录后的 JWT：

```http
Authorization: Bearer <JWT>
```

没有这个头、token 无效或过期时，接口返回 `401`。

示例表名：`dynamic_settlement_statement_suqian`。示例行 id：`896`。其他动态表把路径里的表名换成 `ssc_list_menus` 返回的 `dynamicTableName` 即可，查询参数不变。

| 说明 | 方法与路径 | 结果 |
| --- | --- | --- |
| [上传 Excel](/common-apis/upload) | `POST /api/admin/excel/{tableName}` | JSON。表单字段 `file`（xlsx）、`name`（文件名） |
| [短驳发票](/common-apis/export-short-haul) | `GET /api/admin/excel-exports/{tableName}?ids={id}&type=shortHaulInvoice` | xlsx 文件下载 |
| [拆箱发票](/common-apis/export-unpacking) | `GET /api/admin/excel-exports/{tableName}?ids={id}&type=invoice_unpacking` | xlsx 文件下载 |
| [清关发票](/common-apis/export-clearance) | `GET /api/admin/excel-exports/{tableName}?ids={id}&type=invoice_clearance_only` | xlsx 文件下载 |
| [运费发票](/common-apis/export-freight) | `GET /api/admin/excel-exports/{tableName}?ids={id}&type=invoice_freight` | xlsx 文件下载 |
| [删除一行](/common-apis/delete-row) | `DELETE /api/admin/excel/{tableName}/{id}` | JSON |

四类导出的响应体是 Excel 文件，不是 JSON。浏览器或 HTTP 客户端应按附件保存。

## 和 MCP 的关系

这几条都已经有 MCP 工具：上传是 `ssc_upload_excel`，删除是 `ssc_delete_excel_row`，四类导出发票是 `ssc_export_excel`。工具使用本分组里的路径和参数。接入步骤见 [MCP 接入与使用](/mcp)，工具说明见 [API 与 MCP](/api-and-mcp#tools)。
