import { Button, Result } from 'antd'
import { useNavigate } from 'react-router-dom'
import { HOME_PATH } from './menuHelpers'

const RouteFallback = () => {
  const navigate = useNavigate()

  return (
    <Result
      status="404"
      title="页面不存在"
      subTitle="当前地址没有匹配到可访问的页面。"
      extra={
        <Button type="primary" onClick={() => navigate(HOME_PATH, { replace: true })}>
          返回首页
        </Button>
      }
    />
  )
}

export default RouteFallback
