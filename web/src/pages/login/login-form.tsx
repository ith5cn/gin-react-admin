import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button, message } from 'antd'
import {
    LockOutlined,
    UserOutlined,
} from '@ant-design/icons'
import { useAuthStore } from '@/store/auth'
import { loginApi } from '@/api/auth'

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
            await useAuthStore.getState().login(res)
            await useAuthStore.getState().initUserContext()
            message.success('登录成功')
            navigate('/dashboard', { replace: true })
        } catch {
            return
        } finally {
            setLoading(false)
        }
    }

    return (
        <section className="w-full max-w-md rounded-lg border border-[#d8dee4] bg-white p-8 shadow-[0_12px_30px_rgba(31,35,40,0.08)]">
            <div className="mb-8 text-center">
                <h1 className="mb-2 text-[30px] font-semibold leading-[1.15] tracking-normal text-[#1f2328]">
                    登录管理后台
                </h1>
                <p className="m-0 text-sm leading-6 text-[#57606a]">请输入账号密码进入控制台</p>
            </div>

            <form className="space-y-6" onSubmit={onSubmit}>
                <div className="space-y-2">
                    <label className="block text-sm font-medium leading-none text-[#1f2328]" htmlFor="username">
                        用户名
                    </label>
                    <div className="relative">
                        <UserOutlined className="pointer-events-none absolute left-3 top-1/2 z-10 -translate-y-1/2 text-sm text-[#57606a]" />
                        <input
                            id="username"
                            className="h-11 w-full rounded-md border border-[#d0d7de] bg-white py-2 pl-10 pr-3 text-base leading-6 text-[#1f2328] outline-none transition-shadow duration-200 placeholder:text-[#8c959f] focus:border-[#1677ff] focus:shadow-[0_0_0_2px_rgba(22,119,255,0.15)]"
                            placeholder="请输入用户名"
                            autoComplete="username"
                            required
                            value={values.username}
                            onChange={(event) => setValues((prev) => ({ ...prev, username: event.target.value }))}
                        />
                    </div>
                </div>

                <div className="space-y-2">
                    <div className="flex items-center justify-between">
                        <label className="block text-sm font-medium leading-none text-[#1f2328]" htmlFor="password">
                            密码
                        </label>
                        <a className="text-xs leading-[1.45] text-[#1677ff] hover:underline" href="#">
                            忘记密码？
                        </a>
                    </div>
                    <div className="relative">
                        <LockOutlined className="pointer-events-none absolute left-3 top-1/2 z-10 -translate-y-1/2 text-sm text-[#57606a]" />
                        <input
                            id="password"
                            type="password"
                            className="h-11 w-full rounded-md border border-[#d0d7de] bg-white py-2 pl-10 pr-3 text-base leading-6 text-[#1f2328] outline-none transition-shadow duration-200 placeholder:text-[#8c959f] focus:border-[#1677ff] focus:shadow-[0_0_0_2px_rgba(22,119,255,0.15)]"
                            placeholder="请输入密码"
                            autoComplete="current-password"
                            required
                            value={values.password}
                            onChange={(event) => setValues((prev) => ({ ...prev, password: event.target.value }))}
                        />
                    </div>
                </div>

                <div className="flex items-center">
                    <input
                        id="remember-me"
                        name="remember-me"
                        type="checkbox"
                        className="h-4 w-4 cursor-pointer rounded border-[#d0d7de] text-[#1677ff] focus:ring-[#1677ff]"
                    />
                    <label className="ml-2 block cursor-pointer text-sm leading-6 text-[#57606a]" htmlFor="remember-me">
                        记住我
                    </label>
                </div>

                <Button
                    htmlType="submit"
                    block
                    loading={loading}
                    className="!h-10 !rounded-md !border-none !bg-[#1677ff] !px-4 !py-2 !text-sm !font-medium !leading-none !text-white !shadow-sm transition-colors duration-200 hover:!bg-[#0958d9] hover:!text-white"
                >
                    立即登录
                </Button>
            </form>

            {/* <div className="relative mt-8">
                <div aria-hidden="true" className="absolute inset-0 flex items-center">
                    <div className="w-full border-t border-[#dfdfdf]" />
                </div>
                <div className="relative flex justify-center">
                    <span className="bg-white px-2 text-xs leading-[1.45] text-[#3d4a41]">或者</span>
                </div>
            </div>

            <div className="mt-8 space-y-4">
                <button
                    className="flex h-10 w-full items-center justify-center gap-2 rounded-md border border-[#c7c7c7] bg-white px-4 py-2 text-sm font-medium leading-none text-[#1b1c1c] transition-colors duration-200 hover:bg-[#f5f3f3]"
                    type="button"
                >
                    <GithubOutlined className="text-xl" />
                    使用 GitHub 登录
                </button>
                <div className="text-center">
                    <a className="text-base leading-6 text-[#3d4a41] transition-colors duration-200 hover:text-[#3ecf8e]" href="#">
                        单点登录 (SSO)
                    </a>
                </div>
            </div> */}
        </section>
    )
}

export default LoginForm
