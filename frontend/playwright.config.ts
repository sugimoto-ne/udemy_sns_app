import { defineConfig, devices } from '@playwright/test';

// CI環境では8080、ローカルでは8081（テスト環境）を使用
const API_PORT = process.env.CI ? '8080' : '8081';
const API_BASE_URL = `http://localhost:${API_PORT}/api/v1`;

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
    // CI環境: port 8080 (GitHub Actionsで起動済みのバックエンド)
    // ローカル: port 8081 (docker-compose --profile test で起動するテスト用API)
    command: `VITE_API_BASE_URL=${API_BASE_URL} npm run dev`,
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
    timeout: 120 * 1000,
  },
});
