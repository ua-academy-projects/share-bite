import { useEffect, useState } from "react";
import { Link, useNavigate, useSearchParams } from "react-router-dom";
import { Loader2 } from "lucide-react";
import { apiClient } from "@/api/client";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";

export function OAuthCallback() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [error, setError] = useState("");

  useEffect(() => {
    const code = searchParams.get("code");
    const oauthError = searchParams.get("error");

    if (oauthError) {
      setError(`OAuth error: ${oauthError}`);
      return;
    }
    if (!code) {
      setError("Missing authorization code.");
      return;
    }

    const slug =
      (sessionStorage.getItem("oauth_role_slug") as "user" | "business") ||
      "user";

    void (async () => {
      try {
        await apiClient.oauthCallback("google", code, slug);
        sessionStorage.removeItem("oauth_role_slug");
        if (
          slug === "user" &&
          localStorage.getItem("guest_has_customer") !== "1"
        ) {
          navigate("/profile/create", { replace: true });
        } else {
          navigate("/", { replace: true });
        }
      } catch (err: unknown) {
        const e = err as {
          response?: { data?: { error?: string } };
          message?: string;
        };
        setError(
          e?.response?.data?.error ||
            e?.message ||
            "OAuth login failed. Please try again."
        );
      }
    })();
  }, [searchParams, navigate]);

  if (error) {
    return (
      <div className="flex min-h-[calc(100vh-73px)] items-center justify-center px-4">
        <Card className="max-w-md rounded-3xl bg-card-solid">
          <CardContent className="p-8 text-center">
            <p className="text-destructive">{error}</p>
            <Button asChild className="mt-6 rounded-xl">
              <Link to="/auth">Back to login</Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="flex min-h-[calc(100vh-73px)] items-center justify-center">
      <div className="flex flex-col items-center gap-3 text-muted-foreground">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
        <p>Completing sign-in…</p>
      </div>
    </div>
  );
}
