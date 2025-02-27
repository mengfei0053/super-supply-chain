import { styled } from "@mui/material/styles";
import { httpClient } from "../../dataProvider";
import { useParams } from "react-router-dom";
import { CreateButton, useRefresh } from "react-admin";
import React from "react";

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

export default function Upload({
  type,
}: {
  type: "excel" | "excel-export-rule";
}) {
  const { tableName } = useParams();
  const ref = React.useRef<HTMLInputElement>(null);
  const refresh = useRefresh();

  const uploadFile = async (files: FileList) => {
    const file = files[0];
    const formData = new FormData();
    formData.append("file", file);
    formData.append("name", file.name);
    httpClient(`${import.meta.env.VITE_JSON_SERVER_URL}/${type}/${tableName}`, {
      method: "POST",
      body: formData,
    }).then(() => {
      refresh();
    });
  };

  return (
    <>
      <CreateButton
        onClick={(e) => {
          e.preventDefault();
          e.stopPropagation();
          ref.current?.click();
        }}
      ></CreateButton>

      <VisuallyHiddenInput
        ref={ref}
        type="file"
        onChange={(event) => {
          uploadFile(event.target.files as FileList);
        }}
      />
    </>
  );
}
