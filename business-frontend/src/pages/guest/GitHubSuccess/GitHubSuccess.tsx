import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Loader2 } from "lucide-react";

export function GitHubSuccess() {
  const navigate = useNavigate();

  useEffect(() => {
    const slug =
      (sessionStorage.getItem("oauth_role_slug") as "user" | "business") ||
      "user";
    sessionStorage.removeItem("oauth_role_slug");
    if (slug === "user" && localStorage.getItem("guest_has_customer") !== "1") {
      navigate("/profile/create", { replace: true });
    } else {
      navigate("/", { replace: true });
    }
  }, [navigate]);

  return (
    <div className="flex min-h-[calc(100vh-73px)] items-center justify-center">
      <div className="flex flex-col items-center gap-3 text-muted-foreground">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
        <p>GitHub sign-in successful. Redirecting…</p>
      </div>
    </div>
  );
}
