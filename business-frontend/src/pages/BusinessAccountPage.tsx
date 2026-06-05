import { Link } from "react-router-dom";
import { Building2, KeyRound, PackagePlus, UserRound } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export function BusinessAccountPage() {
  const hasToken = Boolean(localStorage.getItem("token"));

  return (
    <div className="min-h-screen bg-[#F9F7F2] dark:bg-[#0d241d] p-8 md:p-12 w-full transition-colors duration-300">
      <div className="max-w-6xl mx-auto space-y-8">
        <div className="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
          <div>
            <p className="text-sm font-semibold uppercase tracking-wide text-emerald-700 dark:text-[#98FF98]">
              Account
            </p>
            <h1 className="text-4xl md:text-5xl font-bold text-[#1A3C34] dark:text-white tracking-tight mt-2">
              Business profile
            </h1>
            <p className="text-gray-600 dark:text-gray-400 text-lg mt-3 max-w-2xl">
              A lightweight workspace for the current business session and common owner actions.
            </p>
          </div>
          <div className="inline-flex items-center gap-2 rounded-full bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] px-4 py-2 text-sm text-gray-700 dark:text-gray-200">
            <KeyRound className="w-4 h-4 text-emerald-600 dark:text-[#98FF98]" />
            {hasToken ? "Auth token detected" : "No auth token in this browser"}
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-[1.1fr_0.9fr] gap-5">
          <Card className="bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-2xl shadow-sm">
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 rounded-2xl bg-[#FFD700] text-[#1A3C34] flex items-center justify-center">
                  <UserRound className="w-6 h-6" />
                </div>
                <div>
                  <CardTitle className="text-xl font-bold text-[#1A3C34] dark:text-white">
                    Owner workspace
                  </CardTitle>
                  <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                    Use owned venue profiles as the starting point for posts and magic boxes.
                  </p>
                </div>
              </div>
            </CardHeader>
            <CardContent className="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <Button asChild className="bg-[#163d32] text-white hover:bg-[#1A3C34] rounded-xl">
                <Link to="/venues/mine">
                  <Building2 className="w-4 h-4 mr-2" />
                  My venues
                </Link>
              </Button>
              <Button asChild className="bg-[#FFD700] text-[#1A3C34] hover:bg-[#e6c200] rounded-xl">
                <Link to="/create">
                  <PackagePlus className="w-4 h-4 mr-2" />
                  Create
                </Link>
              </Button>
            </CardContent>
          </Card>

          <Card className="bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-2xl shadow-sm">
            <CardHeader>
              <div className="w-12 h-12 rounded-2xl bg-[#0d241d] text-[#98FF98] flex items-center justify-center mb-2">
                <Building2 className="w-6 h-6" />
              </div>
              <CardTitle className="text-xl font-bold text-[#1A3C34] dark:text-white">
                Business data
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3 text-sm text-gray-600 dark:text-gray-300">
              <p>Organization-level editing can be added here when that workflow is ready.</p>
              <p>For now, My Venues keeps owner actions scoped to locations that belong to the selected brand.</p>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
