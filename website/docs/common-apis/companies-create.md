---
sidebar_position: 10
title: 新增公司
slug: /common-apis/companies-create
---

# 新增公司

往公司基础表 `base_companies_infos` 新增一行。娃哈哈发票导出按 `name` 精确匹配，或按 `alias` 包含匹配查找公司；`target_addr` 是发票目标地址。

| 项 | 值 |
| --- | --- |
| 方法 | `POST` |
| 路径 | `/api/admin/companies` |
| 认证 | `Authorization: Bearer <token>` |

`<token>` 与本分组其他接口相同：登录 JWT，或个人访问令牌。

```http
POST https://ssc.mengfei.tech/api/admin/companies
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "杭州测试公司",
  "alias": "杭州测试",
  "target_addr": "杭州",
  "unified_social_credit_code": "91330100TEST"
}
```

```bash
curl -X POST \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{"name":"杭州测试公司","alias":"杭州测试","target_addr":"杭州","unified_social_credit_code":"91330100TEST"}' \
  'https://ssc.mengfei.tech/api/admin/companies'
```

## 请求体

JSON 对象。首尾空白会去掉。

| 字段 | 说明 |
| --- | --- |
| `name` | 必填。公司名称，最长 255 个字符，全表唯一（含已软删除的行）。 |
| `unified_social_credit_code` | 必填。统一社会信用代码，最长 100 个字符，全表唯一。现有表结构要求这一列，不能省略。 |
| `alias` | 可选。公司别名，最长 255 个字符。不传或空字符串都存成空。 |
| `target_addr` | 可选。发票目标地址，最长 255 个字符。不传或空字符串都存成空。 |

地址、银行代码、电话不在这个接口里填写，新行这些列为空。不要传表上没有的字段名；多出来的字段会被忽略。

成功时 HTTP `200`，响应体是刚写入的那一行，列与 [公司查询](/common-apis/companies) 相同。`deleted_at` 为 `null`。

```json
{
  "id": 12,
  "created_at": "2026-10-09T00:00:00Z",
  "updated_at": "2026-10-09T00:00:00Z",
  "deleted_at": null,
  "name": "杭州测试公司",
  "addr_country": "",
  "addr_province": "",
  "addr_city": "",
  "addr_street": "",
  "unified_social_credit_code": "91330100TEST",
  "bank_code": "",
  "phone_num": "",
  "alias": "杭州测试",
  "target_addr": "杭州"
}
```

缺 `name` 或 `unified_social_credit_code`、字段过长时返回 `400`。名称或统一社会信用代码已存在时返回 `409`。未登录返回 `401`。

## MCP

对应工具 `ssc_create_company`。参数是 `name`、`unifiedSocialCreditCode`，可选 `alias`、`targetAddr`。工具请求本页的 `POST`，并把写入的行放在结果的 `company` 里。工具说明见 [API 与 MCP](/api-and-mcp#tools)。接入见 [MCP 接入与使用](/mcp)。
