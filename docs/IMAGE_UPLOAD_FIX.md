# 画像アップロード問題の修正レポート

## 📅 修正日
2026-02-16

## 🐛 発生した問題

### 問題1: 環境変数が設定されていない
**エラー**: `{error: 'local storage upload not implemented in Phase 2'}`

**原因**: `docker-compose.yml`にFirebase Storageの環境変数が設定されていなかった

**症状**:
- Firebase Storageが初期化されない
- `useFirebase = false` になり、ローカルストレージモードにフォールバック
- ローカルストレージはPhase 2で実装していないため、エラーが発生

### 問題2: 署名付きURLが長すぎる
**エラー**: `ERROR: value too long for type character varying(500) (SQLSTATE 22001)`

**原因**: データベースの`media_url`カラムが`VARCHAR(500)`で、署名付きURLが収まらない

**症状**:
- Firebase Storageへのアップロードは成功
- データベースへの保存時にエラーが発生
- 署名付きURLは通常1000-1500文字以上

---

## ✅ 実施した修正

### 修正1: docker-compose.ymlに環境変数を追加

**ファイル**: `/Users/sugimoto/Desktop/udemy/sns/docker-compose.yml`

**変更内容** (Lines 61-67):
```yaml
environment:
  # ... 既存の環境変数 ...

  # Firebase Storage (Phase 2)
  FIREBASE_CREDENTIALS_PATH: ./service_account.json
  FIREBASE_STORAGE_BUCKET: udemy-sns-b9e40.firebasestorage.app

  # Resend API (Phase 2)
  RESEND_API_KEY: re_S4bMnDF3_N4ptZzYx4wmrLqyPMWrAhL7A
  FROM_EMAIL: sugimo.324@gmail.com
  FRONTEND_URL: http://localhost:5173
```

**実行したコマンド**:
```bash
# コンテナを再作成して環境変数を反映
docker compose up -d api

# 環境変数が設定されたことを確認
docker compose exec api printenv | grep FIREBASE
# Output:
# FIREBASE_STORAGE_BUCKET=udemy-sns-b9e40.firebasestorage.app
# FIREBASE_CREDENTIALS_PATH=./service_account.json
```

---

### 修正2: media_urlカラムサイズを拡張

#### データベーススキーマ変更

**実行したコマンド**:
```bash
docker compose exec db psql -U postgres -d sns_db -c \
  "ALTER TABLE media ALTER COLUMN media_url TYPE VARCHAR(2000);"
```

**変更前**:
```sql
media_url | character varying(500) | not null
```

**変更後**:
```sql
media_url | character varying(2000) | not null
```

#### モデル定義の更新

**ファイル**: `/Users/sugimoto/Desktop/udemy/sns/backend/internal/models/media.go`

**変更内容** (Line 11):
```go
// Before
MediaURL   string    `gorm:"type:varchar(500);not null" json:"media_url"`

// After
MediaURL   string    `gorm:"type:varchar(2000);not null" json:"media_url"` // 署名付きURL対応
```

---

## 📊 署名付きURLのサイズについて

### 通常のURL vs 署名付きURL

**公開URL（Phase 1）**:
```
https://storage.googleapis.com/bucket-name/uploads/abc123.jpg
```
- 長さ: 約60-100文字
- VARCHAR(500)で十分

**署名付きURL（Phase 2）**:
```
https://storage.googleapis.com/bucket-name/uploads/abc123.jpg?GoogleAccessId=...&Expires=1234567890&Signature=very_long_base64_string...
```
- 長さ: 約1000-1500文字
- VARCHAR(2000)が必要

### URLの構成要素

1. **ベースURL**: `https://storage.googleapis.com/bucket-name/uploads/file.jpg` (約100文字)
2. **GoogleAccessId**: サービスアカウントのメールアドレス (約50-100文字)
3. **Expires**: Unix timestamp (10文字)
4. **Signature**: Base64エンコードされた署名 (約800-1200文字)

**合計**: 約1000-1500文字

---

## 🧪 修正後のテスト

### 1. 環境変数の確認
```bash
docker compose exec api printenv | grep -E "FIREBASE|RESEND"
```

**期待される出力**:
```
FIREBASE_STORAGE_BUCKET=udemy-sns-b9e40.firebasestorage.app
FIREBASE_CREDENTIALS_PATH=./service_account.json
RESEND_API_KEY=re_S4bMnDF3_N4ptZzYx4wmrLqyPMWrAhL7A
```

### 2. データベーススキーマの確認
```bash
docker compose exec db psql -U postgres -d sns_db -c "\d media"
```

**期待される出力**:
```
media_url | character varying(2000) | not null
```

### 3. サーバーログの確認
```bash
docker compose logs api | grep -i "Server starting"
```

**期待される出力**:
```
{"level":"info","port":"8080","time":1771251003,"message":"Server starting"}
⇨ http server started on [::]:8080
```

### 4. 画像アップロードテスト（実機）

1. ブラウザで http://localhost:5173 を開く
2. ログイン
3. 投稿フォームに文字を入力
4. 画像ファイルを選択（JPG, PNG, GIF, HEIC - 5MB以下）
5. 「投稿」ボタンをクリック

**期待される動作**:
- ✅ 投稿が作成される
- ✅ 画像が投稿と一緒に表示される
- ✅ 画像をクリックすると正常に表示される（403エラーなし）
- ✅ 画像URLが署名付きURL形式

**画像URLの確認方法**:
1. 投稿された画像を右クリック→「新しいタブで開く」
2. URLバーを確認：
   ```
   https://storage.googleapis.com/.../...?GoogleAccessId=...&Expires=...&Signature=...
   ```
3. 画像が正常に表示されればテスト成功

---

## 🔧 トラブルシューティング

### エラー: "local storage upload not implemented in Phase 2"

**原因**: Firebase Storageの環境変数が設定されていない

**解決方法**:
1. `docker-compose.yml`に環境変数を追加
2. `docker compose up -d api`でコンテナを再作成
3. `docker compose exec api printenv | grep FIREBASE`で確認

### エラー: "value too long for type character varying(500)"

**原因**: `media_url`カラムが小さすぎる

**解決方法**:
1. データベーススキーマを変更：
   ```sql
   ALTER TABLE media ALTER COLUMN media_url TYPE VARCHAR(2000);
   ```
2. モデルファイルを更新（`models/media.go`）

### エラー: "failed to initialize Firebase Storage"

**原因**: Firebase credentials ファイルが見つからない、または無効

**解決方法**:
1. `backend/service_account.json`が存在するか確認
2. ファイルの内容が正しいJSONか確認
3. サービスアカウントの権限を確認（Cloud Storage管理者権限が必要）

### エラー: 403 Forbidden (画像アクセス時)

**原因**: 署名付きURLが正しく生成されていない

**解決方法**:
1. `firebase_storage_service.go`の`SignedURL`実装を確認
2. サービスアカウントの権限を確認
3. バケットのIAM設定を確認

---

## 📝 今後の改善案

### 1. マイグレーションファイルの作成

現在はSQLコマンドで直接変更しましたが、本番環境では以下のようなマイグレーションファイルを作成すべきです：

```sql
-- migrations/000X_extend_media_url.up.sql
ALTER TABLE media ALTER COLUMN media_url TYPE VARCHAR(2000);

-- migrations/000X_extend_media_url.down.sql
ALTER TABLE media ALTER COLUMN media_url TYPE VARCHAR(500);
```

### 2. 環境変数の管理

**開発環境**:
- `docker-compose.yml`で直接指定（現在の実装）

**本番環境**:
- 環境変数はシークレット管理サービスで管理
- `.env`ファイルは`.gitignore`に追加
- サービスアカウントキーはKubernetes Secretsなどで管理

### 3. URL短縮の検討（将来的な最適化）

署名付きURLは長いため、以下の最適化を検討できます：

1. **短縮URLサービス**: データベースにマッピングを保存
2. **CDN統合**: Cloud CDNを使用して短いURLを生成
3. **バケット公開設定**: セキュリティポリシーで公開URLを許可

---

## ✅ 修正完了チェックリスト

- [x] docker-compose.ymlに環境変数を追加
- [x] コンテナを再作成
- [x] 環境変数が正しく設定されたことを確認
- [x] データベーススキーマを変更（media_url: VARCHAR(500) → VARCHAR(2000)）
- [x] モデルファイルを更新
- [x] サーバーが正常に起動することを確認
- [ ] 実機での画像アップロードテスト（ユーザーによる確認待ち）

---

## 📁 変更ファイル

```
/Users/sugimoto/Desktop/udemy/sns/
├── docker-compose.yml                        (環境変数追加)
├── backend/
│   └── internal/
│       └── models/
│           └── media.go                      (MediaURL: VARCHAR(2000))
└── docs/
    └── IMAGE_UPLOAD_FIX.md                  (このファイル)
```

---

## 🎯 次のステップ

1. **実機テスト**: ユーザーに画像アップロードを試してもらう
2. **複数画像テスト**: 最大4枚の画像を同時にアップロード
3. **エラーハンドリング**: さまざまなファイル形式・サイズでテスト
4. **パフォーマンス**: 大きな画像のアップロード時間を確認

---

**修正日**: 2026-02-16
**担当**: Claude Code
**ステータス**: ✅ 修正完了（実機テスト待ち）
