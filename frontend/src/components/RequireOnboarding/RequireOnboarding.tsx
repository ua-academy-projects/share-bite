import { Navigate, useLocation } from "react-router-dom";
import { Loader2 } from "lucide-react";
import { useOnboardingStatus } from "@/hooks/useOnboardingStatus";
import { isOnboardingPath } from "@/utils/navigation";

type RequireOnboardingProps = {
  children: React.ReactNode;
};

export function RequireOnboarding({ children }: RequireOnboardingProps) {
  const location = useLocation();
  const token = localStorage.getItem("token");
  const { data, isLoading, isError } = useOnboardingStatus(!!token);

  if (!token || isOnboardingPath(location.pathname)) {
    return <>{children}</>;
  }

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-10 w-10 animate-spin text-emerald-500 dark:text-[#98FF98]" />
      </div>
    );
  }

  if (isError || !data) {
    return <>{children}</>;
  }

  if (data.needsCustomerSetup && location.pathname !== "/profile/create") {
    return <Navigate to="/profile/create" replace state={{ from: location }} />;
  }

  if (data.needsBusinessSetup && location.pathname !== "/business/setup") {
    return <Navigate to="/business/setup" replace state={{ from: location }} />;
  }

  return <>{children}</>;
}
