import React, { useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { UserCog, ArrowLeft } from "lucide-react";
import { toast } from "sonner";

import { apiClient } from "@/api/client";
import { useCurrentCustomer } from "@/hooks/useCurrentCustomer";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import {
  updateCustomerSchema,
  type UpdateCustomerFormValues,
} from "./customerSchemas";
import { getCustomerApiErrorMessage } from "./customerApiErrors";

export const EditCustomerPage: React.FC = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const {
    data: customer,
    isLoading,
    isError,
    error,
  } = useCurrentCustomer();

  const form = useForm<UpdateCustomerFormValues>({
    resolver: zodResolver(updateCustomerSchema),
    defaultValues: {
      userName: "",
      firstName: "",
      lastName: "",
      bio: "",
    },
    mode: "onTouched",
  });

  useEffect(() => {
    if (!customer) return;
    form.reset({
      userName: customer.userName,
      firstName: customer.firstName,
      lastName: customer.lastName,
      bio: customer.bio ?? "",
    });
  }, [customer, form]);

  useEffect(() => {
    if (!isError) return;
    const status = (error as { response?: { status?: number } })?.response
      ?.status;
    if (status === 403 || status === 404) {
      navigate("/profile/create", { replace: true });
    }
  }, [isError, error, navigate]);

  const updateMutation = useMutation({
    mutationFn: (values: UpdateCustomerFormValues) =>
      apiClient.updateCustomer({
        userName: values.userName,
        firstName: values.firstName,
        lastName: values.lastName,
        bio: values.bio?.trim() || undefined,
      }),
    onSuccess: (updated) => {
      queryClient.invalidateQueries({ queryKey: ["currentCustomer"] });
      queryClient.invalidateQueries({ queryKey: ["user", updated.userName] });
      toast.success("Profile updated");
      navigate("/profile", { replace: true });
    },
    onError: (err) => {
      toast.error(getCustomerApiErrorMessage(err, "Failed to update profile"));
    },
  });

  const onSubmit = form.handleSubmit((values) => {
    updateMutation.mutate(values);
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[50vh] text-muted-foreground">
        Loading profile...
      </div>
    );
  }

  if (!customer) {
    return null;
  }

  return (
    <div className="flex flex-col items-center w-full min-h-screen bg-background pt-8 pb-16 px-4">
      <div className="max-w-lg w-full">
        <Button
          variant="ghost"
          className="mb-6 -ml-2 gap-2 text-muted-foreground"
          asChild
        >
          <Link to="/profile">
            <ArrowLeft size={16} />
            Back to profile
          </Link>
        </Button>

        <header className="mb-8 flex items-center gap-3">
          <div className="p-3 bg-primary/10 text-primary rounded-full">
            <UserCog size={24} />
          </div>
          <div>
            <h1 className="text-3xl font-serif font-bold tracking-tight text-foreground">
              Edit guest profile
            </h1>
            <p className="text-sm text-muted-foreground mt-1">
              Update how you appear to other guests on ShareBite.
            </p>
          </div>
        </header>

        <form
          onSubmit={onSubmit}
          className="bg-card p-8 rounded-3xl border border-border shadow-lg space-y-5"
        >
          <div className="space-y-2">
            <Label htmlFor="userName">Username</Label>
            <Input
              id="userName"
              className="bg-muted/50 border-border h-12 rounded-xl"
              {...form.register("userName")}
            />
            {form.formState.errors.userName && (
              <p className="text-sm text-destructive">
                {form.formState.errors.userName.message}
              </p>
            )}
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="firstName">First name</Label>
              <Input
                id="firstName"
                className="bg-muted/50 border-border h-12 rounded-xl"
                {...form.register("firstName")}
              />
              {form.formState.errors.firstName && (
                <p className="text-sm text-destructive">
                  {form.formState.errors.firstName.message}
                </p>
              )}
            </div>
            <div className="space-y-2">
              <Label htmlFor="lastName">Last name</Label>
              <Input
                id="lastName"
                className="bg-muted/50 border-border h-12 rounded-xl"
                {...form.register("lastName")}
              />
              {form.formState.errors.lastName && (
                <p className="text-sm text-destructive">
                  {form.formState.errors.lastName.message}
                </p>
              )}
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="bio">Bio (optional)</Label>
            <Textarea
              id="bio"
              className="bg-muted/50 border-border min-h-[100px] rounded-xl resize-none"
              {...form.register("bio")}
            />
            {form.formState.errors.bio && (
              <p className="text-sm text-destructive">
                {form.formState.errors.bio.message}
              </p>
            )}
          </div>

          <div className="flex gap-3 pt-2">
            <Button
              type="button"
              variant="outline"
              className="flex-1 h-12 rounded-xl"
              asChild
            >
              <Link to="/profile">Cancel</Link>
            </Button>
            <Button
              type="submit"
              disabled={updateMutation.isPending}
              className="flex-1 h-12 rounded-xl bg-accent text-accent-foreground hover:bg-[#e6c200] font-bold"
            >
              {updateMutation.isPending ? "Saving..." : "Save changes"}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
};
