import { Link } from "react-router"

import { useActivePath } from "@/hooks/use-active-path"

import { Separator } from "./ui/separator"
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "./ui/sidebar"

const navigationsData = {
  navMain: [
    {
      title: "Approvals",
      url: "/",
      items: [
        {
          title: "Requests",
          url: "/",
        },
        {
          title: "Approvers",
          url: "/approvers",
        },
      ]
    },
    {
      title: "Invitations",
      url: "/invitations",
      items: [
        {
          title: "Sent",
          url: "/invitations/sent",
        },
        {
          title: "Received",
          url: "/invitations/received",
        },
      ]
    }
  ]
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const { isActive } = useActivePath()

  return (
    <Sidebar {...props}>
      <SidebarHeader className="h-14 flex items-start justify-center">
        <div className="group flex items-center gap-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground text-sm font-bold shadow-sm transition-transform group-hover:scale-105">
            GA
          </div>
          <p className="text-primary text-lg font-semibold">GoApprove</p>
        </div>
      </SidebarHeader>
      <Separator />
      <SidebarContent>
        {navigationsData.navMain.map((navItem) => (
          <SidebarGroup key={navItem.title}>
            <SidebarGroupLabel>{navItem.title}</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {navItem.items.map((item) => (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton
                      render={<Link to={item.url} />}
                      isActive={isActive(item.url, {
                        end: item.url === "/" || item.url.startsWith("/invitations/"),
                      })}
                    >
                      {item.title}
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        ))}
      </SidebarContent>
    </Sidebar>
  )
}
