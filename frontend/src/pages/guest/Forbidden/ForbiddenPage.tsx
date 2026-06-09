import { Link } from "react-router-dom";
import { PageLayout } from "@/components/layout/PageLayout";
import { pageBtnPrimary } from "@/components/layout/pageStyles";
import { Button } from "@/components/ui/button";
import { getDefaultHomePath } from "@/utils/navigation";
import { cn } from "@/lib/utils";

export function ForbiddenPage() {
  return (
    <PageLayout center>
      <div className="max-w-md text-center">
        <h1 className="mb-3 text-3xl font-bold text-[#1A3C34] dark:text-white">
          Access denied
        </h1>
        <p className="mb-6 text-gray-600 dark:text-gray-400">
          You do not have permission to view this page.
        </p>
        <Button asChild className={cn(pageBtnPrimary)}>
          <Link to={getDefaultHomePath()}>Back to feed</Link>
        </Button>
      </div>
    </PageLayout>
  );
}
