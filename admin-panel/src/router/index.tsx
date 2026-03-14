import React, { lazy, Suspense } from 'react';
import { createBrowserRouter, Navigate } from 'react-router-dom';
import { Spin } from 'antd';
import AdminLayout from '../layouts/AdminLayout';

const Login = lazy(() => import('../pages/Login'));
const Dashboard = lazy(() => import('../pages/Dashboard'));
const Products = lazy(() => import('../pages/Products'));
const Categories = lazy(() => import('../pages/Categories'));
const Orders = lazy(() => import('../pages/Orders'));
const CardKeys = lazy(() => import('../pages/CardKeys'));
const PaymentChannels = lazy(() => import('../pages/PaymentChannels'));
const OAuthProviders = lazy(() => import('../pages/OAuthProviders'));

const Loading = () => (
  <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '50vh' }}>
    <Spin size="large" />
  </div>
);

const withSuspense = (Component: React.LazyExoticComponent<React.FC>) => (
  <Suspense fallback={<Loading />}>
    <Component />
  </Suspense>
);

// 路由守卫：检查是否已登录
const RequireAuth: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const token = localStorage.getItem('token');
  if (!token) {
    return <Navigate to="/login" replace />;
  }
  return <>{children}</>;
};

const router = createBrowserRouter([
  {
    path: '/login',
    element: withSuspense(Login),
  },
  {
    path: '/',
    element: (
      <RequireAuth>
        <AdminLayout />
      </RequireAuth>
    ),
    children: [
      { index: true, element: <Navigate to="/dashboard" replace /> },
      { path: 'dashboard', element: withSuspense(Dashboard) },
      { path: 'products', element: withSuspense(Products) },
      { path: 'categories', element: withSuspense(Categories) },
      { path: 'orders', element: withSuspense(Orders) },
      { path: 'card-keys', element: withSuspense(CardKeys) },
      { path: 'payment-channels', element: withSuspense(PaymentChannels) },
      { path: 'oauth-providers', element: withSuspense(OAuthProviders) },
    ],
  },
  {
    path: '*',
    element: <Navigate to="/dashboard" replace />,
  },
]);

export default router;
