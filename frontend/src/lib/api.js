const BASE = '/api/v1';

let _token = localStorage.getItem('lm_token') || '';

export function setToken(t) {
  _token = t;
  localStorage.setItem('lm_token', t);
}

export function clearToken() {
  _token = '';
  localStorage.removeItem('lm_token');
}

export function getToken() {
  return _token;
}

async function req(method, path, body) {
  const res = await fetch(BASE + path, {
    method,
    headers: {
      'Content-Type': 'application/json',
      ...((_token) ? { Authorization: `Bearer ${_token}` } : {}),
    },
    body: body ? JSON.stringify(body) : undefined,
  });
  const json = await res.json();
  if (!json.success) throw new Error(json.error || 'Request failed');
  return json.data;
}

export const api = {
  login:  (email, password)      => req('POST', '/admin/login', { email, password }),

  apps:   {
    list:   ()                   => req('GET',  '/admin/apps'),
    create: (body)               => req('POST', '/admin/apps', body),
    update: (id, body)           => req('PUT',  `/admin/apps/${id}`, body),
    delete: (id)                 => req('DELETE', `/admin/apps/${id}`),
  },

  licenses: {
    list:   (appId)              => req('GET',  `/admin/licenses${appId ? `?appId=${appId}` : ''}`),
    create: (body)               => req('POST', '/admin/licenses', body),
    update: (id, body)           => req('PUT',  `/admin/licenses/${id}`, body),
    delete: (id)                 => req('DELETE', `/admin/licenses/${id}`),
    hwids:  (id)                 => req('GET',  `/admin/licenses/${id}/hwids`),
    addHwid: (id, body)          => req('POST', `/admin/licenses/${id}/hwids`, body),
  },

  hwids: {
    delete: (id)                 => req('DELETE', `/admin/hwids/${id}`),
  },
};
