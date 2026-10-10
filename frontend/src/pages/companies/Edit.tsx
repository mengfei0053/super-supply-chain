import { Edit, SimpleForm, TextInput, required } from "react-admin";
import * as React from "react";
import { useNavigate } from "react-router-dom";

const CompanyEdit: React.FunctionComponent = () => {
  const navigate = useNavigate();

  return (
    <Edit
      title="编辑公司"
      mutationMode="pessimistic"
      mutationOptions={{
        onSuccess() {
          navigate("/companies");
        },
      }}
    >
      <SimpleForm>
        <TextInput source="id" label="ID" disabled />
        <TextInput source="name" label="公司名称" validate={required()} fullWidth />
        <TextInput
          source="unified_social_credit_code"
          label="统一社会信用代码"
          disabled
          fullWidth
        />
        <TextInput source="alias" label="别名" fullWidth />
        <TextInput source="target_addr" label="发票目标地址" fullWidth />
        <TextInput source="addr_country" label="国家" fullWidth />
        <TextInput source="addr_province" label="省" fullWidth />
        <TextInput source="addr_city" label="市" fullWidth />
        <TextInput source="addr_street" label="街道地址" fullWidth />
        <TextInput source="bank_code" label="银行代码" fullWidth />
        <TextInput source="phone_num" label="电话" fullWidth />
      </SimpleForm>
    </Edit>
  );
};

export default CompanyEdit;
