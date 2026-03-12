import { Link, useNavigate } from "react-router"
import { useMutation } from "@tanstack/react-query"

import { LoginForm } from "@/components/login-form"
import { setAccessToken } from "@/lib/auth"
import { postJson } from "@/lib/api"

type SignInResponse = { access_token: string }

export function LoginPage() {
  const navigate = useNavigate()
  const signInMutation = useMutation({
    mutationFn: (values: { email: string; password: string }) =>
      postJson<{ email: string; password: string }, SignInResponse>(
        "/auth/sign-in",
        values
      ),
  })

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
          onSubmitValues={async (values) => {
            const res = await signInMutation.mutateAsync(values)

            if (res.error) {
              return { ok: false as const, error: res.error }
            }

            const token = res.data?.access_token
            if (!token) {
              return { ok: false as const, error: "Missing access token." }
            }

            setAccessToken(token)
            navigate("/", { replace: true })
            return { ok: true as const }
          }}
        />
      </div>
    </main>
  )
}

