import { lazy, Suspense } from 'react'
import type { ComponentType, ReactNode } from 'react'
import * as AntdIcons from '@ant-design/icons'
import * as LucideIcons from 'lucide-react'
import { AppstoreOutlined } from '@ant-design/icons'
import RoutePlaceholder from '@/pages/common/route-placeholder'
import type { AppMenuItem, AppRouteObject, NormalizedMenuNode } from '@/types/router'
import { HOME_PATH } from './menuHelpers'

const pageModules = import.meta.glob('../pages/**/*.{tsx,ts}')
const FALLBACK_MENU_ICON = <AppstoreOutlined />

const createMenuItemKey = (menuNode: NormalizedMenuNode) => {
  return menuNode.path || menuNode.code || String(menuNode.id)
}

export const resolveMenuIcon = (icon?: string, size = 16): ReactNode => {
  if (!icon) {
    return FALLBACK_MENU_ICON
  }

  if (icon.startsWith('lucide:')) {
    const name = icon.replace('lucide:', '')
    const IconComponent = (LucideIcons as unknown as Record<string, ComponentType<any>>)[name]
    return IconComponent ? <IconComponent size={size} /> : FALLBACK_MENU_ICON
  }

  const IconComponent = (AntdIcons as unknown as Record<string, ComponentType<any>>)[icon]
  return IconComponent ? <IconComponent style={{ fontSize: size }} /> : FALLBACK_MENU_ICON
}

const normalizeComponentPath = (componentPath: string) => {
  return componentPath.replace(/^\/+/, '').replace(/\.(tsx|ts|jsx|js)$/, '')
}

const resolveComponentCandidates = (componentPath: string) => {
  const normalizedPath = normalizeComponentPath(componentPath)

  return [
    `../pages/${normalizedPath}.tsx`,
    `../pages/${normalizedPath}.ts`,
    `../pages/${normalizedPath}/index.tsx`,
    `../pages/${normalizedPath}/index.ts`,
  ]
}

export const createRouteElement = (routePath: string, componentPath: string) => {
  const matchedModulePath = resolveComponentCandidates(componentPath).find((candidatePath) => candidatePath in pageModules)

  if (!matchedModulePath) {
    return <RoutePlaceholder routePath={routePath} componentPath={componentPath} />
  }

  const loader = pageModules[matchedModulePath] as () => Promise<{ default: ComponentType }>
  const LazyComponent = lazy(loader)

  return (
    <Suspense fallback={<div className="p-6 text-slate-500">页面加载中...</div>}>
      <LazyComponent />
    </Suspense>
  )
}

const createIframeElement = (routePath: string) => (
  <iframe src={routePath} className="h-full w-full border-0" title={routePath} />
)

const flattenRouteNodes = (menuTree: NormalizedMenuNode[]): NormalizedMenuNode[] => {
  return menuTree.flatMap((menuNode) => {
    const isPageRoute =
      menuNode.meta.type === 'M' && menuNode.meta.layout && menuNode.path && menuNode.component
    const isIframeRoute = menuNode.meta.type === 'I' && menuNode.meta.layout && menuNode.path
    const currentNode = isPageRoute || isIframeRoute ? [menuNode] : []

    if (menuNode.children?.length) {
      return [...currentNode, ...flattenRouteNodes(menuNode.children)]
    }

    return currentNode
  })
}

export const transformMenuToRoutes = (menuTree: NormalizedMenuNode[]): AppRouteObject[] => {
  return flattenRouteNodes(menuTree)
    .filter((menuNode) => menuNode.path !== HOME_PATH)
    .map((menuNode) => ({
      id: String(menuNode.id),
      path: menuNode.path,
      element:
        menuNode.meta.type === 'I'
          ? createIframeElement(menuNode.path)
          : createRouteElement(menuNode.path, menuNode.component),
      handle: {
        name: menuNode.name,
        meta: {
          title: menuNode.meta.title,
          icon: menuNode.meta.icon,
          hidden: menuNode.meta.hidden,
          hiddenBreadcrumb: menuNode.meta.hiddenBreadcrumb,
          type: menuNode.meta.type,
          backendComponent: menuNode.component,
        },
      },
    }))
}

const isMenuNodeVisible = (menuNode: NormalizedMenuNode) => menuNode.meta.type !== 'B' && !menuNode.meta.hidden

export const transformMenuToSiderItems = (menuTree: NormalizedMenuNode[]): AppMenuItem[] => {
  return menuTree.flatMap((menuNode) => {
    if (menuNode.meta.type === 'B') {
      return []
    }

    const children = menuNode.children?.length ? transformMenuToSiderItems(menuNode.children) : []
    const hasChildren = children.length > 0
    const isClickablePage = Boolean(menuNode.path)

    if (!isMenuNodeVisible(menuNode)) {
      return []
    }

    if (menuNode.children?.length && hasChildren) {
      return [
        {
          key: createMenuItemKey(menuNode),
          path: menuNode.path,
          title: menuNode.meta.title,
          icon: resolveMenuIcon(menuNode.meta.icon),
          hidden: menuNode.meta.hidden,
          external: menuNode.meta.external,
          children,
        },
      ]
    }

    if (!isClickablePage) {
      return []
    }

    return [
      {
        key: createMenuItemKey(menuNode),
        path: menuNode.path,
        title: menuNode.meta.title,
        icon: resolveMenuIcon(menuNode.meta.icon),
        hidden: menuNode.meta.hidden,
        external: menuNode.meta.external,
      },
    ]
  })
}
