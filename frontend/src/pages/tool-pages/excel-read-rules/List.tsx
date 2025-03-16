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
import DataField from "../../excels/DataField";

export const ListActions = () => {
  return <TopToolbar>{<CreateButton></CreateButton>}</TopToolbar>;
};

const ListPage: React.FunctionComponent = () => {
  return (
    <List actions={<ListActions></ListActions>}>
      <Datagrid rowClick={false}>
        <TextField source="id"></TextField>
        <TextField source="menuName"></TextField>
        <TextField source="dynamicTableName"></TextField>
        <TextField source="sheetIndex"></TextField>

        <DataField source="rules"></DataField>

        <>
          <EditButton></EditButton>
        </>
      </Datagrid>
    </List>
  );
};

export default ListPage;
