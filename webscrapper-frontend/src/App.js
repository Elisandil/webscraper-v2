import React, { useEffect, useState, useCallback } from "react";
import LoginView from "./components/LoginView";
import RegisterModal from "./components/RegisterModal";
import MainView from "./components/MainView";
import Alert from "./components/Alert";
import { apiRequest } from "./api/client";

function App() {
  const [user, setUser] = useState(null);
  const [results, setResults] = useState([]);
  const [selected, setSelected] = useState(null);
  const [alert, setAlert] = useState(null);
  const [activeTab, setActiveTab] = useState("scraping");
  const [usePagination, setUsePagination] = useState(false);
  
  // Función para recargar resultados (sin paginación - para compatibilidad)
  const reloadResults = useCallback(() => {
    if (!usePagination) {
      apiRequest("/results").then(({ ok, data }) => {
        if (ok) {
          setResults(data.data || []);
        }
      });
    }
  }, [usePagination]);

  useEffect(() => {
    const token = localStorage.getItem("jwtToken");
    if (token) {
      apiRequest("/results").then(({ ok, data }) => {
        if (ok) {
          setUser(true);
          // Solo cargar resultados si no usamos paginación
          if (!usePagination) {
            setResults(data.data || []);
          }
        } else {
          localStorage.removeItem("jwtToken");
          setUser(null);
        }
      });
    }
    
    if (!usePagination) {
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
    }
  }, [reloadResults, usePagination]);

  const handleLogout = () => {
    localStorage.removeItem("jwtToken");
    setUser(null);
    setResults([]);
    setSelected(null);
    setAlert({ type: "success", message: "Sesión cerrada correctamente" });
  };

  const togglePaginationMode = () => {
    setUsePagination(!usePagination);
    setAlert({ 
      type: "info", 
      message: `Modo ${!usePagination ? 'paginación' : 'lista completa'} activado` 
    });
  };

  if (!user) {
    return (
      <>
        <LoginView
          onLogin={() => {
            setUser(true);
            if (!usePagination) {
              reloadResults();
            }
          }}
          onAlert={setAlert}
        />
        <RegisterModal onAlert={setAlert} />
        <Alert alert={alert} onClose={() => setAlert(null)} />
      </>
    );
  }

  return (
    <MainView
      user={user}
      results={results}
      selected={selected}
      alert={alert}
      activeTab={activeTab}
      usePagination={usePagination}
      setSelected={setSelected}
      setAlert={setAlert}
      setActiveTab={setActiveTab}
      reloadResults={reloadResults}
      handleLogout={handleLogout}
      togglePaginationMode={togglePaginationMode}
    />
  );
}

export default App;