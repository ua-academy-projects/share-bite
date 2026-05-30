import React, { useEffect } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { UserPlus, ArrowLeft } from "lucide-react";
import { toast } from "sonner";

import { apiClient } from "@/api/client";
import { useCurrentCustomer } from "@/hooks/useCurrentCustomer";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import {
  createCustomerSchema,
  type CreateCustomerFormValues,
} from "./customerSchemas";
import { getCustomerApiErrorMessage } from "./customerApiErrors";

type CreateLocationState = {
  email?: string;
  firstName?: string;
  lastName?: string;
  userName?: string;
};

export const CreateCustomerPage: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const queryClient = useQueryClient();
  const state = (location.state as CreateLocationState | null) ?? {};

  const { data: existingCustomer, isSuccess: hasCustomer } = useCurrentCustomer();

  const form = useForm<CreateCustomerFormValues>({
    resolver: zodResolver(createCustomerSchema),
    defaultValues: {
      email: state.email ?? "",
      userName: state.userName ?? "",
      firstName: state.firstName ?? "",
      lastName: state.lastName ?? "",
      bio: "",
    },
    mode: "onTouched",
  });

  useEffect(() => {
    if (hasCustomer && existingCustomer) {
      navigate("/profile", { replace: true });
    }
  }, [hasCustomer, existingCustomer, navigate]);

  useEffect(() => {
    if (!state.email && !state.userName) return;
    form.reset({
      email: state.email ?? form.getValues("email"),
      userName: state.userName ?? form.getValues("userName"),
      firstName: state.firstName ?? form.getValues("firstName"),
      lastName: state.lastName ?? form.getValues("lastName"),
      bio: form.getValues("bio"),
    });
  }, [state.email, state.firstName, state.lastName, state.userName, form]);

  const createMutation = useMutation({
    mutationFn: (values: CreateCustomerFormValues) =>
      apiClient.createCustomer({
        email: values.email,
        userName: values.userName,
        firstName: values.firstName,
        lastName: values.lastName,
        bio: values.bio?.trim() || undefined,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["currentCustomer"] });
      toast.success("Guest profile created");
      navigate("/profile", { replace: true });
    },
    onError: (error) => {
      toast.error(
        getCustomerApiErrorMessage(error, "Failed to create profile")
      );
    },
  });

  const onSubmit = form.handleSubmit((values) => {
    createMutation.mutate(values);
  });

  return (
    <div className="flex flex-col items-center w-full min-h-screen bg-background pt-8 pb-16 px-4">
      <div className="max-w-lg w-full">
        <Button
          variant="ghost"
          className="mb-6 -ml-2 gap-2 text-muted-foreground"
          asChild
        >
          <Link to="/">
            <ArrowLeft size={16} />
            Back
          </Link>
        </Button>

        <header className="mb-8 flex items-center gap-3">
          <div className="p-3 bg-primary/10 text-primary rounded-full">
            <UserPlus size={24} />
          </div>
          <div>
            <h1 className="text-3xl font-serif font-bold tracking-tight text-foreground">
              Create guest profile
            </h1>
            <p className="text-sm text-muted-foreground mt-1">
              Set up your public username and display name to post, like, and save collections.
            </p>
          </div>
        </header>

        <form
          onSubmit={onSubmit}
          className="bg-card p-8 rounded-3xl border border-border shadow-lg space-y-5"
        >
          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              placeholder="you@example.com"
              className="bg-muted/50 border-border h-12 rounded-xl"
              {...form.register("email")}
            />
            {form.formState.errors.email && (
              <p className="text-sm text-destructive">
                {form.formState.errors.email.message}
              </p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="userName">Username</Label>
            <Input
              id="userName"
              placeholder="foodie123"
              className="bg-muted/50 border-border h-12 rounded-xl"
              {...form.register("userName")}
            />
            <p className="text-xs text-muted-foreground">
              Letters and numbers only, 3–30 characters.
            </p>
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
                placeholder="John"
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
                placeholder="Doe"
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
              placeholder="I love trying new restaurants..."
              className="bg-muted/50 border-border min-h-[100px] rounded-xl resize-none"
              {...form.register("bio")}
            />
            {form.formState.errors.bio && (
              <p className="text-sm text-destructive">
                {form.formState.errors.bio.message}
              </p>
            )}
          </div>

          <Button
            type="submit"
            disabled={createMutation.isPending}
            className="w-full h-12 rounded-xl bg-accent text-accent-foreground hover:bg-[#e6c200] font-bold text-lg shadow-lg mt-2"
          >
            {createMutation.isPending ? "Creating..." : "Create profile"}
          </Button>
        </form>
      </div>
    </div>
  );
};
