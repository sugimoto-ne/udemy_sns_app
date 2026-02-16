import React from 'react';
import {
  Container,
  Paper,
  Typography,
  Box,
  Alert,
} from '@mui/material';
import HourglassEmptyIcon from '@mui/icons-material/HourglassEmpty';
import { useAuth } from '../contexts/AuthContext';

export const EmailVerificationPendingPage: React.FC = () => {
  const { user } = useAuth();

  return (
    <Container maxWidth="sm" sx={{ py: 8 }}>
      <Paper sx={{ p: 4, textAlign: 'center' }}>
        <Box sx={{ mb: 3 }}>
          <HourglassEmptyIcon sx={{ fontSize: 80, color: 'warning.main' }} />
        </Box>

        <Typography variant="h4" gutterBottom>
          アカウント承認待ち
        </Typography>

        <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
          {user?.email || 'ご登録いただいたアカウント'}は現在、管理者による承認待ちです。
        </Typography>

        <Alert severity="info" sx={{ mb: 3, textAlign: 'left' }}>
          <Typography variant="body2" sx={{ mb: 1 }}>
            <strong>ご登録ありがとうございます</strong>
          </Typography>
          <Typography variant="body2" component="div">
            管理者によるアカウント承認が完了次第、ログインできるようになります。
            <br /><br />
            通常、1〜2営業日以内に承認が完了します。
            <br />
            承認完了後、再度ログインをお試しください。
          </Typography>
        </Alert>

        <Typography variant="body2" color="text.secondary">
          承認に関するご質問がある場合は、管理者までお問い合わせください。
        </Typography>
      </Paper>
    </Container>
  );
};
