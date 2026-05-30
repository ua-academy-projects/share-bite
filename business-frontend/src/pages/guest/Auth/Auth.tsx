import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import { AlertCircle, ArrowRight } from "lucide-react";
import { apiClient } from "@/api/client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";
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

  const handleGoogleLogin = () => {
    const clientId = import.meta.env.VITE_GOOGLE_CLIENT_ID as string | undefined;
    if (!clientId) {
      setOauthError("Google Sign-In is not configured.");
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

  return (
    <div className="relative flex min-h-[calc(100vh-73px)] items-center justify-center px-4 py-12">
      <div className="pointer-events-none absolute inset-0 overflow-hidden">
        <div className="absolute left-1/4 top-1/4 h-96 w-96 rounded-full bg-primary/10 blur-[100px]" />
        <div className="absolute bottom-1/4 right-1/4 h-96 w-96 rounded-full bg-secondary/10 blur-[100px]" />
      </div>

      <Card className="relative z-10 w-full max-w-md rounded-3xl border-border bg-card-solid shadow-2xl">
        <CardContent className="p-8">
          <div className="mb-8 text-center">
            <h1 className="mb-2 text-3xl font-bold text-foreground">
              {isLogin ? "Welcome back" : "Join Share Bite"}
            </h1>
            <p className="text-muted-foreground">
              {isLogin
                ? "Sign in to see what your friends are eating."
                : "Create an account to start sharing your food journey."}
            </p>
          </div>

          {!isLogin && (
            <div className="mb-6 flex rounded-xl border border-border p-1">
              {(["user", "business"] as const).map((slug) => (
                <button
                  key={slug}
                  type="button"
                  className={cn(
                    "flex-1 rounded-lg py-2 text-sm font-semibold transition-colors",
                    roleSlug === slug
                      ? "bg-[#2f5e50] text-white"
                      : "text-muted-foreground hover:text-foreground"
                  )}
                  onClick={() => setRoleSlug(slug)}
                >
                  {slug === "user" ? "Guest" : "Business"}
                </button>
              ))}
            </div>
          )}

          <form
            onSubmit={(e) => {
              e.preventDefault();
              authMutation.mutate();
            }}
            className="space-y-4"
          >
            {!isLogin && (
              <div className="space-y-2">
                <Label htmlFor="name">Full name</Label>
                <Input
                  id="name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="h-11 rounded-xl"
                  required={!isLogin}
                />
              </div>
            )}
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="h-11 rounded-xl"
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="h-11 rounded-xl"
                required
              />
            </div>

            {authMutation.isError && (
              <div className="flex items-center gap-2 rounded-xl bg-destructive/10 p-3 text-sm text-destructive">
                <AlertCircle className="h-4 w-4 shrink-0" />
                Authentication failed. Check your credentials.
              </div>
            )}

            <Button
              type="submit"
              disabled={authMutation.isPending}
              className="h-11 w-full rounded-xl bg-accent font-bold text-accent-foreground hover:bg-accent/90"
            >
              {authMutation.isPending
                ? "Please wait…"
                : isLogin
                  ? "Sign in"
                  : "Create account"}
              <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
          </form>

          <div className="my-6 flex items-center gap-3">
            <div className="h-px flex-1 bg-border" />
            <span className="text-xs text-muted-foreground">or continue with</span>
            <div className="h-px flex-1 bg-border" />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <Button
              type="button"
              variant="outline"
              className="rounded-xl"
              onClick={handleGoogleLogin}
            >
              Google
            </Button>
            <Button
              type="button"
              variant="outline"
              className="rounded-xl"
              onClick={handleGitHubLogin}
            >
              GitHub
            </Button>
          </div>

          {oauthError && (
            <p className="mt-4 text-center text-sm text-destructive">{oauthError}</p>
          )}

          <p className="mt-6 text-center text-sm text-muted-foreground">
            {isLogin ? "Don't have an account?" : "Already have an account?"}{" "}
            <button
              type="button"
              className="font-semibold text-accent hover:underline"
              onClick={() => setIsLogin(!isLogin)}
            >
              {isLogin ? "Sign up" : "Log in"}
            </button>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
