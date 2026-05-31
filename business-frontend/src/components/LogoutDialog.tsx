import { useState } from "react";
import { LogOut, Monitor, ShieldAlert } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { apiClient } from "@/api/client";
import { clearSessionStorage } from "@/utils/auth";
import {
  pageBtnPrimary,
  pageBtnSecondary,
  pagePanelLg,
} from "@/components/layout/pageStyles";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { cn } from "@/lib/utils";

type LogoutDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function LogoutDialog({ open, onOpenChange }: LogoutDialogProps) {
  const navigate = useNavigate();
  const [loading, setLoading] = useState<"device" | "all" | null>(null);

  const finishLogout = () => {
    clearSessionStorage();
    onOpenChange(false);
    navigate("/auth", { replace: true });
  };

  const handleCurrentDevice = async () => {
    setLoading("device");
    try {
      await apiClient.logout();
      toast.success("Logged out on this device");
    } catch {
      toast.error("Logout failed, clearing local session");
    } finally {
      setLoading(null);
      finishLogout();
    }
  };

  const handleAllDevices = async () => {
    setLoading("all");
    try {
      await apiClient.revokeAllSessions();
      toast.success("Logged out on all devices");
    } catch {
      toast.error("Failed to revoke all sessions");
    } finally {
      setLoading(null);
      finishLogout();
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent
        className={cn(
          pagePanelLg,
          "gap-0 overflow-hidden border-0 p-0 shadow-2xl sm:max-w-md"
        )}
      >
        <DialogHeader className="space-y-2 border-b border-gray-200 px-6 py-5 dark:border-[#2f5e50]">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-emerald-500/15 text-emerald-600 dark:bg-[#98FF98]/15 dark:text-[#98FF98]">
              <LogOut className="h-5 w-5" />
            </div>
            <DialogTitle className="text-xl font-bold text-[#1A3C34] dark:text-white">
              Log out
            </DialogTitle>
          </div>
          <DialogDescription className="text-sm leading-relaxed text-gray-500 dark:text-gray-400">
            Choose whether to sign out on this device only or on all devices.
          </DialogDescription>
        </DialogHeader>

        <div className="flex flex-col gap-3 p-6">
          <Button
            className={cn(pageBtnPrimary, "h-12 w-full gap-2")}
            onClick={handleCurrentDevice}
            disabled={loading !== null}
          >
            <Monitor className="h-4 w-4" />
            {loading === "device" ? "Signing out…" : "This device"}
          </Button>
          <Button
            className={cn(
              pageBtnSecondary,
              "h-12 w-full gap-2 border border-destructive/30 text-destructive hover:bg-destructive/10 hover:text-destructive dark:border-destructive/40 dark:text-red-300 dark:hover:bg-destructive/15"
            )}
            onClick={handleAllDevices}
            disabled={loading !== null}
          >
            <ShieldAlert className="h-4 w-4" />
            {loading === "all" ? "Signing out…" : "All devices"}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
