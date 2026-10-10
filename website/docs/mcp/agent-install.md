---
title: Agent 安装说明
sidebar_position: 2
slug: /mcp/agent-install
description: 给 AI Agent 的娃哈哈发票模板导出 MCP 安装清单。
---

# Agent 安装说明

按顺序执行。完成一步再做下一步。连接之前必须先向用户要到个人访问令牌。

## 1. 确认范围

这是 **娃哈哈发票模板导出流程**（Wahaha invoice template export）的一部分。

- Super Supply Chain（SSC）负责娃哈哈结算 Excel 上传，以及四类发票模板导出：`shortHaulInvoice`、`invoice_unpacking`、`invoice_clearance_only`、`invoice_freight`。
- 远程 MCP 的 `ssc_upload_excel`、`ssc_list_excel`、`ssc_delete_excel_row`、`ssc_export_excel`、`ssc_search_companies`、`ssc_create_company`、`ssc_update_company`、`ssc_delete_company` 是这条流水线面向 Agent 的入口。
- 公网地址用 `https://ssc.mengfei.tech/mcp`。

## 2. 取得并安装 Skill

1. 下载 Skill zip：`https://ssc.mengfei.tech/docs/files/ssc-mcp-skill.zip`。文档页 [MCP 接入与使用](/mcp) 上的「下载 SSC MCP Skill」指向同一文件。
2. 解压。确认路径是 `ssc-mcp-skill/SKILL.md`。
3. 把整个 `ssc-mcp-skill` 文件夹装进 Cursor 的 skills 目录，使 `SKILL.md` 位于文件夹根下：
   - 用户级：`~/.cursor/skills/ssc-mcp-skill/SKILL.md`
   - 仅当前项目：`<project>/.cursor/skills/ssc-mcp-skill/SKILL.md`
4. 安装的是解压后的文件夹，不是 zip 本身。

## 3. 向用户索取 PAT

连接和写配置之前先做这一步。

1. 告诉用户：登录管理后台 `https://ssc.mengfei.tech/super-supply-chain/`，打开 **工具 → 个人访问令牌**，创建一把令牌。明文只显示一次，前缀是 `ssc_pat_`。
2. 停下来，等用户把令牌发给你。
3. 不要编造、猜测或写死令牌。不要把占位符 `ssc_pat_…` 当作真实令牌去连接。不要把令牌提交进仓库或写进 `SKILL.md`。

## 4. 配置远程 MCP

只在用户已经给出 PAT 之后写入 Cursor 的 MCP 配置（`~/.cursor/mcp.json` 或项目 `.cursor/mcp.json`）。把 `<PAT>` 换成用户提供的那一串。

```json
{
  "mcpServers": {
    "ssc": {
      "url": "https://ssc.mengfei.tech/mcp",
      "headers": {
        "Authorization": "Bearer <PAT>"
      }
    }
  }
}
```

`Authorization` 的值包含 `Bearer ` 前缀和一个空格。

同一串令牌也可以放在 `X-API-Key` 里，这时不要加 `Bearer `：

```json
{
  "mcpServers": {
    "ssc": {
      "url": "https://ssc.mengfei.tech/mcp",
      "headers": {
        "X-API-Key": "<PAT>"
      }
    }
  }
}
```

客户端只能使用旧版 SSE 时，保留同样的 `headers`，把 `url` 换成 `https://ssc.mengfei.tech/sse`。

## 5. 核对工具后停止

1. 用上一步的配置连接 MCP。
2. 列出工具。确认有 `ssc_upload_excel`、`ssc_list_excel`、`ssc_delete_excel_row`、`ssc_export_excel`、`ssc_search_companies`、`ssc_create_company`、`ssc_update_company`、`ssc_delete_company`。
3. 停。在用户下一步明确要求之前，不要上传、列表查询、删除、导出、搜索公司、新增公司、更新公司或删除公司。
