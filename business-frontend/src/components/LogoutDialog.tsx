import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { apiClient } from "@/api/client";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

type LogoutDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function LogoutDialog({ open, onOpenChange }: LogoutDialogProps) {
  const navigate = useNavigate();
  const [logoutLoading, setLogoutLoading] = useState(false);
  const [logoutError, setLogoutError] = useState("");

  const clearLocalSession = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("refresh_token");
    localStorage.removeItem("guest_has_customer");
    localStorage.removeItem("business_org_id");
  };

  const finishLogout = () => {
    clearLocalSession();
    navigate("/auth", { replace: true, state: { isLogin: true } });
    onOpenChange(false);
    setLogoutLoading(false);
  };

  const handleLogoutCurrentDevice = async () => {
    setLogoutLoading(true);
    setLogoutError("");
    try {
      await apiClient.logout();
    } catch (error: unknown) {
      const err = error as {
        response?: { data?: { error?: string } };
        message?: string;
      };
      setLogoutError(
        err?.response?.data?.error || err?.message || "Failed to logout."
      );
    } finally {
      finishLogout();
    }
  };

  const handleLogoutAllDevices = async () => {
    setLogoutLoading(true);
    setLogoutError("");
    try {
      await apiClient.revokeAllSessions();
      finishLogout();
    } catch (error: unknown) {
      const err = error as {
        response?: { data?: { error?: string } };
        message?: string;
      };
      setLogoutError(
        err?.response?.data?.error ||
          err?.message ||
          "Failed to revoke all sessions."
      );
      setLogoutLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-sm rounded-3xl">
        <DialogHeader>
          <DialogTitle>Log out</DialogTitle>
          <DialogDescription>
            Do you want to log out only on this device, or on all devices?
          </DialogDescription>
        </DialogHeader>
        {logoutError && (
          <p className="text-sm text-destructive">{logoutError}</p>
        )}
        <div className="flex flex-col gap-3">
          <Button
            className="w-full rounded-xl bg-accent font-bold text-accent-foreground"
            onClick={handleLogoutAllDevices}
            disabled={logoutLoading}
          >
            {logoutLoading ? "Logging out…" : "All devices"}
          </Button>
          <Button
            variant="secondary"
            className="w-full rounded-xl font-bold"
            onClick={handleLogoutCurrentDevice}
            disabled={logoutLoading}
          >
            This device
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
