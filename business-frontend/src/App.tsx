// src/App.tsx
import { Routes, Route } from "react-router-dom";
import { Sidebar } from "@/components/ui/Sidebar";
import { BoxesPage } from "@/pages/BoxesPage";
import CreatePostPage from "@/pages/CreatePostPage";
import { CreateBoxPage } from "@/pages/CreateBoxPage";
import { VenueSearchPage } from "@/pages/VenueSearchPage";
import { VenueProfilePage } from "@/pages/VenueProfilePage";
import DiscoverPage from "@/pages/DiscoverPage";
import HomePage from "@/pages/HomePage";
import { BrandProfilePage } from "@/pages/BrandProfilePage";

function App() {
  return (
    <div className="flex min-h-screen bg-background text-foreground">
      <Sidebar />
      <main className="flex-1">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/boxes" element={<BoxesPage />} />
          <Route path="/discover" element={<DiscoverPage />} />
          <Route path="/venues/search" element={<VenueSearchPage />} />
          <Route path="/venue/:id/create-post" element={<CreatePostPage />} />
          <Route path="/venue/:id/create-box" element={<CreateBoxPage />} />
          <Route path="/venue/:id" element={<VenueProfilePage />} />
          <Route path="/brand/:id" element={<BrandProfilePage />} />
        </Routes>
      </main>
    </div>
  );
}

export default App;
