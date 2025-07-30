import React, { useState, useEffect } from "react";
import { apiRequest } from "../api/client";

export default function ScheduleList({ onEdit, onAlert }) {
  const [schedules, setSchedules] = useState([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadSchedules();
  }, []);

  const loadSchedules = async () => {
    setIsLoading(true);
    
    try {
      const { ok, data } = await apiRequest("/schedules");
      
      if (ok) {
        setSchedules(data.data || []);
      } else {
        onAlert && onAlert({ type: "error", message: "Error al cargar schedules" });
      }
    } catch (error) {
      console.error("Error loading schedules:", error);
      onAlert && onAlert({ type: "error", message: "Error al cargar schedules" });
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    const handleReload = () => {
      loadSchedules();
    };

    window.addEventListener('reloadSchedules', handleReload);
    return () => {
      window.removeEventListener('reloadSchedules', handleReload);
    };
  }, []);

  const toggleScheduleStatus = async (schedule) => {
    try {
      const { ok, data } = await apiRequest(`/schedules/${schedule.id}`, {
        method: 'PUT',
        body: JSON.stringify({
          name: schedule.name,
          url: schedule.url,
          cron_expression: schedule.cron_expression,
          active: !schedule.active,
        }),
      });

      if (ok) {
        onAlert && onAlert({ 
          type: "success", 
          message: `Schedule ${!schedule.active ? 'activado' : 'pausado'} correctamente` 
        });
        loadSchedules();
      } else {
        onAlert && onAlert({ type: "error", message: data.error || "Error al cambiar estado" });
      }
    } catch (error) {
      console.error("Error toggling schedule:", error);
      onAlert && onAlert({ type: "error", message: "Error de conexión" });
    }
  };

  const deleteSchedule = async (schedule) => {
    if (!window.confirm("¿Estás seguro de que quieres eliminar este schedule?")) {
      return;
    }

    try {
      const { ok, data } = await apiRequest(`/schedules/${schedule.id}`, {
        method: 'DELETE',
      });

      if (ok) {
        onAlert && onAlert({ 
          type: "success", 
          message: "Schedule eliminado correctamente" 
        });
        loadSchedules();
      } else {
        onAlert && onAlert({ type: "error", message: data.error || "Error al eliminar schedule" });
      }
    } catch (error) {
      console.error("Error deleting schedule:", error);
      onAlert && onAlert({ type: "error", message: "Error de conexión" });
    }
  };

  const getStatusBadge = (active) => {
    return (
      <span
        className={`px-2 py-1 text-xs font-medium rounded-full ${
          active
            ? "bg-green-500/20 text-green-400 border border-green-500/30"
            : "bg-red-500/20 text-red-400 border border-red-500/30"
        }`}
      >
        {active ? "Activo" : "Inactivo"}
      </span>
    );
  };

  const formatDate = (dateString) => {
    if (!dateString) return "Nunca";
    return new Date(dateString).toLocaleString("es-ES", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const getNextRunStatus = (nextRun, active) => {
    if (!active) return { text: "Pausado", color: "bg-gray-500/20 text-gray-400" };
    if (!nextRun) return { text: "No programado", color: "bg-gray-500/20 text-gray-400" };
    
    const now = new Date();
    const next = new Date(nextRun);
    const diff = next - now;
    
    if (diff < 0) return { text: "Pendiente", color: "bg-orange-500/20 text-orange-400" };
    if (diff < 60000) return { text: "En breve", color: "bg-yellow-500/20 text-yellow-400" };
    if (diff < 3600000) return { text: `${Math.floor(diff / 60000)}m`, color: "bg-blue-500/20 text-blue-400" };
    if (diff < 86400000) return { text: `${Math.floor(diff / 3600000)}h`, color: "bg-blue-500/20 text-blue-400" };
    return { text: `${Math.floor(diff / 86400000)}d`, color: "bg-blue-500/20 text-blue-400" };
  };

  const cronToHuman = (cronExpr) => {
    const patterns = {
      "0 */10 * * * *": "Cada 10 minutos",
      "0 0 */1 * * *": "Cada hora",
      "0 0 8 * * *": "Diario a las 8:00 AM",
      "0 0 8 * * 1": "Lunes a las 8:00 AM",
      "0 0 0 1 * *": "Primer día del mes",
    };
    return patterns[cronExpr] || cronExpr;
  };

  if (isLoading) {
    return (
      <div className="bg-black/80 backdrop-blur-lg rounded-2xl shadow-2xl border border-white/20 p-6">
        <div className="flex items-center justify-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-400"></div>
          <span className="ml-3 text-gray-300">Cargando schedules...</span>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-black/80 backdrop-blur-lg rounded-2xl shadow-2xl border border-white/20 overflow-hidden">
      <div className="px-6 py-4 border-b border-white/20">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-bold text-white">Schedules Programados</h2>
            <p className="text-gray-400 text-sm mt-1">
              {schedules.length} schedule{schedules.length !== 1 ? "s" : ""} configurado{schedules.length !== 1 ? "s" : ""}
            </p>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
            <span className="text-sm text-gray-400">Scheduler activo</span>
          </div>
        </div>
      </div>

      {schedules.length === 0 ? (
        <div className="p-8 text-center">
          <div className="w-16 h-16 mx-auto mb-4 bg-white/10 rounded-full flex items-center justify-center">
            <svg
              className="w-8 h-8 text-gray-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          </div>
          <h3 className="text-lg font-medium text-gray-300 mb-2">
            No hay schedules configurados
          </h3>
          <p className="text-gray-400 text-sm">
            Crea tu primer schedule para automatizar el web scraping
          </p>
        </div>
      ) : (
        <div className="divide-y divide-white/20">
          {schedules.map((schedule) => {
            const nextRunStatus = getNextRunStatus(schedule.next_run, schedule.active);
            
            return (
              <div
                key={schedule.id}
                className="p-6 hover:bg-white/5 transition-all duration-200"
              >
                <div className="flex items-start justify-between">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-3 mb-2">
                      <h3 className="text-lg font-semibold text-white truncate">
                        {schedule.name}
                      </h3>
                      {getStatusBadge(schedule.active)}
                    </div>
                    
                    <p className="text-blue-400 text-sm mb-3 truncate">
                      {schedule.url}
                    </p>
                    
                    <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 text-sm">
                      <div>
                        <span className="text-gray-400 block mb-1">Programación</span>
                        <div className="text-white">
                          {cronToHuman(schedule.cron_expression)}
                        </div>
                        <code className="text-purple-400 font-mono text-xs">
                          {schedule.cron_expression}
                        </code>
                      </div>
                      
                      <div>
                        <span className="text-gray-400 block mb-1">Ejecuciones</span>
                        <span className="text-white font-medium">
                          {schedule.run_count || 0}
                        </span>
                      </div>
                      
                      <div>
                        <span className="text-gray-400 block mb-1">Última ejecución</span>
                        <span className="text-white">
                          {formatDate(schedule.last_run)}
                        </span>
                      </div>
                      
                      <div>
                        <span className="text-gray-400 block mb-1">Próxima ejecución</span>
                        <div className="flex flex-col gap-1">
                          <span className="text-white text-sm">
                            {formatDate(schedule.next_run)}
                          </span>
                          <span className={`px-2 py-0.5 text-xs rounded-full ${nextRunStatus.color}`}>
                            {nextRunStatus.text}
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                  
                  <div className="ml-4 flex-shrink-0 flex items-center gap-2">
                    {/* Botón de toggle activo/inactivo */}
                    <button
                      onClick={() => toggleScheduleStatus(schedule)}
                      className={`p-2 rounded-lg transition-colors ${
                        schedule.active
                          ? "text-green-400 hover:text-green-300 hover:bg-green-400/10"
                          : "text-gray-400 hover:text-green-400 hover:bg-green-400/10"
                      }`}
                      title={schedule.active ? "Pausar schedule" : "Activar schedule"}
                    >
                      {schedule.active ? (
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                      ) : (
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1m4 0h1m-6 4h.01M15 14h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                      )}
                    </button>
                    
                    {/* Botón de editar */}
                    <button
                      onClick={() => onEdit && onEdit(schedule)}
                      className="p-2 text-gray-400 hover:text-white hover:bg-white/10 rounded-lg transition-colors"
                      title="Editar schedule"
                    >
                      <svg
                        className="w-5 h-5"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth="2"
                          d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                        />
                      </svg>
                    </button>

                    {/* Botón de eliminar */}
                    <button
                      onClick={() => deleteSchedule(schedule)}
                      className="p-2 text-gray-400 hover:text-red-400 hover:bg-red-400/10 rounded-lg transition-colors"
                      title="Eliminar schedule"
                    >
                      <svg
                        className="w-5 h-5"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth="2"
                          d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                        />
                      </svg>
                    </button>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}