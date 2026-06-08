import { Link } from "react-router-dom";
import { Compass, Flame, MapPinned, Package, Search, Sparkles } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const discoveryCards = [
  {
    title: "Venue search",
    description: "Search by keyword, tags, or intent when you know what kind of place you need.",
    to: "/venues/search",
    icon: Search,
    cta: "Search venues",
  },
  {
    title: "Magic boxes",
    description: "Browse discounted food boxes and reserve available offers.",
    to: "/boxes",
    icon: Package,
    cta: "View boxes",
  },
  {
    title: "Recommended feed",
    description: "Return to the feed for recommended business posts and fresh updates.",
    to: "/",
    icon: Flame,
    cta: "Open feed",
  },
];

const intentLinks = [
  "Morning coffee",
  "Work-friendly spots",
  "Family lunch",
  "Dinner plans",
  "Groceries nearby",
  "Dessert rescue",
];

export function DiscoverPage() {
  return (
    <div className="min-h-screen bg-[#F9F7F2] dark:bg-[#0d241d] p-8 md:p-12 w-full transition-colors duration-300">
      <div className="max-w-7xl mx-auto space-y-8">
        <section className="grid grid-cols-1 xl:grid-cols-[1.1fr_0.9fr] gap-6 items-stretch">
          <div className="rounded-3xl bg-[#163d32] text-white p-8 md:p-10 shadow-xl overflow-hidden relative">
            <div className="relative z-10 max-w-2xl">
              <div className="inline-flex items-center gap-2 rounded-full bg-white/10 border border-white/15 px-4 py-2 text-sm text-[#98FF98]">
                <Compass className="w-4 h-4" />
                Discover workspace
              </div>
              <h1 className="text-4xl md:text-5xl font-bold tracking-tight mt-6">
                Explore venues, boxes, and business updates
              </h1>
              <p className="text-gray-200 text-lg mt-4 leading-relaxed">
                Use Discover as the starting point for browsing. Search remains a focused tool,
                while this page gives the business app a broader entry experience.
              </p>
              <div className="flex flex-col sm:flex-row gap-3 mt-8">
                <Button asChild className="bg-[#FFD700] text-[#1A3C34] hover:bg-[#e6c200] rounded-xl font-bold">
                  <Link to="/venues/search">
                    <Search className="w-4 h-4 mr-2" />
                    Search venues
                  </Link>
                </Button>
                <Button asChild className="bg-white/10 text-white hover:bg-white/15 border border-white/15 rounded-xl">
                  <Link to="/boxes">
                    <Package className="w-4 h-4 mr-2" />
                    Browse boxes
                  </Link>
                </Button>
              </div>
            </div>
          </div>

          <Card className="bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-3xl shadow-sm">
            <CardHeader>
              <div className="w-12 h-12 rounded-2xl bg-[#FFD700] text-[#1A3C34] flex items-center justify-center mb-2">
                <Sparkles className="w-6 h-6" />
              </div>
              <CardTitle className="text-2xl font-bold text-[#1A3C34] dark:text-white">
                Quick intents
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex flex-wrap gap-2">
                {intentLinks.map((intent) => (
                  <Link
                    key={intent}
                    to={`/venues/search?intent=${encodeURIComponent(intent)}`}
                    className="rounded-full border border-gray-200 dark:border-[#2f5e50] bg-gray-50 dark:bg-[#0d241d] px-3 py-2 text-sm text-gray-700 dark:text-gray-200 hover:border-[#FFD700] hover:text-[#1A3C34] dark:hover:text-white transition-colors"
                  >
                    {intent}
                  </Link>
                ))}
              </div>
            </CardContent>
          </Card>
        </section>

        <section className="grid grid-cols-1 md:grid-cols-3 gap-5">
          {discoveryCards.map((card) => {
            const Icon = card.icon;
            return (
              <Card
                key={card.title}
                className="bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-2xl shadow-sm"
              >
                <CardHeader>
                  <div className="w-11 h-11 rounded-xl bg-gray-100 dark:bg-[#0d241d] text-emerald-700 dark:text-[#98FF98] flex items-center justify-center mb-2">
                    <Icon className="w-5 h-5" />
                  </div>
                  <CardTitle className="text-xl font-bold text-[#1A3C34] dark:text-white">
                    {card.title}
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-5">
                  <p className="text-gray-600 dark:text-gray-300 leading-relaxed">
                    {card.description}
                  </p>
                  <Button asChild className="w-full bg-[#163d32] text-white hover:bg-[#1A3C34] rounded-xl">
                    <Link to={card.to}>{card.cta}</Link>
                  </Button>
                </CardContent>
              </Card>
            );
          })}
        </section>

        <section className="rounded-2xl border border-emerald-200 dark:border-[#2f5e50] bg-white dark:bg-[#163d32] p-5 flex flex-col md:flex-row gap-4 md:items-center md:justify-between">
          <div className="flex items-start gap-3">
            <MapPinned className="w-5 h-5 text-emerald-700 dark:text-[#98FF98] mt-1" />
            <div>
              <h2 className="font-bold text-[#1A3C34] dark:text-white">Focused search lives separately</h2>
              <p className="text-sm text-gray-600 dark:text-gray-300 mt-1">
                Keep using Venue Search for filters and exact lookups. Discover is now the broader entry point.
              </p>
            </div>
          </div>
          <Button asChild variant="outline" className="rounded-xl">
            <Link to="/venues/search">Open Venue Search</Link>
          </Button>
        </section>
      </div>
    </div>
  );
}
