import { cn } from "@/lib/utils";

type BrandTab = {
  id: string;
  label: string;
  count?: number;
  loading?: boolean;
};

type BrandTabsProps = {
  tabs: BrandTab[];
  activeTab: string;
  onChange: (tabId: string) => void;
};

export function BrandTabs({ tabs, activeTab, onChange }: BrandTabsProps) {
  return (
    <div className="flex flex-wrap gap-3 border-b border-white/10 pb-4">
      {tabs.map((tab) => {
        const isActive = tab.id === activeTab;
        return (
          <button
            key={tab.id}
            type="button"
            onClick={() => onChange(tab.id)}
            className={cn(
              "px-4 py-2 rounded-full text-sm font-medium transition",
              isActive
                ? "bg-[#98FF98] text-[#0b0f0e]"
                : "bg-[#0f1b17] text-[#cbd5cf] hover:bg-[#1b2e27]"
            )}
          >
            <span className="mr-2">{tab.label}</span>
            {tab.loading ? (
              <span className="inline-block h-4 w-6 rounded-full bg-black/10" />
            ) : typeof tab.count === "number" ? (
              <span className={cn(
                "inline-flex min-w-[28px] justify-center rounded-full px-2 py-0.5 text-xs",
                isActive ? "bg-black/15" : "bg-black/20"
              )}>
                {tab.count}
              </span>
            ) : null}
          </button>
        );
      })}
    </div>
  );
}
