import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import styles from './GitHubSuccess.module.css';

export const GitHubSuccess: React.FC = () => {
  const [error, setError] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    const getCookie = (name: string): string | null => {
      const match = document.cookie
        .split('; ')
        .find(row => row.startsWith(name + '='));
      return match ? decodeURIComponent(match.split('=')[1]) : null;
    };

    const token = getCookie('session');
    if (!token) {
      setError('No session cookie received from GitHub.');
      return;
    }

    localStorage.setItem('token', token);

    // Clear the cookie now that the token is in localStorage.
    document.cookie = 'session=; Max-Age=0; path=/';

    navigate('/', { replace: true });
  }, [navigate]);

  if (error) {
    return (
      <div className={styles.container}>
        <div className={styles.box}>
          <p className={styles.error}>{error}</p>
          <a href="/auth" className={styles.backLink}>Back to login</a>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.box}>
        <div className={styles.spinner} />
        <p className={styles.message}>Completing GitHub sign-in…</p>
      </div>
    </div>
  );
};
