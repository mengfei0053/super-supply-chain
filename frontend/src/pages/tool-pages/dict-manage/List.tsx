import {
  List,
  Datagrid,
  TextField,
  TopToolbar,
  CreateButton,
  DeleteButton,
  EditButton,
} from "react-admin";
import * as React from "react";

export const ListActions = () => {
  return <TopToolbar>{<CreateButton></CreateButton>}</TopToolbar>;
};

const ListPage: React.FunctionComponent = () => {
  return (
    <List actions={<ListActions></ListActions>}>
      <Datagrid>
        <TextField source="id" />
        <TextField source="type" label="类型" />
        <TextField source="key" label="港口原文" />
        <TextField source="value" label="起始地" />
        <TextField source="port_name" label="口岸" />
        <TextField source="extra_pay" label="补差" />
        <>
          <DeleteButton />
          <EditButton />
        </>
      </Datagrid>
    </List>
  );
};

export default ListPage;
