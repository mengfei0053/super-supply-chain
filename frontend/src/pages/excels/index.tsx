import { Box, Typography } from "@mui/material";
import * as React from "react";
import {
  Datagrid,
  List,
  TextField,
  TopToolbar,
  EditButton,
  DeleteButton,
  BulkDeleteButton,
  DateInput,
  FunctionField,
} from "react-admin";
import { useParams } from "react-router-dom";
import DataField from "./DataField";
import Upload from "./Upload";
import SingleExport from "./SingleExport";
import BatchExport from "./BatchExport";
import dayjs from "dayjs";

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
  const { tableName } = useParams();

  return (
    <>
      {(tableName === "dynamic_settlement_statement_suqian" ||
        tableName === "dynamic_Integrity_packaging_invoice") && (
        <BatchExport></BatchExport>
      )}
      <BulkDeleteButton></BulkDeleteButton>
      <></>
    </>
  );
};

const ExcelPage: React.FunctionComponent = () => {
  const { tableName } = useParams();

  const filterts = [
    <DateInput source="start" alwaysOn key={"start"}></DateInput>,
    <DateInput source="end" alwaysOn key={"end"}></DateInput>,
  ];

  return (
    <>
      <List
        resource={`excel/${tableName}`}
        actions={<ListAction></ListAction>}
        empty={<Empty></Empty>}
        filters={filterts}
        filterDefaultValues={{
          start: dayjs().format("YYYY-MM-DD"),
          end: dayjs().format("YYYY-MM-DD"),
        }}
      >
        <Datagrid bulkActionButtons={<BulkActions></BulkActions>}>
          <TextField source="id"></TextField>
          <TextField source="fileName"></TextField>
          <TextField source="uploadFilePath"></TextField>
          <FunctionField
            label="Created At"
            render={(r) => dayjs(r.CreatedAt).format("YYYY-MM-DD HH:mm:ss")}
          ></FunctionField>
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
    </>
  );
};

export default ExcelPage;
