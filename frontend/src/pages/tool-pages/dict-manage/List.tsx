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
        <TextField source="id"></TextField>
        <TextField source="key"></TextField>
        <TextField source="value"></TextField>
        <>
          <DeleteButton></DeleteButton>
          <EditButton></EditButton>
        </>
      </Datagrid>
    </List>
  );
};

export default ListPage;
