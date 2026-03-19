import type { ReactNode } from "react"
import { useQuery } from "@tanstack/react-query"
import { Link, useSearchParams } from "react-router"
import { Button } from "@/components/ui/button"
import { buttonVariants } from "@/components/ui/button-variants"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { api, ApiError } from "@/lib/api"
import { cn } from "@/lib/utils"
import { getAuthHeaders } from "@/lib/auth"
import { Plus } from "lucide-react"

type SortField = "due_at" | "created_at"
type SortOrder = "asc" | "desc"

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

export function HomePage() {
  type ApprovalRoom = {
    id: string
    title: string
    due_at: string
    created_at: string
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

  const { data: rooms = [], isLoading, error } = useQuery({
    queryKey: ["my-approval-rooms", sort, order, limit, offset],
    queryFn: async () => {
      const params = new URLSearchParams({
        sort,
        order,
        limit: String(limit),
        offset: String(offset),
      })

      const res = await api.get<ApprovalRoom[]>(`/approval-room?${params.toString()}`, {
        headers: getAuthHeaders(),
      })

      return res
    },
  })

  const canPrevious = offset > 0
  const canNext = rooms.length === limit
  const computedShowingStart = rooms.length === 0 ? 0 : offset + 1
  const computedShowingEnd = rooms.length === 0 ? 0 : offset + rooms.length

  return (
    <main className="m-7 space-y-6">
      <div className="flex items-start justify-between gap-4">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold tracking-tight">Approvals</h1>
          <p className="text-sm text-muted-foreground">Review approval rooms you created.</p>
        </div>

        <div className="flex items-center gap-2">
          <Link to="/new-approval-room" className={buttonVariants({ size: "lg" })}>
            <Plus />
          </Link>

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
      </div>


      {isLoading && (
        <div className="pb-4 text-sm text-muted-foreground">Loading...</div>
      )}
      {error && (
        <div className="pb-4 text-sm text-red-600">
          {error instanceof ApiError
            ? error.message
            : "Failed to load approval rooms"}
        </div>
      )}

      <div className="rounded-lg border overflow-hidden">
        <Table className="bg-background">
          <TableHeader>
            <TableRow className="bg-muted/30 hover:bg-muted/30">
              <TableHead className="px-4">TITLE</TableHead>
              <TableHead>DUE AT</TableHead>
              <TableHead>CREATED AT</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {rooms.map((room) => (
              <TableRow key={room.id}>
                <TableCell className="px-4">
                  <Link
                    to={`/${room.id}`}
                    className="block max-w-[260px] truncate font-medium text-primary hover:underline"
                    title={room.title}
                  >
                    {room.title}
                  </Link>
                </TableCell>
                <TableCell className="text-muted-foreground">
                  <Link to={`/${room.id}`} className="hover:underline">
                    {formatDateTime(room.due_at)}
                  </Link>
                </TableCell>
                <TableCell className="text-muted-foreground">
                  <Link to={`/${room.id}`} className="hover:underline">
                    {formatDateTime(room.created_at)}
                  </Link>
                </TableCell>
              </TableRow>
            ))}

            {!isLoading && rooms.length === 0 && (
              <TableRow>
                <TableCell colSpan={3} className="py-10 text-center">
                  <p className="text-sm text-muted-foreground">No approval rooms found.</p>
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      <div className="mt-4 flex items-center justify-between text-sm text-muted-foreground">
        <div>
          Showing {computedShowingStart}-{computedShowingEnd} requests
        </div>
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
    </main>
  )
}
