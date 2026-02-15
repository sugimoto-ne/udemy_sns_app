#!/bin/bash

################################################################################
# OWASP ZAP セキュリティスキャンスクリプト（JWT認証対応）
################################################################################
#
# 使い方:
#   ./scripts/security/run-zap-scan.sh
#
# 必要な環境:
#   - Docker Compose がインストールされていること
#   - APIコンテナ (sns_api) が起動していること
#
################################################################################

set -e

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 設定
API_URL="http://api:8080"
ZAP_API="http://localhost:8090"
REPORT_DIR="./docs/security-reports"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# テストユーザー情報（既存のユーザーを使用）
TEST_EMAIL="test@test.com"
TEST_PASSWORD="password"

################################################################################
# ヘルパー関数
################################################################################

print_header() {
    echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║  OWASP ZAP セキュリティスキャン - SNS API                 ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
    echo ""
}

print_step() {
    echo -e "${GREEN}▶ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

################################################################################
# 前提条件チェック
################################################################################

check_prerequisites() {
    print_step "前提条件をチェック中..."

    # Docker Composeのチェック
    if ! command -v docker &> /dev/null; then
        print_error "Docker がインストールされていません"
        exit 1
    fi

    # APIコンテナのチェック
    if ! docker compose ps api | grep -q "Up"; then
        print_error "APIコンテナが起動していません"
        echo ""
        echo "以下のコマンドでAPIを起動してください:"
        echo "  docker compose up -d api"
        exit 1
    fi

    print_success "前提条件OK"
    echo ""
}

################################################################################
# ZAPコンテナ起動
################################################################################

start_zap() {
    print_step "ZAPコンテナを起動中..."

    if docker compose ps zap | grep -q "Up"; then
        print_info "ZAPは既に起動しています"
    else
        docker compose --profile security up -d zap
        print_info "ZAPの起動を待機中..."
        sleep 10

        # ZAP APIが応答するまで待機
        for i in {1..30}; do
            if curl -s "$ZAP_API" > /dev/null 2>&1; then
                print_success "ZAP起動完了"
                echo ""
                return
            fi
            echo -n "."
            sleep 2
        done

        print_error "ZAPの起動がタイムアウトしました"
        exit 1
    fi
    echo ""
}

################################################################################
# JWT認証トークン取得
################################################################################

get_jwt_token() {
    print_step "JWT認証トークンを取得中..."

    # jqがインストールされているか確認
    if ! command -v jq &> /dev/null; then
        print_error "jq がインストールされていません"
        echo "macOSの場合: brew install jq"
        exit 1
    fi

    # ログインを試行（ホストOSから直接APIにアクセス）
    print_info "ログインを試行中..."
    LOGIN_RESPONSE=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}" \
        http://localhost:8080/api/v1/auth/login 2>/dev/null || echo "")

    JWT_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token // empty')

    # トークン取得の確認
    if [ -z "$JWT_TOKEN" ]; then
        print_error "JWTトークンの取得に失敗しました"
        print_error "テストユーザー（$TEST_EMAIL）でログインできませんでした"
        echo ""
        echo "APIレスポンス:"
        echo "$LOGIN_RESPONSE"
        echo ""
        echo "以下を確認してください:"
        echo "  1. テストユーザーが存在するか"
        echo "  2. メールアドレスとパスワードが正しいか"
        echo "  3. APIが正常に動作しているか"
        echo ""
        exit 1
    fi

    print_success "JWTトークン取得成功"
    print_info "トークン（先頭20文字）: ${JWT_TOKEN:0:20}..."
    echo ""
}

################################################################################
# ZAPスキャン設定
################################################################################

configure_zap() {
    print_step "ZAPスキャン設定を適用中..."

    # 認証ヘッダーを設定（トークンがある場合のみ）
    if [ -n "$JWT_TOKEN" ]; then
        print_info "Authorizationヘッダーを設定中..."

        # Replacerアドオンを使用してグローバルヘッダーを追加
        curl -s "$ZAP_API/JSON/replacer/action/addRule/" \
            --data-urlencode "description=JWT Authorization" \
            --data-urlencode "enabled=true" \
            --data-urlencode "matchType=REQ_HEADER" \
            --data-urlencode "matchString=Authorization" \
            --data-urlencode "replacement=Bearer $JWT_TOKEN" \
            --data-urlencode "initiators=" \
            > /dev/null

        print_success "認証ヘッダー設定完了"
    fi

    # コンテキストの設定
    print_info "スキャンコンテキストを設定中..."
    curl -s "$ZAP_API/JSON/context/action/newContext/?contextName=SNS_API" > /dev/null
    curl -s "$ZAP_API/JSON/context/action/includeInContext/?contextName=SNS_API&regex=http://api:8080/api/v1/.*" > /dev/null

    # OpenAPI定義をインポート
    print_info "OpenAPI定義をインポート中..."
    OPENAPI_URL="http://api:8080/swagger/doc.json"
    curl -s "$ZAP_API/JSON/openapi/action/importUrl/?url=$OPENAPI_URL&hostOverride=api:8080" > /dev/null
    print_success "OpenAPI定義をインポート完了"

    print_success "ZAP設定完了"
    echo ""
}

################################################################################
# ZAPスキャン実行
################################################################################

run_zap_scan() {
    print_step "ZAPスキャンを実行中..."
    print_info "ターゲット: $API_URL/api/v1"
    echo ""

    # スパイダースキャン（クローリング）
    print_info "[1/3] スパイダースキャン（URL探索）を開始..."

    # jqがインストールされているか確認
    if ! command -v jq &> /dev/null; then
        print_error "jq がインストールされていません"
        echo "macOSの場合: brew install jq"
        exit 1
    fi

    SPIDER_SCAN_ID=$(curl -s "$ZAP_API/JSON/spider/action/scan/?url=$API_URL/api/v1&contextName=SNS_API" | jq -r '.scan')

    # スパイダー進行状況の監視
    while true; do
        SPIDER_STATUS=$(curl -s "$ZAP_API/JSON/spider/view/status/?scanId=$SPIDER_SCAN_ID" | jq -r '.status')
        if [ "$SPIDER_STATUS" = "100" ]; then
            break
        fi
        echo -ne "\r  進行状況: $SPIDER_STATUS%   "
        sleep 2
    done
    echo ""
    print_success "スパイダースキャン完了"
    echo ""

    # パッシブスキャン完了を待機
    print_info "[2/3] パッシブスキャン（自動検出）実行中..."
    sleep 5
    while true; do
        RECORDS=$(curl -s "$ZAP_API/JSON/pscan/view/recordsToScan/" | jq -r '.recordsToScan')
        if [ "$RECORDS" = "0" ]; then
            break
        fi
        echo -ne "\r  残りレコード: $RECORDS   "
        sleep 2
    done
    echo ""
    print_success "パッシブスキャン完了"
    echo ""

    # アクティブスキャン
    print_info "[3/3] アクティブスキャン（脆弱性検査）を開始..."
    print_warning "これには5〜10分程度かかる場合があります..."

    ACTIVE_SCAN_ID=$(curl -s "$ZAP_API/JSON/ascan/action/scan/?url=$API_URL/api/v1&contextName=SNS_API" | jq -r '.scan')

    # アクティブスキャン進行状況の監視
    while true; do
        ACTIVE_STATUS=$(curl -s "$ZAP_API/JSON/ascan/view/status/?scanId=$ACTIVE_SCAN_ID" | jq -r '.status')
        if [ "$ACTIVE_STATUS" = "100" ]; then
            break
        fi
        echo -ne "\r  進行状況: $ACTIVE_STATUS%   "
        sleep 3
    done
    echo ""
    print_success "アクティブスキャン完了"
    echo ""
}

################################################################################
# レポート生成
################################################################################

generate_reports() {
    print_step "レポートを生成中..."

    # レポートディレクトリ作成
    mkdir -p "$REPORT_DIR"

    # HTMLレポート
    print_info "HTMLレポートを生成中..."
    curl -s "$ZAP_API/OTHER/core/other/htmlreport/" > "$REPORT_DIR/zap-report-$TIMESTAMP.html"
    print_success "HTML: $REPORT_DIR/zap-report-$TIMESTAMP.html"

    # JSONレポート
    print_info "JSONレポートを生成中..."
    curl -s "$ZAP_API/JSON/core/view/alerts/" > "$REPORT_DIR/zap-report-$TIMESTAMP.json"
    print_success "JSON: $REPORT_DIR/zap-report-$TIMESTAMP.json"

    # XMLレポート
    print_info "XMLレポートを生成中..."
    curl -s "$ZAP_API/OTHER/core/other/xmlreport/" > "$REPORT_DIR/zap-report-$TIMESTAMP.xml"
    print_success "XML:  $REPORT_DIR/zap-report-$TIMESTAMP.xml"

    echo ""
}

################################################################################
# サマリー表示
################################################################################

show_summary() {
    print_step "スキャン結果サマリー"
    echo ""

    # アラート数を集計
    SUMMARY=$(curl -s "$ZAP_API/JSON/core/view/alertsSummary/")
    HIGH=$(echo "$SUMMARY" | jq -r '.High // 0')
    MEDIUM=$(echo "$SUMMARY" | jq -r '.Medium // 0')
    LOW=$(echo "$SUMMARY" | jq -r '.Low // 0')
    INFO=$(echo "$SUMMARY" | jq -r '.Informational // 0')

    echo "  脆弱性検出数:"
    echo -e "    ${RED}High (高):         $HIGH${NC}"
    echo -e "    ${YELLOW}Medium (中):       $MEDIUM${NC}"
    echo -e "    ${BLUE}Low (低):          $LOW${NC}"
    echo -e "    Informational:  $INFO"
    echo ""

    # 重要な脆弱性がある場合は警告
    if [ "$HIGH" -gt 0 ]; then
        print_warning "高リスクの脆弱性が検出されました！早急に対応してください。"
    elif [ "$MEDIUM" -gt 0 ]; then
        print_warning "中リスクの脆弱性が検出されました。対応を検討してください。"
    else
        print_success "高・中リスクの脆弱性は検出されませんでした。"
    fi
    echo ""
}

################################################################################
# クリーンアップと終了メッセージ
################################################################################

print_completion() {
    echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║              スキャン完了！                               ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo "📊 レポートを確認:"
    echo "   open $REPORT_DIR/zap-report-$TIMESTAMP.html"
    echo ""
    echo "🛑 ZAPを停止:"
    echo "   docker compose --profile security down"
    echo ""
    echo "🔄 再スキャン:"
    echo "   ./scripts/security/run-zap-scan.sh"
    echo ""
}

################################################################################
# メイン処理
################################################################################

main() {
    print_header
    check_prerequisites
    start_zap
    get_jwt_token
    configure_zap
    run_zap_scan
    generate_reports
    show_summary
    print_completion
}

# スクリプト実行
main
