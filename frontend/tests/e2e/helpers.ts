import { Page, expect } from '@playwright/test';

// テスト用ユーザー生成
export function generateTestUser(prefix: string = 'test') {
  // タイムスタンプの下6桁のみ使用してユーザー名を20文字以内に収める
  const timestamp = Date.now();
  const shortId = timestamp.toString().slice(-6);
  return {
    email: `${prefix}${timestamp}@example.com`,
    username: `${prefix}${shortId}`,  // 例: test123456 (最大10文字程度)
    password: 'TestPass123!',  // バックエンドの要求に合わせて8文字以上
  };
}

// ユーザー登録
export async function registerUser(page: Page, email: string, username: string, password: string) {
  await page.goto('/register');

  await page.locator('[data-testid="username-input"]').fill(username);
  await page.locator('[data-testid="email-input"]').fill(email);
  await page.locator('[data-testid="password-input"]').fill(password);
  await page.locator('[data-testid="password-confirm-input"]').fill(password);

  await page.locator('[data-testid="register-submit-button"]').click();

  // 登録後、ホームページにリダイレクトされるのを待つ
  await page.waitForURL('/', { timeout: 10000 }).catch(() => {});
}

// ログイン
export async function login(page: Page, email: string, password: string) {
  await page.goto('/login');

  await page.locator('[data-testid="email-input"]').fill(email);
  await page.locator('[data-testid="password-input"]').fill(password);
  await page.locator('[data-testid="login-submit-button"]').click();

  // ログイン成功を待つ（ホームページにリダイレクトされる）
  await page.waitForURL('/', { timeout: 10000 }).catch(() => {});
}

// ログアウト
export async function logout(page: Page) {
  // ユーザーメニューを開く
  await page.click('[data-testid="user-menu-button"]');
  // ログアウトボタンをクリック
  await page.click('[data-testid="logout-button"]');
  // ログインページにリダイレクトされるのを待つ
  await page.waitForURL('/login');
}

// 投稿を作成
export async function createPost(page: Page, content: string) {
  await page.fill('[data-testid="post-input"]', content);
  await page.click('[data-testid="post-submit-button"]');

  // 投稿が表示されるのを待つ
  const postCardSelector = getPostCardSelector();
  await page.waitForSelector(postCardSelector, { timeout: 10000 });
  // 投稿内容が表示されることを確認
  await expect(page.locator('text=' + content).first()).toBeVisible({ timeout: 5000 });
}

// 投稿をいいね
export async function likePost(page: Page, postId: string) {
  await page.click(`[data-testid="like-button-${postId}"]`);
}

// ユーザーをフォロー
export async function followUser(page: Page, username: string) {
  await page.goto(`/users/${username}`);

  // フォローボタンが表示されるまで待つ
  await page.waitForSelector('[data-testid="follow-button"]', { timeout: 10000 });
  await page.click('[data-testid="follow-button"]');

  // フォローボタンがアンフォローに変わるのを待つ
  await page.waitForSelector('[data-testid="unfollow-button"]', { timeout: 5000 });
  await expect(page.locator('[data-testid="unfollow-button"]')).toBeVisible();
}

// ユーザーをアンフォロー
export async function unfollowUser(page: Page, username: string) {
  await page.goto(`/users/${username}`);

  // アンフォローボタンが表示されるまで待つ
  await page.waitForSelector('[data-testid="unfollow-button"]', { timeout: 10000 });
  await page.click('[data-testid="unfollow-button"]');

  // アンフォローボタンがフォローに変わるのを待つ
  await page.waitForSelector('[data-testid="follow-button"]', { timeout: 5000 });
  await expect(page.locator('[data-testid="follow-button"]')).toBeVisible();
}

// localStorageをクリア
export async function clearLocalStorage(page: Page) {
  // ページに移動してからlocalStorageをクリア
  await page.goto('/');
  await page.evaluate(() => localStorage.clear());
}

// 認証トークンが存在するか確認
export async function hasAuthToken(page: Page): Promise<boolean> {
  const token = await page.evaluate(() => localStorage.getItem('token'));
  return token !== null;
}

// 投稿カードセレクタを取得（post-inputやpost-submit-buttonを除外）
export function getPostCardSelector(): string {
  return 'div[data-testid^="post-"]:not([data-testid="post-input"]):not([data-testid="post-submit-button"])';
}
