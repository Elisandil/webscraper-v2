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
    if (!isAuthenticated) return;        setIsLoading(true);
        try {
            const { ok, data } = await apiRequest('/results');
            if (ok) {
                setResults(data.data || []);
            }
        } catch (error) {
            console.error('Error loading results:', error);
        } finally {
            setIsLoading(false);
        }
    }, [isAuthenticated]);

    // Toggle pagination mode
    const togglePaginationMode = useCallback(() => {
        setUsePagination(prev => {
            const newValue = !prev;
            // If switching to non-pagination mode, load results immediately
            if (!newValue) {
                // Use setTimeout to ensure state update happens first
                setTimeout(() => {
                    loadResults();
                }, 0);
            }
            return newValue;
        });
    }, [loadResults]);

    // Auto-refresh results every 15 seconds (when not in pagination mode and authenticated)
    useEffect(() => {
        if (usePagination || !isAuthenticated) return;

        loadResults();
        const interval = setInterval(loadResults, 15000);

        return () => clearInterval(interval);
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
