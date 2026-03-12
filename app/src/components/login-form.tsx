import { zodResolver } from "@hookform/resolvers/zod"
import { useState } from "react"
import { useForm } from "react-hook-form"
import { z } from "zod"

import { Button } from "@/components/ui/button"
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
      <div className="rounded-lg border border-border/30 px-5 py-6">
        <div className="mb-5 text-center">
          <h1 className="text-xl font-semibold tracking-tight">{title}</h1>
          {description ? (
            <p className="mt-1 text-sm text-muted-foreground">{description}</p>
          ) : null}
        </div>
        <form
          onSubmit={form.handleSubmit(async (values) => {
            setServerError(null)
            const result = await onSubmitValues(values)
            if (!result.ok) setServerError(result.error)
          })}
          className="flex flex-col gap-4"
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
                className="h-10 rounded-md"
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
                className="h-10 rounded-md"
                aria-invalid={!!form.formState.errors.password}
                {...form.register("password")}
              />
              <FieldError errors={[form.formState.errors.password]} />
            </Field>
            <Field>
              <Button type="submit" className="h-10 rounded-md" disabled={form.formState.isSubmitting}>
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
      </div>
    </div>
  )
}
