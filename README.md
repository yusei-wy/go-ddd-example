# go_ddd_example-server

Go言語でDDD（Domain-Driven Design）、CQRS（Command Query Responsibility Segregation）、FSD（Feature Sliced Design）パターンを採用したWebアプリケーションのサンプルプロジェクトです。

## 📚 ドキュメント

### 🚀 クイックスタート
新規参加者は以下の順序でドキュメントを確認してください：

1. **[📋 開発ルール総合インデックス](./docs/index.md)** - 全ドキュメントへの案内
2. **[🏗️ アーキテクチャ設計書](./docs/architecture.md)** - プロジェクト全体設計
3. **[📖 コーディング規約](./docs/rules/coding-standards.md)** - 基本的な開発ルール
4. **[🔄 Gitワークフロー](./docs/rules/git-workflow.md)** - 開発フロー

### 📁 ドキュメント構成

#### 基本ルール・ガイドライン (`docs/rules/`)
- **[コーディング規約](./docs/rules/coding-standards.md)** - Go + DDD の基本ルール
- **[API設計ルール](./docs/rules/api-design-rules.md)** - REST API 設計標準
- **[エラーハンドリング](./docs/rules/error-handling.md)** - レイヤー別エラー処理
- **[テストガイドライン](./docs/rules/testing-guidelines.md)** - TDD実践方法
- **[セキュリティガイドライン](./docs/rules/security-guidelines.md)** - セキュリティ要件
- **[パフォーマンス規則](./docs/rules/performance-rules.md)** - 最適化指針
- **[データベース規約](./docs/rules/database-conventions.md)** - DB設計・実装
- **[ドキュメント規則](./docs/rules/documentation-rules.md)** - ドキュメント作成標準
- **[Gitワークフロー](./docs/rules/git-workflow.md)** - ブランチ戦略・リリース
- **[コードレビューチェックリスト](./docs/rules/review-checklist.md)** - レビュー基準

#### AI開発支援 (`docs/ai-development/`)
- **[GitHub Copilot利用指針](./docs/ai-development/copilot-instructions.md)** - AI活用ガイド
- **[プロンプトテンプレート](./docs/ai-development/prompt-templates.md)** - AI支援開発用テンプレート
- **[AI生成コードレビュー](./docs/ai-development/ai-review-checklist.md)** - AI生成コード品質管理

## Requirements

- Docker
- [golang-migrate](https://github.com/golang-migrate/migrate)

## Directories

This project is a sample web application created using Go, employing the DDD, CQRS, and Feature Sliced Design (FSD) patterns.

```
.
├── app
│   └── env
├── cmd
├── db
│   └── migrations
├── docker
│   └── app
├── feature
│   └── <domain name>
│       ├── domain
│       │   └── model
│       │       └── command
│       │       └── query
│       │       └── value_object
│       ├── infrastructure
│       │       └── in_memory_repository
│       │       └── psql_repository
│       └── usecase
├── server
├── share
│   ├── custom_error
│   ├── domain
│   │   └── model
│   │       └── value_object
│   └── usecase
└── tmp
```

## commands

```
| Command  | Description |
| --- | --- |
| `make dev` | run dev server |
| `make reset-db` | clean db & create db & migrate |
| `make lint` | run golangci-lint |
| `make format` | run golangci-lint --fix |
| `make test` | run go test ./... |
```

## TODO

### 完了済み
- [x] DDD architecture
- [x] Custom error handler
- [x] Migration system
- [x] 包括的開発ルール・ガイドライン策定
- [x] AI開発支援ドキュメント整備

### 今後の実装予定
- [ ] ロギング機能
- [ ] 依存性注入（Dependency Injection）
- [ ] OpenAPI による Request/Response 型自動生成
- [ ] メトリクス収集機能
- [ ] 分散トレーシング
- [ ] CI/CDパイプライン

## 🤝 開発参加

このプロジェクトへの貢献を歓迎します：

1. **[開発ルール総合インデックス](./docs/index.md)** で開発方針を確認
2. **[Gitワークフロー](./docs/rules/git-workflow.md)** に従ってブランチ作成
3. **[コーディング規約](./docs/rules/coding-standards.md)** に準拠して実装
4. **[テストガイドライン](./docs/rules/testing-guidelines.md)** でTDD実践
5. **[コードレビューチェックリスト](./docs/rules/review-checklist.md)** でセルフチェック
6. Pull Request を作成

## 📞 サポート

- **ドキュメント**: [総合インデックス](./docs/index.md) を参照
- **質問・議論**: GitHub Issues を活用
- **改善提案**: GitHub Issues または Pull Request で提案
