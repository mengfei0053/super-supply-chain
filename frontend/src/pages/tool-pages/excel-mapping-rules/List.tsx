import {
  List,
  Datagrid,
  TextField,
  TopToolbar,
  CreateButton,
  ArrayField,
  DeleteButton,
  EditButton,
} from "react-admin";
import * as React from "react";

export const ListActions = () => {
  return <TopToolbar>{<CreateButton></CreateButton>}</TopToolbar>;
};

const ToolPageList: React.FunctionComponent = () => {
  return (
    <List actions={<ListActions></ListActions>}>
      <Datagrid>
        <TextField source="id"></TextField>
        <TextField source="name"></TextField>
        <TextField source="type"></TextField>
        <TextField source="startRow"></TextField>
        <ArrayField source="rules">
          {
            <Datagrid bulkActionButtons={false}>
              <TextField source="desc"></TextField>
              <TextField source="jsonKey"></TextField>
              <TextField source="excelKey"></TextField>
            </Datagrid>
          }
        </ArrayField>
        <>
          <DeleteButton></DeleteButton>
          <EditButton></EditButton>
        </>
      </Datagrid>
    </List>
  );
};

export default ToolPageList;
