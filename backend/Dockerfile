# 開発環境用のDockerfile（Air使用）
FROM golang:1.21-alpine

# 必要なパッケージをインストール
RUN apk add --no-cache git

# Airをインストール（ホットリロードツール）
# v1.49.0はGo 1.21で動作する安定版
RUN go install github.com/cosmtrek/air@v1.49.0

# GOPATHのbinディレクトリをPATHに追加
ENV PATH="/go/bin:${PATH}"

# GOTOOLCHAINを設定（必要に応じて新しいGoバージョンを自動ダウンロード）
ENV GOTOOLCHAIN=auto

# 作業ディレクトリを設定
WORKDIR /app

# ポートを公開
EXPOSE 8080

# デフォルトコマンド（docker-compose.ymlで上書き）
CMD ["air", "-c", ".air.toml"]
