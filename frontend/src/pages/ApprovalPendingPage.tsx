import React from 'react';
import { Link as RouterLink } from 'react-router-dom';
import {
  Box,
  Paper,
  Typography,
  Button,
  Alert,
} from '@mui/material';
import { CheckCircleOutline } from '@mui/icons-material';

export const ApprovalPendingPage: React.FC = () => {
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
          maxWidth: 500,
          width: '100%',
          mx: 2,
          textAlign: 'center',
        }}
      >
        <CheckCircleOutline
          sx={{
            fontSize: 80,
            color: 'success.main',
            mb: 2,
          }}
        />

        <Typography variant="h4" component="h1" gutterBottom>
          登録が完了しました
        </Typography>

        <Alert severity="info" sx={{ mt: 3, mb: 3, textAlign: 'left' }}>
          <Typography variant="body1" gutterBottom>
            ご登録ありがとうございます。
          </Typography>
          <Typography variant="body2" paragraph>
            現在、管理者による承認待ちの状態です。承認が完了次第、ログイン可能になります。
          </Typography>
          <Typography variant="body2">
            承認には通常1〜2営業日かかります。しばらくお待ちください。
          </Typography>
        </Alert>

        <Box sx={{ mt: 4 }}>
          <Button
            component={RouterLink}
            to="/login"
            variant="contained"
            size="large"
            fullWidth
          >
            ログイン画面に戻る
          </Button>
        </Box>
      </Paper>
    </Box>
  );
};
