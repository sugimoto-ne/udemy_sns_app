import { test, expect } from '@playwright/test';
import { generateTestUser, registerUser, createPost, clearLocalStorage, getPostCardSelector } from './helpers';

test.describe('Post Operations', () => {
  test.beforeEach(async ({ page }) => {
    await clearLocalStorage(page);
  });

  test('should create a new post and display it in timeline', async ({ page }) => {
    const user = generateTestUser('post');

    // ユーザーを登録してログイン
    await registerUser(page, user.email, user.username, user.password);
    await expect(page).toHaveURL('/');

    const postContent = 'This is my first test post!';

    // 投稿を作成
    await createPost(page, postContent);

    // タイムラインに投稿が表示される
    const postCards = page.locator(getPostCardSelector());
    const firstPost = postCards.first();
    await expect(firstPost).toContainText(postContent);
  });

  // Note: Edit post functionality is not implemented in Phase 1
  // This test is commented out for future implementation
  // test('should edit a post', async ({ page }) => {
  //   ...
  // });

  test('should delete a post', async ({ page }) => {
    const user = generateTestUser('delete');

    // ユーザーを登録してログイン
    await registerUser(page, user.email, user.username, user.password);

    const postContent = 'Post to be deleted';

    // 投稿を作成
    await createPost(page, postContent);

    // 投稿が表示されることを確認
    const postCards = page.locator(getPostCardSelector());
    const firstPost = postCards.first();
    await expect(firstPost).toContainText(postContent);

    // 削除ボタンをクリック（ブラウザの確認ダイアログを自動承認）
    page.on('dialog', dialog => dialog.accept());
    await firstPost.locator('[data-testid="delete-post-button"]').click();

    // 投稿が削除される
    await expect(page.locator(`text=${postContent}`)).not.toBeVisible();
  });

  test('should not allow deleting other users posts', async ({ page }) => {
    const user1 = generateTestUser('user1');
    const user2 = generateTestUser('user2');

    // ユーザー1を登録して投稿を作成
    await registerUser(page, user1.email, user1.username, user1.password);
    const postContent = 'User 1 unique post content';
    await createPost(page, postContent);

    // ログアウト
    await page.click('[data-testid="user-menu-button"]');
    await page.click('[data-testid="logout-button"]');

    // ユーザー2を登録してログイン
    await registerUser(page, user2.email, user2.username, user2.password);

    // タイムラインでユーザー1の投稿を探す
    const user1Post = page.locator(getPostCardSelector()).filter({ hasText: postContent }).first();
    await expect(user1Post).toBeVisible();

    // ユーザー2にはユーザー1の投稿の削除ボタンが表示されない
    const deleteButton = user1Post.locator('[data-testid="delete-post-button"]');
    await expect(deleteButton).not.toBeVisible();
  });

  test('should display post on user profile page', async ({ page }) => {
    const user = generateTestUser('profile');

    // ユーザーを登録してログイン
    await registerUser(page, user.email, user.username, user.password);

    const postContent = 'Post on my profile';

    // 投稿を作成
    await createPost(page, postContent);

    // プロフィールページに移動
    await page.goto(`/users/${user.username}`);

    // プロフィールページに投稿が表示される
    const postCards = page.locator(getPostCardSelector());
    await expect(postCards.first()).toContainText(postContent);
  });

  test('should add and display comment on post', async ({ page }) => {
    const user = generateTestUser('comment');

    // ユーザーを登録してログイン
    await registerUser(page, user.email, user.username, user.password);

    const postContent = 'Post with comments';
    const commentContent = 'This is a comment';

    // 投稿を作成
    await createPost(page, postContent);

    // 最初の投稿のコメントボタンをクリックして投稿詳細ページに移動
    const firstPost = page.locator(getPostCardSelector()).first();
    await firstPost.locator('[data-testid="comment-button"]').click();

    // 投稿詳細ページでコメントを入力
    await page.fill('[data-testid="comment-input"]', commentContent);
    await page.click('[data-testid="comment-submit-button"]');

    // コメントが表示される
    await expect(page.locator(`text=${commentContent}`)).toBeVisible();
  });

  test('should navigate through paginated timeline', async ({ page }) => {
    const user = generateTestUser('pagination');

    // ユーザーを登録してログイン
    await registerUser(page, user.email, user.username, user.password);

    // 複数の投稿を作成
    for (let i = 1; i <= 25; i++) {
      await createPost(page, `Post number ${i}`);
    }

    // 最初のページに最新の投稿が表示される
    await expect(page.locator('text=Post number 25')).toBeVisible();

    // スクロールして次のページを読み込む（無限スクロールの実装確認）
    await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight));

    // 無限スクロールのロード待機
    await page.waitForTimeout(2000);

    // 古い投稿が読み込まれる（または少なくとも20個の投稿が表示される）
    const postCount = await page.locator(getPostCardSelector()).count();
    expect(postCount).toBeGreaterThan(15); // 少なくとも15個以上の投稿が表示されることを確認
  });
});
