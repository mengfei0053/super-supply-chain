import * as React from "react";
import { Create, CreateButton, SelectInput, SimpleForm } from "react-admin";
import { useParams } from "react-router-dom";
import Drawer from "@mui/material/Drawer";
import { Box } from "@mui/material";
import { CloudDownloadOutlined } from "@mui/icons-material";
import AdSelectInput from "../../components/AdSelectInput";

const Export: React.FunctionComponent = () => {
  const [visible, setVisible] = React.useState(false);
  const choices = ["MAP_RULE", "ITERATE_RULE"];

  const { tableName } = useParams();

  const toggle = (bool: boolean) => {
    setVisible(bool || !visible);
  };

  return (
    <>
      <CreateButton
        label="Export"
        icon={<CloudDownloadOutlined></CloudDownloadOutlined>}
        onClick={(e) => {
          e.preventDefault();
          e.stopPropagation();
          toggle(true);
        }}
      ></CreateButton>
      <Drawer
        anchor="right"
        sx={{ "& .MuiDrawer-paper": { width: "40%" } }}
        open={visible}
        onClose={() => toggle(false)}
      >
        <Box sx={{ padding: "20px" }}>
          <Create
            resource={`excel-export-rule/${tableName}/export`}
            mutationOptions={{
              onSuccess() {
                toggle(false);
              },
            }}
          >
            <SimpleForm>
              <AdSelectInput
                source="templateId"
                required
                URL="export-templates"
                params={{
                  associated_table: "dynamic_settlement_statement",
                }}
              ></AdSelectInput>
            </SimpleForm>
          </Create>
        </Box>
      </Drawer>
    </>
  );
};

export default Export;
