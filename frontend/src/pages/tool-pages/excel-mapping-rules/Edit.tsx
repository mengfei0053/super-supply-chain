import {
  Edit,
  SimpleForm,
  TextInput,
  SimpleFormIterator,
  ArrayInput,
  TopToolbar,
  NumberInput,
  SelectInput,
} from "react-admin";
import * as React from "react";
import { Button } from "@mui/material";

const choices = ["MAP_RULE", "ITERATE_RULE"];

const ToolPagesEdit: React.FunctionComponent = () => {
  return (
    <Edit
      onSubmit={() => {
        console.log("onSubmit");
      }}
      actions={
        <TopToolbar>
          <Button color="primary">Save</Button>
        </TopToolbar>
      }
    >
      <SimpleForm>
        <TextInput source="name"></TextInput>
        <SelectInput source="type" required choices={choices}></SelectInput>
        <NumberInput source="startRow" min={0} step={1}></NumberInput>
        <ArrayInput source="rules">
          <SimpleFormIterator inline>
            <TextInput source="desc" helperText={false}></TextInput>
            <TextInput source="excelKey" helperText={false}></TextInput>
            <TextInput source="jsonKey" helperText={false}></TextInput>
          </SimpleFormIterator>
        </ArrayInput>
      </SimpleForm>
    </Edit>
  );
};

export default ToolPagesEdit;
