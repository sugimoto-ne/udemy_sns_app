# Phase 2 - フロントエンド開発TODO（中優先度機能）

## 🎯 目標
ユーザー体験を向上させる追加機能のUI/UX実装

---

## 📊 Phase 2 機能概要

- ✅ ハッシュタグ機能
- ✅ 複数画像添付機能
- ✅ ブックマーク機能
- ✅ パスワードリセット機能
- ✅ メールアドレス認証機能

---

## 🏷️ ハッシュタグ機能

### 1. 型定義
- [ ] `src/types/hashtag.ts`
```typescript
export interface Hashtag {
  id: number;
  name: string;
  posts_count?: number;
}
```

- [ ] `src/types/post.ts` 更新
  - [ ] Post型に `hashtags: string[]` を追加

### 2. ハッシュタグAPI
- [ ] `src/api/hashtags.ts`
  - [ ] `getPostsByHashtag(hashtagName: string, limit, cursor)`
  - [ ] `getTrendingHashtags(limit)`

### 3. ハッシュタグカスタムフック
- [ ] `src/hooks/useHashtags.ts`
  - [ ] `useHashtagPosts(hashtagName)` - Infinite Query
  - [ ] `useTrendingHashtags()`

### 4. ハッシュタグ表示コンポーネント
- [ ] `src/components/hashtag/HashtagChip.tsx`
  - [ ] クリック可能なハッシュタグチップ
  - [ ] MUI Chip使用
  - [ ] クリックでハッシュタグページに遷移

- [ ] `src/components/hashtag/HashtagList.tsx`
  - [ ] ハッシュタグ配列を表示

### 5. PostCardにハッシュタグ表示
- [ ] `src/components/post/PostCard.tsx` 更新
  - [ ] 投稿内容からハッシュタグを抽出して強調表示（正規表現）
  - [ ] または、バックエンドから返されたハッシュタグを下部に表示

### 6. トレンドハッシュタグウィジェット
- [ ] `src/components/hashtag/TrendingHashtags.tsx`
  - [ ] トレンドハッシュタグ一覧表示
  - [ ] サイドバーに配置
  - [ ] クリックでハッシュタグページへ

### 7. ハッシュタグページ
- [ ] `src/pages/HashtagPage.tsx`
  - [ ] パラメータからハッシュタグ名取得
  - [ ] ハッシュタグに関連する投稿一覧表示
  - [ ] 無限スクロール

### 8. ルート追加
- [ ] `src/App.tsx`
  - [ ] `/hashtag/:name` ルート追加

---

## 🖼️ 複数画像添付機能

### 9. PostForm更新
- [ ] `src/components/post/PostForm.tsx` 更新
  - [ ] 複数ファイル選択対応（最大4件）
  - [ ] プレビュー表示（グリッドレイアウト）
  - [ ] 個別削除ボタン
  - [ ] ドラッグ&ドロップ対応（オプション）

### 10. MediaPreview更新
- [ ] `src/components/post/MediaPreview.tsx` 更新
  - [ ] グリッドレイアウト
    - [ ] 1枚: 大きく表示
    - [ ] 2枚: 2カラム
    - [ ] 3枚: 2+1レイアウト
    - [ ] 4枚: 2x2グリッド
  - [ ] クリックで拡大表示（ライトボックス）

### 11. ライトボックス（画像拡大表示）
- [ ] ライブラリインストール（オプション）
```bash
npm install react-image-lightbox
```

- [ ] `src/components/common/ImageLightbox.tsx`
  - [ ] 画像をフルスクリーン表示
  - [ ] 左右スワイプで切り替え

---

## 🔖 ブックマーク機能

### 12. ブックマークAPI
- [ ] `src/api/bookmarks.ts`
  - [ ] `getBookmarks(limit, cursor)`
  - [ ] `bookmarkPost(postId)`
  - [ ] `unbookmarkPost(postId)`

### 13. ブックマークカスタムフック
- [ ] `src/hooks/useBookmarks.ts`
  - [ ] `useBookmarks()` - Infinite Query
  - [ ] `useBookmarkPost()` - Mutation（楽観的更新）
  - [ ] `useUnbookmarkPost()` - Mutation（楽観的更新）

### 14. PostCardにブックマークボタン追加
- [ ] `src/components/post/PostCard.tsx` 更新
  - [ ] ブックマークアイコンボタン追加（BookmarkBorder / Bookmark）
  - [ ] クリックでブックマーク追加/削除
  - [ ] `is_bookmarked` に応じてアイコン切り替え

### 15. ブックマークページ
- [ ] `src/pages/BookmarksPage.tsx`
  - [ ] ブックマークした投稿一覧表示
  - [ ] PostList使用
  - [ ] 無限スクロール

### 16. ナビゲーションにブックマークリンク追加
- [ ] `src/components/common/Sidebar.tsx` 更新
  - [ ] ブックマークページへのリンク追加

### 17. ルート追加
- [ ] `src/App.tsx`
  - [ ] `/bookmarks` ルート追加（プライベート）

---

## 🔑 パスワードリセット機能

### 18. パスワードリセットAPI
- [ ] `src/api/auth.ts` 更新
  - [ ] `requestPasswordReset(email)`
  - [ ] `confirmPasswordReset(token, newPassword)`

### 19. パスワードリセットフォーム
- [ ] `src/components/auth/PasswordResetRequestForm.tsx`
  - [ ] メールアドレス入力
  - [ ] 送信ボタン
  - [ ] 成功メッセージ表示

- [ ] `src/components/auth/PasswordResetConfirmForm.tsx`
  - [ ] 新しいパスワード入力
  - [ ] パスワード確認入力
  - [ ] バリデーション
  - [ ] 送信ボタン

### 20. パスワードリセットページ
- [ ] `src/pages/PasswordResetRequestPage.tsx`
  - [ ] PasswordResetRequestForm表示
  - [ ] 送信後の案内表示

- [ ] `src/pages/PasswordResetConfirmPage.tsx`
  - [ ] URLパラメータからトークン取得
  - [ ] PasswordResetConfirmForm表示
  - [ ] 成功後にログインページへリダイレクト

### 21. ログインページにリンク追加
- [ ] `src/pages/LoginPage.tsx` 更新
  - [ ] 「パスワードをお忘れですか？」リンク追加

### 22. ルート追加
- [ ] `src/App.tsx`
  - [ ] `/password-reset/request` ルート追加
  - [ ] `/password-reset/confirm/:token` ルート追加

---

## ✉️ メールアドレス認証機能

### 23. メール認証API
- [ ] `src/api/auth.ts` 更新
  - [ ] `verifyEmail(token)`
  - [ ] `resendVerificationEmail()`

### 24. メール認証バナー
- [ ] `src/components/auth/EmailVerificationBanner.tsx`
  - [ ] 未認証時に表示するバナー
  - [ ] 「確認メールを再送信」ボタン
  - [ ] 閉じるボタン

### 25. メール認証ページ
- [ ] `src/pages/EmailVerificationPage.tsx`
  - [ ] URLパラメータからトークン取得
  - [ ] 自動的に認証API呼び出し
  - [ ] 成功/失敗メッセージ表示
  - [ ] ホームページへのリンク

### 26. ホームページにバナー表示
- [ ] `src/pages/HomePage.tsx` 更新
  - [ ] 未認証ユーザーにEmailVerificationBanner表示
  - [ ] `user.email_verified` で判定

### 27. ルート追加
- [ ] `src/App.tsx`
  - [ ] `/email/verify/:token` ルート追加

---

## 🎨 UI/UX改善

### 28. 通知・フィードバック
- [ ] トースト通知ライブラリ導入（オプション）
```bash
npm install notistack
```

- [ ] `src/App.tsx` または `src/main.tsx` 更新
  - [ ] SnackbarProvider追加

- [ ] 各Mutation成功時にトースト表示
  - [ ] 投稿作成成功
  - [ ] ブックマーク追加/削除
  - [ ] フォロー/フォロー解除
  - [ ] プロフィール更新成功

### 29. スケルトンローディング
- [ ] MUI Skeleton使用
- [ ] `src/components/common/PostSkeleton.tsx`
  - [ ] PostCardのスケルトン
- [ ] Timeline等でローディング中にスケルトン表示

### 30. エラーページ
- [ ] `src/pages/NotFoundPage.tsx`
  - [ ] 404エラーページ

- [ ] `src/pages/ErrorPage.tsx`
  - [ ] 一般的なエラーページ

### 31. 画像最適化
- [ ] 画像遅延読み込み
  - [ ] `loading="lazy"` 属性追加
- [ ] サムネイル生成（バックエンド対応必要）

---

## ✅ テスト

### 32. Phase 2機能のテスト
- [ ] ハッシュタグ機能
  - [ ] 投稿内のハッシュタグがクリック可能
  - [ ] ハッシュタグページで関連投稿が表示される
  - [ ] トレンドハッシュタグが表示される

- [ ] 複数画像添付
  - [ ] 最大4枚の画像を選択できる
  - [ ] プレビューがグリッド表示される
  - [ ] 個別削除ができる

- [ ] ブックマーク機能
  - [ ] ブックマークボタンで追加/削除できる
  - [ ] ブックマークページで一覧が表示される

- [ ] パスワードリセット
  - [ ] リセットリクエストが送信できる
  - [ ] トークンでパスワードを変更できる

- [ ] メール認証
  - [ ] 未認証時にバナーが表示される
  - [ ] 確認メールを再送信できる
  - [ ] トークンで認証できる

---

## 📚 ドキュメント更新

### 33. README更新
- [ ] Phase 2機能の説明追加
- [ ] スクリーンショット追加（オプション）

---

## 🚀 デプロイ

### 34. ビルド・デプロイ
- [ ] `npm run build`
- [ ] Firebase Hosting にデプロイ

---

## ✅ Phase 2 完了チェックリスト

- [ ] 投稿にハッシュタグが表示される
- [ ] ハッシュタグをクリックして検索できる
- [ ] トレンドハッシュタグが表示される
- [ ] 複数画像を投稿できる（最大4枚）
- [ ] 画像がグリッドレイアウトで表示される
- [ ] 投稿をブックマークできる
- [ ] ブックマーク一覧を表示できる
- [ ] パスワードリセット機能が動作する
- [ ] メール認証機能が動作する
- [ ] すべてのPhase 2機能が正常に動作する

---

## 📝 開発の進め方

1. **ハッシュタグ機能** (項目1-8)
2. **複数画像添付** (項目9-11)
3. **ブックマーク機能** (項目12-17)
4. **パスワードリセット** (項目18-22)
5. **メール認証** (項目23-27)
6. **UI/UX改善** (項目28-31)
7. **テスト** (項目32)
8. **ドキュメント** (項目33)
9. **デプロイ** (項目34)

---

**作成日**: 2026-02-14
**最終更新**: 2026-02-14
