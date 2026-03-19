import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Item, ItemContent, ItemDescription, ItemGroup, ItemMedia, ItemTitle } from "@/components/ui/item";
import { Progress } from "@/components/ui/progress";
import { Calendar, ChartColumnIncreasing, File, Paperclip, TrendingUp, Users } from "lucide-react";
import { cn } from "@/lib/utils";
import { Empty, EmptyHeader, EmptyMedia, EmptyTitle } from "@/components/ui/empty";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router";
import { api, getFileUrl } from "@/lib/api";
import { getAuthHeaders } from "@/lib/auth";

interface Approver {
  name: string;
  handle: string;
  decision: string;
}

interface Document {
  link: string;
  display_file_name: string;
  size: number;
}

interface ApprovalRoom {
  title: string;
  created_at: string;
  submitter_handle: string;
  due_at: string;
  documents: Document[];
  approvers: Approver[];
  aggregates: {
    file_uploaded: number;
    approval_progress: number;
  };
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return "0 Bytes";
  const k = 1024;
  const sizes = ["Bytes", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
}

export function ApprovalRoomDetailPage() {
  const { id } = useParams();

  const { data: approvalRoom, isLoading, error } = useQuery({
    queryKey: ["approval-room", id],
    queryFn: async () => {
      const response = await api.get<ApprovalRoom>(`/approval-room/${id}`, {
        headers: getAuthHeaders(),
      });
      return response;
    },
    enabled: !!id,
  });

  if (isLoading) {
    return (
      <div className="m-7 space-y-6">
        <div className="space-y-1">
          <div className="h-8 w-96 bg-muted animate-pulse rounded" />
          <div className="h-5 w-64 bg-muted animate-pulse rounded" />
        </div>
      </div>
    );
  }

  if (error || !approvalRoom) {
    return (
      <div className="m-7 space-y-6">
        <p className="text-red-500">Failed to load approval room details</p>
      </div>
    );
  }

  const { title, created_at, submitter_handle, due_at, documents, approvers, aggregates } = approvalRoom;

  return (
    <div className="m-7 space-y-6">
      <div className="space-y-1">
        <h1 className="text-2xl font-semibold tracking-tight">{title}</h1>
        <p className="text-sm text-muted-foreground">
          Created on {new Date(created_at).toDateString()} by {["@", submitter_handle].join("")}
        </p>
      </div>

      <Card className="shadow-sm">
        <CardContent className="pt-6 space-y-6">

          <div className="flex gap-4">
            <Card className="flex-1">
              <CardHeader>
                <CardTitle className="flex items-end gap-x-3">
                  <Calendar />
                  <div className="font-medium">Due Date</div>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-xl font-semibold">{new Date(due_at).toDateString()}</div>
              </CardContent>
            </Card>

            <Card className="flex-1">
              <CardHeader>
                <CardTitle className="flex items-end gap-x-3">
                  <Paperclip />
                  <div className="font-medium">File(s) Uploaded</div>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-xl font-semibold">{aggregates.file_uploaded} File(s)</div>
              </CardContent>
            </Card>

            <Card className="flex-1">
              <CardHeader>
                <CardTitle className="flex items-end gap-x-3">
                  <ChartColumnIncreasing />
                  <div className="font-medium">Approval Progress</div>
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-2">
                <div className="flex gap-x-4 justify-between items-end">
                  <div className="text-xl font-semibold">{aggregates.approval_progress}%</div>
                  <div className="font-semibold text-green-500 flex gap-x-1">
                    <TrendingUp />
                    {aggregates.approval_progress}%
                  </div>
                </div>
                <Progress value={aggregates.approval_progress} />
              </CardContent>
            </Card>
          </div>

          <div className="mt-8 flex gap-x-3">
            <div className="flex-1">
              <h2 className="text-xl font-medium mb-3">Documents for Review</h2>
              <div className="border rounded-lg p-4">
                <ItemGroup>
                  {documents.map((doc, index) => (
                    <Item key={index} render={<a href={getFileUrl(doc.link)} />}>
                      <ItemMedia variant="icon">
                        <File className="w-12 h-12" />
                      </ItemMedia>
                      <ItemContent>
                        <ItemTitle>{doc.display_file_name}</ItemTitle>
                        <ItemDescription>{formatFileSize(doc.size)}</ItemDescription>
                      </ItemContent>
                    </Item>
                  ))}
                </ItemGroup>
              </div>
            </div>
            <div className="flex-1">
              <h2 className="text-xl font-medium mb-3">Approvers Status</h2>
              <div className="border rounded-lg p-4">
                {approvers.length === 0 ? (
                  <Empty>
                    <EmptyHeader>
                      <EmptyMedia>
                        <Users />
                      </EmptyMedia>
                      <EmptyTitle>No Approvers yet</EmptyTitle>
                    </EmptyHeader>
                  </Empty>
                ) : (
                  <div className="space-y-3">
                    {approvers.map((approver, index) => (
                      <ItemGroup key={index}>
                        <Item className="gap-3">
                          <div className="flex h-10 w-10 items-center justify-center rounded-full bg-primary text-lg font-medium text-primary-foreground">
                            {approver.name.charAt(0)}
                          </div>
                          <ItemContent className="gap-0">
                            <ItemTitle className="text-base font-semibold">{approver.name}</ItemTitle>
                            <ItemDescription className="flex items-center justify-between gap-3">
                              <span>{["@", approver.handle].join("")}</span>
                              <Badge
                                className={cn(
                                  "capitalize px-3 py-1 text-sm font-medium",
                                  approver.decision === "approved" && "bg-green-100 text-green-700 border-green-200 hover:bg-green-100",
                                  approver.decision === "rejected" && "bg-red-100 text-red-700 border-red-200 hover:bg-red-100",
                                  approver.decision === "pending" && "bg-yellow-100 text-yellow-700 border-yellow-200 hover:bg-yellow-100"
                                )}
                              >
                                {approver.decision}
                              </Badge>
                            </ItemDescription>
                          </ItemContent>
                        </Item>
                      </ItemGroup>
                    ))}
                  </div>
                )}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
