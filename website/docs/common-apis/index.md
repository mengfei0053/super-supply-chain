---
sidebar_position: 1
title: 一览
slug: /common-apis
---

# 常用接口

宿迁结算表上最常调用的是：上传 Excel、按日期查询列表、按类型下载四类发票，以及按 id 删除一行。公司基础表可以按关键字查询，也可以新增、修改和软删除公司基本信息。下面每条都是独立说明，侧边栏「常用接口」可以直接点开。

源站是 `https://ssc.mengfei.tech`。路径在 `/api/admin` 下，不要写成 `/super-supply-chain/...`。都要：

```http
Authorization: Bearer <token>
```

`<token>` 是登录 JWT，或个人访问令牌。没有这个头、token 无效或过期时，接口返回 `401`。

示例表名：`dynamic_settlement_statement_suqian`。示例行 id：`896`。其他动态表把路径里的表名换成 `ssc_list_menus` 返回的 `dynamicTableName` 即可，查询参数不变。

| 说明 | 方法与路径 | 结果 |
| --- | --- | --- |
| [上传 Excel](/common-apis/upload) | `POST /api/admin/excel/{tableName}` | JSON。表单字段 `file`（xlsx）、`name`（文件名） |
| [查询列表](/common-apis/list) | `GET /api/admin/excel/{tableName}?filter=...&sort=...&range=0&range=50` | JSON 数组。`Content-Range` 可能是匹配总数 |
| [短驳发票](/common-apis/export-short-haul) | `GET /api/admin/excel-exports/{tableName}?ids={id}&type=shortHaulInvoice` | xlsx 文件下载 |
| [拆箱发票](/common-apis/export-unpacking) | `GET /api/admin/excel-exports/{tableName}?ids={id}&type=invoice_unpacking` | xlsx 文件下载 |
| [清关发票](/common-apis/export-clearance) | `GET /api/admin/excel-exports/{tableName}?ids={id}&type=invoice_clearance_only` | xlsx 文件下载 |
| [运费发票](/common-apis/export-freight) | `GET /api/admin/excel-exports/{tableName}?ids={id}&type=invoice_freight` | xlsx 文件下载 |
| [删除一行](/common-apis/delete-row) | `DELETE /api/admin/excel/{tableName}/{id}` | JSON |
| [公司查询](/common-apis/companies) | `GET /api/admin/companies?keyword={keyword}` | JSON 数组。匹配 `name` 或 `alias` |
| [新增公司](/common-apis/companies-create) | `POST /api/admin/companies` | JSON 对象。必填 `name`、`unified_social_credit_code`，可选别名/地址等 |
| [更新公司](/common-apis/companies-update) | `PUT /api/admin/companies/{id}` | JSON 对象。可改名称、别名、目标地址、地址、银行、电话 |
| [删除公司](/common-apis/companies-delete) | `DELETE /api/admin/companies/{id}` | 软删除。写入 `deleted_at` |

四类导出的响应体是 Excel 文件，不是 JSON。浏览器或 HTTP 客户端应按附件保存。

## 和 MCP 的关系

这几条都已经有 MCP 工具：上传是 `ssc_upload_excel`，查询列表是 `ssc_list_excel`，删除是 `ssc_delete_excel_row`，四类导出发票是 `ssc_export_excel`，公司查询是 `ssc_search_companies`，新增公司是 `ssc_create_company`，更新公司是 `ssc_update_company`，删除公司是 `ssc_delete_company`。工具使用本分组里的路径和参数。接入步骤见 [MCP 接入与使用](/mcp)，工具说明见 [API 与 MCP](/api-and-mcp#tools)。
