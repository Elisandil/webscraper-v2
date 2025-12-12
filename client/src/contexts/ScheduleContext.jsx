import React, { createContext, useContext, useState, useCallback, useEffect } from 'react';
import { apiRequest } from '../api/client';
import { useAuth } from './AuthContext';

const ScheduleContext = createContext(null);

export function ScheduleProvider({ children }) {
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [selectedSchedule, setSelectedSchedule] = useState(null);
    const [prefilledUrl, setPrefilledUrl] = useState(null);
    const [schedules, setSchedules] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const { isAuthenticated } = useAuth();

  const loadSchedules = useCallback(async () => {
        if (!isAuthenticated) return;

        setIsLoading(true);
        try {
            const { ok, data } = await apiRequest('/schedules');
            if (ok) {
                setSchedules(data.data || []);
            }
        } catch (error) {
            console.error('Error loading schedules:', error);
        } finally {
            setIsLoading(false);
        }
    }, [isAuthenticated]);

    useEffect(() => {
        if (isAuthenticated) {
            loadSchedules();
        }
    }, [isAuthenticated, loadSchedules]);

    const openScheduleModal = useCallback((url = null) => {
        setPrefilledUrl(url);
        setSelectedSchedule(null);
        setIsModalOpen(true);
    }, []);

    const editSchedule = useCallback((schedule) => {
        setSelectedSchedule(schedule);
        setPrefilledUrl(null);
        setIsModalOpen(true);
    }, []);

    const closeModal = useCallback(() => {
        setIsModalOpen(false);
        setSelectedSchedule(null);
        setPrefilledUrl(null);
    }, []);

    const value = {
        isModalOpen,
        selectedSchedule,
        prefilledUrl,
        schedules,
        isLoading,
        openScheduleModal,
        editSchedule,
        closeModal,
        refreshSchedules: loadSchedules,
    };

    return <ScheduleContext.Provider value={value}>{children}</ScheduleContext.Provider>;
}

export function useSchedule() {
    const context = useContext(ScheduleContext);
    if (!context) {
        throw new Error('useSchedule must be used within a ScheduleProvider');
    }
    return context;
}
