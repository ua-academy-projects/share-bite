import React, { useState, useEffect } from 'react';
import { Button } from '../../components/Button/Button';
import styles from './Auth.module.css';
import { clsx } from 'clsx';
import { useNavigate, useLocation } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { apiClient } from '../../api/client';
import axios from 'axios';
import { AlertCircle } from 'lucide-react';

export const Auth: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const [isLogin, setIsLogin] = useState(location.state?.isLogin ?? true);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [name, setName] = useState('');

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
    onSuccess: (data) => {
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
