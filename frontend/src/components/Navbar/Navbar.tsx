import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Moon, Sun, Search, Shield, LogOut, Plus } from 'lucide-react';
import { useTheme } from '@/components/theme-provider';
import { isAdminOrModerator } from '../../utils/auth';
import { apiClient } from '../../api/client';
import { useCurrentCustomer } from '../../hooks/useCurrentCustomer';
import { NotificationBell } from '../Notifications/NotificationBell';

export const Navbar: React.FC = () => {
    const { theme, setTheme } = useTheme();
    const navigate = useNavigate();
    const { data: currentCustomer } = useCurrentCustomer();
    const [showLogoutDialog, setShowLogoutDialog] = useState(false);
    const [logoutLoading, setLogoutLoading] = useState(false);
    const [logoutError, setLogoutError] = useState('');

    const clearLocalSession = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('refresh_token');
    };

    const handleLogoutCurrentDevice = async () => {
        setLogoutLoading(true);
        setLogoutError('');
        try {
            await apiClient.logout();
        } catch (error: any) {
            console.error('Logout failed on server, cleaning local session anyway:', error);
            setLogoutError(error?.response?.data?.error || error?.message || 'Failed to logout.');
        } finally {
            clearLocalSession();
            navigate('/auth', { replace: true, state: { isLogin: true } });
            setShowLogoutDialog(false);
            setLogoutLoading(false);
        }
    };

    const handleLogoutAllDevices = async () => {
        setLogoutLoading(true);
        setLogoutError('');
        try {
            await apiClient.revokeAllSessions();
            clearLocalSession();
            navigate('/auth', { replace: true, state: { isLogin: true } });
            setShowLogoutDialog(false);
        } catch (error: any) {
            setLogoutError(error?.response?.data?.error || error?.message || 'Failed to revoke all sessions.');
        } finally {
            setLogoutLoading(false);
        }
    };

    const isAuthenticated = !!localStorage.getItem('token');
    const currentDate = new Date().toLocaleDateString('en-US', { weekday: 'long', month: 'short', day: 'numeric' }).toUpperCase();

    return (
        <>
            <nav style={{ backgroundColor: 'var(--navbar-bg)', borderColor: 'var(--navbar-border)' }} className="sticky top-0 z-50 w-full border-b backdrop-blur-md py-3 px-6 lg:px-8">
                <div className="flex items-center justify-between w-full max-w-7xl mx-auto">

                    {/* Left: Logo and Date */}
                    <div className="flex items-center gap-6">
                        <Link to="/" className="text-3xl font-serif font-bold tracking-tight" style={{ color: 'var(--navbar-foreground)' }}>
                            ShareBite
                        </Link>
                        <span className="hidden md:inline-flex text-[11px] font-black tracking-[0.2em] px-3 py-1.5 rounded-full border" style={{ color: 'var(--navbar-muted)', backgroundColor: 'rgba(170,206,195,0.1)', borderColor: 'rgba(170,206,195,0.2)' }}>
              {currentDate}
            </span>
                    </div>

                    {/* Center: Search Bar */}
                    <div className="hidden lg:flex flex-1 max-w-lg mx-8">
                        <div className="relative w-full flex items-center">
                            <Search className="absolute left-4" size={18} style={{ color: 'var(--navbar-muted)' }} />
                            <input
                                type="text"
                                placeholder="Search restaurants, users..."
                                className="w-full h-11 rounded-full pl-12 pr-4 focus:outline-none focus:ring-2 transition-all shadow-inner"
                                style={{ backgroundColor: 'rgba(255,255,255,0.07)', color: 'var(--navbar-foreground)', borderColor: 'rgba(170,206,195,0.2)', border: '1px solid rgba(170,206,195,0.2)' }}
                            />
                        </div>
                    </div>

                    {/* Right: Actions and Profile */}
                    <div className="flex items-center gap-4">
                        <button
                            onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
                            className="p-2.5 rounded-full transition-colors"
                            style={{ color: 'var(--navbar-muted)' }}
                            onMouseEnter={e => (e.currentTarget.style.color = 'var(--navbar-foreground)')}
                            onMouseLeave={e => (e.currentTarget.style.color = 'var(--navbar-muted)')}
                        >
                            {theme === 'light' ? <Moon size={20} /> : <Sun size={20} />}
                        </button>

                        {isAuthenticated ? (
                            <>
                                {isAdminOrModerator() && (
                                    <Link to="/admin" className="p-2.5 rounded-full transition-colors" style={{ color: 'var(--navbar-muted)' }}>
                                        <Shield size={20} />
                                    </Link>
                                )}

                                <Link to="/post/create" className="bg-accent text-accent-foreground rounded-full px-5 py-2 text-sm font-bold shadow-lg hover:bg-[#e6c200] transition-all flex items-center gap-1">
                                    <Plus size={18} />
                                    Post
                                </Link>

                                <NotificationBell />

                                <div className="h-8 w-px bg-border/50 mx-2 hidden sm:block"></div>

                                <div className="flex items-center gap-4">
                                    {/* Profile Link - STRICTLY SEPARATE */}
                                    <Link to={currentCustomer?.userName ? `/user/${currentCustomer.userName}` : '/profile'} className="flex items-center gap-3 group hover:opacity-80">
                                        <img
                                            src={currentCustomer?.avatarURL || 'https://via.placeholder.com/40'}
                                            alt="Avatar"
                                            className="w-10 h-10 rounded-full border-2 border-border object-cover group-hover:border-primary transition-colors"
                                        />
                                        <div className="hidden sm:flex flex-col items-start">
                                            <span className="text-sm font-bold" style={{ color: 'var(--navbar-foreground)' }}>@{currentCustomer?.userName || 'user'}</span>
                                            <span className="text-[10px] font-black tracking-wider uppercase" style={{ color: 'var(--navbar-muted)' }}>Guest</span>
                                        </div>
                                    </Link>

                                    {/* Logout Button - STRICTLY SEPARATE */}
                                    <button
                                        onClick={() => setShowLogoutDialog(true)}
                                        className="p-2.5 rounded-full transition-colors"
                                        style={{ color: 'var(--navbar-muted)' }}
                                        aria-label="Logout"
                                    >
                                        <LogOut size={20} />
                                    </button>
                                </div>
                            </>
                        ) : (
                            <div className="flex items-center gap-3">
                                <Link to="/auth" state={{ isLogin: true }} className="text-sm font-bold transition-colors px-4" style={{ color: 'var(--navbar-muted)' }}>
                                    Log in
                                </Link>
                                <Link to="/auth" state={{ isLogin: false }} className="px-6 py-2.5 rounded-full font-bold shadow-md transition-all" style={{ backgroundColor: 'var(--navbar-foreground)', color: 'var(--navbar-bg)' }}>
                                    Sign Up
                                </Link>
                            </div>
                        )}
                    </div>
                </div>
            </nav>

            {/* Logout Dialog */}
            {showLogoutDialog && (
                <div className="fixed inset-0 z-[100] flex items-center justify-center bg-black/60 backdrop-blur-sm" onClick={() => setShowLogoutDialog(false)}>
                    <div className="bg-card text-card-foreground border border-border p-8 rounded-3xl max-w-sm w-full shadow-2xl" onClick={(e) => e.stopPropagation()}>
                        <h3 className="text-2xl font-serif font-bold mb-2">Log out</h3>
                        <p className="text-muted-foreground mb-6 text-sm">
                            Do you want to log out only on this device, or on all devices?
                        </p>
                        {logoutError && <p className="text-sm text-destructive mb-4">{logoutError}</p>}
                        <div className="flex flex-col gap-3">
                            <button
                                className="w-full bg-accent text-accent-foreground font-bold py-3 rounded-xl shadow-md transition-transform hover:scale-[1.02]"
                                onClick={handleLogoutAllDevices}
                                disabled={logoutLoading}
                            >
                                {logoutLoading ? 'Logging out…' : 'All devices'}
                            </button>
                            <button
                                className="w-full bg-muted text-foreground hover:bg-muted/80 font-bold py-3 rounded-xl transition-colors"
                                onClick={handleLogoutCurrentDevice}
                                disabled={logoutLoading}
                            >
                                This device
                            </button>
                            <button
                                className="w-full text-sm text-muted-foreground hover:text-foreground font-bold py-2 mt-2"
                                onClick={() => setShowLogoutDialog(false)}
                                disabled={logoutLoading}
                            >
                                Cancel
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </>
    );
};