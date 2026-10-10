---
sidebar_position: 9
title: 公司查询
slug: /common-apis/companies
---

# 公司查询

按关键字查询公司基础表 `base_companies_infos`。`name` 和 `alias` 都会做包含匹配，任一字段命中即返回该行。

| 项 | 值 |
| --- | --- |
| 方法 | `GET` |
| 路径 | `/api/admin/companies` |
| 认证 | `Authorization: Bearer <token>` |

`<token>` 与本分组其他接口相同：登录 JWT，或个人访问令牌。

示例（关键字 `南阳`）：

```http
GET https://ssc.mengfei.tech/api/admin/companies?keyword=南阳
Authorization: Bearer <token>
```

## 查询参数

| 参数 | 说明 |
| --- | --- |
| `keyword` | 必填。对 `name`、`alias` 做 `LIKE %keyword%`。`%` 和 `_` 按字面量匹配，不会变成通配符。空关键字返回 400。 |
| `includeDeleted` | 可选。默认不传，只返回 `deleted_at IS NULL` 的行。设为 `true`（或 `1`、`yes`）时连同已软删除的行一起返回。 |

```bash
curl -G \
  -H 'Authorization: Bearer <token>' \
  --data-urlencode 'keyword=南阳' \
  'https://ssc.mengfei.tech/api/admin/companies'
```

包含已删除行：

```bash
curl -G \
  -H 'Authorization: Bearer <token>' \
  --data-urlencode 'keyword=南阳' \
  --data-urlencode 'includeDeleted=true' \
  'https://ssc.mengfei.tech/api/admin/companies'
```

成功时 HTTP `200`，响应体是 JSON 数组。每一行带出表上的全部列，而不是摘要：`id`、`created_at`、`updated_at`、`deleted_at`、`name`、`addr_country`、`addr_province`、`addr_city`、`addr_street`、`unified_social_credit_code`、`bank_code`、`phone_num`、`alias`、`target_addr`。没有命中时是 `[]`。

```json
[
  {
    "id": 1,
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-01T00:00:00Z",
    "deleted_at": null,
    "name": "南阳食品有限公司",
    "addr_country": "中国",
    "addr_province": "河南",
    "addr_city": "南阳",
    "addr_street": "示例路 1 号",
    "unified_social_credit_code": "91330000NAME",
    "bank_code": "BANK001",
    "phone_num": "0377-0000000",
    "alias": "食品",
    "target_addr": "南阳"
  }
]
```

## MCP

对应工具 `ssc_search_companies`。参数是 `keyword`，可选 `includeDeleted`。工具请求本页的 `GET`，并把完整行数组放在结果的 `rows` 里。工具说明见 [API 与 MCP](/api-and-mcp#tools)。接入见 [MCP 接入与使用](/mcp)。

不带 `keyword`、带 react-admin 的 `range` 时，同一路径返回分页列表，并写 `Content-Range`；`filter.q` / `filter.keyword` / `filter.name` 可按名称或别名筛选（管理后台「公司基本信息」用）。`GET /api/admin/companies/{id}` 取一行。

新增一行见 [新增公司](/common-apis/companies-create)，修改见 [更新公司](/common-apis/companies-update)，软删除见 [删除公司](/common-apis/companies-delete)。
