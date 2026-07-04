import request from "@/utils/request";

/**
 * 角色列表
 */
export const roleListApi = (params?: Record<string, unknown>) => {
    return request.get("/system/role/index", { params })
}

/**
 * 新增角色
 */
export const roleCreateApi = (data: Record<string, unknown>) => {
    return request.post("/system/role/create", data)
}

/**
 * 更新角色
 */
export const roleUpdateApi = (id: string | number, data: Record<string, unknown>) => {
    return request.put(`/system/role/${id}`, data)
}

/**
 * 删除角色
 */
export const roleDeleteApi = (id: string | number) => {
    return request.delete(`/system/role/${id}`)
}

/**
 * 角色绑定菜单
 */
export const roleBindMenuApi = (id: string | number, data: { ids: number[] }) => {
    return request.post(`/system/role/${id}/menu`, data)
}

/**
 * 角色已授权的部门ID（自定义数据权限回显）
 */
export const roleDeptApi = (id: string | number) => {
    return request.get<number[]>(`/system/role/${id}/dept`)
}

/**
 * 允许访问的角色
 */
export const roleAccessApi = () => {
    return request.get("/system/role/access")
}