export const queryKeys = {
  all: ["auth"] as const,
  user: () => [...queryKeys.all, "user"] as const,
}
