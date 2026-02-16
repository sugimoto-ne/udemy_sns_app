import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ThemeProvider } from './contexts/ThemeContext';
import { AuthProvider } from './contexts/AuthContext';
import { ProtectedRoute } from './components/common/ProtectedRoute';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';
import { HomePage } from './pages/HomePage';
import { PostDetailPage } from './pages/PostDetailPage';
import { UserProfilePage } from './pages/UserProfilePage';
import { NotificationsPage } from './pages/NotificationsPage';
import { FollowersPage } from './pages/FollowersPage';
import { FollowingPage } from './pages/FollowingPage';
import { SettingsPage } from './pages/SettingsPage';
import { BookmarksPage } from './pages/BookmarksPage';
import { HashtagPage } from './pages/HashtagPage';
import { PasswordResetPage } from './pages/PasswordResetPage';
import { PasswordResetConfirmPage } from './pages/PasswordResetConfirmPage';
import { EmailVerificationPage } from './pages/EmailVerificationPage';
import { EmailVerificationPendingPage } from './pages/EmailVerificationPendingPage';
import { ApprovalPendingPage } from './pages/ApprovalPendingPage';

// React Query クライアント作成
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider>
        <AuthProvider>
          <BrowserRouter>
            <Routes>
              {/* 公開ルート */}
              <Route path="/login" element={<LoginPage />} />
              <Route path="/register" element={<RegisterPage />} />
              <Route path="/auth/password-reset" element={<PasswordResetPage />} />
              <Route path="/auth/password-reset/confirm" element={<PasswordResetConfirmPage />} />
              <Route path="/auth/email/verify" element={<EmailVerificationPage />} />
              <Route path="/auth/email/verify-pending" element={<EmailVerificationPendingPage />} />
              <Route path="/auth/approval-pending" element={<ApprovalPendingPage />} />

              {/* 保護されたルート */}
              <Route
                path="/"
                element={
                  <ProtectedRoute>
                    <HomePage />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/posts/:postId"
                element={
                  <ProtectedRoute>
                    <PostDetailPage />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/users/:username"
                element={
                  <ProtectedRoute>
                    <UserProfilePage />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/notifications"
                element={
                  <ProtectedRoute>
                    <NotificationsPage />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/followers"
                element={
                  <ProtectedRoute>
                    <FollowersPage />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/following"
                element={
                  <ProtectedRoute>
                    <FollowingPage />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/settings"
                element={
                  <ProtectedRoute>
                    <SettingsPage />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/bookmarks"
                element={
                  <ProtectedRoute>
                    <BookmarksPage />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/hashtags/:hashtagName"
                element={
                  <ProtectedRoute>
                    <HashtagPage />
                  </ProtectedRoute>
                }
              />

              {/* 未定義のルートはホームにリダイレクト */}
              <Route path="*" element={<Navigate to="/" replace />} />
            </Routes>
          </BrowserRouter>
        </AuthProvider>
      </ThemeProvider>
    </QueryClientProvider>
  );
}

export default App;
