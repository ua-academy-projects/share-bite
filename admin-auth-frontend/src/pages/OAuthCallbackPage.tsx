import { useEffect, useState } from "react";
import { useNavigate, useParams, useSearchParams } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";
import { authApi } from "@/api/auth";
import { Loader2 } from "lucide-react";

export function OAuthCallbackPage() {
  const { provider } = useParams<{ provider: string }>();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { saveTokens } = useAuth();
  const [error, setError] = useState("");

  useEffect(() => {
    const code = searchParams.get("code");
    if (!code || !provider) {
      setError("Missing authorization code or provider");
      return;
    }

    const exchange = async () => {
      try {
        const tokens = await authApi.oauthCallback(provider, code, "user");
        saveTokens(tokens.access_token, tokens.refresh_token);
        navigate("/", { replace: true });
      } catch (err) {
        setError(err instanceof Error ? err.message : "OAuth login failed");
      }
    };

    exchange();
  }, [navigate, provider, searchParams, saveTokens]);

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="text-center">
          <p className="text-destructive text-lg mb-4">{error}</p>
          <a href="/login" className="text-primary hover:underline">
            Back to login
          </a>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-background">
      <div className="text-center">
        <Loader2 className="w-12 h-12 text-emerald-500 dark:text-[#98FF98] animate-spin mx-auto mb-4" />
        <p className="text-muted-foreground">Completing login...</p>
      </div>
    </div>
  );
}
