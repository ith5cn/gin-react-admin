import request from "@/utils/request";

export type ConfigGroupPayload = {
  id?: number;
  name?: string;
  code?: string;
  sort?: number;
  remark?: string;
};

export type ConfigItemPayload = {
  id?: number;
  groupId?: number;
  key?: string;
  value?: unknown;
  name?: string;
  inputType?: string;
  configSelectData?: string;
  sort?: number;
  remark?: string;
};

export const getConfigGroupListApi = (params?: Record<string, unknown>) => {
  return request.get("/system/config-group/index", { params });
};

export const configGroupCreateApi = (data: ConfigGroupPayload) => {
  return request.post("/system/config-group", data);
};

export const configGroupUpdateApi = (id: number | string, data: ConfigGroupPayload) => {
  return request.put(`/system/config-group/${id}`, data);
};

export const configGroupDeleteApi = (id: number | string) => {
  return request.delete(`/system/config-group/${id}`);
};

export const getConfigListApi = (params?: Record<string, unknown>) => {
  return request.get("/system/config/index", { params });
};

export const configCreateApi = (data: ConfigItemPayload) => {
  return request.post("/system/config", data);
};

export const configUpdateApi = (id: number | string, data: ConfigItemPayload) => {
  return request.put(`/system/config/${id}`, data);
};

export const configDeleteApi = (id: number | string) => {
  return request.delete(`/system/config/${id}`);
};

export const configBatchUpdateApi = (data: { groupId: number; config: ConfigItemPayload[] }) => {
  return request.post("/system/config/batch-update", data);
};

export const getConfigInfoApi = (code: string) => {
  return request.get("/system/config/get-config-info", { params: { code } });
};
