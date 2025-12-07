import React, { useState } from "react";
import { AppProviders } from "./contexts";
import { useAuth } from "./contexts/AuthContext";
import { useAlert } from "./contexts/AlertContext";
import Landing from "./pages/Landing";
import Login from "./pages/Login";
import Dashboard from "./pages/Dashboard";
import RegisterModal from "./components/modals/RegisterModal";
import Alert from "./components/ui/Alert";

function AppContent() {
  const { isAuthenticated, isLoading, logout } = useAuth();
  const { alert, clearAlert, showSuccess } = useAlert();
  const [showLanding, setShowLanding] = useState(true);

  const handleLogout = () => {
    logout();
    showSuccess("Sesión cerrada correctamente");
  };

  const handleGetStarted = () => {
    setShowLanding(false);
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-cyan-500 shadow-lg shadow-cyan-500/40"></div>
      </div>
    );
  }

  if (!isAuthenticated) {
    if (showLanding) {
      return <Landing onGetStarted={handleGetStarted} />;
    }
    
    return (
      <>
        <Login />
        <RegisterModal />
        <Alert alert={alert} onClose={clearAlert} />
      </>
    );
  }

  return (
    <>
      <Dashboard handleLogout={handleLogout} />
      <Alert alert={alert} onClose={clearAlert} />
    </>
  );
}

function App() {
  return (
    <AppProviders>
      <AppContent />
    </AppProviders>
  );
}

export default App;