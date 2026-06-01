// src/App.tsx
import { Routes, Route } from "react-router-dom";
import { Sidebar } from "@/components/ui/Sidebar";
import { BoxesPage } from "@/pages/BoxesPage";
import  CreatePostPage from "@/pages/CreatePostPage";
import { CreateBoxPage } from "@/pages/CreateBoxPage";
import { VenueSearchPage } from "@/pages/VenueSearchPage";
import { VenueProfilePage } from "@/pages/VenueProfilePage";
import { HomeFeedPage } from "@/pages/HomeFeedPage";
import { CreateHubPage } from "@/pages/CreateHubPage";
import { BusinessAccountPage } from "@/pages/BusinessAccountPage";
import { DiscoverPage } from "@/pages/DiscoverPage";
import { MyVenuesPage } from "@/pages/MyVenuesPage";
import { QRCodeModalProvider } from "@/contexts/QRCodeModalContext";
import { QRCodeModalContainer } from "@/components/ui/QRCodeModal";

function App() {
  return (
    <QRCodeModalProvider>
      <div className="flex min-h-screen bg-background text-foreground">
        <Sidebar />
        <main className="flex-1">
          <Routes>
            <Route path="/" element={<HomeFeedPage />} />
            <Route path="/boxes" element={<BoxesPage />} />
            <Route path="/create" element={<CreateHubPage />} />
            <Route path="/account" element={<BusinessAccountPage />} />
            <Route path="/discover" element={<DiscoverPage />} />
            <Route path="/venues/mine" element={<MyVenuesPage />} />
            <Route path="/venues/search" element={<VenueSearchPage />} />
            <Route path="/venue/:id/create-post" element={<CreatePostPage />} />
            <Route path="/venue/:id/create-box" element={<CreateBoxPage />} />
            <Route path="/venue/:id" element={<VenueProfilePage />} />
          </Routes>
        </main>
      </div>
      <QRCodeModalContainer />
    </QRCodeModalProvider>
  );
}

export default App;
