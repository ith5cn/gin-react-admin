import { create } from 'zustand'
import { createJSONStorage, persist } from 'zustand/middleware'
import { userApi } from '@/api/auth'
import type { LoginResult, UserInfo } from '@/api/auth'
import { buildTagFromPath, createHomeTag, HOME_PATH, HOME_TAG_KEY } from '@/routers/menuHelpers'
import { normalizeBackendMenuTree } from '@/routers/normalizeMenu'
import { transformMenuToRoutes, transformMenuToSiderItems } from '@/routers/transformRoutes'
import type { AppMenuItem, AppRouteObject, NormalizedMenuNode, RawBackendMenuNode, TagViewItem } from '@/types/router'

type AuthState = {
  token: string
  refreshToken: string
  expiresIn: string
  userInfo: UserInfo | null
  roles: Array<number | string>
  codes: string[]
  posts: unknown
  depts: unknown
  rawMenuTree: RawBackendMenuNode[]
  normalizedMenuTree: NormalizedMenuNode[]
  sideMenuItems: AppMenuItem[]
  dynamicRoutes: AppRouteObject[]
  visitedTags: TagViewItem[]
  activeTagKey: string
  isUserInitialized: boolean
  isInitializing: boolean
  login: (session: LoginResult) => Promise<void>
  logout: () => void
  initUserContext: () => Promise<void>
  resetUserContext: () => void
  syncTagByPath: (pathname: string) => void
  activateTag: (tagKey: string) => void
  closeTag: (tagKey: string) => string | null
  closeLeftTags: (tagKey: string) => string | null
  closeRightTags: (tagKey: string) => string | null
  closeOtherTags: (tagKey: string) => string | null
  closeAllTags: () => string
}

const resetAuthState = {
  token: '',
  refreshToken: '',
  expiresIn: '',
  userInfo: null,
  roles: [],
  codes: [],
  posts: null,
  depts: null,
  rawMenuTree: [],
  normalizedMenuTree: [],
  sideMenuItems: [],
  dynamicRoutes: [],
  visitedTags: [createHomeTag()],
  activeTagKey: HOME_TAG_KEY,
  isUserInitialized: false,
  isInitializing: false,
} satisfies Pick<
  AuthState,
  | 'token'
  | 'refreshToken'
  | 'expiresIn'
  | 'userInfo'
  | 'roles'
  | 'codes'
  | 'posts'
  | 'depts'
  | 'rawMenuTree'
  | 'normalizedMenuTree'
  | 'sideMenuItems'
  | 'dynamicRoutes'
  | 'visitedTags'
  | 'activeTagKey'
  | 'isUserInitialized'
  | 'isInitializing'
>

const normalizePermissionCodes = (codes: string | string[] | null | undefined): string[] => {
  if (!codes) {
    return []
  }

  if (Array.isArray(codes)) {
    return codes.map((code) => String(code)).filter(Boolean)
  }

  return [String(codes)].filter(Boolean)
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      ...resetAuthState,
      login: async (session) => {
        set({
          ...resetAuthState,
          token: session.accessToken,
          refreshToken: session.refreshToken ?? '',
          expiresIn: session.expiresIn ?? '',
        })
      },
      logout: () => {
        set(resetAuthState)
      },
      initUserContext: async () => {
        const { token, isInitializing, isUserInitialized } = get()

        if (!token || isInitializing || isUserInitialized) {
          return
        }

        set({ isInitializing: true })

        try {
          const userContext = await userApi()
          const rawMenuTree = userContext.routers as RawBackendMenuNode[]
          const normalizedMenuTree = normalizeBackendMenuTree(rawMenuTree)
          const dynamicRoutes = transformMenuToRoutes(normalizedMenuTree)
          const sideMenuItems = transformMenuToSiderItems(normalizedMenuTree)

          set({
            userInfo: userContext.user,
            roles: userContext.roles ?? [],
            codes: normalizePermissionCodes(userContext.codes),
            posts: userContext.posts ?? null,
            depts: userContext.depts ?? null,
            rawMenuTree,
            normalizedMenuTree,
            sideMenuItems,
            dynamicRoutes,
            visitedTags: [createHomeTag()],
            activeTagKey: HOME_TAG_KEY,
            isUserInitialized: true,
            isInitializing: false,
          })
        } catch (error) {
          set({
            ...resetAuthState,
            token: get().token,
            refreshToken: get().refreshToken,
            expiresIn: get().expiresIn,
          })
          throw error
        }
      },
      resetUserContext: () => {
        set({
          userInfo: null,
          roles: [],
          codes: [],
          posts: null,
          depts: null,
          rawMenuTree: [],
          normalizedMenuTree: [],
          sideMenuItems: [],
          dynamicRoutes: [],
          visitedTags: [createHomeTag()],
          activeTagKey: HOME_TAG_KEY,
          isUserInitialized: false,
          isInitializing: false,
        })
      },
      syncTagByPath: (pathname: string) => {
        const { sideMenuItems } = get()
        const matchedTag = buildTagFromPath(pathname, sideMenuItems)

        if (!matchedTag) {
          return
        }

        set((state) => {
          const hasTag = state.visitedTags.some((item) => item.key === matchedTag.key)

          return {
            visitedTags: hasTag ? state.visitedTags : [...state.visitedTags, matchedTag],
            activeTagKey: matchedTag.key,
          }
        })
      },
      activateTag: (tagKey: string) => {
        set({ activeTagKey: tagKey })
      },
      closeTag: (tagKey: string) => {
        const { visitedTags, activeTagKey } = get()
        const targetIndex = visitedTags.findIndex((item) => item.key === tagKey)

        if (targetIndex === -1 || tagKey === HOME_TAG_KEY) {
          return null
        }

        const nextVisitedTags = visitedTags.filter((item) => item.key !== tagKey)

        if (activeTagKey !== tagKey) {
          set({ visitedTags: nextVisitedTags })
          return null
        }

        const fallbackTag = visitedTags[targetIndex + 1] ?? visitedTags[targetIndex - 1] ?? createHomeTag()

        set({
          visitedTags: nextVisitedTags,
          activeTagKey: fallbackTag.key,
        })

        return fallbackTag.path
      },
      closeLeftTags: (tagKey: string) => {
        const { visitedTags, activeTagKey } = get()
        const targetIndex = visitedTags.findIndex((item) => item.key === tagKey)

        if (targetIndex <= 1) {
          return null
        }

        const nextVisitedTags = visitedTags.filter((_, index) => index === 0 || index >= targetIndex)
        const nextActiveTag =
          nextVisitedTags.find((item) => item.key === activeTagKey) ??
          nextVisitedTags.find((item) => item.key === tagKey) ??
          createHomeTag()

        set({
          visitedTags: nextVisitedTags,
          activeTagKey: nextActiveTag.key,
        })

        return nextActiveTag.key === activeTagKey ? null : nextActiveTag.path
      },
      closeRightTags: (tagKey: string) => {
        const { visitedTags, activeTagKey } = get()
        const targetIndex = visitedTags.findIndex((item) => item.key === tagKey)

        if (targetIndex === -1 || targetIndex >= visitedTags.length - 1) {
          return null
        }

        const nextVisitedTags = visitedTags.filter((_, index) => index <= targetIndex)
        const nextActiveTag =
          nextVisitedTags.find((item) => item.key === activeTagKey) ??
          nextVisitedTags.find((item) => item.key === tagKey) ??
          createHomeTag()

        set({
          visitedTags: nextVisitedTags,
          activeTagKey: nextActiveTag.key,
        })

        return nextActiveTag.key === activeTagKey ? null : nextActiveTag.path
      },
      closeOtherTags: (tagKey: string) => {
        const { visitedTags, activeTagKey } = get()
        const targetTag = visitedTags.find((item) => item.key === tagKey)

        if (!targetTag) {
          return null
        }

        const nextVisitedTags = visitedTags.filter((item) => item.key === HOME_TAG_KEY || item.key === tagKey)
        const nextActiveTag = nextVisitedTags.find((item) => item.key === activeTagKey) ?? targetTag ?? createHomeTag()

        set({
          visitedTags: nextVisitedTags,
          activeTagKey: nextActiveTag.key,
        })

        return nextActiveTag.key === activeTagKey ? null : nextActiveTag.path
      },
      closeAllTags: () => {
        set({
          visitedTags: [createHomeTag()],
          activeTagKey: HOME_TAG_KEY,
        })

        return HOME_PATH
      },
    }),
    {
      name: 'web-auth',
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        token: state.token,
        refreshToken: state.refreshToken,
        expiresIn: state.expiresIn,
      }),
    },
  ),
)

export const clearAuthSession = () => {
  useAuthStore.getState().logout()
  useAuthStore.persist.clearStorage()
}
