import React, { useState } from "react";
import { apiRequest } from "../api/client";

export default function LoginView({ onLogin, onAlert }) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const submit = async (e) => {
    e.preventDefault();
    setError("");
    const { ok, data } = await apiRequest("/auth/login", {
      method: "POST",
      body: JSON.stringify({ username, password }),
    });
    if (ok) {
      const token = data?.data?.token;
      if (token) {
        localStorage.setItem("jwtToken", token);
        onLogin();
      } else {
        setError("No se recibió token de autenticación");
      }
    }
  };

  return (
    <div
      id="loginView"
      className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900"
    >
      <form
        onSubmit={submit}
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
            className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
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
            className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
            placeholder="Ingresa tu contraseña"
          />
        </div>

        <button
          type="submit"
          id="loginBtn"
          className="w-full py-3 bg-green-600/90 hover:bg-green-600 text-white font-medium rounded-lg transition-all duration-200 shadow-lg backdrop-blur-sm border border-green-500/20"
        >
          Iniciar Sesión
        </button>

        <div className="text-center text-gray-400 mt-4">
          <span>¿No tienes cuenta? </span>
          <button
            id="openRegisterBtn"
            type="button"
            className="text-blue-400 hover:text-blue-300 hover:underline transition-colors"
          >
            Regístrate
          </button>
        </div>

        {error && (
          <div id="loginError" className="text-red-400 text-center text-sm bg-red-500/10 border border-red-500/20 rounded-lg p-3">
            {error}
          </div>
        )}
      </form>
    </div>
  );
}