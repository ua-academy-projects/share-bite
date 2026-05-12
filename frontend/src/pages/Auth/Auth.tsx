import React, { useState, useEffect } from 'react';
import { Button } from '../../components/Button/Button';
import styles from './Auth.module.css';
import { clsx } from 'clsx';
import { useLocation } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { apiClient } from '../../api/client';
import axios from 'axios';
import { AlertCircle } from 'lucide-react';

const GOOGLE_REDIRECT_URI = `${window.location.origin}/oauth/google/callback`;

function buildGoogleAuthUrl(): string {
  const clientId = import.meta.env.VITE_GOOGLE_CLIENT_ID as string;
  if (!clientId) {
    console.warn('VITE_GOOGLE_CLIENT_ID is not set');
  }
  const params = new URLSearchParams({
    client_id: clientId ?? '',
    redirect_uri: GOOGLE_REDIRECT_URI,
    response_type: 'code',
    scope: 'openid email profile',
    access_type: 'offline',
    prompt: 'consent',
  });
  return `https://accounts.google.com/o/oauth2/v2/auth?${params.toString()}`;
}

export const Auth: React.FC = () => {
  const location = useLocation();
  const [isLogin, setIsLogin] = useState(location.state?.isLogin ?? true);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [name, setName] = useState('');
  const [oauthError, setOauthError] = useState('');

  useEffect(() => {
    if (location.state?.isLogin !== undefined) {
      setIsLogin(location.state.isLogin);
    }
  }, [location.state]);

  const authMutation = useMutation({
    mutationFn: async () => {
      if (isLogin) {
        const authData = await apiClient.login({ email, password });
        localStorage.setItem('token', authData.access_token);
        return authData;
      } else {
        // Register the auth account
        const authData = await apiClient.register({ email, password, slug: 'user' });
        
        const token = authData.access_token;
        // Immediately use the token to create a customer profile
        localStorage.setItem('token', token);
        // Also set axios default header exactly as requested
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
        
        // Parse the name field into first/last name
        const nameParts = name.trim().split(' ');
        const firstName = nameParts[0] || 'Unknown';
        const lastName = nameParts.slice(1).join(' ') || 'User';
        
        // Generate a random unique username
        const userName = email.split('@')[0].replace(/[^a-zA-Z0-9]/g, '') + Date.now();
        
        await apiClient.createCustomer({
          userName,
          firstName,
          lastName,
          bio: 'Food lover!'
        });
        
        return authData;
      }
    },
    onSuccess: () => {
      // Token is already set in mutationFn, but we force a reload to redirect to feed
      window.location.href = '/';
    },
    onError: (error: any) => {
      console.error("Auth error:", error);
      // Inline error UI handles displaying the message
    }
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    authMutation.mutate();
  };

  const handleGoogleLogin = () => {
    const clientId = import.meta.env.VITE_GOOGLE_CLIENT_ID as string | undefined;
    if (!clientId) {
      setOauthError('Google Sign-In is not configured. Missing VITE_GOOGLE_CLIENT_ID.');
      return;
    }
    setOauthError('');
    window.location.href = buildGoogleAuthUrl();
  };

  const handleGitHubLogin = () => {
    window.location.href = '/api/auth/github';
  };

  return (
    <div className={styles.container}>
      <div className={clsx(styles.card, 'glass-panel')}>
        <h1 className={styles.title}>{isLogin ? 'Welcome Back' : 'Join ShareBite'}</h1>
        <p className={styles.subtitle}>
          {isLogin ? 'Sign in to see what your friends are eating.' : 'Create an account to start sharing your food journey.'}
        </p>

        <form className={styles.form} onSubmit={handleSubmit}>
          {!isLogin && (
            <div className={styles.inputGroup}>
              <label htmlFor="name" className={styles.label}>Name</label>
              <input 
                type="text" 
                id="name" 
                className={styles.input} 
                placeholder="John Doe" 
                value={name}
                onChange={e => setName(e.target.value)}
              />
            </div>
          )}
          <div className={styles.inputGroup}>
            <label htmlFor="email" className={styles.label}>Email</label>
            <input 
              type="email" 
              id="email" 
              className={styles.input} 
              placeholder="john@example.com" 
              required
              value={email}
              onChange={e => setEmail(e.target.value)}
            />
          </div>
          <div className={styles.inputGroup}>
            <label htmlFor="password" className={styles.label}>Password</label>
            <input 
              type="password" 
              id="password" 
              className={styles.input} 
              placeholder="••••••••" 
              required
              value={password}
              onChange={e => setPassword(e.target.value)}
            />
          </div>

          <div className={clsx(styles.errorMessageWrapper, !authMutation.isError && styles.hidden)}>
            <AlertCircle size={16} className={styles.errorIcon} />
            <span className={styles.errorMessage}>
              {authMutation.error?.response?.data?.error || 
                (authMutation.error?.response?.status === 409 
                  ? "User already exists. Please try logging in." 
                  : "Authentication failed. Please check your inputs.")}
            </span>
          </div>

          <Button type="submit" fullWidth className={styles.submitBtn} disabled={authMutation.isPending}>
            {authMutation.isPending ? 'Loading...' : (isLogin ? 'Sign In' : 'Create Account')}
          </Button>
        </form>

        <div className={styles.divider}>
          <span className={styles.dividerText}>or continue with</span>
        </div>

        <div className={styles.oauthButtons}>
          <button type="button" className={styles.oauthBtn} onClick={handleGoogleLogin}>
            <svg className={styles.oauthIcon} viewBox="0 0 24 24" aria-hidden="true">
              <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" fill="#4285F4" />
              <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" />
              <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05" />
              <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" />
            </svg>
            Continue with Google
          </button>

          <button type="button" className={styles.oauthBtn} onClick={handleGitHubLogin}>
            <svg className={styles.oauthIcon} viewBox="0 0 24 24" aria-hidden="true" fill="currentColor">
              <path d="M12 0C5.37 0 0 5.37 0 12c0 5.3 3.438 9.8 8.205 11.387.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61-.546-1.385-1.335-1.755-1.335-1.755-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.418-1.305.762-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 21.795 24 17.295 24 12c0-6.63-5.37-12-12-12z" />
            </svg>
            Continue with GitHub
          </button>
        </div>

        {oauthError && (
          <div className={styles.oauthError}>{oauthError}</div>
        )}

        <div className={styles.toggle}>
          <span className={styles.toggleText}>
            {isLogin ? "Don't have an account?" : "Already have an account?"}
          </span>
          <button 
            type="button" 
            className={styles.toggleBtn}
            onClick={() => setIsLogin(!isLogin)}
          >
            {isLogin ? 'Sign Up' : 'Log In'}
          </button>
        </div>
      </div>
    </div>
  );
};
