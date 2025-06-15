import { Hono } from 'hono';
import { authStore } from '../stores/auth';
import { sendError, ErrorCode } from '../types/errors';
import { authMiddleware } from '../middleware/auth';
import type { components } from '../types/api';

type LoginRequest = components['schemas']['LoginRequest'];
type LoginResponse = components['schemas']['LoginResponse'];
type CurrentUserResponse = components['schemas']['CurrentUserResponse'];

const authRoutes = new Hono();

// POST /auth/login
authRoutes.post('/login', async (c) => {
  const body = await c.req.json<LoginRequest>();
  
  if (!body.email || !body.password) {
    return sendError(c, 400, ErrorCode.BAD_REQUEST, 'Email and password are required');
  }

  const session = authStore.login(body.email, body.password);
  
  if (!session) {
    return sendError(c, 401, ErrorCode.UNAUTHORIZED, 'Invalid email or password');
  }

  const response: LoginResponse = {
    accessToken: session.token,
    expiresIn: 86400, // 24 hours in seconds
    tokenType: 'Bearer',
    user: session.user,
  };

  return c.json(response);
});

// POST /auth/logout
authRoutes.post('/logout', authMiddleware, async (c) => {
  const authHeader = c.req.header('Authorization');
  const token = authHeader?.split(' ')[1];
  
  if (token) {
    authStore.logout(token);
  }

  return c.json({ message: 'Logged out successfully' });
});

// GET /auth/me
authRoutes.get('/me', authMiddleware, async (c) => {
  const user = c.get('user');
  
  const response: CurrentUserResponse = {
    user,
  };

  return c.json(response);
});

export default authRoutes;