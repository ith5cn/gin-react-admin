import { CloseOutlined, FullscreenOutlined, HomeOutlined, LeftOutlined, ReloadOutlined, RightOutlined } from '@ant-design/icons'
import { message } from 'antd'
import { useEffect, useMemo, useRef, useState } from 'react'
import type { CSSProperties } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { HOME_TAG_KEY } from '@/routers/menuHelpers'
import { useAuthStore } from '@/store/auth'

type TagContextAction =
  | 'reload'
  | 'close-current'
  | 'fullscreen-layout'
  | 'fullscreen-content'
  | 'close-left'
  | 'close-right'
  | 'close-others'
  | 'close-all'

type ContextMenuState = {
  visible: boolean
  x: number
  y: number
  tagKey: string | null
}

const SCROLL_OFFSET = 160

const activeTagStyle: CSSProperties = {
  background: 'var(--ith5-primary-bg)',
  boxShadow: 'inset 0 -1px 0 var(--ith5-primary-shadow)',
  color: 'var(--ith5-primary-color)',
}

const activeCloseStyle: CSSProperties = {
  color: 'var(--ith5-primary-color)',
}

export const Tags = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const visitedTags = useAuthStore((state) => state.visitedTags)
  const activeTagKey = useAuthStore((state) => state.activeTagKey)
  const syncTagByPath = useAuthStore((state) => state.syncTagByPath)
  const activateTag = useAuthStore((state) => state.activateTag)
  const closeTag = useAuthStore((state) => state.closeTag)
  const closeLeftTags = useAuthStore((state) => state.closeLeftTags)
  const closeRightTags = useAuthStore((state) => state.closeRightTags)
  const closeOtherTags = useAuthStore((state) => state.closeOtherTags)
  const closeAllTags = useAuthStore((state) => state.closeAllTags)

  const [contextMenu, setContextMenu] = useState<ContextMenuState>({
    visible: false,
    x: 0,
    y: 0,
    tagKey: null,
  })
  const [canScrollLeft, setCanScrollLeft] = useState(false)
  const [canScrollRight, setCanScrollRight] = useState(false)
  const [reloadTick, setReloadTick] = useState(0)
  const scrollerRef = useRef<HTMLDivElement | null>(null)
  const contextMenuRef = useRef<HTMLDivElement | null>(null)

  useEffect(() => {
    syncTagByPath(location.pathname)
  }, [location.pathname, syncTagByPath])

  const updateScrollState = () => {
    const scroller = scrollerRef.current

    if (!scroller) {
      setCanScrollLeft(false)
      setCanScrollRight(false)
      return
    }

    const { scrollLeft, clientWidth, scrollWidth } = scroller
    setCanScrollLeft(scrollLeft > 0)
    setCanScrollRight(scrollLeft + clientWidth < scrollWidth - 1)
  }

  useEffect(() => {
    updateScrollState()
  }, [visitedTags, reloadTick, activeTagKey])

  useEffect(() => {
    const scroller = scrollerRef.current

    if (!scroller) {
      return
    }

    const handleScroll = () => updateScrollState()
    scroller.addEventListener('scroll', handleScroll)
    window.addEventListener('resize', handleScroll)

    return () => {
      scroller.removeEventListener('scroll', handleScroll)
      window.removeEventListener('resize', handleScroll)
    }
  }, [])

  useEffect(() => {
    if (!contextMenu.visible) {
      return
    }

    const handlePointerDown = (event: MouseEvent | Event) => {
      const target = event.target as Node | null

      if (target && contextMenuRef.current?.contains(target)) {
        return
      }

      setContextMenu((prev) => ({ ...prev, visible: false }))
    }

    window.addEventListener('mousedown', handlePointerDown)
    window.addEventListener('scroll', handlePointerDown, true)

    return () => {
      window.removeEventListener('mousedown', handlePointerDown)
      window.removeEventListener('scroll', handlePointerDown, true)
    }
  }, [contextMenu.visible])

  const currentContextTag = useMemo(
    () => visitedTags.find((item) => item.key === contextMenu.tagKey) ?? visitedTags[0],
    [contextMenu.tagKey, visitedTags],
  )

  const runCloseAction = (nextPath: string | null) => {
    if (nextPath) {
      navigate(nextPath)
      return
    }

    activateTag(location.pathname === '/dashboard' ? HOME_TAG_KEY : location.pathname)
  }

  const handleContextAction = (action: TagContextAction, tagKey: string) => {
    switch (action) {
      case 'reload':
        setReloadTick((value) => value + 1)
        message.info(`重新加载：${visitedTags.find((item) => item.key === tagKey)?.label ?? ''}`)
        break
      case 'close-current':
        runCloseAction(closeTag(tagKey))
        break
      case 'close-left':
        runCloseAction(closeLeftTags(tagKey))
        break
      case 'close-right':
        runCloseAction(closeRightTags(tagKey))
        break
      case 'close-others':
        runCloseAction(closeOtherTags(tagKey))
        break
      case 'close-all':
        runCloseAction(closeAllTags())
        break
      case 'fullscreen-layout':
        message.info('全屏主体区域待实现')
        break
      case 'fullscreen-content':
        message.info('全屏内容区域待实现')
        break
      default:
        break
    }

    setContextMenu((prev) => ({ ...prev, visible: false }))
  }

  const scrollTags = (direction: 'left' | 'right') => {
    scrollerRef.current?.scrollBy({
      left: direction === 'left' ? -SCROLL_OFFSET : SCROLL_OFFSET,
      behavior: 'smooth',
    })
  }

  const menuDisabledState = useMemo(() => {
    const targetIndex = visitedTags.findIndex((item) => item.key === currentContextTag.key)
    const hasClosableLeft = visitedTags.slice(1, targetIndex).some((item) => item.closable)
    const hasRight = visitedTags.slice(targetIndex + 1).some((item) => item.closable)
    const hasOthers = visitedTags.some((item) => item.key !== HOME_TAG_KEY && item.key !== currentContextTag.key)
    const hasClosableTags = visitedTags.some((item) => item.closable)

    return {
      closeCurrent: !currentContextTag.closable,
      closeLeft: !hasClosableLeft,
      closeRight: !hasRight,
      closeOthers: !hasOthers,
      closeAll: !hasClosableTags,
    }
  }, [currentContextTag, visitedTags])

  const menuItems: Array<
    | { type: 'separator'; key: string }
    | { type: 'item'; key: TagContextAction; label: string; disabled?: boolean }
  > = [
    { type: 'item', key: 'reload', label: '重新加载' },
    { type: 'item', key: 'close-current', label: '关闭标签页', disabled: menuDisabledState.closeCurrent },
    { type: 'separator', key: 'separator-1' },
    { type: 'item', key: 'fullscreen-layout', label: '全屏主体区域' },
    { type: 'item', key: 'fullscreen-content', label: '全屏内容区域' },
    { type: 'separator', key: 'separator-2' },
    { type: 'item', key: 'close-left', label: '关闭左侧标签页', disabled: menuDisabledState.closeLeft },
    { type: 'item', key: 'close-right', label: '关闭右侧标签页', disabled: menuDisabledState.closeRight },
    { type: 'item', key: 'close-others', label: '关闭其它标签页', disabled: menuDisabledState.closeOthers },
    { type: 'item', key: 'close-all', label: '关闭全部标签页', disabled: menuDisabledState.closeAll },
  ]

  return (
    <div className="relative border-t border-[#eef0f4] bg-white px-2 py-2">
      <div className="flex items-center gap-2">
        <button
          type="button"
          className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg border border-[#edf1f7] bg-white text-[#64748b] transition hover:bg-[#f7f8fb] disabled:cursor-not-allowed disabled:opacity-40"
          onClick={() => scrollTags('left')}
          disabled={!canScrollLeft}
        >
          <LeftOutlined />
        </button>

        <div className="flex-1 overflow-hidden">
          <div ref={scrollerRef} className="flex items-end gap-1 overflow-x-auto scroll-smooth scrollbar-none">
            {visitedTags.map((item) => {
              const isActive = item.key === activeTagKey

              return (
                <div
                  key={item.key}
                  onContextMenu={(event) => {
                    event.preventDefault()
                    setContextMenu({
                      visible: true,
                      x: event.clientX,
                      y: event.clientY,
                      tagKey: item.key,
                    })
                  }}
                  className={[
                    'group relative flex  py-[5px] min-w-[100px] max-w-[180px] shrink-0 items-center gap-2 rounded-[5px] px-3 text-sm transition-all duration-150',
                    isActive
                      ? ''
                      : 'bg-white text-[#1f1f1f] hover:bg-[#f7f8fb]',
                  ].join(' ')}
                  style={isActive ? activeTagStyle : undefined}
                >
                  <button
                    type="button"
                    className="flex min-w-0 flex-1 items-center gap-2 text-left"
                    onClick={() => {
                      activateTag(item.key)
                      navigate(item.path)
                    }}
                    title={item.label}
                  >
                    {item.icon === 'home' ? <HomeOutlined className="text-sm" /> : null}
                    <span className="truncate">{item.label}</span>
                  </button>

                  {item.closable ? (
                    <button
                      type="button"
                      aria-label={`关闭${item.label}`}
                      className={[
                        'flex h-5 w-5 items-center justify-center rounded-full text-[10px] transition',
                        isActive ? 'hover:bg-[var(--ith5-primary-hover-bg)]' : 'text-[#666] hover:bg-black/5',
                      ].join(' ')}
                      style={isActive ? activeCloseStyle : undefined}
                      onClick={(event) => {
                        event.stopPropagation()
                        runCloseAction(closeTag(item.key))
                      }}
                    >
                      <CloseOutlined />
                    </button>
                  ) : null}
                </div>
              )
            })}
          </div>
        </div>

        <button
          type="button"
          className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg border border-[#edf1f7] bg-white text-[#64748b] transition hover:bg-[#f7f8fb] disabled:cursor-not-allowed disabled:opacity-40"
          onClick={() => scrollTags('right')}
          disabled={!canScrollRight}
        >
          <RightOutlined />
        </button>
      </div>

      {contextMenu.visible && currentContextTag ? (
        <div
          ref={contextMenuRef}
          className="fixed z-50 min-w-[188px] rounded-xl border border-[#e8edf7] bg-white py-2 shadow-[0_12px_30px_rgba(15,23,42,0.12)]"
          style={{ left: contextMenu.x, top: contextMenu.y }}
        >
          {menuItems.map((menuItem) => {
            if (menuItem.type === 'separator') {
              return <div key={menuItem.key} className="my-1 border-t border-[#eef0f4]" />
            }

            return (
              <button
                key={menuItem.key}
                type="button"
                disabled={menuItem.disabled}
                className="flex w-full items-center justify-between px-3 py-2 text-left text-sm text-[#1f2937] transition hover:bg-[#f7f8fb] disabled:cursor-not-allowed disabled:text-[#a8b0bf] disabled:hover:bg-white"
                onClick={() => handleContextAction(menuItem.key, currentContextTag.key)}
              >
                <span>{menuItem.label}</span>
                {menuItem.key === 'reload' ? <ReloadOutlined className="text-xs text-[#94a3b8]" /> : null}
                {menuItem.key === 'fullscreen-layout' || menuItem.key === 'fullscreen-content' ? (
                  <FullscreenOutlined className="text-xs text-[#94a3b8]" />
                ) : null}
              </button>
            )
          })}
        </div>
      ) : null}
    </div>
  )
}
