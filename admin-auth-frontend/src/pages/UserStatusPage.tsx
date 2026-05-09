import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";
import { authApi } from "@/api/auth";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Loader2, ArrowLeft, CheckCircle } from "lucide-react";

export function UserStatusPage() {
  const { userId } = useParams<{ userId: string }>();
  const { accessToken } = useAuth();
  const [status, setStatus] = useState("");
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  useEffect(() => {
    if (!accessToken || !userId) return;
    const load = async () => {
      try {
        const res = await authApi.getUserStatus(accessToken, userId);
        setStatus(res.status);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load status");
      } finally {
        setLoading(false);
      }
    };
    load();
  }, [accessToken, userId]);

  const handleUpdate = async (newStatus: string) => {
    if (!accessToken || !userId) return;
    setSaving(true);
    setError("");
    setSuccess("");
    try {
      await authApi.updateUserStatus(accessToken, userId, newStatus);
      setStatus(newStatus);
      setSuccess(`Status updated to "${newStatus}"`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update status");
    } finally {
      setSaving(false);
    }
  };

  const statusOptions = ["active", "muted", "suspended"] as const;

  const statusStyle = (s: string) => {
    switch (s) {
      case "active": return "bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-[#98FF98] border-emerald-300 dark:border-emerald-700";
      case "muted": return "bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400 border-yellow-300 dark:border-yellow-700";
      case "suspended": return "bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400 border-red-300 dark:border-red-700";
      default: return "bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-400 border-gray-300 dark:border-gray-600";
    }
  };

  return (
    <div className="min-h-screen bg-[#F9F7F2] dark:bg-[#0d241d] p-8 md:p-12 w-full transition-colors duration-300">
      <div className="max-w-2xl mx-auto">
        <Link
          to="/users"
          className="inline-flex items-center gap-1 text-sm text-primary hover:underline mb-6"
        >
          <ArrowLeft className="w-4 h-4" />
          Back to users
        </Link>

        <div className="mb-8">
          <h1 className="text-3xl font-bold text-[#1A3C34] dark:text-white tracking-tight mb-2">
            User Status
          </h1>
          <p className="text-gray-600 dark:text-gray-400 text-sm font-mono">{userId}</p>
        </div>

        {loading ? (
          <div className="flex justify-center items-center h-40">
            <Loader2 className="w-10 h-10 text-emerald-500 dark:text-[#98FF98] animate-spin" />
          </div>
        ) : (
          <Card className="bg-white dark:bg-[#163d32] border-gray-200 dark:border-[#2f5e50]">
            <CardHeader>
              <CardTitle className="text-[#1A3C34] dark:text-white">Current Status</CardTitle>
              <CardDescription>
                <span className={`inline-block px-3 py-1 rounded-full text-sm font-medium ${statusStyle(status)}`}>
                  {status}
                </span>
              </CardDescription>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground mb-4">Change status:</p>
              <div className="flex flex-wrap gap-3">
                {statusOptions.map((opt) => (
                  <Button
                    key={opt}
                    variant={status === opt ? "default" : "outline"}
                    disabled={saving || status === opt}
                    onClick={() => handleUpdate(opt)}
                    className={`capitalize ${status === opt ? "ring-2 ring-ring" : ""}`}
                  >
                    {saving ? (
                      <Loader2 className="w-4 h-4 animate-spin mr-1" />
                    ) : status === opt ? (
                      <CheckCircle className="w-4 h-4 mr-1" />
                    ) : null}
                    {opt}
                  </Button>
                ))}
              </div>

              {error && <p className="text-destructive text-sm mt-4">{error}</p>}
              {success && <p className="text-emerald-600 dark:text-[#98FF98] text-sm mt-4">{success}</p>}
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
