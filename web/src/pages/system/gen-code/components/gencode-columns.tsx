import { Button, Checkbox, Input, InputNumber, Select, Table, Tag } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useState } from "react";
import type { ComponentCapability, GencodeColumnRecord, OptionRoute } from "./gen-code-modal";
import OptionConfigModal from "./option-config-modal";
import { queryType } from "../js/var";

interface GencodeColumnsProps {
  data: GencodeColumnRecord[];
  onChange: (rows: GencodeColumnRecord[]) => void;
  dictTypeOptions: Array<{ label: string; value: string }>;
  componentOptions: ComponentCapability[];
  optionRoutes: OptionRoute[];
}

const DICT_VIEW_TYPES = ["saSelect", "radio", "checkbox"];
const OPTION_VIEW_TYPES = ["select", "treeSelect", "cascader"];

const isColumnConfigRequired = (record: GencodeColumnRecord) =>
  record.is_insert === 2 || record.is_edit === 2 || record.is_query === 2;

const hasInvalidOptionNode = (nodes: Array<{ label?: unknown; value?: unknown; children?: unknown }> = []): boolean =>
  nodes.some((node) => {
    if (!String(node.label ?? "").trim() || node.value === undefined || node.value === null || node.value === "") return true;
    return Array.isArray(node.children) && hasInvalidOptionNode(node.children);
  });

export const getColumnConfigError = (record: GencodeColumnRecord, capabilities: ComponentCapability[]) => {
  const capability = capabilities.find((item) => item.value === record.view_type);
  if (!capability) return "不支持的页面组件";
  if (!isColumnConfigRequired(record)) return "";
  if (record.is_query === 2 && capability.queryDisabled) return "该组件不支持查询";
  if (capability.configType === "dict" && !record.dict_type) return "请选择数据字典";
  if (capability.configType === "options") {
    if (!record.option_source || !record.option_config) return "请完善组件数据源";
    try {
      const config = JSON.parse(record.option_config);
      if (record.option_source === "static" && (!config.options?.length || hasInvalidOptionNode(config.options))) return "请完整填写静态选项";
      if (record.option_source === "route" && !config.path) return "请选择数据源路由";
    } catch {
      return "组件配置格式错误";
    }
  }
  return "";
};

const GencodeColumns = ({ data, onChange, dictTypeOptions, componentOptions, optionRoutes }: GencodeColumnsProps) => {
  const [configIndex, setConfigIndex] = useState<number | null>(null);

  const updateRow = (index: number, patch: Partial<GencodeColumnRecord>) => {
    const next = [...data];
    const viewChanged = patch.view_type && patch.view_type !== next[index].view_type;
    next[index] = { ...next[index], ...patch };
    if (viewChanged) {
      if (!DICT_VIEW_TYPES.includes(next[index].view_type)) next[index].dict_type = undefined;
      if (!OPTION_VIEW_TYPES.includes(next[index].view_type)) {
        next[index].option_source = undefined;
        next[index].option_config = undefined;
      }
      const capability = componentOptions.find((item) => item.value === next[index].view_type);
      if (capability?.queryDisabled) {
        next[index].is_query = 1;
        next[index].query_type = "eq";
      }
    }
    onChange(next);
  };

  const toggleFlag = (index: number, key: keyof GencodeColumnRecord, checked: boolean) => {
    if (key === "is_edit") {
      updateRow(index, { is_insert: checked ? 2 : 1, is_edit: checked ? 2 : 1 });
      return;
    }
    updateRow(index, { [key]: checked ? 2 : 1 } as Partial<GencodeColumnRecord>);
  };

  const columns: ColumnsType<GencodeColumnRecord> = [
    { title: "排序", dataIndex: "sort", key: "sort", width: 90, render: (value, _record, index) => <InputNumber value={value} min={0} onChange={(next) => updateRow(index, { sort: Number(next || 0) })} /> },
    { title: "字段名称", dataIndex: "column_name", key: "column_name", width: 180, render: (value) => <Input value={value} disabled /> },
    { title: "字段描述", dataIndex: "column_comment", key: "column_comment", width: 180, render: (value, _record, index) => <Input value={value} onChange={(event) => updateRow(index, { column_comment: event.target.value })} /> },
    { title: "字段类型", dataIndex: "column_type", key: "column_type", width: 140, render: (value) => <Input value={value} disabled /> },
    { title: "必填", dataIndex: "is_required", key: "is_required", width: 80, render: (value, _record, index) => <Checkbox checked={value === 2} onChange={(event) => toggleFlag(index, "is_required", event.target.checked)} /> },
    { title: "表单", dataIndex: "is_edit", key: "is_edit", width: 80, render: (value, record, index) => <Checkbox checked={value === 2} disabled={record.is_pk === 2} onChange={(event) => toggleFlag(index, "is_edit", event.target.checked)} /> },
    { title: "列表", dataIndex: "is_list", key: "is_list", width: 80, render: (value, _record, index) => <Checkbox checked={value === 2} onChange={(event) => toggleFlag(index, "is_list", event.target.checked)} /> },
    {
      title: "查询", dataIndex: "is_query", key: "is_query", width: 80,
      render: (value, record, index) => {
        const disabled = componentOptions.find((item) => item.value === record.view_type)?.queryDisabled;
        return <Checkbox checked={value === 2} disabled={disabled} onChange={(event) => toggleFlag(index, "is_query", event.target.checked)} />;
      },
    },
    { title: "排序字段", dataIndex: "is_sort", key: "is_sort", width: 100, render: (value, _record, index) => <Checkbox checked={value === 2} onChange={(event) => toggleFlag(index, "is_sort", event.target.checked)} /> },
    { title: "查询方式", dataIndex: "query_type", key: "query_type", width: 130, render: (value, record, index) => <Select style={{ width: 110 }} value={value} disabled={record.is_query !== 2} options={queryType} onChange={(next) => updateRow(index, { query_type: next })} /> },
    { title: "页面组件", dataIndex: "view_type", key: "view_type", width: 160, render: (value, _record, index) => <Select style={{ width: 140 }} value={value} options={componentOptions.map((item) => ({ label: item.label, value: item.value }))} onChange={(next) => updateRow(index, { view_type: next })} /> },
    {
      title: "数据字典", dataIndex: "dict_type", key: "dict_type", width: 180,
      render: (value, record, index) => <Select style={{ width: 160 }} value={value} status={isColumnConfigRequired(record) && DICT_VIEW_TYPES.includes(record.view_type) && !value ? "error" : undefined} allowClear disabled={!DICT_VIEW_TYPES.includes(record.view_type)} options={dictTypeOptions} onChange={(next) => updateRow(index, { dict_type: next })} />,
    },
    {
      title: "组件配置", key: "component_config", width: 150,
      render: (_value, record, index) => {
        if (!OPTION_VIEW_TYPES.includes(record.view_type)) return "-";
        const required = isColumnConfigRequired(record);
        const error = getColumnConfigError(record, componentOptions);
        const configured = Boolean(record.option_source && record.option_config);
        return <Button size="small" danger={required && Boolean(error)} onClick={() => setConfigIndex(index)}>{required && error ? "待完善" : configured ? "已配置" : "配置"}</Button>;
      },
    },
    {
      title: "状态", key: "config_status", width: 130,
      render: (_value, record) => {
        const error = getColumnConfigError(record, componentOptions);
        if (!error && !isColumnConfigRequired(record)) return <Tag>无需配置</Tag>;
        return error ? <Tag color="error">{error}</Tag> : <Tag color="success">可生成</Tag>;
      },
    },
  ];

  return (
    <>
      <Table rowKey={(record) => String(record.id ?? record.column_name)} columns={columns} dataSource={data} pagination={false} scroll={{ x: 1850 }} />
      <OptionConfigModal
        open={configIndex !== null}
        record={configIndex === null ? undefined : data[configIndex]}
        routes={optionRoutes}
        onCancel={() => setConfigIndex(null)}
        onSave={(patch) => {
          if (configIndex !== null) updateRow(configIndex, patch);
          setConfigIndex(null);
        }}
      />
    </>
  );
};

export default GencodeColumns;
