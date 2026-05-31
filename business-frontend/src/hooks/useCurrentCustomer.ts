import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/api/client";

export function useCurrentCustomer() {
  const token = localStorage.getItem("token");
  return useQuery({
    queryKey: ["currentCustomer"],
    queryFn: apiClient.getCurrentCustomer,
    staleTime: 5 * 60 * 1000,
    retry: false,
    enabled: !!token,
  });
}
