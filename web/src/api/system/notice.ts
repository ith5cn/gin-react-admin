import request from "@/utils/request";

/**
 * 公告列表
 */
export const noticeListApi = (params?: Record<string, unknown>) => {
  return request.get("/system/notice/index", { params });
};

/**
 * 创建公告
 */
export const noticeCreateApi = (data: Record<string, unknown>) => {
  return request.post("/system/notice", data);
};

/**
 * 更新公告
 */
export const noticeUpdateApi = (id: number | string, data: Record<string, unknown>) => {
  return request.put(`/system/notice/${id}`, data);
};

/**
 * 删除公告
 */
export const noticeDeleteApi = (id: number | string) => {
  return request.delete(`/system/notice/${id}`);
};
