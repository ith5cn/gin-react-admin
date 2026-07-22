import request from "@/utils/request";

export const loginLogListApi = (params?: Record<string, unknown>) => request.get("/system/login-log/index", { params });
export const loginLogDeleteApi = (id: string | number) => request.delete(`/system/login-log/${id}`);
