import { useNavigate } from "react-router"

import { Button } from "@/components/ui/button"
import { useUser, useLogout } from "@/hooks/use-auth"
import { SidebarTrigger } from "./ui/sidebar"

export function Header() {
  const navigate = useNavigate()
  const { data: user } = useUser()
  const logout = useLogout()

  const handleLogout = async () => {
    await logout.mutateAsync()
    navigate("/login", { replace: true })
  }

  return (
    <header className="border-b bg-background/95 backdrop-blur supports-backdrop-filter:bg-background">
      <div className="container mx-auto flex h-14 items-center justify-between px-4">
        <SidebarTrigger />
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-2">
            <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary text-sm font-medium text-primary-foreground">
              {user?.name?.charAt(0).toUpperCase() ?? "?"}
            </div>
            <span className="text-sm font-medium">{user?.name ?? "Loading..."}</span>
          </div>
          <Button
            variant="outline"
            size="sm"
            className={"cursor-pointer"}
            onClick={handleLogout}
            disabled={logout.isPending}
          >
            {logout.isPending ? "Logging out..." : "Logout"}
          </Button>
        </div>
      </div>
    </header>
  )
}
