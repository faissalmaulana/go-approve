import { Link, useNavigate } from "react-router"
import { useMutation } from "@tanstack/react-query"

import { RegisterForm } from "@/components/register-form"
import { postJson } from "@/lib/api"
import { setAccessToken } from "@/lib/auth"

type SignUpResponse = { access_token: string }

export function RegisterPage() {
  const navigate = useNavigate()
  const signUpMutation = useMutation({
    mutationFn: (values: {
      email: string
      name: string
      handler: string
      password: string
    }) =>
      postJson<
        { email: string; name: string; handler: string; password: string },
        SignUpResponse
      >("/auth/sign-up", values),
  })

  return (
    <main className="min-h-dvh px-4 py-10">
      <div className="mx-auto flex w-full max-w-md flex-col gap-6">
        <RegisterForm
          onSubmitValues={async (values) => {
            const res = await signUpMutation.mutateAsync({
              email: values.email,
              name: values.name,
              handler: values.handler,
              password: values.password,
            })

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
          bottomText={
            <>
              Already have an account?{" "}
              <Link to="/login" className="underline underline-offset-4">
                Login
              </Link>
            </>
          }
        />
      </div>
    </main>
  )
}

