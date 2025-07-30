import React, { useEffect, useState, useCallback } from "react";
import LoginView from "./components/LoginView";
import RegisterModal from "./components/RegisterModal";
import ScrapeForm from "./components/ScrapeForm";
import ResultsList from "./components/ResultsList";
import ScheduleSection from "./components/ScheduleSection";
import DetailModal from "./components/DetailModal";
import HealthIndicator from "./components/HealthIndicator";
import Alert from "./components/Alert";
import { apiRequest } from "./api/client";

function App() {
  const [user, setUser] = useState(null);
  const [results, setResults] = useState([]);
  const [selected, setSelected] = useState(null);
  const [alert, setAlert] = useState(null);
  const [activeTab, setActiveTab] = useState("scraping");
  const reloadResults = useCallback(() => {
    apiRequest("/results").then(({ ok, data }) => {
      if (ok) {
        setResults(data.data || []);
      }
    });
  }, []);

  useEffect(() => {
    const token = localStorage.getItem("jwtToken");
    if (token) {
      apiRequest("/results").then(({ ok, data }) => {
        if (ok) {
          setUser(true);
          setResults(data.data || []);
        } else {
          localStorage.removeItem("jwtToken");
          setUser(null);
        }
      });
    }
    reloadResults();
    const pollInterval = setInterval(reloadResults, 15000);

    const handleSwitchToSchedules = () => {
      setActiveTab("schedules");
    };
    window.addEventListener('switchToSchedules', handleSwitchToSchedules);

    return () => {
      clearInterval(pollInterval);
      window.removeEventListener('switchToSchedules', handleSwitchToSchedules);
    };
  }, [reloadResults]);

  const handleLogout = () => {
    localStorage.removeItem("jwtToken");
    setUser(null);
    setResults([]);
    setSelected(null);
    setAlert({ type: "success", message: "Sesión cerrada correctamente" });
  };

  if (!user) {
    return (
      <>
        <LoginView
          onLogin={() => {
            setUser(true);
            reloadResults();
          }}
          onAlert={setAlert}
        />
        <RegisterModal onAlert={setAlert} />
        <Alert alert={alert} onClose={() => setAlert(null)} />
      </>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-gray-100 overflow-x-hidden">
      <HealthIndicator />

      <div className="container mx-auto px-6 py-8">
        <header className="flex items-center justify-between mb-8">
          <h1 className="text-3xl font-bold text-gray-100">WebScraper 1.0</h1>
          <button
            onClick={handleLogout}
            className="text-gray-400 hover:text-white transition-colors flex items-center gap-2"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
            Cerrar Sesión
          </button>
        </header>

        {/* Navegación por pestañas */}
        <div className="mb-8">
          <nav className="flex space-x-8 border-b border-white/20">
            <button
              onClick={() => setActiveTab("scraping")}
              className={`py-4 px-1 border-b-2 font-medium text-sm transition-colors ${
                activeTab === "scraping"
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
              className={`py-4 px-1 border-b-2 font-medium text-sm transition-colors ${
                activeTab === "schedules"
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
            <ScrapeForm
              onSuccess={(msg) => {
                setAlert({ type: "success", message: msg });
                reloadResults();
              }}
              onError={(msg) => setAlert({ type: "error", message: msg })}
            />

            <ResultsList
              results={results}
              onView={(r) => setSelected(r)}
              onDelete={(id) => {
                setAlert({
                  type: "success",
                  message: "Resultado eliminado con éxito",
                });
                reloadResults();
              }}
            />
          </>
        ) : (
          <ScheduleSection onAlert={setAlert} />
        )}

        <DetailModal result={selected} onClose={() => setSelected(null)} />
        <Alert alert={alert} onClose={() => setAlert(null)} />
      </div>
    </div>
  );
}

export default App;