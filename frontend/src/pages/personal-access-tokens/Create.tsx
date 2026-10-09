import { Title, useNotify, useRefresh } from "react-admin";
import {
  Alert,
  Box,
  Button,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import * as React from "react";
import { useNavigate } from "react-router-dom";
import { httpClient } from "../../dataProvider";

type IssuedToken = {
  token: string;
  prefix: string;
  name: string;
};

const errorMessage = (error: unknown) => {
  if (typeof error === "object" && error && "body" in error) {
    const body = (error as { body?: { error?: string; message?: string } })
      .body;
    if (body?.error) {
      return body.error;
    }
    if (body?.message) {
      return body.message;
    }
  }
  if (error instanceof Error && error.message) {
    return error.message;
  }
  return "创建失败";
};

const CreatePage: React.FunctionComponent = () => {
  const notify = useNotify();
  const refresh = useRefresh();
  const navigate = useNavigate();
  const [name, setName] = React.useState("");
  const [scopes, setScopes] = React.useState("");
  const [expiresAt, setExpiresAt] = React.useState("");
  const [submitting, setSubmitting] = React.useState(false);
  const [issued, setIssued] = React.useState<IssuedToken | null>(null);

  const submit = async () => {
    const trimmed = name.trim();
    if (!trimmed) {
      notify("请填写名称", { type: "warning" });
      return;
    }
    setSubmitting(true);
    try {
      const payload: { name: string; scopes?: string; expiresAt?: string } = {
        name: trimmed,
      };
      if (scopes.trim()) {
        payload.scopes = scopes.trim();
      }
      if (expiresAt) {
        payload.expiresAt = expiresAt;
      }
      const { json } = await httpClient(
        `${import.meta.env.VITE_JSON_SERVER_URL}/personal-access-tokens`,
        {
          method: "POST",
          body: JSON.stringify(payload),
        },
      );
      const token = typeof json?.token === "string" ? json.token : "";
      if (!token) {
        notify("创建失败：服务器没有返回令牌", { type: "warning" });
        return;
      }
      setIssued({
        token,
        prefix: typeof json.prefix === "string" ? json.prefix : "",
        name: typeof json.name === "string" ? json.name : trimmed,
      });
    } catch (error) {
      notify(errorMessage(error), { type: "error" });
    } finally {
      setSubmitting(false);
    }
  };

  const copyToken = async () => {
    if (!issued) {
      return;
    }
    try {
      await navigator.clipboard.writeText(issued.token);
      notify("已复制令牌", { type: "success" });
    } catch {
      notify("复制失败，请手动选择令牌复制", { type: "warning" });
    }
  };

  const finish = () => {
    setIssued(null);
    refresh();
    navigate("/personal-access-tokens");
  };

  return (
    <Box sx={{ p: 2, maxWidth: 720 }}>
      <Title title="创建个人访问令牌" />
      {issued ? (
        <Stack spacing={2}>
          <Alert severity="warning">
            令牌只显示这一次。关闭或刷新页面后不能再查看，请立刻复制并妥善保存。
          </Alert>
          <Typography>
            {issued.name}
            {issued.prefix ? `（前缀 ${issued.prefix}）` : ""}
          </Typography>
          <TextField
            label="令牌"
            value={issued.token}
            fullWidth
            multiline
            InputProps={{ readOnly: true }}
            inputProps={{
              autoComplete: "off",
              spellCheck: "false",
              "aria-label": "个人访问令牌",
            }}
          />
          <Stack direction="row" spacing={1}>
            <Button variant="contained" onClick={() => void copyToken()}>
              复制
            </Button>
            <Button variant="outlined" onClick={finish}>
              我已保存，返回列表
            </Button>
          </Stack>
        </Stack>
      ) : (
        <Box
          component="form"
          onSubmit={(event) => {
            event.preventDefault();
            void submit();
          }}
        >
          <Stack spacing={2}>
            <Alert severity="info">
              调用接口时优先使用请求头 Authorization: Bearer。查询参数 token
              可用，但容易进入访问日志。
            </Alert>
            <TextField
              label="名称"
              value={name}
              required
              fullWidth
              autoComplete="off"
              onChange={(event) => setName(event.target.value)}
              helperText="用来辨认这把令牌，例如 cursor-mcp"
            />
            <TextField
              label="权限范围"
              value={scopes}
              fullWidth
              autoComplete="off"
              onChange={(event) => setScopes(event.target.value)}
              helperText="可选。留空表示可调用全部管理接口。多个用英文逗号分隔：excel、dict、settlement、meta"
            />
            <TextField
              label="过期日期"
              type="date"
              value={expiresAt}
              fullWidth
              onChange={(event) => setExpiresAt(event.target.value)}
              helperText="可选。留空表示不过期。到期日当天结束前有效。"
              InputLabelProps={{ shrink: true }}
            />
            <Stack direction="row" spacing={1}>
              <Button type="submit" variant="contained" disabled={submitting}>
                创建
              </Button>
              <Button
                type="button"
                onClick={() => navigate("/personal-access-tokens")}
              >
                返回
              </Button>
            </Stack>
          </Stack>
        </Box>
      )}
    </Box>
  );
};

export default CreatePage;
