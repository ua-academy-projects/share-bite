import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/api/client";
import { businessApi } from "@/api/business";
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
    if (context.venueId) {
      setBusinessOrgId(context.venueId);
    }

    const hasVenue = context.venueId != null && context.venueId > 0;
    return {
      needsCustomerSetup: false,
      needsBusinessSetup: !hasVenue,
      isComplete: hasVenue,
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
  return useQuery({
    queryKey: ["onboardingStatus"],
    queryFn: fetchOnboardingStatus,
    enabled: enabled && !!token,
    staleTime: 60_000,
    retry: false,
  });
}
