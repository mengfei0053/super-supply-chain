---
sidebar_position: 4
title: 个人访问令牌
slug: personal-access-tokens
---

# 个人访问令牌

个人访问令牌（PAT）用来让脚本、HTTP 客户端和 MCP 长期调用 `/api/admin`，不用把账号密码交给它们，也不用每 24 小时换一次登录 JWT。

登录管理后台后，在侧栏 **工具 → 个人访问令牌** 里创建、查看和撤销。界面文案是中文。

## 创建和撤销

1. 用账号密码登录管理后台。管理令牌的接口只认这次登录会话，不认 PAT 本身，所以一把令牌不能再签发新令牌。
2. 打开 **个人访问令牌**，点击 **创建**。
3. 填写名称，方便以后辨认，例如 `cursor-mcp`。
4. 权限范围可以留空，表示与该账号一样可以调用全部管理接口。需要限制时再填写，多个值用英文逗号分隔。
5. 过期日期可以留空。填写后，到期日当天 23:59:59（服务器本地时间）前有效。
6. 创建成功后，明文令牌只出现这一次。复制并放到安全的地方。关闭或刷新页面后，系统不会再显示它。
7. 列表里只能看到名称、前缀、权限范围、创建时间、最近使用和过期时间。不需要时点击 **撤销**，令牌会立即失效。

数据库表是 `personal_access_tokens`。保存的是令牌的 SHA-256 十六进制摘要，不保存明文。前缀（`ssc_pat_` 加少量字符）只用于在列表里辨认。

进程启动时会尝试自动建这张表。数据库账号如果没有建表权限，先在库里执行 [`ssc-sqls/personal_access_tokens.sql`](https://github.com/mengfei0053/super-supply-chain/blob/master/ssc-sqls/personal_access_tokens.sql)。

## 调用 HTTP API

受保护接口都在 `/api/admin`。原来的登录 JWT 仍然可用，请求头同样是 `Authorization: Bearer`。

优先使用请求头，不要把令牌放进 URL：

```bash
curl -H "Authorization: Bearer ssc_pat_替换为你的令牌" \
  "https://ssc.mengfei.tech/api/admin/menus"
```

也可以把同一串令牌放在 `X-API-Key` 里：

```bash
curl -H "X-API-Key: ssc_pat_替换为你的令牌" \
  "https://ssc.mengfei.tech/api/admin/excel/dynamic_settlement_statement_suqian"
```

查询参数 `token` 可用，但 URL 会进入浏览器历史、反向代理日志和 `Referer`。只在客户端实在不能设置请求头时使用：

```bash
curl "https://ssc.mengfei.tech/api/admin/menus?token=ssc_pat_替换为你的令牌"
```

本服务自己的访问日志会去掉查询参数里的 `token`。上游代理不一定会这样做，所以仍然优先用请求头。

上传 Excel、删除一行、导出发票，以及其它已经登录才能访问的接口，把原来的 JWT 换成 PAT 即可。路径和参数不变，见 [常用接口](/common-apis) 和 [API 与 MCP](/api-and-mcp)。

## 权限范围

| 值 | 可以访问 |
| --- | --- |
| 留空 | 全部 `/api/admin` 接口 |
| `*` 或 `all` | 全部 `/api/admin` 接口 |
| `excel` | Excel 上传、读取、更新、删除、导出，以及读取规则、导出规则 |
| `dict` | 字典管理 |
| `settlement` | 结算单 |
| `meta` | 菜单和选项 |

范围不匹配时接口返回 403。留空即可覆盖上传、删除和导出。

## 给 MCP 用

`ssc-mcp` 不关心令牌是 JWT 还是 PAT。公网服务把 PAT 放在服务端的 `SSC_TOKEN` 里，访问源站时使用 `Authorization: Bearer`。远程连接器若自行带请求头，同样是 `Authorization: Bearer <token>`。这样不用在 MCP 配置里写密码，也不会在 24 小时后失效（除非你给令牌设了过期时间）。

不要把真实令牌提交进仓库。连接地址见 [MCP 接入与使用](/mcp)。
