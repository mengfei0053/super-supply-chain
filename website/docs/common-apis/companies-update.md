---
sidebar_position: 11
title: 更新公司
slug: /common-apis/companies-update
---

# 更新公司

修改公司基础表 `base_companies_infos` 里一行的基本信息（名称、别名、发票目标地址，以及地址、银行代码、电话）。统一社会信用代码不能改。已软删除的行不能改。

| 项 | 值 |
| --- | --- |
| 方法 | `PUT` |
| 路径 | `/api/admin/companies/{id}` |
| 认证 | `Authorization: Bearer <token>` |

`<token>` 与本分组其他接口相同：登录 JWT，或个人访问令牌。`{id}` 是 [公司查询](/common-apis/companies) 返回的 `id`。

```http
PUT https://ssc.mengfei.tech/api/admin/companies/12
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "杭州测试公司",
  "alias": "杭州测试",
  "target_addr": "海宁"
}
```

```bash
curl -X PUT \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{"name":"杭州测试公司","alias":"杭州测试","target_addr":"海宁"}' \
  'https://ssc.mengfei.tech/api/admin/companies/12'
```

只改目标地址时，请求体里可以只放 `target_addr`：

```bash
curl -X PUT \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{"target_addr":"海宁"}' \
  'https://ssc.mengfei.tech/api/admin/companies/12'
```

## 请求体

JSON 对象。至少要有 `name`、`alias`、`target_addr` 之一。未出现的字段保持原值。`null` 与不传相同。首尾空白会去掉。

| 字段 | 说明 |
| --- | --- |
| `name` | 可选。新的公司名称，不能为空，最长 255 个字符，不能与其它行重复。 |
| `alias` | 可选。新的公司别名，最长 255 个字符。空字符串会清空别名。 |
| `target_addr` | 可选。新的发票目标地址，最长 255 个字符。空字符串会清空。 |

`unified_social_credit_code`、地址、银行代码、电话即使出现在 JSON 里也不会被修改。

成功时 HTTP `200`，响应体是更新后的整行，列与 [公司查询](/common-apis/companies) 相同。

```json
{
  "id": 12,
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-10-09T00:00:00Z",
  "deleted_at": null,
  "name": "杭州测试公司",
  "addr_country": "中国",
  "addr_province": "浙江",
  "addr_city": "杭州",
  "addr_street": "示例路 1 号",
  "unified_social_credit_code": "91330100TEST",
  "bank_code": "BANK001",
  "phone_num": "0571-0000000",
  "alias": "杭州测试",
  "target_addr": "海宁"
}
```

`id` 不是正整数、三个字段都没传、或 `name` 为空时返回 `400`。行不存在或已软删除时返回 `404`。新名称与其它行重复时返回 `409`。未登录返回 `401`。

## MCP

对应工具 `ssc_update_company`。参数是 `id`，以及可选的 `name`、`alias`、`targetAddr`。工具请求本页的 `PUT`，只提交调用时给出的字段，并把更新后的行放在结果的 `company` 里。工具说明见 [API 与 MCP](/api-and-mcp#tools)。接入见 [MCP 接入与使用](/mcp)。

软删除见 [删除公司](/common-apis/companies-delete)。
