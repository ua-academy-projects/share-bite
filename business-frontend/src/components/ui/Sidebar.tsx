// src/components/Sidebar.tsx
import { NavLink } from "react-router-dom";
import { Button } from "@/components/ui/button";

export function Sidebar() {
  const linkClass = ({ isActive }: { isActive: boolean }) =>
    `px-3 py-2 rounded-lg transition ${
      isActive ? "bg-[#2f5e50] text-white" : "text-gray-200 hover:bg-[#2f5e50]"
    }`;

  return (
    <aside className="w-64 bg-[#163d32] p-6 flex flex-col justify-between">
      <div>
        <div className="flex items-center gap-3 mb-6">
          <div className="w-10 h-10 bg-black rounded-full flex items-center justify-center text-white font-bold">
            SB
          </div>
          <div>
            <h2 className="text-white font-semibold">Share Bite</h2>
            <p className="text-gray-300 text-xs">The Art of Dining</p>
          </div>
        </div>

        <Button className="w-full bg-green-500 text-black rounded-full mb-8 hover:bg-green-400">
          + Share a Bite
        </Button>

        <nav className="flex flex-col gap-2">
          <NavLink to="/" end className={linkClass}>Home Feed</NavLink>
          <NavLink to="/boxes" className={linkClass}>Magic Boxes</NavLink>
          <NavLink to="/discover" className={linkClass}>Discover</NavLink>
          <span className="text-gray-400 px-3 py-2 text-sm">Social Bites</span>
          <span className="text-gray-400 px-3 py-2 text-sm">Settings</span>
        </nav>
      </div>

      <div className="text-gray-400 text-xs flex gap-4">
        <span>Support</span>
        <span>Privacy</span>
      </div>
    </aside>
  );
}