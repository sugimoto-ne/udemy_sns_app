import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: false, // E2Eテストは順次実行（DBの競合を避ける）
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1, // 並列実行数を1に設定（DBの競合を避ける）
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],

  webServer: {
    // テスト用環境変数を使用してフロントエンドを起動
    // VITE_API_BASE_URLを直接指定してテスト用APIサーバー(port 8081)に接続
    command: 'VITE_API_BASE_URL=http://localhost:8081/api/v1 npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
    timeout: 120 * 1000,
  },
});
