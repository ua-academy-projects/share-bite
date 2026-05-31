import React, { useState } from 'react';
import { Shield, Smartphone, AlertTriangle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { apiClient } from '@/api/client';
import { toast } from 'sonner';
import { PageHeader } from '@/components/layout/PageHeader';
import { PageLayout } from '@/components/layout/PageLayout';
import { pagePanel } from '@/components/layout/pageStyles';
import { cn } from '@/lib/utils';

export const SecurityPage: React.FC = () => {
  const [isRevoking, setIsRevoking] = useState(false);

  const handleRevokeSessions = async () => {
    setIsRevoking(true);
    try {
      await apiClient.revokeAllSessions();
      toast.success("All sessions revoked successfully. You will be logged out of other devices.");
    } catch (error: unknown) {
      const e = error as { response?: { data?: { error?: string } } };
      toast.error(e?.response?.data?.error || "Failed to revoke sessions.");
    } finally {
      setIsRevoking(false);
    }
  };

  return (
    <PageLayout>
      <PageHeader
        title="Security Settings"
        description="Manage active sessions and account security"
      />

      <div className={cn(pagePanel, "mx-auto max-w-2xl overflow-hidden")}>
        <div className="border-b border-gray-200 p-6 dark:border-[#2f5e50]">
          <div className="flex items-start gap-4">
            <div className="rounded-lg bg-gray-100 p-2 dark:bg-[#0d241d]">
              <Smartphone size={24} className="text-[#1A3C34] dark:text-white" />
            </div>
            <div className="flex-1">
              <h3 className="text-lg font-semibold text-[#1A3C34] dark:text-white">
                Active Sessions
              </h3>
              <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                Manage the devices that are currently logged into your account.
                Revoking all sessions will log out all other devices.
              </p>
            </div>
            <Shield className="h-6 w-6 text-emerald-500 dark:text-[#98FF98]" />
          </div>
        </div>
        <div className="flex flex-col items-center justify-between gap-4 bg-gray-50/50 p-6 sm:flex-row dark:bg-[#0d241d]/30">
          <div className="flex items-center gap-2 text-destructive">
            <AlertTriangle size={18} />
            <span className="text-sm font-medium">This action cannot be undone.</span>
          </div>
          <Button
            variant="destructive"
            onClick={handleRevokeSessions}
            disabled={isRevoking}
            className="rounded-xl font-bold"
          >
            {isRevoking ? "Revoking..." : "Revoke All Sessions"}
          </Button>
        </div>
      </div>
    </PageLayout>
  );
};
