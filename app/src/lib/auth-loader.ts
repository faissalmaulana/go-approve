import { redirect } from "react-router"

import { getAccessToken } from "@/lib/auth"

export function requireAuth() {
  const token = getAccessToken()
  if (!token) {
    throw redirect("/login")
  }
  return token
}

export function requireGuest() {
  const token = getAccessToken()
  if (token) {
    throw redirect("/")
  }
  return null
}
