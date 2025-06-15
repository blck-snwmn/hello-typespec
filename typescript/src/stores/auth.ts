import { v4 as uuidv4 } from 'uuid';

export interface AuthUser {
  id: string;
  email: string;
  name: string;
}

export interface AuthSession {
  token: string;
  user: AuthUser;
  expiresAt: Date;
}

export class AuthStore {
  private sessions: Map<string, AuthSession> = new Map();
  private users: Map<string, { password: string; user: AuthUser }> = new Map([
    [
      'alice@example.com',
      {
        password: 'password123', // In production, this would be hashed
        user: {
          id: '550e8400-e29b-41d4-a716-446655440001',
          email: 'alice@example.com',
          name: 'Alice Johnson',
        },
      },
    ],
    [
      'bob@example.com',
      {
        password: 'password456', // In production, this would be hashed
        user: {
          id: '550e8400-e29b-41d4-a716-446655440002',
          email: 'bob@example.com',
          name: 'Bob Smith',
        },
      },
    ],
  ]);

  login(email: string, password: string): AuthSession | null {
    const userRecord = this.users.get(email);
    if (!userRecord || userRecord.password !== password) {
      return null;
    }

    const token = uuidv4();
    const session: AuthSession = {
      token,
      user: userRecord.user,
      expiresAt: new Date(Date.now() + 24 * 60 * 60 * 1000), // 24 hours
    };

    this.sessions.set(token, session);
    return session;
  }

  logout(token: string): void {
    this.sessions.delete(token);
  }

  validateToken(token: string): AuthUser | null {
    const session = this.sessions.get(token);
    if (!session) {
      return null;
    }

    // Check if token is expired
    if (new Date() > session.expiresAt) {
      this.sessions.delete(token);
      return null;
    }

    return session.user;
  }

  // Clean up expired tokens
  cleanupExpiredTokens(): void {
    const now = new Date();
    for (const [token, session] of this.sessions.entries()) {
      if (now > session.expiresAt) {
        this.sessions.delete(token);
      }
    }
  }
}

// Singleton instance
export const authStore = new AuthStore();