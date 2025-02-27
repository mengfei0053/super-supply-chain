import * as React from "react";
import { useController } from "react-hook-form";

interface IDatasInputProps {
  source: string;
}

const DatasInput: React.FunctionComponent<IDatasInputProps> = () => {
  const input1 = useController({
    name: "datas.baseData.name",
    defaultValue: "",
  });
  const input2 = useController({
    name: "datas.baseData.total_count",
    defaultValue: "",
  });

  return (
    <>
      <div>{input1.field.value}</div>
      <div>{input2.field.value}</div>
    </>
  );
};

export default DatasInput;
