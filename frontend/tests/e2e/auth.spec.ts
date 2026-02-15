import { test, expect } from '@playwright/test';
import { generateTestUser, registerUser, login, logout, clearLocalStorage, hasAuthToken } from './helpers';

test.describe('Authentication Flow', () => {
  test.beforeEach(async ({ page }) => {
    // 各テストの前にlocalStorageをクリア
    await clearLocalStorage(page);
  });

  test('should register a new user successfully', async ({ page }) => {
    const user = generateTestUser('register');

    await registerUser(page, user.email, user.username, user.password);

    // 登録後、ホームページにリダイレクトされる
    await expect(page).toHaveURL('/');

    // 認証トークンが保存されている
    const hasToken = await hasAuthToken(page);
    expect(hasToken).toBeTruthy();
  });

  test('should not register with duplicate email', async ({ page }) => {
    const user = generateTestUser('duplicate');

    // 最初のユーザーを登録
    await registerUser(page, user.email, user.username, user.password);
    await expect(page).toHaveURL('/');

    // ログアウト
    await logout(page);

    // 同じメールで再度登録を試みる
    await page.goto('/register');
    await page.locator('[data-testid="username-input"]').fill('anotheruser');
    await page.locator('[data-testid="email-input"]').fill(user.email);
    await page.locator('[data-testid="password-input"]').fill(user.password);
    await page.locator('[data-testid="password-confirm-input"]').fill(user.password);

    // APIレスポンスを待つ（固定時間待機を削除）
    const responsePromise = page.waitForResponse(response =>
      response.url().includes('/api/v1/auth/register') && response.status() >= 400
    );
    await page.locator('[data-testid="register-submit-button"]').click();
    await responsePromise;

    // ページに留まっていることを確認（リダイレクトされていない）
    await expect(page).toHaveURL('/register');

    // エラーメッセージが表示される
    const errorAlert = page.locator('[role="alert"]');
    await expect(errorAlert).toBeVisible({ timeout: 3000 });

    // エラーメッセージの内容を確認（より具体的に）
    const errorText = await errorAlert.textContent();
    expect(errorText).toMatch(/メールアドレス|email/i);
  });

  test('should login with valid credentials', async ({ page }) => {
    const user = generateTestUser('login');

    // ユーザーを登録
    await registerUser(page, user.email, user.username, user.password);
    await expect(page).toHaveURL('/');

    // ログアウト
    await logout(page);

    // ログイン
    await login(page, user.email, user.password);

    // ホームページにリダイレクトされる
    await expect(page).toHaveURL('/');

    // 認証トークンが保存されている
    const hasToken = await hasAuthToken(page);
    expect(hasToken).toBeTruthy();
  });

  test('should not login with invalid email', async ({ page }) => {
    await page.goto('/login');
    await page.locator('[data-testid="email-input"]').fill('nonexistent@example.com');
    await page.locator('[data-testid="password-input"]').fill('Password123!');

    // APIレスポンスを待つ
    const responsePromise = page.waitForResponse(response =>
      response.url().includes('/api/v1/auth/login') && response.status() >= 400
    );
    await page.locator('[data-testid="login-submit-button"]').click();
    await responsePromise;

    // ログインページに留まる
    await expect(page).toHaveURL('/login');

    // エラーメッセージを確認
    const errorAlert = page.locator('[role="alert"]');
    await expect(errorAlert).toBeVisible({ timeout: 3000 });

    const errorText = await errorAlert.textContent();
    expect(errorText).toMatch(/メールアドレス|パスワード|email|password/i);
  });

  test('should not login with invalid password', async ({ page }) => {
    const user = generateTestUser('wrongpass');

    // ユーザーを登録
    await registerUser(page, user.email, user.username, user.password);
    await logout(page);

    // 間違ったパスワードでログイン
    await page.goto('/login');
    await page.locator('[data-testid="email-input"]').fill(user.email);
    await page.locator('[data-testid="password-input"]').fill('WrongPassword123!');

    // APIレスポンスを待つ
    const responsePromise = page.waitForResponse(response =>
      response.url().includes('/api/v1/auth/login') && response.status() >= 400
    );
    await page.locator('[data-testid="login-submit-button"]').click();
    await responsePromise;

    // ログインページに留まる
    await expect(page).toHaveURL('/login');

    // エラーメッセージを確認
    const errorAlert = page.locator('[role="alert"]');
    await expect(errorAlert).toBeVisible({ timeout: 3000 });

    const errorText = await errorAlert.textContent();
    expect(errorText).toMatch(/メールアドレス|パスワード|email|password/i);
  });

  test('should logout successfully', async ({ page }) => {
    const user = generateTestUser('logout');

    // ユーザーを登録してログイン
    await registerUser(page, user.email, user.username, user.password);
    await expect(page).toHaveURL('/');

    // ログアウト
    await logout(page);

    // ログインページにリダイレクトされる
    await expect(page).toHaveURL('/login');

    // 認証トークンが削除されている
    const hasToken = await hasAuthToken(page);
    expect(hasToken).toBeFalsy();
  });

  test('should redirect to login page when accessing protected route without authentication', async ({ page }) => {
    await page.goto('/');

    // ログインページにリダイレクトされる
    await expect(page).toHaveURL('/login');
  });

  test('should access protected route after authentication', async ({ page }) => {
    const user = generateTestUser('protected');

    // ユーザーを登録
    await registerUser(page, user.email, user.username, user.password);

    // ホームページにアクセスできる
    await expect(page).toHaveURL('/');

    // プロフィールページにアクセスできる
    await page.goto(`/users/${user.username}`);
    await expect(page).toHaveURL(`/users/${user.username}`);
  });
});
