---
sidebar_position: 7
title: 删除一行
slug: /common-apis/delete-row
---

# 删除一行

按表名和 id 删除动态 Excel 表中的一行。这是硬删除（`Unscoped`），没有回收站。

| 项 | 值 |
| --- | --- |
| 方法 | `DELETE` |
| 路径 | `/api/admin/excel/{tableName}/{id}` |
| 认证 | `Authorization: Bearer <JWT>` |

示例：

```http
DELETE https://ssc.mengfei.tech/api/admin/excel/dynamic_settlement_statement_suqian/896
Authorization: Bearer <JWT>
```

成功时 HTTP `200`，响应 JSON：

```json
{"message": "Delete successfully"}
```

```bash
curl -X DELETE \
  -H 'Authorization: Bearer <JWT>' \
  'https://ssc.mengfei.tech/api/admin/excel/dynamic_settlement_statement_suqian/896'
```

## MCP

这条路径已经有 MCP 工具 `ssc_delete_excel_row`。参数是 `tableName` 和 `id`，例如表名 `dynamic_settlement_statement_suqian`、id `896`。工具说明见 [API 与 MCP](/api-and-mcp#tools)。如何把 `ssc-mcp` 接到 Cursor，见 [MCP 接入与使用](/mcp)。
