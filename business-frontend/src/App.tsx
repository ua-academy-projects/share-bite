import { Routes, Route } from "react-router-dom";
import { Sidebar } from "@/components/ui/Sidebar";
import { BoxesPage } from "@/pages/BoxesPage";

function Home() {
  return <div className="p-8 text-white"><h1>Home Feed 🔥</h1></div>;
}

function Discover() {
  return <div className="p-8 text-white"><h1>Discover 🌍</h1></div>;
}

function App() {
  return (
    <div className="flex min-h-screen bg-[#0b0f0e]">
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