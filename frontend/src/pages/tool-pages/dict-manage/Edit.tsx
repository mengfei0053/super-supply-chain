import { Edit, SelectInput, SimpleForm, TextInput } from "react-admin";
import * as React from "react";
import { useNavigate } from "react-router";

const ToolPageEdit: React.FunctionComponent = () => {
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
        <TextInput source="key" label="港口原文" helperText="结算单 arrival_port，如 上海洋山 / 上海洋山保税" />
        <TextInput source="value" label="起始地城市" helperText="发票起始地，如 上海" />
        <TextInput source="port_name" label="口岸名称" helperText="如 上海口岸；可空则用「城市+口岸」" />
        <TextInput source="extra_pay" label="补差项" helperText="如 洋山补差；无可留空" />
        <SelectInput source="type" choices={["港口字典"]} />
      </SimpleForm>
    </Edit>
  );
};

export default ToolPageEdit;
