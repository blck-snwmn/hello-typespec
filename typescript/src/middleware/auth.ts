import { MiddlewareHandler } from 'hono';
import { authStore } from '../stores/auth';
import { sendError } from '../types/errors';

declare module 'hono' {
  interface ContextVariableMap {
    user: {
      id: string;
      email: string;
      name: string;
    };
  }
}

export const authMiddleware: MiddlewareHandler = async (c, next) => {
  const authHeader = c.req.header('Authorization');
  
  if (!authHeader) {
    return sendError(c, 401, 'UNAUTHORIZED', 'Missing Authorization header');
  }

  const parts = authHeader.split(' ');
  if (parts.length !== 2 || parts[0] !== 'Bearer') {
    return sendError(c, 401, 'UNAUTHORIZED', 'Invalid Authorization header format');
  }

  const token = parts[1];
  const user = authStore.validateToken(token);

  if (!user) {
    return sendError(c, 401, 'UNAUTHORIZED', 'Invalid or expired token');
  }

  // Add user to context
  c.set('user', user);
  
  await next();
};