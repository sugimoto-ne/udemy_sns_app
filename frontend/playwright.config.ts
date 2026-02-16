import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: false, // E2Eテストは順次実行（DBの競合を避ける）
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1, // 並列実行数を1に設定（DBの競合を避ける）
  reporter: 'html',
  use: {
    baseURL: process.env.CI ? 'http://localhost:5173' : 'http://localhost:5174',
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
    command: process.env.CI ? 'npm run dev' : 'npm run dev:test',
    url: process.env.CI ? 'http://localhost:5173' : 'http://localhost:5174',
    reuseExistingServer: false,  // 常に新しいサーバーを起動（テスト用環境変数を確実に反映）
    timeout: 120 * 1000,
  },
});
