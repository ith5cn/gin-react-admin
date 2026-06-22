import type { AppMenuItem, TagViewItem } from '@/types/router'

export const HOME_TAG_KEY = 'dashboard'
export const HOME_PATH = '/dashboard'

export const createHomeTag = (): TagViewItem => ({
  key: HOME_TAG_KEY,
  label: '首页',
  path: HOME_PATH,
  closable: false,
  icon: 'home',
})

export const findMenuPathChain = (items: AppMenuItem[], pathname: string): string[] => {
  for (const item of items) {
    if (item.path === pathname) {
      return [item.key]
    }

    if (item.children?.length) {
      const childChain = findMenuPathChain(item.children, pathname)

      if (childChain.length) {
        return [item.key, ...childChain]
      }
    }
  }

  return []
}

export const findMenuKeyChain = (items: AppMenuItem[], targetKey: string): string[] => {
  for (const item of items) {
    if (item.key === targetKey) {
      return [item.key]
    }

    if (item.children?.length) {
      const childChain = findMenuKeyChain(item.children, targetKey)

      if (childChain.length) {
        return [item.key, ...childChain]
      }
    }
  }

  return []
}

export const findMenuItemByPath = (items: AppMenuItem[], pathname: string): AppMenuItem | null => {
  for (const item of items) {
    if (item.path === pathname) {
      return item
    }

    if (item.children?.length) {
      const childItem = findMenuItemByPath(item.children, pathname)

      if (childItem) {
        return childItem
      }
    }
  }

  return null
}

export const buildTagFromPath = (pathname: string, sideMenuItems: AppMenuItem[]): TagViewItem | null => {
  if (pathname === HOME_PATH) {
    return createHomeTag()
  }

  const menuItem = findMenuItemByPath(sideMenuItems, pathname)

  if (!menuItem || menuItem.external) {
    return null
  }

  return {
    key: menuItem.path,
    label: menuItem.title,
    path: menuItem.path,
    closable: true,
  }
}
