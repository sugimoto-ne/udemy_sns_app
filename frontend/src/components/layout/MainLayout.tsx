import React from 'react';
import { Box, useTheme, useMediaQuery } from '@mui/material';
import { Header } from './Header';
import { Sidebar } from './Sidebar';
import { CreatePostFab } from '../common/CreatePostFab';

interface MainLayoutProps {
  children: React.ReactNode;
}

export const MainLayout: React.FC<MainLayoutProps> = ({ children }) => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh', width: '100%' }}>
      <Header />
      {/* サイドバー（PCのみ表示） */}
      {!isMobile && <Sidebar />}

      {/* メインコンテンツエリア */}
      <Box
        sx={{
          flexGrow: 1,
          ml: { xs: 0, md: '280px' },
          width: { xs: '100%', md: 'calc(100% - 280px)' },
        }}
      >
        <Box
          sx={{
            py: 3,
            px: { xs: 2, sm: 3, md: 4 },
            maxWidth: { xs: '100%', md: '900px', lg: '1100px' },
            margin: '0 auto',
          }}
        >
          {children}
        </Box>
      </Box>

      {/* フローティングアクションボタン（モバイルのみ） */}
      <CreatePostFab />
    </Box>
  );
};
