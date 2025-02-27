import { Box, Stack, Typography } from "@mui/material";
import * as React from "react";
import {
  Datagrid,
  List,
  TextField,
  TopToolbar,
  ArrayField,
  FunctionField,
  WrapperField,
  TextInput,
  EditButton,
} from "react-admin";
import { useParams } from "react-router-dom";
import DataField from "./DataField";
import Upload from "./Upload";
import ManageExportRule from "./ManageExportRule";
import Export from "./Export";
import SingleExport from "./SingleExport";
import { JsonInput } from "react-admin-json-view";

const Empty = () => {
  return (
    <Box textAlign={"center"} m={1}>
      <Typography variant="h4">No Excel/test dynamic tables yet.</Typography>
      <Upload type="excel"></Upload>
    </Box>
  );
};

const ListAction = () => {
  return (
    <TopToolbar>
      <Export></Export>
      <ManageExportRule></ManageExportRule>
      <Upload type="excel"></Upload>
    </TopToolbar>
  );
};

const ExcelPage: React.FunctionComponent = () => {
  const { tableName } = useParams();

  return (
    <List
      resource={`excel/${tableName}`}
      actions={<ListAction></ListAction>}
      empty={<Empty></Empty>}
    >
      <Datagrid>
        <TextField source="id"></TextField>
        <TextField source="fileName"></TextField>
        <TextField source="uploadFilePath"></TextField>
        <DataField source="datas"></DataField>
        <>
          <EditButton></EditButton>
          <SingleExport></SingleExport>
        </>
      </Datagrid>
    </List>
  );
};

export default ExcelPage;
