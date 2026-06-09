import { cn } from "@/lib/utils";
import { pageShell } from "./pageStyles";

type PageLayoutProps = {
  children: React.ReactNode;
  maxWidth?: "md" | "lg" | "xl" | "2xl" | "3xl" | "5xl" | "6xl" | "7xl";
  className?: string;
  center?: boolean;
};

const maxWidthClass: Record<NonNullable<PageLayoutProps["maxWidth"]>, string> = {
  md: "max-w-md",
  lg: "max-w-lg",
  xl: "max-w-xl",
  "2xl": "max-w-2xl",
  "3xl": "max-w-3xl",
  "5xl": "max-w-5xl",
  "6xl": "max-w-6xl",
  "7xl": "max-w-7xl",
};

export function PageLayout({
  children,
  maxWidth = "7xl",
  className,
  center = false,
}: PageLayoutProps) {
  return (
    <div className={cn(pageShell, center && "flex min-h-screen items-center justify-center")}>
      <div className={cn(maxWidthClass[maxWidth], "mx-auto w-full", className)}>
        {children}
      </div>
    </div>
  );
}
