# Phase 2 実装完了レポート

## 📅 完了日
2026-02-16

## ✅ 実装完了機能

### 1. メール認証機能（Email Verification）

#### 実装内容
- ✅ 新規登録後、メール確認待ちページへ自動遷移
- ✅ 未認証ユーザーのログイン防止（403エラー）
- ✅ 認証メール再送信機能
- ✅ メール認証後のログイン許可

#### 実装ファイル
**バックエンド:**
- `backend/internal/services/auth_service.go` (lines 80-83)
  - ログイン時の`email_verified`チェック追加
- `backend/internal/handlers/auth_handler.go` (lines 121-123)
  - 403エラーレスポンス実装

**フロントエンド:**
- `frontend/src/pages/EmailVerificationPendingPage.tsx` (新規作成)
  - メール確認待ちページ
  - 再送信機能付き
- `frontend/src/components/auth/RegisterForm.tsx` (line 42)
  - 登録後に認証待ちページへリダイレクト
- `frontend/src/components/auth/LoginForm.tsx` (lines 38-42)
  - 403エラー時に認証待ちページへリダイレクト
- `frontend/src/App.tsx` (line 45)
  - `/auth/email/verify-pending` ルート追加

#### 動作フロー
```
1. ユーザー登録
   ↓
2. メール送信（Resend API）
   ↓
3. /auth/email/verify-pending へ遷移
   ↓
4. メール内のリンクをクリック
   ↓
5. email_verified = true に更新
   ↓
6. ログイン可能
```

---

### 2. ブックマーク機能（Bookmark）

#### 実装内容
- ✅ 投稿のブックマーク追加/解除
- ✅ バックエンドから`is_bookmarked`フィールドを返却
- ✅ 楽観的UI更新（即座にON/OFF切り替え）
- ✅ エラー時の自動ロールバック
- ✅ ブックマーク一覧ページでの表示

#### 実装ファイル
**バックエンド:**
- `backend/internal/services/post_service.go`
  - Lines 181-194: タイムライン取得時のブックマーク状態を一括取得
  - Lines 222-225: 投稿詳細取得時のブックマーク状態をチェック

**フロントエンド:**
- `frontend/src/hooks/useBookmarks.ts` (lines 12-82)
  - `useBookmarkPost`: 楽観的更新実装
  - `useUnbookmarkPost`: 楽観的更新実装
  - `onMutate`: UI即座更新
  - `onError`: エラー時ロールバック
  - `onSettled`: サーバーと同期

#### 楽観的更新パターン
```typescript
onMutate: async (postId) => {
  // 1. 進行中のクエリをキャンセル
  await queryClient.cancelQueries({ queryKey: ['timeline'] });

  // 2. 現在のデータを保存（ロールバック用）
  const previousData = queryClient.getQueryData(['timeline']);

  // 3. UIを即座に更新
  queryClient.setQueriesData({ queryKey: ['timeline'] }, (old: any) => {
    // is_bookmarked を即座に更新
  });

  return { previousData };
},

onError: (err, postId, context) => {
  // エラー時は元のデータに戻す
  if (context?.previousData) {
    queryClient.setQueryData(['timeline'], context.previousData);
  }
},

onSettled: () => {
  // 最終的にサーバーと同期
  queryClient.invalidateQueries({ queryKey: ['bookmarks'] });
  queryClient.invalidateQueries({ queryKey: ['timeline'] });
}
```

---

### 3. いいね機能の楽観的更新（Like with Optimistic Updates）

#### 実装内容
- ✅ いいねボタンクリック時の即座の反映
- ✅ いいね数のリアルタイム更新
- ✅ エラー時の自動ロールバック

#### 実装ファイル
**フロントエンド:**
- `frontend/src/hooks/usePosts.ts`
  - Lines 77-110: `useLikePost` 楽観的更新
  - Lines 119-146: `useUnlikePost` 楽観的更新

---

### 4. 画像アップロード機能（Firebase Storage）

#### 実装内容
- ✅ Firebase Storageへの画像アップロード
- ✅ 署名付きURL生成（7日間有効）
- ✅ 投稿と画像の同期表示
- ✅ アップロード失敗時の投稿ロールバック
- ✅ 403 Forbiddenエラーの解決

#### 実装ファイル
**バックエンド:**
- `backend/internal/services/firebase_storage_service.go`
  - Line 31: `config.AppConfig` を使用（`config.LoadConfig()` から変更）
  - Lines 92-102: 署名付きURL生成実装

```go
// 署名付きURL（有効期限あり）を生成
url, err := s.bucket.SignedURL(objectPath, &storage.SignedURLOptions{
    Method:  "GET",
    Expires: time.Now().Add(7 * 24 * time.Hour), // 7日間有効
})
```

**フロントエンド:**
- `frontend/src/pages/HomePage.tsx` (lines 28-62)
  - 画像アップロードの同期処理実装
  - アップロード失敗時の投稿削除（ロールバック）

```typescript
const handleCreatePost = async (data: CreatePostRequest, files?: File[]) => {
  let createdPostId: number | null = null;

  try {
    // 1. 投稿作成
    const newPost = await createPost.mutateAsync(data);
    createdPostId = newPost.id;

    // 2. 画像がある場合はアップロード
    if (files && files.length > 0) {
      try {
        await uploadMedia.mutateAsync({ postId: newPost.id, files });
      } catch (uploadError) {
        // 失敗時は投稿を削除（トランザクション的な動作）
        if (createdPostId) {
          await deletePost.mutateAsync(createdPostId);
        }
        throw new Error('画像のアップロードに失敗しました。投稿を取り消しました。');
      }
    } else {
      // 画像なしの場合は即座にタイムライン更新
      queryClient.invalidateQueries({ queryKey: ['timeline'] });
    }
  } catch (error: any) {
    throw new Error(errorMessage);
  }
};
```

- `frontend/src/hooks/usePosts.ts` (lines 32-37)
  - `useCreatePost` から自動タイムライン更新を削除
  - 呼び出し側で手動制御

- `frontend/src/hooks/useMedia.ts` (lines 15-17)
  - クエリキーを修正（`['posts']` → `['timeline']`, `['post']`, `['userPosts']`）

#### 画像アップロードフロー
```
1. ユーザーが画像付き投稿を作成
   ↓
2. 投稿をデータベースに作成（post_id取得）
   ↓
3. Firebase Storageに画像をアップロード
   ↓ 成功
4. 署名付きURLをデータベースに保存
   ↓
5. タイムラインを更新（画像付き投稿が表示される）

   ↓ 失敗
4. 作成した投稿を削除（ロールバック）
   ↓
5. エラーメッセージを表示
```

#### 署名付きURLの重要性
- **問題**: Firebase Storageのバケットが非公開（セキュリティ設定）
- **解決**: 署名付きURLを生成することで、一時的なアクセス権限を付与
- **有効期限**: 7日間
- **URL形式**: `https://storage.googleapis.com/bucket-name/uploads/file.jpg?GoogleAccessId=...&Expires=...&Signature=...`

---

## 🐛 解決した問題

### 1. Firebase Storage初期化失敗
**問題**: `{error: 'local storage upload not implemented in Phase 2'}`

**原因**: `media_service.go` が `config.LoadConfig()` を呼び出し、新しいconfig インスタンスを作成していた

**解決**: `config.AppConfig` を使用するように修正

**修正箇所**: `backend/internal/services/media_service.go:26`

---

### 2. メール未認証でもログイン可能
**問題**: 新規登録後、メール認証なしでログインできてしまう

**原因**: ログイン時に`email_verified`フィールドをチェックしていなかった

**解決**:
- `auth_service.go` にチェック追加
- `auth_handler.go` で403エラーレスポンス
- フロントエンドで認証待ちページへリダイレクト

**修正箇所**:
- `backend/internal/services/auth_service.go:80-83`
- `backend/internal/handlers/auth_handler.go:121-123`
- `frontend/src/components/auth/LoginForm.tsx:38-42`

---

### 3. ブックマーク状態が表示されない
**問題**: バックエンドが`is_bookmarked`フィールドを返していない

**原因**: `post_service.go` でブックマーク状態をクエリしていなかった

**解決**: タイムライン取得時にブックマーク状態を一括取得

**修正箇所**: `backend/internal/services/post_service.go:181-194, 222-225`

---

### 4. いいね/ブックマークの即座反映なし
**問題**: ボタンクリック後、サーバーレスポンスを待つまでUIが更新されない

**原因**: 楽観的更新を実装していなかった

**解決**: React Queryの`onMutate`, `onError`, `onSettled`を使用して楽観的更新を実装

**修正箇所**:
- `frontend/src/hooks/usePosts.ts:77-146`
- `frontend/src/hooks/useBookmarks.ts:12-82`

---

### 5. 画像アップロード失敗時に投稿が残る
**問題**: 画像アップロードに失敗しても、投稿だけが作成されてしまう

**原因**: エラーハンドリングとロールバック処理がなかった

**解決**: try-catchで画像アップロード失敗時に投稿を削除

**修正箇所**: `frontend/src/pages/HomePage.tsx:28-62`

---

### 6. 画像が投稿と別々に表示される
**問題**: 投稿が先に表示され、後から画像が追加される

**原因**: `useCreatePost` が即座にタイムラインを更新していた

**解決**:
- `useCreatePost` から自動タイムライン更新を削除
- `HomePage.tsx` で画像アップロード完了後に手動更新

**修正箇所**:
- `frontend/src/hooks/usePosts.ts:32-37`
- `frontend/src/pages/HomePage.tsx:28-62`

---

### 7. アップロード画像が403エラー
**問題**: 画像はアップロードされるが、表示時に403 Forbidden

**原因**: Firebase Storageバケットが非公開、公開URLではアクセスできない

**解決**: 署名付きURL（有効期限7日）を生成

**修正箇所**: `backend/internal/services/firebase_storage_service.go:92-102`

---

## 📊 Firebase Storage設定確認

### 環境変数（`.env`）
```bash
FIREBASE_STORAGE_BUCKET=udemy-sns-b9e40.firebasestorage.app
FIREBASE_CREDENTIALS_PATH=./service_account.json
```

### 設定確認コマンド
```bash
cd backend
bash scripts/test_image_upload.sh
```

### 確認結果
- ✅ `.env` ファイル存在
- ✅ `FIREBASE_STORAGE_BUCKET` 設定済み: `udemy-sns-b9e40.firebasestorage.app`
- ✅ `FIREBASE_CREDENTIALS_PATH` 設定済み: `./service_account.json`
- ✅ Firebase credentials ファイル存在

---

## 🧪 テスト手順

### 自動テストスクリプト
```bash
cd /Users/sugimoto/Desktop/udemy/sns/backend
bash scripts/test_image_upload.sh
```

### 手動テスト（推奨）

#### 1. 画像アップロードテスト
1. ブラウザで http://localhost:5173 を開く
2. メール認証済みユーザーでログイン
3. 投稿フォームに文字を入力
4. 画像ファイルを選択（JPG, PNG, GIF, HEIC）
5. 「投稿」ボタンをクリック

**期待結果:**
- ✅ 投稿が作成される
- ✅ 画像が投稿と一緒に表示される
- ✅ 画像をクリックすると正常に表示される（403エラーなし）
- ✅ 画像URLが署名付きURL形式（`...&Expires=...&Signature=...`）

#### 2. ブックマークテスト
1. タイムラインの投稿のブックマークアイコンをクリック
2. **期待結果:**
   - アイコンが即座に塗りつぶしに変化
   - ページリロード後もブックマーク状態が保持される
   - ブックマーク一覧ページで確認できる

#### 3. いいねテスト
1. タイムラインの投稿のいいねボタンをクリック
2. **期待結果:**
   - ハートアイコンが即座に赤くなる
   - いいね数が即座に+1される
   - ページリロード後もいいね状態が保持される

#### 4. メール認証テスト
1. 新規アカウントを作成
2. **期待結果:**
   - メール確認待ちページへ遷移
   - メールが送信される（Resend API）
3. 認証前にログインを試みる
4. **期待結果:**
   - 403エラー
   - メール確認待ちページへリダイレクト

---

## 📁 主要な変更ファイル一覧

### バックエンド
```
backend/
├── internal/
│   ├── services/
│   │   ├── firebase_storage_service.go  (署名付きURL生成)
│   │   ├── media_service.go             (config修正)
│   │   ├── auth_service.go              (メール認証チェック)
│   │   └── post_service.go              (ブックマーク状態取得)
│   └── handlers/
│       └── auth_handler.go              (403エラーハンドリング)
└── scripts/
    └── test_image_upload.sh             (テストスクリプト)
```

### フロントエンド
```
frontend/src/
├── pages/
│   ├── HomePage.tsx                     (画像アップロード同期)
│   └── EmailVerificationPendingPage.tsx (新規: メール確認待ち)
├── hooks/
│   ├── usePosts.ts                      (楽観的更新)
│   ├── useBookmarks.ts                  (楽観的更新)
│   └── useMedia.ts                      (クエリキー修正)
├── components/
│   └── auth/
│       ├── RegisterForm.tsx             (リダイレクト修正)
│       └── LoginForm.tsx                (403エラーハンドリング)
└── App.tsx                              (ルート追加)
```

### ドキュメント
```
docs/
├── PHASE2_IMPLEMENTATION_COMPLETE.md   (このファイル)
├── PHASE2_VERIFICATION_TEST.md          (検証テストガイド)
└── NON_FUNCTIONAL_REQUIREMENTS_DIAGNOSTIC.md
```

---

## 🎯 次のステップ

### Phase 2 残りの機能（オプション）
- [ ] ハッシュタグ機能（バックエンド実装済み）
- [ ] パスワードリセット機能（バックエンド実装済み、フロントエンド未実装）

### Phase 3 高度な機能
- [ ] ユーザー検索機能
- [ ] リツイート機能
- [ ] リアルタイム通知
- [ ] ダイレクトメッセージ
- [ ] トレンド投稿

---

## 📝 技術的なポイント

### React Queryの楽観的更新パターン
- `onMutate`: UI即座更新 + 現在のデータ保存
- `onError`: ロールバックで元のデータに戻す
- `onSettled`: サーバーと最終的に同期

### トランザクション的な処理（画像アップロード）
1. 投稿作成（ID取得）
2. 画像アップロード（try-catch）
3. 成功 → タイムライン更新
4. 失敗 → 投稿削除（ロールバック）

### Firebase Storage署名付きURL
- 非公開バケットでも安全にファイルを共有
- 有効期限を設定可能（7日間）
- URL漏洩時のセキュリティリスクを軽減

---

## ✅ Phase 2 実装完了チェックリスト

- [x] メール認証機能
- [x] ブックマーク機能
- [x] 楽観的UI更新
- [x] 画像アップロード（Firebase Storage）
- [x] 署名付きURL生成
- [x] エラーハンドリングとロールバック
- [x] すべての403エラー解決
- [x] テストスクリプト作成
- [x] ドキュメント作成

---

**実装完了日**: 2026-02-16
**担当**: Claude Code
**ステータス**: ✅ Phase 2 完了（画像アップロード動作テスト待ち）
