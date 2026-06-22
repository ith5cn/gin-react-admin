import { Outlet } from 'react-router-dom'

export const WorkerArea = () => {
  return (
    <div className="worker-area h-full flex-1 overflow-auto">
      <Outlet />
    </div>
  )
}
