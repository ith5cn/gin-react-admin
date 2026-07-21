import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'
import { message as Message } from 'antd'
import { clearAuthSession, useAuthStore } from '@/store/auth'
import { navigateTo } from '@/utils/navigateHelper'

export interface Result<T = any> {
  code: number
  message: string
  msg?: string
  data: T
  total?: number
}

const SUCCESS_CODES = [0, 200]
const LOGIN_REQUIRED_CODES = [401, 402, 40102]
const ACCESS_TOKEN_EXPIRED_CODE = 40101
const SYSTEM_NOT_INSTALLED_CODE = 50301

type RetryRequestConfig = AxiosRequestConfig & {
  _retry?: boolean
}

const instance: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_APP_BASE_API || '/api',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json;charset=utf-8' },
})

const normalizeMessage = (payload: any) => payload?.message || payload?.msg || '服务异常'

let refreshPromise: Promise<string> | null = null

const redirectToLogin = () => {
  clearAuthSession()
  // navigateTo 使用 React Router navigate（SPA 跳转，无全页刷新）。
  // 对于 Layout 外的路由（如 /install），Layout 不会自动接管，因此需要主动跳转。
  navigateTo('/login')
}

const redirectToInstall = () => {
  if (window.location.pathname !== '/install') {
    navigateTo('/install')
  }
}

const updateTokenSession = (tokenPair: any) => {
  useAuthStore.setState({
    token: tokenPair.access_token,
    refreshToken: tokenPair.refresh_token,
    expiresIn: String(tokenPair.expires_in ?? ''),
  })
}

const refreshAccessToken = async () => {
  const refreshToken = useAuthStore.getState().refreshToken
  if (!refreshToken) {
    throw new Error('refresh token missing')
  }

  if (!refreshPromise) {
    refreshPromise = axios
      .post(
        '/base/token/refresh',
        { refresh_token: refreshToken },
        {
          baseURL: import.meta.env.VITE_APP_BASE_API || '/api',
          timeout: 15000,
          headers: { 'Content-Type': 'application/json;charset=utf-8' },
        },
      )
      .then((res) => {
        const payload = res.data
        if (!SUCCESS_CODES.includes(payload?.code)) {
          throw new Error(normalizeMessage(payload))
        }

        updateTokenSession(payload.data)
        return payload.data.access_token as string
      })
      .finally(() => {
        refreshPromise = null
      })
  }

  return refreshPromise
}

const retryWithNewToken = async (response: AxiosResponse<any>) => {
  const originalConfig = response.config as RetryRequestConfig
  if (originalConfig._retry || originalConfig.url?.includes('/base/token/refresh')) {
    redirectToLogin()
    return Promise.reject(new Error(normalizeMessage(response.data)))
  }

  originalConfig._retry = true

  try {
    const token = await refreshAccessToken()
    originalConfig.headers = originalConfig.headers ?? {}
    originalConfig.headers.Authorization = `Bearer ${token}`
    return instance(originalConfig)
  } catch (error) {
    redirectToLogin()
    return Promise.reject(error)
  }
}

instance.interceptors.request.use(
  (config) => {
    const token = useAuthStore.getState().token
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }

    return config
  },
  (error) => Promise.reject(error),
)

instance.interceptors.response.use(
  (response: AxiosResponse<any>) => {
    if (response.config.responseType === 'blob') {
      return response
    }

    const payload = response.data
    const code = payload?.code

    if (SUCCESS_CODES.includes(code)) {
      if (payload.message === undefined && payload.msg !== undefined) {
        payload.message = payload.msg
      }
      return response
    }

    if (code === ACCESS_TOKEN_EXPIRED_CODE) {
      return retryWithNewToken(response)
    }

    if (code === SYSTEM_NOT_INSTALLED_CODE) {
      redirectToInstall()
      return Promise.reject(new Error(normalizeMessage(payload)))
    }

    if (LOGIN_REQUIRED_CODES.includes(code) || response.status === 401) {
      redirectToLogin()
    }

    const errorMessage = normalizeMessage(payload)
    Message.error(errorMessage)
    return Promise.reject(new Error(errorMessage))
  },
  (error) => {
    const code = error.response?.data?.code

    if (code === ACCESS_TOKEN_EXPIRED_CODE && error.response) {
      return retryWithNewToken(error.response)
    }

    if (code === SYSTEM_NOT_INSTALLED_CODE) {
      redirectToInstall()
      return Promise.reject(error)
    }

    if (error.response?.status === 401 || LOGIN_REQUIRED_CODES.includes(code)) {
      redirectToLogin()
    }

    Message.error(error.response?.data ? normalizeMessage(error.response.data) : error.message || '网络异常')
    return Promise.reject(error)
  },
)

const request = {
  get<T = any>(url: string, config?: AxiosRequestConfig): Promise<Result<T>> {
    return instance.get(url, config).then((res) => res.data)
  },

  post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<Result<T>> {
    return instance.post(url, data, config).then((res) => res.data)
  },

  put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<Result<T>> {
    return instance.put(url, data, config).then((res) => res.data)
  },

  delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<Result<T>> {
    return instance.delete(url, config).then((res) => res.data)
  },

  uploadFile<T = any>(url: string, formData: FormData): Promise<Result<T>> {
    return instance
      .post(url, formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
      .then((res) => res.data)
  },

  downloadFile(url: string): Promise<Blob> {
    return instance.get(url, { responseType: 'blob' }).then((res) => res.data)
  },
}

export default request
