import React from 'react'
import {
    BranchesOutlined,
    CustomerServiceOutlined,
    DatabaseOutlined,
    FundProjectionScreenOutlined,
    ThunderboltOutlined,
} from '@ant-design/icons'

const nodes = [
    { label: '机器人编排', icon: <BranchesOutlined />, className: 'left-1/2 top-3 -translate-x-1/2' },
    { label: '知识库召回', icon: <DatabaseOutlined />, className: 'right-4 top-1/3 -translate-y-1/2' },
    { label: '客户触达', icon: <CustomerServiceOutlined />, className: 'right-12 bottom-14' },
    { label: '任务自动化', icon: <ThunderboltOutlined />, className: 'left-12 bottom-14' },
    { label: '转化分析', icon: <FundProjectionScreenOutlined />, className: 'left-4 top-1/3 -translate-y-1/2' },
]

const telemetry = [
    { label: '活跃机器人', value: '128', unit: 'bots' },
    { label: '会话接管率', value: '96.8', unit: '%' },
    { label: '自动跟进', value: '2,436', unit: 'tasks' },
    { label: '转化线索', value: '684', unit: 'leads' },
]

const streams = ['客户意图识别', '知识库检索', '机器人响应', '自动任务派发']

const HeroSection: React.FC = () => {
    return (
        <div className="relative mx-auto aspect-[1/1.05] w-full max-w-[680px] min-w-0 sm:aspect-square">
            <div className="absolute inset-0 rounded-full border border-cyan-200/10 bg-[radial-gradient(circle,rgba(14,165,233,0.18),transparent_48%)]" />
            <div className="absolute inset-[7%] rounded-full border border-cyan-200/15" />
            <div className="absolute inset-[15%] rounded-full border border-emerald-200/15" />
            <div className="absolute inset-[25%] rounded-full border border-cyan-200/10" />

            <div className="absolute left-1/2 top-1/2 h-[72%] w-px -translate-x-1/2 -translate-y-1/2 bg-[linear-gradient(180deg,transparent,rgba(125,211,252,0.32),transparent)]" />
            <div className="absolute left-1/2 top-1/2 h-px w-[72%] -translate-x-1/2 -translate-y-1/2 bg-[linear-gradient(90deg,transparent,rgba(125,211,252,0.32),transparent)]" />
            <div className="absolute left-1/2 top-1/2 h-px w-[70%] -translate-x-1/2 -translate-y-1/2 rotate-45 bg-[linear-gradient(90deg,transparent,rgba(34,197,94,0.22),transparent)]" />
            <div className="absolute left-1/2 top-1/2 h-px w-[70%] -translate-x-1/2 -translate-y-1/2 -rotate-45 bg-[linear-gradient(90deg,transparent,rgba(34,197,94,0.22),transparent)]" />

            <div className="absolute inset-[21%] animate-[spin_26s_linear_infinite] rounded-full border border-dashed border-cyan-200/20" />
            <div className="absolute inset-[32%] animate-[spin_18s_linear_infinite_reverse] rounded-full border border-dashed border-emerald-200/20" />

            <div className="absolute left-1/2 top-1/2 z-10 flex h-40 w-40 -translate-x-1/2 -translate-y-1/2 flex-col items-center justify-center rounded-full border border-cyan-200/35 bg-[#06111f]/95 text-center shadow-[0_0_70px_rgba(34,211,238,0.26)] sm:h-44 sm:w-44">
                <div className="absolute inset-3 rounded-full border border-cyan-200/10" />
                <div className="text-xs tracking-[0.28em] text-cyan-100">AI CORE</div>
                <div className="mt-2 text-xl font-semibold text-white">运营中枢</div>
                <div className="mt-2 text-xs text-slate-400">Command Engine</div>
                <span className="mt-4 h-1.5 w-16 overflow-hidden rounded-full bg-slate-800">
                    <span className="block h-full w-2/3 animate-pulse rounded-full bg-cyan-300" />
                </span>
            </div>

            {nodes.map((node) => (
                <div key={node.label} className={`absolute z-20 ${node.className}`}>
                    <div className="flex min-w-28 items-center gap-2 rounded-xl border border-white/12 bg-[#071827]/90 px-3 py-2 text-xs text-slate-200 shadow-[0_18px_44px_rgba(2,6,23,0.34)] backdrop-blur-md">
                        <span className="flex h-8 w-8 items-center justify-center rounded-lg border border-cyan-300/20 bg-cyan-300/10 text-cyan-100">
                            {node.icon}
                        </span>
                        <span className="whitespace-nowrap">{node.label}</span>
                    </div>
                </div>
            ))}

            <div className="absolute inset-x-3 bottom-0 z-30 rounded-2xl border border-white/10 bg-[#06111f]/92 p-3 shadow-[0_24px_80px_rgba(0,0,0,0.42)] backdrop-blur-xl sm:inset-x-6 sm:p-4">
                <div className="mb-3 flex items-center justify-between">
                    <div>
                        <div className="text-sm font-semibold text-white">实时运营态势</div>
                        <div className="mt-1 text-xs text-slate-500">AI routing telemetry</div>
                    </div>
                    <span className="rounded-md border border-emerald-300/20 px-2 py-1 text-xs text-emerald-100">LIVE</span>
                </div>
                <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
                    {telemetry.map((item) => (
                        <div key={item.label}>
                            <div className="text-lg font-semibold text-white">{item.value}</div>
                            <div className="mt-0.5 text-[11px] uppercase tracking-wide text-slate-500">{item.unit}</div>
                            <div className="mt-1 text-xs text-slate-400">{item.label}</div>
                        </div>
                    ))}
                </div>
            </div>

            <div className="absolute left-5 top-12 hidden w-44 rounded-xl border border-white/10 bg-white/[0.045] p-3 text-xs text-slate-400 backdrop-blur-md sm:block">
                {streams.map((item, index) => (
                    <div key={item} className="flex items-center justify-between border-b border-white/8 py-2 last:border-b-0">
                        <span>{item}</span>
                        <span className="text-cyan-100">{String(index + 1).padStart(2, '0')}</span>
                    </div>
                ))}
            </div>
        </div>
    )
}

export default HeroSection
