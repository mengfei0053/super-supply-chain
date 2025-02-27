import * as React from "react";
import { Button, List, Datagrid, TextField, TopToolbar } from "react-admin";
import Drawer from "@mui/material/Drawer";
import { Box, Typography, Modal } from "@mui/material";
import Upload from "./Upload";
import { useParams } from "react-router-dom";
import { AddOutlined } from "@mui/icons-material";
import CreateTemplate from "./CreateTemplate";

const Empty = () => {
  return (
    <Box textAlign={"center"} m={1}>
      <Typography variant="h4">创建导出模板</Typography>
      <CreateTemplate></CreateTemplate>
      {/* <Upload type="excel-export-rule"></Upload> */}
    </Box>
  );
};

const ManageExportRule: React.FunctionComponent = () => {
  const [visible, setVisible] = React.useState(false);

  const { tableName } = useParams();

  const toggle = (bool: boolean) => {
    setVisible(bool || !visible);
  };

  return (
    <>
      <Button onClick={() => toggle(true)} label="创建导出规则"></Button>
      <Drawer
        anchor="right"
        sx={{ "& .MuiDrawer-paper": { width: "40%" } }}
        open={visible}
        onClose={() => toggle(false)}
      >
        <Box sx={{ padding: "20px" }}>
          <List
            resource={`excel-export-rule/template/${tableName}`}
            empty={<Empty></Empty>}
            actions={
              <TopToolbar>
                <CreateTemplate></CreateTemplate>
              </TopToolbar>
            }
          >
            <Datagrid>
              <TextField source="id"></TextField>
              <TextField source="alias"></TextField>
              <TextField source="associatedTable"></TextField>
              <TextField source="uploadFilePath"></TextField>
            </Datagrid>
          </List>
        </Box>
      </Drawer>
    </>
  );
};

export default ManageExportRule;
