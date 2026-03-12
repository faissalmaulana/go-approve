import { Link } from "react-router"

import { LoginForm } from "@/components/login-form"

export function LoginPage() {
  return (
    <main className="min-h-dvh px-4 py-10">
      <div className="mx-auto flex w-full max-w-md flex-col gap-6">
        <LoginForm
          title="Welcome back"
          description="Login with your email and password."
          submitLabel="Login"
          bottomText={
            <>
              Don&apos;t have an account?{" "}
              <Link to="/register" className="underline underline-offset-4">
                Sign up
              </Link>
            </>
          }
        />
      </div>
    </main>
  )
}
