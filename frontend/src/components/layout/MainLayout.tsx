import React from 'react';
import { Box, useTheme, useMediaQuery } from '@mui/material';
import { Header } from './Header';
import { Sidebar } from './Sidebar';
import { CreatePostFab } from '../common/CreatePostFab';
import { TrendingHashtags } from '../common/TrendingHashtags';

interface MainLayoutProps {
  children: React.ReactNode;
}

export const MainLayout: React.FC<MainLayoutProps> = ({ children }) => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const isLargeScreen = useMediaQuery(theme.breakpoints.up('lg'));

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh', width: '100%' }}>
      <Header />
      {/* 左サイドバー（PCのみ表示） */}
      {!isMobile && <Sidebar />}

      {/* メインコンテンツエリア */}
      <Box
        sx={{
          flexGrow: 1,
          ml: { xs: 0, md: '280px' },
          mr: { xs: 0, lg: '320px' },
          width: { xs: '100%', md: 'calc(100% - 280px)', lg: 'calc(100% - 600px)' },
        }}
      >
        <Box
          sx={{
            py: 3,
            px: { xs: 2, sm: 3, md: 4 },
            maxWidth: { xs: '100%', md: '900px', lg: '700px' },
            margin: '0 auto',
          }}
        >
          {children}
        </Box>
      </Box>

      {/* 右サイドバー（大画面のみ表示） */}
      {isLargeScreen && (
        <Box
          sx={{
            position: 'fixed',
            top: '64px',
            right: 0,
            width: '320px',
            height: 'calc(100vh - 64px)',
            overflowY: 'auto',
            p: 2,
          }}
        >
          <TrendingHashtags />
        </Box>
      )}

      {/* フローティングアクションボタン（モバイルのみ） */}
      <CreatePostFab />
    </Box>
  );
};
