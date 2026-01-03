import { getAuthToken, saveAuth, type AuthUser } from './auth';

const API_BASE = import.meta.env.VITE_API_URL ?? '/api';

type AuthResponse = {
  token: string;
  user: AuthUser;
};

export type Subscription = {
  id: string;
  service_name: string;
  bank_name: string;
  card_last4: string;
  billing_cycle: 'monthly' | 'yearly';
  charge_date: string;
};

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = getAuthToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
  });

  const contentType = res.headers.get('content-type');
  const isJSON = contentType?.includes('application/json');
  const payload = isJSON ? await res.json() : null;

  if (!res.ok) {
    const message = payload?.error || 'Request failed';
    throw new Error(message);
  }

  return payload as T;
}

export async function registerUser(name: string, email: string, password: string) {
  const data = await request<AuthResponse>('/auth/register', {
    method: 'POST',
    body: JSON.stringify({ name, email, password }),
  });
  saveAuth(data.token, data.user);
  return data.user;
}

export async function loginUser(email: string, password: string) {
  const data = await request<AuthResponse>('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  });
  saveAuth(data.token, data.user);
  return data.user;
}

export async function getSubscriptions() {
  return request<Subscription[]>('/subscriptions');
}

export async function createSubscription(payload: Omit<Subscription, 'id'>) {
  return request<Subscription>('/subscriptions', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function updateSubscription(id: string, payload: Omit<Subscription, 'id'>) {
  return request<Subscription>(`/subscriptions/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  });
}

export async function deleteSubscription(id: string) {
  return request<{ id: string }>(`/subscriptions/${id}`, {
    method: 'DELETE',
  });
}
