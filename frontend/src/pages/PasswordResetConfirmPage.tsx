import React, { useState } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { Container, TextField, Button, Typography, Box, Alert } from '@mui/material';
import { useConfirmPasswordReset } from '../hooks/useAuth';

export const PasswordResetConfirmPage: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const token = searchParams.get('token') || '';

  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');

  const mutation = useConfirmPasswordReset();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (newPassword !== confirmPassword) {
      setError('パスワードが一致しません');
      return;
    }

    if (newPassword.length < 8) {
      setError('パスワードは8文字以上である必要があります');
      return;
    }

    try {
      await mutation.mutateAsync({ token, new_password: newPassword });
      // 成功したらログインページへ
      setTimeout(() => {
        navigate('/login');
      }, 2000);
    } catch (error) {
      setError('パスワードのリセットに失敗しました');
    }
  };

  if (!token) {
    return (
      <Container maxWidth="sm" sx={{ py: 8 }}>
        <Alert severity="error">無効なリンクです</Alert>
      </Container>
    );
  }

  return (
    <Container maxWidth="sm" sx={{ py: 8 }}>
      <Typography variant="h4" gutterBottom align="center">
        新しいパスワードを設定
      </Typography>

      {mutation.isSuccess && (
        <Alert severity="success" sx={{ mb: 2 }}>
          パスワードをリセットしました。ログインページに移動します...
        </Alert>
      )}

      {(error || mutation.isError) && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error || 'エラーが発生しました'}
        </Alert>
      )}

      <Box component="form" onSubmit={handleSubmit}>
        <TextField
          label="新しいパスワード"
          type="password"
          fullWidth
          required
          value={newPassword}
          onChange={(e) => setNewPassword(e.target.value)}
          sx={{ mb: 2 }}
          disabled={mutation.isPending || mutation.isSuccess}
        />

        <TextField
          label="パスワード（確認）"
          type="password"
          fullWidth
          required
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          sx={{ mb: 2 }}
          disabled={mutation.isPending || mutation.isSuccess}
        />

        <Button
          type="submit"
          variant="contained"
          fullWidth
          size="large"
          disabled={mutation.isPending || mutation.isSuccess}
        >
          {mutation.isPending ? '変更中...' : 'パスワードを変更'}
        </Button>
      </Box>
    </Container>
  );
};
