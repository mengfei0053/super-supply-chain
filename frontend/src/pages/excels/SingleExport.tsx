import * as React from "react";
import { Button, useRecordContext } from "react-admin";
import { useParams } from "react-router-dom";
import { httpClient } from "../../dataProvider";

const SingleExport: React.FunctionComponent = () => {
  const record = useRecordContext();
  const { tableName } = useParams();

  const exportExcel = () => {
    window.open(
      `${import.meta.env.VITE_JSON_SERVER_URL}/excel-export-rule/${tableName}/export/${record?.id}`,
    );
  };

  return (
    <>
      <Button
        label="Export"
        onClick={() => {
          exportExcel();
        }}
      ></Button>
    </>
  );
};

export default SingleExport;
