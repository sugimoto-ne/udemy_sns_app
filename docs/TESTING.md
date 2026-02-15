# テストドキュメント

このドキュメントでは、TwitterライクSNSアプリケーションのテスト方法について説明します。

## 目次

1. [バックエンドテスト（ユニットテスト）](#バックエンドテストユニットテスト)
2. [フロントエンドテスト（E2Eテスト）](#フロントエンドテストe2eテスト)
3. [テスト環境のセットアップ](#テスト環境のセットアップ)
4. [トラブルシューティング](#トラブルシューティング)

---

## バックエンドテスト（ユニットテスト）

### 🔒 重要: 開発環境への影響なし

バックエンド単体テストは**完全に分離されたテスト専用環境**で実行されます:

| 環境 | データベース | APIサーバー | ポート |
|------|------------|-----------|--------|
| **開発環境** | `sns_db` | `sns_api` | DB: 5432, API: 8080 |
| **テスト環境** | `sns_db_test` | `sns_api_test` | DB: 5433, API: 8081 |

✅ **テストを実行しても開発環境のデータは一切影響を受けません**

### 概要

バックエンドテストは、Go言語の標準テストフレームワークを使用したユニット/インテグレーションテストです。

**テスト対象:**
- 認証（ユーザー登録、ログイン、JWT生成・検証）
- 認可（他人のリソースへのアクセス制御）
- 投稿のCRUD操作
- ユーティリティ関数

**実行環境:**
- テスト専用のAPIコンテナ (`sns_api_test`) で実行
- テスト専用のデータベース (`sns_db_test`) に接続
- 開発環境とは完全分離

### テスト実行方法

#### 🚀 推奨: Makefileを使う（ワンコマンド）

**すべて自動で実行されます（環境起動→テスト実行→環境停止）**

```bash
# バックエンド単体テストを実行（すべて自動）
make test-backend

# 利用可能なコマンドを表示
make help
```

#### 📝 手動で実行する場合（3ステップ）

<details>
<summary>クリックして展開</summary>

##### 1. テスト用環境を起動

```bash
# テスト用DBとテスト用APIを起動（profileを使用）
docker compose --profile test up -d db_test api_test

# または
make test-setup

# 起動確認（sns_db_test と sns_api_test が起動していることを確認）
docker compose ps
```

##### 2. テストを実行

```bash
# すべてのテストを実行（テスト用APIコンテナで実行）
docker compose exec api_test go test ./...

# 特定のパッケージをテスト
docker compose exec api_test go test ./internal/services

# 詳細な出力でテスト実行
docker compose exec api_test go test -v ./...

# カバレッジ付きでテスト実行
docker compose exec api_test go test -cover ./...

# カバレッジレポートを生成
docker compose exec api_test go test -coverprofile=coverage.out ./...
docker compose exec api_test go tool cover -html=coverage.out -o coverage.html
```

##### 3. テスト後のクリーンアップ

```bash
# テスト用コンテナを停止（テスト用DBとテスト用APIを停止）
docker compose --profile test down

# または
make test-teardown
```

</details>

### テストファイルの構成

```
backend/
├── internal/
│   ├── services/
│   │   ├── auth_service.go
│   │   ├── auth_service_test.go          # 認証テスト
│   │   ├── post_service.go
│   │   ├── post_service_test.go          # 投稿CRUDテスト
│   │   └── post_authorization_test.go    # 認可テスト
│   ├── utils/
│   │   ├── jwt.go
│   │   └── jwt_test.go                   # JWTテスト
│   └── testutil/
│       ├── testdb.go                     # テストDB接続
│       └── helpers.go                    # テストヘルパー関数
```

### テスト環境変数

テストでは以下の環境変数を使用します（デフォルト値あり）:

```bash
DB_TEST_HOST=localhost
DB_TEST_PORT=5433
DB_TEST_USER=postgres
DB_TEST_PASSWORD=postgres
DB_TEST_NAME=sns_db_test
```

### テストの特徴

1. **自動マイグレーション**: テスト実行前に自動でDBマイグレーションが実行されます
2. **データクリーンアップ**: 各テストケースの前後でDBをクリーンアップ
3. **独立したテスト**: 各テストは他のテストに依存せず独立して実行可能
4. **ヘルパー関数**: テストデータの作成用ヘルパー関数を提供

---

## フロントエンドテスト（E2Eテスト）

### 🔒 重要: 開発環境への影響なし

E2Eテストは**完全に分離されたテスト専用環境**で実行されます:

| 環境 | データベース | APIサーバー | ポート |
|------|------------|-----------|--------|
| **開発環境** | `sns_db` | `sns_api` | DB: 5432, API: 8080 |
| **テスト環境** | `sns_db_test` | `sns_api_test` | DB: 5433, API: 8081 |

✅ **テストを実行しても開発環境のデータは一切影響を受けません**

### ⚠️ 注意: フロントエンド実装が必要

E2Eテストは**フロントエンドが実装された後に実行可能**になります。現在、テストコードは作成済みですが、以下の実装が必要です:

1. フロントエンドコンポーネントの実装
2. `data-testid`属性の追加
3. ルーティングの設定

### 概要

フロントエンドテストは、Playwrightを使用したEnd-to-End（E2E）テストです。実際のブラウザを使用して、ユーザーの操作をシミュレートします。

**テストシナリオ:**
- ユーザー登録 → ログイン → ログアウト
- ツイート投稿 → 編集 → 削除
- いいね → いいね解除
- フォロー → フォロー解除
- プロフィール編集
- 認証が必要なページへのアクセス制御

### テスト実行方法

#### 1. 前提条件

**重要**: E2Eテストは開発環境とは**完全に分離**されています。

- テスト用DB (`sns_db_test`) とテスト用API (`sns_api_test`) を使用
- 開発環境 (`sns_db`, `sns_api`) には**一切影響しません**

#### 2. Playwrightブラウザをインストール（初回のみ）

```bash
cd frontend
npx playwright install
```

#### 3. テストを実行

##### 🚀 推奨: Makefileを使う（ワンコマンド）

**すべて自動で実行されます（環境起動→テスト実行→環境停止）**

```bash
# プロジェクトルートで実行
make test-e2e
```

##### 📝 手動で実行する場合

<details>
<summary>クリックして展開</summary>

###### a. テスト用サービスを起動

```bash
cd frontend

# テスト用サービスを起動（テスト用DBとテスト用API）
npm run test:setup

# または、プロジェクトルートで
make test-setup

# 起動確認
docker compose ps
# sns_db_test と sns_api_test が起動していることを確認
```

###### b. テストを実行

```bash
cd frontend

# すべてのE2Eテストを実行
npm run test:e2e

# UIモードでテスト実行（視覚的に確認しながら実行）
npm run test:e2e:ui

# デバッグモードでテスト実行
npm run test:e2e:debug

# 特定のテストファイルのみ実行
npx playwright test tests/e2e/auth.spec.ts

# ヘッドレスモードを無効化してブラウザを表示
npx playwright test --headed

# 特定のブラウザでテスト実行
npx playwright test --project=chromium
```

###### c. テスト用サービスを停止

```bash
cd frontend
npm run test:teardown

# または、プロジェクトルートで
make test-teardown
```

</details>

#### 6. テストレポートを表示

```bash
# HTMLレポートを表示
npx playwright show-report
```

### テストファイルの構成

```
frontend/
├── tests/
│   └── e2e/
│       ├── helpers.ts              # テストヘルパー関数
│       ├── auth.spec.ts            # 認証フローのテスト
│       ├── posts.spec.ts           # 投稿操作のテスト
│       ├── social.spec.ts          # ソーシャルインタラクションのテスト
│       └── profile.spec.ts         # プロフィール管理のテスト
├── playwright.config.ts            # Playwright設定ファイル
└── package.json
```

### テストの特徴

1. **自動サーバー起動**: テスト実行時にVite開発サーバーが自動起動
2. **スクリーンショット**: テスト失敗時に自動でスクリーンショットを撮影
3. **トレース記録**: 最初のリトライ時にトレースを記録
4. **順次実行**: DBの競合を避けるため、テストは順次実行

### E2Eテストのベストプラクティス

1. **data-testid属性を使用**: セレクターは `data-testid` 属性を優先的に使用
   ```tsx
   <button data-testid="submit-button">送信</button>
   ```

2. **独立したテストケース**: 各テストケースは独立して実行可能にする

3. **テストデータの管理**: ユニークなテストデータを生成する（タイムスタンプ使用）

4. **待機処理**: 適切な待機処理を使用（`waitForSelector`, `waitForURL`）

---

## テスト環境のセットアップ

### 初回セットアップ

1. **バックエンドのセットアップ**

```bash
# プロジェクトルートで実行
docker compose build
docker compose up -d
```

2. **フロントエンドのセットアップ**

```bash
cd frontend
npm install
npx playwright install
```

3. **テスト環境の確認**

```bash
# プロジェクトルートで実行
# テスト環境は必要な時だけ起動します（通常は起動不要）
docker compose --profile test up -d db_test api_test

# 確認したら停止してOK
docker compose --profile test down
```

### テスト実行フロー

#### バックエンドテストの実行フロー

**重要**: バックエンド単体テストも**テスト専用環境**で実行します。開発環境には影響しません。

```
1. テスト用コンテナを起動（テスト用DBとテスト用API）
   └─> docker compose --profile test up -d db_test api_test

2. テストを実行（テスト用APIコンテナで実行）
   └─> docker compose exec api_test go test ./...

3. テスト用コンテナを停止（オプション）
   └─> docker compose --profile test down
```

**環境分離:**
- **開発環境** (`sns_db`, `sns_api`): ポート5432, 8080 ← **影響なし**
- **テスト環境** (`sns_db_test`, `sns_api_test`): ポート5433, 8081 ← **テスト実行**

#### E2Eテストの実行フロー

**重要**: E2Eテストは**テスト専用のAPIサーバー**と**テスト用データベース**を使用します。開発環境には影響を与えません。

```
1. テスト用サービスを起動（テスト用DBとテスト用API）
   └─> npm run test:setup
   または
   └─> docker compose --profile test up -d db_test api_test

2. E2Eテストを実行（Viteサーバーは自動起動、テスト用APIに接続）
   └─> cd frontend && npm run test:e2e

3. テストレポートを確認
   └─> npx playwright show-report

4. テスト用サービスを停止（オプション）
   └─> npm run test:teardown
   または
   └─> docker compose --profile test down
```

**テスト環境の構成:**
- **テスト用DB**: `localhost:5433` (コンテナ: `sns_db_test`)
- **テスト用API**: `localhost:8081` (コンテナ: `sns_api_test`)
- **フロントエンド**: `localhost:5173` (`.env.test`でテスト用APIに接続)
- **開発環境DB**: `localhost:5432` (コンテナ: `sns_db`) ← **影響なし**
- **開発環境API**: `localhost:8080` (コンテナ: `sns_api`) ← **影響なし**

---

## トラブルシューティング

### バックエンドテスト

#### エラー: テスト用DBに接続できない

```bash
# テスト用DBコンテナが起動しているか確認
docker compose ps

# テスト用DBコンテナを起動
docker compose --profile test up -d db_test

# ログを確認
docker compose logs db_test
```

#### エラー: "dial tcp: lookup db_test"

テスト用DBのホスト名は `localhost` です（`db_test` ではありません）。
環境変数 `DB_TEST_HOST=localhost` を設定してください。

#### エラー: ポートが既に使用されている

```bash
# ポート5433を使用しているプロセスを確認
lsof -i :5433

# テスト用DBコンテナを再起動
docker compose --profile test restart db_test
```

### E2Eテスト

#### エラー: "Error: page.goto: net::ERR_CONNECTION_REFUSED"

テスト用APIサーバーが起動していません。

```bash
# テスト用サービスが起動しているか確認
docker compose ps

# テスト用サービスを起動
docker compose --profile test up -d db_test api_test
# または
cd frontend && npm run test:setup

# 起動確認（sns_db_test と sns_api_test が起動していることを確認）
docker compose ps
```

#### エラー: "Test timeout of 30000ms exceeded"

テストがタイムアウトしました。テスト用APIのレスポンスが遅い可能性があります。

```bash
# テスト用APIのログを確認
docker compose logs -f api_test

# テストのタイムアウト時間を延長（playwright.config.ts）
timeout: 60000  // 60秒に延長
```

#### エラー: Playwrightブラウザがインストールされていない

```bash
cd frontend
npx playwright install
```

#### E2Eテストが不安定

- ネットワークの遅延やDBの競合が原因の可能性があります
- `playwright.config.ts` で `workers: 1` に設定し、順次実行を確認
- リトライ回数を増やす: `retries: 2`

---

## CI/CD統合

### GitHub Actions の例

```yaml
name: Tests

on: [push, pull_request]

jobs:
  backend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Start test database
        run: docker compose --profile test up -d db_test
      - name: Run backend tests
        run: docker compose exec -T api go test ./...

  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Start services
        run: docker compose up -d
      - name: Install dependencies
        run: cd frontend && npm install
      - name: Install Playwright
        run: cd frontend && npx playwright install --with-deps
      - name: Run E2E tests
        run: cd frontend && npm run test:e2e
      - uses: actions/upload-artifact@v3
        if: failure()
        with:
          name: playwright-report
          path: frontend/playwright-report/
```

---

## テストカバレッジの目標

### バックエンド

- **目標カバレッジ**: 80%以上
- **重要パッケージ**:
  - `services`: 90%以上
  - `utils`: 85%以上
  - `handlers`: 70%以上

### フロントエンド

- **E2Eテスト**: 主要ユーザーフロー100%カバー
  - ✅ ユーザー登録・ログイン
  - ✅ 投稿CRUD操作
  - ✅ いいね・フォロー
  - ✅ プロフィール管理
  - ✅ 認証・認可

---

## 参考リンク

- [Go Testing Package](https://pkg.go.dev/testing)
- [Playwright Documentation](https://playwright.dev/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [GORM Documentation](https://gorm.io/docs/)

---

## まとめ

このプロジェクトでは、バックエンドとフロントエンドの両方で包括的なテストを実装しています。

- **バックエンド**: ✅ Go標準のテストフレームワークでユニット/インテグレーションテスト（実装完了）
- **フロントエンド**: 🚧 PlaywrightでE2Eテスト（テストコード作成済み、フロントエンド実装後に実行可能）

### 現在のテスト状況

#### ✅ 実装完了（すぐに実行可能）
- バックエンドユニットテスト
  - 認証（登録、ログイン、JWT）
  - 認可（アクセス制御）
  - 投稿CRUD操作
- テスト用DBコンテナ
- テストヘルパー関数

#### 🚧 準備中（フロントエンド実装が必要）
- E2Eテスト
  - テストコードは作成済み
  - フロントエンド実装後に`data-testid`属性を追加する必要あり
  - 詳細は `/docs/E2E_TEST_IMPLEMENTATION_GUIDE.md` を参照

テストを定期的に実行し、コードの品質を維持してください。
