import * as React from "react";
import {
  Button,
  Create,
  FileField,
  FileInput,
  SimpleForm,
  TextInput,
} from "react-admin";
import { Modal } from "@mui/material";
import { AddOutlined } from "@mui/icons-material";
import { useParams } from "react-router-dom";

export interface CreateTemplateParams {
  alias: string;
  file: {
    rawFile?: File;
    src: string;
    title: string;
  };
}
const CreateTemplate: React.FunctionComponent = () => {
  const [open, setOpen] = React.useState(false);
  const { tableName } = useParams();

  return (
    <>
      <Button
        label="CREATE"
        onClick={() => setOpen(true)}
        startIcon={<AddOutlined></AddOutlined>}
      ></Button>
      <Modal open={open} onClose={() => setOpen(false)}>
        <Create
          mutationOptions={{
            onSuccess: () => {
              setOpen(false);
            },
          }}
          resource={`excel-export-rule/template/${tableName}`}
        >
          <SimpleForm>
            <TextInput source="alias"></TextInput>
            <FileInput source="file">
              <FileField source="src" title="title" />
            </FileInput>
          </SimpleForm>
        </Create>
      </Modal>
    </>
  );
};

export default CreateTemplate;
