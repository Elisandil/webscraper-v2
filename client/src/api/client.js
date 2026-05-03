const API_BASE = '/api';

export async function apiRequest(path, options = {}) {
  const { noRedirectOn401, ...fetchOptions } = options;

  const res = await fetch(`${API_BASE}${path}`, {
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      ...(fetchOptions.headers || {}),
    },
    ...fetchOptions,
  });

  if (res.status === 401 && !noRedirectOn401) {
    window.location.href = '/';
    return { ok: false, status: 401, data: { error: 'Sesión expirada' } };
  }

  const data = await res.json().catch(() => ({}));
  return { ok: res.ok, status: res.status, data };
}

export async function publicApiRequest(path, options = {}) {
  const res = await fetch(`${API_BASE}/public${path}`, {
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers || {}),
    },
    ...options,
  });

  const data = await res.json().catch(() => ({}));
  return { ok: res.ok, status: res.status, data };
}
