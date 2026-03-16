import { zodResolver } from "@hookform/resolvers/zod"
import { useDeferredValue, useState } from "react"
import { useForm } from "react-hook-form"
import { useMutation, useQuery } from "@tanstack/react-query"
import { z } from "zod"
import { Badge } from "@/components/ui/badge"
import {
  AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent,
  AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger
} from "@/components/ui/alert-dialog"
import { buttonVariants } from "@/components/ui/button-variants"
import { cn } from "@/lib/utils"
import { Calendar } from "@/components/ui/calendar"
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { InputGroup, InputGroupAddon, InputGroupButton, InputGroupInput } from "@/components/ui/input-group"
import { Item, ItemContent, ItemDescription, ItemTitle } from "@/components/ui/item"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { Separator } from "@base-ui/react"
import { CalendarIcon, File, SearchIcon, X } from "lucide-react"
import { api } from "@/lib/api"
import { getAuthHeaders } from "@/lib/auth"
import { useNavigate } from "react-router"

type SearchUser = {
  id: string
  name: string
  handler: string
}

type CreateApprovalRoomResponse = {
  message: string
}

const approvalRoomSchema = z.object({
  title: z.string().min(1, "Title is required."),
  due_date: z.string(),
  approvers: z.array(z.string()).min(1, "At least one approver is required."),
})

type ApprovalRoomValues = z.infer<typeof approvalRoomSchema>

function formatDate(date: Date | undefined) {
  if (!date) {
    return ""
  }
  return date.toLocaleDateString("en-US", {
    day: "2-digit",
    month: "long",
    year: "numeric",
  })
}

export function NewApprovalRoom() {
  const [open, setOpen] = useState(false)
  const [alertOpen, setAlertOpen] = useState(false)
  const [serverError, setServerError] = useState<string | null>(null)
  const [date, setDate] = useState<Date | undefined>(
    new Date()
  )

  const [month, setMonth] = useState<Date | undefined>(date)
  const [value, setValue] = useState(formatDate(date))

  const [files, setFiles] = useState<File[]>([])
  const [selectedApprovers, setSelectedApprovers] = useState<{ id: string; name: string; handler: string }[]>([])
  const [searchQuery, setSearchQuery] = useState("")
  const deferredSearchQuery = useDeferredValue(searchQuery)

  const navigate = useNavigate();

  const form = useForm<ApprovalRoomValues>({
    resolver: zodResolver(approvalRoomSchema),
    defaultValues: {
      title: "",
      due_date: "",
      approvers: [],
    },
    mode: "onSubmit",
  })

  const createApprovalRoom = useMutation({
    mutationFn: async (data: ApprovalRoomValues) => {
      const formData = new FormData()
      const jsonData = {
        title: data.title,
        due_at: data.due_date,
        approvers: data.approvers,
      }
      formData.append("json_data", JSON.stringify(jsonData))
      files.forEach((file) => {
        formData.append("documents", file)
      })

      const response = await api.post<CreateApprovalRoomResponse>("/approval-room", formData, {
        headers: {
          ...getAuthHeaders(),
        },
      })
      return response
    },
    onSuccess: () => {
      setAlertOpen(false)
      setFiles([])
      setSelectedApprovers([])
      setValue("")
      setDate(undefined)
      form.reset()

      navigate("/")
    },
    onError: (error) => {
      console.error(error)
      setServerError("Failed to create approval room")
    },
  })

  const { data: searchResults = [] } = useQuery({
    queryKey: ["search-users", deferredSearchQuery],
    queryFn: async () => {
      if (!deferredSearchQuery.trim()) return []
      const users = await api.get<SearchUser[]>(`/search-users?handle=${encodeURIComponent(deferredSearchQuery)}`, {
        headers: getAuthHeaders(),
      })
      return users
    },
    enabled: deferredSearchQuery.trim().length > 0,
  })

  const filteredResults = searchResults.filter((r) => !selectedApprovers.find((a) => a.id === r.id))

  const handleSelectApprover = (approver: { id: string; name: string; handler: string }) => {
    if (!selectedApprovers.find((a) => a.id === approver.id)) {
      const newApprovers = [...selectedApprovers, approver]
      setSelectedApprovers(newApprovers)
      form.setValue("approvers", newApprovers.map((a) => a.id), { shouldValidate: true })
    }
    setSearchQuery("")
  }

  const handleRemoveApprover = (id: string) => {
    const newApprovers = selectedApprovers.filter((a) => a.id !== id)
    setSelectedApprovers(newApprovers)
    form.setValue("approvers", newApprovers.map((a) => a.id), { shouldValidate: true })
  }

  const handleConfirmCreate = async () => {
    const isTitleValid = await form.trigger("title")
    const isApproversValid = await form.trigger("approvers")

    if (!date) {
      form.setError("due_date", { message: "Due date is required." })
      return
    }

    if (files.length === 0) {
      setServerError("At least one document is required.")
      return
    }

    if (!isTitleValid || !isApproversValid) {
      return
    }

    setServerError(null)
    try {
      await createApprovalRoom.mutateAsync({
        title: form.getValues("title"),
        due_date: date.toISOString(),
        approvers: selectedApprovers.map((a) => a.id),
      })
    } catch (error) {
      console.error(error)
    }
  }

  return (
    <div className="m-7 shadow-sm">
      <form className="min-h-full bg-background p-6 space-y-8">
        <FieldGroup className="grid grid-cols-2">
          <Field>
            <FieldLabel htmlFor="title">Title</FieldLabel>
            <Input
              id="title"
              type="text"
              placeholder="Title Proposal"
              aria-invalid={!!form.formState.errors.title}
              {...form.register("title")}
            />
            <FieldError errors={[form.formState.errors.title]} />
          </Field>
          <Field>
            <FieldLabel htmlFor="due_date">Due Date</FieldLabel>
            <InputGroup >
              <InputGroupInput
                id="due_date"
                value={value}
                placeholder="Select date"
                onClick={() => setOpen(true)}
                onKeyDown={(e) => {
                  if (e.key === "ArrowDown" || e.key === "Enter") {
                    e.preventDefault()
                    setOpen(true)
                  }
                }}
              />
              <InputGroupAddon align="inline-end">
                <Popover open={open} onOpenChange={setOpen}>
                  <PopoverTrigger render={<InputGroupButton id="date-picker" variant="ghost" size="icon-xs" aria-label="Select date"><CalendarIcon /><span className="sr-only">Select date</span></InputGroupButton>} />
                  <PopoverContent
                    className="w-auto overflow-hidden p-0"
                    align="end"
                    alignOffset={-8}
                    sideOffset={10}
                  >
                    <Calendar
                      mode="single"
                      selected={date}
                      month={month}
                      onMonthChange={setMonth}
                      onSelect={(date) => {
                        setDate(date)
                        setValue(formatDate(date))
                        setOpen(false)
                      }}
                    />
                  </PopoverContent>
                </Popover>
              </InputGroupAddon>
            </InputGroup>
            <FieldError errors={[form.formState.errors.due_date]} />
          </Field>
        </FieldGroup>

        <FieldGroup className="grid grid-cols-2">
          <div className="space-y-4">
            <Field>
              <FieldLabel htmlFor="documents">Documents</FieldLabel>
              <Input
                id="documents"
                type="file"
                multiple
                accept=".pdf,.doc,.docx,application/pdf,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
                onChange={(e) => {
                  if (e.target.files) {
                    setFiles((prev) => [...prev, ...Array.from(e.target.files!)]);
                  }
                }}
              />
            </Field>
            {serverError && (
              <p className="text-sm text-destructive">{serverError}</p>
            )}
            <div className="space-y-2">
              {files.map((file) => (
                <div key={file.name}>
                  <div className="flex items-center justify-between gap-4">
                    <div className="flex items-center gap-2 min-w-0">
                      <File className="size-4 shrink-0 text-muted-foreground" />
                      <span className="truncate text-sm">{file.name}</span>
                    </div>
                    <div className="flex items-center gap-4 shrink-0">
                      <span className="text-sm text-muted-foreground">{(file.size / 1024).toFixed(1)} KB</span>
                      <button
                        type="button"
                        onClick={() => setFiles((prev) => prev.filter((f) => f.name !== file.name))}
                        className="rounded-full hover:bg-muted p-1"
                        aria-label="Remove file"
                      >
                        <X className="size-3.5" />
                      </button>
                    </div>
                  </div>
                  <Separator className="mt-2 border" />
                </div>
              ))}
            </div>
          </div>

          <div className="space-y-4">
            <Field>
              <FieldLabel>Search Approvers</FieldLabel>
              <div className="relative">
                <InputGroup>
                  <InputGroupInput
                    id="search-input"
                    placeholder="Search..."
                    value={searchQuery}
                    className="h-10 rounded-md"
                    onChange={(e) => setSearchQuery(e.target.value)}
                  />
                  <InputGroupAddon align="inline-start">
                    <SearchIcon className="text-muted-foreground" />
                  </InputGroupAddon>
                </InputGroup>
                {filteredResults.length > 0 && (
                  <div className="absolute left-0 right-0 border rounded-lg mt-1 bg-popover shadow-md z-50 max-h-60 overflow-y-auto">
                    {filteredResults.map((result) => (
                      <div
                        key={result.id}
                        className="cursor-pointer hover:bg-muted p-2"
                        onClick={() => handleSelectApprover(result)}
                      >
                        <Item>
                          <ItemContent>
                            <ItemTitle className="truncate">{result.name}</ItemTitle>
                            <ItemDescription className="truncate">{["@", result.handler ?? ""].join("")}</ItemDescription>
                          </ItemContent>
                        </Item>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </Field>
            <div className="space-y-2 mt-8">
              <FieldLabel>Selected Approvers</FieldLabel>
              {selectedApprovers.length > 0 ? (
                <div className="flex flex-wrap gap-2">
                  {selectedApprovers.map((approver) => (
                    <Badge key={approver.id} variant="secondary" className="p-3">
                      {approver.name}
                      <button
                        type="button"
                        onClick={() => handleRemoveApprover(approver.id)}
                        className="ml-0.5 rounded-full hover:bg-secondary/80 p-0.5"
                        aria-label="Remove approver"
                      >
                        <X data-icon="inline-end" className="size-3" />
                      </button>
                    </Badge>
                  ))}
                </div>
              ) : (
                <p className="text-sm text-muted-foreground">No approvers picked yet.</p>
              )}
              <FieldError errors={[form.formState.errors.approvers]} />
            </div>
          </div>
        </FieldGroup>

        <Field orientation={"horizontal"} className="mt-16 flex justify-end">
          <AlertDialog open={alertOpen} onOpenChange={setAlertOpen}>
            <AlertDialogTrigger>
              <span className={cn(buttonVariants({ variant: "default" }), "p-5")}>
                Create Approval Room
              </span>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Create Approval Room?</AlertDialogTitle>
                <AlertDialogDescription>
                  Are you sure you want to create a new approval room? This action cannot be undone.
                </AlertDialogDescription>
              </AlertDialogHeader>
              {serverError ? (
                <Field>
                  <FieldError>{serverError}</FieldError>
                </Field>
              ) : null}
              <AlertDialogFooter>
                <AlertDialogCancel>Cancel</AlertDialogCancel>
                <AlertDialogAction
                  onClick={handleConfirmCreate}
                  disabled={createApprovalRoom.isPending}
                >
                  {createApprovalRoom.isPending ? "Creating..." : "Create"}
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </Field>
      </form>
    </div>
  )
}
