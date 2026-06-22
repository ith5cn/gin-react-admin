import useAppStore from '@/store/useAppStore'
import type { AppMenuItem } from '@/types/router'
import { Operation } from '../operation'
import { TopMenu } from './top-menu'

type MixedHeaderProps = {
  topMenus: AppMenuItem[]
  activeTopKey: string
  onTopMenuChange: (menuItem: AppMenuItem) => void
}

export const MixedHeader = ({ topMenus, activeTopKey, onTopMenuChange }: MixedHeaderProps) => {
  const siteConfig = useAppStore((state) => state.siteConfig)
  const title = siteConfig?.site_name || siteConfig?.name || 'Nest Admin'

  return (
    <div className="shrink-0 border-b border-slate-200 bg-white">
      <div className="flex h-[56px] items-center justify-between gap-6 px-5">
        <div className="flex min-w-0 flex-1 items-center gap-8 overflow-hidden">
          <div className="flex shrink-0 items-center gap-3">
            <div className="text-sm font-semibold text-slate-900">{title}</div>
          </div>

          <TopMenu items={topMenus} activeKey={activeTopKey} onChange={onTopMenuChange} />
        </div>

        <Operation />
      </div>
    </div>
  )
}
