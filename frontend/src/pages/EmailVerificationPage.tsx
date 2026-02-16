import React, { useEffect, useState } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { Container, Typography, Box, CircularProgress, Alert, Button } from '@mui/material';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import ErrorIcon from '@mui/icons-material/Error';
import { useVerifyEmail } from '../hooks/useAuth';

export const EmailVerificationPage: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const token = searchParams.get('token') || '';

  const mutation = useVerifyEmail();
  const [hasVerified, setHasVerified] = useState(false);

  useEffect(() => {
    if (token && !hasVerified) {
      setHasVerified(true);
      mutation.mutate({ token });
    }
  }, [token, hasVerified, mutation]);

  if (!token) {
    return (
      <Container maxWidth="sm" sx={{ py: 8 }}>
        <Alert severity="error">無効なリンクです</Alert>
      </Container>
    );
  }

  if (mutation.isPending) {
    return (
      <Container maxWidth="sm" sx={{ py: 8 }}>
        <Box display="flex" flexDirection="column" alignItems="center" gap={2}>
          <CircularProgress size={60} />
          <Typography variant="h6">メールアドレスを認証中...</Typography>
        </Box>
      </Container>
    );
  }

  if (mutation.isError) {
    return (
      <Container maxWidth="sm" sx={{ py: 8 }}>
        <Box display="flex" flexDirection="column" alignItems="center" gap={2}>
          <ErrorIcon color="error" sx={{ fontSize: 80 }} />
          <Typography variant="h5">認証に失敗しました</Typography>
          <Typography color="text.secondary" align="center">
            トークンが無効か、有効期限が切れています。
          </Typography>
          <Button variant="contained" onClick={() => navigate('/login')}>
            ログインページへ
          </Button>
        </Box>
      </Container>
    );
  }

  if (mutation.isSuccess) {
    return (
      <Container maxWidth="sm" sx={{ py: 8 }}>
        <Box display="flex" flexDirection="column" alignItems="center" gap={2}>
          <CheckCircleIcon color="success" sx={{ fontSize: 80 }} />
          <Typography variant="h5">認証が完了しました！</Typography>
          <Typography color="text.secondary" align="center">
            メールアドレスの認証が完了しました。
            <br />
            すべての機能をご利用いただけます。
          </Typography>
          <Button variant="contained" onClick={() => navigate('/')}>
            ホームへ
          </Button>
        </Box>
      </Container>
    );
  }

  return null;
};
