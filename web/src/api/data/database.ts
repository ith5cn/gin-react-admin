import request from '@/utils/request'

export interface DatabaseTableRecord {
  tableName: string
  tableComment: string
  engine: string
  updateTime?: string
  tableRows: number
  fragmentSize: number
  dataSize: number
  indexSize: number
  tableCharset: string
  createTime?: string
}

export interface DatabaseColumnRecord {
  columnName: string
  columnComment: string
  columnType: string
  isNullable: string
  columnKey: string
  columnDefault?: string
  extra: string
  ordinal: number
}

export interface RecycleRecord {
  id: string | number
  deleteTime?: string
  content: Record<string, any>
}

export const databaseTableListApi = (params?: any) => request.get('/data/database/index', { params })

export const databaseTableColumnsApi = (tableName: string) =>
  request.get<DatabaseColumnRecord[]>(`/data/database/columns/${tableName}`)

export const databaseClearFragmentApi = (tables: string[]) =>
  request.post('/data/database/fragment', { tables })

export const databaseOptimizeApi = (tables: string[]) =>
  request.post('/data/database/optimize', { tables })

export const databaseRecycleListApi = (params?: any) =>
  request.get<{ list: RecycleRecord[]; total: number }>('/data/database/recycle', { params })

export const databaseRecycleRecoverApi = (tableName: string, ids: Array<string | number>) =>
  request.post('/data/database/recycle/recover', { tableName, ids })

export const databaseRecycleDestroyApi = (tableName: string, ids: Array<string | number>) =>
  request.post('/data/database/recycle/destroy', { tableName, ids })
