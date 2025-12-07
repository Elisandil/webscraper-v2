import React from 'react';
import { AuthProvider } from './AuthContext';
import { AlertProvider } from './AlertContext';
import { ResultsProvider } from './ResultsContext';
import { ScheduleProvider } from './ScheduleContext';
import { ChatProvider } from './ChatContext';

export function AppProviders({ children }) {
    return (
        <AuthProvider>
            <AlertProvider>
                <ResultsProvider>
                    <ScheduleProvider>
                        <ChatProvider>
                            {children}
                        </ChatProvider>
                    </ScheduleProvider>
                </ResultsProvider>
            </AlertProvider>
        </AuthProvider>
    );
}
