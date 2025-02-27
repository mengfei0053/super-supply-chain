import * as React from "react";
import { SelectArrayInput, SelectInput } from "react-admin";
import { httpClient } from "../dataProvider";
import qs from "qs";

interface IAdSelectInputProps {
  URL: string;
  source: string;
  isArraySelectInput?: boolean;
  params?: Record<string, any>;
  required?: boolean;
}

const AdSelectInput: React.FunctionComponent<IAdSelectInputProps> = ({
  URL: URL,
  source,
  isArraySelectInput,
  params,
}) => {
  const [options, setOptions] = React.useState<
    {
      id: number;
      name: string;
    }[]
  >([]);

  const getOptions = React.useCallback(async () => {
    if (!URL) return;

    const res = await httpClient(
      import.meta.env.VITE_JSON_SERVER_URL +
        `/options/${URL}${params ? `?${qs.stringify(params)}` : ""}`,
      {},
    );
    setOptions(res.json);
  }, [URL, params]);

  React.useEffect(() => {
    getOptions();
  }, [getOptions]);
  if (isArraySelectInput) {
    <SelectArrayInput source={source} choices={options}></SelectArrayInput>;
  }

  return <SelectInput source={source} choices={options}></SelectInput>;
};

export default AdSelectInput;
