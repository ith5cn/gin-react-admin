import { useAuthStore } from '@/store/auth'
import type { StaticRouteConfig } from '@/types/router'
import { createRouteElement } from './transformRoutes'

const dashboardMapper: Record<string, string> = {
  dashboard: 'system/home/index',
  home: 'system/home/index',
  statistics: 'statistics/index',
  work: 'work/index',
}

const DynamicDashboard = () => {
  const userInfo = useAuthStore((state) => state.userInfo)
  const rawDashboard = userInfo?.dashboard || ''
  const dashboardComponent = dashboardMapper[rawDashboard] || rawDashboard || 'system/home/index'

  return createRouteElement('/dashboard', dashboardComponent)
}

export const staticRoutes: StaticRouteConfig[] = [
  {
    path: '/dashboard',
    name: '工作台',
    isLayout: true,
    meta: {
      title: '工作台',
      icon: 'HomeOutlined',
      affix: true,
      keepAlive: true,
    },
    element: <DynamicDashboard />,
  },
  {
    // 个人中心从头像下拉进入，不在侧边菜单展示。
    path: '/profile',
    name: '个人中心',
    isLayout: true,
    meta: {
      title: '个人中心',
      icon: 'UserOutlined',
      hidden: true,
    },
    element: createRouteElement('/profile', 'profile/index'),
  },
]
