import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from "@/components/ui/alert-dialog";
import { Item, ItemContent, ItemDescription, ItemGroup, ItemMedia, ItemTitle } from "@/components/ui/item";
import { Progress } from "@/components/ui/progress";
import { Calendar, ChartColumnIncreasing, File, Paperclip, TrendingUp, Users } from "lucide-react";
import { cn } from "@/lib/utils";
import { Empty, EmptyHeader, EmptyMedia, EmptyTitle } from "@/components/ui/empty";
import { buttonVariants } from "@/components/ui/button-variants";
import { getAuthHeaders } from "@/lib/auth";
import { api } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router";
import formatFileSize from "@/utils/formatfilesize";
import { useUser } from "@/hooks/use-auth";

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


export function ApproverRoomDetailPage() {
  const { id } = useParams();

  const { data: currentUser, isLoading: isLoadingCurrentUser } = useUser()

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

  if (isLoading && isLoadingCurrentUser) {
    return (
      <div className="m-7 shadow-sm">
        <div className="bg-background p-6 space-y-6">
          <div className="space-y-2">
            <div className="h-9 w-96 bg-muted animate-pulse rounded" />
            <div className="h-5 w-64 bg-muted animate-pulse rounded" />
          </div>
        </div>
      </div>
    );
  }

  if (error || !approvalRoom || !currentUser) {
    return (
      <div className="m-7 shadow-sm">
        <div className="bg-background p-6">
          <p className="text-red-500">Failed to load approval room details</p>
        </div>
      </div>
    );
  }


  const { title, created_at, submitter_handle, due_at, documents, approvers, aggregates } = approvalRoom;

  return (
    <div className="m-7 shadow-sm">
      <div className="bg-background p-6 space-y-6">
        <div className="space-y-2">
          <h2 className="text-3xl font-bold tracking-tight">{title}</h2>
          <div className="flex items-center gap-x-4 text-muted-foreground">
            <p>Created On {new Date(created_at).toDateString()}</p>
            <span className="text-border">|</span>
            <p>Submitting by {["@", submitter_handle].join("")}</p>
          </div>
        </div>

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
            <h2 className="text-xl font-medium mb-3">Documents to reviews</h2>
            <div className="border-2 rounded-lg p-4">
              <ItemGroup>
                {documents.map((doc, index) => (
                  <Item key={index} render={<a href={doc.link} />}>
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
            <h2 className="text-xl font-medium mb-3">Other's approvers status</h2>
            <div className="border-2 rounded-lg p-4">
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
                  {approvers.filter((approver) => approver.handle !== currentUser.handler).map((approver, index) => (
                    <ItemGroup key={index}>
                      <Item className="gap-3">
                        <div className="flex h-10 w-10 items-center justify-center rounded-full bg-primary text-lg font-medium text-primary-foreground">
                          {approver.name.charAt(0)}
                        </div>
                        <ItemContent className="gap-0">
                          <ItemTitle className="text-base font-semibold">{approver.name}</ItemTitle>
                          <ItemDescription className="flex items-center justify-between gap-3">
                            <span>{approver.handle}</span>
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
        <div className="w-full grid justify-end">
          <h2 className="text-xl font-medium mb-3">Submit Approver or Reject</h2>
          <div className="border-2 rounded-lg p-4 flex gap-x-4">
            <AlertDialog>
              <AlertDialogTrigger>
                <span className={cn(buttonVariants({ variant: "default" }), "p-5")}>
                  Approve
                </span>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Confirm Approval</AlertDialogTitle>
                  <AlertDialogDescription>
                    Are you sure you want to approve this request?
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>Cancel</AlertDialogCancel>
                  <AlertDialogAction onClick={() => console.log("Approved")}>Confirm</AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
            <AlertDialog>
              <AlertDialogTrigger>
                <span className={cn(buttonVariants({ variant: "destructive" }), "p-5")}>
                  Reject
                </span>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Confirm Rejection</AlertDialogTitle>
                  <AlertDialogDescription>
                    Are you sure you want to reject this request?
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>Cancel</AlertDialogCancel>
                  <AlertDialogAction onClick={() => console.log("Rejected")}>Confirm</AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </div>
        </div>
      </div>
    </div>
  )
}
