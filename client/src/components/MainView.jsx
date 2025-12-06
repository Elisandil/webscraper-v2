import React, { useState } from "react";
import ScrapeForm from "./ScrapeForm";
import ResultsList from "./ResultsList";
import PaginatedResultsList from "./PaginatedResultsList";
import ScheduleSection from "./ScheduleSection";
import DetailModal from "./DetailModal";
import HealthIndicator from "./HealthIndicator";
import { useAlert } from "../contexts/AlertContext";
import { useResults } from "../contexts/ResultsContext";

export default function MainView({ handleLogout }) {
  const [activeTab, setActiveTab] = useState("scraping");
  const { showInfo } = useAlert();
  const { results, selected, setSelected, usePagination, togglePaginationMode, loadResults } = useResults();

  const handleTogglePagination = () => {
    togglePaginationMode();
    showInfo(`Modo ${!usePagination ? 'paginación' : 'lista completa'} activado`);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-gray-100 overflow-x-hidden">
      <HealthIndicator />

      <div className="container mx-auto px-6 py-8">
        <header className="flex items-center justify-between mb-8">
          <h1 className="text-3xl font-bold text-gray-100">WebScraper 1.0</h1>

          <div className="flex items-center space-x-4">
            {/* Toggle para paginación (solo visible en tab scraping) */}
            {activeTab === "scraping" && (
              <div className="flex items-center space-x-3">
                <span className="text-sm text-gray-400">Paginación:</span>
                <button
                  onClick={handleTogglePagination}
                  className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${usePagination ? "bg-blue-600" : "bg-gray-600"
                    }`}
                  title={`${usePagination ? 'Desactivar' : 'Activar'} paginación`}
                >
                  <span
                    className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${usePagination ? "translate-x-6" : "translate-x-1"
                      }`}
                  />
                </button>
                <span className="text-xs text-gray-500">
                  {usePagination ? 'Paginado' : 'Lista completa'}
                </span>
              </div>
            )}

            <button
              onClick={handleLogout}
              className="px-6 py-3 bg-red-600/90 hover:bg-red-600 text-white font-medium rounded-lg transition-all duration-200 flex items-center gap-2 shadow-lg backdrop-blur-sm border border-red-500/20"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              Cerrar Sesión
            </button>
          </div>
        </header>

        {/* Navegación por pestañas */}
        <div className="mb-8">
          <nav className="flex space-x-8 border-b border-white/20">
            <button
              onClick={() => setActiveTab("scraping")}
              className={`py-4 px-1 border-b-2 font-medium text-sm transition-colors ${activeTab === "scraping"
                ? "border-blue-500 text-blue-400"
                : "border-transparent text-gray-400 hover:text-gray-300 hover:border-gray-300"
                }`}
            >
              <div className="flex items-center gap-2">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
                Scraping Manual
              </div>
            </button>
            <button
              onClick={() => setActiveTab("schedules")}
              className={`py-4 px-1 border-b-2 font-medium text-sm transition-colors ${activeTab === "schedules"
                ? "border-blue-500 text-blue-400"
                : "border-transparent text-gray-400 hover:text-gray-300 hover:border-gray-300"
                }`}
            >
              <div className="flex items-center gap-2">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Schedules Automáticos
              </div>
            </button>
          </nav>
        </div>

        {/* Contenido según la pestaña activa */}
        {activeTab === "scraping" ? (
          <>
            <ScrapeForm />

            {/* Renderizar lista según el modo */}
            {usePagination ? (
              <PaginatedResultsList onView={(r) => setSelected(r)} />
            ) : (
              <ResultsList
                results={results}
                onView={(r) => setSelected(r)}
                onDelete={(id) => {
                  loadResults();
                }}
              />
            )}
          </>
        ) : (
          <ScheduleSection />
        )}

        <DetailModal result={selected} onClose={() => setSelected(null)} />
      </div>
    </div>
  );
}