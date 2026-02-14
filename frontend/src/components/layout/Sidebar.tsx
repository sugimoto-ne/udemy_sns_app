import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  Box,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Divider,
  Avatar,
  Typography,
} from '@mui/material';
import {
  Home as HomeIcon,
  Notifications as NotificationsIcon,
  Person as PersonIcon,
  People as PeopleIcon,
  PersonAdd as PersonAddIcon,
  Settings as SettingsIcon,
} from '@mui/icons-material';
import { useAuth } from '../../contexts/AuthContext';

export const Sidebar: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { user } = useAuth();

  const menuItems = [
    { text: 'ホーム', icon: <HomeIcon />, path: '/' },
    { text: '通知', icon: <NotificationsIcon />, path: '/notifications' },
    { text: 'プロフィール', icon: <PersonIcon />, path: user ? `/users/${user.username}` : '/profile' },
    { text: 'フォロワー', icon: <PeopleIcon />, path: '/followers' },
    { text: 'フォロー中', icon: <PersonAddIcon />, path: '/following' },
  ];

  const isActive = (path: string) => {
    if (path === '/') {
      return location.pathname === '/';
    }
    return location.pathname.startsWith(path);
  };

  return (
    <Box
      sx={{
        width: 280,
        height: 'calc(100vh - 64px)',
        position: 'fixed',
        top: 64,
        left: 0,
        borderRight: 1,
        borderColor: 'divider',
        display: 'flex',
        flexDirection: 'column',
        overflowY: 'auto',
        bgcolor: 'background.paper',
      }}
    >
      {/* ユーザー情報 */}
      {user && (
        <Box sx={{ p: 2, borderBottom: 1, borderColor: 'divider' }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
            <Avatar
              src={user.avatar_url || undefined}
              alt={user.username}
              sx={{ width: 48, height: 48 }}
            >
              {user.username.charAt(0).toUpperCase()}
            </Avatar>
            <Box sx={{ overflow: 'hidden' }}>
              <Typography variant="subtitle1" fontWeight="bold" noWrap>
                {user.display_name || user.username}
              </Typography>
              <Typography variant="body2" color="text.secondary" noWrap>
                @{user.username}
              </Typography>
            </Box>
          </Box>
        </Box>
      )}

      {/* メニュー */}
      <List sx={{ flexGrow: 1, py: 1 }}>
        {menuItems.map((item) => (
          <ListItem key={item.text} disablePadding>
            <ListItemButton
              selected={isActive(item.path)}
              onClick={() => navigate(item.path)}
              sx={{
                py: 1.5,
                '&.Mui-selected': {
                  bgcolor: 'action.selected',
                  borderRight: 3,
                  borderColor: 'primary.main',
                },
                '&:hover': {
                  bgcolor: 'action.hover',
                },
              }}
            >
              <ListItemIcon sx={{ color: isActive(item.path) ? 'primary.main' : 'inherit' }}>
                {item.icon}
              </ListItemIcon>
              <ListItemText
                primary={item.text}
                primaryTypographyProps={{
                  fontWeight: isActive(item.path) ? 'bold' : 'normal',
                }}
              />
            </ListItemButton>
          </ListItem>
        ))}
      </List>

      <Divider />

      {/* 設定 */}
      <List>
        <ListItem disablePadding>
          <ListItemButton
            selected={location.pathname === '/settings'}
            onClick={() => navigate('/settings')}
            sx={{
              py: 1.5,
              '&.Mui-selected': {
                bgcolor: 'action.selected',
                borderRight: 3,
                borderColor: 'primary.main',
              },
            }}
          >
            <ListItemIcon sx={{ color: location.pathname === '/settings' ? 'primary.main' : 'inherit' }}>
              <SettingsIcon />
            </ListItemIcon>
            <ListItemText
              primary="設定"
              primaryTypographyProps={{
                fontWeight: location.pathname === '/settings' ? 'bold' : 'normal',
              }}
            />
          </ListItemButton>
        </ListItem>
      </List>
    </Box>
  );
};
