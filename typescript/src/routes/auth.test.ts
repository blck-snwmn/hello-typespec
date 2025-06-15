import { describe, it, expect, beforeEach } from 'vitest';
import app from '../index';
import { authStore } from '../stores/auth';

describe('Auth Routes', () => {
  beforeEach(() => {
    // Clear all sessions before each test
    authStore.logout('');
    // @ts-ignore - accessing private property for testing
    authStore.sessions.clear();
  });

  describe('POST /auth/login', () => {
    it('should login with valid credentials', async () => {
      const res = await app.request('/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: 'alice@example.com',
          password: 'password123',
        }),
      });

      expect(res.status).toBe(200);
      const body = await res.json();
      expect(body).toHaveProperty('accessToken');
      expect(body).toHaveProperty('tokenType', 'Bearer');
      expect(body).toHaveProperty('expiresIn', 86400);
      expect(body).toHaveProperty('user');
      expect(body.user).toMatchObject({
        id: '550e8400-e29b-41d4-a716-446655440001',
        email: 'alice@example.com',
        name: 'Alice Johnson',
      });
    });

    it('should fail with invalid credentials', async () => {
      const res = await app.request('/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: 'alice@example.com',
          password: 'wrongpassword',
        }),
      });

      expect(res.status).toBe(401);
      const body = await res.json();
      expect(body.error).toMatchObject({
        code: 'UNAUTHORIZED',
        message: 'Invalid email or password',
      });
    });

    it('should fail with missing credentials', async () => {
      const res = await app.request('/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: 'alice@example.com',
        }),
      });

      expect(res.status).toBe(400);
      const body = await res.json();
      expect(body.error).toMatchObject({
        code: 'BAD_REQUEST',
        message: 'Email and password are required',
      });
    });
  });

  describe('POST /auth/logout', () => {
    it('should logout with valid token', async () => {
      // First login to get a token
      const loginRes = await app.request('/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: 'alice@example.com',
          password: 'password123',
        }),
      });
      const { accessToken } = await loginRes.json();

      // Then logout
      const res = await app.request('/auth/logout', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${accessToken}`,
        },
      });

      expect(res.status).toBe(200);
      const body = await res.json();
      expect(body).toMatchObject({
        message: 'Logged out successfully',
      });

      // Verify token is invalidated
      const meRes = await app.request('/auth/me', {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${accessToken}`,
        },
      });
      expect(meRes.status).toBe(401);
    });

    it('should fail without authentication', async () => {
      const res = await app.request('/auth/logout', {
        method: 'POST',
      });

      expect(res.status).toBe(401);
      const body = await res.json();
      expect(body.error).toMatchObject({
        code: 'UNAUTHORIZED',
        message: 'Missing Authorization header',
      });
    });
  });

  describe('GET /auth/me', () => {
    it('should return current user with valid token', async () => {
      // First login to get a token
      const loginRes = await app.request('/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: 'alice@example.com',
          password: 'password123',
        }),
      });
      const { accessToken } = await loginRes.json();

      // Get current user
      const res = await app.request('/auth/me', {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${accessToken}`,
        },
      });

      expect(res.status).toBe(200);
      const body = await res.json();
      expect(body).toHaveProperty('user');
      expect(body.user).toMatchObject({
        id: '550e8400-e29b-41d4-a716-446655440001',
        email: 'alice@example.com',
        name: 'Alice Johnson',
      });
    });

    it('should fail without authentication', async () => {
      const res = await app.request('/auth/me', {
        method: 'GET',
      });

      expect(res.status).toBe(401);
      const body = await res.json();
      expect(body.error).toMatchObject({
        code: 'UNAUTHORIZED',
        message: 'Missing Authorization header',
      });
    });

    it('should fail with invalid token', async () => {
      const res = await app.request('/auth/me', {
        method: 'GET',
        headers: {
          'Authorization': 'Bearer invalid-token',
        },
      });

      expect(res.status).toBe(401);
      const body = await res.json();
      expect(body.error).toMatchObject({
        code: 'UNAUTHORIZED',
        message: 'Invalid or expired token',
      });
    });
  });
});