import React, { createContext, useContext, useState, useCallback, useEffect } from 'react';
import { apiRequest } from '../api/client';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('jwtToken');
    if (token) {
      apiRequest('/results')
        .then(({ ok }) => {
          if (ok) {
            setUser(true);
          } else {
            localStorage.removeItem('jwtToken');
            setUser(null);
          }
        })
        .finally(() => setIsLoading(false));
    } else {
      setIsLoading(false);
    }
  }, []);

  const login = useCallback((token) => {
    localStorage.setItem('jwtToken', token);
    setUser(true);
  }, []);

  const logout = useCallback(async () => {
    const token = localStorage.getItem('jwtToken');
    if (token) {
      try {
        await apiRequest('/auth/logout', {
          method: 'POST',
        });
      } catch (error) {
        console.error('Logout error:', error);
      }
    }
    localStorage.removeItem('jwtToken');
    setUser(null);
  }, []);

  const value = {
    user,
    isAuthenticated: !!user,
    isLoading,
    login,
    logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
