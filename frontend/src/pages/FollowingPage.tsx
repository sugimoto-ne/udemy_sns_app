import React from 'react';
import { Box, Typography, Paper, Avatar, List, ListItem, ListItemAvatar, ListItemText, Button } from '@mui/material';
import { MainLayout } from '../components/layout/MainLayout';
import { useAuth } from '../contexts/AuthContext';
import { useFollowing } from '../hooks/useUsers';
import { CircularProgress, Alert } from '@mui/material';

export const FollowingPage: React.FC = () => {
  const { user } = useAuth();
  const { data: followingData, isLoading, error } = useFollowing(user?.username || '');

  return (
    <MainLayout>
      <Box>
        <Typography variant="h4" gutterBottom fontWeight="bold" sx={{ mb: 3 }}>
          フォロー中
        </Typography>

        {isLoading && (
          <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
            <CircularProgress />
          </Box>
        )}

        {error && (
          <Alert severity="error">
            フォロー中のユーザーの取得に失敗しました
          </Alert>
        )}

        {followingData?.data && followingData.data.length === 0 && (
          <Box sx={{ textAlign: 'center', py: 8 }}>
            <Typography variant="h6" color="text.secondary">
              フォロー中のユーザーはいません
            </Typography>
          </Box>
        )}

        {followingData?.data && followingData.data.length > 0 && (
          <Paper>
            <List sx={{ p: 0 }}>
              {followingData.data.map((following, index) => (
                <React.Fragment key={following.id}>
                  <ListItem
                    alignItems="flex-start"
                    sx={{
                      py: 2,
                      '&:hover': {
                        bgcolor: 'action.hover',
                      },
                    }}
                    secondaryAction={
                      <Button variant="outlined" size="small">
                        プロフィール
                      </Button>
                    }
                  >
                    <ListItemAvatar>
                      <Avatar
                        src={following.avatar_url || undefined}
                        alt={following.username}
                        sx={{ width: 48, height: 48 }}
                      >
                        {following.username.charAt(0).toUpperCase()}
                      </Avatar>
                    </ListItemAvatar>
                    <ListItemText
                      primary={
                        <Typography variant="subtitle1" fontWeight="bold">
                          {following.display_name || following.username}
                        </Typography>
                      }
                      secondary={
                        <>
                          <Typography variant="body2" color="text.secondary">
                            @{following.username}
                          </Typography>
                          {following.bio && (
                            <Typography variant="body2" sx={{ mt: 0.5 }}>
                              {following.bio}
                            </Typography>
                          )}
                        </>
                      }
                    />
                  </ListItem>
                  {index < followingData.data.length - 1 && <Box sx={{ borderBottom: 1, borderColor: 'divider' }} />}
                </React.Fragment>
              ))}
            </List>
          </Paper>
        )}
      </Box>
    </MainLayout>
  );
};
