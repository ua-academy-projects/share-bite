import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { apiClient } from "@/api/client";
import { clearSessionStorage } from "@/utils/auth";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

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
      <DialogContent className="rounded-3xl border-border bg-card-solid sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Log out</DialogTitle>
          <DialogDescription>
            Choose whether to sign out on this device only or on all devices.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter className="flex-col gap-2 sm:flex-col">
          <Button
            className="w-full rounded-xl"
            onClick={handleCurrentDevice}
            disabled={loading !== null}
          >
            {loading === "device" ? "Signing out…" : "This device"}
          </Button>
          <Button
            variant="destructive"
            className="w-full rounded-xl"
            onClick={handleAllDevices}
            disabled={loading !== null}
          >
            {loading === "all" ? "Signing out…" : "All devices"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
