# Voice Link Backend Makefile
# 開発・ビルド・テスト・デプロイ用のコマンドを定義

# 変数定義
BINARY_NAME := main
BUILD_DIR := bin

# Go関連の変数
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOMOD := $(GOCMD) mod

# デフォルトターゲット
.PHONY: help
help: ## 利用可能なコマンドを表示
	@echo "Voice Link Backend - 利用可能なコマンド:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# 開発関連
.PHONY: dev
dev: ## 開発サーバーを起動（ホットリロード）
	@echo "🚀 開発サーバーを起動中..."
	@docker compose up -d db
	@sleep 3
	@air

.PHONY: run
run: ## アプリケーションを直接実行
	@echo "🏃 アプリケーションを実行中..."
	@$(GOCMD) run main.go

.PHONY: build
build: ## アプリケーションをビルド
	@echo "🔨 アプリケーションをビルド中..."
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "✅ ビルド完了: $(BUILD_DIR)/$(BINARY_NAME)"


# テスト関連
.PHONY: test
test: ## 全テストを実行
	@echo "🧪 全テストを実行中..."
	@$(GOTEST) ./...

.PHONY: test-verbose
test-verbose: ## 詳細なテスト結果を表示
	@echo "🧪 詳細テストを実行中..."
	@$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage: ## テストカバレッジを測定
	@echo "📊 テストカバレッジを測定中..."
	@$(GOTEST) -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✅ カバレッジレポート生成完了: coverage.html"

.PHONY: test-unit
test-unit: ## ユニットテストのみ実行
	@echo "🔬 ユニットテストを実行中..."
	@$(GOTEST) ./usecase/... ./interface/handler/... ./interface/middleware/...

.PHONY: test-integration
test-integration: ## 統合テストのみ実行
	@echo "🔗 統合テストを実行中..."
	@$(GOTEST) -v ./integration_test.go

# 依存関係管理
.PHONY: deps
deps: ## 依存関係をインストール
	@echo "📦 依存関係をインストール中..."
	@$(GOMOD) download
	@$(GOMOD) tidy

.PHONY: deps-update
deps-update: ## 依存関係を更新
	@echo "🔄 依存関係を更新中..."
	@$(GOMOD) get -u ./...
	@$(GOMOD) tidy

# コード品質
.PHONY: fmt
fmt: ## コードをフォーマット
	@echo "🎨 コードをフォーマット中..."
	@$(GOCMD) fmt ./...

.PHONY: vet
vet: ## コードを静的解析
	@echo "🔍 コードを静的解析中..."
	@$(GOCMD) vet ./...

.PHONY: lint
lint: ## リンターを実行
	@echo "🔧 リンターを実行中..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint がインストールされていません"; \
		echo "   インストール: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: check
check: fmt vet lint test ## コード品質チェック（フォーマット、静的解析、リンター、テスト）

# データベース関連
.PHONY: db-up
db-up: ## データベースを起動
	@echo "🗄️  データベースを起動中..."
	@docker compose up -d db

.PHONY: db-down
db-down: ## データベースを停止
	@echo "🛑 データベースを停止中..."
	@docker compose down

.PHONY: db-reset
db-reset: ## データベースをリセット
	@echo "🔄 データベースをリセット中..."
	@docker compose down -v
	@docker compose up -d db

# クリーンアップ
.PHONY: clean
clean: ## ビルドファイルとキャッシュを削除
	@echo "🧹 クリーンアップ中..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@rm -f tmp/main


# 開発環境セットアップ
.PHONY: setup
setup: ## 開発環境をセットアップ
	@echo "⚙️  開発環境をセットアップ中..."
	@$(GOMOD) download
	@docker compose up -d db
	@sleep 5
	@echo "✅ セットアップ完了！"
	@echo "   開発サーバー起動: make dev"
	@echo "   テスト実行: make test"


# デフォルトターゲット
.DEFAULT_GOAL := help
