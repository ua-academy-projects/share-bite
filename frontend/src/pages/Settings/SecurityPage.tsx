import React, { useState } from 'react';
import { Shield, Smartphone, AlertTriangle } from 'lucide-react';
import { Button } from '../../components/ui/button';
import { apiClient } from '../../api/client';
import { toast } from 'sonner';

export const SecurityPage: React.FC = () => {
  const [isRevoking, setIsRevoking] = useState(false);

  const handleRevokeSessions = async () => {
    setIsRevoking(true);
    try {
      await apiClient.revokeAllSessions();
      toast.success("All sessions revoked successfully. You will be logged out of other devices.");
    } catch (error: any) {
      toast.error(error?.response?.data?.error || "Failed to revoke sessions.");
    } finally {
      setIsRevoking(false);
    }
  };

  return (
    <div className="flex flex-col items-center w-full min-h-screen bg-background pt-8 pb-16 px-4">
      <div className="max-w-2xl w-full">
        <header className="mb-8 flex items-center gap-3">
          <div className="p-3 bg-primary/10 text-primary rounded-full">
            <Shield size={24} />
          </div>
          <h1 className="text-3xl font-bold tracking-tight text-foreground">Security Settings</h1>
        </header>

        <div className="bg-card rounded-xl border border-border shadow-sm overflow-hidden">
          <div className="p-6 border-b border-border">
            <div className="flex items-start gap-4">
              <div className="p-2 bg-muted rounded-lg">
                <Smartphone size={24} className="text-foreground" />
              </div>
              <div className="flex-1">
                <h3 className="text-lg font-semibold text-foreground">Active Sessions</h3>
                <p className="text-sm text-muted-foreground mt-1">
                  Manage the devices that are currently logged into your account. Revoking all sessions will immediately log out all devices except this one (if supported by backend) or force a full re-login.
                </p>
              </div>
            </div>
          </div>
          <div className="p-6 bg-muted/10 flex flex-col sm:flex-row items-center justify-between gap-4">
            <div className="flex items-center gap-2 text-destructive">
              <AlertTriangle size={18} />
              <span className="text-sm font-medium">This action cannot be undone.</span>
            </div>
            <Button 
              variant="destructive" 
              onClick={handleRevokeSessions}
              disabled={isRevoking}
            >
              {isRevoking ? "Revoking..." : "Revoke All Sessions"}
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
};
