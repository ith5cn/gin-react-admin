import { useMemo, useState } from 'react'
import { Menu } from 'antd'
import type { MenuProps } from 'antd'
import Sider from 'antd/es/layout/Sider'
import { useLocation, useNavigate } from 'react-router-dom'
import { findMenuKeyChain, findMenuPathChain } from '@/routers/menuHelpers'
import { useAuthStore } from '@/store/auth'
import type { AppMenuItem } from '@/types/router'

const mapMenuItems = (items: AppMenuItem[]): MenuProps['items'] => {
  return items.map((item) => ({
    key: item.key,
    icon: item.icon,
    label: item.title,
    children: item.children?.length ? mapMenuItems(item.children) : undefined,
  }))
}

type ClassicSiderMenuProps = {
  collapsed: boolean
  pathname: string
  sideMenuItems: AppMenuItem[]
}

const ClassicSiderMenu = ({ collapsed, pathname, sideMenuItems }: ClassicSiderMenuProps) => {
  const navigate = useNavigate()
  const selectedKeyChain = useMemo(
    () => findMenuPathChain(sideMenuItems, pathname),
    [pathname, sideMenuItems],
  )
  const selectedKeys = selectedKeyChain.length ? [selectedKeyChain[selectedKeyChain.length - 1]] : []
  const routeOpenKeys = selectedKeyChain.slice(0, -1)
  const [openKeys, setOpenKeys] = useState<string[]>(routeOpenKeys)
  const menuItems = useMemo(() => mapMenuItems(sideMenuItems), [sideMenuItems])

  const handleOpenChange: MenuProps['onOpenChange'] = (nextOpenKeys) => {
    const latestOpenedKey = nextOpenKeys.find((key) => !openKeys.includes(String(key)))

    if (!latestOpenedKey) {
      setOpenKeys(nextOpenKeys.map(String))
      return
    }

    const latestKeyChain = findMenuKeyChain(sideMenuItems, String(latestOpenedKey))
    setOpenKeys(latestKeyChain)
  }

  const handleClick: MenuProps['onClick'] = ({ key }) => {
    const selectedChain = findMenuKeyChain(sideMenuItems, String(key))
    const selectedItemKey = String(key)
    const targetItem = selectedChain.reduce<AppMenuItem | undefined>((current, chainKey, index) => {
      const source = index === 0 ? sideMenuItems : current?.children ?? []
      return source.find((item) => item.key === chainKey)
    }, undefined)

    setOpenKeys(selectedChain.slice(0, -1))

    if (targetItem?.external || /^https?:\/\//i.test(selectedItemKey)) {
      window.open(targetItem?.path || selectedItemKey, '_blank')
      return
    }

    navigate(targetItem?.path || selectedItemKey)
  }

  return (
    <Menu
      theme="dark"
      mode="inline"
      className="h-[calc(100%-64px)] overflow-y-auto overflow-x-hidden border-e-0"
      selectedKeys={selectedKeys}
      openKeys={collapsed ? [] : openKeys}
      onOpenChange={handleOpenChange}
      onClick={handleClick}
      items={menuItems}
    />
  )
}

export const ClassicSider = ({ collapsed }: { collapsed: boolean }) => {
  const location = useLocation()
  const sideMenuItems = useAuthStore((state) => state.sideMenuItems)

  return (
    <Sider trigger={null} collapsible collapsed={collapsed} className="h-full overflow-hidden">
      <div className="m-4 flex h-8 items-center justify-center rounded-md bg-white/10 text-xs font-semibold text-white">
        {!collapsed ? 'Ith5 Admin' : 'I'}
      </div>
      <ClassicSiderMenu
        key={`${collapsed ? 'collapsed' : 'expanded'}-${location.pathname}`}
        collapsed={collapsed}
        pathname={location.pathname}
        sideMenuItems={sideMenuItems}
      />
    </Sider>
  )
}
