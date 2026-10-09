---
sidebar_position: 2
title: 上传 Excel
slug: /common-apis/upload
---

# 上传 Excel

把一份 xlsx 解析后写入动态表。管理后台的上传按钮走的就是这条接口。

| 项 | 值 |
| --- | --- |
| 方法 | `POST` |
| 路径 | `/api/admin/excel/{tableName}` |
| 认证 | `Authorization: Bearer <JWT>` |
| 请求体 | `multipart/form-data` |

示例（宿迁结算表）：

```http
POST https://ssc.mengfei.tech/api/admin/excel/dynamic_settlement_statement_suqian
Authorization: Bearer <JWT>
Content-Type: multipart/form-data
```

## 表单字段

| 字段 | 内容 |
| --- | --- |
| `file` | xlsx 文件。服务端用 `FormFile("file")` 读取，缺这个字段会失败。 |
| `name` | 文件名。管理后台上传时与 `file` 一起提交，值为所选文件的 `name`。 |

入库时的 `FileName` 取自 `file` 这一部分自带的文件名。

```bash
curl -X POST \
  'https://ssc.mengfei.tech/api/admin/excel/dynamic_settlement_statement_suqian' \
  -H 'Authorization: Bearer <JWT>' \
  -F 'file=@./statement.xlsx' \
  -F 'name=statement.xlsx'
```

成功时 HTTP `200`，响应体是解析后的表数据 JSON。

文件落在服务端配置的上传目录。`UPLOAD_SERVER` 为 `file://` 或 `local://`（生产示例 `file:///data/ssc-uploads`）时写本机目录；设成可访问的 WebDAV 地址时才走 NAS。部署机访问不到家里局域网的 NAS。

## MCP

还没有对应的 MCP 工具。`ssc-mcp` 不能代替这次上传。以后若包装，应调用本页的 `POST`，表单仍是 `file` 与 `name`。接入见 [MCP 接入与使用](/mcp)。
