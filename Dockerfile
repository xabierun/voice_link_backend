FROM golang:1.25-alpine

WORKDIR /app

# 必要なパッケージのインストール
RUN apk add --no-cache gcc musl-dev

# Airのインストール
RUN go install github.com/air-verse/air@latest

# アプリケーションの依存関係をコピー
COPY go.mod go.sum ./
RUN go mod download

# アプリケーションのソースコードをコピー
COPY . .

# ポートの公開
EXPOSE 8080

# Airを使用してアプリケーションを起動
CMD ["air"]
