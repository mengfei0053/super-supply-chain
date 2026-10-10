---
sidebar_position: 12
title: 删除公司
slug: /common-apis/companies-delete
---

# 删除公司

软删除公司基础表 `base_companies_infos` 中的一行（写入 `deleted_at`）。发票导出之后不再匹配该公司。已删除或不存在的 id 返回 404。

| 项 | 值 |
| --- | --- |
| 方法 | `DELETE` |
| 路径 | `/api/admin/companies/{id}` |
| 认证 | `Authorization: Bearer <token>` |

`<token>` 与本分组其他接口相同：登录 JWT，或个人访问令牌。`{id}` 是 [公司查询](/common-apis/companies) 返回的 `id`。

```http
DELETE https://ssc.mengfei.tech/api/admin/companies/12
Authorization: Bearer <token>
```

```bash
curl -X DELETE   -H 'Authorization: Bearer <token>'   'https://ssc.mengfei.tech/api/admin/companies/12'
```

成功时 HTTP `200`：

```json
{
  "message": "Delete company successfully",
  "id": 12
}
```

## MCP

对应工具 `ssc_delete_company`。参数是 `id`。工具说明见 [API 与 MCP](/api-and-mcp#tools)。接入见 [MCP 接入与使用](/mcp)。
