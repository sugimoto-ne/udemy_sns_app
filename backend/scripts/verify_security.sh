#!/bin/bash

echo "=========================================="
echo "セキュリティヘッダー検証スクリプト"
echo "=========================================="

API_URL="http://localhost:8080"

# カラーコード
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo ""
echo "📡 サーバーのヘルスチェック..."
HEALTH_RESPONSE=$(curl -s -I "${API_URL}/health")

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ サーバーが起動していません${NC}"
    echo "   docker compose up -d を実行してください"
    exit 1
fi

echo -e "${GREEN}✅ サーバーが起動しています${NC}"
echo ""

echo "=========================================="
echo "🔒 セキュリティヘッダーの確認"
echo "=========================================="

# ヘッダーを取得
HEADERS=$(curl -s -I "${API_URL}/health")

# Content-Security-Policy
echo -n "1. Content-Security-Policy: "
CSP=$(echo "$HEADERS" | grep -i "Content-Security-Policy:" | sed 's/Content-Security-Policy: //')
if [ -n "$CSP" ]; then
    echo -e "${GREEN}✅ 設定されています${NC}"
    echo "   $CSP"
else
    echo -e "${RED}❌ 設定されていません${NC}"
fi

# X-Frame-Options
echo ""
echo -n "2. X-Frame-Options: "
XFO=$(echo "$HEADERS" | grep -i "X-Frame-Options:" | sed 's/X-Frame-Options: //' | tr -d '\r')
if [ "$XFO" = "DENY" ]; then
    echo -e "${GREEN}✅ DENY${NC}"
else
    echo -e "${RED}❌ 設定されていません (期待値: DENY)${NC}"
fi

# X-Content-Type-Options
echo ""
echo -n "3. X-Content-Type-Options: "
XCTO=$(echo "$HEADERS" | grep -i "X-Content-Type-Options:" | sed 's/X-Content-Type-Options: //' | tr -d '\r')
if [ "$XCTO" = "nosniff" ]; then
    echo -e "${GREEN}✅ nosniff${NC}"
else
    echo -e "${RED}❌ 設定されていません (期待値: nosniff)${NC}"
fi

# Strict-Transport-Security
echo ""
echo -n "4. Strict-Transport-Security: "
HSTS=$(echo "$HEADERS" | grep -i "Strict-Transport-Security:" | sed 's/Strict-Transport-Security: //')
if [ -n "$HSTS" ]; then
    echo -e "${GREEN}✅ 設定されています${NC}"
    echo "   $HSTS"
else
    echo -e "${RED}❌ 設定されていません${NC}"
fi

# X-XSS-Protection
echo ""
echo -n "5. X-XSS-Protection: "
XXP=$(echo "$HEADERS" | grep -i "X-XSS-Protection:" | sed 's/X-XSS-Protection: //' | tr -d '\r')
if [ "$XXP" = "0" ]; then
    echo -e "${GREEN}✅ 0 (CSPを優先)${NC}"
else
    echo -e "${RED}❌ 設定されていません (期待値: 0)${NC}"
fi

echo ""
echo "=========================================="
echo "⏱️ レートリミットの確認"
echo "=========================================="

# 認証エンドポイントのレートリミット（5回/分）
echo ""
echo "認証エンドポイント (/api/v1/auth/register) - 5回/分"
echo ""

SUCCESS_COUNT=0
for i in {1..7}; do
    RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}\nREMAINING:%{header_x_ratelimit_remaining}" \
        -X POST "${API_URL}/api/v1/auth/register" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"test${i}@example.com\",\"username\":\"test${i}\",\"password\":\"pass123\"}")

    HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)
    REMAINING=$(echo "$RESPONSE" | grep "REMAINING:" | cut -d: -f2 | tr -d '\r')

    echo -n "  リクエスト ${i}: "
    if [ "$HTTP_CODE" = "429" ]; then
        echo -e "${YELLOW}429 Too Many Requests${NC} (残り: ${REMAINING})"
    elif [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "201" ]; then
        echo -e "${GREEN}${HTTP_CODE} Success${NC} (残り: ${REMAINING})"
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    else
        echo -e "${GREEN}${HTTP_CODE}${NC} (残り: ${REMAINING})"
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    fi

    # 6回目以降はレートリミットに引っかかるはず
    if [ $i -ge 6 ] && [ "$HTTP_CODE" != "429" ]; then
        echo -e "    ${RED}⚠️  6回目以降は429エラーが期待されます${NC}"
    fi
done

echo ""
if [ $SUCCESS_COUNT -eq 5 ]; then
    echo -e "${GREEN}✅ レートリミットが正しく動作しています（5回成功、6回目以降は429）${NC}"
elif [ $SUCCESS_COUNT -lt 5 ]; then
    echo -e "${YELLOW}⚠️  認証エラーにより5回未満の成功（レートリミット自体は動作）${NC}"
else
    echo -e "${RED}❌ レートリミットが正しく動作していません${NC}"
fi

echo ""
echo "=========================================="
echo "📊 検証結果サマリー"
echo "=========================================="
echo ""
echo "セキュリティヘッダー:"
echo "  - Content-Security-Policy"
echo "  - X-Frame-Options"
echo "  - X-Content-Type-Options"
echo "  - Strict-Transport-Security"
echo "  - X-XSS-Protection"
echo ""
echo "レートリミット:"
echo "  - 認証系API: 5回/分"
echo "  - 一般API: 60回/分"
echo ""
echo "詳細レポート: docs/SECURITY_IMPROVEMENT_REPORT.md"
echo ""
echo "=========================================="
