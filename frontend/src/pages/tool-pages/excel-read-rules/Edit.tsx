import {
  Edit,
  SimpleForm,
  TextInput,
  TopToolbar,
  SelectArrayInput,
  NumberInput,
} from "react-admin";
import * as React from "react";
import { Button } from "@mui/material";
import { httpClient } from "../../../dataProvider";
import AdSelectInput from "../../../components/AdSelectInput";

const ToolPagesEdit: React.FunctionComponent = () => {
  const [options, setOptions] = React.useState<
    {
      id: number;
      name: string;
    }[]
  >([]);

  const getOptions = async () => {
    const res = await httpClient(
      import.meta.env.VITE_JSON_SERVER_URL + "/options/excel-mapping-rule",
    );
    setOptions(res.json);
  };

  React.useEffect(() => {
    getOptions();
  }, []);

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
        <TextInput source="menuName"></TextInput>
        <TextInput source="desc"></TextInput>
        <TextInput source="dynamicTableName"></TextInput>
        <NumberInput min={0} step={1} source="sheetIndex"></NumberInput>
        <AdSelectInput
          source="mapRuleId"
          URL="excel-mapping-rule"
        ></AdSelectInput>
        <AdSelectInput
          source="iterateRuleId"
          URL="excel-mapping-rule"
        ></AdSelectInput>
      </SimpleForm>
    </Edit>
  );
};

export default ToolPagesEdit;
