import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { ArrowRight, Building2, Store, UserRound, Users } from "lucide-react";
import { Link } from "react-router-dom";

const BUSINESS_FRONTEND_URL = import.meta.env.VITE_BUSINESS_FRONTEND_URL || "http://localhost:5173";
const USERS_FEED_URL = import.meta.env.VITE_USERS_FEED_URL || "http://localhost:5174";

export function LandingPage() {
  return (
    <div className="min-h-screen bg-[#F9F7F2] p-8 md:p-12 transition-colors duration-300 dark:bg-[#0d241d]">
      <div className="mx-auto max-w-5xl">
        <div className="mb-10">
          <h1 className="mb-3 text-4xl font-bold tracking-tight text-[#1A3C34] md:text-5xl dark:text-white">
            Where do you want to go?
          </h1>
          <p className="text-lg text-gray-600 dark:text-gray-400">
            Select your destination after login.
          </p>
        </div>

        <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
          <Card className="group overflow-hidden border-emerald-100 bg-gradient-to-br from-white to-emerald-50/40 transition-all duration-300 hover:-translate-y-1 hover:shadow-xl hover:shadow-emerald-200/50 dark:border-[#2f5e50] dark:bg-gradient-to-br dark:from-[#163d32] dark:to-[#113328] dark:hover:shadow-[#98FF98]/10">
            <CardHeader className="space-y-4">
              <div className="flex items-center justify-between">
                <div className="inline-flex h-12 w-12 items-center justify-center rounded-xl bg-emerald-100 text-emerald-700 dark:bg-[#1d4f42] dark:text-[#98FF98]">
                  <Building2 className="h-6 w-6" />
                </div>
                <Store className="h-5 w-5 text-emerald-400/80 dark:text-[#98FF98]/80" />
              </div>
              <CardTitle className="flex items-center gap-2 text-[#1A3C34] dark:text-white">
                Explore Business
              </CardTitle>
              <CardDescription className="text-gray-600 dark:text-gray-300">
                Open business portal to manage posts and operations
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Button
                className="h-10 bg-[#1A3C34] text-white hover:bg-[#1A3C34]/90 dark:bg-[#98FF98] dark:text-[#0b0f0e] dark:hover:bg-[#98FF98]/80"
                onClick={() => {
                  window.location.href = BUSINESS_FRONTEND_URL;
                }}
              >
                Open Business
                <ArrowRight className="ml-2 h-4 w-4 transition-transform group-hover:translate-x-0.5" />
              </Button>
            </CardContent>
          </Card>

          <Card className="group overflow-hidden border-cyan-100 bg-gradient-to-br from-white to-cyan-50/40 transition-all duration-300 hover:-translate-y-1 hover:shadow-xl hover:shadow-cyan-200/40 dark:border-[#2f5e50] dark:bg-gradient-to-br dark:from-[#163d32] dark:to-[#0f2f3c] dark:hover:shadow-[#98FF98]/10">
            <CardHeader className="space-y-4">
              <div className="flex items-center justify-between">
                <div className="inline-flex h-12 w-12 items-center justify-center rounded-xl bg-cyan-100 text-cyan-700 dark:bg-[#1a4250] dark:text-cyan-200">
                  <Users className="h-6 w-6" />
                </div>
                <UserRound className="h-5 w-5 text-cyan-400/80 dark:text-cyan-200/80" />
              </div>
              <CardTitle className="flex items-center gap-2 text-[#1A3C34] dark:text-white">
                Explore Users Feed
              </CardTitle>
              <CardDescription className="text-gray-600 dark:text-gray-300">
                Open the user-facing feed experience
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Button
                className="h-10 bg-[#1A3C34] text-white hover:bg-[#1A3C34]/90 dark:bg-cyan-200 dark:text-[#07242b] dark:hover:bg-cyan-100"
                onClick={() => {
                  window.location.href = USERS_FEED_URL;
                }}
              >
                Open Feed
                <ArrowRight className="ml-2 h-4 w-4 transition-transform group-hover:translate-x-0.5" />
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
