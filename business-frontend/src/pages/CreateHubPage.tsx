import { Link } from "react-router-dom";
import { Building2, FilePlus2, PackagePlus, Store } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const createActions = [
  {
    title: "Create a venue post",
    description: "Pick one of your business venues first, then publish a post from that location.",
    icon: FilePlus2,
    to: "/venues/mine",
    cta: "Choose my venue",
  },
  {
    title: "Create a magic box",
    description: "Open an owned venue profile and create a discounted box for that location.",
    icon: PackagePlus,
    to: "/venues/mine",
    cta: "Choose my venue",
  },
  {
    title: "Browse existing boxes",
    description: "Review live boxes, availability, and reservation behavior.",
    icon: Store,
    to: "/boxes",
    cta: "View boxes",
  },
];

export function CreateHubPage() {
  return (
    <div className="min-h-screen bg-[#F9F7F2] dark:bg-[#0d241d] p-8 md:p-12 w-full transition-colors duration-300">
      <div className="max-w-6xl mx-auto space-y-8">
        <div className="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
          <div>
            <p className="text-sm font-semibold uppercase tracking-wide text-emerald-700 dark:text-[#98FF98]">
              Business workspace
            </p>
            <h1 className="text-4xl md:text-5xl font-bold text-[#1A3C34] dark:text-white tracking-tight mt-2">
              Create content
            </h1>
            <p className="text-gray-600 dark:text-gray-400 text-lg mt-3 max-w-2xl">
              Start from a venue so posts and boxes stay attached to the right business location.
            </p>
          </div>
          <Button asChild className="bg-[#163d32] text-white hover:bg-[#1A3C34] rounded-xl px-5">
            <Link to="/venues/mine">
              <Building2 className="w-4 h-4 mr-2" />
              My venues
            </Link>
          </Button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-5">
          {createActions.map((action) => {
            const Icon = action.icon;
            return (
              <Card
                key={action.title}
                className="bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-2xl shadow-sm"
              >
                <CardHeader>
                  <div className="w-11 h-11 rounded-xl bg-[#FFD700] text-[#1A3C34] flex items-center justify-center mb-2">
                    <Icon className="w-5 h-5" />
                  </div>
                  <CardTitle className="text-xl font-bold text-[#1A3C34] dark:text-white">
                    {action.title}
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-5">
                  <p className="text-gray-600 dark:text-gray-300 leading-relaxed">
                    {action.description}
                  </p>
                  <Button asChild className="w-full bg-[#163d32] text-white hover:bg-[#1A3C34] rounded-xl">
                    <Link to={action.to}>{action.cta}</Link>
                  </Button>
                </CardContent>
              </Card>
            );
          })}
        </div>

        <div className="rounded-2xl border border-emerald-200 dark:border-[#2f5e50] bg-white dark:bg-[#163d32] p-5 text-sm text-gray-600 dark:text-gray-300">
          Direct create pages are still available after an owned venue is selected. Public venue search stays
          separate from business management actions.
        </div>
      </div>
    </div>
  );
}
