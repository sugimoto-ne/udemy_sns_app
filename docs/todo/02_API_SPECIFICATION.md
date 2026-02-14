# API仕様書

## 📡 基本情報

- **ベースURL（開発）**: `http://localhost:8080/api/v1`
- **ベースURL（本番）**: `https://your-api.example.com/api/v1`
- **プロトコル**: HTTPS（本番）/ HTTP（開発）
- **データ形式**: JSON
- **文字コード**: UTF-8
- **認証方式**: JWT (Bearer Token)

---

## 🔐 認証ヘッダー

認証が必要なエンドポイントは以下のヘッダーが必須：

```
Authorization: Bearer <JWT_TOKEN>
```

---

## 📋 共通レスポンス形式

### 成功レスポンス

```json
{
  "data": { /* レスポンスデータ */ },
  "message": "Success"
}
```

### エラーレスポンス

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message"
  }
}
```

### ページネーションレスポンス

```json
{
  "data": [ /* データ配列 */ ],
  "pagination": {
    "has_more": true,
    "next_cursor": "1234567890",
    "limit": 20
  }
}
```

---

## 🚀 エンドポイント一覧

### Phase 1 - MVP

#### 認証 (Authentication)
- `POST /auth/register` - ユーザー登録
- `POST /auth/login` - ログイン
- `POST /auth/logout` - ログアウト
- `GET /auth/me` - 現在のユーザー情報取得

#### ユーザー (Users)
- `GET /users/:username` - ユーザー情報取得
- `PUT /users/me` - プロフィール更新
- `GET /users/:username/posts` - ユーザーの投稿一覧
- `GET /users/:username/followers` - フォロワー一覧
- `GET /users/:username/following` - フォロー中一覧

#### 投稿 (Posts)
- `GET /posts` - タイムライン取得（全体 / フォロー中）
- `GET /posts/:id` - 投稿詳細取得
- `POST /posts` - 投稿作成
- `PUT /posts/:id` - 投稿更新
- `DELETE /posts/:id` - 投稿削除

#### コメント (Comments)
- `GET /posts/:id/comments` - コメント一覧取得
- `POST /posts/:id/comments` - コメント作成
- `DELETE /comments/:id` - コメント削除

#### いいね (Likes)
- `POST /posts/:id/like` - いいね追加
- `DELETE /posts/:id/like` - いいね削除
- `GET /posts/:id/likes` - いいね一覧取得

#### フォロー (Follows)
- `POST /users/:username/follow` - フォロー
- `DELETE /users/:username/follow` - フォロー解除

#### メディア (Media)
- `POST /media/upload` - メディアアップロード

### Phase 2

#### ハッシュタグ (Hashtags)
- `GET /hashtags/:name/posts` - ハッシュタグ別投稿一覧
- `GET /hashtags/trending` - トレンドハッシュタグ

#### ブックマーク (Bookmarks)
- `GET /bookmarks` - ブックマーク一覧
- `POST /posts/:id/bookmark` - ブックマーク追加
- `DELETE /posts/:id/bookmark` - ブックマーク削除

#### パスワードリセット
- `POST /auth/password-reset/request` - パスワードリセット要求
- `POST /auth/password-reset/confirm` - パスワードリセット実行

#### メール認証
- `POST /auth/email/verify` - メール認証実行
- `POST /auth/email/resend` - 確認メール再送信

### Phase 3

#### 通知 (Notifications)
- `GET /notifications` - 通知一覧
- `PUT /notifications/:id/read` - 通知既読
- `PUT /notifications/read-all` - 全通知既読

#### ユーザー検索
- `GET /users/search` - ユーザー検索

---

## 📖 詳細仕様

---

## 🔐 認証 (Authentication)

### POST /auth/register
ユーザー登録

**Phase**: 1
**認証**: 不要

**リクエストボディ**:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "username": "john_doe"
}
```

**バリデーション**:
- `email`: 必須、メール形式、最大255文字
- `password`: 必須、最小8文字
- `username`: 必須、英数字とアンダースコアのみ、3〜50文字

**レスポンス** (201 Created):
```json
{
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "username": "john_doe",
      "display_name": null,
      "avatar_url": null,
      "created_at": "2026-02-14T10:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "message": "User registered successfully"
}
```

**エラーレスポンス**:
- `400 Bad Request`: バリデーションエラー
- `409 Conflict`: メールアドレスまたはユーザー名が既に存在

---

### POST /auth/login
ログイン

**Phase**: 1
**認証**: 不要

**リクエストボディ**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**レスポンス** (200 OK):
```json
{
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "username": "john_doe",
      "display_name": "John Doe",
      "avatar_url": "https://example.com/avatar.jpg",
      "created_at": "2026-02-14T10:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "message": "Login successful"
}
```

**エラーレスポンス**:
- `400 Bad Request`: リクエストボディが不正
- `401 Unauthorized`: メールアドレスまたはパスワードが間違っている

---

### GET /auth/me
現在のユーザー情報取得

**Phase**: 1
**認証**: 必須

**レスポンス** (200 OK):
```json
{
  "data": {
    "id": 1,
    "email": "user@example.com",
    "username": "john_doe",
    "display_name": "John Doe",
    "bio": "Software Engineer",
    "avatar_url": "https://example.com/avatar.jpg",
    "header_url": "https://example.com/header.jpg",
    "website": "https://johndoe.com",
    "birth_date": "1990-01-01",
    "occupation": "Engineer",
    "followers_count": 150,
    "following_count": 100,
    "created_at": "2026-02-14T10:00:00Z"
  }
}
```

**エラーレスポンス**:
- `401 Unauthorized`: トークンが無効または期限切れ

---

## 👤 ユーザー (Users)

### GET /users/:username
ユーザー情報取得

**Phase**: 1
**認証**: 任意（認証済みの場合、フォロー状態を含む）

**パスメータ**:
- `username`: ユーザー名

**レスポンス** (200 OK):
```json
{
  "data": {
    "id": 1,
    "username": "john_doe",
    "display_name": "John Doe",
    "bio": "Software Engineer",
    "avatar_url": "https://example.com/avatar.jpg",
    "header_url": "https://example.com/header.jpg",
    "website": "https://johndoe.com",
    "birth_date": "1990-01-01",
    "occupation": "Engineer",
    "followers_count": 150,
    "following_count": 100,
    "is_following": true,
    "is_followed_by": false,
    "created_at": "2026-02-14T10:00:00Z"
  }
}
```

**エラーレスポンス**:
- `404 Not Found`: ユーザーが存在しない

---

### PUT /users/me
プロフィール更新

**Phase**: 1
**認証**: 必須

**リクエストボディ**:
```json
{
  "display_name": "John Doe",
  "bio": "Software Engineer & Tech Enthusiast",
  "website": "https://johndoe.com",
  "birth_date": "1990-01-01",
  "occupation": "Software Engineer",
  "avatar_url": "https://storage.example.com/avatar.jpg",
  "header_url": "https://storage.example.com/header.jpg"
}
```

**バリデーション**:
- 全フィールド任意
- `display_name`: 最大100文字
- `bio`: 最大500文字
- `website`: URL形式、最大255文字

**レスポンス** (200 OK):
```json
{
  "data": {
    "id": 1,
    "username": "john_doe",
    "display_name": "John Doe",
    "bio": "Software Engineer & Tech Enthusiast",
    "avatar_url": "https://storage.example.com/avatar.jpg",
    "header_url": "https://storage.example.com/header.jpg",
    "website": "https://johndoe.com",
    "birth_date": "1990-01-01",
    "occupation": "Software Engineer",
    "updated_at": "2026-02-14T11:00:00Z"
  },
  "message": "Profile updated successfully"
}
```

**エラーレスポンス**:
- `400 Bad Request`: バリデーションエラー
- `401 Unauthorized`: 認証エラー

---

### GET /users/:username/followers
フォロワー一覧取得

**Phase**: 1
**認証**: 任意

**クエリパラメータ**:
- `limit`: 取得件数（デフォルト: 20、最大: 100）
- `cursor`: ページネーションカーソル

**レスポンス** (200 OK):
```json
{
  "data": [
    {
      "id": 2,
      "username": "jane_smith",
      "display_name": "Jane Smith",
      "avatar_url": "https://example.com/avatar2.jpg",
      "bio": "Designer",
      "is_following": false,
      "followed_at": "2026-02-10T10:00:00Z"
    }
  ],
  "pagination": {
    "has_more": true,
    "next_cursor": "1234567890",
    "limit": 20
  }
}
```

---

### GET /users/:username/following
フォロー中一覧取得

**Phase**: 1
**認証**: 任意

**クエリパラメータ**:
- `limit`: 取得件数（デフォルト: 20、最大: 100）
- `cursor`: ページネーションカーソル

**レスポンス** (200 OK):
```json
{
  "data": [
    {
      "id": 3,
      "username": "bob_martin",
      "display_name": "Bob Martin",
      "avatar_url": "https://example.com/avatar3.jpg",
      "bio": "Writer",
      "is_following": true,
      "followed_at": "2026-02-12T10:00:00Z"
    }
  ],
  "pagination": {
    "has_more": false,
    "next_cursor": null,
    "limit": 20
  }
}
```

---

## 📝 投稿 (Posts)

### GET /posts
タイムライン取得

**Phase**: 1
**認証**: 必須（フォロー中タイムライン）/ 任意（全体タイムライン）

**クエリパラメータ**:
- `type`: タイムライン種別（`following` または `all`、デフォルト: `all`）
- `limit`: 取得件数（デフォルト: 20、最大: 100）
- `cursor`: ページネーションカーソル（投稿IDベース）

**レスポンス** (200 OK):
```json
{
  "data": [
    {
      "id": 1,
      "user": {
        "id": 1,
        "username": "john_doe",
        "display_name": "John Doe",
        "avatar_url": "https://example.com/avatar.jpg"
      },
      "content": "Hello World! This is my first post.",
      "media": [
        {
          "id": 1,
          "media_type": "image",
          "media_url": "https://storage.example.com/image.jpg",
          "file_size": 1024000,
          "order_index": 0
        }
      ],
      "likes_count": 10,
      "comments_count": 5,
      "is_liked": true,
      "is_bookmarked": false,
      "created_at": "2026-02-14T10:00:00Z",
      "updated_at": "2026-02-14T10:00:00Z"
    }
  ],
  "pagination": {
    "has_more": true,
    "next_cursor": "1234567890",
    "limit": 20
  }
}
```

---

### GET /posts/:id
投稿詳細取得

**Phase**: 1
**認証**: 任意

**パスパラメータ**:
- `id`: 投稿ID

**レスポンス** (200 OK):
```json
{
  "data": {
    "id": 1,
    "user": {
      "id": 1,
      "username": "john_doe",
      "display_name": "John Doe",
      "avatar_url": "https://example.com/avatar.jpg"
    },
    "content": "Hello World! This is my first post.",
    "media": [
      {
        "id": 1,
        "media_type": "image",
        "media_url": "https://storage.example.com/image.jpg",
        "file_size": 1024000,
        "order_index": 0
      }
    ],
    "likes_count": 10,
    "comments_count": 5,
    "is_liked": true,
    "is_bookmarked": false,
    "created_at": "2026-02-14T10:00:00Z",
    "updated_at": "2026-02-14T10:00:00Z"
  }
}
```

**エラーレスポンス**:
- `404 Not Found`: 投稿が存在しない

---

### POST /posts
投稿作成

**Phase**: 1
**認証**: 必須

**リクエストボディ**:
```json
{
  "content": "Hello World! This is my first post.",
  "media_urls": [
    "https://storage.example.com/image1.jpg",
    "https://storage.example.com/image2.jpg"
  ]
}
```

**バリデーション**:
- `content`: 必須、最大280文字
- `media_urls`: 任意、配列、最大4件（Phase 2）

**レスポンス** (201 Created):
```json
{
  "data": {
    "id": 1,
    "user": {
      "id": 1,
      "username": "john_doe",
      "display_name": "John Doe",
      "avatar_url": "https://example.com/avatar.jpg"
    },
    "content": "Hello World! This is my first post.",
    "media": [
      {
        "id": 1,
        "media_type": "image",
        "media_url": "https://storage.example.com/image1.jpg",
        "file_size": 1024000,
        "order_index": 0
      }
    ],
    "likes_count": 0,
    "comments_count": 0,
    "created_at": "2026-02-14T10:00:00Z"
  },
  "message": "Post created successfully"
}
```

**エラーレスポンス**:
- `400 Bad Request`: バリデーションエラー
- `401 Unauthorized`: 認証エラー

---

### PUT /posts/:id
投稿更新

**Phase**: 1
**認証**: 必須（自分の投稿のみ）

**パスパラメータ**:
- `id`: 投稿ID

**リクエストボディ**:
```json
{
  "content": "Updated content"
}
```

**レスポンス** (200 OK):
```json
{
  "data": {
    "id": 1,
    "content": "Updated content",
    "updated_at": "2026-02-14T11:00:00Z"
  },
  "message": "Post updated successfully"
}
```

**エラーレスポンス**:
- `400 Bad Request`: バリデーションエラー
- `401 Unauthorized`: 認証エラー
- `403 Forbidden`: 他人の投稿を編集しようとした
- `404 Not Found`: 投稿が存在しない

---

### DELETE /posts/:id
投稿削除（論理削除）

**Phase**: 1
**認証**: 必須（自分の投稿のみ）

**パスパラメータ**:
- `id`: 投稿ID

**レスポンス** (200 OK):
```json
{
  "message": "Post deleted successfully"
}
```

**エラーレスポンス**:
- `401 Unauthorized`: 認証エラー
- `403 Forbidden`: 他人の投稿を削除しようとした
- `404 Not Found`: 投稿が存在しない

---

## 💬 コメント (Comments)

### GET /posts/:id/comments
コメント一覧取得

**Phase**: 1
**認証**: 任意

**パスパラメータ**:
- `id`: 投稿ID

**クエリパラメータ**:
- `limit`: 取得件数（デフォルト: 20、最大: 100）
- `cursor`: ページネーションカーソル

**レスポンス** (200 OK):
```json
{
  "data": [
    {
      "id": 1,
      "user": {
        "id": 2,
        "username": "jane_smith",
        "display_name": "Jane Smith",
        "avatar_url": "https://example.com/avatar2.jpg"
      },
      "content": "Great post!",
      "created_at": "2026-02-14T10:30:00Z"
    }
  ],
  "pagination": {
    "has_more": false,
    "next_cursor": null,
    "limit": 20
  }
}
```

---

### POST /posts/:id/comments
コメント作成

**Phase**: 1
**認証**: 必須

**パスパラメータ**:
- `id`: 投稿ID

**リクエストボディ**:
```json
{
  "content": "Great post!"
}
```

**バリデーション**:
- `content`: 必須、最大280文字

**レスポンス** (201 Created):
```json
{
  "data": {
    "id": 1,
    "user": {
      "id": 2,
      "username": "jane_smith",
      "display_name": "Jane Smith",
      "avatar_url": "https://example.com/avatar2.jpg"
    },
    "post_id": 1,
    "content": "Great post!",
    "created_at": "2026-02-14T10:30:00Z"
  },
  "message": "Comment created successfully"
}
```

**エラーレスポンス**:
- `400 Bad Request`: バリデーションエラー
- `401 Unauthorized`: 認証エラー
- `404 Not Found`: 投稿が存在しない

---

### DELETE /comments/:id
コメント削除

**Phase**: 1
**認証**: 必須（自分のコメントのみ）

**パスパラメータ**:
- `id`: コメントID

**レスポンス** (200 OK):
```json
{
  "message": "Comment deleted successfully"
}
```

**エラーレスポンス**:
- `401 Unauthorized`: 認証エラー
- `403 Forbidden`: 他人のコメントを削除しようとした
- `404 Not Found`: コメントが存在しない

---

## ❤️ いいね (Likes)

### POST /posts/:id/like
いいね追加

**Phase**: 1
**認証**: 必須

**パスパラメータ**:
- `id`: 投稿ID

**レスポンス** (201 Created):
```json
{
  "data": {
    "post_id": 1,
    "likes_count": 11
  },
  "message": "Post liked successfully"
}
```

**エラーレスポンス**:
- `401 Unauthorized`: 認証エラー
- `404 Not Found`: 投稿が存在しない
- `409 Conflict`: 既にいいね済み

---

### DELETE /posts/:id/like
いいね削除

**Phase**: 1
**認証**: 必須

**パスパラメータ**:
- `id`: 投稿ID

**レスポンス** (200 OK):
```json
{
  "data": {
    "post_id": 1,
    "likes_count": 10
  },
  "message": "Post unliked successfully"
}
```

**エラーレスポンス**:
- `401 Unauthorized`: 認証エラー
- `404 Not Found`: 投稿またはいいねが存在しない

---

### GET /posts/:id/likes
いいね一覧取得

**Phase**: 1
**認証**: 任意

**パスパラメータ**:
- `id`: 投稿ID

**クエリパラメータ**:
- `limit`: 取得件数（デフォルト: 20、最大: 100）
- `cursor`: ページネーションカーソル

**レスポンス** (200 OK):
```json
{
  "data": [
    {
      "id": 2,
      "username": "jane_smith",
      "display_name": "Jane Smith",
      "avatar_url": "https://example.com/avatar2.jpg",
      "liked_at": "2026-02-14T10:15:00Z"
    }
  ],
  "pagination": {
    "has_more": false,
    "next_cursor": null,
    "limit": 20
  }
}
```

---

## 👥 フォロー (Follows)

### POST /users/:username/follow
フォロー

**Phase**: 1
**認証**: 必須

**パスパラメータ**:
- `username`: フォローするユーザー名

**レスポンス** (201 Created):
```json
{
  "data": {
    "username": "jane_smith",
    "is_following": true,
    "followed_at": "2026-02-14T10:00:00Z"
  },
  "message": "User followed successfully"
}
```

**エラーレスポンス**:
- `400 Bad Request`: 自分自身をフォローしようとした
- `401 Unauthorized`: 認証エラー
- `404 Not Found`: ユーザーが存在しない
- `409 Conflict`: 既にフォロー済み

---

### DELETE /users/:username/follow
フォロー解除

**Phase**: 1
**認証**: 必須

**パスパラメータ**:
- `username`: フォロー解除するユーザー名

**レスポンス** (200 OK):
```json
{
  "data": {
    "username": "jane_smith",
    "is_following": false
  },
  "message": "User unfollowed successfully"
}
```

**エラーレスポンス**:
- `401 Unauthorized`: 認証エラー
- `404 Not Found`: ユーザーまたはフォロー関係が存在しない

---

## 📷 メディア (Media)

### POST /media/upload
メディアアップロード

**Phase**: 1
**認証**: 必須

**リクエスト形式**: `multipart/form-data`

**フォームデータ**:
- `file`: メディアファイル（画像/動画/音声）

**レスポンス** (201 Created):
```json
{
  "data": {
    "media_url": "https://storage.example.com/uploads/abc123.jpg",
    "media_type": "image",
    "file_size": 1024000
  },
  "message": "Media uploaded successfully"
}
```

**エラーレスポンス**:
- `400 Bad Request`: ファイルが添付されていない、サイズ超過、非対応形式
- `401 Unauthorized`: 認証エラー

**バリデーション**:
- 画像: 最大5MB、jpg/png/gif/heic
- 動画: 最大50MB、30秒以内、mp4/mov
- 音声: mp3

---

## 🔖 ブックマーク (Bookmarks) - Phase 2

### GET /bookmarks
ブックマーク一覧取得

**Phase**: 2
**認証**: 必須

**クエリパラメータ**:
- `limit`: 取得件数（デフォルト: 20、最大: 100）
- `cursor`: ページネーションカーソル

**レスポンス** (200 OK):
```json
{
  "data": [
    {
      "id": 1,
      "user": {
        "id": 1,
        "username": "john_doe",
        "display_name": "John Doe",
        "avatar_url": "https://example.com/avatar.jpg"
      },
      "content": "Bookmarked post content",
      "media": [],
      "likes_count": 50,
      "comments_count": 10,
      "bookmarked_at": "2026-02-14T12:00:00Z",
      "created_at": "2026-02-14T10:00:00Z"
    }
  ],
  "pagination": {
    "has_more": true,
    "next_cursor": "1234567890",
    "limit": 20
  }
}
```

---

### POST /posts/:id/bookmark
ブックマーク追加

**Phase**: 2
**認証**: 必須

**パスパラメータ**:
- `id`: 投稿ID

**レスポンス** (201 Created):
```json
{
  "data": {
    "post_id": 1,
    "is_bookmarked": true,
    "bookmarked_at": "2026-02-14T12:00:00Z"
  },
  "message": "Post bookmarked successfully"
}
```

---

### DELETE /posts/:id/bookmark
ブックマーク削除

**Phase**: 2
**認証**: 必須

**パスパラメータ**:
- `id`: 投稿ID

**レスポンス** (200 OK):
```json
{
  "data": {
    "post_id": 1,
    "is_bookmarked": false
  },
  "message": "Bookmark removed successfully"
}
```

---

## 📊 エラーコード一覧

| HTTPステータス | エラーコード | 説明 |
|--------------|-------------|------|
| 400 | `VALIDATION_ERROR` | バリデーションエラー |
| 400 | `INVALID_REQUEST` | リクエストが不正 |
| 401 | `UNAUTHORIZED` | 認証エラー |
| 401 | `INVALID_TOKEN` | トークンが無効 |
| 401 | `TOKEN_EXPIRED` | トークンの期限切れ |
| 403 | `FORBIDDEN` | アクセス権限なし |
| 404 | `NOT_FOUND` | リソースが存在しない |
| 409 | `CONFLICT` | リソースの競合（既に存在など） |
| 413 | `FILE_TOO_LARGE` | ファイルサイズ超過 |
| 415 | `UNSUPPORTED_MEDIA_TYPE` | 非対応のファイル形式 |
| 429 | `RATE_LIMIT_EXCEEDED` | レート制限超過 |
| 500 | `INTERNAL_SERVER_ERROR` | サーバーエラー |

---

## 🔄 ページネーション仕様

**カーソルベースページネーション**を採用

### リクエスト
```
GET /posts?limit=20&cursor=1234567890
```

### レスポンス
```json
{
  "data": [ /* データ配列 */ ],
  "pagination": {
    "has_more": true,
    "next_cursor": "9876543210",
    "limit": 20
  }
}
```

### 次のページ取得
```
GET /posts?limit=20&cursor=9876543210
```

---

## 🛡️ レート制限

| エンドポイント | 制限 |
|--------------|------|
| 認証（ログイン/登録） | 10回/分 |
| 投稿作成 | 30回/時間 |
| いいね/フォロー | 100回/時間 |
| その他（読み取り） | 300回/15分 |

**レート制限超過時のレスポンス** (429 Too Many Requests):
```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Rate limit exceeded. Please try again later.",
    "retry_after": 60
  }
}
```

---

**作成日**: 2026-02-14
**最終更新**: 2026-02-14
