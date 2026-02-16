import React, { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';
import type { User } from '../types/user';
import type { LoginRequest, RegisterRequest } from '../types/api';
import * as authApi from '../api/auth';

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (data: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  updateUser: (user: User) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUserState] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // 初期化：Cookieからユーザー情報を取得
  useEffect(() => {
    const initAuth = async () => {
      try {
        // Cookieにトークンがあれば、現在のユーザー情報を取得
        const currentUser = await authApi.getCurrentUser();
        setUserState(currentUser);
      } catch (error) {
        // Cookieにトークンがない、または無効な場合は何もしない
        setUserState(null);
      } finally {
        setIsLoading(false);
      }
    };

    initAuth();
  }, []);

  // ログイン
  const login = async (data: LoginRequest): Promise<void> => {
    try {
      const user = await authApi.login(data);
      setUserState(user);
    } catch (error) {
      throw error;
    }
  };

  // 新規登録（管理者承認制: ログイン状態にしない）
  const register = async (data: RegisterRequest): Promise<void> => {
    try {
      await authApi.register(data);
      // 管理者承認制: ユーザー状態は更新しない（承認されるまでログイン不可）
      setUserState(null);
    } catch (error) {
      throw error;
    }
  };

  // ログアウト
  const logout = async (): Promise<void> => {
    try {
      await authApi.logout();
    } catch (error) {
      console.error('Logout API error:', error);
    } finally {
      setUserState(null);
    }
  };

  // ユーザー情報更新
  const updateUser = (updatedUser: User): void => {
    setUserState(updatedUser);
  };

  const value: AuthContextType = {
    user,
    isLoading,
    isAuthenticated: !!user,
    login,
    register,
    logout,
    updateUser,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// カスタムフック
export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
