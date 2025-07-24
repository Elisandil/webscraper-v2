import React, { useEffect, useState } from "react";
import LoginView from "./components/LoginView";
import RegisterModal from "./components/RegisterModal";
import ScrapeForm from "./components/ScrapeForm";
import ResultsList from "./components/ResultsList";
import DetailModal from "./components/DetailModal";
import HealthIndicator from "./components/HealthIndicator";
import Alert from "./components/Alert";
import { apiRequest } from "./api/client";

function App() {
  const [user, setUser] = useState(null);
  const [results, setResults] = useState([]);
  const [selected, setSelected] = useState(null);
  const [alert, setAlert] = useState(null);

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
    const iv = setInterval(() => apiRequest("/health"), 30000);
    return () => clearInterval(iv);
  }, []);

  const reloadResults = () =>
    apiRequest("/results").then(
      ({ ok, data }) => ok && setResults(data.data || [])
    );
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
        </header>

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
          onLogout={handleLogout}
        />

        <DetailModal result={selected} onClose={() => setSelected(null)} />
        <Alert alert={alert} onClose={() => setAlert(null)} />
      </div>
    </div>
  );
}

export default App;
