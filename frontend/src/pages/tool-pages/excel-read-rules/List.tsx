import {
  List,
  Datagrid,
  TextField,
  TopToolbar,
  CreateButton,
  ArrayField,
  DeleteButton,
  EditButton,
  FunctionField,
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
        <TextField source="menuName"></TextField>
        <TextField source="desc"></TextField>
        <TextField source="sheetIndex"></TextField>
        <FunctionField
          label="MapRule"
          render={(r) => r.mapRule?.name}
        ></FunctionField>
        <FunctionField
          label="iterateRule"
          render={(r) => r.iterateRule?.name}
        ></FunctionField>

        <>
          <DeleteButton></DeleteButton>
          <EditButton></EditButton>
        </>
      </Datagrid>
    </List>
  );
};

export default ListPage;
