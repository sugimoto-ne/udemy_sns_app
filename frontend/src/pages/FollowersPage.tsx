import React from 'react';
import { Box, Typography, Paper, Avatar, List, ListItem, ListItemAvatar, ListItemText, Button } from '@mui/material';
import { MainLayout } from '../components/layout/MainLayout';
import { useAuth } from '../contexts/AuthContext';
import { useFollowers } from '../hooks/useUsers';
import { CircularProgress, Alert } from '@mui/material';

export const FollowersPage: React.FC = () => {
  const { user } = useAuth();
  const { data: followersData, isLoading, error } = useFollowers(user?.username || '');

  return (
    <MainLayout>
      <Box>
        <Typography variant="h4" gutterBottom fontWeight="bold" sx={{ mb: 3 }}>
          フォロワー
        </Typography>

        {isLoading && (
          <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
            <CircularProgress />
          </Box>
        )}

        {error && (
          <Alert severity="error">
            フォロワーの取得に失敗しました
          </Alert>
        )}

        {followersData?.data && followersData.data.length === 0 && (
          <Box sx={{ textAlign: 'center', py: 8 }}>
            <Typography variant="h6" color="text.secondary">
              フォロワーはいません
            </Typography>
          </Box>
        )}

        {followersData?.data && followersData.data.length > 0 && (
          <Paper>
            <List sx={{ p: 0 }}>
              {followersData.data.map((follower, index) => (
                <React.Fragment key={follower.id}>
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
                        src={follower.avatar_url || undefined}
                        alt={follower.username}
                        sx={{ width: 48, height: 48 }}
                      >
                        {follower.username.charAt(0).toUpperCase()}
                      </Avatar>
                    </ListItemAvatar>
                    <ListItemText
                      primary={
                        <Typography variant="subtitle1" fontWeight="bold">
                          {follower.display_name || follower.username}
                        </Typography>
                      }
                      secondary={
                        <>
                          <Typography variant="body2" color="text.secondary">
                            @{follower.username}
                          </Typography>
                          {follower.bio && (
                            <Typography variant="body2" sx={{ mt: 0.5 }}>
                              {follower.bio}
                            </Typography>
                          )}
                        </>
                      }
                    />
                  </ListItem>
                  {index < followersData.data.length - 1 && <Box sx={{ borderBottom: 1, borderColor: 'divider' }} />}
                </React.Fragment>
              ))}
            </List>
          </Paper>
        )}
      </Box>
    </MainLayout>
  );
};
