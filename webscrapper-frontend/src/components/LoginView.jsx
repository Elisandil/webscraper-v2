import React, { useState } from 'react';
import { apiRequest } from '../api/client';

export default function LoginView({ onLogin, onAlert }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError]       = useState('');

  const submit = async (e) => {
    e.preventDefault();
    setError('');
    const { ok, data } = await apiRequest('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    });
    if (ok) {
      // Aquí sacamos el token del campo data.token
      const token = data?.data?.token;
      if (token) {
        localStorage.setItem('jwtToken', token);
        onLogin();
      } else {
        // Por si acaso no viene el token
        setError('No se recibió token de autenticación');
      }
    }
  };

  return (
    <div id="loginView" className="min-h-screen flex items-center justify-center">
      <form onSubmit={submit} id="loginForm"
        className="bg-gray-700 p-8 rounded-lg shadow-lg border border-gray-600 w-full max-w-sm space-y-6">
        <h2 className="text-2xl font-bold text-center text-gray-100">Iniciar Sesión</h2>

        <div>
          <label htmlFor="username" className="block mb-1 text-gray-300">Usuario</label>
          <input id="username" type="text"
            value={username} onChange={e => setUsername(e.target.value)}
            required
            className="w-full px-4 py-2 bg-gray-700/50 border border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-gray-500 focus:border-transparent text-gray-100" />
        </div>
        <div>
          <label htmlFor="password" className="block mb-1 text-gray-300">Contraseña</label>
          <input id="password" type="password" required
            value={password} onChange={e => setPassword(e.target.value)}
            className="w-full px-4 py-2 bg-gray-700/50 border border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-gray-500 focus:border-transparent text-gray-100" />
        </div>

        <button type="submit" id="loginBtn"
          className="w-full py-2 bg-green-600 rounded hover:bg-green-500 transition text-white font-semibold">
          Iniciar Sesión
        </button>

        <div className="text-center text-gray-400 mt-4">
          <span>¿No tienes cuenta? </span>
          <button id="openRegisterBtn" type="button"
            className="text-blue-400 hover:underline">
            Regístrate
          </button>
        </div>

        {error && <div id="loginError" className="text-red-400 text-center">{error}</div>}
      </form>
    </div>
  );
}
