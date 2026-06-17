import { Navigate } from "react-router-dom";
import { isAdminOrModerator } from "@/utils/auth";

export function RequireAdmin({ children }: { children: React.ReactNode }) {
  const token = localStorage.getItem("token");

  if (!token) {
    return <Navigate to="/auth" replace />;
  }

  if (!isAdminOrModerator()) {
    return <Navigate to="/forbidden" replace />;
  }

  return <>{children}</>;
}
