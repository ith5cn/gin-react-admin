import Login from '@/pages/login'
import Install from '@/pages/install'
import Register from '@/pages/register'

export const publicRoutes = [
  {
    path: '/install',
    element: <Install />,
  },
  {
    path: '/login',
    element: <Login />,
  },
  {
    path: '/register',
    element: <Register />,
  },
]
