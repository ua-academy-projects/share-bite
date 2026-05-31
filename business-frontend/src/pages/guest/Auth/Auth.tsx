import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import { AlertCircle, ArrowRight } from "lucide-react";
import { apiClient } from "@/api/client";
import {
  pageShell,
  pageBtnPrimary,
  pageBtnSecondary,
  pageInput,
  pageLabel,
  pagePanelLg,
} from "@/components/layout/pageStyles";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

const GOOGLE_REDIRECT_URI =
  (import.meta.env.VITE_GOOGLE_REDIRECT_URI as string | undefined) ||
  `${window.location.origin}/oauth/google/callback`;

function buildGoogleAuthUrl(): string {
  const clientId = import.meta.env.VITE_GOOGLE_CLIENT_ID as string;
  const params = new URLSearchParams({
    client_id: clientId ?? "",
    redirect_uri: GOOGLE_REDIRECT_URI,
    response_type: "code",
    scope: "openid email profile",
    access_type: "offline",
    prompt: "consent",
  });
  return `https://accounts.google.com/o/oauth2/v2/auth?${params.toString()}`;
}

export function Auth() {
  const location = useLocation();
  const navigate = useNavigate();
  const [isLogin, setIsLogin] = useState(location.state?.isLogin ?? true);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const [roleSlug, setRoleSlug] = useState<"user" | "business">("user");
  const [oauthError, setOauthError] = useState("");

  useEffect(() => {
    if (location.state?.isLogin !== undefined) {
      setIsLogin(location.state.isLogin);
    }
  }, [location.state]);

  const authMutation = useMutation({
    mutationFn: async () => {
      if (isLogin) {
        return apiClient.login({ email, password });
      }
      const authData = await apiClient.register({
        email,
        password,
        slug: roleSlug,
      });
      if (roleSlug === "user") {
        const prefix =
          email.split("@")[0].replace(/[^a-zA-Z0-9]/g, "") || "user";
        const userName = `${prefix}${Math.random().toString(36).slice(2, 7)}`;
        const nameParts = name.trim().split(" ");
        await apiClient.createCustomer({
          userName,
          firstName: nameParts[0] || "Unknown",
          lastName: nameParts.slice(1).join(" ") || "User",
          bio: "Food lover!",
        });
      }
      return authData;
    },
    onSuccess: () => {
      if (
        !isLogin &&
        roleSlug === "user" &&
        localStorage.getItem("guest_has_customer") !== "1"
      ) {
        navigate("/profile/create", { replace: true });
        return;
      }
      navigate("/");
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    authMutation.mutate();
  };

  const handleGoogleLogin = () => {
    const clientId = import.meta.env.VITE_GOOGLE_CLIENT_ID as string | undefined;
    if (!clientId) {
      setOauthError("Google Sign-In is not configured. Missing VITE_GOOGLE_CLIENT_ID.");
      return;
    }
    sessionStorage.setItem("oauth_role_slug", roleSlug);
    setOauthError("");
    window.location.href = buildGoogleAuthUrl();
  };

  const handleGitHubLogin = () => {
    sessionStorage.setItem("oauth_role_slug", roleSlug);
    window.location.href = "/api/auth/github";
  };

  const authError = authMutation.error as {
    response?: { data?: { error?: string }; status?: number };
  } | null;

  return (
    <div
      className={cn(
        pageShell,
        "relative flex min-h-screen items-center justify-center overflow-hidden px-4 py-12"
      )}
    >
      <div className="pointer-events-none absolute top-1/4 left-1/4 -z-10 h-96 w-96 rounded-full bg-emerald-500/10 blur-[100px]" />
      <div className="pointer-events-none absolute right-1/4 bottom-1/4 -z-10 h-96 w-96 rounded-full bg-[#98FF98]/10 blur-[100px]" />

      <div className={cn(pagePanelLg, "relative z-10 w-full max-w-md p-10 shadow-2xl")}>
        <div className="mb-8 text-center">
          <h1 className="mb-3 text-4xl font-bold tracking-tight text-[#1A3C34] dark:text-white">
            {isLogin ? "Welcome Back" : "Join Share Bite"}
          </h1>
          <p className="text-lg text-gray-600 dark:text-gray-400">
            {isLogin
              ? "Sign in to see what your friends are eating."
              : "Create an account to start sharing your food journey."}
          </p>
        </div>

        {!isLogin ? (
          <div className="mb-5 grid grid-cols-2 gap-2 rounded-xl bg-gray-100 p-1 dark:bg-[#0d241d]">
            {(["user", "business"] as const).map((role) => (
              <Button
                key={role}
                type="button"
                variant="ghost"
                className={cn(
                  "h-9 rounded-lg px-3 py-2 text-sm font-semibold capitalize transition-colors",
                  roleSlug === role
                    ? "bg-emerald-500 text-black dark:bg-[#98FF98] dark:text-[#1A3C34]"
                    : "text-gray-500 hover:text-[#1A3C34] dark:text-gray-400 dark:hover:text-white"
                )}
                onClick={() => setRoleSlug(role)}
              >
                {role}
              </Button>
            ))}
          </div>
        ) : null}

        <form className="space-y-5" onSubmit={handleSubmit}>
          {!isLogin ? (
            <div className="space-y-2">
              <label htmlFor="name" className={pageLabel}>
                Full Name
              </label>
              <input
                id="name"
                placeholder="John Doe"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className={cn(pageInput, "h-12")}
              />
            </div>
          ) : null}
          <div className="space-y-2">
            <label htmlFor="email" className={pageLabel}>
              Email
            </label>
            <input
              type="email"
              id="email"
              placeholder="name@example.com"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className={cn(pageInput, "h-12")}
            />
          </div>
          <div className="space-y-2">
            <label htmlFor="password" className={pageLabel}>
              Password
            </label>
            <input
              type="password"
              id="password"
              placeholder="••••••••"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className={cn(pageInput, "h-12")}
            />
          </div>

          {authMutation.isError ? (
            <div className="flex items-center gap-2 rounded-xl bg-destructive/10 p-3 text-sm text-destructive">
              <AlertCircle size={16} />
              <span>
                {authError?.response?.data?.error ||
                  (authError?.response?.status === 409
                    ? "User already exists. Please try logging in."
                    : "Authentication failed. Please check your inputs.")}
              </span>
            </div>
          ) : null}

          <Button
            type="submit"
            className={cn(pageBtnPrimary, "group flex h-12 w-full items-center justify-center gap-2 text-lg")}
            disabled={authMutation.isPending}
          >
            {authMutation.isPending
              ? "Loading..."
              : isLogin
                ? "Sign In"
                : "Create Account"}
            {!authMutation.isPending ? (
              <ArrowRight
                size={18}
                className="transition-transform group-hover:translate-x-1"
              />
            ) : null}
          </Button>
        </form>

        <div className="my-8 flex items-center gap-4">
          <div className="h-px flex-1 bg-gray-200 dark:bg-[#2f5e50]" />
          <span className="text-[10px] font-bold tracking-widest whitespace-nowrap text-gray-500 uppercase">
            Or continue with
          </span>
          <div className="h-px flex-1 bg-gray-200 dark:bg-[#2f5e50]" />
        </div>

        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <Button
            type="button"
            className={pageBtnSecondary}
            onClick={handleGoogleLogin}
          >
            Google
          </Button>
          <Button
            type="button"
            className={pageBtnSecondary}
            onClick={handleGitHubLogin}
          >
            GitHub
          </Button>
        </div>

        {oauthError ? (
          <div className="mt-4 rounded-xl bg-destructive/10 p-3 text-center text-sm text-destructive">
            {oauthError}
          </div>
        ) : null}

        <div className="mt-8 text-center text-sm">
          <span className="mr-2 text-gray-500 dark:text-gray-400">
            {isLogin ? "Don't have an account?" : "Already have an account?"}
          </span>
          <button
            type="button"
            className="font-bold text-emerald-600 transition-colors hover:text-emerald-700 dark:text-[#98FF98] dark:hover:text-emerald-300"
            onClick={() => setIsLogin(!isLogin)}
          >
            {isLogin ? "Sign Up" : "Log In"}
          </button>
        </div>
      </div>
    </div>
  );
}
