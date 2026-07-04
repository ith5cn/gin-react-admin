import request from "@/utils/request";

export interface OnlineUserItem {
  userId: number;
  username: string;
  ip: string;
  os: string;
  browser: string;
  loginTime: string;
  accessJti: string;
}

/**
 * 在线用户列表
 */
export const onlineUserListApi = (params?: Record<string, unknown>) => {
  return request.get("/system/online/index", { params });
};

/**
 * 踢下线（accessJti 是该会话 access token 的唯一编号）
 */
export const kickOnlineUserApi = (accessJti: string) => {
  return request.delete(`/system/online/kick/${accessJti}`);
};
