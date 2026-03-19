import type { ReactNode } from "react"
import { useEffect, useMemo, useState } from "react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { cn } from "@/lib/utils"
import { api, ApiError } from "@/lib/api"
import { getAuthHeaders } from "@/lib/auth"

type InvitationStatus = "pending" | "accepted" | "rejected"

type Invitation = {
  id: string
  roomId: string
  requestedBy: {
    name?: string
    handle: string
  }
  createdAt: string
  status: InvitationStatus
}

type Filter = "all" | InvitationStatus

function formatDateTime(iso: string) {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleString(undefined, {
    month: "short",
    day: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  })
}

function statusBadgeClassName(status: InvitationStatus) {
  if (status === "accepted") return "bg-green-100 text-green-700 border-green-200"
  if (status === "rejected") return "bg-red-100 text-red-700 border-red-200"
  return "bg-yellow-100 text-yellow-700 border-yellow-200"
}

function statusLabel(status: InvitationStatus) {
  return status.charAt(0).toUpperCase() + status.slice(1)
}


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
        active &&
        "rounded-md bg-background text-foreground shadow-sm ring-1 ring-foreground/10"
      )}
    >
      {children}
    </button>
  )
}

export function InvitationsReceivedPage() {
  const [filter, setFilter] = useState<Filter>("all")
  const queryClient = useQueryClient()

  const limit = 10
  const [offset, setOffset] = useState(0)

  type ApiInvitation = {
    id: string
    room_id: string
    status: InvitationStatus
    created_at: string
    requested_by: {
      id: string
      name: string
      handler: string
    }
  }

  useEffect(() => {
    setOffset(0)
  }, [filter])

  const { data: items = [], isLoading, error } = useQuery({
    queryKey: ["request-review", "received", filter, limit, offset],
    queryFn: async () => {
      const params = new URLSearchParams({ as: "received" })
      if (filter !== "all") params.set("status", filter)
      params.set("limit", String(limit))
      params.set("offset", String(offset))

      const res = await api.get<ApiInvitation[]>(`/request-review?${params.toString()}`, {
        headers: getAuthHeaders(),
      })

      return res.map<Invitation>((r) => ({
        id: r.id,
        roomId: r.room_id,
        createdAt: r.created_at,
        status: r.status,
        requestedBy: {
          name: r.requested_by.name,
          handle: r.requested_by.handler,
        },
      }))
    },
  })

  const filtered = useMemo(() => {
    if (filter === "all") return items
    return items.filter((i) => i.status === filter)
  }, [filter, items])

  const showingStart = filtered.length === 0 ? 0 : offset + 1
  const showingEnd = offset + filtered.length
  const canPrevious = offset > 0
  const canNext = filtered.length === limit

  const [loadingId, setLoadingId] = useState<string | null>(null)

  const confirmMutation = useMutation({
    mutationFn: async ({ id, status }: { id: string; status: InvitationStatus }) => {
      await api.put(
        `/request-review/${id}/confirm`,
        { status },
        { headers: getAuthHeaders() }
      )
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["request-review", "received"] })
    },
  })

  const handleAccept = async (id: string) => {
    setLoadingId(id)
    try {
      await confirmMutation.mutateAsync({ id, status: "accepted" })
    } catch (err) {
      const message = err instanceof ApiError ? err.message : "Failed to accept request"
      console.error(message)
    } finally {
      setLoadingId(null)
    }
  }

  const handleReject = async (id: string) => {
    setLoadingId(id)
    try {
      await confirmMutation.mutateAsync({ id, status: "rejected" })
    } catch (err) {
      const message = err instanceof ApiError ? err.message : "Failed to reject request"
      console.error(message)
    } finally {
      setLoadingId(null)
    }
  }

  return (
    <div className="m-7 space-y-6">
      <div className="flex items-start justify-between gap-4">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold tracking-tight">
            Review Requests
          </h1>
          <p className="text-sm text-muted-foreground">
            Manage incoming invitations to review approval rooms.
          </p>
        </div>

        <div className="inline-flex items-center rounded-lg border bg-muted/30 p-1">
          <FilterButton active={filter === "all"} onClick={() => setFilter("all")}>
            All
          </FilterButton>
          <FilterButton
            active={filter === "pending"}
            onClick={() => setFilter("pending")}
          >
            Pending
          </FilterButton>
          <FilterButton
            active={filter === "accepted"}
            onClick={() => setFilter("accepted")}
          >
            Accepted
          </FilterButton>
          <FilterButton
            active={filter === "rejected"}
            onClick={() => setFilter("rejected")}
          >
            Rejected
          </FilterButton>
        </div>
      </div>

      <Card className="shadow-sm">
        <CardHeader className="border-b">
          <CardTitle className="text-base">Requests</CardTitle>
          <CardDescription>Incoming invitations awaiting your action.</CardDescription>
        </CardHeader>
        <CardContent className="pt-4">
          {isLoading && (
            <div className="pb-4 text-sm text-muted-foreground">Loading...</div>
          )}
          {error && (
            <div className="pb-4 text-sm text-red-600">
              Failed to load invitations
            </div>
          )}
          <div className="rounded-lg border">
            <Table>
              <TableHeader>
                <TableRow className="bg-muted/30 hover:bg-muted/30">
                  <TableHead className="px-4">STATUS</TableHead>
                  <TableHead>ROOM ID</TableHead>
                  <TableHead>REQUESTED BY</TableHead>
                  <TableHead>DATE CREATED</TableHead>
                  <TableHead className="text-right pr-4">ACTIONS</TableHead>
                </TableRow>
              </TableHeader>

              <TableBody>
                {filtered.map((inv) => {
                  const isDone = inv.status !== "pending"

                  return (
                    <TableRow key={inv.id}>
                      <TableCell className="px-4">
                        <Badge
                          variant="outline"
                          className={cn(
                            "capitalize px-3 py-1",
                            statusBadgeClassName(inv.status)
                          )}
                        >
                          <span
                            className={cn(
                              "mr-2 inline-block h-1.5 w-1.5 rounded-full",
                              inv.status === "pending" && "bg-yellow-500",
                              inv.status === "accepted" && "bg-green-500",
                              inv.status === "rejected" && "bg-red-500"
                            )}
                          />
                          {statusLabel(inv.status)}
                        </Badge>
                      </TableCell>

                      <TableCell className="font-medium text-muted-foreground">
                        {inv.roomId}
                      </TableCell>

                      <TableCell>
                        <span className="text-sm text-foreground">
                          {['@', inv.requestedBy.handle].join("")}
                        </span>
                      </TableCell>

                      <TableCell className="text-muted-foreground">
                        {formatDateTime(inv.createdAt)}
                      </TableCell>

                      <TableCell className="text-right pr-4">
                        {isDone ? (
                          <span className="text-sm italic text-muted-foreground">
                            Action completed
                          </span>
                        ) : (
                          <div className="inline-flex items-center gap-2">
                            <Button
                              size="sm"
                              onClick={() => handleAccept(inv.id)}
                              disabled={loadingId !== null}
                            >
                              {loadingId === inv.id ? "Loading..." : "Accept"}
                            </Button>
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => handleReject(inv.id)}
                              disabled={loadingId !== null}
                            >
                              {loadingId === inv.id ? "Loading..." : "Reject"}
                            </Button>
                          </div>
                        )}
                      </TableCell>
                    </TableRow>
                  )
                })}

                {filtered.length === 0 && (
                  <TableRow>
                    <TableCell colSpan={5} className="py-10 text-center">
                      <p className="text-sm text-muted-foreground">
                        No requests found.
                      </p>
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>

          <div className="mt-4 flex items-center justify-between text-sm text-muted-foreground">
            <div>
              Showing {showingStart}-{showingEnd} requests
            </div>
            <div className="flex items-center gap-2">
              <Button
                size="sm"
                variant="outline"
                disabled={!canPrevious || isLoading}
                onClick={() => setOffset((prev) => Math.max(0, prev - limit))}
              >
                Previous
              </Button>
              <Button
                size="sm"
                variant="outline"
                disabled={!canNext || isLoading}
                onClick={() => setOffset((prev) => prev + limit)}
              >
                Next
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
