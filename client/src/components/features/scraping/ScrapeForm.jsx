import React, { useState } from "react";
import { apiRequest } from "../../../api/client";
import { useAlert } from "../../../contexts/AlertContext";
import { useSchedule } from "../../../contexts/ScheduleContext";
import { useResults } from "../../../contexts/ResultsContext";

export default function ScrapeForm() {
  const { showSuccess, showError } = useAlert();
  const { openScheduleModal } = useSchedule();
  const { usePagination, loadResults } = useResults();
  const [url, setUrl] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [showScheduleOption, setShowScheduleOption] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!url.trim()) return;

    setIsLoading(true);
    try {
      const { ok, data } = await apiRequest("/scrape", {
        method: "POST",
        body: JSON.stringify({ url: url.trim() }),
      });

      if (ok) {
        showSuccess("URL scrapeada exitosamente");
        setUrl("");
        setShowScheduleOption(true);
        setTimeout(() => setShowScheduleOption(false), 5000);

        // Reload results if not in pagination mode
        if (!usePagination) {
          loadResults();
        }
      } else {
        showError(data.error || "Error al scrapear la URL");
      }
    } catch (error) {
      console.error("Error:", error);
      showError("Error de conexión");
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateSchedule = () => {
    openScheduleModal(url.trim());
    setShowScheduleOption(false);
  };

  return (
    <div className="mb-8">
      <div className="bg-black/80 backdrop-blur-lg rounded-2xl shadow-2xl border border-white/20 p-6">
        <div className="mb-6">
          <h2 className="text-2xl font-bold text-white mb-2">Scraping Manual</h2>
          <p className="text-gray-400 text-sm">
            Ingresa una URL para extraer información inmediatamente
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="flex gap-4">
            <div className="flex-1">
              <input
                type="url"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                placeholder="https://ejemplo.com"
                className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-cyan-500 focus:border-transparent transition-all duration-200"
                required
                disabled={isLoading}
              />
            </div>

            <button
              type="submit"
              disabled={isLoading || !url.trim()}
              className="px-8 py-3 bg-gradient-to-r from-cyan-600 to-teal-600 hover:from-cyan-500 hover:to-teal-500 text-white font-medium rounded-lg transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2 shadow-lg shadow-cyan-500/30 backdrop-blur-sm border border-cyan-500/20"
            >
              {isLoading ? (
                <>
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                  Scrapeando...
                </>
              ) : (
                <>
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                  Scrap
                </>
              )}
            </button>
          </div>

          {/* Opción para crear schedule después de un scraping exitoso */}
          {showScheduleOption && (
            <div className="bg-gradient-to-r from-teal-500/10 to-cyan-500/10 border border-teal-500/20 rounded-lg p-4 animate-pulse">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 bg-teal-500/20 rounded-full flex items-center justify-center">
                    <svg className="w-4 h-4 text-teal-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                  <div>
                    <p className="text-white font-medium">¿Quieres automatizar este scraping?</p>
                    <p className="text-gray-400 text-sm">Crea un schedule para scrapear esta URL automáticamente</p>
                  </div>
                </div>

                <div className="flex items-center gap-2">
                  <button
                    type="button"
                    onClick={handleCreateSchedule}
                    className="px-4 py-2 bg-cyan-600/90 hover:bg-cyan-600 text-white text-sm font-medium rounded-lg transition-all duration-200 shadow-lg shadow-cyan-500/20 backdrop-blur-sm border border-cyan-500/20"
                  >
                    Crear Schedule
                  </button>
                  <button
                    type="button"
                    onClick={() => setShowScheduleOption(false)}
                    className="text-gray-400 hover:text-white p-1 transition-colors"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          )}
        </form>

        <div className="mt-6 flex items-center gap-6 text-sm text-gray-400">
          <div className="flex items-center gap-2">
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>Soporta HTTP y HTTPS</span>
          </div>
          <div className="flex items-center gap-2">
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>Extrae metadata, links e imágenes</span>
          </div>
        </div>
      </div>
    </div>
  );
}