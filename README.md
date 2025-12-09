# Job Detail API (Go)

簡易的な求人情報取得APIです。DynamoDBから複数のJob IDに対応する求人情報を一括取得します。

## アーキテクチャ

- **言語**: Go
- **データベース**: Amazon DynamoDB
- **ローカル開発環境**: LocalStack (Docker)

### ディレクトリ構成
- `cmd/main.go`: アプリケーションのエントリポイント。HTTPサーバーの起動とAWS SDKの設定を行います。LocalStack接続用のロジックもここに含まれます。
- `internal/job.go`: データモデル (`Job` 構造体) の定義。
- `internal/handler.go`: ビジネスロジック。DynamoDBへのアクセス (`BatchGetItem`) を担当します。
- `docker/`: LocalStack用のDocker構成ファイル。初期化スクリプトも含まれます。

## 起動方法

### 1. ローカル環境の起動
Dockerを使用してDynamoDB (LocalStack) を起動します。

```bash
cd docker
docker-compose up -d
cd ..
```

### 2. データの確認・登録
LocalStackが起動したら、以下のURLからデータの確認や登録が可能です。
[LocalStack Dashboard - Jobs Table](https://app.localstack.cloud/inst/default/resources/dynamodb/tables/Jobs/items)

※ 初回起動時は `docker/init/01_create_table.sh` により自動的にテーブル作成とテストデータ (`job_a`, `job_b`) の投入が行われます。

### 3. アプリケーションの起動
LocalStackのエンドポイントを指定してGoアプリケーションを起動します。

```bash
AWS_ENDPOINT_URL=http://localhost:4566 go run cmd/main.go
```

### 4. 動作確認
curlコマンドを使用してAPIを呼び出します。

```bash
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{"jobIds": ["job_a", "job_b"]}'
```

成功すると以下のようなJSONが返却されます。

```json
{"job":[{"jobId":"job_a","jobTitle":"Software Engineer","jobContent":"Develop awesome software."},{"jobId":"job_b","jobTitle":"Product Manager","jobContent":"Manage product roadmap."}]}
```
