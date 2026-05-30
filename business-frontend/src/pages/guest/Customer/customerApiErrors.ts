import axios from "axios";

export function getCustomerApiErrorMessage(error: unknown, fallback: string): string {
  if (!axios.isAxiosError(error)) {
    return fallback;
  }
  const data = error.response?.data as
    | { message?: string; error?: string; details?: { field?: string; message?: string }[] }
    | undefined;
  if (data?.details?.length) {
    return data.details.map((d) => d.message || d.field).filter(Boolean).join(". ");
  }
  return data?.message || data?.error || fallback;
}
