# データベーススキーマ設計

## 📊 ER図（Entity Relationship Diagram）

```
┌─────────────────┐
│     Users       │
├─────────────────┤
│ id (PK)         │
│ email           │──┐
│ password        │  │
│ username        │  │    ┌──────────────────┐
│ display_name    │  │    │   Follows        │
│ bio             │  ├───<│ follower_id (FK) │
│ avatar_url      │  │    │ following_id(FK) │
│ header_url      │  │    │ created_at       │
│ website         │  │    └──────────────────┘
│ birth_date      │  │
│ occupation      │  │
│ email_verified  │  │
│ created_at      │  │
│ updated_at      │  │
│ deleted_at      │  │
└─────────────────┘  │
         │           │
         │ (1)       │
         │           │
         │ (*)       │
         ▼           │
┌─────────────────┐  │    ┌─────────────────┐
│     Posts       │  │    │   Post_Likes    │
├─────────────────┤  │    ├─────────────────┤
│ id (PK)         │  │    │ id (PK)         │
│ user_id (FK)    │──┘    │ user_id (FK)    │──┐
│ content         │       │ post_id (FK)    │──┤
│ created_at      │──┐    │ created_at      │  │
│ updated_at      │  │    └─────────────────┘  │
│ deleted_at      │  │                          │
└─────────────────┘  │                          │
         │           │                          │
         │ (1)       │                          │
         │           │    ┌─────────────────┐  │
         │ (*)       │    │   Comments      │  │
         ▼           │    ├─────────────────┤  │
┌─────────────────┐  │    │ id (PK)         │  │
│     Media       │  │    │ post_id (FK)    │──┤
├─────────────────┤  │    │ user_id (FK)    │──┤
│ id (PK)         │  │    │ content         │  │
│ post_id (FK)    │──┘    │ created_at      │  │
│ media_type      │       │ updated_at      │  │
│ media_url       │       │ deleted_at      │  │
│ file_size       │       └─────────────────┘  │
│ duration        │                             │
│ order_index     │       ┌─────────────────┐  │
│ created_at      │       │   Bookmarks     │  │
└─────────────────┘       ├─────────────────┤  │
                          │ id (PK)         │  │
         ┌────────────────│ user_id (FK)    │──┘
         │                │ post_id (FK)    │──┐
         │                │ created_at      │  │
         │                └─────────────────┘  │
         │                                     │
         │                ┌─────────────────┐  │
         │                │   Hashtags      │  │
         │                ├─────────────────┤  │
         │                │ id (PK)         │  │
         │                │ name            │  │
         │                │ created_at      │  │
         │                └─────────────────┘  │
         │                         │           │
         │                         │ (*)       │
         │                         │           │
         │                         │ (*)       │
         │                         ▼           │
         │                ┌─────────────────┐  │
         │                │ Post_Hashtags   │  │
         │                ├─────────────────┤  │
         └───────────────>│ post_id (FK)    │  │
                          │ hashtag_id (FK) │  │
                          │ created_at      │  │
                          └─────────────────┘  │
                                                │
                          ┌─────────────────┐  │
                          │ Notifications   │  │
                          ├─────────────────┤  │
                          │ id (PK)         │  │
                          │ user_id (FK)    │──┘
                          │ actor_id (FK)   │
                          │ type            │
                          │ post_id (FK)    │
                          │ comment_id (FK) │
                          │ is_read         │
                          │ created_at      │
                          └─────────────────┘
```

---

## 📋 テーブル定義

### 1. users（ユーザー）

**Phase**: 1

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|----------|------|
| id | BIGSERIAL | NOT NULL | AUTO | ユーザーID（主キー） |
| email | VARCHAR(255) | NOT NULL | - | メールアドレス（一意） |
| password | VARCHAR(255) | NOT NULL | - | パスワード（bcryptハッシュ） |
| username | VARCHAR(50) | NOT NULL | - | ユーザー名（@username、一意） |
| display_name | VARCHAR(100) | NULL | NULL | 表示名（ニックネーム） |
| bio | TEXT | NULL | NULL | 自己紹介 |
| avatar_url | VARCHAR(500) | NULL | NULL | アイコン画像URL |
| header_url | VARCHAR(500) | NULL | NULL | ヘッダー画像URL |
| website | VARCHAR(255) | NULL | NULL | ウェブサイトURL |
| birth_date | DATE | NULL | NULL | 誕生日 |
| occupation | VARCHAR(100) | NULL | NULL | 職業 |
| email_verified | BOOLEAN | NOT NULL | FALSE | メール認証済みフラグ（Phase 2） |
| created_at | TIMESTAMP | NOT NULL | NOW() | 作成日時 |
| updated_at | TIMESTAMP | NOT NULL | NOW() | 更新日時 |
| deleted_at | TIMESTAMP | NULL | NULL | 削除日時（論理削除） |

**インデックス**:
- PRIMARY KEY (id)
- UNIQUE (email) WHERE deleted_at IS NULL
- UNIQUE (username) WHERE deleted_at IS NULL
- INDEX (created_at)

**制約**:
- email: メール形式バリデーション
- username: 英数字とアンダースコアのみ、3〜50文字
- password: 最小8文字（アプリケーション層でバリデーション）

---

### 2. posts（投稿）

**Phase**: 1

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|----------|------|
| id | BIGSERIAL | NOT NULL | AUTO | 投稿ID（主キー） |
| user_id | BIGINT | NOT NULL | - | 投稿者ID（外部キー → users.id） |
| content | TEXT | NOT NULL | - | 投稿内容（最大280文字） |
| created_at | TIMESTAMP | NOT NULL | NOW() | 作成日時 |
| updated_at | TIMESTAMP | NOT NULL | NOW() | 更新日時 |
| deleted_at | TIMESTAMP | NULL | NULL | 削除日時（論理削除） |

**インデックス**:
- PRIMARY KEY (id)
- INDEX (user_id)
- INDEX (created_at DESC)
- INDEX (deleted_at)

**外部キー制約**:
- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE

---

### 3. media（メディアファイル）

**Phase**: 1

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|----------|------|
| id | BIGSERIAL | NOT NULL | AUTO | メディアID（主キー） |
| post_id | BIGINT | NOT NULL | - | 投稿ID（外部キー → posts.id） |
| media_type | VARCHAR(20) | NOT NULL | - | メディア種別（image/video/audio） |
| media_url | VARCHAR(500) | NOT NULL | - | メディアファイルURL |
| file_size | BIGINT | NOT NULL | - | ファイルサイズ（バイト） |
| duration | INTEGER | NULL | NULL | 動画・音声の長さ（秒） |
| order_index | INTEGER | NOT NULL | 0 | 表示順序（Phase 2の複数画像対応） |
| created_at | TIMESTAMP | NOT NULL | NOW() | 作成日時 |

**インデックス**:
- PRIMARY KEY (id)
- INDEX (post_id, order_index)

**外部キー制約**:
- FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE

**制約**:
- media_type: 'image', 'video', 'audio' のいずれか
- file_size: 画像5MB、動画50MB以内（アプリケーション層でバリデーション）

---

### 4. comments（コメント）

**Phase**: 1

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|----------|------|
| id | BIGSERIAL | NOT NULL | AUTO | コメントID（主キー） |
| post_id | BIGINT | NOT NULL | - | 投稿ID（外部キー → posts.id） |
| user_id | BIGINT | NOT NULL | - | コメント投稿者ID（外部キー → users.id） |
| content | TEXT | NOT NULL | - | コメント内容（最大280文字） |
| created_at | TIMESTAMP | NOT NULL | NOW() | 作成日時 |
| updated_at | TIMESTAMP | NOT NULL | NOW() | 更新日時 |
| deleted_at | TIMESTAMP | NULL | NULL | 削除日時（論理削除） |

**インデックス**:
- PRIMARY KEY (id)
- INDEX (post_id, created_at)
- INDEX (user_id)
- INDEX (deleted_at)

**外部キー制約**:
- FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE

---

### 5. post_likes（投稿いいね）

**Phase**: 1

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|----------|------|
| id | BIGSERIAL | NOT NULL | AUTO | いいねID（主キー） |
| post_id | BIGINT | NOT NULL | - | 投稿ID（外部キー → posts.id） |
| user_id | BIGINT | NOT NULL | - | いいねしたユーザーID（外部キー → users.id） |
| created_at | TIMESTAMP | NOT NULL | NOW() | 作成日時 |

**インデックス**:
- PRIMARY KEY (id)
- UNIQUE (post_id, user_id)
- INDEX (post_id, created_at)
- INDEX (user_id)

**外部キー制約**:
- FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE

**制約**:
- 同じユーザーが同じ投稿に複数回いいねできない（UNIQUE制約）

---

### 6. follows（フォロー関係）

**Phase**: 1

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|----------|------|
| id | BIGSERIAL | NOT NULL | AUTO | フォローID（主キー） |
| follower_id | BIGINT | NOT NULL | - | フォローする側のユーザーID（外部キー → users.id） |
| following_id | BIGINT | NOT NULL | - | フォローされる側のユーザーID（外部キー → users.id） |
| created_at | TIMESTAMP | NOT NULL | NOW() | 作成日時 |

**インデックス**:
- PRIMARY KEY (id)
- UNIQUE (follower_id, following_id)
- INDEX (follower_id)
- INDEX (following_id)

**外部キー制約**:
- FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE
- FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE

**制約**:
- 自分自身をフォローできない（アプリケーション層でバリデーション）
- 同じユーザーを複数回フォローできない（UNIQUE制約）

---

### 7. hashtags（ハッシュタグ）

**Phase**: 2

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|----------|------|
| id | BIGSERIAL | NOT NULL | AUTO | ハッシュタグID（主キー） |
| name | VARCHAR(100) | NOT NULL | - | ハッシュタグ名（#を除く） |
| created_at | TIMESTAMP | NOT NULL | NOW() | 作成日時 |

**インデックス**:
- PRIMARY KEY (id)
- UNIQUE (name)
- INDEX (created_at DESC)

**制約**:
- name: 英数字、日本語、アンダースコアのみ

---

### 8. post_hashtags（投稿とハッシュタグの中間テーブル）

**Phase**: 2

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|----------|------|
| id | BIGSERIAL | NOT NULL | AUTO | ID（主キー） |
| post_id | BIGINT | NOT NULL | - | 投稿ID（外部キー → posts.id） |
| hashtag_id | BIGINT | NOT NULL | - | ハッシュタグID（外部キー → hashtags.id） |
| created_at | TIMESTAMP | NOT NULL | NOW() | 作成日時 |

**インデックス**:
- PRIMARY KEY (id)
- UNIQUE (post_id, hashtag_id)
- INDEX (hashtag_id, created_at DESC)

**外部キー制約**:
- FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
- FOREIGN KEY (hashtag_id) REFERENCES hashtags(id) ON DELETE CASCADE

---

### 9. bookmarks（ブックマーク）

**Phase**: 2

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|----------|------|
| id | BIGSERIAL | NOT NULL | AUTO | ブックマークID（主キー） |
| user_id | BIGINT | NOT NULL | - | ユーザーID（外部キー → users.id） |
| post_id | BIGINT | NOT NULL | - | 投稿ID（外部キー → posts.id） |
| created_at | TIMESTAMP | NOT NULL | NOW() | 作成日時 |

**インデックス**:
- PRIMARY KEY (id)
- UNIQUE (user_id, post_id)
- INDEX (user_id, created_at DESC)

**外部キー制約**:
- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
- FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE

---

### 10. notifications（通知）

**Phase**: 3

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|----------|------|
| id | BIGSERIAL | NOT NULL | AUTO | 通知ID（主キー） |
| user_id | BIGINT | NOT NULL | - | 通知を受け取るユーザーID（外部キー → users.id） |
| actor_id | BIGINT | NOT NULL | - | 通知を発生させたユーザーID（外部キー → users.id） |
| type | VARCHAR(20) | NOT NULL | - | 通知種別（like/comment/follow） |
| post_id | BIGINT | NULL | NULL | 関連投稿ID（外部キー → posts.id） |
| comment_id | BIGINT | NULL | NULL | 関連コメントID（外部キー → comments.id） |
| is_read | BOOLEAN | NOT NULL | FALSE | 既読フラグ |
| created_at | TIMESTAMP | NOT NULL | NOW() | 作成日時 |

**インデックス**:
- PRIMARY KEY (id)
- INDEX (user_id, is_read, created_at DESC)
- INDEX (actor_id)

**外部キー制約**:
- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
- FOREIGN KEY (actor_id) REFERENCES users(id) ON DELETE CASCADE
- FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
- FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE

**制約**:
- type: 'like', 'comment', 'follow' のいずれか

---

## 🔧 マイグレーション戦略

### GORM AutoMigrate

開発初期はGORMの `AutoMigrate` を使用してスキーマを自動生成

```go
db.AutoMigrate(
    &User{},
    &Post{},
    &Media{},
    &Comment{},
    &PostLike{},
    &Follow{},
)
```

### 本番環境

本番環境ではマイグレーションツール（golang-migrate等）を使用して管理

---

## 🔍 クエリ最適化

### よく使われるクエリとインデックス

1. **タイムライン取得（フォロー中）**
```sql
SELECT p.* FROM posts p
INNER JOIN follows f ON f.following_id = p.user_id
WHERE f.follower_id = ? AND p.deleted_at IS NULL
ORDER BY p.created_at DESC
LIMIT 20;
```
→ インデックス: `follows(follower_id)`, `posts(created_at DESC)`

2. **タイムライン取得（全体）**
```sql
SELECT * FROM posts
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT 20;
```
→ インデックス: `posts(deleted_at, created_at DESC)`

3. **いいね数カウント**
```sql
SELECT COUNT(*) FROM post_likes
WHERE post_id = ?;
```
→ インデックス: `post_likes(post_id)`

4. **コメント一覧取得**
```sql
SELECT * FROM comments
WHERE post_id = ? AND deleted_at IS NULL
ORDER BY created_at ASC;
```
→ インデックス: `comments(post_id, created_at)`

---

## 📝 マイグレーションファイル例

### Phase 1

```
migrations/
├── 001_create_users_table.up.sql
├── 001_create_users_table.down.sql
├── 002_create_posts_table.up.sql
├── 002_create_posts_table.down.sql
├── 003_create_media_table.up.sql
├── 003_create_media_table.down.sql
├── 004_create_comments_table.up.sql
├── 004_create_comments_table.down.sql
├── 005_create_post_likes_table.up.sql
├── 005_create_post_likes_table.down.sql
├── 006_create_follows_table.up.sql
└── 006_create_follows_table.down.sql
```

---

**作成日**: 2026-02-14
**最終更新**: 2026-02-14
