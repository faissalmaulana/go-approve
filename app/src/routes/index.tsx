import { Navigate, createBrowserRouter } from "react-router"

import { HomePage } from "@/pages/home-page"
import { LoginPage } from "@/pages/login-page"
import { RegisterPage } from "@/pages/register-page"
import { requireAuth, requireGuest } from "@/lib/auth-loader"

const router = createBrowserRouter([
  {
    id: "root",
    children: [
      {
        path: "/",
        loader: requireAuth,
        Component: HomePage,
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
