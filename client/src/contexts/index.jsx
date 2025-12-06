import React from 'react';
import { AuthProvider } from './AuthContext';
import { AlertProvider } from './AlertContext';
import { ResultsProvider } from './ResultsContext';
import { ScheduleProvider } from './ScheduleContext';

export function AppProviders({ children }) {
    return (
        <AuthProvider>
            <AlertProvider>
                <ResultsProvider>
                    <ScheduleProvider>
                        {children}
                    </ScheduleProvider>
                </ResultsProvider>
            </AlertProvider>
        </AuthProvider>
    );
}
