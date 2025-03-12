import * as React from "react";
import { DateRangePicker } from "@mui/x-date-pickers-pro/DateRangePicker";

interface IDateRangeInputProps {
  source: string;
}

const DateRangeInput: React.FunctionComponent<IDateRangeInputProps> = (
  props,
) => {
  return (
    <div>
      <DateRangePicker></DateRangePicker>
    </div>
  );
};

export default DateRangeInput;
