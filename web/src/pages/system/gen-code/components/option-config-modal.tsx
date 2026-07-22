import { Button, Divider, Input, Modal, Radio, Select, Space, Typography, message } from "antd";
import { DeleteOutlined, PlusOutlined } from "@ant-design/icons";
import { useEffect, useState } from "react";
import type { GencodeColumnRecord, OptionConfig, OptionNode, OptionRoute } from "./gen-code-modal";

interface OptionConfigModalProps {
  open: boolean;
  record?: GencodeColumnRecord;
  routes: OptionRoute[];
  onCancel: () => void;
  onSave: (patch: Partial<GencodeColumnRecord>) => void;
}

type ParamRow = { key: string; value: string };

const emptyNode = (): OptionNode => ({ label: "", value: "", children: [] });

const parseConfig = (raw?: string): OptionConfig => {
  if (!raw) return {};
  try {
    return JSON.parse(raw) as OptionConfig;
  } catch {
    return {};
  }
};

const isNumericColumn = (columnType?: string) =>
  /^(tinyint|smallint|mediumint|int|bigint|decimal|float|double)/i.test(columnType || "");

const parseScalar = (value: string): string | number | boolean => {
  const trimmed = value.trim();
  if (trimmed === "true") return true;
  if (trimmed === "false") return false;
  if (trimmed !== "" && Number.isFinite(Number(trimmed))) return Number(trimmed);
  return value;
};

const hasInvalidNode = (nodes: OptionNode[], numeric: boolean): boolean =>
  nodes.some((node) =>
    !node.label.trim() ||
    node.value === "" ||
    (numeric && Number.isNaN(Number(node.value))) ||
    (node.children?.length ? hasInvalidNode(node.children, numeric) : false),
  );

interface NodeEditorProps {
  nodes: OptionNode[];
  nested: boolean;
  numeric: boolean;
  depth?: number;
  onChange: (nodes: OptionNode[]) => void;
}

const NodeEditor = ({ nodes, nested, numeric, depth = 0, onChange }: NodeEditorProps) => {
  const updateNode = (index: number, patch: Partial<OptionNode>) => {
    const next = [...nodes];
    next[index] = { ...next[index], ...patch };
    onChange(next);
  };

  return (
    <Space direction="vertical" size={8} style={{ width: "100%" }}>
      {nodes.map((node, index) => (
        <div key={`${depth}-${index}`} style={{ marginLeft: depth * 24, padding: 10, border: "1px solid #e5e7eb", borderRadius: 8 }}>
          <Space wrap style={{ width: "100%" }}>
            <Input
              value={node.label}
              placeholder="显示文字"
              style={{ width: 180 }}
              onChange={(event) => updateNode(index, { label: event.target.value })}
            />
            <Input
              value={String(node.value ?? "")}
              placeholder={numeric ? "数字值" : "选项值"}
              style={{ width: 180 }}
              onChange={(event) => updateNode(index, { value: event.target.value })}
            />
            {nested ? (
              <Button size="small" icon={<PlusOutlined />} onClick={() => updateNode(index, { children: [...(node.children || []), emptyNode()] })}>
                子节点
              </Button>
            ) : null}
            <Button
              size="small"
              danger
              icon={<DeleteOutlined />}
              onClick={() => onChange(nodes.filter((_, nodeIndex) => nodeIndex !== index))}
            />
          </Space>
          {nested && node.children?.length ? (
            <div style={{ marginTop: 8 }}>
              <NodeEditor
                nodes={node.children}
                nested
                numeric={numeric}
                depth={depth + 1}
                onChange={(children) => updateNode(index, { children })}
              />
            </div>
          ) : null}
        </div>
      ))}
      <Button type="dashed" block icon={<PlusOutlined />} onClick={() => onChange([...nodes, emptyNode()])}>
        添加选项
      </Button>
    </Space>
  );
};

const OptionConfigModal = ({ open, record, routes, onCancel, onSave }: OptionConfigModalProps) => {
  const [source, setSource] = useState<"static" | "route">("static");
  const [config, setConfig] = useState<OptionConfig>({ options: [emptyNode()] });
  const [params, setParams] = useState<ParamRow[]>([]);

  useEffect(() => {
    if (!open || !record) return;
    const parsed = parseConfig(record.option_config);
    setSource(record.option_source === "route" ? "route" : "static");
    setConfig({
      options: parsed.options?.length ? parsed.options : [emptyNode()],
      path: parsed.path,
      dataPath: parsed.dataPath || "",
      labelField: parsed.labelField || "label",
      valueField: parsed.valueField || "value",
      childrenField: parsed.childrenField || "children",
    });
    setParams(Object.entries(parsed.params || {}).map(([key, value]) => ({ key, value: String(value) })));
  }, [open, record]);

  const handleSave = () => {
    if (!record) return;
    const numeric = isNumericColumn(record.column_type);
    if (source === "static") {
      const normalizeNodes = (nodes: OptionNode[]): OptionNode[] =>
        nodes.map((node) => ({
          label: node.label.trim(),
          value: numeric ? Number(node.value) : String(node.value),
          children: node.children?.length ? normalizeNodes(node.children) : undefined,
        }));
      const options = normalizeNodes(config.options || []);
      if (!options.length || hasInvalidNode(options, numeric)) {
        message.error("请完整填写静态选项的标签和值");
        return;
      }
      onSave({ option_source: "static", option_config: JSON.stringify({ options }) });
      return;
    }

    if (!config.path) {
      message.error("请选择数据源路由");
      return;
    }
    const paramObject = Object.fromEntries(params.filter((item) => item.key.trim()).map((item) => [item.key.trim(), parseScalar(item.value)]));
    onSave({
      option_source: "route",
      option_config: JSON.stringify({
        path: config.path,
        dataPath: config.dataPath || "",
        labelField: config.labelField || "label",
        valueField: config.valueField || "value",
        childrenField: config.childrenField || "children",
        params: paramObject,
      }),
    });
  };

  const nested = record?.view_type === "treeSelect" || record?.view_type === "cascader";

  return (
    <Modal title={`组件配置 - ${record?.column_comment || record?.column_name || ""}`} open={open} width={760} onCancel={onCancel} onOk={handleSave}>
      <Radio.Group value={source} optionType="button" onChange={(event) => setSource(event.target.value)}>
        <Radio value="static">静态选项</Radio>
        <Radio value="route">系统路由</Radio>
      </Radio.Group>
      <Divider />
      {source === "static" ? (
        <NodeEditor
          nodes={config.options || []}
          nested={nested}
          numeric={isNumericColumn(record?.column_type)}
          onChange={(options) => setConfig((current) => ({ ...current, options }))}
        />
      ) : (
        <Space direction="vertical" size={12} style={{ width: "100%" }}>
          <Select
            showSearch
            value={config.path}
            placeholder="选择无路径参数的系统 GET 路由"
            options={routes.map((route) => ({ label: route.label, value: route.path }))}
            onChange={(path) => setConfig((current) => ({ ...current, path }))}
          />
          <Space wrap>
            <Input value={config.dataPath} placeholder="数据路径，如 list" onChange={(event) => setConfig((current) => ({ ...current, dataPath: event.target.value }))} />
            <Input value={config.labelField} placeholder="标签字段" onChange={(event) => setConfig((current) => ({ ...current, labelField: event.target.value }))} />
            <Input value={config.valueField} placeholder="值字段" onChange={(event) => setConfig((current) => ({ ...current, valueField: event.target.value }))} />
            {nested ? <Input value={config.childrenField} placeholder="子节点字段" onChange={(event) => setConfig((current) => ({ ...current, childrenField: event.target.value }))} /> : null}
          </Space>
          <Typography.Text type="secondary">静态请求参数</Typography.Text>
          {params.map((item, index) => (
            <Space key={index}>
              <Input value={item.key} placeholder="参数名" onChange={(event) => setParams((rows) => rows.map((row, rowIndex) => rowIndex === index ? { ...row, key: event.target.value } : row))} />
              <Input value={item.value} placeholder="参数值" onChange={(event) => setParams((rows) => rows.map((row, rowIndex) => rowIndex === index ? { ...row, value: event.target.value } : row))} />
              <Button danger icon={<DeleteOutlined />} onClick={() => setParams((rows) => rows.filter((_, rowIndex) => rowIndex !== index))} />
            </Space>
          ))}
          <Button type="dashed" icon={<PlusOutlined />} onClick={() => setParams((rows) => [...rows, { key: "", value: "" }])}>添加参数</Button>
        </Space>
      )}
    </Modal>
  );
};

export default OptionConfigModal;
