import { Separator } from "./ui/separator";
import { Sidebar, SidebarContent, SidebarGroup, SidebarGroupContent, SidebarGroupLabel, SidebarHeader, SidebarMenu, SidebarMenuBadge, SidebarMenuButton, SidebarMenuItem } from "./ui/sidebar";

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
          withBadge: true,
        },
      ]
    }
  ]
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar {...props}>
      <SidebarHeader className="h-14 flex justify-center">
        <div className="font-semibold text-xl text-blue-500">Go Approve</div>
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
                    <SidebarMenuButton >
                      <a href={item.url}>{item.title}</a>
                    </SidebarMenuButton>
                    {item.withBadge && (
                      <SidebarMenuBadge>25</SidebarMenuBadge>
                    )
                    }
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
