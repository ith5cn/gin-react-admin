import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { message } from 'antd'
import { ArrowRightOutlined, GlobalOutlined, LockOutlined, SafetyCertificateOutlined, UserOutlined } from '@ant-design/icons'

import { loginApi } from '@/api/auth'
import { useAuthStore } from '@/store/auth'

type LoginFormValues = {
  username: string
  password: string
}

const LoginForm: React.FC = () => {
  const [loading, setLoading] = useState(false)
  const [values, setValues] = useState<LoginFormValues>({
    username: 'admin',
    password: '123456',
  })
  const navigate = useNavigate()

  const onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setLoading(true)
    try {
      const res = await loginApi(values)
      // 只负责写入 token 并跳转，用户上下文（菜单/权限）由 Layout 统一加载，
      // 避免登录页在此处重复初始化、与 Layout 的初始化相互竞争导致页面闪烁。
      await useAuthStore.getState().login(res)
      message.success('登录成功')
      navigate('/dashboard', { replace: true })
    } catch {
      // 错误提示已由请求拦截器统一弹出，这里不再重复处理
    } finally {
      setLoading(false)
    }
  }

  return (
    <div
      className="w-full max-w-[480px] p-10 rounded-[2rem] relative overflow-hidden shadow-[0px_40px_80px_rgba(0,105,113,0.08)]"
      style={{
        background: 'rgba(255, 255, 255, 0.7)',
        backdropFilter: 'blur(20px)',
        border: '1px solid rgba(187, 201, 203, 0.15)',
      }}
    >
      {/* Subtle interior glow */}
      <div className="absolute -top-24 -right-24 w-48 h-48 bg-[#006971]/10 rounded-full blur-3xl pointer-events-none" />

      <div className="relative z-10">
        <div className="mb-10 text-center">
          <h2 className="text-2xl font-bold text-[#191c1d] mb-2 font-headline">欢迎回来</h2>
          <p className="text-[#3c494b]/60 text-sm">请通过管理凭据验证身份</p>
        </div>

        <form className="space-y-6" onSubmit={onSubmit}>
          <div className="space-y-2">
            <label
              className="block text-xs font-semibold tracking-widest text-[#006971] uppercase ml-1"
              htmlFor="username"
            >
              访问标识 / 用户名
            </label>
            <div className="relative">
              <UserOutlined className="absolute left-4 top-1/2 -translate-y-1/2 text-[#3c494b]/40 text-base pointer-events-none" />
              <input
                id="username"
                className="w-full pl-12 pr-4 py-4 bg-[#e1e3e4]/30 border-none rounded-xl text-[#191c1d] placeholder:text-[#3c494b]/30 outline-none focus:ring-2 focus:ring-[#006971]/20 transition-all"
                placeholder="Access ID or Username"
                autoComplete="username"
                required
                value={values.username}
                onChange={(e) => setValues((prev) => ({ ...prev, username: e.target.value }))}
              />
            </div>
          </div>

          <div className="space-y-2">
            <label
              className="block text-xs font-semibold tracking-widest text-[#006971] uppercase ml-1"
              htmlFor="password"
            >
              安全令牌
            </label>
            <div className="relative">
              <LockOutlined className="absolute left-4 top-1/2 -translate-y-1/2 text-[#3c494b]/40 text-base pointer-events-none" />
              <input
                id="password"
                type="password"
                className="w-full pl-12 pr-4 py-4 bg-[#e1e3e4]/30 border-none rounded-xl text-[#191c1d] placeholder:text-[#3c494b]/30 outline-none focus:ring-2 focus:ring-[#006971]/20 transition-all"
                placeholder="Access Cipher"
                autoComplete="current-password"
                required
                value={values.password}
                onChange={(e) => setValues((prev) => ({ ...prev, password: e.target.value }))}
              />
            </div>
            <div className="flex justify-end">
              <a className="text-xs text-[#3c494b]/60 hover:text-[#006971] transition-colors" href="#">
                忘记访问令牌?
              </a>
            </div>
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full py-4 bg-[#006971] text-white rounded-xl font-semibold shadow-xl shadow-[#006971]/25 hover:bg-[#005a61] hover:shadow-[#006971]/30 transition-all active:scale-[0.98] flex items-center justify-center space-x-2 disabled:opacity-70 disabled:cursor-not-allowed"
          >
            {loading ? (
              <span>正在验证...</span>
            ) : (
              <>
                <span>初始化访问</span>
                <ArrowRightOutlined className="text-sm" />
              </>
            )}
          </button>
        </form>

        <div className="relative my-8">
          <div className="absolute inset-0 flex items-center">
            <div className="w-full border-t border-[#bbc9cb]/20" />
          </div>
          <div className="relative flex justify-center text-xs uppercase tracking-wider">
            <span className="px-4 text-[#3c494b]/40 bg-white/60">其他认证方式</span>
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <button
            type="button"
            className="flex items-center justify-center space-x-2 py-3 px-4 bg-white border border-[#bbc9cb]/20 rounded-xl hover:bg-[#f2f4f5] transition-all active:scale-95 group"
          >
            <SafetyCertificateOutlined className="text-[#4646d8] text-lg group-hover:scale-110 transition-transform" />
            <span className="text-sm font-medium text-[#191c1d]">生物识别</span>
          </button>
          <button
            type="button"
            className="flex items-center justify-center space-x-2 py-3 px-4 bg-white border border-[#bbc9cb]/20 rounded-xl hover:bg-[#f2f4f5] transition-all active:scale-95 group"
          >
            <GlobalOutlined className="text-[#006971] text-lg group-hover:scale-110 transition-transform" />
            <span className="text-sm font-medium text-[#191c1d]">OAuth 登录</span>
          </button>
        </div>
      </div>
    </div>
  )
}

export default LoginForm
