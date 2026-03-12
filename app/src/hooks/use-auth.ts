import { useMutation, useQueryClient } from "@tanstack/react-query"

import { queryKeys } from "@/lib/query-keys"
import { api } from "@/lib/api"
import { setAccessToken } from "@/lib/auth"

export type SignInCredentials = {
  email: string
  password: string
}

export type SignUpCredentials = {
  email: string
  name: string
  handler: string
  password: string
}

type SignInResponse = { access_token: string }
type SignUpResponse = { access_token: string }

export function useSignIn() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (credentials: SignInCredentials) => {
      const response = await api.post<SignInResponse>("/auth/sign-in", credentials)
      return response.access_token
    },
    onSuccess: (token) => {
      setAccessToken(token)
      queryClient.invalidateQueries({ queryKey: queryKeys.user() })
    },
  })
}

export function useSignUp() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (credentials: SignUpCredentials) => {
      const response = await api.post<SignUpResponse>("/auth/sign-up", credentials)
      return response.access_token
    },
    onSuccess: (token) => {
      setAccessToken(token)
      queryClient.invalidateQueries({ queryKey: queryKeys.user() })
    },
  })
}
