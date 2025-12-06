const API_BASE = '/api';

export async function apiRequest(path, options = {}) {
  const token = localStorage.getItem('jwtToken');
  
  // Si se requiere autenticación y no hay token, rechazar la petición
  if (path !== '/auth/login' && path !== '/auth/register' && !token) {
    return { ok: false, status: 401, data: { error: 'No autenticado' } };
  }

  const res = await fetch(`${API_BASE}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { 'Authorization': `Bearer ${token}` } : {}),
      ...(options.headers || {})
    },
    ...options,
  });

  // Si recibimos un 401, limpiar la sesión
  if (res.status === 401) {
    localStorage.removeItem('jwtToken');
    window.location.reload();
    return { ok: false, status: 401, data: { error: 'Sesión expirada' } };
  }

  const data = await res.json().catch(() => ({}));
  return { ok: res.ok, status: res.status, data };
}