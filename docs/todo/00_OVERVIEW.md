# SNSアプリケーション - プロジェクト概要

## 📋 プロジェクト概要

ユーザーがテキスト、画像、動画、音声を投稿し、いいねやコメントで交流できるTwitterライクなSNSアプリケーション。

---

## 🎯 主な機能

### Phase 1（高優先度 - MVP）
- ✅ ユーザー認証（登録・ログイン・JWT）
- ✅ プロフィール表示・編集
- ✅ フォロー/フォロワー機能
- ✅ 投稿機能（テキスト、画像、動画、音声）
- ✅ 投稿の編集・削除（論理削除）
- ✅ コメント（リプライ）機能
- ✅ タイムライン表示（全体 / フォロー中の切り替え）
- ✅ 無限スクロール（20件ずつ取得）
- ✅ いいね機能（追加・取り消し・ユーザー一覧表示）

### Phase 2（中優先度）
- 🔶 ハッシュタグ機能
- 🔶 複数画像添付機能（1投稿に複数枚）
- 🔶 ブックマーク機能
- 🔶 パスワードリセット機能
- 🔶 メールアドレス認証機能

### Phase 3（低優先度 - 将来的な機能）
- 🔹 ユーザー検索機能
- 🔹 リツイート（再投稿/シェア）機能
- 🔹 通知機能（いいね・コメント通知）
- 🔹 ダイレクトメッセージ（DM）機能
- 🔹 トレンド/人気投稿表示
- 🔹 ソーシャルログイン（Google、Twitter等）

---

## 🛠️ 技術スタック

### フロントエンド
- **言語**: TypeScript
- **フレームワーク**: React
- **UIライブラリ**: Material-UI (MUI)
- **状態管理**: React Query / Context API
- **ルーティング**: React Router
- **型生成**: openapi-typescript (OpenAPI定義から自動生成)
- **APIクライアント**: openapi-fetch (型安全なHTTPクライアント)
- **レスポンシブ対応**: 必須

### バックエンド
- **言語**: Go
- **フレームワーク**: Echo
- **ORM**: GORM
- **認証**: JWT (JSON Web Token)
- **バリデーション**: validator/v10
- **API仕様生成**: swaggo/echo-swagger (OpenAPI/Swagger定義の自動生成)
- **ドキュメント**: Swagger UI

### データベース
- **RDBMS**: PostgreSQL
- **削除方式**: 論理削除（deleted_atカラム使用）

### インフラ・デプロイ
- **ローカル開発**: Docker + Docker Compose
- **バックエンドデプロイ**: Render / Google Cloud Run
- **フロントエンドデプロイ**: Firebase Hosting
- **メディアストレージ**:
  - ローカル: Docker Volume
  - 本番: Firebase Storage

---

## 📐 アーキテクチャ

```
┌─────────────────────────────────────────────────┐
│           Frontend (React + TypeScript)         │
│              Firebase Hosting                   │
└────────────────┬────────────────────────────────┘
                 │ REST API (JWT)
┌────────────────▼────────────────────────────────┐
│         Backend (Go + Echo + GORM)              │
│              Render / Cloud Run                 │
└────────┬───────────────────────┬─────────────────┘
         │                       │
┌────────▼──────────┐   ┌────────▼─────────────────┐
│   PostgreSQL      │   │   Firebase Storage       │
│   (Database)      │   │   (Images/Videos/Audio)  │
└───────────────────┘   └──────────────────────────┘
```

---

## 📊 データモデル概要

### 主要エンティティ
1. **Users** - ユーザー情報
2. **Posts** - 投稿
3. **Comments** - コメント（投稿へのリプライ）
4. **Likes** - いいね
5. **Follows** - フォロー関係
6. **Media** - メディアファイル（画像・動画・音声）
7. **Hashtags** - ハッシュタグ（Phase 2）
8. **Bookmarks** - ブックマーク（Phase 2）
9. **Notifications** - 通知（Phase 3）

詳細は `01_DATABASE_SCHEMA.md` を参照

---

## 🔐 認証・セキュリティ

### ユーザー登録（サインアップ）
- メールアドレス（一意）
- パスワード（bcryptでハッシュ化）
- ユーザー名（一意、@username形式）

### JWT認証フロー
```
1. POST /api/auth/register (登録)
2. POST /api/auth/login (ログイン) → JWT発行
3. Authorization: Bearer <token> でAPI呼び出し
4. バックエンドでJWT検証 → ユーザー情報取得
```

### メール認証（Phase 2）
- 登録時に確認メール送信
- メール内のリンクをクリックして認証完了

---

## 📁 メディアファイル制限

| メディア種別 | 最大サイズ | 最大長さ | 対応フォーマット |
|------------|----------|---------|----------------|
| 画像       | 5 MB     | -       | jpg, png, gif, heic |
| 動画       | 50 MB    | 30秒    | mp4, mov       |
| 音声       | -        | -       | mp3            |

---

## 🌐 API設計原則

- RESTful API
- エンドポイント形式: `/api/v1/resource`
- レスポンス形式: JSON
- ページネーション: カーソルベース（無限スクロール対応）
- エラーハンドリング: 統一されたエラーレスポンス
- **OpenAPI 3.0準拠**: バックエンドでSwagger自動生成 → フロントエンドで型自動生成
- **型安全性**: TypeScriptとGoの型システムをOpenAPIで連携

### 型生成フロー
```
Backend (Go)
  ↓ swaggo annotations
swagger.yaml (OpenAPI 3.0)
  ↓ openapi-typescript
Frontend (TypeScript)
  ↓ 型安全なAPI呼び出し
```

詳細は `02_API_SPECIFICATION.md` を参照

---

## 📱 画面構成（フロントエンド）

### 認証画面
- `/login` - ログイン
- `/register` - ユーザー登録

### メイン画面
- `/` - ホームタイムライン（フォロー中 / 全体切り替え）
- `/post/:id` - 投稿詳細・コメント表示
- `/profile/:username` - ユーザープロフィール
- `/profile/edit` - プロフィール編集
- `/bookmarks` - ブックマーク一覧（Phase 2）
- `/notifications` - 通知一覧（Phase 3）

---

## 🚀 開発フェーズ

### Phase 1（MVP - 基本機能）
**目標**: 投稿・いいね・コメント・フォローができる最小限のSNS

**バックエンド**
- 認証API（登録・ログイン）
- ユーザープロフィールAPI
- 投稿CRUD API
- コメントAPI
- いいねAPI
- フォローAPI
- タイムラインAPI

**フロントエンド**
- 認証画面
- タイムライン表示（無限スクロール）
- 投稿作成・編集・削除
- コメント表示・投稿
- いいね機能
- プロフィール表示・編集
- フォロー機能

### Phase 2（拡張機能）
**目標**: ユーザー体験の向上

- ハッシュタグ機能
- 複数画像添付
- ブックマーク機能
- パスワードリセット
- メール認証

### Phase 3（高度な機能）
**目標**: エンゲージメント向上

- ユーザー検索
- リツイート機能
- リアルタイム通知
- ダイレクトメッセージ
- トレンド表示

---

## 🔄 データ削除ポリシー

**論理削除（Soft Delete）**を採用

- `deleted_at` カラムを使用
- 削除されたデータはNULLではなく削除日時を記録
- クエリ時は `deleted_at IS NULL` で有効なデータのみ取得
- 物理削除は管理者のみが実行可能

**対象テーブル**
- Users
- Posts
- Comments

---

## 📝 命名規則

### データベース
- テーブル名: スネークケース複数形（`users`, `posts`, `post_likes`）
- カラム名: スネークケース（`user_id`, `created_at`）

### バックエンド（Go）
- 構造体: パスカルケース（`User`, `Post`）
- 関数: キャメルケース（`GetUserByID`, `CreatePost`）
- パッケージ: 小文字（`auth`, `models`）

### フロントエンド（TypeScript）
- コンポーネント: パスカルケース（`TimelinePost`, `UserProfile`）
- 関数: キャメルケース（`fetchPosts`, `handleLike`）
- 定数: スネークケース大文字（`API_BASE_URL`）

---

## 🧪 テスト方針

### バックエンド
- ユニットテスト: Go標準のtestingパッケージ
- APIテスト: HTTPテスト

### フロントエンド
- コンポーネントテスト: React Testing Library
- E2Eテスト: Playwright（Phase 2以降）

---

## 📚 参考ドキュメント

- [データベーススキーマ](./01_DATABASE_SCHEMA.md)
- [API仕様](./02_API_SPECIFICATION.md)
- [Phase 1 バックエンドTODO](./03_PHASE1_BACKEND.md)
- [Phase 1 フロントエンドTODO](./04_PHASE1_FRONTEND.md)
- [Phase 2 バックエンドTODO](./05_PHASE2_BACKEND.md)
- [Phase 2 フロントエンドTODO](./06_PHASE2_FRONTEND.md)
- [Phase 3 将来的な機能](./07_PHASE3_FUTURE.md)

---

## 🎯 成功基準

### Phase 1完了条件
- [ ] ユーザー登録・ログインができる
- [ ] プロフィールを編集できる
- [ ] 投稿（テキスト・画像・動画・音声）を作成できる
- [ ] 投稿を編集・削除できる
- [ ] 投稿にコメントできる
- [ ] 投稿にいいねできる
- [ ] ユーザーをフォロー/フォロー解除できる
- [ ] タイムラインを表示できる（全体 / フォロー中）
- [ ] 無限スクロールが機能する
- [ ] レスポンシブデザインが実装されている
- [ ] Docker環境で動作する

---

**作成日**: 2026-02-14
**最終更新**: 2026-02-14
