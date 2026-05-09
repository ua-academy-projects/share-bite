import { useState } from "react";
import { useAuth } from "@/hooks/useAuth";
import { authApi } from "@/api/auth";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { AlertTriangle, Globe, Link2, Loader2, LogOut } from "lucide-react";

export function SettingsPage() {
  const { accessToken, refreshToken, clearTokens } = useAuth();
  const [oauthCode, setOauthCode] = useState("");
  const [linkProvider, setLinkProvider] = useState("google");
  const [linkLoading, setLinkLoading] = useState(false);
  const [revokeLoading, setRevokeLoading] = useState(false);
  const [logoutLoading, setLogoutLoading] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  const handleLinkProvider = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!accessToken) return;
    setLinkLoading(true);
    setError("");
    setMessage("");
    try {
      const res = await authApi.linkProvider(accessToken, linkProvider, oauthCode);
      setMessage(res.message);
      setOauthCode("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to link");
    } finally {
      setLinkLoading(false);
    }
  };

  const handleLogout = async () => {
    if (!accessToken || !refreshToken) return;
    setLogoutLoading(true);
    setError("");
    try {
      await authApi.logout(accessToken, refreshToken);
      clearTokens();
      window.location.href = "/login";
    } catch (err) {
      setError(err instanceof Error ? err.message : "Logout failed");
      setLogoutLoading(false);
    }
  };

  const handleRevokeAll = async () => {
    if (!accessToken) return;
    setRevokeLoading(true);
    setError("");
    setMessage("");
    try {
      const res = await authApi.revokeAllSessions(accessToken);
      setMessage(res.message);
      clearTokens();
      window.location.href = "/login";
    } catch (err) {
      setError(err instanceof Error ? err.message : "Revoke failed");
      setRevokeLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-[#F9F7F2] dark:bg-[#0d241d] p-8 md:p-12 w-full transition-colors duration-300">
      <div className="max-w-2xl mx-auto">
        <div className="mb-10">
          <h1 className="text-4xl md:text-5xl font-bold text-[#1A3C34] dark:text-white tracking-tight mb-3">
            Account Settings
          </h1>
          <p className="text-gray-600 dark:text-gray-400 text-lg">Manage your account and security.</p>
        </div>

        <div className="flex flex-col gap-6">
          <Card className="bg-white dark:bg-[#163d32] border-gray-200 dark:border-[#2f5e50]">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-[#1A3C34] dark:text-white">
                <Link2 className="w-5 h-5 text-emerald-500 dark:text-[#98FF98]" />
                Link Social Account
              </CardTitle>
              <CardDescription>Connect a Google or GitHub account</CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleLinkProvider} className="flex flex-col gap-4">
                <div className="flex flex-col gap-2">
                  <Label htmlFor="provider">Provider</Label>
                  <select
                    id="provider"
                    value={linkProvider}
                    onChange={(e) => setLinkProvider(e.target.value)}
                    className="h-10 w-full rounded-lg border border-input bg-transparent px-2.5 py-1 text-sm outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 dark:bg-input/30"
                  >
                    <option value="google">Google</option>
                    <option value="github">GitHub</option>
                  </select>
                </div>
                <div className="flex flex-col gap-2">
                  <Label htmlFor="oauth-code">Authorization Code</Label>
                  <Input
                    id="oauth-code"
                    placeholder="Paste OAuth code"
                    value={oauthCode}
                    onChange={(e) => setOauthCode(e.target.value)}
                    required
                    className="h-10"
                  />
                </div>
                <Button type="submit" disabled={linkLoading} className="h-10 bg-[#1A3C34] text-white hover:bg-[#1A3C34]/90 dark:bg-[#98FF98] dark:text-[#0b0f0e] dark:hover:bg-[#98FF98]/80">
                  {linkLoading ? <Loader2 className="w-4 h-4 animate-spin mr-2" /> : <Globe className="w-4 h-4 mr-2" />}
                  Link Account
                </Button>
              </form>
            </CardContent>
          </Card>

          <Card className="bg-white dark:bg-[#163d32] border-gray-200 dark:border-[#2f5e50]">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-[#1A3C34] dark:text-white">
                <LogOut className="w-5 h-5 text-emerald-500 dark:text-[#98FF98]" />
                Session Management
              </CardTitle>
              <CardDescription>Logout or revoke all sessions</CardDescription>
            </CardHeader>
            <CardContent className="flex flex-col gap-3">
              <Button
                variant="outline"
                className="h-10 w-full"
                onClick={handleLogout}
                disabled={logoutLoading}
              >
                {logoutLoading ? <Loader2 className="w-4 h-4 animate-spin mr-2" /> : <LogOut className="w-4 h-4 mr-2" />}
                Logout (this session)
              </Button>
              <Button
                variant="destructive"
                className="h-10 w-full"
                onClick={handleRevokeAll}
                disabled={revokeLoading}
              >
                {revokeLoading ? <Loader2 className="w-4 h-4 animate-spin mr-2" /> : <AlertTriangle className="w-4 h-4 mr-2" />}
                Revoke All Sessions
              </Button>
            </CardContent>
          </Card>

          {message && (
            <p className="text-emerald-600 dark:text-[#98FF98] text-sm text-center">{message}</p>
          )}
          {error && (
            <p className="text-destructive text-sm text-center">{error}</p>
          )}
        </div>
      </div>
    </div>
  );
}
