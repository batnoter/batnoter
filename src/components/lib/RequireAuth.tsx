import React from 'react';
import { Navigate } from 'react-router-dom';
import { User } from '../../reducer/userSlice';

const RequireAuth = ({ user, children }: { user: User | null, children: JSX.Element }) => {
  if (!user) {
    return <Navigate to="/login" replace />;
  }

  return children;
}

export default RequireAuth;
