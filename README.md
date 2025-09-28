# Voice Link Backend API

Voice Linkは、リアルタイム翻訳機能を提供するアプリケーションのバックエンドAPIです。将来的にはFlutterアプリケーションと連携して、音声のリアルタイム翻訳機能を実現する予定です。

## 🚀 プロジェクト概要

このプロジェクトは、ユーザー認証機能を備えたRESTful APIを提供します。現在は基本的なユーザー管理機能が実装されており、今後リアルタイム翻訳機能の追加を予定しています。

### 主な機能（現在実装済み）
- ユーザー登録・認証（JWT認証）
- ユーザー情報の取得・更新
- PostgreSQLデータベース連携

### 将来の機能（予定）
- リアルタイム音声翻訳
- WebSocket対応
- Flutterアプリケーションとの連携
- 多言語翻訳API統合

## 🛠 技術スタック

- **言語**: Go 1.24.3
- **Webフレームワーク**: Echo v4
- **データベース**: PostgreSQL 16
- **ORM**: GORM
- **認証**: JWT
- **コンテナ化**: Docker & Docker Compose
- **開発環境**: Air（ホットリロード）

## 📁 プロジェクト構造

```
voice_link_backend/
├── domain/                          # ドメイン層
│   └── model/
│       └── user.go                  # ユーザーモデル（エンティティ・リポジトリインターフェース）
├── infrastructure/                  # インフラストラクチャ層
│   └── persistence/
│       └── user_repository.go       # データベース操作の実装
├── interface/                       # インターフェース層
│   ├── handler/                     # HTTPハンドラー
│   │   ├── auth/                    # 認証関連ハンドラー
│   │   │   ├── auth_handler.go      # 認証ハンドラー（登録・ログイン・パスワードリセット）
│   │   │   └── auth_handler_test.go # 認証ハンドラーのテスト
│   │   ├── user/                    # ユーザー管理ハンドラー
│   │   │   ├── user_handler.go      # ユーザー管理ハンドラー（CRUD操作）
│   │   │   └── user_handler_test.go # ユーザー管理ハンドラーのテスト
│   │   └── common/                  # 共通機能
│   │       ├── request_types.go     # 共通のリクエスト・レスポンス型
│   │       ├── response_helper.go   # 共通のレスポンス処理
│   │       └── mock_usecase.go      # テスト用モック実装
│   ├── middleware/                  # ミドルウェア
│   │   ├── auth.go                  # JWT認証ミドルウェア
│   │   └── auth_test.go             # 認証ミドルウェアのテスト
│   └── router/                      # ルーティング設定
│       └── router.go                # APIルート定義
├── usecase/                         # ユースケース層（ビジネスロジック）
│   ├── user_usecase.go              # ユーザー関連のビジネスロジック
│   └── user_usecase_test.go         # ユースケースのテスト
├── main.go                          # アプリケーションエントリーポイント
├── integration_test.go              # 統合テスト
├── openapi.yml                      # API仕様書（OpenAPI 3.1）
├── docker-compose.yml               # Docker Compose設定
├── Dockerfile                       # Dockerイメージ設定
├── go.mod                           # Go依存関係管理
└── go.sum                           # Go依存関係のチェックサム
```

### 🏗️ アーキテクチャパターン

このプロジェクトは **Clean Architecture** の原則に従って設計されています：

- **Domain Layer**: ビジネスルールとエンティティを定義
- **UseCase Layer**: アプリケーション固有のビジネスロジック
- **Interface Layer**: 外部との入出力（HTTP、データベース）を処理
- **Infrastructure Layer**: 外部システムとの具体的な実装

### 📦 パッケージ構成

- **`domain/model`**: エンティティとリポジトリインターフェース
- **`usecase`**: ビジネスロジックの実装
- **`infrastructure/persistence`**: データベース操作の実装
- **`interface/handler`**: HTTPリクエストの処理
- **`interface/middleware`**: 共通処理（認証など）
- **`interface/router`**: ルーティング設定


## 🚀 セットアップ

### 前提条件
- Go 1.24.3以上
- Docker & Docker Compose
- PostgreSQL（Dockerを使用する場合は不要）

### 1. リポジトリのクローン
```bash
git clone <repository-url>
cd voice_link_backend
```

### 2. 環境変数の設定
```bash
# .envファイルを作成（オプション）
cp .env.example .env

# 必要な環境変数
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=voice_link
DB_PORT=5432
JWT_SECRET=your-secret-key-change-in-production
PORT=8080
```

### 3. Docker Composeを使用した起動（推奨）
```bash
# アプリケーションとデータベースを起動
docker-compose up -d

# ログの確認
docker-compose logs -f app
```

### 4. ローカル開発環境での起動

#### Makefileを使用（推奨）
```bash
# 開発環境のセットアップ
make setup

# 開発サーバーを起動（ホットリロード）
make dev

# アプリケーションを直接実行
make run

# 全テストを実行
make test

# 利用可能なコマンドを確認
make help
```

#### 手動での起動
```bash
# 依存関係のインストール
go mod download

# データベースの起動（Docker Composeを使用）
docker-compose up -d db

# アプリケーションの起動
go run main.go
```

## 🛠️ 開発用コマンド

このプロジェクトには便利なMakefileが含まれています：

### 基本コマンド
- `make help` - 利用可能なコマンドを表示
- `make setup` - 開発環境をセットアップ
- `make dev` - 開発サーバーを起動（ホットリロード）
- `make run` - アプリケーションを直接実行
- `make build` - アプリケーションをビルド

### テストコマンド
- `make test` - 全テストを実行
- `make test-unit` - ユニットテストのみ実行
- `make test-integration` - 統合テストのみ実行
- `make test-coverage` - テストカバレッジを測定

### コード品質
- `make fmt` - コードをフォーマット
- `make vet` - コードを静的解析
- `make lint` - リンターを実行
- `make check` - コード品質チェック（全実行）

### データベース
- `make db-up` - データベースを起動
- `make db-down` - データベースを停止
- `make db-reset` - データベースをリセット

### Docker
- `make docker-build` - Dockerイメージをビルド
- `make docker-compose-up` - Docker Composeで全サービスを起動
- `make docker-logs` - Dockerコンテナのログを表示

## 📚 API仕様

### 認証エンドポイント

#### ユーザー登録
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "name": "ユーザー名",
  "email": "user@example.com",
  "password": "password123"
}
```

#### ユーザーログイン
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### ユーザーエンドポイント

#### 現在のユーザー情報取得
```http
GET /api/v1/users/me
Authorization: Bearer <jwt-token>
```

#### ユーザー情報更新
```http
PUT /api/v1/users/me
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "name": "新しいユーザー名",
  "email": "newemail@example.com"
}
```

詳細なAPI仕様は [openapi.yml](./openapi.yml) を参照してください。

## 🧪 テスト

```bash
# 全テストの実行
go test ./...

# 特定のテストファイルの実行
go test ./usecase
go test ./interface/handler
go test ./interface/middleware
```

## 🔧 開発

### ホットリロード
Docker Composeを使用している場合、Airによるホットリロードが有効になっています。ソースコードを変更すると自動的にアプリケーションが再起動されます。

### データベースマイグレーション
アプリケーション起動時に自動的にマイグレーションが実行されます。

## 🐳 Docker

### イメージのビルド
```bash
docker build -t voice-link-backend .
```

### コンテナの実行
```bash
docker run -p 8080:8080 voice-link-backend
```

## 📝 今後の開発予定

### Phase 1: 基本機能の拡張
- [ ] パスワードリセット機能
- [ ] メール認証機能
- [ ] ユーザープロフィール機能
- [ ] ユーザー情報の取得・更新
- [ ] ユーザー情報の削除
- [ ] ユーザー情報の更新
- [ ] 2段階認証

### Phase 2: リアルタイム機能
- [ ] WebSocket接続の実装
- [ ] 音声ストリーミング機能
- [ ] リアルタイム翻訳API統合
