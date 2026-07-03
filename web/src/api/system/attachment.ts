import request from "@/utils/request";

export const attachmentListApi = (params?: any) => request.get("/system/attachment/index", { params });
// 表格选中的 rowKey 可能是 string，后端 ids 是 uint 数组，这里统一转数字。
export const attachmentDeleteApi = (data: { ids: Array<number | string>; removeSource?: boolean }) =>
  request.post("/system/attachment/delete", { ...data, ids: data.ids.map(Number) });
