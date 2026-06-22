import request from '@/utils/request'

export interface InstallMysqlConfig {
  host: string
  port: string
  user: string
  password?: string
  dbname: string
}

export interface InstallRedisConfig {
  mode: string
  addr: string
  addrs?: string
  password?: string
  db: number
}

export interface InstallPayload {
  mysql: InstallMysqlConfig
  redis: InstallRedisConfig
  sqlFiles: string[]
  jwtSecret: string
}

export const installStatusApi = () => request.get<{ installed: boolean; sqlFiles: string[] }>('/install/status')

export const installCheckApi = (data: Pick<InstallPayload, 'mysql' | 'redis'>) =>
  request.post<{ mysqlOk: boolean; redisOk: boolean }>('/install/check', data)

export const installRunApi = (data: InstallPayload) =>
  request.post<{ installed: boolean; sqlFiles: string[] }>('/install/run', data)
