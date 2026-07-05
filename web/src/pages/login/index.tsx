import React from 'react'
import { Navigate } from 'react-router-dom'
import { useAuthStore } from '@/store/auth'

import LoginForm from './login-form'

const stats = [
  { label: 'Latency', value: '2.4', unit: 'ms', color: 'text-[#006971]' },
  { label: 'Uptime', value: '99.9', unit: '%', color: 'text-[#4646d8]' },
  { label: 'Data', value: '128', unit: 'TB', color: 'text-[#006971]' },
]

const avatarLetters = ['A', 'B', 'C']

const Login: React.FC = () => {
  const token = useAuthStore((state) => state.token)

  if (token) {
    return <Navigate to="/dashboard" replace />
  }

  return (
    <main className="relative min-h-screen flex items-center justify-center px-6 overflow-hidden bg-[#f8fafb]">
      {/* Background auras */}
      <div className="absolute top-[-10%] right-[-5%] w-[600px] h-[600px] rounded-full blur-[120px] bg-[#006971]/5 pointer-events-none" />
      <div className="absolute bottom-[-10%] left-[-5%] w-[500px] h-[500px] rounded-full blur-[120px] bg-[#4646d8]/5 pointer-events-none" />

      {/* Dot grid overlay */}
      <div
        className="absolute inset-0 pointer-events-none"
        style={{
          backgroundImage: 'radial-gradient(circle at 2px 2px, rgba(0,105,113,0.03) 1px, transparent 0)',
          backgroundSize: '40px 40px',
        }}
      />

      <div className="relative z-10 w-full max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-2 gap-20 items-center py-16">
        {/* Left: Branding & Stats */}
        <div className="hidden lg:flex flex-col space-y-12">
          <div className="space-y-6">
            <h1 className="text-6xl font-bold text-[#191c1d] tracking-tight leading-tight font-headline">
              Gin React<br />
              <span className="text-[#006971]">管理后台</span>
            </h1>
            <p className="text-xl text-[#3c494b]/80 max-w-md leading-relaxed">
              进入天枢架构，调度智能代理，管理高维数据集。
            </p>
          </div>

          {/* Stats grid */}
          <div className="grid grid-cols-3 gap-6 max-w-lg">
            {stats.map((stat) => (
              <div key={stat.label} className="p-6 rounded-xl bg-[#f2f4f5] border border-[#bbc9cb]/10">
                <div className={`text-xs font-medium tracking-wider mb-2 uppercase ${stat.color}`}>
                  {stat.label}
                </div>
                <div className="text-3xl font-bold text-[#191c1d] font-headline">
                  {stat.value}
                  <span className="text-sm font-normal text-[#3c494b] ml-1">{stat.unit}</span>
                </div>
              </div>
            ))}
          </div>

          {/* Social proof */}
          <div className="flex items-center space-x-4">
            <div className="flex -space-x-3">
              {avatarLetters.map((letter) => (
                <div
                  key={letter}
                  className="w-10 h-10 rounded-full border-2 border-[#f8fafb] bg-[#006971]/20 flex items-center justify-center text-[#006971] text-xs font-bold"
                >
                  {letter}
                </div>
              ))}
              <div className="w-10 h-10 rounded-full border-2 border-[#f8fafb] bg-[#3abbc9] flex items-center justify-center text-white text-[10px] font-bold">
                +5k
              </div>
            </div>
            <div className="text-sm text-[#3c494b] font-medium">
              当前已有 5,000+ 企业部署天枢集群
            </div>
          </div>
        </div>

        {/* Right: Login card */}
        <div className="flex justify-center lg:justify-end">
          <LoginForm />
        </div>
      </div>
    </main>
  )
}

export default Login
