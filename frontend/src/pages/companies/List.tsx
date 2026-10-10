import {
  CreateButton,
  Datagrid,
  DeleteButton,
  EditButton,
  List,
  TextField,
  TextInput,
  TopToolbar,
} from "react-admin";
import * as React from "react";

const filters = [
  <TextInput key="q" source="q" label="搜索名称/别名" alwaysOn resettable />,
];

const ListActions = () => (
  <TopToolbar>
    <CreateButton label="新建" />
  </TopToolbar>
);

const ListPage: React.FunctionComponent = () => {
  return (
    <List
      title="公司基本信息"
      filters={filters}
      actions={<ListActions />}
      exporter={false}
      sort={{ field: "id", order: "ASC" }}
    >
      <Datagrid bulkActionButtons={false} rowClick="edit">
        <TextField source="id" label="ID" />
        <TextField source="name" label="公司名称" />
        <TextField source="alias" label="别名" />
        <TextField source="target_addr" label="发票目标地址" />
        <TextField
          source="unified_social_credit_code"
          label="统一社会信用代码"
        />
        <TextField source="addr_city" label="城市" />
        <TextField source="phone_num" label="电话" />
        <EditButton label="编辑" />
        <DeleteButton
          label="删除"
          mutationMode="pessimistic"
          confirmTitle="删除公司"
          confirmContent="删除后发票导出将无法再匹配该公司（软删除）。确定删除？"
        />
      </Datagrid>
    </List>
  );
};

export default ListPage;
