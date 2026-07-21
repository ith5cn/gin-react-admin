import { useCallback, useEffect, useMemo, useState, type ReactNode } from 'react'
import { Button, Result, Spin } from 'antd'
import { Navigate, useLocation } from 'react-router-dom'
import { installStatusApi } from '@/api/install'
import { InstallBootstrapContext, type InstallBootstrapContextValue } from './installBootstrapContext'

type InstallState = 'loading' | 'installed' | 'not-installed' | 'error'

let statusRequest: Promise<boolean> | null = null

const loadInstallStatus = () => {
  if (!statusRequest) {
    statusRequest = installStatusApi()
      .then((response) => response.data.installed)
      .catch((error) => {
        statusRequest = null
        throw error
      })
  }

  return statusRequest
}

interface InstallBootstrapProps {
  children: ReactNode
}

const InstallBootstrap = ({ children }: InstallBootstrapProps) => {
  const location = useLocation()
  const [state, setState] = useState<InstallState>('loading')

  useEffect(() => {
    let cancelled = false

    void loadInstallStatus()
      .then((installed) => {
        if (!cancelled) {
          setState(installed ? 'installed' : 'not-installed')
        }
      })
      .catch(() => {
        if (!cancelled) {
          setState('error')
        }
      })

    return () => {
      cancelled = true
    }
  }, [])

  const retryCheck = useCallback(() => {
    setState('loading')
    void loadInstallStatus()
      .then((installed) => setState(installed ? 'installed' : 'not-installed'))
      .catch(() => setState('error'))
  }, [])

  const contextValue = useMemo<InstallBootstrapContextValue>(
    () => ({
      markInstalled: () => setState('installed'),
    }),
    [],
  )

  if (state === 'loading') {
    return (
      <div className="grid min-h-screen place-items-center bg-white">
        <Spin size="large" tip="正在检查安装状态..." />
      </div>
    )
  }

  if (state === 'error') {
    return (
      <div className="grid min-h-screen place-items-center bg-[#f6f7f9] px-4">
        <Result
          status="error"
          title="无法获取安装状态"
          subTitle="请确认后端服务已经启动并可正常访问。"
          extra={<Button type="primary" onClick={retryCheck}>重新检查</Button>}
        />
      </div>
    )
  }

  if (state === 'not-installed' && location.pathname !== '/install') {
    return <Navigate to="/install" replace />
  }

  return <InstallBootstrapContext.Provider value={contextValue}>{children}</InstallBootstrapContext.Provider>
}

export default InstallBootstrap