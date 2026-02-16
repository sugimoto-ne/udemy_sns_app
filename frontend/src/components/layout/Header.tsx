import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  AppBar,
  Toolbar,
  Typography,
  IconButton,
  Avatar,
  Menu,
  MenuItem,
  Box,
  useMediaQuery,
  useTheme as useMuiTheme,
} from '@mui/material';
import {
  Home as HomeIcon,
  Person as PersonIcon,
  Logout as LogoutIcon,
  Bookmark as BookmarkIcon,
} from '@mui/icons-material';
import { useAuth } from '../../contexts/AuthContext';
import { ThemeSwitcher } from '../common/ThemeSwitcher';

export const Header: React.FC = () => {
  const navigate = useNavigate();
  const { user, logout } = useAuth();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const muiTheme = useMuiTheme();
  const isMobile = useMediaQuery(muiTheme.breakpoints.down('sm'));

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = async () => {
    handleMenuClose();
    await logout();
    navigate('/login');
  };

  const handleProfile = () => {
    handleMenuClose();
    if (user) {
      navigate(`/users/${user.username}`);
    }
  };

  const handleHome = () => {
    navigate('/');
  };

  const handleBookmarks = () => {
    navigate('/bookmarks');
  };

  return (
    <AppBar position="sticky">
      <Toolbar sx={{ gap: { xs: 0.5, sm: 1 } }}>
        {/* ロゴ/アプリ名 */}
        <Typography
          variant="h6"
          component="div"
          sx={{
            flexGrow: 0,
            mr: { xs: 1, sm: 3 },
            cursor: 'pointer',
            fontSize: { xs: '1rem', sm: '1.25rem' },
          }}
          onClick={handleHome}
        >
          {isMobile ? 'SNS' : 'SNS App'}
        </Typography>

        {/* ナビゲーションボタン */}
        <Box sx={{ flexGrow: 1, display: 'flex', gap: 1 }}>
          <IconButton color="inherit" onClick={handleHome} size={isMobile ? 'small' : 'medium'}>
            <HomeIcon />
          </IconButton>
          <IconButton color="inherit" onClick={handleBookmarks} size={isMobile ? 'small' : 'medium'}>
            <BookmarkIcon />
          </IconButton>
        </Box>

        {/* テーマスイッチャー */}
        <ThemeSwitcher />

        {/* ユーザーメニュー */}
        <IconButton onClick={handleMenuOpen} sx={{ p: 0 }} data-testid="user-menu-button">
          <Avatar
            src={user?.avatar_url || undefined}
            alt={user?.username}
            sx={{ width: { xs: 32, sm: 40 }, height: { xs: 32, sm: 40 } }}
          >
            {user?.username?.charAt(0).toUpperCase()}
          </Avatar>
        </IconButton>

        <Menu
          anchorEl={anchorEl}
          open={Boolean(anchorEl)}
          onClose={handleMenuClose}
          anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'right',
          }}
          transformOrigin={{
            vertical: 'top',
            horizontal: 'right',
          }}
        >
          <MenuItem onClick={handleProfile}>
            <PersonIcon sx={{ mr: 1 }} />
            プロフィール
          </MenuItem>
          <MenuItem onClick={handleLogout} data-testid="logout-button">
            <LogoutIcon sx={{ mr: 1 }} />
            ログアウト
          </MenuItem>
        </Menu>
      </Toolbar>
    </AppBar>
  );
};
