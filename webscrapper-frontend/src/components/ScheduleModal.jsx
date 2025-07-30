import React, { useState, useEffect } from "react";
import { apiRequest } from "../api/client";

export default function ScheduleModal({ isOpen, onClose, schedule, onSuccess, onError }) {
  const [formData, setFormData] = useState({
    name: "",
    url: "",
    cronExpr: "",
    active: true,
  });
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (schedule) {
      setFormData({
        name: schedule.name || "",
        url: schedule.url || "",
        cronExpr: schedule.cron_expression || "",
        active: schedule.active !== undefined ? schedule.active : true,
      });
    } else {
      setFormData({
        name: "",
        url: "",
        cronExpr: "",
        active: true,
      });
    }
  }, [schedule]);

  // Escuchar evento para pre-llenar URL desde ScrapeForm
  useEffect(() => {
    const handleOpenScheduleModal = (event) => {
      const { url } = event.detail;
      setFormData(prev => ({
        ...prev,
        url: url || "",
        name: url ? `Schedule para ${new URL(url).hostname}` : ""
      }));
    };

    window.addEventListener('openScheduleModal', handleOpenScheduleModal);
    return () => {
      window.removeEventListener('openScheduleModal', handleOpenScheduleModal);
    };
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const requestData = {
        name: formData.name.trim(),
        url: formData.url.trim(),
        cron_expression: formData.cronExpr.trim(),
        active: formData.active,
      };

      const url = schedule ? `/schedules/${schedule.id}` : "/schedules";
      const method = schedule ? "PUT" : "POST";

      const { ok, data } = await apiRequest(url, {
        method,
        body: JSON.stringify(requestData),
      });

      if (ok) {
        onSuccess(
          schedule
            ? "Schedule actualizado correctamente"
            : "Schedule creado correctamente"
        );
        // Disparar evento para recargar la lista
        window.dispatchEvent(new Event('reloadSchedules'));
        onClose();
      } else {
        onError(data.error || "Error al procesar el schedule");
      }
    } catch (error) {
      console.error("Error submitting schedule:", error);
      onError("Error de conexión");
    } finally {
      setIsLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!schedule || !window.confirm("¿Estás seguro de que quieres eliminar este schedule?")) {
      return;
    }

    setIsLoading(true);
    try {
      const { ok, data } = await apiRequest(`/schedules/${schedule.id}`, {
        method: "DELETE",
      });

      if (ok) {
        onSuccess("Schedule eliminado correctamente");
        // Disparar evento para recargar la lista
        window.dispatchEvent(new Event('reloadSchedules'));
        onClose();
      } else {
        onError(data.error || "Error al eliminar el schedule");
      }
    } catch (error) {
      console.error("Error deleting schedule:", error);
      onError("Error de conexión");
    } finally {
      setIsLoading(false);
    }
  };

  const cronExamples = [
    { value: "0 */10 * * * *", label: "Cada 10 minutos" },
    { value: "0 0 */1 * * *", label: "Cada hora" },
    { value: "0 0 8 * * *", label: "Diario a las 8:00 AM" },
    { value: "0 0 8 * * 1", label: "Lunes a las 8:00 AM" },
    { value: "0 0 0 1 * *", label: "Primer día del mes a medianoche" },
  ];

  if (!isOpen) return null;

  return (
    <div 
      className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      onClick={(e) => e.target === e.currentTarget && onClose()}
    >
      <div className="bg-black/80 backdrop-blur-lg rounded-2xl shadow-2xl border border-white/20 w-full max-w-md max-h-[90vh] overflow-y-auto">
        <div className="p-6 border-b border-white/20">
          <div className="flex items-center justify-between">
            <h2 className="text-2xl font-bold text-white">
              {schedule ? "Editar Schedule" : "Nuevo Schedule"}
            </h2>
            <button
              onClick={onClose}
              className="text-gray-400 hover:text-white transition-colors p-1"
              disabled={isLoading}
              aria-label="Cerrar"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Nombre del Schedule
            </label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
              placeholder="Ej: Scraping diario de noticias"
              required
              disabled={isLoading}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              URL a scrapear
            </label>
            <input
              type="url"
              value={formData.url}
              onChange={(e) => setFormData({ ...formData, url: e.target.value })}
              className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
              placeholder="https://ejemplo.com"
              required
              disabled={isLoading}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Expresión Cron
            </label>
            <input
              type="text"
              value={formData.cronExpr}
              onChange={(e) => setFormData({ ...formData, cronExpr: e.target.value })}
              className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 font-mono text-sm"
              placeholder="0 0 8 * * *"
              required
              disabled={isLoading}
            />
            <p className="text-xs text-gray-400 mt-1">
              Formato: segundo minuto hora día mes día_semana
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-3">
              Ejemplos comunes
            </label>
            <div className="grid grid-cols-1 gap-2">
              {cronExamples.map((example, index) => (
                <button
                  key={index}
                  type="button"
                  onClick={() => setFormData({ ...formData, cronExpr: example.value })}
                  className="text-left p-3 bg-white/5 hover:bg-white/10 rounded-lg transition-colors text-sm disabled:opacity-50"
                  disabled={isLoading}
                >
                  <code className="text-purple-400 font-mono text-xs block mb-1">{example.value}</code>
                  <span className="text-gray-300">{example.label}</span>
                </button>
              ))}
            </div>
          </div>

          <div className="flex items-center pt-2">
            <input
              type="checkbox"
              id="active"
              checked={formData.active}
              onChange={(e) => setFormData({ ...formData, active: e.target.checked })}
              className="w-4 h-4 text-blue-600 bg-white/10 border-white/20 rounded focus:ring-blue-500 focus:ring-2"
              disabled={isLoading}
            />
            <label htmlFor="active" className="ml-2 text-sm font-medium text-gray-300">
              Schedule activo
            </label>
          </div>

          <div className="flex gap-3 pt-6">
            <button
              type="submit"
              disabled={isLoading}
              className="flex-1 bg-blue-600/90 hover:bg-blue-600 text-white font-medium py-3 px-4 rounded-lg transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 shadow-lg backdrop-blur-sm border border-blue-500/20"
            >
              {isLoading ? (
                <>
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                  Procesando...
                </>
              ) : (
                <>
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
                  </svg>
                  {schedule ? "Actualizar" : "Crear Schedule"}
                </>
              )}
            </button>

            {schedule && (
              <button
                type="button"
                onClick={handleDelete}
                disabled={isLoading}
                className="px-4 py-3 bg-red-600/90 hover:bg-red-600 text-white font-medium rounded-lg transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed shadow-lg backdrop-blur-sm border border-red-500/20"
                title="Eliminar schedule"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            )}
          </div>
        </form>
      </div>
    </div>
  );
}