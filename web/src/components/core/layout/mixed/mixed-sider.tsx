import { Menu } from 'antd'
import Sider from 'antd/es/layout/Sider'
import type { MenuProps } from 'antd'
import { useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { findMenuKeyChain, findMenuPathChain } from '@/routers/menuHelpers'
import type { AppMenuItem } from '@/types/router'

const mapMenuItems = (items: AppMenuItem[]): MenuProps['items'] => {
  return items.map((item) => ({
    key: item.key,
    icon: item.icon,
    label: item.title,
    children: item.children?.length ? mapMenuItems(item.children) : undefined,
  }))
}

type MixedSiderProps = {
  menuItems: AppMenuItem[]
  pathname: string
}

export const MixedSider = ({ menuItems, pathname }: MixedSiderProps) => {
  const navigate = useNavigate()

  const selectedKeyChain = useMemo(() => findMenuPathChain(menuItems, pathname), [menuItems, pathname])
  const selectedKeys = selectedKeyChain.length ? [selectedKeyChain[selectedKeyChain.length - 1]] : []
  const routeOpenKeys = selectedKeyChain.slice(0, -1)
  const [openKeys, setOpenKeys] = useState<string[]>(routeOpenKeys)
  const mappedItems = useMemo(() => mapMenuItems(menuItems), [menuItems])

  const handleOpenChange: MenuProps['onOpenChange'] = (nextOpenKeys) => {
    const latestOpenedKey = nextOpenKeys.find((key) => !openKeys.includes(String(key)))

    if (!latestOpenedKey) {
      setOpenKeys(nextOpenKeys.map(String))
      return
    }

    const latestKeyChain = findMenuKeyChain(menuItems, String(latestOpenedKey))
    setOpenKeys(latestKeyChain)
  }

  const handleClick: MenuProps['onClick'] = ({ key }) => {
    const selectedChain = findMenuKeyChain(menuItems, String(key))
    const selectedItemKey = String(key)
    const targetItem = selectedChain.reduce<AppMenuItem | undefined>((current, chainKey, index) => {
      const source = index === 0 ? menuItems : current?.children ?? []
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
    <Sider width={228} theme="light" trigger={null} className="h-full overflow-hidden border-r border-slate-200 !bg-white">
      <Menu
        mode="inline"
        className="h-full overflow-y-auto overflow-x-hidden border-e-0 bg-transparent px-2 pt-2"
        selectedKeys={selectedKeys}
        openKeys={openKeys}
        onOpenChange={handleOpenChange}
        onClick={handleClick}
        items={mappedItems}
      />
    </Sider>
  )
}
