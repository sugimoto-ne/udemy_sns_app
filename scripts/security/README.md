# セキュリティスキャン

OWASP ZAPを使用したGoバックエンドAPIの脆弱性診断ツール

## クイックスタート

```bash
# 1. APIを起動
docker compose up -d
あらかじ起動しておくこと

# 2. スキャン実行
./scripts/security/run-zap-scan.sh

# 3. レポート確認
open security-reports/zap-report-*.html

# 4. ZAP停止
docker compose --profile security down
```

## 特徴

✅ **JWT認証対応** - 認証が必要なエンドポイントも自動スキャン
✅ **ワンコマンド実行** - セットアップ不要、すぐに使える
✅ **詳細なレポート** - HTML/JSON/XML形式で出力
✅ **OWASP Top 10** - 主要な脆弱性を網羅的に検出

## 検出可能な脆弱性

- SQLインジェクション
- XSS（クロスサイトスクリプティング）
- CSRF（クロスサイトリクエストフォージェリ）
- 認証・認可の不備
- セキュリティヘッダーの欠落
- 機密情報の漏洩
- その他 OWASP Top 10

## 詳細ドキュメント

詳しい使い方は [docs/SECURITY_SCANNING.md](../../docs/SECURITY_SCANNING.md) を参照してください。

## トラブルシューティング

### APIコンテナが起動していません

```bash
docker compose up -d api
```

### ZAPコンテナの起動に失敗

```bash
docker compose --profile security down
docker compose --profile security up -d zap
```

### ポート8090が使用中

```bash
lsof -i :8090
```

---

**作成日**: 2026-02-15
