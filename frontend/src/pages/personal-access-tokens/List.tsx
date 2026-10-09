import {
  CreateButton,
  Datagrid,
  DateField,
  DeleteButton,
  FunctionField,
  List,
  TextField,
  TopToolbar,
} from "react-admin";
import { Box, Typography } from "@mui/material";
import * as React from "react";

const ListActions = () => (
  <TopToolbar>
    <CreateButton label="创建" />
  </TopToolbar>
);

const EmptyTokens = () => (
  <Box sx={{ p: 2 }}>
    <Typography sx={{ mb: 2 }}>
      还没有个人访问令牌。创建后明文只显示一次，请立刻复制保存。
    </Typography>
    <CreateButton label="创建" />
  </Box>
);

const scopeLabel = (scopes?: string) => (scopes ? scopes : "全部接口");

const ListPage: React.FunctionComponent = () => {
  return (
    <List
      title="个人访问令牌"
      actions={<ListActions />}
      empty={<EmptyTokens />}
      exporter={false}
    >
      <Datagrid bulkActionButtons={false} rowClick={false}>
        <TextField source="name" label="名称" sortable={false} />
        <TextField source="prefix" label="前缀" sortable={false} />
        <FunctionField
          label="权限范围"
          sortable={false}
          render={(record: { scopes?: string }) => scopeLabel(record?.scopes)}
        />
        <DateField
          source="createdAt"
          label="创建时间"
          showTime
          locales="zh-CN"
          sortable={false}
        />
        <DateField
          source="lastUsedAt"
          label="最近使用"
          showTime
          locales="zh-CN"
          sortable={false}
        />
        <DateField
          source="expiresAt"
          label="过期时间"
          showTime
          locales="zh-CN"
          sortable={false}
        />
        <DeleteButton
          label="撤销"
          mutationMode="pessimistic"
          confirmTitle="撤销令牌"
          confirmContent="撤销后该令牌会立即失效，且不能恢复。确定撤销？"
        />
      </Datagrid>
    </List>
  );
};

export default ListPage;
