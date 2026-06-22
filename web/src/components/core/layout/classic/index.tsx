import { Layout } from 'antd'
import { useState } from 'react'
import { WorkerArea } from '../worker-area'
import { ClassicHeader } from './classic-header'
import { ClassicSider } from './classic-sider'

export const ClassicLayout = () => {
  const [collapsed, setCollapsed] = useState(false)

  return (
    <Layout className="flex h-full min-h-0 justify-between">
      <ClassicSider collapsed={collapsed} />
      <Layout className="min-h-0 flex-1 overflow-hidden bg-[#f4f7fb]">
        <ClassicHeader collapsed={collapsed} setCollapsed={setCollapsed} />
        <div className="min-h-0 flex-1 overflow-hidden p-4">
          <div className="h-full overflow-hidden bg-white p-4">
            <WorkerArea />
          </div>
        </div>
      </Layout>
    </Layout>
  )
}
