import { Layout } from 'antd'
import { useMemo } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { findMenuPathChain, HOME_PATH } from '@/routers/menuHelpers'
import { useAuthStore } from '@/store/auth'
import type { AppMenuItem } from '@/types/router'
import { Tags } from '../tags'
import { WorkerArea } from '../worker-area'
import { MixedHeader } from './mixed-header'
import { MixedSider } from './mixed-sider'

const findTopLevelMenuByPath = (items: AppMenuItem[], pathname: string) => {
  const keyChain = findMenuPathChain(items, pathname)
  const topLevelKey = keyChain[0]

  if (!topLevelKey) {
    return items[0] ?? null
  }

  return items.find((item) => item.key === topLevelKey) ?? items[0] ?? null
}

const findFirstNavigablePath = (menuItem: AppMenuItem): string => {
  if (menuItem.path && !menuItem.external) {
    return menuItem.path
  }

  if (!menuItem.children?.length) {
    return ''
  }

  for (const child of menuItem.children) {
    const childPath = findFirstNavigablePath(child)

    if (childPath) {
      return childPath
    }
  }

  return ''
}

export const MixedLayout = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const sideMenuItems = useAuthStore((state) => state.sideMenuItems)

  const activeTopMenu = useMemo(
    () => findTopLevelMenuByPath(sideMenuItems, location.pathname),
    [location.pathname, sideMenuItems],
  )

  const handleTopMenuChange = (menuItem: AppMenuItem) => {
    if (menuItem.external) {
      window.open(menuItem.path, '_blank')
      return
    }

    if (menuItem.children?.length) {
      const firstChildPath = findFirstNavigablePath(menuItem)

      if (firstChildPath && location.pathname !== firstChildPath) {
        navigate(firstChildPath)
      }

      return
    }

    if (menuItem.path && location.pathname !== menuItem.path) {
      navigate(menuItem.path)
    }
  }

  const siderMenuItems = activeTopMenu?.children ?? []
  const shouldShowSider = siderMenuItems.length > 0 && location.pathname !== HOME_PATH

  return (
    <Layout className="flex h-full min-h-0 flex-col bg-[#f4f7fb]">
      <MixedHeader
        topMenus={sideMenuItems}
        activeTopKey={activeTopMenu?.key ?? ''}
        onTopMenuChange={handleTopMenuChange}
      />

      <Layout className="min-h-0 flex-1 bg-transparent">
        {shouldShowSider ? (
          <MixedSider key={`${activeTopMenu?.key ?? 'empty'}-${location.pathname}`} menuItems={siderMenuItems} pathname={location.pathname} />
        ) : null}
        <Layout className="min-h-0 flex-1 overflow-hidden bg-transparent">
          <div className="flex min-h-0 flex-1 flex-col overflow-hidden">
            <div className="shrink-0 border-b border-slate-200 bg-white px-4">
              <Tags />
            </div>
            <div className="min-h-0 flex-1 overflow-hidden p-4">
              <div className="h-full overflow-hidden bg-white p-4">
                <WorkerArea />
              </div>
            </div>
          </div>
        </Layout>
      </Layout>
    </Layout>
  )
}
