import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { Navbar } from './components/Navbar/Navbar';
import { HomeFeed } from './pages/HomeFeed/HomeFeed';
import { Explore } from './pages/Explore/Explore';
import { RestaurantProfile } from './pages/RestaurantProfile/RestaurantProfile';
import { UserProfile } from './pages/UserProfile/UserProfile';
import { Auth } from './pages/Auth/Auth';
import { CreatePost } from './pages/CreatePost/CreatePost';
import { RequireAuth } from './components/RequireAuth/RequireAuth';
import { ThemeProvider } from './context/ThemeContext';
import './styles/variables.css';

function App() {
  return (
    <ThemeProvider>
      <div className="app-container">
        <Navbar />
        <main className="main-content">
          <Routes>
            <Route path="/" element={<HomeFeed />} />
            <Route path="/explore" element={<Navigate to="/" replace />} />
            <Route path="/restaurant/:id" element={<RestaurantProfile />} />
            <Route path="/profile" element={<RequireAuth><UserProfile /></RequireAuth>} />
            <Route path="/user/:id" element={<RequireAuth><UserProfile /></RequireAuth>} />
            <Route path="/auth" element={<Auth />} />
            <Route path="/post/create" element={<RequireAuth><CreatePost /></RequireAuth>} />
          </Routes>
        </main>
      </div>
    </ThemeProvider>
  );
}

export default App;
