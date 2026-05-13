import type { ReactNode } from "react";

import { cn } from "@/lib/utils";

type BrandSectionProps = {
  title: string;
  subtitle?: string;
  action?: ReactNode;
  className?: string;
  children: ReactNode;
};

export function BrandSection({ title, subtitle, action, className, children }: BrandSectionProps) {
  return (
    <section className={cn("space-y-4", className)}>
      <div className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h2 className="text-lg md:text-xl font-semibold text-[#F9F7F2] tracking-tight">
            {title}
          </h2>
          {subtitle ? (
            <p className="mt-1 text-sm text-[#c7d2cc]">
              {subtitle}
            </p>
          ) : null}
        </div>
        {action ? <div>{action}</div> : null}
      </div>
      {children}
    </section>
  );
}
