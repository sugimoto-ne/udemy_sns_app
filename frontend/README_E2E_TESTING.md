# E2Eテストクイックガイド

## 🔒 重要: 開発環境への影響はありません

E2Eテストは**完全に分離されたテスト環境**で実行されます。開発環境のデータベースやAPIサーバーには一切影響しません。

## テスト環境の構成

| 環境 | データベース | APIサーバー | フロントエンド |
|------|------------|-----------|---------------|
| **開発** | localhost:5432 | localhost:8080 | localhost:5173 |
| **テスト** | localhost:5433 | localhost:8081 | localhost:5173 (.env.test使用) |

## テスト実行手順

### 1. テスト用サービスを起動

```bash
# frontendディレクトリで実行
npm run test:setup

# または手動で起動
docker compose --profile test up -d db_test api_test
```

### 2. テストを実行

```bash
# すべてのテストを実行
npm run test:e2e

# UIモードで実行（推奨）
npm run test:e2e:ui

# デバッグモードで実行
npm run test:e2e:debug
```

### 3. テスト用サービスを停止

```bash
npm run test:teardown

# または手動で停止
docker compose --profile test down
```

## 利用可能なコマンド

```bash
# テスト環境のセットアップ
npm run test:setup

# E2Eテスト実行
npm run test:e2e          # ヘッドレスモードで実行
npm run test:e2e:ui       # UIモードで実行（視覚的に確認）
npm run test:e2e:debug    # デバッグモードで実行

# テスト環境のクリーンアップ
npm run test:teardown

# 開発サーバー（通常の開発）
npm run dev               # 開発環境API (port 8080) に接続

# テスト用開発サーバー
npm run dev:test          # テスト環境API (port 8081) に接続
```

## トラブルシューティング

### テストが接続エラーで失敗する

```bash
# テスト用サービスの状態を確認
docker compose ps

# sns_db_test と sns_api_test が起動しているか確認
# 起動していない場合:
npm run test:setup
```

### テスト用APIのログを確認

```bash
docker compose logs -f api_test
```

### テスト用データベースをリセット

```bash
# テスト用サービスを停止して再起動
npm run test:teardown
npm run test:setup
```

## 詳細ドキュメント

- テスト実装ガイド: `/docs/E2E_TEST_IMPLEMENTATION_GUIDE.md`
- 総合テストドキュメント: `/docs/TESTING.md`
