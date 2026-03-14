import { Outlet } from "react-router"
import { Header } from "@/components/header"
import { SidebarProvider } from "./ui/sidebar"
import { AppSidebar } from "./app-sidebar"

export function ProtectedLayout() {
  return (
    <SidebarProvider>
      <div className="flex w-full">
        <AppSidebar />
        <div className="min-h-dvh flex flex-col w-full">
          <Header />
          <main className="flex-1">
            <Outlet />
          </main>
        </div>
      </div>
    </SidebarProvider>
  )
}
