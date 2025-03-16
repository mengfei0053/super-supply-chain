import {
  Create,
  SimpleForm,
  TextInput,
  useStore,
  NumberInput,
  ArrayInput,
  SimpleFormIterator,
} from "react-admin";
import * as React from "react";
import { useNavigate } from "react-router";
import { httpClient } from "../../../dataProvider";
import AdSelectInput from "../../../components/AdSelectInput";

const ToolPageCreate: React.FunctionComponent = () => {
  const navigate = useNavigate();

  const [, setMenus] = useStore<
    { id: number; dynamicTableName: string; menuName: string }[]
  >("global.menu", []);

  const getMenus = async () => {
    const res = await httpClient(
      import.meta.env.VITE_JSON_SERVER_URL + "/menus",
    );
    setMenus(res.json);
  };

  return (
    <Create
      mutationOptions={{
        onSuccess() {
          navigate("/excel-read-rules");

          getMenus();
        },
      }}
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
    </Create>
  );
};

export default ToolPageCreate;
