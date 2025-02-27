import * as React from "react";
import { useRecordContext } from "react-admin";
import { Button } from "@mui/material";
import Drawer from "@mui/material/Drawer";
import JSONPretty from "react-json-pretty";

interface IDataFieldProps {
  source: string;
}

const DataField: React.FunctionComponent<IDataFieldProps> = ({ source }) => {
  const record = useRecordContext();
  const [visible, setVisible] = React.useState(false);

  const toggle = (bool: boolean) => {
    setVisible(bool || !visible);
  };

  return (
    <div>
      <Button onClick={() => toggle(true)}>查看数据</Button>
      <Drawer
        anchor="right"
        sx={{ "& .MuiDrawer-paper": { width: "40%" } }}
        open={visible}
        onClose={() => toggle(false)}
      >
        <div
          style={{
            padding: "20px",
          }}
        >
          {record && <JSONPretty data={record[source]}></JSONPretty>}
        </div>
      </Drawer>
    </div>
  );
};

export default DataField;
