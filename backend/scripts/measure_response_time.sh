#!/bin/bash

# レスポンスタイム計測スクリプト
# タイムライン取得API（GET /api/v1/posts）のパフォーマンスを計測

set -e

API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
ITERATIONS="${ITERATIONS:-10}"

# 色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}レスポンスタイム計測スクリプト${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 1. テストユーザーを作成してトークンを取得
echo -e "${YELLOW}📝 ステップ1: テストユーザーを作成...${NC}"

TIMESTAMP=$(date +%s)
TEST_EMAIL="perf-test-${TIMESTAMP}@example.com"
TEST_USERNAME="perftest${TIMESTAMP}"
TEST_PASSWORD="TestPassword123"

# ユーザー登録
REGISTER_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE_URL}/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"${TEST_EMAIL}\",
    \"username\": \"${TEST_USERNAME}\",
    \"password\": \"${TEST_PASSWORD}\"
  }")

HTTP_CODE=$(echo "$REGISTER_RESPONSE" | tail -n 1)
RESPONSE_BODY=$(echo "$REGISTER_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" != "201" ]; then
  echo -e "${RED}❌ ユーザー登録に失敗しました (HTTP ${HTTP_CODE})${NC}"
  echo "$RESPONSE_BODY"
  exit 1
fi

# トークンを抽出
TOKEN=$(echo "$RESPONSE_BODY" | grep -o '"token":"[^"]*"' | sed 's/"token":"\([^"]*\)"/\1/')

if [ -z "$TOKEN" ]; then
  echo -e "${RED}❌ トークンの取得に失敗しました${NC}"
  echo "$RESPONSE_BODY"
  exit 1
fi

echo -e "${GREEN}✅ テストユーザー作成完了${NC}"
echo "   ユーザー名: ${TEST_USERNAME}"
echo "   トークン: ${TOKEN:0:20}..."
echo ""

# 2. テスト投稿を作成（タイムラインにデータがあることを確認）
echo -e "${YELLOW}📝 ステップ2: テスト投稿を作成...${NC}"

for i in {1..5}; do
  curl -s -X POST "${API_BASE_URL}/api/v1/posts" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer ${TOKEN}" \
    -d "{\"content\": \"パフォーマンステスト投稿 #${i} at $(date)\"}" > /dev/null
done

echo -e "${GREEN}✅ テスト投稿作成完了（5件）${NC}"
echo ""

# 3. レスポンスタイムを計測
echo -e "${YELLOW}📊 ステップ3: レスポンスタイム計測（${ITERATIONS}回）...${NC}"
echo ""

TOTAL_TIME=0
MIN_TIME=999999
MAX_TIME=0
SUCCESS_COUNT=0
FAIL_COUNT=0

for i in $(seq 1 $ITERATIONS); do
  # curlでレスポンスタイムを計測
  START_TIME=$(date +%s%N)

  HTTP_CODE=$(curl -s -w "%{http_code}" -o /dev/null \
    -H "Authorization: Bearer ${TOKEN}" \
    "${API_BASE_URL}/api/v1/posts?limit=20")

  END_TIME=$(date +%s%N)

  # ミリ秒に変換
  RESPONSE_TIME=$(( (END_TIME - START_TIME) / 1000000 ))

  if [ "$HTTP_CODE" = "200" ]; then
    SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    TOTAL_TIME=$((TOTAL_TIME + RESPONSE_TIME))

    # 最小・最大を記録
    if [ $RESPONSE_TIME -lt $MIN_TIME ]; then
      MIN_TIME=$RESPONSE_TIME
    fi
    if [ $RESPONSE_TIME -gt $MAX_TIME ]; then
      MAX_TIME=$RESPONSE_TIME
    fi

    # 色分け（500ms以上は赤、300ms以上は黄色、それ以下は緑）
    if [ $RESPONSE_TIME -ge 500 ]; then
      COLOR=$RED
      STATUS="⚠️ "
    elif [ $RESPONSE_TIME -ge 300 ]; then
      COLOR=$YELLOW
      STATUS="⚠️ "
    else
      COLOR=$GREEN
      STATUS="✅"
    fi

    echo -e "${COLOR}${STATUS} リクエスト${i}: ${RESPONSE_TIME}ms${NC}"
  else
    FAIL_COUNT=$((FAIL_COUNT + 1))
    echo -e "${RED}❌ リクエスト${i}: 失敗 (HTTP ${HTTP_CODE})${NC}"
  fi

  # サーバーへの負荷を避けるため少し待つ
  sleep 0.1
done

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}📊 計測結果${NC}"
echo -e "${BLUE}========================================${NC}"

if [ $SUCCESS_COUNT -gt 0 ]; then
  AVG_TIME=$((TOTAL_TIME / SUCCESS_COUNT))

  echo -e "${GREEN}成功: ${SUCCESS_COUNT}/${ITERATIONS} リクエスト${NC}"

  if [ $FAIL_COUNT -gt 0 ]; then
    echo -e "${RED}失敗: ${FAIL_COUNT}/${ITERATIONS} リクエスト${NC}"
  fi

  echo ""
  echo "📈 レスポンスタイム統計:"
  echo "   平均: ${AVG_TIME}ms"
  echo "   最小: ${MIN_TIME}ms"
  echo "   最大: ${MAX_TIME}ms"
  echo ""

  # 要件との比較（500ms以内）
  if [ $AVG_TIME -le 500 ]; then
    echo -e "${GREEN}✅ 非機能要件達成: 平均レスポンスタイムが500ms以内です${NC}"
  else
    echo -e "${RED}❌ 非機能要件未達成: 平均レスポンスタイムが500msを超えています${NC}"
    echo -e "${YELLOW}   改善が必要です（目標: 500ms以内）${NC}"
  fi
else
  echo -e "${RED}❌ 計測失敗: すべてのリクエストが失敗しました${NC}"
  exit 1
fi

echo ""
echo -e "${BLUE}========================================${NC}"
echo ""

# 4. クリーンアップ（オプション）
if [ "${CLEANUP:-true}" = "true" ]; then
  echo -e "${YELLOW}🧹 クリーンアップ: テストユーザーは残したままにします${NC}"
  echo "   後で手動で削除する場合は、ユーザー名: ${TEST_USERNAME}"
fi

echo ""
echo -e "${GREEN}✅ 計測完了${NC}"
