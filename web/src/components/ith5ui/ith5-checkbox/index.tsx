import { Checkbox } from "antd";
import type { CheckboxGroupProps } from "antd/es/checkbox";
import { useEffect } from "react";
import useDictStore from "@/store/useDictStore";

interface Ith5CheckboxProps extends Omit<CheckboxGroupProps, "options"> {
  dict: string;
  valueType?: "string" | "number";
}

const Ith5Checkbox = ({ dict, valueType = "string", ...rest }: Ith5CheckboxProps) => {
  const { dictMap, isLoaded, fetchDictAll } = useDictStore();

  useEffect(() => {
    if (!isLoaded) fetchDictAll();
  }, [isLoaded, fetchDictAll]);

  const options = (dictMap[dict] ?? []).map((item) => ({
    label: item.label,
    value: valueType === "number" ? Number(item.value) : item.value,
  }));

  return <Checkbox.Group options={options} {...rest} />;
};

export default Ith5Checkbox;
