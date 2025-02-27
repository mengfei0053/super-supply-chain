import * as React from "react";
import { Edit, SimpleForm, TextInput } from "react-admin";
import { useParams } from "react-router-dom";
import { JsonInput } from "react-admin-json-view";

const ExcelEditPage: React.FunctionComponent = () => {
  const { tableName } = useParams();

  return (
    <>
      <Edit resource={`excel/${tableName}`}>
        <SimpleForm>
          <TextInput source="fileName" disabled></TextInput>
          <TextInput source="uploadFilePath" disabled></TextInput>
          <JsonInput source="datas"></JsonInput>
        </SimpleForm>
      </Edit>
    </>
  );
};

export default ExcelEditPage;
