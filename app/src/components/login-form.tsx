import { zodResolver } from "@hookform/resolvers/zod"
import { useState } from "react"
import { useForm } from "react-hook-form"
import { z } from "zod"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { cn } from "@/lib/utils"

const loginSchema = z.object({
  email: z.string().email("Please enter a valid email."),
  password: z.string().min(1, "Password is required."),
})

export type LoginValues = z.infer<typeof loginSchema>

type SubmitResult = { ok: true } | { ok: false; error: string }

export function LoginForm({
  className,
  title = "Login",
  description,
  submitLabel = "Continue",
  bottomText,
  onSubmitValues,
  ...props
}: React.ComponentProps<"div"> & {
  title?: string
  description?: React.ReactNode
  submitLabel?: string
  bottomText?: React.ReactNode
  onSubmitValues: (values: LoginValues) => Promise<SubmitResult>
}) {
  const [serverError, setServerError] = useState<string | null>(null)

  const form = useForm<LoginValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: "", password: "" },
    mode: "onSubmit",
  })

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">{title}</CardTitle>
          {description ? <CardDescription>{description}</CardDescription> : null}
        </CardHeader>
        <CardContent>
          <form
            onSubmit={form.handleSubmit(async (values) => {
              setServerError(null)
              const result = await onSubmitValues(values)
              if (!result.ok) setServerError(result.error)
            })}
            className="flex flex-col gap-5"
          >
            <FieldGroup>
              {serverError ? (
                <Field>
                  <FieldError>{serverError}</FieldError>
                </Field>
              ) : null}
              <Field>
                <FieldLabel htmlFor="email">Email</FieldLabel>
                <Input
                  id="email"
                  type="email"
                  placeholder="m@example.com"
                  autoComplete="email"
                  aria-invalid={!!form.formState.errors.email}
                  {...form.register("email")}
                />
                <FieldError errors={[form.formState.errors.email]} />
              </Field>
              <Field>
                <FieldLabel htmlFor="password">Password</FieldLabel>
                <Input
                  id="password"
                  type="password"
                  autoComplete="current-password"
                  aria-invalid={!!form.formState.errors.password}
                  {...form.register("password")}
                />
                <FieldError errors={[form.formState.errors.password]} />
              </Field>
              <Field>
                <Button type="submit" disabled={form.formState.isSubmitting}>
                  {submitLabel}
                </Button>
                {bottomText ? (
                  <FieldDescription className="text-center">
                    {bottomText}
                  </FieldDescription>
                ) : null}
              </Field>
            </FieldGroup>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
