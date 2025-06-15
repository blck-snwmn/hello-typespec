import type { Context } from "hono";
import type { components } from "../types/api";

type AuthUser = components["schemas"]["AuthUser"];

// Test token that works in test environment
export const TEST_AUTH_TOKEN = "test-auth-token-12345";

// Test users
export const TEST_USERS = {
  alice: {
    id: "550e8400-e29b-41d4-a716-446655440001",
    email: "alice@example.com",
    name: "Alice Johnson",
  },
  bob: {
    id: "550e8400-e29b-41d4-a716-446655440002",
    email: "bob@example.com",
    name: "Bob Smith",
  },
} as const;

// Mock auth service for tests
export class MockAuthService {
  private currentUser: AuthUser | null = null;

  setCurrentUser(user: AuthUser | null) {
    this.currentUser = user;
  }

  async getUserByToken(token: string): Promise<AuthUser | null> {
    if (token === TEST_AUTH_TOKEN && this.currentUser) {
      return this.currentUser;
    }
    return null;
  }

  reset() {
    this.currentUser = null;
  }
}

// Global mock auth service instance
export const mockAuthService = new MockAuthService();

// Create authenticated request headers
export function createAuthHeaders(token: string = TEST_AUTH_TOKEN): Record<string, string> {
  return {
    "Authorization": `Bearer ${token}`,
  };
}

// Setup authenticated context for tests
export function setupAuthenticatedTest(user: AuthUser = TEST_USERS.alice) {
  mockAuthService.setCurrentUser(user);
  return {
    token: TEST_AUTH_TOKEN,
    user,
    headers: createAuthHeaders(),
  };
}

// Mock auth middleware for tests
export async function mockAuthMiddleware(c: Context, next: () => Promise<void>) {
  const authHeader = c.req.header("Authorization");
  
  if (!authHeader || !authHeader.startsWith("Bearer ")) {
    return c.json({ error: { code: "UNAUTHORIZED", message: "Missing or invalid Authorization header" } }, 401);
  }

  const token = authHeader.substring(7);
  const user = await mockAuthService.getUserByToken(token);

  if (!user) {
    return c.json({ error: { code: "UNAUTHORIZED", message: "Invalid or expired token" } }, 401);
  }

  c.set("user", user);
  await next();
}

// Helper to make authenticated fetch requests in tests
export function createAuthenticatedFetch(baseUrl: string, token: string = TEST_AUTH_TOKEN) {
  return (path: string, options: RequestInit = {}) => {
    const headers = {
      ...createAuthHeaders(token),
      "Content-Type": "application/json",
      ...(options.headers || {}),
    };

    return fetch(`${baseUrl}${path}`, {
      ...options,
      headers,
    });
  };
}