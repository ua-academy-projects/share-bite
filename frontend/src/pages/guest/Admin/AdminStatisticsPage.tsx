import { useEffect, useState } from "react";
import {
  BarChart3,
  Box,
  Building2,
  FileText,
  Heart,
  Loader2,
  MessageSquare,
  Shield,
  Users,
} from "lucide-react";
import { apiClient } from "@/api/client";
import type { PlatformStatistics } from "@/types/api";
import { PageLayout } from "@/components/layout/PageLayout";
import { pageLoader } from "@/components/layout/pageStyles";
import { Card, CardContent } from "@/components/ui/card";
import { cn } from "@/lib/utils";

type Metric = {
  label: string;
  value: number;
  isFloat?: boolean;
};

type Section = {
  title: string;
  icon: typeof Users;
  metrics: Metric[];
};

function formatValue(value: number, isFloat?: boolean) {
  if (value == null || Number.isNaN(value)) return "—";
  if (isFloat) return value.toFixed(2);
  return value.toLocaleString();
}

function buildSections(s: PlatformStatistics): Section[] {
  return [
    {
      title: "Users & Accounts",
      icon: Users,
      metrics: [
        { label: "Total users", value: s.total_users },
        { label: "Admins", value: s.total_admin_users },
        { label: "Moderators", value: s.total_moderator_users },
        { label: "Regular users", value: s.total_regular_users },
        { label: "Business accounts", value: s.total_business_role_users },
        { label: "Customers", value: s.total_customers },
      ],
    },
    {
      title: "Account Status",
      icon: Shield,
      metrics: [
        { label: "Active", value: s.total_active_users },
        { label: "Muted", value: s.total_muted_users },
        { label: "Suspended", value: s.total_suspended_users },
      ],
    },
    {
      title: "Guest Activity",
      icon: FileText,
      metrics: [
        { label: "Posts", value: s.total_guest_posts },
        { label: "Comments", value: s.total_guest_comments },
        { label: "Post likes", value: s.total_guest_post_likes },
        { label: "Collections", value: s.total_collections },
        {
          label: "Collections w/ collaborators",
          value: s.collections_with_collaborators,
        },
        { label: "Posts w/ collaborators", value: s.posts_with_collaborators },
        { label: "Avg posts / customer", value: s.avg_posts_per_customer, isFloat: true },
        {
          label: "Avg comments / customer",
          value: s.avg_comments_per_customer,
          isFloat: true,
        },
        { label: "Avg comments / post", value: s.avg_comments_per_post, isFloat: true },
      ],
    },
    {
      title: "Business Activity",
      icon: Building2,
      metrics: [
        { label: "Org units", value: s.total_business_org_units },
        { label: "Posts", value: s.total_business_posts },
        { label: "Comments", value: s.total_business_comments },
        { label: "Likes", value: s.total_business_likes },
        { label: "Rescue boxes", value: s.total_business_boxes },
        { label: "Box items", value: s.total_business_box_items },
        { label: "Avg posts / business", value: s.avg_posts_per_business, isFloat: true },
        {
          label: "Avg comments / business",
          value: s.avg_comments_per_business,
          isFloat: true,
        },
        {
          label: "Avg comments / post",
          value: s.avg_business_comments_per_post,
          isFloat: true,
        },
      ],
    },
  ];
}

const HIGHLIGHT_ICONS = [Users, FileText, Building2, Box, MessageSquare, Heart];

export function AdminStatisticsPage() {
  const [stats, setStats] = useState<PlatformStatistics | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let mounted = true;
    const load = async () => {
      setLoading(true);
      setError("");
      try {
        const data = await apiClient.adminGetStatistics();
        if (mounted) setStats(data);
      } catch (err: unknown) {
        const e = err as { response?: { data?: { error?: string } }; message?: string };
        if (mounted) {
          setError(e?.response?.data?.error || e?.message || "Failed to load statistics.");
        }
      } finally {
        if (mounted) setLoading(false);
      }
    };
    void load();
    return () => {
      mounted = false;
    };
  }, []);

  const highlights = stats
    ? [
        { label: "Total users", value: stats.total_users },
        { label: "Guest posts", value: stats.total_guest_posts },
        { label: "Businesses", value: stats.total_business_org_units },
        { label: "Rescue boxes", value: stats.total_business_boxes },
        { label: "Comments", value: stats.total_guest_comments + stats.total_business_comments },
        { label: "Likes", value: stats.total_guest_post_likes + stats.total_business_likes },
      ]
    : [];

  return (
    <PageLayout className="space-y-8">
      <div>
        <h1 className="mb-3 text-4xl font-bold tracking-tight text-[#1A3C34] dark:text-white md:text-5xl">
          Admin — Statistics{" "}
          <span className="text-emerald-500 dark:text-[#98FF98]">📊</span>
        </h1>
        <p className="text-lg text-gray-600 dark:text-gray-400">
          All-time aggregated platform metrics
        </p>
      </div>

      {error ? (
        <div className="rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm font-medium text-red-400">
          {error}
        </div>
      ) : null}

      {loading ? (
        <div className="flex h-64 items-center justify-center">
          <Loader2 className={cn(pageLoader, "h-12 w-12")} />
        </div>
      ) : !stats ? (
        <div className="rounded-3xl border border-gray-200 bg-white p-16 text-center shadow-sm dark:border-[#2f5e50] dark:bg-[#163d32]">
          <p className="text-xl font-bold text-[#1A3C34] dark:text-gray-200">
            No statistics available
          </p>
        </div>
      ) : (
        <>
          <div className="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-6">
            {highlights.map((h, i) => {
              const Icon = HIGHLIGHT_ICONS[i % HIGHLIGHT_ICONS.length];
              return (
                <Card
                  key={h.label}
                  className="rounded-3xl border border-gray-200 bg-white shadow-sm dark:border-[#2f5e50] dark:bg-[#163d32]"
                >
                  <CardContent className="space-y-2 p-5">
                    <div className="flex h-10 w-10 items-center justify-center rounded-2xl border border-[#2f5e50] bg-[#0d241d] text-[#98FF98]">
                      <Icon className="h-5 w-5" />
                    </div>
                    <p className="text-2xl font-bold text-[#1A3C34] dark:text-white">
                      {h.value.toLocaleString()}
                    </p>
                    <p className="text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">
                      {h.label}
                    </p>
                  </CardContent>
                </Card>
              );
            })}
          </div>

          <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
            {buildSections(stats).map((section) => {
              const Icon = section.icon;
              return (
                <Card
                  key={section.title}
                  className="rounded-3xl border border-gray-200 bg-white shadow-sm dark:border-[#2f5e50] dark:bg-[#163d32]"
                >
                  <CardContent className="space-y-4 p-6">
                    <div className="flex items-center gap-2 font-semibold text-[#1A3C34] dark:text-white">
                      <Icon size={20} className="text-emerald-500 dark:text-[#98FF98]" />
                      <span>{section.title}</span>
                    </div>
                    <dl className="divide-y divide-gray-100 dark:divide-[#2f5e50]">
                      {section.metrics.map((m) => (
                        <div
                          key={m.label}
                          className="flex items-center justify-between py-2.5"
                        >
                          <dt className="text-sm text-gray-600 dark:text-gray-300">
                            {m.label}
                          </dt>
                          <dd className="text-sm font-bold text-[#1A3C34] dark:text-white">
                            {formatValue(m.value, m.isFloat)}
                          </dd>
                        </div>
                      ))}
                    </dl>
                  </CardContent>
                </Card>
              );
            })}
          </div>

          <p className="flex items-center gap-2 text-sm text-gray-400">
            <BarChart3 className="h-4 w-4" />
            Metrics are aggregated across auth, guest, and business domains.
          </p>
        </>
      )}
    </PageLayout>
  );
}
