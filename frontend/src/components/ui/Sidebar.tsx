// src/components/ui/sidebar.tsx
import { NavLink } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { useTheme } from "@/components/theme-provider";
import { Moon, Sun } from "lucide-react";

export function Sidebar() {
  const { theme, setTheme } = useTheme();

  const linkClass = ({ isActive }: { isActive: boolean }) =>
    `px-3 py-2 rounded-lg transition-colors duration-200 flex items-center gap-2 ${
      isActive 
        ? "bg-secondary text-secondary-foreground" 
        : "text-muted-foreground hover:bg-secondary/50 hover:text-foreground"
    }`;

  return (
    <aside className="w-64 bg-card border-r border-border p-6 flex flex-col justify-between">
      <div>
        <div className="flex items-center gap-3 mb-6">
          <div className="w-10 h-10 bg-primary rounded-full flex items-center justify-center text-primary-foreground font-bold">
            SB
          </div>
          <div>
            <h2 className="text-foreground font-semibold">Share Bite</h2>
            <p className="text-muted-foreground text-xs">The Art of Dining</p>
          </div>
        </div>

        <Button className="w-full bg-accent text-accent-foreground hover:bg-accent/90 rounded-full mb-8 font-bold shadow-lg shadow-accent/20">
          + Share a Bite
        </Button>

        <nav className="flex flex-col gap-2">
          <NavLink to="/" end className={linkClass}>Home Feed</NavLink>
          <NavLink to="/boxes" className={linkClass}>Magic Boxes</NavLink>
          <NavLink to="/discover" className={linkClass}>Discover</NavLink>
          <NavLink to="/venues/search" className={linkClass}>Venue Search</NavLink>
          
          <div className="mt-6 flex flex-col gap-2">
            <h3 className="text-[10px] font-black uppercase tracking-widest text-muted-foreground/60 px-3 py-1">Social Bites</h3>
            <h3 className="text-[10px] font-black uppercase tracking-widest text-muted-foreground/60 px-3 py-1">Settings</h3>
          </div>
        </nav>
      </div>

      <div className="flex flex-col gap-4">
        {/* Кнопка перемикання теми */}
        <Button 
          variant="ghost" 
          className="justify-start px-3 text-muted-foreground hover:text-foreground hover:bg-secondary/50"
          onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
        >
          {theme === "dark" ? <Sun className="mr-2 h-4 w-4" /> : <Moon className="mr-2 h-4 w-4" />}
          {theme === "dark" ? "Light Mode" : "Dark Mode"}
        </Button>

        <div className="text-muted-foreground/60 text-[10px] flex gap-4 px-3 uppercase font-bold">
          <span className="cursor-pointer hover:text-foreground transition-colors">Support</span>
          <span className="cursor-pointer hover:text-foreground transition-colors">Privacy</span>
        </div>
      </div>
    </aside>
  );
}
