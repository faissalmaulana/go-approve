import { matchPath, useLocation } from "react-router"

type IsActiveOptions = {
  /**
   * When true, requires an exact match (no nested routes).
   * Defaults to true for "/", otherwise false.
   */
  end?: boolean
}

function normalizePathname(pathname: string) {
  if (pathname.length > 1 && pathname.endsWith("/")) {
    return pathname.slice(0, -1)
  }
  return pathname
}

export function useActivePath() {
  const location = useLocation()
  const pathname = normalizePathname(location.pathname)

  const isActive = (to: string, options?: IsActiveOptions) => {
    const target = normalizePathname(to)
    const end = options?.end ?? (target === "/")

    if (target === "/") {
      return pathname === "/"
    }

    return !!matchPath({ path: target, end }, pathname)
  }

  return { pathname, isActive }
}

