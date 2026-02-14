import React, { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';
import type { User } from '../types/user';
import type { AuthResponse, LoginRequest, RegisterRequest } from '../types/api';
import { getToken, setToken, getUser, setUser, clearAuth } from '../utils/storage';
import * as authApi from '../api/auth';

interface AuthContextType {
  user: User | null;
  token: string | null;
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
  const [token, setTokenState] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // 初期化：LocalStorageからユーザー情報を復元
  useEffect(() => {
    const initAuth = async () => {
      try {
        const storedToken = getToken();
        const storedUser = getUser();

        if (storedToken && storedUser) {
          setTokenState(storedToken);
          setUserState(storedUser);

          // トークンが有効か確認（現在のユーザー情報を取得）
          try {
            const currentUser = await authApi.getCurrentUser();
            setUserState(currentUser);
            setUser(currentUser);
          } catch (error) {
            // トークンが無効な場合はクリア
            clearAuth();
            setTokenState(null);
            setUserState(null);
          }
        }
      } catch (error) {
        console.error('Failed to initialize auth:', error);
      } finally {
        setIsLoading(false);
      }
    };

    initAuth();
  }, []);

  // ログイン
  const login = async (data: LoginRequest): Promise<void> => {
    try {
      const response: AuthResponse = await authApi.login(data);
      setToken(response.token);
      setUser(response.user);
      setTokenState(response.token);
      setUserState(response.user);
    } catch (error) {
      throw error;
    }
  };

  // 新規登録
  const register = async (data: RegisterRequest): Promise<void> => {
    try {
      const response: AuthResponse = await authApi.register(data);
      setToken(response.token);
      setUser(response.user);
      setTokenState(response.token);
      setUserState(response.user);
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
      clearAuth();
      setTokenState(null);
      setUserState(null);
    }
  };

  // ユーザー情報更新
  const updateUser = (updatedUser: User): void => {
    setUser(updatedUser);
    setUserState(updatedUser);
  };

  const value: AuthContextType = {
    user,
    token,
    isLoading,
    isAuthenticated: !!user && !!token,
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
