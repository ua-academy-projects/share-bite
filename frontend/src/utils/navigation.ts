import { isBusinessRole } from "@/utils/auth";

export function getDefaultHomePath(): string {
  return isBusinessRole() ? "/feed/business" : "/feed/users";
}

export function resolvePostAuthDestination(fromPath?: string | null): string {
  if (fromPath && fromPath !== "/auth" && !fromPath.startsWith("/oauth/")) {
    return fromPath;
  }
  return getDefaultHomePath();
}

export const ONBOARDING_PATHS = [
  "/profile/create",
  "/business/setup",
  "/auth",
] as const;

export function isOnboardingPath(pathname: string): boolean {
  if (ONBOARDING_PATHS.includes(pathname as (typeof ONBOARDING_PATHS)[number])) {
    return true;
  }
  return pathname.startsWith("/oauth/");
}
