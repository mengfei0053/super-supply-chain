# Super Supply Chain CLI

`ssc` 是调用 `backend/` 服务接口的命令行工具。当前第一阶段规划并实现两个功能：登录和状态检查。

## 功能规划

### `login`

目标：调用后端公开接口 `POST /api/login`，保存返回的 JWT，供后续 CLI 命令使用。

- 请求地址：`{baseUrl}/api/login`
- 请求体：`{"username":"...","password":"..."}`
- 成功响应：后端返回用户信息和 `token`
- 本地状态：保存到 OS 用户配置目录下的 `super-supply-chain/cli.json`
- 安全约束：配置文件权限使用 `0600`，不在输出中打印密码或 token

### `status`

目标：展示当前 CLI 登录状态，并验证 token 是否仍可访问后端受保护接口。

- 读取本地配置文件
- 解析 JWT 的 `exp` 字段，提示是否过期
- 默认调用 `GET /api/admin/menus` 做远程认证检查
- 如果后端不可达、token 失效或未登录，返回非 0 退出码

> 后端目前没有专用 health/status 接口，因此第一版 `status` 使用现有受保护接口 `/api/admin/menus` 验证认证状态。

## 使用

```sh
cd cli
go run . login --base-url http://localhost:8081 --username <账号> --password <密码>
go run . status
```

也可以用环境变量减少参数：

```sh
export SSC_BASE_URL=http://localhost:8081
export SSC_USERNAME=<账号>
export SSC_PASSWORD=<密码>
go run . login
```

构建二进制：

```sh
cd cli
go build -o ssc .
./ssc status
```

## 配置

| 变量 | 说明 | 默认值 |
| --- | --- | --- |
| `SSC_BASE_URL` | 后端服务基础地址 | `http://localhost:8081` |
| `SSC_USERNAME` | 登录用户名 | 空 |
| `SSC_PASSWORD` | 登录密码 | 空 |
| `SSC_CONFIG` | CLI 配置文件路径 | OS 用户配置目录下 `super-supply-chain/cli.json` |

## 后续扩展建议

- 添加 `logout`：删除本地 token。
- 添加专用后端 `GET /api/admin/status` 或 `GET /api/health`，避免 `status` 依赖菜单接口。
- 将动态 Excel、字典、结算单等接口按资源拆分为子命令。
