import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../api/client';

export const useCurrentCustomer = () => {
  const token = localStorage.getItem('token');
  return useQuery({
    queryKey: ['currentCustomer'],
    queryFn: apiClient.getCurrentCustomer,
    staleTime: 5 * 60 * 1000, // 5 minutes
    retry: false, // Don't retry on 401s
    enabled: !!token,
  });
};
