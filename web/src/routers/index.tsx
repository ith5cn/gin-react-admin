import { useMemo } from 'react'
import { useRoutes } from 'react-router-dom'
import { useAuthStore } from '@/store/auth'
import { baseRouters } from './baseRouters'
import { fullScreenRoutes, getLayoutRouters } from './layoutRouters'
import { publicRoutes } from './publicRouters'

const RenderRoutes = () => {
  const dynamicRoutes = useAuthStore((state) => state.dynamicRoutes)

  const routes = useMemo(
    () => [...publicRoutes, ...getLayoutRouters(dynamicRoutes), ...fullScreenRoutes, ...baseRouters],
    [dynamicRoutes],
  )

  return useRoutes(routes)
}

export { publicRoutes }
export default RenderRoutes
