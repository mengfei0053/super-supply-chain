import { Box, Typography } from "@mui/material";
import * as React from "react";
import {
  Datagrid,
  List,
  TextField,
  TopToolbar,
  EditButton,
  DeleteButton,
} from "react-admin";
import { useParams } from "react-router-dom";
import DataField from "./DataField";
import Upload from "./Upload";
import SingleExport from "./SingleExport";
import BatchExport from "./BatchExport";

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
      {/* <Export></Export>
      <ManageExportRule></ManageExportRule> */}
      <Upload type="excel"></Upload>
    </TopToolbar>
  );
};

const BulkActions = () => {
  return (
    <>
      <BatchExport></BatchExport>
    </>
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
      <Datagrid bulkActionButtons={<BulkActions></BulkActions>}>
        <TextField source="id"></TextField>
        <TextField source="fileName"></TextField>
        <TextField source="uploadFilePath"></TextField>
        <DataField source="datas"></DataField>
        <>
          <EditButton></EditButton>
          {tableName !== "dynamic_settlement_statement_suqian" && (
            <SingleExport></SingleExport>
          )}

          <DeleteButton></DeleteButton>
        </>
      </Datagrid>
    </List>
  );
};

export default ExcelPage;
