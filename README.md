# TwitterライクSNSアプリケーション

モダンな技術スタックで構築されたTwitterライクなSNSアプリケーション。

## 技術スタック

- **フロントエンド**: React + TypeScript + Material-UI
- **バックエンド**: Go + Echo + GORM
- **データベース**: PostgreSQL
- **認証**: JWT
- **インフラ**: Docker + Docker Compose

## クイックスタート

### 開発環境の起動

```bash
# ワンコマンドで起動
make dev-up

# または
docker compose up -d
```

アクセス:
- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html
- データベース: localhost:5432

### 開発環境の停止

```bash
make dev-down
```

## テスト実行

### ワンコマンドで実行（推奨）

```bash
# すべてのテストを実行
make test

# バックエンド単体テストのみ
make test-backend

# E2Eテストのみ
make test-e2e

# 利用可能なコマンド一覧
make help
```

**すべて自動で実行されます:**
1. テスト環境を自動起動（テスト用DB + テスト用API）
2. テストを実行
3. テスト環境を自動停止

### テスト環境について

テストは**完全に分離されたテスト専用環境**で実行されます。

| 環境 | データベース | APIサーバー | ポート |
|------|------------|-----------|--------|
| **開発環境** | `sns_db` | `sns_api` | DB: 5432, API: 8080 |
| **テスト環境** | `sns_db_test` | `sns_api_test` | DB: 5433, API: 8081 |

✅ **開発環境のデータに一切影響しません**

## 利用可能なMakeコマンド

```bash
make help          # コマンド一覧を表示
make dev-up        # 開発環境を起動
make dev-down      # 開発環境を停止
make dev-logs      # 開発環境のログを表示
make test          # すべてのテストを実行
make test-backend  # バックエンド単体テストを実行
make test-e2e      # E2Eテストを実行
make test-setup    # テスト環境を起動
make test-teardown # テスト環境を停止
```

## プロジェクト構成

```
.
├── backend/              # Go バックエンド
│   ├── cmd/server/       # エントリポイント
│   ├── internal/         # 内部パッケージ
│   │   ├── handlers/     # HTTPハンドラー
│   │   ├── services/     # ビジネスロジック
│   │   ├── models/       # GORMモデル
│   │   ├── middleware/   # ミドルウェア
│   │   └── utils/        # ユーティリティ
│   └── docs/             # Swagger生成ファイル
├── frontend/             # React フロントエンド
│   ├── src/
│   │   ├── components/   # Reactコンポーネント
│   │   ├── pages/        # ページコンポーネント
│   │   ├── api/          # APIクライアント
│   │   └── hooks/        # カスタムフック
│   └── tests/e2e/        # E2Eテスト
├── docs/                 # ドキュメント
├── docker-compose.yml    # Docker Compose設定
├── Makefile              # 便利なコマンド
└── README.md             # このファイル
```

## ドキュメント

- **プロジェクト概要**: `docs/todo/00_OVERVIEW.md`
- **データベーススキーマ**: `docs/todo/01_DATABASE_SCHEMA.md`
- **API仕様**: `docs/todo/02_API_SPECIFICATION.md`
- **テストガイド**: `docs/TESTING.md`
- **テストクイックリファレンス**: `TESTING_QUICK_REFERENCE.md`
- **E2E実装ガイド**: `docs/E2E_TEST_IMPLEMENTATION_GUIDE.md`
- **プロジェクトガイド**: `CLAUDE.md`

## 開発フロー

1. **機能開発**
   ```bash
   make dev-up        # 開発環境起動
   # コードを書く...
   make dev-logs      # ログを確認
   ```

2. **テスト実行**
   ```bash
   make test-backend  # バックエンドテスト
   make test-e2e      # E2Eテスト（フロントエンド実装後）
   ```

3. **環境停止**
   ```bash
   make dev-down      # 開発環境停止
   ```

## トラブルシューティング

### ポートが使用中

```bash
# ポートを使用しているプロセスを確認
lsof -i :8080   # 開発環境API
lsof -i :5432   # 開発環境DB
lsof -i :8081   # テスト環境API
lsof -i :5433   # テスト環境DB
```

### コンテナの状態確認

```bash
# すべてのコンテナを確認
docker compose ps

# テスト環境も含めて確認
docker compose --profile test ps

# ログを確認
make dev-logs
docker compose logs -f api_test  # テスト環境のログ
```

### データベースのリセット

```bash
# ⚠️ 警告: データが削除されます
docker compose down -v  # ボリュームも削除
docker compose up -d    # 再起動
```

## CI/CD

### GitHub Actions

このプロジェクトはGitHub Actionsで自動テストを実行します。

**ワークフロー**: `.github/workflows/ci.yml`

#### 実行されるテスト

| ジョブ | 内容 |
|-------|------|
| **Backend Tests** | Go単体テスト (PostgreSQL使用) |
| **Frontend Build** | TypeScript/Viteビルドチェック |
| **E2E Tests** | Playwright E2Eテスト |

#### トリガー

- `main`ブランチへのpush
- `main`ブランチへのPull Request

#### Pull Request作成前のチェックリスト

```bash
# ローカルでテストを実行
make test-backend  # バックエンドテスト
npm run build      # フロントエンドビルド

# E2Eテスト（フロントエンド実装後）
make test-e2e
```

すべてのCIチェックが✅グリーンであることを確認してからマージしてください。

---

## ライセンス

MIT

## 貢献

プルリクエストを歓迎します！

### 貢献の流れ

1. このリポジトリをFork
2. 機能ブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add amazing feature'`)
4. ローカルでテストを実行 (`make test-backend`)
5. ブランチにPush (`git push origin feature/amazing-feature`)
6. Pull Requestを作成
7. CIチェックが全てグリーンになるのを確認
8. レビューを待つ
