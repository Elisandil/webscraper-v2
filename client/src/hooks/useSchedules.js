import { useState, useEffect, useCallback } from 'react';
import { apiRequest } from '../api/client';

export function useSchedules() {
  const [schedules, setSchedules] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  const loadSchedules = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const { ok, data } = await apiRequest('/schedules');
      if (ok) {
        setSchedules(data.data || []);
      } else {
        setError(data.error || 'Error al cargar schedules');
      }
    } catch (err) {
      setError('Error de conexión al cargar schedules');
    } finally {
      setIsLoading(false);
    }
  }, []);

  const createSchedule = useCallback(async (scheduleData) => {
    try {
      const { ok, data } = await apiRequest('/schedules', {
        method: 'POST',
        body: JSON.stringify({
          name: scheduleData.name,
          url: scheduleData.url,
          cron_expression: scheduleData.cronExpr,
          active: scheduleData.active,
        }),
      });

      if (ok) {
        await loadSchedules();
        return { success: true, data: data.data };
      } else {
        return { success: false, error: data.error || 'Error al crear el schedule' };
      }
    } catch (err) {
      return { success: false, error: 'Error de conexión al crear el schedule' };
    }
  }, [loadSchedules]);

  const updateSchedule = useCallback(async (scheduleId, scheduleData) => {
    try {
      const { ok, data } = await apiRequest(`/schedules/${scheduleId}`, {
        method: 'PUT',
        body: JSON.stringify({
          name: scheduleData.name,
          url: scheduleData.url,
          cron_expression: scheduleData.cronExpr,
          active: scheduleData.active,
        }),
      });

      if (ok) {
        await loadSchedules();
        return { success: true, data: data.data };
      } else {
        return { success: false, error: data.error || 'Error al actualizar el schedule' };
      }
    } catch (err) {
      return { success: false, error: 'Error de conexión al actualizar el schedule' };
    }
  }, [loadSchedules]);

  const deleteSchedule = useCallback(async (scheduleId) => {
    try {
      const { ok, data } = await apiRequest(`/schedules/${scheduleId}`, {
        method: 'DELETE',
      });

      if (ok) {
        await loadSchedules();
        return { success: true };
      } else {
        return { success: false, error: data.error || 'Error al eliminar el schedule' };
      }
    } catch (err) {
      return { success: false, error: 'Error de conexión al eliminar el schedule' };
    }
  }, [loadSchedules]);

  const getSchedule = useCallback(async (scheduleId) => {
    try {
      const { ok, data } = await apiRequest(`/schedules/${scheduleId}`);
      if (ok) {
        return { success: true, data: data.data };
      } else {
        return { success: false, error: data.error || 'Error al obtener el schedule' };
      }
    } catch (err) {
      return { success: false, error: 'Error de conexión al obtener el schedule' };
    }
  }, []);

  useEffect(() => {
    loadSchedules();
  }, [loadSchedules]);

  useEffect(() => {
    const handleReload = () => {
      loadSchedules();
    };

    window.addEventListener('reloadSchedules', handleReload);
    return () => {
      window.removeEventListener('reloadSchedules', handleReload);
    };
  }, [loadSchedules]);

  return {
    schedules,
    isLoading,
    error,
    loadSchedules,
    createSchedule,
    updateSchedule,
    deleteSchedule,
    getSchedule,
  };
}