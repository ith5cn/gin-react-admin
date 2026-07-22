import { userAuthListApi } from "@/api/system/user";
import { Select, message } from "antd";
import type { SelectProps } from "antd";
import { useEffect, useState } from "react";

type UserValue = string | number;
type UserOption = { label: string; value: UserValue };

interface AuthSelectProps extends Omit<SelectProps<UserValue>, "options" | "loading" | "value" | "onChange"> {
  value?: UserValue;
  onChange?: (value?: UserValue) => void;
}

const AuthSelect = ({ value, onChange, ...rest }: AuthSelectProps) => {
  const [options, setOptions] = useState<UserOption[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let mounted = true;
    userAuthListApi()
      .then((res) => { if (mounted) setOptions((res.data || []) as UserOption[]); })
      .catch((error) => { if (mounted) message.error(error?.message || "用户选项加载失败"); })
      .finally(() => { if (mounted) setLoading(false); });
    return () => { mounted = false; };
  }, []);

  return (
    <Select<UserValue>
      allowClear
      showSearch
      optionFilterProp="label"
      placeholder="请选择用户"
      {...rest}
      options={options}
      loading={loading}
      value={value}
      onChange={(next) => onChange?.(next)}
    />
  );
};

export default AuthSelect;
