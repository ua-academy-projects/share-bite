import React, { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export const GitHubSuccess: React.FC = () => {
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

  return (
    <div className="flex min-h-[60vh] items-center justify-center p-8">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>
            {error ? "Sign-in failed" : "Completing GitHub sign-in"}
          </CardTitle>
          <CardDescription>
            {error
              ? "We could not read a session from the callback."
              : "Redirecting you to the app."}
          </CardDescription>
        </CardHeader>
        <CardContent className="flex flex-col items-center gap-4">
          {error ? (
            <>
              <p className="text-center text-sm text-destructive">{error}</p>
              <Button asChild variant="outline">
                <Link to="/auth">Back to login</Link>
              </Button>
            </>
          ) : (
            <Loader2 className="h-8 w-8 animate-spin text-primary" />
          )}
        </CardContent>
      </Card>
    </div>
  );
};
