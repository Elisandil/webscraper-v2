import React, { createContext, useContext, useState, useCallback, useEffect } from 'react';
import { apiRequest } from '../api/client';
import { useAuth } from './AuthContext';

const ResultsContext = createContext(null);

export function ResultsProvider({ children }) {
    const [results, setResults] = useState([]);
    const [selected, setSelected] = useState(null);
    const [usePagination, setUsePagination] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const { isAuthenticated } = useAuth();

  const loadResults = useCallback(async () => {
    if (!isAuthenticated) return;
    setIsLoading(true);
    try {
      const { ok, data } = await apiRequest('/results?page=1&per_page=50');
      if (ok) {
        setResults(data.data?.data || []);
      }
    } catch (error) {
      console.error('Error loading results:', error);
    } finally {
      setIsLoading(false);
    }
  }, [isAuthenticated]);

    const togglePaginationMode = useCallback(() => {
        setUsePagination(prev => !prev);
    }, []);

    // Subscribe to SSE for real-time updates. Falls back to 60s polling if the
    // connection drops before EventSource auto-reconnects.
    useEffect(() => {
        if (usePagination || !isAuthenticated) return;

        loadResults();

        const es = new EventSource('/api/results/events', { withCredentials: true });
        es.onmessage = () => loadResults();

        // Slow fallback in case SSE is temporarily disconnected
        const fallback = setInterval(loadResults, 60000);

        return () => {
            es.close();
            clearInterval(fallback);
        };
    }, [usePagination, isAuthenticated, loadResults]);

    const value = {
        results,
        selected,
        usePagination,
        isLoading,
        setResults,
        setSelected,
        togglePaginationMode,
        loadResults,
    };

    return <ResultsContext.Provider value={value}>{children}</ResultsContext.Provider>;
}

export function useResults() {
    const context = useContext(ResultsContext);
    if (!context) {
        throw new Error('useResults must be used within a ResultsProvider');
    }
    return context;
}
