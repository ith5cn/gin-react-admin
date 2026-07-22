import { Input, InputNumber, Space } from "antd";

type RangeValue = [string | number | undefined, string | number | undefined];

interface QueryRangeProps {
  value?: RangeValue;
  onChange?: (value: RangeValue) => void;
  numeric?: boolean;
}

const QueryRange = ({ value = [undefined, undefined], onChange, numeric = false }: QueryRangeProps) => {
  const update = (index: number, next: string | number | null) => {
    const result: RangeValue = [...value];
    result[index] = next ?? undefined;
    onChange?.(result);
  };
  return (
    <Space.Compact>
      {numeric ? (
        <>
          <InputNumber value={value[0] as number | undefined} placeholder="最小值" onChange={(next) => update(0, next)} />
          <InputNumber value={value[1] as number | undefined} placeholder="最大值" onChange={(next) => update(1, next)} />
        </>
      ) : (
        <>
          <Input value={value[0] as string | undefined} placeholder="起始值" onChange={(event) => update(0, event.target.value)} />
          <Input value={value[1] as string | undefined} placeholder="结束值" onChange={(event) => update(1, event.target.value)} />
        </>
      )}
    </Space.Compact>
  );
};

export default QueryRange;
