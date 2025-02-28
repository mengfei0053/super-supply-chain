import * as React from "react";
import { BulkExportButton, useListContext } from "react-admin";
import { useParams } from "react-router-dom";
import qs from "qs";

const BatchExport: React.FunctionComponent = () => {
  const { tableName } = useParams();
  const { selectedIds } = useListContext();

  return (
    <>
      <BulkExportButton
        onClick={(e) => {
          e.stopPropagation();
          e.preventDefault();
          const query = qs.stringify(
            {
              ids: selectedIds,
            },
            {
              encode: false,
              arrayFormat: "comma",
            },
          );
          console.log(query);

          window.open(
            `${import.meta.env.VITE_JSON_SERVER_URL}/excel/${tableName}/exports?${query}`,
          );
        }}
      ></BulkExportButton>
    </>
  );
};

export default BatchExport;
