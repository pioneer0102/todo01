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

### 6. クライアントの方でCRUDの操作
```
make client-create
make client-list
make client-update
make client-delete
```
