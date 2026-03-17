import { Navigate, createBrowserRouter } from "react-router"

import { HomePage } from "@/pages/home-page"
import { LoginPage } from "@/pages/login-page"
import { RegisterPage } from "@/pages/register-page"
import { ProtectedLayout } from "@/components/protected-layout"
import { requireAuth, requireGuest } from "@/lib/auth-loader"
import { NewApprovalRoom } from "@/pages/new-approval-room-page"
import { ApprovalRoomDetailPage } from "@/pages/approval-room-detail-page"
import { ApproverRoomDetailPage } from "@/pages/approver-room-detail-page"

const router = createBrowserRouter([
  {
    id: "root",
    children: [
      {
        element: <ProtectedLayout />,
        loader: requireAuth,
        children: [
          {
            path: "/",
            Component: HomePage,
          },
          {
            path: ":id",
            Component: ApprovalRoomDetailPage,
          },
        ],
      },
      {
        element: <ProtectedLayout />,
        loader: requireAuth,
        children: [
          {
            path: "/new-approval-room",
            Component: NewApprovalRoom,
          },
        ],
      },
      {
        element: <ProtectedLayout />,
        loader: requireAuth,
        path: "/approvers",
        children: [
          {
            index: true,
            element: <div>Approvers List</div>,
          },
          {
            path: ":id",
            Component: ApproverRoomDetailPage,
          },
        ],
      },
      {
        path: "/login",
        loader: requireGuest,
        Component: LoginPage,
      },
      {
        path: "/register",
        loader: requireGuest,
        Component: RegisterPage,
      },
      {
        path: "*",
        element: <Navigate to="/" replace />,
      },
    ],
  },
])

export { router }
