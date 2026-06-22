import type { NormalizedMenuNode, RawBackendMenuNode } from '@/types/router'

const normalizeRoutePath = (route?: string | null, isExternal = false) => {
  if (!route) {
    return ''
  }

  if (isExternal || /^https?:\/\//i.test(route)) {
    return route
  }

  return route.startsWith('/') ? route : `/${route}`
}

export const normalizeBackendMenuNode = (menuNode: RawBackendMenuNode): NormalizedMenuNode => {
  const isExternal = menuNode.type === 'L'

  return {
    id: menuNode.id,
    name: menuNode.name,
    code: menuNode.code,
    path: normalizeRoutePath(menuNode.route, isExternal),
    component: menuNode.component ?? '',
    redirect: menuNode.redirect ?? null,
    meta: {
      title: menuNode.name,
      type: menuNode.type,
      hidden: menuNode.isHidden === 1,
      layout: menuNode.isLayout === 1,
      hiddenBreadcrumb: false,
      icon: menuNode.icon ?? undefined,
      external: isExternal,
    },
    children: menuNode.children?.filter((child) => child.status !== 2).map(normalizeBackendMenuNode) ?? [],
  }
}

export const normalizeBackendMenuTree = (menuTree: RawBackendMenuNode[]) => {
  return menuTree.filter((menuNode) => menuNode.status !== 2).map(normalizeBackendMenuNode)
}
