import { useState } from "react";
import { AlertTriangle, Shield, Smartphone } from "lucide-react";
import { toast } from "sonner";
import { apiClient } from "@/api/client";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { PageHeader } from "@/components/layout/PageHeader";

export function SecurityPage() {
  const [isRevoking, setIsRevoking] = useState(false);

  const handleRevokeSessions = async () => {
    setIsRevoking(true);
    try {
      await apiClient.revokeAllSessions();
      toast.success("All sessions revoked successfully.");
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      toast.error(err?.response?.data?.error || "Failed to revoke sessions.");
    } finally {
      setIsRevoking(false);
    }
  };

  return (
    <div className="px-6 py-8 lg:px-10">
      <PageHeader title="Security Settings" icon={Shield} />

      <Card className="mx-auto max-w-2xl overflow-hidden rounded-2xl bg-card-solid">
        <CardContent className="p-0">
          <div className="flex gap-4 border-b border-border p-6">
            <div className="rounded-xl bg-muted p-2">
              <Smartphone className="h-6 w-6" />
            </div>
            <div>
              <h3 className="text-lg font-semibold">Active sessions</h3>
              <p className="mt-1 text-sm text-muted-foreground">
                Revoke all sessions to log out every device using your account.
              </p>
            </div>
          </div>
          <div className="flex flex-col items-center justify-between gap-4 bg-muted/10 p-6 sm:flex-row">
            <div className="flex items-center gap-2 text-destructive">
              <AlertTriangle className="h-4 w-4" />
              <span className="text-sm font-medium">This action cannot be undone.</span>
            </div>
            <Button
              variant="destructive"
              onClick={handleRevokeSessions}
              disabled={isRevoking}
            >
              {isRevoking ? "Revoking…" : "Revoke all sessions"}
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
