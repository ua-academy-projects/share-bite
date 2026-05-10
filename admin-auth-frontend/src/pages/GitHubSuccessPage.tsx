import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";
import { Loader2 } from "lucide-react";

export function GitHubSuccessPage() {
  const navigate = useNavigate();
  const { saveTokens } = useAuth();
  const [error, setError] = useState("");

  useEffect(() => {
    const hash = window.location.hash.replace(/^#/, "");
    if (!hash) {
      setError("Missing tokens in callback URL");
      return;
    }

    const params = new URLSearchParams(hash);
    const accessToken = params.get("access_token");
    const refreshToken = params.get("refresh_token");

    if (!accessToken || !refreshToken) {
      setError("Missing access_token or refresh_token");
      return;
    }

    saveTokens(accessToken, refreshToken);
    window.history.replaceState(null, "", window.location.pathname);
    navigate("/", { replace: true });
  }, [navigate, saveTokens]);

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
        <p className="text-muted-foreground">Finishing GitHub sign in...</p>
      </div>
    </div>
  );
}
