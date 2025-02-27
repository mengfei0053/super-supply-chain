import {
  Create,
  SimpleForm,
  TextInput,
  SimpleFormIterator,
  ArrayInput,
  SelectInput,
  NumberInput,
} from "react-admin";
import * as React from "react";
import { useNavigate } from "react-router";

export const choices = ["MAP_RULE", "ITERATE_RULE"];

const ToolPageCreate: React.FunctionComponent = () => {
  const navigate = useNavigate();

  return (
    <Create
      mutationOptions={{
        onSuccess() {
          navigate("/excel-mapping-rule");
        },
      }}
    >
      <SimpleForm>
        <TextInput source="name"></TextInput>
        <SelectInput source="type" required choices={choices}></SelectInput>
        <NumberInput source="startRow" min={0} step={1}></NumberInput>
        <ArrayInput source="rules" fullWidth={false}>
          <SimpleFormIterator inline>
            <TextInput source="desc" helperText={false}></TextInput>
            <TextInput source="excelKey" helperText={false}></TextInput>
            <TextInput source="jsonKey" helperText={false}></TextInput>
          </SimpleFormIterator>
        </ArrayInput>
      </SimpleForm>
    </Create>
  );
};

export default ToolPageCreate;
