import React from 'react'
import { Navigate } from 'react-router-dom'
import { useAuthStore } from '@/store/auth'
import LoginForm from './login-form'

const Login: React.FC = () => {
    const token = useAuthStore((state) => state.token)

    if (token) {
        return <Navigate to="/dashboard" replace />
    }

    return (
        <div className="flex min-h-screen flex-col bg-[#f6f8fb] font-sans text-[#1f2328] antialiased">
            <header className="fixed top-0 z-50 w-full border-b border-[#d8dee4] bg-white/95 backdrop-blur">
                <div className="mx-auto flex h-16 w-full max-w-[1280px] items-center justify-between px-4 md:px-8">
                    <div className="text-[22px] font-bold leading-[1.2] text-[#1f2328]">Gin React Admin</div>
                    <a className="text-sm text-[#57606a] transition-colors hover:text-[#1677ff]" href="/install">
                        首次安装
                    </a>
                </div>
            </header>

            <main className="flex flex-1 items-center justify-center px-4 pb-16 pt-24 md:px-8">
                <LoginForm />
            </main>

            <footer className="w-full border-t border-[#d8dee4] bg-white py-2" />
        </div>
    )
}

export default Login
