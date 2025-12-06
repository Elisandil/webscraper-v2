import React from "react";
import { AppProviders } from "./contexts";
import { useAuth } from "./contexts/AuthContext";
import { useAlert } from "./contexts/AlertContext";
import LoginView from "./components/LoginView";
import RegisterModal from "./components/RegisterModal";
import MainView from "./components/MainView";
import Alert from "./components/Alert";

function AppContent() {
  const { isAuthenticated, isLoading, logout } = useAuth();
  const { alert, clearAlert, showSuccess } = useAlert();

  const handleLogout = () => {
    logout();
    showSuccess("Sesión cerrada correctamente");
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <>
        <LoginView />
        <RegisterModal />
        <Alert alert={alert} onClose={clearAlert} />
      </>
    );
  }

  return (
    <>
      <MainView handleLogout={handleLogout} />
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