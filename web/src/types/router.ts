import type { ReactNode } from 'react'
import type { RouteObject } from 'react-router-dom'

export type MenuType = 'M' | 'B' | 'L' | 'I'

export type RawBackendMenuNode = {
  id: number
  parentId: number | null
  name: string
  code: string
  icon?: string | null
  route?: string | null
  component?: string | null
  redirect?: string | null
  isHidden: number
  isLayout: number
  type: MenuType
  status?: number
  sort?: number
  children?: RawBackendMenuNode[]
}

export type NormalizedMenuMeta = {
  title: string
  type: MenuType
  hidden: boolean
  layout: boolean
  hiddenBreadcrumb: boolean
  icon?: string
  external?: boolean
}

export type NormalizedMenuNode = {
  id: number
  name: string
  code: string
  path: string
  component: string
  redirect?: string | null
  meta: NormalizedMenuMeta
  children?: NormalizedMenuNode[]
}

export type AppRouteMeta = {
  title: string
  icon?: string
  hidden?: boolean
  keepAlive?: boolean
  affix?: boolean
  type?: string
  hiddenBreadcrumb?: boolean
  backendComponent?: string
}

export type AppMenuItem = {
  key: string
  path: string
  title: string
  icon?: ReactNode
  hidden?: boolean
  external?: boolean
  children?: AppMenuItem[]
}

export type TagViewItem = {
  key: string
  label: string
  path: string
  closable: boolean
  icon?: 'home'
}

export type StaticRouteConfig = {
  path: string
  name: string
  isLayout: boolean
  meta: AppRouteMeta
  element: ReactNode
}

export type AppRouteObject = RouteObject & {
  id?: string
  handle?: {
    name?: string
    meta?: AppRouteMeta
  }
}
