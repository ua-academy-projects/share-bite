import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { isAdminOrModerator } from '../../utils/auth';

export const RequireAdmin: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const location = useLocation();
  const token = localStorage.getItem('token');

  if (!token) {
    return <Navigate to="/auth" state={{ from: location }} replace />;
  }

  if (!isAdminOrModerator()) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
};
