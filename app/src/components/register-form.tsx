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

const registerSchema = z
  .object({
    name: z.string().min(1, "Name is required."),
    handler: z
      .string()
      .min(1, "Username is required.")
      .regex(/^[a-zA-Z0-9_]+$/, "Only letters, numbers, and underscore."),
    email: z.string().email("Please enter a valid email."),
    password: z.string().min(1, "Password is required."),
    confirmPassword: z.string().min(1, "Please confirm your password."),
  })
  .refine((v) => v.password === v.confirmPassword, {
    message: "Passwords do not match.",
    path: ["confirmPassword"],
  })

export type RegisterValues = z.infer<typeof registerSchema>

type SubmitResult = { ok: true } | { ok: false; error: string }

export function RegisterForm({
  className,
  title = "Create account",
  description = "Register to start using Go Approve.",
  submitLabel = "Create account",
  bottomText,
  onSubmitValues,
  ...props
}: React.ComponentProps<"div"> & {
  title?: string
  description?: React.ReactNode
  submitLabel?: string
  bottomText?: React.ReactNode
  onSubmitValues: (values: RegisterValues) => Promise<SubmitResult>
}) {
  const [serverError, setServerError] = useState<string | null>(null)

  const form = useForm<RegisterValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      name: "",
      handler: "",
      email: "",
      password: "",
      confirmPassword: "",
    },
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
                <FieldLabel htmlFor="name">Name</FieldLabel>
                <Input
                  id="name"
                  autoComplete="name"
                  aria-invalid={!!form.formState.errors.name}
                  {...form.register("name")}
                />
                <FieldError errors={[form.formState.errors.name]} />
              </Field>

              <Field>
                <FieldLabel htmlFor="handler">Username</FieldLabel>
                <Input
                  id="handler"
                  autoComplete="username"
                  aria-invalid={!!form.formState.errors.handler}
                  {...form.register("handler")}
                />
                <FieldError errors={[form.formState.errors.handler]} />
              </Field>

              <Field>
                <FieldLabel htmlFor="email">Email</FieldLabel>
                <Input
                  id="email"
                  type="email"
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
                  autoComplete="new-password"
                  aria-invalid={!!form.formState.errors.password}
                  {...form.register("password")}
                />
                <FieldError errors={[form.formState.errors.password]} />
              </Field>

              <Field>
                <FieldLabel htmlFor="confirmPassword">Confirm password</FieldLabel>
                <Input
                  id="confirmPassword"
                  type="password"
                  autoComplete="new-password"
                  aria-invalid={!!form.formState.errors.confirmPassword}
                  {...form.register("confirmPassword")}
                />
                <FieldError errors={[form.formState.errors.confirmPassword]} />
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

