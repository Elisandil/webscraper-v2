import React, { createContext, useContext, useState, useCallback } from 'react';

const AlertContext = createContext(null);

export function AlertProvider({ children }) {
    const [alert, setAlert] = useState(null);

    const showSuccess = useCallback((message) => {
        setAlert({ type: 'success', message });
    }, []);

    const showError = useCallback((message) => {
        setAlert({ type: 'error', message });
    }, []);

    const showInfo = useCallback((message) => {
        setAlert({ type: 'info', message });
    }, []);

    const showWarning = useCallback((message) => {
        setAlert({ type: 'warning', message });
    }, []);

    const clearAlert = useCallback(() => {
        setAlert(null);
    }, []);

    const value = {
        alert,
        showSuccess,
        showError,
        showInfo,
        showWarning,
        clearAlert,
    };

    return <AlertContext.Provider value={value}>{children}</AlertContext.Provider>;
}

export function useAlert() {
    const context = useContext(AlertContext);
    if (!context) {
        throw new Error('useAlert must be used within an AlertProvider');
    }
    return context;
}
