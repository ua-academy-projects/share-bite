// src/components/ui/Sidebar.tsx
import { NavLink } from "react-router-dom";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { useTheme } from "@/components/theme-provider";
import { Moon, Sun, Terminal } from "lucide-react";
import { DeveloperSettingsModal } from "@/components/DeveloperSettingsModal";

export function Sidebar() {
  const { theme, setTheme } = useTheme();
  const [isDevModalOpen, setIsDevModalOpen] = useState(false);

  const linkClass = ({ isActive }: { isActive: boolean }) =>
    `px-3 py-2 rounded-lg transition-colors duration-200 flex items-center gap-2 ${
      isActive 
        ? "bg-[#2f5e50] text-white" 
        : "text-gray-300 hover:bg-[#2f5e50]/50 hover:text-white"
    }`;

  return (
    <aside className="w-64 bg-[#163d32] border-r border-[#2f5e50] p-6 flex flex-col justify-between">
      <div>
        <div className="flex items-center gap-3 mb-6">
          <div className="w-10 h-10 bg-[#0b0f0e] rounded-full flex items-center justify-center text-[#98FF98] font-bold">
            SB
          </div>
          <div>
            <h2 className="text-white font-semibold">Share Bite</h2>
            <p className="text-gray-400 text-xs">The Art of Dining</p>
          </div>
        </div>

        <Button className="w-full bg-[#FFD700] text-[#1A3C34] hover:bg-[#FFD700]/80 rounded-full mb-8 font-bold">
          + Share a Bite
        </Button>

        <nav className="flex flex-col gap-2">
          <NavLink to="/" end className={linkClass}>Home Feed</NavLink>
          <NavLink to="/boxes" className={linkClass}>Magic Boxes</NavLink>
          <NavLink to="/discover" className={linkClass}>Discover</NavLink>
          <NavLink to="/venue/search" className={linkClass}>Venue Search</NavLink>
          
          <div className="mt-4 flex flex-col gap-2">
            <span className="text-gray-400 px-3 py-2 text-sm font-medium">Identity</span>
            {/* 
                Strategy: brandId should be dynamically resolved via AuthContext/Me-endpoint.
                Using '1' as a safe default for development/prototype phases.
            */}
            <NavLink to="/brand/1" className={linkClass}>Brand Profile</NavLink>
            <span className="text-gray-400 px-3 py-2 text-sm font-medium">Settings</span>
          </div>
        </nav>
      </div>

      <div className="flex flex-col gap-4">
        {/* Developer Mode Button */}
        <Button 
          variant="ghost" 
          className="justify-start px-3 text-emerald-400 hover:text-emerald-300 hover:bg-[#2f5e50]/50"
          onClick={() => setIsDevModalOpen(true)}
        >
          <Terminal className="mr-2 h-4 w-4" />
          Developer Mode
        </Button>

        {/* Кнопка перемикання теми */}
        <Button 
          variant="ghost" 
          className="justify-start px-3 text-gray-300 hover:text-white hover:bg-[#2f5e50]/50"
          onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
        >
          {theme === "dark" ? <Sun className="mr-2 h-4 w-4" /> : <Moon className="mr-2 h-4 w-4" />}
          {theme === "dark" ? "Light Mode" : "Dark Mode"}
        </Button>

        <div className="text-gray-400 text-xs flex gap-4 px-3">
          <span className="cursor-pointer hover:text-white transition-colors">Support</span>
          <span className="cursor-pointer hover:text-white transition-colors">Privacy</span>
        </div>
      </div>

      <DeveloperSettingsModal 
        isOpen={isDevModalOpen} 
        onOpenChange={setIsDevModalOpen} 
      />
    </aside>
  );
}
