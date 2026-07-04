import request from "@/utils/request";

export const userListApi = (params?: Record<string, unknown>) => {
    return request.get("/system/user/index", { params })
}

export const userUpdateApi = (id: string | number, data: Record<string, unknown>) => {
    return request.put(`/system/user/${id}`, data)
}

export const userAddApi = (data: Record<string, unknown>) => {
    return request.post("/system/user", data)
}

export const userDeleteApi = (id: string | number) => {
    return request.delete(`/system/user/${id}`)
}

export const userRefreshCacheApi = (id: string | number) => {
    return request.put(`/system/user/${id}/refresh-cache`)
}

export const userAuthListApi = () => {
    return request.get(`/system/user/auth-list`)
}

export const getCurrentUserApi = () => {
    return request.get(`/system/user`)
}
// 个人中心：更新自己的资料（后端以 JWT 里的用户为准，不传 id）。
export const updateProfileApi = (data: {
    nickname?: string;
    phone?: string;
    email?: string;
    avatar?: string;
    signed?: string;
}) => {
    return request.put("/system/user/profile", data)
}

// 个人中心：修改自己的密码，需携带原密码。
export const changePasswordApi = (data: { oldPassword: string; newPassword: string }) => {
    return request.put("/system/user/password", data)
}
