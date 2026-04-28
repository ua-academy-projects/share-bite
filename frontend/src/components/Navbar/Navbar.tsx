import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Moon, Sun, Search, User } from 'lucide-react';
import { useTheme } from '../../context/ThemeContext';
import styles from './Navbar.module.css';
import { clsx } from 'clsx';

export const Navbar: React.FC = () => {
  const { theme, toggleTheme } = useTheme();
  const location = useLocation();

  const isActive = (path: string) => location.pathname === path;

  return (
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
          
          {localStorage.getItem('token') ? (
            <>
              <Link to="/profile" className={styles.iconButton} aria-label="User profile">
                <User size={20} />
              </Link>
              <Link to="/post/create" className={styles.createButton}>
                + Post
              </Link>
              <button 
                onClick={() => { localStorage.removeItem('token'); window.location.reload(); }} 
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
  );
};
