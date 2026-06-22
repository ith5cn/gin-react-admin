import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import enUS from 'antd/locale/en_US'
import { BrowserRouter } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import RenderRoutes from '@/routers'
import useAppStore from '@/store/useAppStore'

function App() {
  const { i18n } = useTranslation()
  const antdLocale = i18n.language === 'en-US' ? enUS : zhCN
  const primaryColor = useAppStore((state) => state.primaryColor)

  return (
    <ConfigProvider
      locale={antdLocale}
      theme={{
        token: {
          colorPrimary: primaryColor,
          colorLink: primaryColor,
          colorLinkActive: primaryColor,
          colorLinkHover: primaryColor,
        },
      }}
    >
      <style>
        {`:root {
          --ith5-primary-color: ${primaryColor};
          --ith5-primary-bg: color-mix(in srgb, ${primaryColor} 10%, white);
          --ith5-primary-hover-bg: color-mix(in srgb, ${primaryColor} 16%, white);
          --ith5-primary-shadow: color-mix(in srgb, ${primaryColor} 16%, transparent);
        }`}
      </style>
      <BrowserRouter>
        <RenderRoutes />
      </BrowserRouter>
    </ConfigProvider>
  )
}

export default App
