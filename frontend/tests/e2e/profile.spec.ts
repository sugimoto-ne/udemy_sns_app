import { test, expect } from '@playwright/test';
import { generateTestUser, registerUser, createPost, clearLocalStorage } from './helpers';

test.describe('Profile Management', () => {
  test.beforeEach(async ({ page }) => {
    await clearLocalStorage(page);
  });

  test('should display user profile information', async ({ page }) => {
    const user = generateTestUser('profile');

    // ユーザーを登録してログイン
    await registerUser(page, user.email, user.username, user.password);

    // プロフィールページに移動
    await page.goto(`/users/${user.username}`);

    // ユーザー名が表示される
    await expect(page.locator('[data-testid="profile-username"]')).toContainText(user.username);

    // メールアドレスは表示されない（プライバシー）
    await expect(page.locator(`text=${user.email}`)).not.toBeVisible();
  });

  // Note: Profile editing is not implemented in Phase 1
  // test.skip('should edit profile information', async ({ page }) => {
  //   ...
  // });

  test('should display user posts on profile page', async ({ page }) => {
    const user = generateTestUser('posts');

    // ユーザーを登録してログイン
    await registerUser(page, user.email, user.username, user.password);

    const post1Content = 'First post on my profile';
    const post2Content = 'Second post on my profile';

    // 2つの投稿を作成
    await createPost(page, post1Content);
    await createPost(page, post2Content);

    // プロフィールページに移動
    await page.goto(`/users/${user.username}`);

    // 両方の投稿が表示される
    await expect(page.locator(`text=${post1Content}`)).toBeVisible();
    await expect(page.locator(`text=${post2Content}`)).toBeVisible();
  });

  test('should display post count, followers, and following count', async ({ page }) => {
    const user1 = generateTestUser('user1');
    const user2 = generateTestUser('user2');

    // ユーザー1を登録して投稿を作成
    await registerUser(page, user1.email, user1.username, user1.password);
    await createPost(page, 'Test post');

    // ログアウト
    await page.click('[data-testid="user-menu-button"]');
    await page.click('[data-testid="logout-button"]');

    // ユーザー2を登録してユーザー1をフォロー
    await registerUser(page, user2.email, user2.username, user2.password);
    await page.goto(`/users/${user1.username}`);

    // フォローボタンが表示されるまで待つ
    await page.waitForSelector('[data-testid="follow-button"]', { timeout: 5000 });
    await page.click('[data-testid="follow-button"]');

    // フォロー完了を待つ（unfollowボタンが表示される）
    await page.waitForSelector('[data-testid="unfollow-button"]', { timeout: 5000 });

    // ユーザー1のプロフィールを確認
    await expect(page.locator('[data-testid="followers-count"]')).toContainText('1');
    await expect(page.locator('[data-testid="following-count"]')).toContainText('0');

    // ユーザー2のプロフィールに移動
    await page.goto(`/users/${user2.username}`);

    await expect(page.locator('[data-testid="followers-count"]')).toContainText('0');
    await expect(page.locator('[data-testid="following-count"]')).toContainText('1');
  });

  test('should show follow button on other users profiles', async ({ page }) => {
    const user1 = generateTestUser('user1');
    const user2 = generateTestUser('user2');

    // ユーザー1を登録
    await registerUser(page, user1.email, user1.username, user1.password);
    await page.click('[data-testid="user-menu-button"]');
    await page.click('[data-testid="logout-button"]');

    // ユーザー2を登録してログイン
    await registerUser(page, user2.email, user2.username, user2.password);

    // ユーザー1のプロフィールページに移動
    await page.goto(`/users/${user1.username}`);

    // フォローボタンが表示されるまで待つ
    await page.waitForSelector('[data-testid="follow-button"]', { timeout: 5000 });
    await expect(page.locator('[data-testid="follow-button"]')).toBeVisible();
  });

  // Note: Profile avatar upload is not implemented in Phase 1
  // test.skip('should upload and display profile avatar', async ({ page }) => {
  //   ...
  // });

  test('should view other users profile', async ({ page }) => {
    const user1 = generateTestUser('viewer');
    const user2 = generateTestUser('viewed');

    // ユーザー2を登録して投稿を作成
    await registerUser(page, user2.email, user2.username, user2.password);
    await createPost(page, 'Public post');
    await page.click('[data-testid="user-menu-button"]');
    await page.click('[data-testid="logout-button"]');

    // ユーザー1を登録してログイン
    await registerUser(page, user1.email, user1.username, user1.password);

    // ユーザー2のプロフィールページに移動
    await page.goto(`/users/${user2.username}`);

    // ユーザー2のユーザー名が表示される
    await expect(page.locator('[data-testid="profile-username"]')).toContainText(user2.username);

    // ユーザー2の投稿が表示される
    await expect(page.locator('text=Public post')).toBeVisible();

    // フォローボタンが表示されるまで待つ
    await page.waitForSelector('[data-testid="follow-button"]', { timeout: 5000 });
    await expect(page.locator('[data-testid="follow-button"]')).toBeVisible();
  });
});
