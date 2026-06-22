import { ColorPicker, Divider, Drawer, Select } from 'antd'
import type { Color } from 'antd/es/color-picker'
import useAppStore, { type LayoutMode } from '@/store/useAppStore'

type ThemeDrawerProps = {
  visible: boolean
  onClose: () => void
}

const layoutOptions: Array<{ label: string; value: LayoutMode }> = [
  { label: '混合', value: 'mixed' },
  { label: '经典', value: 'classic' },
]

export const ThemeDrawer = ({ visible, onClose }: ThemeDrawerProps) => {
  const primaryColor = useAppStore((state) => state.primaryColor)
  const layoutMode = useAppStore((state) => state.layoutMode)
  const setPrimaryColor = useAppStore((state) => state.setPrimaryColor)
  const setLayoutMode = useAppStore((state) => state.setLayoutMode)

  const handleColorChange = (_color: Color, cssColor: string) => {
    setPrimaryColor(cssColor)
  }

  return (
    <Drawer title="主题设置" open={visible} width={320} onClose={onClose}>
      <Divider plain>系统主色调</Divider>
      <div className="flex items-center justify-between">
        <span className="text-sm text-slate-600">主色调</span>
        <ColorPicker value={primaryColor} showText onChange={handleColorChange} />
      </div>

      <div className="mt-8 flex items-center justify-between gap-4">
        <span className="text-sm text-slate-600">布局</span>
        <Select
          className="w-40"
          value={layoutMode}
          options={layoutOptions}
          onChange={setLayoutMode}
        />
      </div>
    </Drawer>
  )
}
