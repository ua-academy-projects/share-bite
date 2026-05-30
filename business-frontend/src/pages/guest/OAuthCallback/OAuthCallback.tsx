import React, { useEffect, useState } from "react";
import { Link, useNavigate, useSearchParams } from "react-router-dom";
import { Loader2 } from "lucide-react";
import { apiClient } from "@/api/client";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export const OAuthCallback: React.FC = () => {
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

    const exchange = async () => {
      try {
        const data = await apiClient.oauthCallback("google", code, "user");
        localStorage.setItem("token", data.access_token);
        try {
          await apiClient.getCurrentCustomer();
          navigate("/", { replace: true });
        } catch {
          navigate("/profile/create", { replace: true });
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
    };

    void exchange();
  }, [searchParams, navigate]);

  return (
    <div className="flex min-h-[60vh] items-center justify-center p-8">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>
            {error ? "Sign-in failed" : "Completing sign-in"}
          </CardTitle>
          <CardDescription>
            {error
              ? "Something went wrong during Google authentication."
              : "Please wait while we finish connecting your account."}
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
