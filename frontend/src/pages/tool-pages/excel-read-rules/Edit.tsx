import {
  Edit,
  SimpleForm,
  TextInput,
  TopToolbar,
  SelectArrayInput,
  NumberInput,
  ArrayInput,
  SimpleFormIterator,
} from "react-admin";
import * as React from "react";
import { Button } from "@mui/material";
import { httpClient } from "../../../dataProvider";
import AdSelectInput from "../../../components/AdSelectInput";
import RuleInput from "./RuleInput";

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
        <TextInput source="dynamicTableName"></TextInput>
        <NumberInput min={0} step={1} source="sheetIndex"></NumberInput>

        <ArrayInput source="rules.mapRule">
          <SimpleFormIterator inline>
            <TextInput source="desc" helperText={false}></TextInput>
            <TextInput source="excelKey" helperText={false}></TextInput>
            <TextInput source="jsonKey" helperText={false}></TextInput>
          </SimpleFormIterator>
        </ArrayInput>

        <NumberInput
          min={0}
          step={1}
          source="rules.iterateRule.startRow"
        ></NumberInput>

        <ArrayInput source="rules.iterateRule.rules">
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
