import { useEffect, useRef, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { Loader2 } from "lucide-react";
import { apiClient } from "@/api/client";
import {
  pageBtnPrimary,
  pageLoader,
  pagePanelLg,
  pageShell,
} from "@/components/layout/pageStyles";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

function oauthExchangeStorageKey(code: string) {
  return `oauth_exchange_${code}`;
}

export function OAuthCallback() {
  const [searchParams] = useSearchParams();
  const [error, setError] = useState("");
  const exchangeStarted = useRef(false);

  useEffect(() => {
    if (exchangeStarted.current) {
      return;
    }

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

    const exchangeKey = oauthExchangeStorageKey(code);
    if (sessionStorage.getItem(exchangeKey)) {
      if (localStorage.getItem("token")) {
        window.location.href = "/";
      }
      return;
    }

    exchangeStarted.current = true;
    sessionStorage.setItem(exchangeKey, "pending");

    const slug =
      sessionStorage.getItem("oauth_role_slug") === "business"
        ? "business"
        : "user";

    const exchange = async () => {
      try {
        await apiClient.oauthCallback("google", code, slug);
        sessionStorage.setItem(exchangeKey, "done");
        sessionStorage.removeItem("oauth_role_slug");
        window.location.href = "/";
      } catch (err: unknown) {
        sessionStorage.removeItem(exchangeKey);
        exchangeStarted.current = false;
        const e = err as { response?: { data?: { error?: string } }; message?: string };
        setError(
          e?.response?.data?.error ||
            e?.message ||
            "OAuth login failed. Please try again."
        );
      }
    };

    void exchange();
  }, [searchParams]);

  const content = error ? (
    <>
      <p className="text-destructive">{error}</p>
      <Button asChild className={cn(pageBtnPrimary, "mt-4")}>
        <Link to="/auth">Back to login</Link>
      </Button>
    </>
  ) : (
    <>
      <Loader2 className={cn(pageLoader, "h-12 w-12")} />
      <p className="text-gray-500 dark:text-gray-400">Completing sign-in…</p>
    </>
  );

  return (
    <div
      className={cn(
        pageShell,
        "flex min-h-[60vh] items-center justify-center px-4"
      )}
    >
      <div
        className={cn(
          pagePanelLg,
          "flex w-full max-w-md flex-col items-center gap-4 p-8 text-center"
        )}
      >
        {content}
      </div>
    </div>
  );
}
