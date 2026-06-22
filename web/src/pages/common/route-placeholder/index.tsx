type RoutePlaceholderProps = {
  routePath: string
  componentPath?: string
}

const RoutePlaceholder = ({ routePath, componentPath }: RoutePlaceholderProps) => {
  return (
    <div className="flex h-full min-h-[240px] items-center justify-center rounded-lg border border-dashed border-slate-300 bg-slate-50 p-8 text-center text-slate-500">
      <div>
        <div className="text-base font-medium text-slate-700">页面组件尚未创建</div>
        <div className="mt-2 text-sm">路由：{routePath}</div>
        {componentPath ? <div className="mt-1 text-sm">组件：{componentPath}</div> : null}
      </div>
    </div>
  )
}

export default RoutePlaceholder
