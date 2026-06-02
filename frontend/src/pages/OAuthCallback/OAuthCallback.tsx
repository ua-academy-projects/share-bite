import React, { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { apiClient } from '../../api/client';
import styles from './OAuthCallback.module.css';

export const OAuthCallback: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [error, setError] = useState('');

  useEffect(() => {
    const code = searchParams.get('code');
    const oauthError = searchParams.get('error');

    if (oauthError) {
      setError(`OAuth error: ${oauthError}`);
      return;
    }

    if (!code) {
      setError('Missing authorization code.');
      return;
    }

    const exchange = async () => {
      try {
        const data = await apiClient.oauthCallback('google', code, 'user');
        localStorage.setItem('token', data.access_token);
        window.location.href = '/';
      } catch (err: any) {
        setError(
          err?.response?.data?.error ||
          err?.message ||
          'OAuth login failed. Please try again.'
        );
      }
    };

    exchange();
  }, [searchParams, navigate]);

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
        <p className={styles.message}>Completing sign-in…</p>
      </div>
    </div>
  );
};
