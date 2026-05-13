import { useState, useEffect } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Terminal, Save } from "lucide-react";

type DeveloperSettingsModalProps = {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
};

export function DeveloperSettingsModal({
  isOpen,
  onOpenChange,
}: DeveloperSettingsModalProps) {
  const [token, setToken] = useState("");

  useEffect(() => {
    if (isOpen) {
      const savedToken = localStorage.getItem("token") || "";
      setToken(savedToken);
    }
  }, [isOpen]);

  const handleSave = () => {
    localStorage.setItem("token", token);
    onOpenChange(false);
  };

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md border-white/10 bg-[#0d241d] p-8">
        <DialogHeader className="mb-6">
          <DialogTitle className="text-xl font-bold text-[#F9F7F2] flex items-center gap-2">
            <Terminal className="h-5 w-5 text-emerald-400" />
            Developer Mode
          </DialogTitle>
          <DialogDescription className="text-[#cbd5cf]">
            Configure authentication and environment settings for development.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6">
          <div className="space-y-2">
            <Label htmlFor="token" className="text-[#cbd5cf]">
              Auth Header (JWT)
            </Label>
            <Input
              id="token"
              value={token}
              onChange={(e) => setToken(e.target.value)}
              placeholder="Bearer eyJhbGciOiJIUzI1..."
              className="bg-[#163d32] border-white/5 text-white h-11 rounded-xl"
            />
          </div>

          <div className="flex justify-end pt-4">
            <Button
              onClick={handleSave}
              className="bg-emerald-500 text-black hover:bg-emerald-400 font-bold px-6 h-11 rounded-xl flex items-center gap-2"
            >
              <Save className="h-4 w-4" />
              Save Settings
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
