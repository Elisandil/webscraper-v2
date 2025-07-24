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
      className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center"
    >
      <div className="bg-gray-800 rounded-2xl shadow-2xl border border-white/10 w-full max-w-md p-6">
        <div className="flex justify-between mb-4">
          <h3 className="text-xl font-bold text-white">Registro de Usuario</h3>
          <button
            id="closeRegisterBtn"
            onClick={close}
            className="text-gray-400 hover:text-white"
          >
            &times;
          </button>
        </div>
        <form onSubmit={submit} id="registerForm" className="space-y-4">
          <div>
            <label htmlFor="regUsername" className="block mb-1 text-gray-300">
              Usuario
            </label>
            <input
              id="regUsername"
              type="text"
              required
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full px-3 py-2 bg-gray-700 rounded focus:outline-none focus:ring"
            />
          </div>
          <div>
            <label htmlFor="regEmail" className="block mb-1 text-gray-300">
              Correo electrónico
            </label>
            <input
              id="regEmail"
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-3 py-2 bg-gray-700 rounded focus:outline-none focus:ring"
            />
          </div>
          <div>
            <label htmlFor="regPassword" className="block mb-1 text-gray-300">
              Contraseña
            </label>
            <input
              id="regPassword"
              type="password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-3 py-2 bg-gray-700 rounded focus:outline-none focus:ring"
            />
          </div>
          <button
            type="submit"
            className="w-full py-2 bg-green-600 rounded hover:bg-green-500 transition text-white font-semibold"
          >
            Registrarse
          </button>
          {error && (
            <div id="registerError" className="text-red-400 text-center">
              {error}
            </div>
          )}
        </form>
      </div>
    </div>
  );
}
