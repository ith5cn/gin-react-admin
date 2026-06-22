import type { AppMenuItem } from '@/types/router'

type TopMenuProps = {
  items: AppMenuItem[]
  activeKey: string
  onChange: (menuItem: AppMenuItem) => void
}

export const TopMenu = ({ items, activeKey, onChange }: TopMenuProps) => {
  return (
    <div className="flex min-w-0 items-center gap-1 overflow-x-auto">
      {items.map((item) => {
        const isActive = item.key === activeKey

        return (
          <button
            key={item.key}
            type="button"
            onClick={() => onChange(item)}
            style={isActive ? { color: 'var(--ith5-primary-color)' } : undefined}
            className={[
              'relative flex h-10 shrink-0 items-center gap-2 rounded-[5px] px-4 text-sm font-medium transition',
              isActive ? '' : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900',
            ].join(' ')}
          >
            {item.icon ? <span className="text-[15px]">{item.icon}</span> : null}
            <span>{item.title}</span>
          </button>
        )
      })}
    </div>
  )
}
