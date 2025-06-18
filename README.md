# TODOアプリケーション（Go + connect-go + MySQL）

## 概要
Go, connect-go, MySQL, sqlboiler, slog, buf などを用いたgRPCベースのTODO管理アプリです。

## ディレクトリ構成
```
todo01/
├── cmd/            # エントリーポイント（server, client）
├── internal/       # DB, handler, repository, models
├── proto/          # Protocol Buffers定義
├── migrations/     # DBマイグレーション
├── gen/            # 生成コード（sqlboiler, protobuf）
├── Makefile        # ビルド・開発コマンド
├── Dockerfile      # Dockerビルド
├── docker-compose.yml # サービス一括起動
└── README.md       # このファイル
```

## セットアップ手順

### 1. 必要ツールのインストール
- Go（最新版推奨）
- Docker / Docker Compose
- [sqlboiler](https://github.com/volatiletech/sqlboiler) & ドライバ
- [buf](https://buf.build/)

#### sqlboiler, bufのインストール
```
make setup
```

### 2. MySQLサーバーの起動
```
make docker-mysql
```

### 3. 生成コードの作成
- Protobuf/Connect, sqlboilerモデル生成
```
make generate
```

### 4. ビルド
```
make build
```

### 5. サーバー起動
```
make start-server
```
実行結果列:
```
./bin/server.exe
{"time":"2025-06-18T20:25:10.3507316+09:00","level":"INFO","msg":"Successfully connected to database","host":"localhost","port":"3306","database":"todo"}
{"time":"2025-06-18T20:25:10.3512634+09:00","level":"INFO","msg":"Starting server","addr":":8080"}
```

クライアントでのCRUD操作を行うため、別のPowerShellターミナルを開いてください。

### 6. クライアントでのCRUD操作
```
make client-create
make client-list
make client-update
make client-delete
```
