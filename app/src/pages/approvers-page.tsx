import type { ReactNode } from "react"
import { useMemo } from "react"
import { useQuery } from "@tanstack/react-query"
import { Link, useSearchParams } from "react-router"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { api, ApiError } from "@/lib/api"
import { cn } from "@/lib/utils"
import { getAuthHeaders } from "@/lib/auth"

type SortField = "due_at" | "created_at"
type SortOrder = "asc" | "desc"
type ApprovalDecision = "pending" | "approved" | "rejected"

function FilterButton({
  active,
  children,
  onClick,
}: {
  active: boolean
  children: ReactNode
  onClick: () => void
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(
        "px-3 py-1.5 text-sm font-medium text-muted-foreground transition-colors",
        "hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background",
        active && "rounded-md bg-background text-foreground shadow-sm ring-1 ring-foreground/10"
      )}
    >
      {children}
    </button>
  )
}

function statusBadgeClassName(decision: ApprovalDecision) {
  if (decision === "approved") return "bg-green-100 text-green-700 border-green-200"
  if (decision === "rejected") return "bg-red-100 text-red-700 border-red-200"
  return "bg-yellow-100 text-yellow-700 border-yellow-200"
}

function statusLabel(decision: ApprovalDecision) {
  return decision.charAt(0).toUpperCase() + decision.slice(1)
}

function formatDateTime(iso: string) {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleDateString(undefined, {
    year: "numeric",
    month: "short",
    day: "2-digit",
  })
}

export function ApproversPage() {
  const [searchParams, setSearchParams] = useSearchParams()

  const sort = (searchParams.get("sort") as SortField | null) ?? "created_at"
  const order = (searchParams.get("order") as SortOrder | null) ?? "desc"
  const limit = Math.min(100, Math.max(1, Number(searchParams.get("limit") ?? "10") || 10))
  const offset = Math.max(0, Number(searchParams.get("offset") ?? "0") || 0)

  const setSortInUrl = (nextSort: SortField) => {
    const nextParams = new URLSearchParams(searchParams)
    nextParams.set("sort", nextSort)
    nextParams.set("order", order)
    nextParams.set("offset", "0")
    if (!nextParams.get("limit")) nextParams.set("limit", String(limit))
    setSearchParams(nextParams)
  }

  const setOrderInUrl = (nextOrder: SortOrder) => {
    const nextParams = new URLSearchParams(searchParams)
    nextParams.set("order", nextOrder)
    nextParams.set("sort", sort)
    nextParams.set("offset", "0")
    if (!nextParams.get("limit")) nextParams.set("limit", String(limit))
    setSearchParams(nextParams)
  }

  const setOffsetInUrl = (nextOffset: number) => {
    const nextParams = new URLSearchParams(searchParams)
    nextParams.set("offset", String(Math.max(0, nextOffset)))
    if (!nextParams.get("limit")) nextParams.set("limit", String(limit))
    setSearchParams(nextParams)
  }

  type ApiRow = {
    id: string
    title: string
    due_at: string
    created_at: string
    created_by: string
    decision: ApprovalDecision
  }

  const { data: rooms = [], isLoading, error } = useQuery({
    queryKey: ["invited-approvals", sort, order, limit, offset],
    queryFn: async () => {
      const params = new URLSearchParams({
        sort,
        order,
        limit: String(limit),
        offset: String(offset),
      })
      return api.get<ApiRow[]>(`/approval-room/approvers?${params.toString()}`, {
        headers: getAuthHeaders(),
      })
    },
  })

  const canPrevious = offset > 0
  const canNext = rooms.length === limit

  const computedShowingStart = useMemo(() => {
    return rooms.length === 0 ? 0 : offset + 1
  }, [rooms.length, offset])

  const computedShowingEnd = useMemo(() => {
    return rooms.length === 0 ? 0 : offset + rooms.length
  }, [rooms.length, offset])

  return (
    <main className="m-7 space-y-6">
      <div className="space-y-1">
        <h1 className="text-2xl font-semibold tracking-tight">Approvers</h1>
        <p className="text-sm text-muted-foreground">
          Approval rooms you were invited to review.
        </p>
      </div>

      <Card className="shadow-sm">
        <CardHeader className="border-b">
          <CardTitle className="text-base">Invited Approvals</CardTitle>
          <div className="mt-3 flex items-center justify-between gap-4">
            <div className="inline-flex items-center rounded-lg border bg-muted/30 p-1">
              <FilterButton active={sort === "due_at"} onClick={() => setSortInUrl("due_at")}>
                Due At
              </FilterButton>
              <FilterButton active={sort === "created_at"} onClick={() => setSortInUrl("created_at")}>
                Created At
              </FilterButton>
            </div>

            <div className="inline-flex items-center rounded-lg border bg-muted/30 p-1">
              <FilterButton active={order === "desc"} onClick={() => setOrderInUrl("desc")}>
                Latest
              </FilterButton>
              <FilterButton active={order === "asc"} onClick={() => setOrderInUrl("asc")}>
                Oldest
              </FilterButton>
            </div>
          </div>
        </CardHeader>

        <CardContent className="pt-4">
          {isLoading && (
            <div className="pb-4 text-sm text-muted-foreground">Loading...</div>
          )}
          {error && (
            <div className="pb-4 text-sm text-red-600">
              {error instanceof ApiError ? error.message : "Failed to load invitations"}
            </div>
          )}

          <div className="rounded-lg border">
            <Table>
              <TableHeader>
                <TableRow className="bg-muted/30 hover:bg-muted/30">
                  <TableHead className="px-4">TITLE</TableHead>
                  <TableHead>DUE AT</TableHead>
                  <TableHead>CREATED AT</TableHead>
                  <TableHead>CREATED BY</TableHead>
                  <TableHead>DECISION</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {rooms.map((room) => (
                  <TableRow key={room.id}>
                    <TableCell className="px-4">
                      <Link
                        to={`/approvers/${room.id}`}
                        className="block max-w-[260px] truncate font-medium text-primary hover:underline"
                        title={room.title}
                      >
                        {room.title}
                      </Link>
                    </TableCell>
                    <TableCell className="text-muted-foreground">{formatDateTime(room.due_at)}</TableCell>
                    <TableCell className="text-muted-foreground">{formatDateTime(room.created_at)}</TableCell>
                    <TableCell className="text-muted-foreground">
                      {["@", room.created_by].join("")}
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant="outline"
                        className={cn("capitalize px-3 py-1", statusBadgeClassName(room.decision))}
                      >
                        {statusLabel(room.decision)}
                      </Badge>
                    </TableCell>
                  </TableRow>
                ))}

                {!isLoading && rooms.length === 0 && (
                  <TableRow>
                    <TableCell colSpan={5} className="py-10 text-center">
                      <p className="text-sm text-muted-foreground">
                        No invited approvals found.
                      </p>
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>

          <div className="mt-4 flex items-center justify-between text-sm text-muted-foreground">
            <div>Showing {computedShowingStart}-{computedShowingEnd} requests</div>
            <div className="flex items-center gap-2">
              <Button
                size="sm"
                variant="outline"
                disabled={!canPrevious || isLoading}
                onClick={() => setOffsetInUrl(offset - limit)}
              >
                Previous
              </Button>
              <Button
                size="sm"
                variant="outline"
                disabled={!canNext || isLoading}
                onClick={() => setOffsetInUrl(offset + limit)}
              >
                Next
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </main>
  )
}

