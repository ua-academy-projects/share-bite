import { useEffect, useRef, useState, type ChangeEvent } from "react";
import { Link } from "react-router-dom";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  AlertTriangle,
  Bell,
  Camera,
  Loader2,
  Mail,
  Shield,
  Smartphone,
  Tag,
  User,
  Volume2,
} from "lucide-react";
import { toast } from "sonner";
import { apiClient } from "@/api/client";
import {
  fetchNotificationPreferences,
  updateNotificationPreferences,
} from "@/api/notifications";
import { useCurrentCustomer } from "@/hooks/useCurrentCustomer";
import { PageHeader } from "@/components/layout/PageHeader";
import { PageLayout } from "@/components/layout/PageLayout";
import {
  pageBtnPrimary,
  pageInput,
  pageLabel,
  pageLinkAccent,
  pageLoader,
} from "@/components/layout/pageStyles";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  getBusinessOrgId,
  getTokenPayload,
  isBusinessRole,
} from "@/utils/auth";
import { cn } from "@/lib/utils";
import { playNotificationSound } from "@/utils/audio";

const settingsCardClass =
  "rounded-3xl border border-gray-200 bg-white shadow-sm dark:border-[#2f5e50] dark:bg-[#163d32]";

const pageBtnDangerOutline =
  "rounded-xl border border-red-500/50 bg-transparent px-6 py-2.5 font-bold text-red-400 shadow-sm transition-all hover:border-red-500 hover:bg-red-500/10 disabled:opacity-70";

const AVATAR_MAX_BYTES = 5 * 1024 * 1024;
const AVATAR_ACCEPT = "image/jpeg,image/png";

function avatarInitialFrom(
  userName: string | undefined,
  displayName: string
): string {
  return (
    userName?.charAt(0)?.toUpperCase() || displayName.charAt(0).toUpperCase()
  );
}

function AvatarPreview({
  avatarURL,
  initial,
  size = "md",
}: {
  avatarURL?: string | null;
  initial: string;
  size?: "sm" | "md" | "lg";
}) {
  const sizeClass =
    size === "lg"
      ? "h-24 w-24 text-2xl"
      : size === "sm"
        ? "h-14 w-14 text-xl"
        : "h-20 w-20 text-xl";

  return (
    <div
      className={cn(
        "flex shrink-0 items-center justify-center overflow-hidden rounded-2xl border border-gray-200 bg-[#0d241d] font-bold text-[#98FF98] dark:border-[#2f5e50]",
        sizeClass
      )}
    >
      {avatarURL ? (
        <img src={avatarURL} alt="" className="h-full w-full object-cover" />
      ) : (
        initial
      )}
    </div>
  );
}

function rolePillClass(role: string) {
  if (role === "admin") return "border-red-500/40 bg-red-500/10 text-red-300";
  if (role === "moderator") return "border-[#FFD700]/40 bg-[#FFD700]/10 text-[#FFD700]";
  if (role === "business") return "border-emerald-500/40 bg-emerald-500/10 text-emerald-300";
  return "border-[#2f5e50] bg-[#0d241d] text-gray-200";
}

function statusPillClass(status: string) {
  if (status === "active") return "border-emerald-500/40 bg-emerald-500/10 text-emerald-300";
  if (status === "suspended") return "border-red-500/40 bg-red-500/10 text-red-300";
  return "border-[#2f5e50] bg-[#0d241d] text-gray-300";
}

function Switch({
  checked,
  onChange,
  disabled,
}: {
  checked: boolean;
  onChange: (val: boolean) => void;
  disabled?: boolean;
}) {
  return (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      disabled={disabled}
      onClick={() => onChange(!checked)}
      className={cn(
        "relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:ring-offset-2 disabled:opacity-50",
        checked ? "bg-emerald-500" : "bg-gray-200 dark:bg-[#0d241d]"
      )}
    >
      <span
        className={cn(
          "pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out",
          checked ? "translate-x-5" : "translate-x-0"
        )}
      />
    </button>
  );
}

const PREFERENCE_DETAILS = [
  {
    key: "post_liked",
    label: "Post Likes",
    description: "Get notified when someone likes one of your posts",
  },
  {
    key: "invitation_received",
    label: "Collaboration Invitations",
    description: "Get notified when you are invited to co-author a post or collection",
  },
  {
    key: "post_published",
    label: "Collaborator Posts Published",
    description: "Get notified when a collaborative post you contributed to is published",
  },
  {
    key: "post_invitation_accepted",
    label: "Invitations Accepted",
    description: "Get notified when someone accepts your collaborator invitation",
  },
  {
    key: "business_verified",
    label: "Brand Profile Verified",
    description: "Get notified when your brand page is verified by our admins",
  },
  {
    key: "business_rejected",
    label: "Brand Profile Rejected",
    description: "Get notified if your brand page verification request is rejected",
  },
];

export function AccountSettingsPage() {
  const queryClient = useQueryClient();
  const payload = getTokenPayload();
  const { data: customer, isLoading: customerLoading, isError: noCustomer } =
    useCurrentCustomer();
  const businessOrgId = getBusinessOrgId();

  const [form, setForm] = useState({
    userName: "",
    firstName: "",
    lastName: "",
    bio: "",
  });
  const [isRevoking, setIsRevoking] = useState(false);
  const [avatarPreview, setAvatarPreview] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const [soundEnabled, setSoundEnabled] = useState(() => {
    return localStorage.getItem("notification_sound_enabled") !== "false";
  });
  const [soundVolume, setSoundVolume] = useState(() => {
    const vol = localStorage.getItem("notification_sound_volume");
    return vol ? parseInt(vol, 10) : 50;
  });

  const handleToggleSound = (enabled: boolean) => {
    setSoundEnabled(enabled);
    localStorage.setItem("notification_sound_enabled", String(enabled));
    toast.success(enabled ? "Sound notifications enabled" : "Sound notifications muted");
  };

  const handleVolumeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const vol = parseInt(e.target.value, 10);
    setSoundVolume(vol);
    localStorage.setItem("notification_sound_volume", String(vol));
  };

  const handleTestSound = () => {
    playNotificationSound(soundVolume / 100);
  };

  const token = localStorage.getItem("token") || "";

  const { data: preferences, isLoading: prefsLoading } = useQuery({
    queryKey: ["notificationPreferences"],
    queryFn: () => fetchNotificationPreferences(token),
    enabled: !!token,
  });

  const updatePrefsMutation = useMutation({
    mutationFn: (newPrefs: Record<string, boolean>) =>
      updateNotificationPreferences(token, newPrefs),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["notificationPreferences"] });
      toast.success("Notification preferences updated");
    },
    onError: (error: unknown) => {
      const e = error as { response?: { data?: { error?: string } } };
      toast.error(e?.response?.data?.error || "Failed to update preferences");
    },
  });

  const handleTogglePreference = (key: string, currentValue: boolean) => {
    if (!preferences) return;
    const updated = { ...preferences, [key]: !currentValue };
    updatePrefsMutation.mutate(updated);
  };

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

  useEffect(() => {
    return () => {
      if (avatarPreview) URL.revokeObjectURL(avatarPreview);
    };
  }, [avatarPreview]);

  const invalidateProfileQueries = (userName?: string) => {
    void queryClient.invalidateQueries({ queryKey: ["currentCustomer"] });
    if (userName) {
      void queryClient.invalidateQueries({ queryKey: ["user", userName] });
    }
  };

  const uploadAvatarMutation = useMutation({
    mutationFn: (file: File) => apiClient.uploadCustomerAvatar(file),
    onSuccess: (updated) => {
      if (avatarPreview) {
        URL.revokeObjectURL(avatarPreview);
        setAvatarPreview(null);
      }
      invalidateProfileQueries(updated.userName);
      toast.success("Avatar updated");
    },
    onError: (error: unknown) => {
      const e = error as { response?: { data?: { error?: string } } };
      toast.error(e?.response?.data?.error || "Failed to upload avatar");
    },
  });

  const removeAvatarMutation = useMutation({
    mutationFn: () => apiClient.removeCustomerAvatar(),
    onSuccess: (updated) => {
      if (avatarPreview) {
        URL.revokeObjectURL(avatarPreview);
        setAvatarPreview(null);
      }
      invalidateProfileQueries(updated.userName);
      toast.success("Avatar removed");
    },
    onError: (error: unknown) => {
      const e = error as { response?: { data?: { error?: string } } };
      toast.error(e?.response?.data?.error || "Failed to remove avatar");
    },
  });

  const updateMutation = useMutation({
    mutationFn: () =>
      apiClient.updateCustomer({
        userName: form.userName,
        firstName: form.firstName,
        lastName: form.lastName,
        bio: form.bio,
      }),
    onSuccess: (updated) => {
      void queryClient.invalidateQueries({ queryKey: ["currentCustomer"] });
      void queryClient.invalidateQueries({ queryKey: ["user", updated.userName] });
      toast.success("Profile updated successfully");
    },
    onError: (error: unknown) => {
      const e = error as { response?: { data?: { error?: string } } };
      toast.error(e?.response?.data?.error || "Failed to update profile");
    },
  });

  const email = payload?.email ?? "—";
  const role = payload?.role ?? "—";
  const status = payload?.status ?? "—";
  const userId = payload?.sub ?? "—";

  const displayName =
    customer?.firstName && customer?.lastName
      ? `${customer.firstName} ${customer.lastName}`
      : payload?.name ?? email;

  const avatarInitial = avatarInitialFrom(customer?.userName, displayName);
  const displayedAvatarURL = avatarPreview ?? customer?.avatarURL;
  const avatarBusy = uploadAvatarMutation.isPending || removeAvatarMutation.isPending;

  const handleAvatarFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    e.target.value = "";
    if (!file) return;

    if (!["image/jpeg", "image/png"].includes(file.type)) {
      toast.error("Only JPEG and PNG images are supported");
      return;
    }
    if (file.size > AVATAR_MAX_BYTES) {
      toast.error("Image must be 5 MB or smaller");
      return;
    }

    if (avatarPreview) URL.revokeObjectURL(avatarPreview);
    setAvatarPreview(URL.createObjectURL(file));
    uploadAvatarMutation.mutate(file);
  };

  const handleRemoveAvatar = () => {
    if (avatarPreview) {
      URL.revokeObjectURL(avatarPreview);
      setAvatarPreview(null);
    }
    removeAvatarMutation.mutate();
  };

  const handleRevokeSessions = async () => {
    setIsRevoking(true);
    try {
      await apiClient.revokeAllSessions();
      toast.success(
        "All sessions revoked successfully. You will be logged out of other devices."
      );
    } catch (error: unknown) {
      const e = error as { response?: { data?: { error?: string } } };
      toast.error(e?.response?.data?.error || "Failed to revoke sessions.");
    } finally {
      setIsRevoking(false);
    }
  };

  return (
    <PageLayout className="space-y-8">
      <PageHeader
        title={
          <>
            Account Settings{" "}
            <span className="text-emerald-500 dark:text-[#98FF98]">⚙️</span>
          </>
        }
        description="View your account details and manage profile settings"
      />

      <div className="mx-auto flex w-full max-w-3xl flex-col gap-6">
        <Card className={settingsCardClass}>
          <CardContent className="space-y-6 p-6">
            <div className="flex items-center gap-2 font-semibold text-[#1A3C34] dark:text-white">
              <User size={20} className="text-emerald-500 dark:text-[#98FF98]" />
              <span>Account</span>
            </div>

            <div className="flex items-start gap-4">
              <AvatarPreview
                avatarURL={displayedAvatarURL}
                initial={avatarInitial}
                size="sm"
              />
              <div className="min-w-0 flex-1 space-y-3">
                <p className="text-sm font-medium text-[#1A3C34] dark:text-white">
                  {displayName}
                </p>
                <div className="flex flex-wrap gap-2">
                  <span
                    className={cn(
                      "inline-flex items-center gap-1 rounded-full border px-2.5 py-1 text-xs capitalize",
                      rolePillClass(role)
                    )}
                  >
                    <Shield className="h-3 w-3" />
                    {role}
                  </span>
                  <span
                    className={cn(
                      "inline-flex items-center gap-1 rounded-full border px-2.5 py-1 text-xs capitalize",
                      statusPillClass(status)
                    )}
                  >
                    {status}
                  </span>
                </div>
                <div className="flex flex-wrap gap-4 text-sm text-gray-600 dark:text-gray-300">
                  <span className="inline-flex items-center gap-1">
                    <Mail className="h-4 w-4" />
                    {email}
                  </span>
                  <span className="font-mono text-xs text-gray-500">ID: {userId}</span>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className={settingsCardClass}>
          <CardContent className="space-y-6 p-6">
            <div className="flex items-center gap-2 font-semibold text-[#1A3C34] dark:text-white">
              <Tag size={20} className="text-emerald-500 dark:text-[#98FF98]" />
              <span>Profile information</span>
            </div>

            {customerLoading ? (
              <div className="flex h-32 items-center justify-center">
                <Loader2 className={cn(pageLoader, "h-8 w-8")} />
              </div>
            ) : customer ? (
              <div className="space-y-6">
                <div className="flex flex-col gap-4 border-b border-gray-200 pb-6 sm:flex-row sm:items-center dark:border-[#2f5e50]">
                  <AvatarPreview
                    avatarURL={displayedAvatarURL}
                    initial={avatarInitial}
                    size="lg"
                  />
                  <div className="min-w-0 flex-1 space-y-3">
                    <div>
                      <p className="text-sm font-medium text-[#1A3C34] dark:text-white">
                        Profile photo
                      </p>
                      <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                        JPEG or PNG, up to 5 MB. Shown on your profile and in the
                        feed.
                      </p>
                    </div>
                    <div className="flex flex-wrap gap-2">
                      <input
                        ref={fileInputRef}
                        type="file"
                        accept={AVATAR_ACCEPT}
                        className="hidden"
                        onChange={handleAvatarFileChange}
                      />
                      <Button
                        type="button"
                        variant="outline"
                        className="rounded-xl"
                        disabled={avatarBusy}
                        onClick={() => fileInputRef.current?.click()}
                      >
                        {uploadAvatarMutation.isPending ? (
                          <>
                            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                            Uploading…
                          </>
                        ) : (
                          <>
                            <Camera className="mr-2 h-4 w-4" />
                            {customer.avatarURL || avatarPreview
                              ? "Change photo"
                              : "Upload photo"}
                          </>
                        )}
                      </Button>
                      {customer.avatarURL || avatarPreview ? (
                        <Button
                          type="button"
                          variant="outline"
                          className={cn(pageBtnDangerOutline, "h-auto px-4 py-2")}
                          disabled={avatarBusy}
                          onClick={handleRemoveAvatar}
                        >
                          {removeAvatarMutation.isPending
                            ? "Removing…"
                            : "Remove photo"}
                        </Button>
                      ) : null}
                    </div>
                  </div>
                </div>

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
                    onChange={(e) =>
                      setForm((f) => ({ ...f, userName: e.target.value }))
                    }
                    className={pageInput}
                    placeholder="foodie"
                  />
                </div>
                <div className="grid gap-5 sm:grid-cols-2">
                  <div className="space-y-2">
                    <label htmlFor="firstName" className={pageLabel}>
                      First name
                    </label>
                    <input
                      id="firstName"
                      required
                      value={form.firstName}
                      onChange={(e) =>
                        setForm((f) => ({ ...f, firstName: e.target.value }))
                      }
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
                      onChange={(e) =>
                        setForm((f) => ({ ...f, lastName: e.target.value }))
                      }
                      className={pageInput}
                    />
                  </div>
                </div>
                <div className="space-y-2">
                  <label htmlFor="bio" className={pageLabel}>
                    Bio
                  </label>
                  <textarea
                    id="bio"
                    value={form.bio}
                    onChange={(e) =>
                      setForm((f) => ({ ...f, bio: e.target.value }))
                    }
                    className={cn(pageInput, "min-h-[100px] resize-y py-3")}
                    placeholder="Tell others a little about yourself"
                  />
                </div>
                <div className="flex flex-wrap gap-3 border-t border-gray-200 pt-6 dark:border-[#2f5e50]">
                  <Button
                    type="submit"
                    className={pageBtnPrimary}
                    disabled={updateMutation.isPending}
                  >
                    {updateMutation.isPending ? "Saving…" : "Save profile"}
                  </Button>
                  <Button asChild variant="outline" className="rounded-xl">
                    <Link to="/profile">View public profile</Link>
                  </Button>
                </div>
              </form>
              </div>
            ) : (
              <div className="space-y-4 text-sm text-gray-600 dark:text-gray-300">
                <p>
                  {noCustomer
                    ? "Create a guest profile to set your username and display name."
                    : "Profile information is unavailable right now."}
                </p>
                {noCustomer ? (
                  <Button asChild className={pageBtnPrimary}>
                    <Link to="/profile/create">Create profile</Link>
                  </Button>
                ) : null}
              </div>
            )}
          </CardContent>
        </Card>

        <Card className={settingsCardClass}>
          <CardContent className="space-y-6 p-6">
            <div className="flex items-center gap-2 font-semibold text-[#1A3C34] dark:text-white">
              <Bell size={20} className="text-emerald-500 dark:text-[#98FF98]" />
              <span>Notification Preferences</span>
            </div>

            {prefsLoading ? (
              <div className="flex h-20 items-center justify-center">
                <Loader2 className="h-6 w-6 animate-spin text-emerald-500 dark:text-[#98FF98]" />
              </div>
            ) : preferences ? (
              <div className="space-y-5 divide-y divide-gray-100 dark:divide-gray-800">
                {PREFERENCE_DETAILS.filter((item) => {
                  if (item.key.startsWith("business_")) {
                    return isBusinessRole();
                  }
                  return true;
                }).map((pref) => {
                  const isChecked = !!preferences[pref.key];
                  return (
                    <div
                      key={pref.key}
                      className="flex items-center justify-between pt-4 first:pt-0"
                    >
                      <div className="mr-4 space-y-0.5">
                        <label className="text-sm font-medium text-[#1A3C34] dark:text-white">
                          {pref.label}
                        </label>
                        <p className="text-xs text-gray-500 dark:text-gray-400">
                          {pref.description}
                        </p>
                      </div>
                      <Switch
                        checked={isChecked}
                        disabled={updatePrefsMutation.isPending}
                        onChange={() => handleTogglePreference(pref.key, isChecked)}
                      />
                    </div>
                  );
                })}
              </div>
            ) : (
              <p className="text-sm text-gray-500 dark:text-gray-400">
                Unable to load notification preferences.
              </p>
            )}
          </CardContent>
        </Card>

        <Card className={settingsCardClass}>
          <CardContent className="space-y-6 p-6">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2 font-semibold text-[#1A3C34] dark:text-white">
                <Volume2 size={20} className="text-emerald-500 dark:text-[#98FF98]" />
                <span>Notification Sounds</span>
              </div>
              <Button
                type="button"
                variant="outline"
                size="sm"
                className="rounded-xl text-xs font-semibold h-8"
                disabled={!soundEnabled}
                onClick={handleTestSound}
              >
                Test Sound
              </Button>
            </div>

            <div className="space-y-5 divide-y divide-gray-100 dark:divide-gray-800">
              <div className="flex items-center justify-between">
                <div className="mr-4 space-y-0.5">
                  <label className="text-sm font-medium text-[#1A3C34] dark:text-white">
                    Sound Effects
                  </label>
                  <p className="text-xs text-gray-500 dark:text-gray-400">
                    Play a premium chime sound when you receive a new notification
                  </p>
                </div>
                <Switch
                  checked={soundEnabled}
                  onChange={handleToggleSound}
                />
              </div>

              <div className="flex flex-col gap-2 pt-4">
                <div className="flex items-center justify-between">
                  <label className="text-sm font-medium text-[#1A3C34] dark:text-white">
                    Volume Level
                  </label>
                  <span className="text-xs font-semibold text-gray-500 dark:text-gray-400">
                    {soundVolume}%
                  </span>
                </div>
                <input
                  type="range"
                  min="0"
                  max="100"
                  value={soundVolume}
                  disabled={!soundEnabled}
                  onChange={handleVolumeChange}
                  className="h-2 w-full cursor-pointer appearance-none rounded-lg bg-gray-200 accent-emerald-500 disabled:opacity-50 dark:bg-gray-800"
                />
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className={settingsCardClass}>
          <CardContent className="space-y-6 p-6">
            <div className="flex items-center gap-2 font-semibold text-[#1A3C34] dark:text-white">
              <Smartphone size={20} className="text-emerald-500 dark:text-[#98FF98]" />
              <span>Active Sessions</span>
            </div>

            <div className="flex items-start gap-4">
              <div className="flex h-14 w-14 shrink-0 items-center justify-center rounded-2xl border border-gray-200 bg-[#0d241d] dark:border-[#2f5e50]">
                <Shield className="h-7 w-7 text-[#98FF98]" />
              </div>
              <div className="min-w-0 flex-1">
                <p className="text-sm leading-relaxed text-gray-600 dark:text-gray-300">
                  Manage the devices that are currently logged into your account.
                  Revoking all sessions will log out all other devices.
                </p>
              </div>
            </div>

            <div className="flex flex-col gap-4 border-t border-gray-200 pt-6 sm:flex-row sm:items-center sm:justify-between dark:border-[#2f5e50]">
              <div className="inline-flex items-center gap-2 text-sm font-medium text-red-400">
                <AlertTriangle className="h-4 w-4 shrink-0" />
                This action cannot be undone.
              </div>
              <Button
                variant="outline"
                onClick={() => void handleRevokeSessions()}
                disabled={isRevoking}
                className={cn(pageBtnDangerOutline, "h-auto sm:shrink-0")}
              >
                {isRevoking ? "Revoking…" : "Revoke All Sessions"}
              </Button>
            </div>
          </CardContent>
        </Card>

        {isBusinessRole() ? (
          <div className="px-1">
            <Link
              to={businessOrgId ? `/venue/${businessOrgId}` : "/business/setup"}
              className={pageLinkAccent}
            >
              {businessOrgId ? "Venue profile" : "Set up venue"} →
            </Link>
          </div>
        ) : null}
      </div>
    </PageLayout>
  );
}
