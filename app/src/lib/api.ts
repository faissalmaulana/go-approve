const API_BASE_URL = "http://localhost:8080"

export type ApiEnvelope<T> = {
  data: T | null
  error: string | null
}

async function parseJsonSafely(response: Response): Promise<unknown> {
  const text = await response.text()
  if (!text) return null
  try {
    return JSON.parse(text) as unknown
  } catch {
    return null
  }
}

export async function postJson<TReq extends Record<string, unknown>, TRes>(
  path: string,
  body: TReq
): Promise<ApiEnvelope<TRes>> {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  })

  const json = (await parseJsonSafely(res)) as ApiEnvelope<TRes> | null
  if (json && typeof json === "object" && "data" in json && "error" in json) {
    return json
  }

  if (!res.ok) {
    return { data: null, error: `Request failed (${res.status})` }
  }

  return { data: null, error: "Unexpected server response." }
}

