// src/App.tsx
import { Routes, Route } from "react-router-dom";
import { Sidebar } from "@/components/ui/Sidebar";
import { BoxesPage } from "@/pages/BoxesPage";
import  CreatePostPage from "@/pages/CreatePostPage";
import { CreateBoxPage } from "@/pages/CreateBoxPage";
import { VenueSearchPage } from "@/pages/VenueSearchPage";
import { VenueProfilePage } from "@/pages/VenueProfilePage";
import { QRCodeModalProvider } from "@/contexts/QRCodeModalContext";
import { QRCodeModalContainer } from "@/components/ui/QRCodeModal";

function Home() {
  return <div className="p-8"><h1 className="text-2xl font-bold">Home Feed 🔥</h1></div>;
}

function Discover() {
  return <div className="p-8"><h1 className="text-2xl font-bold">Discover 🌍</h1></div>;
}

function App() {
  return (
    <QRCodeModalProvider>
      <div className="flex min-h-screen bg-background text-foreground">
        <Sidebar />
        <main className="flex-1">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/boxes" element={<BoxesPage />} />
            <Route path="/discover" element={<Discover />} />
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
