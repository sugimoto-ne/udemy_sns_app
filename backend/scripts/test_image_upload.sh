#!/bin/bash

# Phase 2 画像アップロード機能テストスクリプト
# このスクリプトは、画像アップロード機能が正常に動作しているかをテストします

echo "==================================="
echo "Phase 2 Image Upload Test"
echo "==================================="
echo ""

# テスト用の認証トークンを取得（既存のユーザーでログイン）
echo "1. Testing Authentication..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }')

# Cookieからアクセストークンを抽出（HttpOnlyの場合）
# または、レスポンスボディからトークンを取得
echo "Login Response: $LOGIN_RESPONSE"
echo ""

# 注意: 実際のテストでは、ブラウザまたはPostmanを使用することを推奨
echo "Note: このスクリプトはAPI構造の確認用です。"
echo "実際の画像アップロードテストは、ブラウザまたはPostmanで行ってください。"
echo ""

echo "2. Checking Firebase Storage Configuration..."
# 環境変数の確認
if [ -f ".env" ]; then
    echo "✓ .env file found"

    if grep -q "FIREBASE_STORAGE_BUCKET" .env; then
        BUCKET=$(grep "FIREBASE_STORAGE_BUCKET" .env | cut -d '=' -f2)
        echo "✓ FIREBASE_STORAGE_BUCKET is configured: $BUCKET"
    else
        echo "✗ FIREBASE_STORAGE_BUCKET is not configured"
    fi

    if grep -q "FIREBASE_CREDENTIALS_PATH" .env; then
        CREDS_PATH=$(grep "FIREBASE_CREDENTIALS_PATH" .env | cut -d '=' -f2)
        echo "✓ FIREBASE_CREDENTIALS_PATH is configured: $CREDS_PATH"

        # 相対パスの場合はbackendディレクトリからの相対パスとして解決
        if [ -f "$CREDS_PATH" ]; then
            echo "✓ Firebase credentials file exists"
        else
            echo "✗ Firebase credentials file not found at: $CREDS_PATH"
        fi
    else
        echo "✗ FIREBASE_CREDENTIALS_PATH is not configured"
    fi
else
    echo "✗ .env file not found"
fi
echo ""

echo "3. Testing Post Creation Endpoint..."
# 投稿作成エンドポイントのテスト（認証が必要）
echo "POST /api/v1/posts"
echo "Note: Requires authentication cookie"
echo ""

echo "4. Testing Media Upload Endpoint..."
# メディアアップロードエンドポイントのテスト（認証が必要）
echo "POST /api/v1/posts/:id/media"
echo "Note: Requires authentication cookie and multipart/form-data"
echo ""

echo "==================================="
echo "Manual Test Instructions:"
echo "==================================="
echo ""
echo "1. ブラウザで http://localhost:5173 を開く"
echo "2. テストユーザーでログイン（メール認証済みのユーザー）"
echo "3. 投稿フォームに文字を入力"
echo "4. 画像ファイルを選択（JPG, PNG, GIF, HEIC）"
echo "5. 「投稿」ボタンをクリック"
echo ""
echo "【期待される動作】"
echo "- 投稿が作成される"
echo "- 画像が投稿と一緒に表示される"
echo "- 画像をクリックすると正常に表示される（403エラーなし）"
echo "- 画像URLが署名付きURL形式になっている"
echo ""
echo "【確認方法】"
echo "1. 投稿が表示されたら、画像を右クリック→「新しいタブで開く」"
echo "2. URLバーを確認："
echo "   - 正: https://storage.googleapis.com/...?...&Expires=...&Signature=..."
echo "   - 誤: https://storage.googleapis.com/bucket-name/uploads/..."
echo "3. 画像が正常に表示されればテスト成功"
echo ""
echo "【トラブルシューティング】"
echo "- 403 Forbidden: 署名付きURLが生成されていない可能性"
echo "  → backend/internal/services/firebase_storage_service.go を確認"
echo "- Upload failed: Firebase Storageの設定を確認"
echo "  → .env の FIREBASE_STORAGE_BUCKET と FIREBASE_CREDENTIALS_PATH"
echo ""
echo "==================================="
