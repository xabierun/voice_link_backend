# Voice Link Backend API

リアルタイム翻訳機能を提供するアプリケーションのバックエンドAPIです。

## 機能

- ユーザー登録・認証（JWT認証）
- ユーザー情報の取得・更新
- PostgreSQLデータベース連携

## 技術スタック

- Go 1.24.3
- Echo v4
- PostgreSQL 16
- GORM
- Docker & Docker Compose

## プロジェクト構造

```
voice_link_backend/
├── domain/                    # ドメイン層
│   └── model/user.go         # ユーザーモデル
├── infrastructure/            # インフラストラクチャ層
│   └── persistence/          # データベース操作
├── interface/                 # インターフェース層
│   ├── handler/              # HTTPハンドラー
│   ├── middleware/           # ミドルウェア
│   └── router/               # ルーティング
├── usecase/                  # ユースケース層
├── main.go                   # エントリーポイント
├── openapi.yml               # API仕様書
└── docker-compose.yml        # Docker設定
```


## セットアップ

### 前提条件
- Go 1.24.3以上
- Docker & Docker Compose

### 起動方法

#### Docker Compose（推奨）
```bash
docker-compose up -d
```

#### ローカル開発
```bash
# 依存関係のインストール
go mod download

# データベースの起動
docker-compose up -d db

# アプリケーションの起動
go run main.go
```

#### Makefileを使用
```bash
make dev    # 開発サーバー起動
make test   # テスト実行
make help   # 利用可能なコマンド表示
```


## API仕様

### 認証
- `POST /api/v1/auth/register` - ユーザー登録
- `POST /api/v1/auth/login` - ログイン

### ユーザー
- `GET /api/v1/users/me` - 現在のユーザー情報取得
- `PUT /api/v1/users/me` - ユーザー情報更新

詳細なAPI仕様は [openapi.yml](./openapi.yml) を参照してください。

## テスト

```bash
go test ./...
```

## 開発

- ホットリロード: Docker Compose使用時に自動で再起動
- データベースマイグレーション: 起動時に自動実行

## 今後の開発予定

### Phase 1: 基本機能の拡張
- [ ] メール認証機能
- [ ] ユーザープロフィール機能
- [ ] 2段階認証

### Phase 2: リアルタイム機能
- [ ] WebSocket接続の実装
- [ ] 音声ストリーミング機能
- [ ] リアルタイム翻訳API統合
- [ ] Flutterアプリケーションとの連携
