import type { components } from "../types/api";

type LoginResponse = components["schemas"]["LoginResponse"];
type AuthUser = components["schemas"]["AuthUser"];

interface StoredUser extends AuthUser {
  password: string;
  token?: string;
  tokenExpiry?: Date;
}

// Mock users storage (in production, this would be a database)
const mockUsers: Map<string, StoredUser> = new Map([
  [
    "alice@example.com",
    {
      id: "550e8400-e29b-41d4-a716-446655440001",
      email: "alice@example.com",
      name: "Alice Johnson",
      password: "password123", // In production, this would be hashed
    },
  ],
  [
    "bob@example.com",
    {
      id: "550e8400-e29b-41d4-a716-446655440002",
      email: "bob@example.com",
      name: "Bob Smith",
      password: "password456", // In production, this would be hashed
    },
  ],
]);

// Token storage (in production, this would be Redis or similar)
const tokenStorage: Map<string, StoredUser> = new Map();

export const authService = {
  /**
   * Authenticate user and generate token
   */
  async login(email: string, password: string): Promise<LoginResponse | null> {
    const user = mockUsers.get(email);
    
    // In production, use proper password hashing comparison
    if (!user || user.password !== password) {
      return null;
    }

    // Generate token
    const token = crypto.randomUUID();
    const expiresIn = 86400; // 24 hours in seconds
    const tokenExpiry = new Date(Date.now() + expiresIn * 1000);

    // Store token
    tokenStorage.set(token, {
      ...user,
      token,
      tokenExpiry,
    });

    return {
      accessToken: token,
      tokenType: "Bearer",
      expiresIn,
      user: {
        id: user.id,
        email: user.email,
        name: user.name,
      },
    };
  },

  /**
   * Invalidate token
   */
  async logout(token: string): Promise<boolean> {
    return tokenStorage.delete(token);
  },

  /**
   * Get user by token
   */
  async getUserByToken(token: string): Promise<AuthUser | null> {
    const storedUser = tokenStorage.get(token);
    
    if (!storedUser || !storedUser.tokenExpiry) {
      return null;
    }

    // Check if token is expired
    if (new Date() > storedUser.tokenExpiry) {
      tokenStorage.delete(token);
      return null;
    }

    return {
      id: storedUser.id,
      email: storedUser.email,
      name: storedUser.name,
    };
  },

  /**
   * Clean up expired tokens (could be run periodically)
   */
  cleanupExpiredTokens(): void {
    const now = new Date();
    for (const [token, user] of tokenStorage.entries()) {
      if (user.tokenExpiry && now > user.tokenExpiry) {
        tokenStorage.delete(token);
      }
    }
  },
};