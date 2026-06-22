import { getConfigInfoApi } from '@/api/system/config';
import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';

export type LayoutMode = 'mixed' | 'classic';

interface AppState {
  siderCollapsed: boolean;
  layoutMode: LayoutMode;
  primaryColor: string;
  setSiderCollapsed: (collapsed: boolean) => void;
  toggleSider: () => void;
  setLayoutMode: (mode: LayoutMode) => void;
  setPrimaryColor: (color: string) => void;
  siteConfig: any;
  setSiteConfig: (config: any) => void;
  initSiteConfig: () => Promise<void>;
}

const useAppStore = create<AppState>()(
  devtools(
    persist(
      (set) => ({
        siderCollapsed: false,
        layoutMode: 'mixed',
        primaryColor: '#1677ff',
        siteConfig: null,

        setSiderCollapsed: (collapsed) => set({ siderCollapsed: collapsed }),

        toggleSider: () => set((state) => ({ siderCollapsed: !state.siderCollapsed })),

        setLayoutMode: (mode) => set({ layoutMode: mode }),

        setPrimaryColor: (color) => set({ primaryColor: color }),

        setSiteConfig: (config) => set({ siteConfig: config }),

        initSiteConfig: async () => {
          const res = await getConfigInfoApi('site_setting');
          console.log("site_setting", res);
          set({ siteConfig: res.data });
        }
      }),
      {
        name: 'app-storage',
        partialize: (state) => ({
          siderCollapsed: state.siderCollapsed,
          layoutMode: state.layoutMode,
          primaryColor: state.primaryColor,
          siteConfig: state.siteConfig,
        }),
      }
    )
  )
);

export default useAppStore;
