import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { InputGroup, InputGroupAddon, InputGroupButton, InputGroupInput } from "@/components/ui/input-group";
import { Item, ItemContent, ItemDescription, ItemTitle } from "@/components/ui/item";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Separator } from "@base-ui/react";
import { CalendarIcon, File, SearchIcon, X } from "lucide-react";
import { useState } from "react";

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

function isValidDate(date: Date | undefined) {
  if (!date) {
    return false
  }
  return !isNaN(date.getTime())
}

export function NewApprovalRoom() {
  const [open, setOpen] = useState(false)
  const [date, setDate] = useState<Date | undefined>(
    new Date("2025-06-01")
  )

  const [month, setMonth] = useState<Date | undefined>(date)
  const [value, setValue] = useState(formatDate(date))

  const [files, setFiles] = useState<File[]>([])
  const [selectedApprovers, setSelectedApprovers] = useState<{ id: string; username: string; email: string }[]>([])
  const [searchQuery, setSearchQuery] = useState("")
  const [searchResults] = useState<{ id: string; username: string; email: string }[]>([
    { id: "1", username: "john_doe", email: "john@example.com" },
    { id: "2", username: "jane_smith", email: "jane@example.com" },
    { id: "3", username: "bob_wilson", email: "bob@example.com" },
  ])

  const filteredResults = searchQuery
    ? searchResults.filter(
      (r) =>
        r.username.toLowerCase().includes(searchQuery.toLowerCase()) &&
        !selectedApprovers.find((a) => a.id === r.id)
    )
    : []

  const handleSelectApprover = (approver: { id: string; username: string; email: string }) => {
    if (!selectedApprovers.find((a) => a.id === approver.id)) {
      setSelectedApprovers((prev) => [...prev, approver])
    }
    setSearchQuery("")
  }

  const handleRemoveApprover = (id: string) => {
    setSelectedApprovers((prev) => prev.filter((a) => a.id !== id))
  }

  return (
    <div className="m-7 shadow-sm">
      <form className="min-h-full bg-background p-6 space-y-8">
        <FieldGroup className="grid grid-cols-2">
          <Field>
            <FieldLabel htmlFor="title">Title</FieldLabel>
            <Input id="title" type="text" placeholder="Title Proposal" />
          </Field>
          <Field>
            <FieldLabel htmlFor="due_date">Due Date</FieldLabel>
            <InputGroup >
              <InputGroupInput
                id="due_date"
                value={value}
                placeholder="June 01, 2025"
                onChange={(e) => {
                  const date = new Date(e.target.value)
                  setValue(e.target.value)
                  if (isValidDate(date)) {
                    setDate(date)
                    setMonth(date)
                  }
                }}
                onKeyDown={(e) => {
                  if (e.key === "ArrowDown") {
                    e.preventDefault()
                    setOpen(true)
                  }
                }} />
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
                accept=".pdf/.docs"
                className="h-10 rounded-md"
                onChange={(e) => {
                  if (e.target.files) {
                    setFiles((prev) => [...prev, ...Array.from(e.target.files!)]);
                  }
                }}
              />
            </Field>
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
                  <div className="absolute left-0 right-0 border rounded-lg mt-1 bg-popover shadow-md z-50">
                    {filteredResults.map((result) => (
                      <div
                        key={result.id}
                        className="cursor-pointer hover:bg-muted p-2"
                        onClick={() => handleSelectApprover(result)}
                      >
                        <Item>
                          <ItemContent>
                            <ItemTitle className="truncate">{result.username}</ItemTitle>
                            <ItemDescription className="truncate">{result.email}</ItemDescription>
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
                      {approver.username}
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
            </div>
          </div>
        </FieldGroup>

        <Field orientation={"horizontal"} className="mt-16 flex justify-end">
          <Button className={"p-5"}>Create Approval Room</Button>
        </Field>
      </form>
    </div>
  )
}
