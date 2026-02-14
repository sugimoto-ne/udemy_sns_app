import React from 'react';
import { Box, Typography, Paper, Avatar, List, ListItem, ListItemAvatar, ListItemText } from '@mui/material';
import { MainLayout } from '../components/layout/MainLayout';
import {
  FavoriteBorder as LikeIcon,
  ChatBubbleOutline as CommentIcon,
  PersonAdd as FollowIcon,
} from '@mui/icons-material';

// ダミーデータ（将来的にAPIから取得）
const dummyNotifications = [
  {
    id: 1,
    type: 'like',
    user: { username: 'user1', display_name: 'ユーザー1', avatar_url: null },
    content: 'があなたの投稿にいいねしました',
    created_at: '2024-01-15T10:30:00Z',
  },
  {
    id: 2,
    type: 'comment',
    user: { username: 'user2', display_name: 'ユーザー2', avatar_url: null },
    content: 'があなたの投稿にコメントしました',
    created_at: '2024-01-15T09:15:00Z',
  },
  {
    id: 3,
    type: 'follow',
    user: { username: 'user3', display_name: 'ユーザー3', avatar_url: null },
    content: 'があなたをフォローしました',
    created_at: '2024-01-14T18:00:00Z',
  },
];

const getNotificationIcon = (type: string) => {
  switch (type) {
    case 'like':
      return <LikeIcon color="error" />;
    case 'comment':
      return <CommentIcon color="primary" />;
    case 'follow':
      return <FollowIcon color="success" />;
    default:
      return null;
  }
};

export const NotificationsPage: React.FC = () => {
  return (
    <MainLayout>
      <Box>
        <Typography variant="h4" gutterBottom fontWeight="bold" sx={{ mb: 3 }}>
          通知
        </Typography>

        <Paper>
          <List sx={{ p: 0 }}>
            {dummyNotifications.map((notification, index) => (
              <React.Fragment key={notification.id}>
                <ListItem
                  alignItems="flex-start"
                  sx={{
                    py: 2,
                    '&:hover': {
                      bgcolor: 'action.hover',
                      cursor: 'pointer',
                    },
                  }}
                >
                  <ListItemAvatar>
                    <Box sx={{ position: 'relative' }}>
                      <Avatar
                        src={notification.user.avatar_url || undefined}
                        alt={notification.user.username}
                      >
                        {notification.user.username.charAt(0).toUpperCase()}
                      </Avatar>
                      <Box
                        sx={{
                          position: 'absolute',
                          bottom: -4,
                          right: -4,
                          bgcolor: 'background.paper',
                          borderRadius: '50%',
                          p: 0.5,
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                        }}
                      >
                        {getNotificationIcon(notification.type)}
                      </Box>
                    </Box>
                  </ListItemAvatar>
                  <ListItemText
                    primary={
                      <Typography variant="body1">
                        <strong>{notification.user.display_name || notification.user.username}</strong>
                        {notification.content}
                      </Typography>
                    }
                    secondary={
                      <Typography variant="caption" color="text.secondary">
                        {new Date(notification.created_at).toLocaleString('ja-JP')}
                      </Typography>
                    }
                  />
                </ListItem>
                {index < dummyNotifications.length - 1 && <Box sx={{ borderBottom: 1, borderColor: 'divider' }} />}
              </React.Fragment>
            ))}
          </List>
        </Paper>

        {dummyNotifications.length === 0 && (
          <Box sx={{ textAlign: 'center', py: 8 }}>
            <Typography variant="h6" color="text.secondary">
              通知はありません
            </Typography>
          </Box>
        )}
      </Box>
    </MainLayout>
  );
};
