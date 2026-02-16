.PHONY: help test test-backend test-e2e test-setup test-teardown dev-up dev-down dev-logs

# デフォルトターゲット: ヘルプを表示
help:
	@echo "📚 利用可能なコマンド:"
	@echo ""
	@echo "🔧 開発環境:"
	@echo "  make dev-up        - 開発環境を起動（DB + API）"
	@echo "  make dev-down      - 開発環境を停止"
	@echo "  make dev-logs      - 開発環境のログを表示"
	@echo ""
	@echo "🧪 テスト:"
	@echo "  make test          - すべてのテストを実行（バックエンド + E2E）"
	@echo "  make test-backend  - バックエンド単体テストを実行"
	@echo "  make test-e2e      - E2Eテストを実行"
	@echo ""
	@echo "⚙️  テスト環境:"
	@echo "  make test-setup    - テスト環境を起動"
	@echo "  make test-teardown - テスト環境を停止"
	@echo ""

# 開発環境
dev-up:
	@echo "🚀 開発環境を起動中..."
	docker compose up -d
	@echo "✅ 開発環境が起動しました"
	@echo "   API: http://localhost:8080"
	@echo "   DB:  localhost:5432"

dev-down:
	@echo "⏸️  開発環境を停止中..."
	docker compose down
	@echo "✅ 開発環境が停止しました"

dev-logs:
	docker compose logs -f api

# テスト環境のセットアップ
test-setup:
	@echo "🔧 テスト環境を起動中..."
	docker compose --profile test up -d db_test api_test
	@echo "⏳ テスト環境の準備を待機中..."
	@sleep 3
	@echo "✅ テスト環境が起動しました"
	@echo "   テスト用API: http://localhost:8081"
	@echo "   テスト用DB:  localhost:5433"

test-teardown:
	@echo "🧹 テスト環境を停止中..."
	docker compose stop api_test db_test
	docker compose rm -f api_test db_test
	@echo "✅ テスト環境が停止しました（開発環境は起動中）"

# バックエンド単体テスト（自動でセットアップ→テスト→クリーンアップ）
test-backend:
	@echo "🧪 バックエンド単体テストを実行します..."
	@echo ""
	@$(MAKE) test-setup
	@echo ""
	@echo "▶️  テストを実行中（並列数: 2、パッケージ順次実行）..."
	@docker compose exec -T api_test go test -v -parallel 2 -p=1 ./... || ($(MAKE) test-teardown && exit 1)
	@echo ""
	@$(MAKE) test-teardown
	@echo ""
	@echo "✅ バックエンド単体テスト完了"

# E2Eテスト（自動でセットアップ→テスト→クリーンアップ）
test-e2e:
	@echo "🧪 E2Eテストを実行します..."
	@echo ""
	@$(MAKE) test-setup
	@echo ""
	@echo "▶️  E2Eテストを実行中..."
	@cd frontend && npm run test:e2e || (cd .. && $(MAKE) test-teardown && exit 1)
	@echo ""
	@$(MAKE) test-teardown
	@echo ""
	@echo "✅ E2Eテスト完了"

# すべてのテストを実行
test:
	@echo "🧪 すべてのテストを実行します..."
	@echo ""
	@$(MAKE) test-backend
	@echo ""
	@$(MAKE) test-e2e
	@echo ""
	@echo "🎉 すべてのテスト完了！"
