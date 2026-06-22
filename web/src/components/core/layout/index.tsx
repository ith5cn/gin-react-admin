import { useEffect } from 'react'
import { Layout as AntdLayout, Spin } from 'antd'
import { Navigate } from 'react-router-dom'
import { useAuthStore } from '@/store/auth'
import useAppStore from '@/store/useAppStore'
import { ClassicLayout } from './classic'
import { MixedLayout } from './mixed'

export const Layout = () => {
  const token = useAuthStore((state) => state.token)
  const isInitializing = useAuthStore((state) => state.isInitializing)
  const isUserInitialized = useAuthStore((state) => state.isUserInitialized)
  const initUserContext = useAuthStore((state) => state.initUserContext)
  const initSiteConfig = useAppStore((state) => state.initSiteConfig)
  const layoutMode = useAppStore((state) => state.layoutMode)

  useEffect(() => {
    if (token && !isUserInitialized && !isInitializing) {
      void initUserContext()
    }
  }, [initUserContext, isInitializing, isUserInitialized, token])

  useEffect(() => {
    if (token && isUserInitialized) {
      void initSiteConfig()
    }
  }, [initSiteConfig, isUserInitialized, token])

  if (!token) {
    return <Navigate to="/login" replace />
  }

  if (isInitializing || !isUserInitialized) {
    return (
      <div className="grid min-h-screen place-items-center bg-white">
        <Spin size="large" tip="正在加载用户上下文..." />
      </div>
    )
  }

  return (
    <AntdLayout className="main-container h-screen overflow-hidden">
      {layoutMode === 'classic' ? <ClassicLayout /> : <MixedLayout />}
    </AntdLayout>
  )
}
