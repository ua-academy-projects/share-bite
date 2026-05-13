// src/App.tsx
import { lazy } from "react";
import { Routes, Route, Navigate } from "react-router-dom";
import { DashboardLayout } from "@/layouts/DashboardLayout";

// Consistent lazy loading
const HomePage = lazy(() => import("./pages/HomePage"));
const DiscoverPage = lazy(() => import("./pages/DiscoverPage"));
const BoxesPage = lazy(() => import("@/pages/BoxesPage").then((module) => ({
  default: module.BoxesPage,
})));
const VenueSearchPage = lazy(() => import("@/pages/VenueSearchPage").then((module) => ({
  default: module.VenueSearchPage,
})));
const VenueProfilePage = lazy(() => import("@/pages/VenueProfilePage").then((module) => ({
  default: module.VenueProfilePage,
})));
const CreatePostPage = lazy(() => import("@/pages/CreatePostPage"));
const CreateBoxPage = lazy(() => import("@/pages/CreateBoxPage").then((module) => ({
  default: module.CreateBoxPage,
})));
const BrandProfilePage = lazy(() => import("@/pages/BrandProfilePage").then((module) => ({
  default: module.BrandProfilePage,
})));

function App() {
  return (
    <Routes>
      <Route element={<DashboardLayout />}>
        {/* Dashboard index */}
        <Route index element={<HomePage />} />
        
        {/* Core business routes */}
        <Route path="boxes" element={<BoxesPage />} />
        <Route path="discover" element={<DiscoverPage />} />
        
        {/* Venue management routes - nested for scalability */}
        <Route path="venue">
          <Route path="search" element={<VenueSearchPage />} />
          <Route path=":id" element={<VenueProfilePage />} />
          <Route path=":id/create-post" element={<CreatePostPage />} />
          <Route path=":id/create-box" element={<CreateBoxPage />} />
        </Route>

        {/* Brand identity */}
        <Route path="brand/:id" element={<BrandProfilePage />} />

        {/* Catch-all redirect to home */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Route>
    </Routes>
  );
}

export default App;
