import app from '../index';
import type { components } from '../types/api';

type LoginResponse = components['schemas']['LoginResponse'];

export async function loginTestUser(email: string, password: string): Promise<string> {
  const res = await app.request('/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password }),
  });

  if (res.status !== 200) {
    throw new Error('Failed to login test user');
  }

  const body = await res.json() as LoginResponse;
  return body.accessToken;
}

export function createAuthHeaders(token: string): Record<string, string> {
  return {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
}

// Default test user credentials
export const TEST_USERS = {
  alice: {
    email: 'alice@example.com',
    password: 'password123',
    id: '550e8400-e29b-41d4-a716-446655440001',
  },
  bob: {
    email: 'bob@example.com',
    password: 'password456',
    id: '550e8400-e29b-41d4-a716-446655440002',
  },
};