import React, { useState, useEffect, useRef } from "react";
import { apiRequest } from "../../api/client";
import { useAlert } from "../../contexts/AlertContext";

export default function RegisterModal() {
  const { showSuccess, showError } = useAlert();
  const [isOpen, setIsOpen] = useState(false);
  const [formData, setFormData] = useState({
    username: "",
    email: "",
    password: "",
    confirmPassword: "",
  });
  const [isLoading, setIsLoading] = useState(false);
  // Tracks whether the mousedown started on the backdrop so that dragging text
  // from inside the form and releasing outside does not close the modal.
  const mouseDownOnBackdrop = useRef(false);

  useEffect(() => {
    const btn = document.getElementById("openRegisterBtn");
    if (btn) btn.addEventListener("click", () => setIsOpen(true));
  }, []);

  const close = () => {
    setIsOpen(false);
    setFormData({ username: "", email: "", password: "", confirmPassword: "" });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!formData.username.trim() || !formData.email.trim() || !formData.password.trim()) {
      showError("Por favor completa todos los campos");
      return;
    }

    if (formData.password !== formData.confirmPassword) {
      showError("Las contraseñas no coinciden");
      return;
    }

    if (formData.password.length < 6) {
      showError("La contraseña debe tener al menos 6 caracteres");
      return;
    }

    setIsLoading(true);
    try {
      const { ok, data } = await apiRequest("/auth/register", {
        method: "POST",
        body: JSON.stringify({
          username: formData.username,
          email: formData.email,
          password: formData.password,
        }),
      });

      if (ok) {
        showSuccess("Usuario registrado exitosamente. Por favor inicia sesión.");
        close();
      } else {
        showError(data.error || "Error al registrar usuario");
      }
    } catch (error) {
      showError("Error de conexión");
    } finally {
      setIsLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div
      id="registerModal"
      className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
      onMouseDown={(e) => { mouseDownOnBackdrop.current = e.target.id === "registerModal"; }}
      onClick={(e) => { if (mouseDownOnBackdrop.current && e.target.id === "registerModal") close(); }}
    >
      <div className="animate-modal-in bg-black/80 backdrop-blur-lg rounded-2xl shadow-2xl border border-white/20 w-full max-w-md">
        <div className="flex justify-between items-center px-6 py-5 border-b border-white/20">
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

        <form onSubmit={handleSubmit} id="registerForm" className="space-y-4 px-6 py-6">
          <div>
            <label htmlFor="regUsername" className="block mb-2 text-sm font-medium text-gray-300">
              Usuario
            </label>
            <input
              id="regUsername"
              type="text"
              required
              value={formData.username}
              onChange={(e) => setFormData({ ...formData, username: e.target.value })}
              disabled={isLoading}
              className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:border-transparent transition-all duration-200"
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
              value={formData.email}
              onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              disabled={isLoading}
              className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:border-transparent transition-all duration-200"
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
              value={formData.password}
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
              disabled={isLoading}
              className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:border-transparent transition-all duration-200"
              placeholder="Mínimo 6 caracteres"
            />
          </div>

          <div>
            <label htmlFor="regConfirmPassword" className="block mb-2 text-sm font-medium text-gray-300">
              Confirmar Contraseña
            </label>
            <input
              id="regConfirmPassword"
              type="password"
              required
              value={formData.confirmPassword}
              onChange={(e) => setFormData({ ...formData, confirmPassword: e.target.value })}
              disabled={isLoading}
              className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:border-transparent transition-all duration-200"
              placeholder="Repite tu contraseña"
            />
          </div>

          <button
            type="submit"
            disabled={isLoading}
            className="w-full py-3 bg-gradient-to-r from-violet-600 to-indigo-600 hover:from-violet-500 hover:to-indigo-500 text-white font-medium rounded-lg transition-all duration-200 shadow-lg shadow-violet-500/30 backdrop-blur-sm border border-violet-500/20 mt-6 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isLoading ? "Registrando..." : "Registrarse"}
          </button>
        </form>
      </div>
    </div>
  );
}