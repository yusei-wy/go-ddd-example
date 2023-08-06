# ddd_go_example

## ディレクトリ構成

```
.
├── cmd
│   └── env      // 環境ごとのエントリーポイント
│   └── main.go  // エントリーポイント
├── config
└── internal
    └── app
        ├── convert               // Input, Output の 変換を行う
        ├── domain                // domain 層
        │   ├── model
        │   │   ├── command_model // 副作用のある操作に関するモデルを生成する
        │   │   ├── query_model   // レスポンスに使用する型
        │   │   └── value_object  // 値オブジェクト
        │   ├── repository
        │   └── service
        ├── infrastructure
        │   ├── model             // DB から取得したデータをマッピングする型
        │   └── psql_repository   // PostgreSQL 用のリポジトリ（DynamoDB などに柔軟に変更できるように
        ├── interface             // HTTP リクエストを受け取り UseCase で処理を行ってレスポンスを返す
        │   ├── respose
        │   └── router
        └── usecase```

