import { Box, Typography, Button } from "@mui/material";
import { styled } from "@mui/material/styles";
import { Datagrid, List, TextField, TopToolbar } from "react-admin";
import CloudUpload from "@mui/icons-material/CloudUpload";
import { httpClient } from "../../../dataProvider";

const VisuallyHiddenInput = styled("input")({
  clip: "rect(0 0 0 0)",
  clipPath: "inset(50%)",
  height: 1,
  overflow: "hidden",
  position: "absolute",
  bottom: 0,
  left: 0,
  whiteSpace: "nowrap",
  width: 1,
});

const Actions = () => {
  return (
    <TopToolbar>
      <Button
        onClick={() => {
          console.log(123, "Create");
        }}
      >
        Create2
      </Button>
    </TopToolbar>
  );
};

const Empty = () => {
  const uploadFile = async (files: FileList) => {
    const file = files[0];
    const formData = new FormData();
    formData.append("file", file);
    formData.append("name", file.name);
    httpClient(
      `${import.meta.env.VITE_JSON_SERVER_URL}/settlement-form-entries`,
      {
        method: "POST",
        body: formData,
      },
    );
  };

  return (
    <Box textAlign="center" m={1}>
      <Typography variant="h4" paragraph>
        No products available
      </Typography>
      <Typography variant="body1">Create one or import from a file</Typography>

      <Button
        component="label"
        role={undefined}
        variant="contained"
        tabIndex={-1}
        startIcon={<CloudUpload />}
      >
        Upload files
        <VisuallyHiddenInput
          type="file"
          onChange={(event) => {
            console.log(event.target.files);
            uploadFile(event.target.files as FileList);
          }}
        />
      </Button>
    </Box>
  );
};

const SettlementFormEntryList = () => {
  return (
    <List actions={<Actions></Actions>} empty={<Empty></Empty>}>
      <Datagrid>
        <TextField source="id"></TextField>
        <TextField source="orderNumber"></TextField>
        <TextField source="arrivalPort"></TextField>
        <TextField source="arrivalDate"></TextField>
        <TextField source="count"></TextField>
        <TextField source="unit"></TextField>
        <TextField source="valueUsd"></TextField>
        <TextField source="valueRmb"></TextField>
        <TextField source="goodsName"></TextField>
      </Datagrid>
    </List>
  );
};

export default SettlementFormEntryList;
