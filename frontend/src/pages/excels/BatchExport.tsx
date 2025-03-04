import * as React from "react";
import { useListContext, Button, useNotify } from "react-admin";
import { useParams } from "react-router-dom";
import Download from "@mui/icons-material/Download";
import qs from "qs";

const BatchExport: React.FunctionComponent = () => {
  const { tableName } = useParams();
  const { selectedIds } = useListContext();
  const notice = useNotify();

  const validate = React.useCallback(() => {
    if (selectedIds.length === 0) {
      notice("请至少选择一条记录", { type: "error" });
      return false;
    }
    return true;
  }, [notice, selectedIds.length]);

  return (
    <>
      <Button
        label="导出发票-运费"
        startIcon={<Download></Download>}
        onClick={(e) => {
          e.stopPropagation();
          e.preventDefault();
          if (validate()) {
            const query = qs.stringify(
              {
                ids: selectedIds,
                type: "invoice_freight",
              },
              {
                encode: false,
                arrayFormat: "repeat",
              },
            );

            window.open(
              `${import.meta.env.VITE_JSON_SERVER_URL}/excel-exports/${tableName}?${query}`,
            );
          }
        }}
      ></Button>
      <Button
        label="导出发票-清关-掏箱"
        startIcon={<Download></Download>}
        onClick={(e) => {
          e.stopPropagation();
          e.preventDefault();
          if (validate()) {
            const query = qs.stringify(
              {
                ids: selectedIds,
                type: "invoice_clearance",
              },
              {
                encode: false,
                arrayFormat: "comma",
              },
            );

            window.open(
              `${import.meta.env.VITE_JSON_SERVER_URL}/excel-exports/${tableName}?${query}`,
            );
          }
        }}
      ></Button>
      <Button
        label="导出-短驳费表"
        startIcon={<Download></Download>}
        onClick={(e) => {
          e.stopPropagation();
          e.preventDefault();
          if (validate()) {
            const query = qs.stringify(
              {
                ids: selectedIds,
                type: "shortHaul",
              },
              {
                encode: false,
                arrayFormat: "comma",
              },
            );

            window.open(
              `${import.meta.env.VITE_JSON_SERVER_URL}/excel-exports/${tableName}?${query}`,
            );
          }
        }}
      ></Button>
    </>
  );
};

export default BatchExport;
