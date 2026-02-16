import React, { useState } from 'react';
import { Container, TextField, Button, Typography, Box, Alert } from '@mui/material';
import { useRequestPasswordReset } from '../hooks/useAuth';

export const PasswordResetPage: React.FC = () => {
  const [email, setEmail] = useState('');
  const mutation = useRequestPasswordReset();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await mutation.mutateAsync({ email });
    } catch (error) {
      // エラーハンドリング
    }
  };

  return (
    <Container maxWidth="sm" sx={{ py: 8 }}>
      <Typography variant="h4" gutterBottom align="center">
        パスワードリセット
      </Typography>

      <Typography variant="body2" color="text.secondary" gutterBottom align="center" mb={4}>
        登録したメールアドレスを入力してください。
        <br />
        パスワードリセット用のリンクをお送りします。
      </Typography>

      {mutation.isSuccess && (
        <Alert severity="success" sx={{ mb: 2 }}>
          パスワードリセット用のメールを送信しました。メールボックスをご確認ください。
        </Alert>
      )}

      {mutation.isError && (
        <Alert severity="error" sx={{ mb: 2 }}>
          エラーが発生しました。もう一度お試しください。
        </Alert>
      )}

      <Box component="form" onSubmit={handleSubmit}>
        <TextField
          label="メールアドレス"
          type="email"
          fullWidth
          required
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          sx={{ mb: 2 }}
          disabled={mutation.isPending}
        />

        <Button
          type="submit"
          variant="contained"
          fullWidth
          size="large"
          disabled={mutation.isPending}
        >
          {mutation.isPending ? '送信中...' : 'リセットメールを送信'}
        </Button>
      </Box>
    </Container>
  );
};
