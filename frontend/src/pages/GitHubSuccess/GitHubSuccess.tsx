import React, { useEffect, useState } from 'react';
import styles from './GitHubSuccess.module.css';

export const GitHubSuccess: React.FC = () => {
  const [error, setError] = useState('');

  useEffect(() => {
    const hash = window.location.hash.slice(1);
    if (!hash) {
      setError('No token data received from GitHub.');
      return;
    }

    const params = new URLSearchParams(hash);
    const accessToken = params.get('access_token');
    const refreshToken = params.get('refresh_token');

    if (!accessToken) {
      setError('Missing access token in redirect.');
      return;
    }

    localStorage.setItem('token', accessToken);
    if (refreshToken) {
      localStorage.setItem('refresh_token', refreshToken);
    }

    window.location.href = '/';
  }, []);

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
