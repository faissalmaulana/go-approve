import { Link } from "react-router"

import { RegisterForm } from "@/components/register-form"

export function RegisterPage() {
  return (
    <main className="min-h-dvh px-4 py-10">
      <div className="mx-auto flex w-full max-w-md flex-col gap-6">
        <RegisterForm
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
