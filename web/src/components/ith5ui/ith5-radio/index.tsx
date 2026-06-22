import { Radio } from "antd";
import type { RadioGroupProps } from "antd";
import { useEffect } from "react";
import useDictStore from "@/store/useDictStore";

interface Ith5RadioProps extends Omit<RadioGroupProps, "options"> {
  dict: string;
}

const Ith5Radio = ({ dict, ...rest }: Ith5RadioProps) => {
  const { dictMap, isLoaded, fetchDictAll } = useDictStore();

  useEffect(() => {
    if (!isLoaded) fetchDictAll();
  }, [isLoaded, fetchDictAll]);

  const options = (dictMap[dict] ?? []).map((item) => ({
    label: item.label,
    value: item.value,
  }));

  return <Radio.Group options={options} {...rest} />;
};

export default Ith5Radio;
