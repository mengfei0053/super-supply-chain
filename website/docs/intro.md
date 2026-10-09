---
sidebar_position: 1
title: 概览
slug: intro
---

# Super Supply Chain 概览

Super Supply Chain（SSC）是一套供应链管理后台。它把结算录入、Excel 动态表、字典和读取规则放在同一个 Web 系统里，并用 Go API 对外提供数据。

生产环境里，管理后台和文档站在同一域名下，路径分开：

| 路径 | 内容 |
| --- | --- |
| `/super-supply-chain/` | React Admin 管理后台 |
| `/docs/` | 本说明站（例如 [`/docs/intro`](/intro)） |
| `/api/` | HTTP API。公开接口在 `/api`，登录后的接口在 `/api/admin` |

## 仓库里有什么

- `frontend/`：React 19 + Vite + React Admin。页面包括结算单、动态 Excel 表、字典、读取规则和翊帆费用计算。
- `backend/`：Go 1.23 + Gin + GORM。负责登录、业务接口，以及生产环境下的静态文件。
- `website/`：本 Docusaurus 站点。
- `ssc-mcp/`：把后端 HTTP API 包成 MCP 工具，供 Cursor 等客户端调用。
- `cli/`：命令行登录和状态检查。
- `ssc-sqls/`：数据库初始化脚本。

更细的开发约定写在仓库根目录的 `AGENTS.md`，以及 `backend/`、`frontend/` 下各自的 `AGENTS.md`。

## 后台能做什么

登录后的侧栏大致分成两类：

- **Excel 处理**：按动态表名进入对应数据页，做录入、查看和导出。
- **工具**：字典管理、Excel 读取规则。

顶栏的「文档」会打开本站首页 `/docs/`。
