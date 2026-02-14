import React from 'react';
import {
  Box,
  Typography,
  Paper,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Switch,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
} from '@mui/material';
import {
  Palette as PaletteIcon,
  Notifications as NotificationsIcon,
  Security as SecurityIcon,
  Language as LanguageIcon,
} from '@mui/icons-material';
import { MainLayout } from '../components/layout/MainLayout';
import { useTheme } from '../contexts/ThemeContext';
import type { ThemeName } from '../theme/themes';

export const SettingsPage: React.FC = () => {
  const { currentTheme, setTheme } = useTheme();
  const [notificationsEnabled, setNotificationsEnabled] = React.useState(true);
  const [language, setLanguage] = React.useState('ja');

  const handleThemeChange = (event: React.ChangeEvent<{ value: unknown }>) => {
    setTheme(event.target.value as ThemeName);
  };

  const themeOptions = [
    { value: 'default', label: 'デフォルト（ブルー）' },
    { value: 'dark', label: 'ダーク' },
    { value: 'purple', label: 'パープル' },
    { value: 'green', label: 'グリーン' },
    { value: 'orange', label: 'オレンジ' },
  ];

  return (
    <MainLayout>
      <Box>
        <Typography variant="h4" gutterBottom fontWeight="bold" sx={{ mb: 3 }}>
          設定
        </Typography>

        {/* テーマ設定 */}
        <Paper sx={{ mb: 3 }}>
          <Box sx={{ p: 2, borderBottom: 1, borderColor: 'divider' }}>
            <Typography variant="h6" fontWeight="bold">
              外観
            </Typography>
          </Box>
          <List sx={{ p: 0 }}>
            <ListItem sx={{ py: 2 }}>
              <ListItemIcon>
                <PaletteIcon />
              </ListItemIcon>
              <ListItemText
                primary="テーマ"
                secondary="アプリケーションの配色テーマを選択"
              />
              <FormControl sx={{ minWidth: 200 }}>
                <InputLabel id="theme-select-label">テーマ</InputLabel>
                <Select
                  labelId="theme-select-label"
                  value={currentTheme}
                  onChange={handleThemeChange as any}
                  label="テーマ"
                  size="small"
                >
                  {themeOptions.map((option) => (
                    <MenuItem key={option.value} value={option.value}>
                      {option.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </ListItem>
          </List>
        </Paper>

        {/* 通知設定 */}
        <Paper sx={{ mb: 3 }}>
          <Box sx={{ p: 2, borderBottom: 1, borderColor: 'divider' }}>
            <Typography variant="h6" fontWeight="bold">
              通知
            </Typography>
          </Box>
          <List sx={{ p: 0 }}>
            <ListItem sx={{ py: 2 }}>
              <ListItemIcon>
                <NotificationsIcon />
              </ListItemIcon>
              <ListItemText
                primary="プッシュ通知"
                secondary="新しい通知を受け取る"
              />
              <Switch
                edge="end"
                checked={notificationsEnabled}
                onChange={(e) => setNotificationsEnabled(e.target.checked)}
              />
            </ListItem>
          </List>
        </Paper>

        {/* 言語設定 */}
        <Paper sx={{ mb: 3 }}>
          <Box sx={{ p: 2, borderBottom: 1, borderColor: 'divider' }}>
            <Typography variant="h6" fontWeight="bold">
              言語
            </Typography>
          </Box>
          <List sx={{ p: 0 }}>
            <ListItem sx={{ py: 2 }}>
              <ListItemIcon>
                <LanguageIcon />
              </ListItemIcon>
              <ListItemText
                primary="表示言語"
                secondary="インターフェースの表示言語"
              />
              <FormControl sx={{ minWidth: 200 }}>
                <InputLabel id="language-select-label">言語</InputLabel>
                <Select
                  labelId="language-select-label"
                  value={language}
                  onChange={(e) => setLanguage(e.target.value)}
                  label="言語"
                  size="small"
                >
                  <MenuItem value="ja">日本語</MenuItem>
                  <MenuItem value="en">English</MenuItem>
                </Select>
              </FormControl>
            </ListItem>
          </List>
        </Paper>

        {/* プライバシー設定 */}
        <Paper>
          <Box sx={{ p: 2, borderBottom: 1, borderColor: 'divider' }}>
            <Typography variant="h6" fontWeight="bold">
              プライバシーとセキュリティ
            </Typography>
          </Box>
          <List sx={{ p: 0 }}>
            <ListItem sx={{ py: 2 }}>
              <ListItemIcon>
                <SecurityIcon />
              </ListItemIcon>
              <ListItemText
                primary="非公開アカウント"
                secondary="フォロワーのみがあなたの投稿を見ることができます"
              />
              <Switch edge="end" />
            </ListItem>
          </List>
        </Paper>

        {/* アプリ情報 */}
        <Box sx={{ mt: 4, textAlign: 'center' }}>
          <Typography variant="body2" color="text.secondary">
            バージョン 1.0.0
          </Typography>
        </Box>
      </Box>
    </MainLayout>
  );
};
