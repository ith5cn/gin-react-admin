import { ClearOutlined, DownOutlined, FullscreenOutlined, LogoutOutlined, SettingOutlined, UserOutlined } from '@ant-design/icons'
import { Avatar, Dropdown, Space } from 'antd'
import type { MenuProps } from 'antd'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { clearAuthSession, useAuthStore } from '@/store/auth'
import { logoutApi } from '@/api/auth'
import { ThemeDrawer } from './theme-drawer'

export const Operation = () => {
  const navigate = useNavigate()
  const userInfo = useAuthStore((state) => state.userInfo)
  const [settingVisible, setSettingVisible] = useState(false)

  const items: MenuProps['items'] = [
    { key: 'clear-cache', label: '清除缓存', icon: <ClearOutlined /> },
    { type: 'divider' },
    { key: 'logout', label: '退出登录', icon: <LogoutOutlined />, danger: true },
  ]

  const handleMenuClick: MenuProps['onClick'] = async ({ key }) => {
    if (key === 'clear-cache') {
      localStorage.clear()
      sessionStorage.clear()
      window.location.reload()
      return
    }

    if (key === 'logout') {
      const refreshToken = useAuthStore.getState().refreshToken
      try {
        await logoutApi(refreshToken)
      } catch {
        // 后端撤销失败时仍清理前端登录态，避免用户卡在当前页面。
      }

      clearAuthSession()
      navigate('/login', { replace: true })
    }
  }

  const handleSettingClick = () => {
    setSettingVisible(true)
  }

  return (
    <>
      <Space className="mr-2 w-full lg:w-auto" size={20}>
        <FullscreenOutlined className="cursor-pointer" />
        <SettingOutlined className="cursor-pointer" onClick={handleSettingClick} />

        <Dropdown menu={{ items, onClick: handleMenuClick }} trigger={['click']}>
          <Space className="cursor-pointer">
            <Avatar icon={<UserOutlined />} />
            <span>{userInfo?.nickname || userInfo?.username || '用户'}</span>
            <DownOutlined />
          </Space>
        </Dropdown>

      </Space>
      <ThemeDrawer visible={settingVisible} onClose={() => setSettingVisible(false)} />
    </>
  )
}
