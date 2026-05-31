import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";
import { apiClient } from "@/api/client";
import { useCurrentCustomer } from "@/hooks/useCurrentCustomer";
import { PageHeader } from "@/components/layout/PageHeader";
import { PageLayout } from "@/components/layout/PageLayout";
import {
  pageBtnPrimary,
  pageInput,
  pageLabel,
  pageLoader,
  pagePanel,
} from "@/components/layout/pageStyles";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export function ProfileEditPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { data: customer, isLoading } = useCurrentCustomer();
  const [form, setForm] = useState({
    userName: "",
    firstName: "",
    lastName: "",
    bio: "",
  });

  useEffect(() => {
    if (customer) {
      setForm({
        userName: customer.userName,
        firstName: customer.firstName,
        lastName: customer.lastName,
        bio: customer.bio || "",
      });
    }
  }, [customer]);

  const updateMutation = useMutation({
    mutationFn: () => apiClient.updateCustomer(form),
    onSuccess: (updated) => {
      queryClient.invalidateQueries({ queryKey: ["currentCustomer"] });
      queryClient.invalidateQueries({ queryKey: ["user", updated.userName] });
      toast.success("Profile updated successfully!");
      navigate("/profile", { replace: true });
    },
    onError: (error: unknown) => {
      const e = error as { response?: { data?: { error?: string } } };
      toast.error(e?.response?.data?.error || "Failed to update profile");
    },
  });

  if (isLoading) {
    return (
      <PageLayout maxWidth="md">
        <div className="flex h-64 items-center justify-center">
          <Loader2 className={cn(pageLoader, "h-12 w-12")} />
        </div>
      </PageLayout>
    );
  }

  if (!customer) {
    navigate("/profile/create", { replace: true });
    return null;
  }

  return (
    <PageLayout maxWidth="md">
      <PageHeader title="Edit Profile" description="Update your public guest profile" />
      <div className={cn(pagePanel, "p-8")}>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            updateMutation.mutate();
          }}
          className="space-y-5"
        >
          <div className="space-y-2">
            <label htmlFor="userName" className={pageLabel}>
              Username
            </label>
            <input
              id="userName"
              required
              value={form.userName}
              onChange={(e) => setForm((f) => ({ ...f, userName: e.target.value }))}
              className={pageInput}
            />
          </div>
          <div className="space-y-2">
            <label htmlFor="firstName" className={pageLabel}>
              First name
            </label>
            <input
              id="firstName"
              required
              value={form.firstName}
              onChange={(e) => setForm((f) => ({ ...f, firstName: e.target.value }))}
              className={pageInput}
            />
          </div>
          <div className="space-y-2">
            <label htmlFor="lastName" className={pageLabel}>
              Last name
            </label>
            <input
              id="lastName"
              required
              value={form.lastName}
              onChange={(e) => setForm((f) => ({ ...f, lastName: e.target.value }))}
              className={pageInput}
            />
          </div>
          <div className="space-y-2">
            <label htmlFor="bio" className={pageLabel}>
              Bio
            </label>
            <textarea
              id="bio"
              value={form.bio}
              onChange={(e) => setForm((f) => ({ ...f, bio: e.target.value }))}
              className={cn(pageInput, "min-h-[100px] resize-y py-3")}
            />
          </div>
          <Button
            type="submit"
            className={cn(pageBtnPrimary, "w-full")}
            disabled={updateMutation.isPending}
          >
            {updateMutation.isPending ? "Saving…" : "Save Changes"}
          </Button>
        </form>
      </div>
    </PageLayout>
  );
}
