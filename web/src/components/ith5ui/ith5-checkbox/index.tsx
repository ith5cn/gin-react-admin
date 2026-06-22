import { Checkbox } from "antd";
import type { CheckboxGroupProps } from "antd/es/checkbox";
import { useEffect } from "react";
import useDictStore from "@/store/useDictStore";

interface Ith5CheckboxProps extends Omit<CheckboxGroupProps, "options"> {
  dict: string;
}

const Ith5Checkbox = ({ dict, ...rest }: Ith5CheckboxProps) => {
  const { dictMap, isLoaded, fetchDictAll } = useDictStore();

  useEffect(() => {
    if (!isLoaded) fetchDictAll();
  }, [isLoaded, fetchDictAll]);

  const options = (dictMap[dict] ?? []).map((item) => ({
    label: item.label,
    value: item.value,
  }));

  return <Checkbox.Group options={options} {...rest} />;
};

export default Ith5Checkbox;
