import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { useAuth } from "@/hooks/useAuth";
import { authApi } from "@/api/auth";
import { Globe, Loader2 } from "lucide-react";

export function RegisterPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [slug, setSlug] = useState("user");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const { saveTokens } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      const tokens = await authApi.register(email, password, slug);
      saveTokens(tokens.access_token, tokens.refresh_token);
      navigate("/");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Registration failed");
    } finally {
      setLoading(false);
    }
  };

  const handleGoogleLogin = () => {
    const clientId = import.meta.env.VITE_GOOGLE_CLIENT_ID;
    const redirectUri = import.meta.env.VITE_GOOGLE_REDIRECT_URI || `${window.location.origin}/oauth/google/callback`;
    const scope = "openid email profile";
    const url = `https://accounts.google.com/o/oauth2/v2/auth?client_id=${clientId}&redirect_uri=${encodeURIComponent(redirectUri)}&response_type=code&scope=${encodeURIComponent(scope)}&access_type=offline&prompt=consent`;
    window.location.href = url;
  };

  const handleGitHubLogin = () => {
    window.location.href = authApi.getGitHubLoginUrl();
  };

  return (
    <div className="relative min-h-screen overflow-hidden bg-[#04110d] p-4 sm:p-6">
      <div className="pointer-events-none absolute inset-0">
        <div className="absolute left-1/2 top-[-18rem] h-[32rem] w-[32rem] -translate-x-1/2 rounded-full bg-[#2f8f74]/25 blur-3xl" />
        <div className="absolute bottom-[-12rem] right-[-10rem] h-[24rem] w-[24rem] rounded-full bg-[#98FF98]/10 blur-3xl" />
      </div>

      <div className="relative flex min-h-screen items-center justify-center">
        <Card className="w-full max-w-md border-white/20 bg-[#03120d]/85 text-white shadow-[0_20px_70px_-30px_rgba(126,255,188,0.55)] backdrop-blur-md">
          <CardHeader className="space-y-3 pb-2 text-center">
            <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-2xl border border-[#98FF98]/35 bg-[#0d3128] text-2xl font-bold text-[#98FF98] shadow-[0_0_30px_-12px_rgba(152,255,152,0.7)]">
            SB
            </div>
            <CardTitle className="text-3xl font-semibold tracking-tight text-white">Create account</CardTitle>
            <CardDescription className="text-base text-slate-300">Join Share Bite</CardDescription>
          </CardHeader>
          <CardContent className="pt-2">
            <form onSubmit={handleSubmit} className="flex flex-col gap-5">
              <div className="flex flex-col gap-2">
                <Label htmlFor="email" className="text-sm font-medium text-slate-200">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="you@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className="h-12 border-white/20 bg-[#0a231b]/95 text-white placeholder:text-slate-400 focus-visible:border-[#98FF98]/60 focus-visible:ring-[#98FF98]/40 autofill:[-webkit-text-fill-color:white] autofill:shadow-[inset_0_0_0px_1000px_rgb(10,35,27)]"
              />
              </div>
              <div className="flex flex-col gap-2">
                <Label htmlFor="password" className="text-sm font-medium text-slate-200">Password</Label>
              <Input
                id="password"
                type="password"
                placeholder="Min 8 characters"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                minLength={8}
                className="h-12 border-white/20 bg-[#0a231b]/95 text-white placeholder:text-slate-400 focus-visible:border-[#98FF98]/60 focus-visible:ring-[#98FF98]/40 autofill:[-webkit-text-fill-color:white] autofill:shadow-[inset_0_0_0px_1000px_rgb(10,35,27)]"
              />
              </div>
              <div className="flex flex-col gap-2">
                <Label htmlFor="role" className="text-sm font-medium text-slate-200">Account type</Label>
              <select
                id="role"
                value={slug}
                onChange={(e) => setSlug(e.target.value)}
                className="h-12 w-full rounded-md border border-white/20 bg-[#0a231b]/95 px-3 text-sm text-white outline-none transition focus-visible:border-[#98FF98]/60 focus-visible:ring-2 focus-visible:ring-[#98FF98]/40"
              >
                <option value="user">User</option>
                <option value="business">Business</option>
              </select>
              </div>

              {error && (
                <p className="rounded-md border border-red-400/40 bg-red-500/10 px-3 py-2 text-sm text-red-300">{error}</p>
              )}

              <Button type="submit" disabled={loading} className="h-12 w-full rounded-xl bg-[#87e68c] text-lg font-semibold text-[#052017] shadow-[0_10px_25px_-10px_rgba(152,255,152,0.7)] transition hover:bg-[#98FF98] hover:shadow-[0_16px_35px_-15px_rgba(152,255,152,0.85)]">
                {loading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
                Create Account
              </Button>
            </form>

            <div className="relative my-7">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-white/20" />
              </div>
              <div className="relative flex justify-center text-xs uppercase tracking-[0.18em]">
                <span className="bg-[#03120d] px-3 text-slate-400">or continue with</span>
              </div>
            </div>

            <div className="flex flex-col gap-3">
              <Button variant="outline" className="h-12 w-full rounded-xl border-white/35 bg-white/[0.03] text-base font-medium text-white hover:bg-white/10" onClick={handleGoogleLogin}>
                <svg className="mr-2 h-4 w-4" viewBox="0 0 24 24">
                  <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" fill="#4285F4" />
                  <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" />
                  <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05" />
                  <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" />
                </svg>
                Continue with Google
              </Button>
              <Button variant="outline" className="h-12 w-full rounded-xl border-white/35 bg-white/[0.03] text-base font-medium text-white hover:bg-white/10" onClick={handleGitHubLogin}>
                <Globe className="mr-2 h-4 w-4" />
                Continue with GitHub
              </Button>
            </div>

            <div className="mt-7 text-center text-sm text-slate-300">
              Already have an account?{" "}
              <Link to="/login" className="font-medium text-[#98FF98] transition hover:text-[#b5ffb5] hover:underline">
                Sign in
              </Link>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
