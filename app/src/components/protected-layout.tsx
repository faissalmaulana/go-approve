import { Outlet } from "react-router"

import { Header } from "@/components/header"

export function ProtectedLayout() {
  return (
    <div className="min-h-dvh flex flex-col">
      <Header />
      <main className="flex-1">
        <Outlet />
      </main>
    </div>
  )
}
