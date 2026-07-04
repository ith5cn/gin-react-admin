import request from '@/utils/request'
import type { RawBackendMenuNode } from '@/types/router'

export interface RegisterPayload {
  username: string
  email: string
  password: string
  nickname?: string
  phone?: string
}

export interface LoginPayload {
  username: string
  password: string
}

export interface LoginResult {
  accessToken: string
  refreshToken: string
  expiresIn: string
}

export interface BackendTokenPair {
  access_token: string
  refresh_token: string
  expires_in: number
}

export interface UserInfo {
  id: number
  username: string
  user_type?: string
  nickname?: string | null
  phone?: string | null
  email?: string | null
  avatar?: string | null
  signed?: string | null
  deptId?: number | null
  postId?: number | null
  status?: number
  dashboard?: string | null
  backendSetting?: string | null
  [key: string]: unknown
}

export interface UserContextResponse {
  user: UserInfo
  roles: Array<number | string>
  routers: RawBackendMenuNode[]
  codes: string | string[]
  posts: unknown
  depts: unknown
}

export const registerApi = (data: RegisterPayload) => {
  return request.post('/auth/system/register', data)
}

export const toLoginResult = (tokenPair: BackendTokenPair): LoginResult => ({
  accessToken: tokenPair.access_token,
  refreshToken: tokenPair.refresh_token,
  expiresIn: String(tokenPair.expires_in),
})

export const loginApi = async (data: LoginPayload) => {
  const res = await request.post<BackendTokenPair>('/base/login', {
    user_name: data.username,
    password: data.password,
  })
  return toLoginResult(res.data)
}

export const refreshTokenApi = async (refreshToken: string) => {
  const res = await request.post<BackendTokenPair>('/base/token/refresh', {
    refresh_token: refreshToken,
  })
  return toLoginResult(res.data)
}

export const logoutApi = async (refreshToken?: string) => {
  return request.post('/base/logout', refreshToken ? { refresh_token: refreshToken } : {})
}

export const userApi = async () => {
  const res = await request.get<UserContextResponse>('/system/user')
  return res.data
}
