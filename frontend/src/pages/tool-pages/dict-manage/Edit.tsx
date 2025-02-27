import { Edit, SelectInput, SimpleForm, TextInput } from "react-admin";
import * as React from "react";
import { useNavigate } from "react-router-dom";

const ToolPagesEdit: React.FunctionComponent = () => {
  const navigate = useNavigate();

  return (
    <Edit
      mutationOptions={{
        onSuccess() {
          navigate("/dict-manage");
        },
      }}
    >
      <SimpleForm>
        <TextInput source="key"></TextInput>
        <TextInput source="value"></TextInput>
        <SelectInput source="type" choices={["港口字典"]}></SelectInput>
      </SimpleForm>
    </Edit>
  );
};

export default ToolPagesEdit;
