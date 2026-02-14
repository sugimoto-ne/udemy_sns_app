# 開発環境用のDockerfile
FROM golang:alpine

# 必要なパッケージをインストール
RUN apk add --no-cache git

# 作業ディレクトリを設定
WORKDIR /app

# ポートを公開
EXPOSE 8080

# デフォルトコマンド（docker-compose.ymlで上書き）
CMD ["go", "run", "cmd/server/main.go"]
