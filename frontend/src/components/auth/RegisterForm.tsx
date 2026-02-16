import React, { useState } from 'react';
import { useNavigate, Link as RouterLink } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import {
  Box,
  Button,
  TextField,
  Typography,
  Alert,
  Paper,
  Link,
} from '@mui/material';
import { useAuth } from '../../contexts/AuthContext';
import type { RegisterRequest } from '../../types/api';

interface RegisterFormData extends RegisterRequest {
  passwordConfirm: string;
}

export const RegisterForm: React.FC = () => {
  const navigate = useNavigate();
  const { register: registerUser } = useAuth();
  const [error, setError] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<RegisterFormData>();

  const password = watch('password');

  const onSubmit = async (data: RegisterFormData) => {
    try {
      setIsLoading(true);
      setError('');
      const { passwordConfirm, ...registerData } = data;
      await registerUser(registerData);
      // 管理者承認制: 承認待ちメッセージページに遷移
      navigate('/auth/approval-pending');
    } catch (err: any) {
      setError(
        err.response?.data?.error?.message || '登録に失敗しました'
      );
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Box
      sx={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '100vh',
        width: '100%',
        bgcolor: 'background.default',
      }}
    >
      <Paper
        elevation={3}
        sx={{
          p: { xs: 3, sm: 4 },
          maxWidth: 400,
          width: '100%',
          mx: 2,
        }}
      >
        <Typography variant="h4" component="h1" gutterBottom textAlign="center">
          新規登録
        </Typography>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        <form onSubmit={handleSubmit(onSubmit)}>
          <TextField
            fullWidth
            label="ユーザー名"
            margin="normal"
            inputProps={{ 'data-testid': 'username-input' }}
            {...register('username', {
              required: 'ユーザー名を入力してください',
              minLength: {
                value: 3,
                message: 'ユーザー名は3文字以上で入力してください',
              },
              maxLength: {
                value: 20,
                message: 'ユーザー名は20文字以内で入力してください',
              },
              pattern: {
                value: /^[a-zA-Z0-9_]+$/,
                message: 'ユーザー名は英数字とアンダースコアのみ使用できます',
              },
            })}
            error={!!errors.username}
            helperText={errors.username?.message}
          />

          <TextField
            fullWidth
            label="メールアドレス"
            type="email"
            margin="normal"
            inputProps={{ 'data-testid': 'email-input' }}
            {...register('email', {
              required: 'メールアドレスを入力してください',
              pattern: {
                value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                message: '有効なメールアドレスを入力してください',
              },
            })}
            error={!!errors.email}
            helperText={errors.email?.message}
          />

          <TextField
            fullWidth
            label="パスワード"
            type="password"
            margin="normal"
            inputProps={{ 'data-testid': 'password-input' }}
            {...register('password', {
              required: 'パスワードを入力してください',
              minLength: {
                value: 8,
                message: 'パスワードは8文字以上で入力してください',
              },
            })}
            error={!!errors.password}
            helperText={errors.password?.message}
          />

          <TextField
            fullWidth
            label="パスワード（確認）"
            type="password"
            margin="normal"
            inputProps={{ 'data-testid': 'password-confirm-input' }}
            {...register('passwordConfirm', {
              required: 'パスワード（確認）を入力してください',
              validate: (value) =>
                value === password || 'パスワードが一致しません',
            })}
            error={!!errors.passwordConfirm}
            helperText={errors.passwordConfirm?.message}
          />

          <Button
            type="submit"
            fullWidth
            variant="contained"
            size="large"
            disabled={isLoading}
            data-testid="register-submit-button"
            sx={{ mt: 3, mb: 2 }}
          >
            {isLoading ? '登録中...' : '新規登録'}
          </Button>

          <Box textAlign="center">
            <Typography variant="body2">
              すでにアカウントをお持ちの方は{' '}
              <Link component={RouterLink} to="/login">
                ログイン
              </Link>
            </Typography>
          </Box>
        </form>
      </Paper>
    </Box>
  );
};
