// src/App.tsx
import { Routes, Route } from "react-router-dom";
import { Sidebar } from "@/components/ui/sidebar"; // Виправлено регістр файлу
import { BoxesPage } from "@/pages/BoxesPage";

function Home() {
  return <div className="p-8"><h1 className="text-2xl font-bold">Home Feed 🔥</h1></div>;
}

function Discover() {
  return <div className="p-8"><h1 className="text-2xl font-bold">Discover 🌍</h1></div>;
}

function App() {
  return (
    <div className="flex min-h-screen bg-background text-foreground">
      <Sidebar />
      <main className="flex-1">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/boxes" element={<BoxesPage />} />
          <Route path="/discover" element={<Discover />} />
        </Routes>
      </main>
    </div>
  );
}

export default App;