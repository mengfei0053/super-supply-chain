---
sidebar_position: 4
title: API 与 MCP
slug: api-and-mcp
---

# API 与 MCP

JSON API 在站点源站的 `/api` 下，**不是** `/super-supply-chain` 下面。`/super-supply-chain/` 只提供管理后台的静态页面。

生产源站示例：`https://ssc.mengfei.tech`。本地则是后端监听的地址，例如 `http://localhost:8081`。

怎么把这套 API 接到 Cursor，见独立说明 [MCP 接入与使用](/mcp)。本页列出接口范围，以及 `ssc-mcp` 已经包好的 20 个工具。

## 认证

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/api/register` | 注册。MCP 未包装此接口。 |
| POST | `/api/login` | 登录，返回 JWT 和用户信息。对应工具 `ssc_login`。 |

其余业务接口在 `/api/admin`，请求头为 `Authorization: Bearer <token>`。JWT 为 HS256，有效期 24 小时。

登录后的主要资源包括结算单、Excel 读取规则、字典、动态 Excel 表、导出规则，以及菜单和选项接口。列表接口沿用 React Admin 的 `range`、`sort`、`filter` 和 `Content-Range`。`range` 为 `[start, end]`，`end` 是开区间端点，不传给 MCP 工具时默认 `[0, 49]`。

完整路径表在仓库的 [`ssc-mcp/API_CATALOG.md`](https://github.com/mengfei0053/super-supply-chain/blob/master/ssc-mcp/API_CATALOG.md)，路由注册在 `backend/main.go`。

## MCP 是什么

`ssc-mcp/` 把上述 HTTP API 包成 MCP 工具，供 Cursor 等客户端调用。默认走 stdio；设置 `SSC_MCP_TRANSPORT=http` 后，Streamable HTTP 在 `/mcp`，旧版 SSE 在 `/sse`。

配置时注意：

- `SSC_BASE_URL` 填源站，例如 `https://ssc.mengfei.tech` 或 `http://localhost:8081`。不要带 `/super-supply-chain`。
- 认证用 `SSC_TOKEN`，或用 `SSC_USERNAME` 加 `SSC_PASSWORD` 让服务自己登录。
- 把 `.env.example` 复制为本地 `.env`。不要把 token 或密码提交进仓库。

安装、两种传输、Cursor 的 `mcp.json` 示例和冒烟命令都写在 [MCP 接入与使用](/mcp)。

## 工具一览 {#tools}

共 20 个工具。说明根据 `ssc-mcp/README.md` 与 `ssc-mcp/src/server.ts` 里的注册信息整理。

### 登录、状态与菜单

| 工具 | 说明 |
| --- | --- |
| `ssc_login` | 调用 `POST /api/login`，把 JWT 留在当前进程里。已设置 `SSC_TOKEN` 时不必再登。 |
| `ssc_status` | 请求 `GET /api/admin/menus`，检查当前凭证是否有效。与 CLI 的 `ssc status` 相同。 |
| `ssc_list_menus` | 列出动态 Excel 菜单。每项含 `id`、`menuName`、`dynamicTableName`，表名交给后面的 Excel 工具。 |

### 结算单

| 工具 | 说明 |
| --- | --- |
| `ssc_list_orders` | 列出结算单摘要：`id`、`orderNumber`、`arrivalDate`、`arrivalPort`。 |
| `ssc_get_order` | 按 id 读取一条结算单摘要。 |

### 字典

| 工具 | 说明 |
| --- | --- |
| `ssc_list_dicts` | 分页列出字典项。 |
| `ssc_get_dict` | 按 id 读取一条字典。 |
| `ssc_get_dict_map` | 按 `type` 读取 key 到 value 的映射。 |
| `ssc_create_dict` | 新建字典，字段为 `key`、`value`、`type`。 |
| `ssc_update_dict` | 按 id 更新字典，字段可以只提交要改的部分。 |
| `ssc_delete_dict` | 按 id 删除字典。 |

### Excel 读取规则

| 工具 | 说明 |
| --- | --- |
| `ssc_list_excel_read_rules` | 分页列出读取规则，也就是动态表的注册信息。 |
| `ssc_get_excel_read_rule` | 按 id 读取一条规则，包含字段映射 `rules`。 |

读取规则的创建、更新、删除 HTTP 接口存在，MCP 目前只包装了查询。

### 动态 Excel 行

| 工具 | 说明 |
| --- | --- |
| `ssc_list_excel_rows` | 按表名分页列出数据行。可用 `filterStart`、`filterEnd` 过滤 `created_at`（`YYYY-MM-DD`）。未传日期时，工具默认使用 `2000-01-01` 到 `2099-12-31`。 |
| `ssc_get_excel_row` | 按表名和 id 读取一行。 |
| `ssc_update_excel_row` | 用 JSON 更新一行，字段与 `DynamicExcelTable` 一致，例如 `fileName`、`datas`。 |
| `ssc_delete_excel_row` | 按表名和 id 硬删除一行。 |

表名来自菜单或读取规则，例如 `dynamic_customs_declaration_form`。目录中出现过的动态表还有 `dynamic_Integrity_packaging_invoice`、`dynamic_settlement_statement_fenchang`、`dynamic_settlement_statement_suqian`、`dynamic_yifan_cost_cal`。以 `ssc_list_menus` 的结果为准。

### 导出规则

| 工具 | 说明 |
| --- | --- |
| `ssc_list_export_rules` | 按表名分页列出导出模板。 |
| `ssc_list_export_template_options` | 按关联表名列出导出模板的下拉选项（查询参数 `associated_table`）。 |
| `ssc_get_export_rule` | 按表名和 id 读取一条导出模板。 |

## 尚未包装的能力 {#not-wrapped}

下列接口在后端存在，MCP **还没有**对应工具：

| 能力 | HTTP | 原因 |
| --- | --- | --- |
| 动态表 Excel 上传 | `POST /api/admin/excel/:tableName` | multipart 文件 |
| 结算单文件上传 | `POST /api/admin/settlement-form-entries` | multipart 文件 |
| 批量导出下载 | `GET /api/admin/excel-exports/:tableName` | 返回二进制文件 |
| 注册账号 | `POST /api/register` | 未包装 |

管理后台的 Excel 上传已经可用：服务端按 `UPLOAD_SERVER` 写入本机目录（`file://` 或 `local://`，生产示例为 `file:///data/ssc-uploads`），不再依赖从部署机访问不到的家里 NAS。这只打通了网站上传。MCP 上传工具仍然待做，现有工具不能代替在网页里传文件。

Cursor 配置示例（stdio 命令与 HTTP `url`）见 [MCP 接入与使用](/mcp)。
