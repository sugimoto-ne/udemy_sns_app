import { createTheme } from '@mui/material/styles';
import type { Theme } from '@mui/material/styles';

// 共通のタイポグラフィ設定
const commonTypography = {
  fontFamily: [
    '-apple-system',
    'BlinkMacSystemFont',
    '"Segoe UI"',
    'Roboto',
    '"Helvetica Neue"',
    'Arial',
    'sans-serif',
  ].join(','),
  h1: {
    fontSize: '2rem',
    fontWeight: 600,
  },
  h2: {
    fontSize: '1.75rem',
    fontWeight: 600,
  },
  h3: {
    fontSize: '1.5rem',
    fontWeight: 600,
  },
  h4: {
    fontSize: '1.25rem',
    fontWeight: 600,
  },
  h5: {
    fontSize: '1.1rem',
    fontWeight: 600,
  },
  h6: {
    fontSize: '1rem',
    fontWeight: 600,
  },
  body1: {
    fontSize: '1rem',
    lineHeight: 1.5,
  },
  body2: {
    fontSize: '0.875rem',
    lineHeight: 1.43,
  },
};

// 共通のブレークポイント設定
const commonBreakpoints = {
  values: {
    xs: 0,      // モバイル
    sm: 600,    // タブレット
    md: 960,    // デスクトップ小
    lg: 1280,   // デスクトップ大
    xl: 1920,
  },
};

// 共通のコンポーネントスタイル
const commonComponents = {
  MuiButton: {
    styleOverrides: {
      root: {
        textTransform: 'none' as const,
        borderRadius: 20,
        padding: '8px 24px',
      },
    },
  },
  MuiCard: {
    styleOverrides: {
      root: {
        borderRadius: 12,
      },
    },
  },
  MuiTextField: {
    styleOverrides: {
      root: {
        '& .MuiOutlinedInput-root': {
          borderRadius: 12,
        },
      },
    },
  },
};

// デフォルトテーマ (Twitter風ブルー)
export const defaultTheme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#1DA1F2',
      light: '#42b0f5',
      dark: '#1a8cd8',
      contrastText: '#fff',
    },
    secondary: {
      main: '#f50057',
      light: '#ff4081',
      dark: '#c51162',
      contrastText: '#fff',
    },
    background: {
      default: '#f5f8fa',
      paper: '#ffffff',
    },
    text: {
      primary: '#14171a',
      secondary: '#657786',
    },
  },
  typography: commonTypography,
  breakpoints: commonBreakpoints,
  components: commonComponents,
});

// ダークテーマ
export const darkTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#1DA1F2',
      light: '#42b0f5',
      dark: '#1a8cd8',
      contrastText: '#fff',
    },
    secondary: {
      main: '#f50057',
      light: '#ff4081',
      dark: '#c51162',
      contrastText: '#fff',
    },
    background: {
      default: '#15202b',
      paper: '#192734',
    },
    text: {
      primary: '#ffffff',
      secondary: '#8899a6',
    },
  },
  typography: commonTypography,
  breakpoints: commonBreakpoints,
  components: {
    ...commonComponents,
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          border: '1px solid rgba(255, 255, 255, 0.12)',
        },
      },
    },
  },
});

// パープルテーマ
export const purpleTheme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#9c27b0',
      light: '#ba68c8',
      dark: '#7b1fa2',
      contrastText: '#fff',
    },
    secondary: {
      main: '#ff6090',
      light: '#ff8fab',
      dark: '#cc4d73',
      contrastText: '#fff',
    },
    background: {
      default: '#faf5ff',
      paper: '#ffffff',
    },
    text: {
      primary: '#2d1b3d',
      secondary: '#6b4f7a',
    },
  },
  typography: commonTypography,
  breakpoints: commonBreakpoints,
  components: commonComponents,
});

// グリーンテーマ
export const greenTheme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#2e7d32',
      light: '#60ad5e',
      dark: '#005005',
      contrastText: '#fff',
    },
    secondary: {
      main: '#ff9800',
      light: '#ffb74d',
      dark: '#f57c00',
      contrastText: '#000',
    },
    background: {
      default: '#f1f8f4',
      paper: '#ffffff',
    },
    text: {
      primary: '#1b3a1f',
      secondary: '#4a6b4e',
    },
  },
  typography: commonTypography,
  breakpoints: commonBreakpoints,
  components: commonComponents,
});

// オレンジテーマ
export const orangeTheme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#ff6b35',
      light: '#ff8c5f',
      dark: '#e55527',
      contrastText: '#fff',
    },
    secondary: {
      main: '#004e89',
      light: '#3375a8',
      dark: '#003761',
      contrastText: '#fff',
    },
    background: {
      default: '#fff8f5',
      paper: '#ffffff',
    },
    text: {
      primary: '#2d1810',
      secondary: '#6b5749',
    },
  },
  typography: commonTypography,
  breakpoints: commonBreakpoints,
  components: commonComponents,
});

// テーマの型定義
export type ThemeName = 'default' | 'dark' | 'purple' | 'green' | 'orange';

// テーママップ
export const themes: Record<ThemeName, Theme> = {
  default: defaultTheme,
  dark: darkTheme,
  purple: purpleTheme,
  green: greenTheme,
  orange: orangeTheme,
};

// テーマの表示名
export const themeDisplayNames: Record<ThemeName, string> = {
  default: 'デフォルト (ブルー)',
  dark: 'ダーク',
  purple: 'パープル',
  green: 'グリーン',
  orange: 'オレンジ',
};
