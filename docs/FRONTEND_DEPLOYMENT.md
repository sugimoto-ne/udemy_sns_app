# フロントエンド デプロイ手順書

このドキュメントでは、ReactフロントエンドをFirebase Hostingにデプロイする手順を説明します。

---

## 📋 前提条件

### 必要なツール

- Node.js 18以上
- Firebase CLI
- Firebaseプロジェクト（`udemy-sns-b9e40`）

### Firebase CLIのインストール

```bash
npm install -g firebase-tools
```

---

## 🔐 Firebase認証

初回のみFirebaseにログインします：

```bash
firebase login
```

ブラウザが開くので、Googleアカウントでログイン。

---

## ⚙️ 環境変数の設定

### 1. `.env.production`ファイルを編集

`frontend/.env.production`を開いて、**バックエンドAPIのURL**を設定します：

```bash
# バックエンドAPIのURL（Cloud Run デプロイ後のURL）
VITE_API_BASE_URL=https://your-backend-service.run.app/api/v1
```

**重要**: Cloud Runにバックエンドをデプロイして、実際のURLを取得してから設定してください。

#### Cloud Run URLの確認方法

```bash
gcloud run services describe sns-backend --region us-central1 --format 'value(status.url)'
```

または、GCPコンソール → Cloud Run → サービス → URLをコピー

---

## 🚀 デプロイ手順

### 方法1: 簡単デプロイ（推奨）

```bash
cd frontend
npm run deploy
```

このコマンドは以下を実行します：
1. TypeScriptのビルド
2. Viteで本番ビルド（`.env.production`を使用）
3. Firebase Hostingにデプロイ

### 方法2: 手動ステップ

```bash
cd frontend

# 1. 本番ビルド
npm run build:prod

# 2. ビルド結果を確認
ls -lh dist/

# 3. Firebase Hostingにデプロイ
firebase deploy --only hosting
```

---

## 📦 デプロイされる内容

- **URL**: `https://udemy-sns-b9e40.web.app/`
- **公開ディレクトリ**: `dist/`
- **SPA対応**: 全てのルートが`index.html`にリダイレクト
- **キャッシュ設定**:
  - HTML: キャッシュなし
  - 画像: 1年間キャッシュ
  - JS/CSS: 1年間キャッシュ（ハッシュ付きファイル名）

---

## ✅ デプロイ後の確認

### 1. デプロイ完了メッセージを確認

```
✔  Deploy complete!

Project Console: https://console.firebase.google.com/project/udemy-sns-b9e40/overview
Hosting URL: https://udemy-sns-b9e40.web.app
```

### 2. ブラウザでアクセス

```
https://udemy-sns-b9e40.web.app/
```

### 3. ネットワークタブで確認

1. ブラウザの開発者ツールを開く
2. Networkタブを開く
3. APIリクエストが正しいバックエンドURLに向かっているか確認

```
Request URL: https://your-backend-service.run.app/api/v1/auth/me
```

### 4. 機能テスト

- [ ] ログイン画面が表示される
- [ ] ログインできる
- [ ] タイムラインが表示される
- [ ] 投稿できる
- [ ] 画像アップロードができる
- [ ] いいね・コメントが機能する

---

## 🔧 トラブルシューティング

### 問題1: APIリクエストが失敗する

**症状**: `Failed to fetch` エラー

**原因**: バックエンドURLが間違っている、またはCORSエラー

**解決策**:
1. `.env.production`のURLを確認
2. バックエンドのCORS設定を確認（フロントエンドのURLを許可）

```go
// backend/internal/middleware/cors_middleware.go
AllowOrigins: []string{
    "http://localhost:5173",
    "https://udemy-sns-b9e40.web.app",  // ← 追加
    "https://udemy-sns-b9e40.firebaseapp.com", // ← 追加
},
```

### 問題2: 404エラー（ページがない）

**症状**: リロードすると404エラー

**原因**: `firebase.json`のrewriteルールが正しく設定されていない

**解決策**: `firebase.json`を確認
```json
"rewrites": [
  {
    "source": "**",
    "destination": "/index.html"
  }
]
```

### 問題3: 環境変数が反映されない

**症状**: APIリクエストが`localhost:8080`に向かう

**原因**: ビルド時に`.env.production`が読み込まれていない

**解決策**:
```bash
# 再ビルド（--mode productionを明示）
npm run build:prod

# distを削除してクリーンビルド
rm -rf dist
npm run build:prod
```

### 問題4: ビルドエラー

**症状**: TypeScriptエラーでビルドが失敗

**解決策**:
```bash
# 型エラーを修正
npm run lint

# TypeScriptの型チェック
npx tsc --noEmit
```

---

## 🔄 更新デプロイ

コードを更新した後の再デプロイ手順：

```bash
cd frontend

# 1. 変更をコミット
git add .
git commit -m "Update: ..."

# 2. デプロイ
npm run deploy
```

---

## 🗑️ ロールバック

前のバージョンに戻す方法：

### Firebase Hostingコンソールで操作

1. [Firebase Console](https://console.firebase.google.com/project/udemy-sns-b9e40/hosting) を開く
2. 「Hosting」→「リリース履歴」
3. 戻したいバージョンを選択
4. 「ロールバック」をクリック

### CLIで操作

```bash
# デプロイ履歴を確認
firebase hosting:channel:list

# 特定のバージョンにロールバック
firebase hosting:clone SOURCE_SITE_ID:SOURCE_CHANNEL_ID TARGET_SITE_ID:live
```

---

## 📊 パフォーマンス最適化

### 画像最適化

```bash
# WebPに変換（オプション）
npm install -D vite-plugin-imagemin
```

### バンドルサイズの確認

```bash
npm run build:prod

# バンドル分析
npx vite-bundle-visualizer
```

---

## 🔐 セキュリティ設定

### Content Security Policy（CSP）

`firebase.json`にCSPヘッダーを追加することもできます：

```json
{
  "key": "Content-Security-Policy",
  "value": "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; connect-src 'self' https://your-backend-service.run.app"
}
```

---

## 📝 デプロイチェックリスト

デプロイ前に確認：

- [ ] `.env.production`にバックエンドURLを設定
- [ ] バックエンドがデプロイ済み
- [ ] バックエンドのCORSにフロントエンドURLを追加
- [ ] ローカルでビルドテスト（`npm run build:prod`）
- [ ] Gitにコミット済み
- [ ] Firebase CLIでログイン済み

---

## 🔗 参考リンク

- [Firebase Hosting ドキュメント](https://firebase.google.com/docs/hosting)
- [Vite デプロイガイド](https://vitejs.dev/guide/static-deploy.html)
- Firebase Console: https://console.firebase.google.com/project/udemy-sns-b9e40

---

**作成日**: 2026-02-16
**最終更新**: 2026-02-16
