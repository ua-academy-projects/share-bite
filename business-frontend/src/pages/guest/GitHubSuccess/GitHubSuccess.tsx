import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Loader2 } from "lucide-react";
import {
  pageBtnPrimary,
  pageLoader,
  pagePanelLg,
  pageShell,
} from "@/components/layout/pageStyles";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export function GitHubSuccess() {
  const [error, setError] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    const getCookie = (name: string): string | null => {
      const match = document.cookie
        .split("; ")
        .find((row) => row.startsWith(name + "="));
      return match ? decodeURIComponent(match.split("=")[1]) : null;
    };

    const token = getCookie("session");
    if (!token) {
      setError("No session cookie received from GitHub.");
      return;
    }

    localStorage.setItem("token", token);
    document.cookie = "session=; Max-Age=0; path=/";
    navigate("/", { replace: true });
  }, [navigate]);

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
      <p className="text-gray-500 dark:text-gray-400">Completing GitHub sign-in…</p>
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
