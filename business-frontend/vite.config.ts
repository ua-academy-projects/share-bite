import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import { fileURLToPath } from "node:url";
import path from "node:path";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      src: path.resolve(__dirname, "./src"),
    },
  },
  server: {
    proxy: {
      "/api/auth": {
        target: "http://localhost:3850",
        changeOrigin: true,
        rewrite: (p) => p.replace(/^\/api/, ""),
      },
      "/api/admin": {
        target: "http://localhost:3850",
        changeOrigin: true,
        rewrite: (p) => p.replace(/^\/api/, ""),
      },
      "/api/users": {
        target: "http://localhost:3850",
        changeOrigin: true,
        rewrite: (p) => p.replace(/^\/api/, ""),
      },
      "/api/user": {
        target: "http://localhost:3850",
        changeOrigin: true,
        rewrite: (p) => p.replace(/^\/api/, ""),
      },
      "/api/guest": {
        target: "http://localhost:3800",
        changeOrigin: true,
        rewrite: (p) => p.replace(/^\/api\/guest/, ""),
      },
      "/api/business": {
        target: "http://localhost:3900",
        changeOrigin: true,
        rewrite: (p) => p.replace(/^\/api/, ""),
      },
      "/api/notifications": {
        target: "http://localhost:4005",
        changeOrigin: true,
        rewrite: (p) => p.replace(/^\/api\/notifications/, ""),
      },
    },
  },
});
