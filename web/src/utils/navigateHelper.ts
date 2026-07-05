type NavigateFn = (path: string, options?: { replace?: boolean }) => void

let _navigate: NavigateFn | null = null

export const setNavigate = (fn: NavigateFn) => {
  _navigate = fn
}

export const navigateTo = (path: string, replace = true) => {
  if (_navigate) {
    _navigate(path, { replace })
  } else {
    // Fallback before React Router mounts (e.g. SSR or early interceptor calls)
    window.location.href = path
  }
}
