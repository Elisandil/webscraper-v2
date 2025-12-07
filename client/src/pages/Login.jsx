import React, { useState } from "react";
import { apiRequest } from "../api/client";
import { useAuth } from "../contexts/AuthContext";
import { useAlert } from "../contexts/AlertContext";

export default function LoginView() {
  const { login } = useAuth();
  const { showSuccess, showError } = useAlert();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleLogin = async (e) => {
    e.preventDefault();
    if (!username.trim() || !password.trim()) {
      showError("Por favor completa todos los campos");
      return;
    }

    setIsLoading(true);
    try {
      const { ok, data } = await apiRequest("/auth/login", {
        method: "POST",
        body: JSON.stringify({ username, password }),
      });

      if (ok && data.data?.token) {
        login(data.data.token);
        showSuccess("Inicio de sesión exitoso");
      } else {
        showError(data.error || "Credenciales inválidas");
      }
    } catch (error) {
      console.error("Error:", error);
      showError("Error de conexión");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div
      id="loginView"
      className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900"
    >
      <form
        onSubmit={handleLogin}
        id="loginForm"
        className="bg-black/80 backdrop-blur-lg rounded-2xl shadow-2xl border border-white/20 w-full max-w-sm p-8 space-y-6"
      >
        <h2 className="text-2xl font-bold text-center text-white">
          Iniciar Sesión
        </h2>

        <div>
          <label htmlFor="username" className="block mb-2 text-sm font-medium text-gray-300">
            Usuario
          </label>
          <input
            id="username"
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
            disabled={isLoading}
            className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-cyan-500 focus:border-transparent transition-all duration-200"
            placeholder="Ingresa tu usuario"
          />
        </div>

        <div>
          <label htmlFor="password" className="block mb-2 text-sm font-medium text-gray-300">
            Contraseña
          </label>
          <input
            id="password"
            type="password"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            disabled={isLoading}
            className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-cyan-500 focus:border-transparent transition-all duration-200"
            placeholder="Ingresa tu contraseña"
          />
        </div>

        <button
          type="submit"
          id="loginBtn"
          disabled={isLoading}
          className="w-full py-3 bg-gradient-to-r from-cyan-600 to-teal-600 hover:from-cyan-500 hover:to-teal-500 text-white font-medium rounded-lg transition-all duration-200 shadow-lg shadow-cyan-500/30 backdrop-blur-sm border border-cyan-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isLoading ? "Iniciando sesión..." : "Iniciar Sesión"}
        </button>

        <div className="text-center text-gray-400 mt-4">
          <span>¿No tienes cuenta? </span>
          <button
            id="openRegisterBtn"
            type="button"
            className="text-cyan-400 hover:text-cyan-300 hover:underline transition-colors"
          >
            Regístrate
          </button>
        </div>
      </form>
    </div>
  );
}