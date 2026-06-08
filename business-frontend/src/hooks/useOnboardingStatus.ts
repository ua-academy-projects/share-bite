import { useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/api/client";
import { businessApi, type BusinessOnboardingContext } from "@/api/business";
import {
  getTokenRole,
  isAdminOrModerator,
  isBusinessRole,
  isUserRole,
  setBusinessOrgId,
} from "@/utils/auth";

export type OnboardingStatus = {
  needsCustomerSetup: boolean;
  needsBusinessSetup: boolean;
  isComplete: boolean;
  context?: BusinessOnboardingContext | null;
};

async function fetchOnboardingStatus(): Promise<OnboardingStatus> {
  const role = getTokenRole();
  if (!role || isAdminOrModerator()) {
    return {
      needsCustomerSetup: false,
      needsBusinessSetup: false,
      isComplete: true,
    };
  }

  if (isUserRole()) {
    const hasCustomer = await apiClient.hasCurrentCustomer();
    return {
      needsCustomerSetup: !hasCustomer,
      needsBusinessSetup: false,
      isComplete: hasCustomer,
    };
  }

  if (isBusinessRole()) {
    const token = localStorage.getItem("token");
    if (!token) {
      return {
        needsCustomerSetup: false,
        needsBusinessSetup: false,
        isComplete: true,
      };
    }

    const context = await businessApi.getMyOnboardingContext(token);
    const hasBrand = context.brandId != null && context.brandId > 0;
    return {
      needsCustomerSetup: false,
      needsBusinessSetup: !hasBrand,
      isComplete: hasBrand,
      context,
    };
  }

  return {
    needsCustomerSetup: false,
    needsBusinessSetup: false,
    isComplete: true,
  };
}

export function useOnboardingStatus(enabled = true) {
  const token = localStorage.getItem("token");
  const query = useQuery({
    queryKey: ["onboardingStatus"],
    queryFn: fetchOnboardingStatus,
    enabled: enabled && !!token,
    staleTime: 60_000,
    retry: false,
  });

  useEffect(() => {
    const context = query.data?.context;
    if (!context) return;
    if (context.venueId) {
      setBusinessOrgId(context.venueId);
    } else if (context.brandId) {
      setBusinessOrgId(context.brandId);
    }
  }, [query.data]);

  return query;
}
