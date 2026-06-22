import type { RouteObject } from 'react-router-dom'
import RouteFallback from './RouteFallback'

export const baseRouters: RouteObject[] = [{ path: '*', element: <RouteFallback /> }]
