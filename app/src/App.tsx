import { Navigate, RouterProvider, createBrowserRouter } from "react-router"

import { HomePage } from "@/pages/home-page"
import { LoginPage } from "@/pages/login-page"
import { RegisterPage } from "@/pages/register-page"

function App() {
  const router = createBrowserRouter([
    { path: "/", element: <HomePage /> },
    { path: "/login", element: <LoginPage /> },
    { path: "/register", element: <RegisterPage /> },
    { path: "*", element: <Navigate to="/" replace /> },
  ])

  return <RouterProvider router={router} />
}

export default App
