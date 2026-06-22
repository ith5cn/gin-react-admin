import { MenuFoldOutlined, MenuUnfoldOutlined } from '@ant-design/icons'
import { Button } from 'antd'
import useAppStore from '@/store/useAppStore'
import { Operation } from '../operation'
import { Tags } from '../tags'

type ClassicHeaderProps = {
  collapsed: boolean
  setCollapsed: (collapsed: boolean) => void
}

export const ClassicHeader = ({ collapsed, setCollapsed }: ClassicHeaderProps) => {
  const siteConfig = useAppStore((state) => state.siteConfig)
  const title = siteConfig?.site_name || siteConfig?.name || 'Nest Admin'

  return (
    <div className="shrink-0 border-b border-slate-200 bg-white">
      <div className="flex h-[48px] items-center justify-between px-2">
        <div className="flex min-w-0 items-center gap-2">
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
          />
          <span className="truncate text-sm font-semibold text-slate-900">{title}</span>
        </div>
        <Operation />
      </div>
      <Tags />
    </div>
  )
}
