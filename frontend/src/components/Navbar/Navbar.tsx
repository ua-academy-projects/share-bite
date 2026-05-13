import React, { useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { Moon, Sun, Search, User, Shield } from 'lucide-react';
import { useTheme } from '../../context/ThemeContext';
import { isAdminOrModerator } from '../../utils/auth';
import { apiClient } from '../../api/client';
import styles from './Navbar.module.css';
import { clsx } from 'clsx';

export const Navbar: React.FC = () => {
  const { theme, toggleTheme } = useTheme();
  const location = useLocation();
  const navigate = useNavigate();
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

  const isActive = (path: string) => location.pathname === path;
  const isAuthenticated = !!localStorage.getItem('token');

  return (
    <>
      <nav className={clsx(styles.navbar, 'glass-panel')}>
        <div className={styles.container}>
        <div className={styles.left}>
          <Link to="/" className={styles.logo}>
            ShareBite
          </Link>
          <div className={styles.links}>
            <Link to="/" className={clsx(styles.link, isActive('/') && styles.active)}>
              Feed
            </Link>
            {/* <Link to="/explore" className={clsx(styles.link, isActive('/explore') && styles.active)}>
              Explore
            </Link> */}
          </div>
        </div>

        <div className={styles.center}>
          <div className={styles.searchBar}>
            <Search size={18} className={styles.searchIcon} />
            <input type="text" placeholder="Search restaurants, users..." className={styles.searchInput} />
          </div>
        </div>

        <div className={styles.right}>
          <button onClick={toggleTheme} className={styles.iconButton} aria-label="Toggle theme">
            {theme === 'light' ? <Moon size={20} /> : <Sun size={20} />}
          </button>
          
          {isAuthenticated ? (
            <>
              {isAdminOrModerator() && (
                <Link
                  to="/admin"
                  className={styles.iconButton}
                  aria-label="Admin panel"
                  title="Admin panel"
                >
                  <Shield size={20} />
                </Link>
              )}
              <Link to="/profile" className={styles.iconButton} aria-label="User profile">
                <User size={20} />
              </Link>
              <Link to="/post/create" className={styles.createButton}>
                + Post
              </Link>
              <button 
                onClick={() => setShowLogoutDialog(true)}
                className={clsx(styles.link, styles.logoutBtn)}
              >
                Logout
              </button>
            </>
          ) : (
            <>
              <Link to="/auth" state={{ isLogin: true }} className={styles.link}>
                Login
              </Link>
              <Link to="/auth" state={{ isLogin: false }} className={styles.createButton}>
                Sign Up
              </Link>
            </>
          )}
        </div>
      </div>
      </nav>

      {showLogoutDialog && (
        <div className={styles.dialogBackdrop} onClick={() => setShowLogoutDialog(false)}>
          <div className={styles.dialog} onClick={(e) => e.stopPropagation()}>
            <h3 className={styles.dialogTitle}>Log out</h3>
            <p className={styles.dialogText}>
              Do you want to log out only on this device, or on all devices?
            </p>
            {logoutError && <p className={styles.dialogError}>{logoutError}</p>}
            <div className={styles.dialogActions}>
              <button
                type="button"
                className={styles.dialogSecondaryBtn}
                onClick={() => setShowLogoutDialog(false)}
                disabled={logoutLoading}
              >
                Cancel
              </button>
              <button
                type="button"
                className={styles.dialogSecondaryBtn}
                onClick={handleLogoutCurrentDevice}
                disabled={logoutLoading}
              >
                This device
              </button>
              <button
                type="button"
                className={styles.dialogPrimaryBtn}
                onClick={handleLogoutAllDevices}
                disabled={logoutLoading}
              >
                {logoutLoading ? 'Logging out…' : 'All devices'}
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
};
