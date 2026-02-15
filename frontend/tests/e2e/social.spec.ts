import { test, expect } from '@playwright/test';
import {
  generateTestUser,
  registerUser,
  createPost,
  followUser,
  unfollowUser,
  clearLocalStorage,
  getPostCardSelector,
  login
} from './helpers';

test.describe('Social Interactions', () => {
  test.beforeEach(async ({ page }) => {
    await clearLocalStorage(page);
  });

  test('should like and unlike a post', async ({ page }) => {
    const user = generateTestUser('like');

    // ユーザーを登録してログイン
    await registerUser(page, user.email, user.username, user.password);

    const postContent = 'Post to be liked';

    // 投稿を作成
    await createPost(page, postContent);

    // 最初の投稿を取得
    const firstPost = page.locator(getPostCardSelector()).first();

    // いいねボタンをクリック
    await firstPost.locator('[data-testid="like-button"]').click();

    // いいねカウントが増える
    await expect(firstPost.locator('[data-testid="like-count"]')).toContainText('1');

    // いいねボタンがアクティブになる
    await expect(firstPost.locator('[data-testid="unlike-button"]')).toBeVisible();

    // いいねを解除
    await firstPost.locator('[data-testid="unlike-button"]').click();

    // いいねカウントが減る
    await expect(firstPost.locator('[data-testid="like-count"]')).toContainText('0');

    // いいねボタンが非アクティブになる
    await expect(firstPost.locator('[data-testid="like-button"]')).toBeVisible();
  });

  test('should like other users posts', async ({ page, context }) => {
    const user1 = generateTestUser('user1');
    const user2 = generateTestUser('user2');

    // ユーザー1を登録して投稿を作成
    await registerUser(page, user1.email, user1.username, user1.password);
    const postContent = 'User 1 post to be liked';
    await createPost(page, postContent);

    // ログアウト
    await page.click('[data-testid="user-menu-button"]');
    await page.click('[data-testid="logout-button"]');

    // ユーザー2を登録してログイン
    await registerUser(page, user2.email, user2.username, user2.password);

    // タイムラインでユーザー1の投稿を見つける
    const user1Post = page.locator(getPostCardSelector()).filter({ hasText: postContent });
    await expect(user1Post).toBeVisible();

    // いいねボタンをクリック
    await user1Post.locator('[data-testid="like-button"]').click();

    // いいねカウントが増える
    await expect(user1Post.locator('[data-testid="like-count"]')).toContainText('1');
  });

  test('should follow and unfollow a user', async ({ page }) => {
    const user1 = generateTestUser('follower');
    const user2 = generateTestUser('following');

    // ユーザー1を登録
    await registerUser(page, user1.email, user1.username, user1.password);

    // ログアウト
    await page.click('[data-testid="user-menu-button"]');
    await page.click('[data-testid="logout-button"]');

    // ユーザー2を登録
    await registerUser(page, user2.email, user2.username, user2.password);

    // ログアウト
    await page.click('[data-testid="user-menu-button"]');
    await page.click('[data-testid="logout-button"]');

    // ユーザー1でログイン
    await login(page, user1.email, user1.password);

    // ユーザー2をフォロー
    await followUser(page, user2.username);

    // ユーザー1のプロフィールページに移動してフォロー数を確認
    await page.goto(`/users/${user1.username}`);
    await page.waitForSelector('[data-testid="following-count"]', { timeout: 10000 });

    // フォロー数が増える
    await expect(page.locator('[data-testid="following-count"]')).toContainText('1');

    // ユーザー2をアンフォロー
    await unfollowUser(page, user2.username);

    // ユーザー1のプロフィールページに移動してフォロー数を確認
    await page.goto(`/users/${user1.username}`);
    await page.waitForSelector('[data-testid="following-count"]', { timeout: 10000 });

    // フォロー数が減る
    await expect(page.locator('[data-testid="following-count"]')).toContainText('0');
  });

  test.skip('should display followed users posts in following timeline', async ({ page }) => {
    const user1 = generateTestUser('follower');
    const user2 = generateTestUser('followed');
    const user3 = generateTestUser('notfollowed');

    // ユーザー2を登録して投稿を作成
    await registerUser(page, user2.email, user2.username, user2.password);
    const post2Content = 'Post from followed user';
    await createPost(page, post2Content);
    await page.click('[data-testid="user-menu-button"]');
    await page.click('[data-testid="logout-button"]');

    // ユーザー3を登録して投稿を作成
    await registerUser(page, user3.email, user3.username, user3.password);
    const post3Content = 'Post from not followed user';
    await createPost(page, post3Content);
    await page.click('[data-testid="user-menu-button"]');
    await page.click('[data-testid="logout-button"]');

    // ユーザー1を登録してログイン
    await registerUser(page, user1.email, user1.username, user1.password);

    // ユーザー2をフォロー
    await followUser(page, user2.username);

    // ホームページに移動
    await page.goto('/');

    // フォロータイムラインに切り替え
    await page.click('[data-testid="following-timeline-tab"]');

    // ユーザー2の投稿が表示される
    await expect(page.locator(`text=${post2Content}`)).toBeVisible();

    // ユーザー3の投稿は表示されない
    await expect(page.locator(`text=${post3Content}`)).not.toBeVisible();

    // グローバルタイムラインに切り替え
    await page.click('[data-testid="global-timeline-tab"]');

    // すべての投稿が表示される
    await expect(page.locator(`text=${post2Content}`)).toBeVisible();
    await expect(page.locator(`text=${post3Content}`)).toBeVisible();
  });

  test.skip('should display followers and following lists', async ({ page }) => {
    const user1 = generateTestUser('user1');
    const user2 = generateTestUser('user2');

    // ユーザー1を登録
    await registerUser(page, user1.email, user1.username, user1.password);
    await page.click('[data-testid="user-menu-button"]');
    await page.click('[data-testid="logout-button"]');

    // ユーザー2を登録してユーザー1をフォロー
    await registerUser(page, user2.email, user2.username, user2.password);
    await followUser(page, user1.username);

    // ユーザー1のプロフィールページに移動
    await page.goto(`/users/${user1.username}`);

    // フォロワー数が表示されるまで待つ
    await page.waitForSelector('[data-testid="followers-count"]', { timeout: 10000 });

    // フォロワー数が1になる
    await expect(page.locator('[data-testid="followers-count"]')).toContainText('1');

    // フォロワーリンクが表示されるまで待つ
    await page.waitForSelector('[data-testid="followers-link"]', { timeout: 5000 });

    // フォロワーリストをクリック
    await page.click('[data-testid="followers-link"]');

    // フォロワーリストにユーザー2が表示される
    await expect(page.getByRole('heading', { name: user2.username })).toBeVisible();

    // ユーザー2のプロフィールページに移動
    await page.goto(`/users/${user2.username}`);

    // フォロー数が表示されるまで待つ
    await page.waitForSelector('[data-testid="following-count"]', { timeout: 10000 });

    // フォロー数が1になる
    await expect(page.locator('[data-testid="following-count"]')).toContainText('1');

    // フォローリンクが表示されるまで待つ
    await page.waitForSelector('[data-testid="following-link"]', { timeout: 5000 });

    // フォローリストをクリック
    await page.click('[data-testid="following-link"]');

    // フォローリストにユーザー1が表示される
    await expect(page.getByRole('heading', { name: user1.username })).toBeVisible();
  });

  test('should prevent following yourself', async ({ page }) => {
    const user = generateTestUser('self');

    // ユーザーを登録してログイン
    await registerUser(page, user.email, user.username, user.password);

    // 自分のプロフィールページに移動
    await page.goto(`/users/${user.username}`);

    // フォローボタンが表示されない
    await expect(page.locator('[data-testid="follow-button"]')).not.toBeVisible();
  });
});
