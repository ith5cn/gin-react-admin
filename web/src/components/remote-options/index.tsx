import { Cascader, Select, TreeSelect, message } from "antd";
import type { CascaderProps, SelectProps, TreeSelectProps } from "antd";
import { useEffect, useState, type ComponentType } from "react";
import request from "@/utils/request";

export interface RemoteOptionsProps {
  url: string;
  params?: Record<string, string | number | boolean>;
  dataPath?: string;
  labelField?: string;
  valueField?: string;
  childrenField?: string;
  valueType?: "string" | "number";
}

type OptionNode = { label: string; value: string | number; children?: OptionNode[] };

const CascaderControl = Cascader as unknown as ComponentType<Record<string, unknown>>;

const cache = new Map<string, OptionNode[]>();
const pending = new Map<string, Promise<OptionNode[]>>();

const readPath = (value: unknown, path?: string): unknown => {
  if (!path) return value;
  return path.split(".").filter(Boolean).reduce<unknown>((current, key) => {
    if (!current || typeof current !== "object") return undefined;
    return (current as Record<string, unknown>)[key];
  }, value);
};

const mapNodes = (
  rows: unknown,
  labelField: string,
  valueField: string,
  childrenField: string,
  valueType: "string" | "number",
): OptionNode[] => {
  if (!Array.isArray(rows)) return [];
  return rows.map((row) => {
    const item = (row || {}) as Record<string, unknown>;
    const rawValue = item[valueField];
    const value = valueType === "number" ? Number(rawValue) : String(rawValue ?? "");
    const children = mapNodes(item[childrenField], labelField, valueField, childrenField, valueType);
    return {
      label: String(item[labelField] ?? ""),
      value,
      children: children.length ? children : undefined,
    };
  });
};

const useRemoteOptions = ({
  url,
  params,
  dataPath,
  labelField = "label",
  valueField = "value",
  childrenField = "children",
  valueType = "string",
}: RemoteOptionsProps) => {
  const key = JSON.stringify({ url, params, dataPath, labelField, valueField, childrenField, valueType });
  const [options, setOptions] = useState<OptionNode[]>(cache.get(key) || []);
  const [loading, setLoading] = useState(!cache.has(key));

  useEffect(() => {
    let mounted = true;
    if (cache.has(key)) {
      setOptions(cache.get(key) || []);
      setLoading(false);
      return () => { mounted = false; };
    }
    let task = pending.get(key);
    if (!task) {
      task = request.get<unknown>(url, { params }).then((response) => {
        const next = mapNodes(readPath(response.data, dataPath), labelField, valueField, childrenField, valueType);
        cache.set(key, next);
        return next;
      });
      pending.set(key, task);
    }
    task
      .then((next) => { if (mounted) setOptions(next); })
      .catch((error) => { if (mounted) message.error(error?.message || "选项数据加载失败"); })
      .finally(() => {
        pending.delete(key);
        if (mounted) setLoading(false);
      });
    return () => { mounted = false; };
  }, [key]);

  return { options, loading };
};

export const RemoteSelect = ({ url, params, dataPath, labelField, valueField, childrenField, valueType, ...rest }: RemoteOptionsProps & Omit<SelectProps, "options" | "loading">) => {
  const state = useRemoteOptions({ url, params, dataPath, labelField, valueField, childrenField, valueType });
  return <Select allowClear showSearch optionFilterProp="label" {...rest} {...state} />;
};

export const RemoteTreeSelect = ({ url, params, dataPath, labelField, valueField, childrenField, valueType, ...rest }: RemoteOptionsProps & Omit<TreeSelectProps, "treeData" | "loading">) => {
  const state = useRemoteOptions({ url, params, dataPath, labelField, valueField, childrenField, valueType });
  return <TreeSelect allowClear showSearch {...rest} treeData={state.options} loading={state.loading} />;
};

export const RemoteCascader = ({ url, params, dataPath, labelField, valueField, childrenField, valueType, ...rest }: RemoteOptionsProps & Omit<CascaderProps, "options" | "loading">) => {
  const state = useRemoteOptions({ url, params, dataPath, labelField, valueField, childrenField, valueType });
  const { multiple, ...cascaderProps } = rest;
  if (multiple) {
    return <CascaderControl allowClear showSearch multiple {...cascaderProps} options={state.options} loading={state.loading} />;
  }
  return <CascaderControl allowClear showSearch {...cascaderProps} options={state.options} loading={state.loading} />;
};