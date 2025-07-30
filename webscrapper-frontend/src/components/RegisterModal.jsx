import React, { useState, useEffect } from "react";
import { apiRequest } from "../api/client";

export default function RegisterModal({ onAlert }) {
  const [open, setOpen] = useState(false);
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    const btn = document.getElementById("openRegisterBtn");
    if (btn) btn.addEventListener("click", () => setOpen(true));
  }, []);

  const close = () => {
    setOpen(false);
    setError("");
    setUsername("");
    setEmail("");
    setPassword("");
  };

  const submit = async (e) => {
    e.preventDefault();
    if (!username || !email || password.length < 6) {
      setError(
        "Todos los campos son requeridos (mín. 6 caracteres en contraseña)."
      );
      return;
    }
    const { ok, data } = await apiRequest("/auth/register", {
      method: "POST",
      body: JSON.stringify({ username, email, password }),
    });
    if (ok) {
      onAlert({
        type: "success",
        message: "Registro exitoso. Por favor inicia sesión.",
      });
      close();
    } else {
      const msg = data.message || data.error || "Error al registrar.";
      setError(msg);
    }
  };

  if (!open) return null;
  
  return (
    <div
      id="registerModal"
      className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
      onClick={(e) => e.target.id === "registerModal" && close()}
    >
      <div className="bg-black/80 backdrop-blur-lg rounded-2xl shadow-2xl border border-white/20 w-full max-w-md p-6">
        <div className="flex justify-between items-center mb-6">
          <h3 className="text-2xl font-bold text-white">Registro de Usuario</h3>
          <button
            id="closeRegisterBtn"
            onClick={close}
            className="text-gray-400 hover:text-white transition-colors p-1"
            aria-label="Cerrar"
          >
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>
        
        <form onSubmit={submit} id="registerForm" className="space-y-4">
          <div>
            <label htmlFor="regUsername" className="block mb-2 text-sm font-medium text-gray-300">
              Usuario
            </label>
            <input
              id="regUsername"
              type="text"
              required
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
              placeholder="Ingresa tu usuario"
            />
          </div>
          
          <div>
            <label htmlFor="regEmail" className="block mb-2 text-sm font-medium text-gray-300">
              Correo electrónico
            </label>
            <input
              id="regEmail"
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
              placeholder="correo@ejemplo.com"
            />
          </div>
          
          <div>
            <label htmlFor="regPassword" className="block mb-2 text-sm font-medium text-gray-300">
              Contraseña
            </label>
            <input
              id="regPassword"
              type="password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
              placeholder="Mínimo 6 caracteres"
            />
          </div>
          
          <button
            type="submit"
            className="w-full py-3 bg-green-600/90 hover:bg-green-600 text-white font-medium rounded-lg transition-all duration-200 shadow-lg backdrop-blur-sm border border-green-500/20 mt-6"
          >
            Registrarse
          </button>
          
          {error && (
            <div id="registerError" className="text-red-400 text-center text-sm bg-red-500/10 border border-red-500/20 rounded-lg p-3 mt-4">
              {error}
            </div>
          )}
        </form>
      </div>
    </div>
  );
}