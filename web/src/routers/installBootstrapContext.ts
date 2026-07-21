import { createContext, useContext } from 'react'

export interface InstallBootstrapContextValue {
  markInstalled: () => void
}

export const InstallBootstrapContext = createContext<InstallBootstrapContextValue>({
  markInstalled: () => undefined,
})

export const useInstallBootstrap = () => useContext(InstallBootstrapContext)