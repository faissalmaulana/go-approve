import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Item, ItemContent, ItemDescription, ItemGroup, ItemMedia, ItemTitle } from "@/components/ui/item";
import { Progress } from "@/components/ui/progress";
import { Calendar, ChartColumnIncreasing, File, Paperclip, TrendingUp } from "lucide-react";

export function ApprovalRoomDetailPage() {
  return (
    <div className="m-7 shadow-sm">
      <div className="bg-background p-6 space-y-6">
        <div className="space-y-2">
          <h2 className="text-3xl font-bold tracking-tight">Q4 Marketing Request</h2>
          <div className="flex items-center gap-x-4 text-muted-foreground">
            <p>Created On {new Date().toDateString()}</p>
            <span className="text-border">|</span>
            <p>Submitting by You</p>
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
              <div className="text-xl font-semibold">{new Date().toDateString()}</div>
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
              <div className="text-xl font-semibold">3 File(s)</div>
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
                <div className="text-xl font-semibold">50%</div>
                <div className="font-semibold text-green-500 flex gap-x-1">
                  <TrendingUp />
                  50%
                </div>
              </div>
              <Progress value={50} />
            </CardContent>
          </Card>
        </div>

        <div className="mt-8 flex gap-x-3">
          <div className="flex-1">
            <h2 className="text-xl font-medium mb-3">Documents for Review</h2>
            <div className="border-2 rounded-lg p-4">
              <ItemGroup>
                <Item render={<a href="/dashboard" />}>
                  <ItemMedia variant="icon">
                    <File className="w-12 h-12" />
                  </ItemMedia>
                  <ItemContent>
                    <ItemTitle>Documents.pdf</ItemTitle>
                    <ItemDescription>2.4MB</ItemDescription>
                  </ItemContent>
                </Item>
                <Item render={<a href="/dashboard" />}>
                  <ItemMedia variant="icon">
                    <File className="w-12 h-12" />
                  </ItemMedia>
                  <ItemContent>
                    <ItemTitle>Documents.pdf</ItemTitle>
                    <ItemDescription>2.4MB</ItemDescription>
                  </ItemContent>
                </Item>
                <Item render={<a href="/dashboard" />}>
                  <ItemMedia variant="icon">
                    <File className="w-12 h-12" />
                  </ItemMedia>
                  <ItemContent>
                    <ItemTitle>Documents.pdf</ItemTitle>
                    <ItemDescription>2.4MB</ItemDescription>
                  </ItemContent>
                </Item>
              </ItemGroup>
            </div>
          </div>
          <div className="flex-1">
            <h2 className="text-xl font-medium mb-3">Approvers Status</h2>
            <div className="border-2 rounded-lg p-4">
              <div className="space-y-3">
                <ItemGroup>
                  <Item className="gap-3">
                    <div className="flex h-10 w-10 items-center justify-center rounded-full bg-primary text-lg font-medium text-primary-foreground">
                      L
                    </div>
                    <ItemContent className="gap-0">
                      <ItemTitle className="text-base font-semibold">Lizzy MCalpine</ItemTitle>
                      <ItemDescription className="flex items-center justify-between gap-3">
                        <span>@lizzymcalpine</span>
                        <Badge variant="secondary">Pending</Badge>
                      </ItemDescription>
                    </ItemContent>
                  </Item>
                </ItemGroup>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
