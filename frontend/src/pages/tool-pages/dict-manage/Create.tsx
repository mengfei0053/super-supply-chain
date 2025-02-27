import { Create, SelectInput, SimpleForm, TextInput } from "react-admin";
import * as React from "react";
import { useNavigate } from "react-router";

const ToolPageCreate: React.FunctionComponent = () => {
  const navigate = useNavigate();

  return (
    <Create
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
    </Create>
  );
};

export default ToolPageCreate;
