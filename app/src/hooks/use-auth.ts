import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"

import { queryKeys } from "@/lib/query-keys"
import { api } from "@/lib/api"
import { clearAccessToken, getAuthHeaders, setAccessToken } from "@/lib/auth"

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

export type User = {
  id: string
  name: string
  handler: string
  email: string
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

export function useUser() {
  return useQuery({
    queryKey: queryKeys.user(),
    queryFn: async () => {
      const response = await api.get<User>("/profile", {
        headers: getAuthHeaders(),
      })
      return response
    },
  })
}

export function useLogout() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async () => {
      await api.post("/logout", undefined, {
        headers: getAuthHeaders(),
      })
    },
    onSuccess: () => {
      clearAccessToken()
      queryClient.clear()
    },
  })
}
