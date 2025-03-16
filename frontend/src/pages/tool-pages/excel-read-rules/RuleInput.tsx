import * as React from "react";
import { useController } from "react-hook-form";
import {
  ArrayInput,
  Labeled,
  SimpleFormIterator,
  TextInput,
  useInput,
} from "react-admin";

interface IRuleInputProps {
  source: string;
}

const RuleInput: React.FunctionComponent<IRuleInputProps> = ({ source }) => {
  const { id, field, fieldState } = useInput({ source });

  console.log(id, field, fieldState);

  return (
    <Labeled label="Rules">
      <div>1</div>
    </Labeled>
  );
};

export default RuleInput;
