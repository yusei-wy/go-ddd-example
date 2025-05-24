# Git ワークフロー規約

## 概要

このドキュメントは、Go DDD プロジェクトにおける Git ワークフローとブランチ管理の規約を定義します。

## ブランチ戦略

### GitFlow ベースのブランチ戦略

```
main (本番環境)
├── develop (開発統合)
│   ├── feature/user-management (機能開発)
│   ├── feature/order-processing (機能開発)
│   └── bugfix/user-validation-fix (バグ修正)
├── release/v1.2.0 (リリース準備)
└── hotfix/critical-security-fix (緊急修正)
```

### ブランチ命名規則

#### 1. 機能開発ブランチ

```bash
# パターン: feature/<domain>-<feature>
feature/user-registration
feature/order-payment
feature/notification-email

# 複数単語はハイフン区切り
feature/user-profile-management
feature/payment-method-validation
```

#### 2. バグ修正ブランチ

```bash
# パターン: bugfix/<domain>-<issue>
bugfix/user-validation-error
bugfix/order-total-calculation
bugfix/email-template-rendering
```

#### 3. リリースブランチ

```bash
# パターン: release/v<major>.<minor>.<patch>
release/v1.0.0
release/v1.2.0
release/v2.0.0-beta.1
```

#### 4. ホットフィックスブランチ

```bash
# パターン: hotfix/<critical-issue>
hotfix/security-vulnerability
hotfix/data-corruption-fix
hotfix/critical-performance-issue
```

## コミットメッセージ規約

### Conventional Commits 準拠

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### コミットタイプ

| タイプ | 説明 | 例 |
|--------|------|-----|
| `feat` | 新機能追加 | `feat(user): add user registration endpoint` |
| `fix` | バグ修正 | `fix(order): fix total calculation error` |
| `docs` | ドキュメント変更 | `docs(api): update user endpoint documentation` |
| `style` | コードフォーマット | `style(user): fix code formatting in user service` |
| `refactor` | リファクタリング | `refactor(domain): extract user validation logic` |
| `test` | テスト追加・修正 | `test(user): add unit tests for user repository` |
| `chore` | ビルドプロセス・補助ツール | `chore(deps): update dependencies` |
| `perf` | パフォーマンス改善 | `perf(db): optimize user query performance` |
| `ci` | CI設定変更 | `ci(github): update test workflow` |

### スコープ（DDD層別）

| スコープ | 説明 | 例 |
|----------|------|-----|
| `domain` | ドメイン層の変更 | `feat(domain): add email value object` |
| `usecase` | ユースケース層の変更 | `feat(usecase): add create user usecase` |
| `infra` | インフラ層の変更 | `feat(infra): add postgres user repository` |
| `handler` | プレゼンテーション層の変更 | `feat(handler): add user creation endpoint` |
| `config` | 設定関連の変更 | `chore(config): update database connection pool` |

### コミットメッセージ例

#### ✅ 推奨例

```bash
# 機能追加
feat(user): add user registration with email validation

Implements user registration feature including:
- Email validation using regex pattern
- Password strength validation
- Duplicate email checking

Closes #123

# バグ修正
fix(order): resolve total calculation rounding error

Fixed floating point precision issue in order total calculation.
Now using decimal library for monetary calculations.

Fixes #456

# リファクタリング
refactor(domain): extract user validation into value objects

- Move email validation to Email value object
- Move name validation to UserName value object
- Improve testability and reusability

# ドキュメント更新
docs(architecture): update DDD layer responsibility

Add detailed explanation of each layer's responsibility
and data flow between layers.

# テスト追加
test(user): add integration tests for user repository

Add comprehensive integration tests covering:
- User creation with valid data
- Duplicate email handling
- Database constraint validation
```

#### ❌ 避けるべき例

```bash
# あいまいな説明
fix: bug fix

# 詳細すぎるコミット
feat: add user name validation and email validation and password validation and...

# タイポやフォーマット違反
feat(usr): ad user registraton
Fix: user bug
FEAT(user): add user
```

## ブランチ運用ルール

### 1. main ブランチ

- **用途**: 本番環境デプロイ用
- **保護設定**: 直接プッシュ禁止
- **マージ方法**: Pull Request のみ
- **必須レビュー**: 2名以上

```bash
# main ブランチへの直接プッシュは禁止
# 必ず Pull Request を作成
git switch main
git pull origin main
git switch -c feature/new-feature
```

### 2. develop ブランチ

- **用途**: 開発版統合
- **マージ元**: feature, bugfix ブランチ
- **マージ先**: release ブランチ
- **必須レビュー**: 1名以上

```bash
# 機能開発の開始
git switch develop
git pull origin develop
git switch -c feature/user-authentication

# 開発完了後、develop へマージ
git switch develop
git merge feature/user-authentication
git push origin develop
```

### 3. feature ブランチ

- **ライフサイクル**: 機能完了まで
- **命名**: `feature/<domain>-<feature>`
- **ベースブランチ**: develop
- **マージ先**: develop

```bash
# 機能ブランチ作成
git switch develop
git pull origin develop
git switch -c feature/order-management

# 開発中のコミット
git add .
git commit -m "feat(order): add order creation logic"
git commit -m "test(order): add order service unit tests"
git commit -m "docs(order): add order API documentation"

# Pull Request 作成前の準備
git switch develop
git pull origin develop
git switch feature/order-management
git rebase develop  # 必要に応じて
```

## Pull Request 規約

### PR タイトル規約

```
[Type] Brief description of changes
```

#### 例

```
[Feature] Add user authentication with JWT
[Fix] Resolve order total calculation error
[Refactor] Extract email validation logic to value object
[Docs] Update API documentation for user endpoints
```

### PR 説明テンプレート

```markdown
## 概要
この PR で実装した内容を簡潔に説明してください。

## 変更内容
- [ ] 新機能追加
- [ ] バグ修正
- [ ] リファクタリング
- [ ] ドキュメント更新
- [ ] テスト追加

## 変更詳細
### 追加・変更されたファイル
- `feature/user/domain/model/user.go` - User エンティティの実装
- `feature/user/usecase/create_user.go` - ユーザー作成ユースケース
- `feature/user/handler.go` - ユーザー作成エンドポイント

### 削除されたファイル
- （該当する場合のみ記載）

## テスト
- [ ] ユニットテスト実行済み (`make test`)
- [ ] 統合テスト実行済み (`make test-integration`)
- [ ] 手動テスト実行済み
- [ ] カバレッジ確認済み (`make test-coverage`)

## 動作確認
### テスト手順
1. アプリケーション起動 (`make dev`)
2. POST /users で新規ユーザー作成
3. GET /users/{id} でユーザー情報取得確認

### 確認項目
- [ ] 正常系の動作確認
- [ ] 異常系の動作確認（バリデーションエラー等）
- [ ] エラーハンドリングの確認

## 関連 Issue
Closes #123
Relates to #456

## レビュー依頼項目
- [ ] アーキテクチャ設計の妥当性
- [ ] DDD パターンの適切な実装
- [ ] エラーハンドリングの適切性
- [ ] テストカバレッジの充分性
- [ ] セキュリティ観点での問題有無

## スクリーンショット（該当する場合）
（UI 変更がある場合のみ）

## その他
その他の注意事項や補足情報があれば記載してください。
```

### PR サイズガイドライン

#### Small PR（推奨）
- 変更行数: ~100行
- 変更ファイル数: ~5ファイル
- レビュー時間: ~30分

#### Medium PR
- 変更行数: ~300行
- 変更ファイル数: ~10ファイル
- レビュー時間: ~1時間

#### Large PR（分割推奨）
- 変更行数: 300行以上
- 変更ファイル数: 10ファイル以上
- レビュー時間: 1時間以上

### PR 作成前チェックリスト

```bash
# 1. テスト実行
make test
make test-integration
make lint

# 2. コードフォーマット
make format

# 3. 最新の develop との同期
git switch develop
git pull origin develop
git switch feature/your-branch
git rebase develop

# 4. コンフリクト解決（必要に応じて）
git add .
git rebase --continue

# 5. プッシュ
git push origin feature/your-branch
```

## リベース vs マージ戦略

### 機能ブランチでの作業

```bash
# ✅ 推奨: リベースを使用（履歴をきれいに保つ）
git switch feature/user-auth
git rebase develop

# ❌ 非推奨: マージコミットによる履歴の複雑化
git switch feature/user-auth
git merge develop
```

### メインブランチへのマージ

```bash
# ✅ 推奨: Squash and Merge（機能単位で1コミット）
# GitHub/GitLab の UI で "Squash and Merge" を使用

# 例: マージ後のコミットメッセージ
feat(user): add user authentication with JWT (#123)

* implement user registration
* add JWT token generation
* add authentication middleware
* add user login endpoint
```

## タグ運用

### セマンティックバージョニング

```bash
# パターン: v<major>.<minor>.<patch>
v1.0.0  # 初回リリース
v1.1.0  # 機能追加
v1.1.1  # バグ修正
v2.0.0  # 破壊的変更
```

### タグ作成手順

```bash
# 1. main ブランチで最新を取得
git switch main
git pull origin main

# 2. タグ作成（注釈付きタグ）
git tag -a v1.0.0 -m "Release version 1.0.0

Features:
- User management system
- Order processing
- Payment integration

Bug fixes:
- Fix user validation error
- Resolve order calculation issue"

# 3. タグをプッシュ
git push origin v1.0.0

# 4. すべてのタグをプッシュ
git push origin --tags
```

## Git フック活用

### Pre-commit フック

```bash
#!/bin/sh
# .git/hooks/pre-commit

# テスト実行
echo "Running tests..."
make test
if [ $? -ne 0 ]; then
    echo "Tests failed. Commit aborted."
    exit 1
fi

# Lint チェック
echo "Running lint..."
make lint
if [ $? -ne 0 ]; then
    echo "Lint failed. Commit aborted."
    exit 1
fi

# フォーマットチェック
echo "Checking format..."
make format-check
if [ $? -ne 0 ]; then
    echo "Code formatting required. Run 'make format' and commit again."
    exit 1
fi

echo "Pre-commit checks passed."
```

### Commit-msg フック

```bash
#!/bin/sh
# .git/hooks/commit-msg

# コミットメッセージフォーマットチェック
commit_regex='^(feat|fix|docs|style|refactor|test|chore|perf|ci)(\(.+\))?: .{1,50}'

if ! grep -qE "$commit_regex" "$1"; then
    echo "Invalid commit message format."
    echo "Format: <type>[optional scope]: <description>"
    echo "Example: feat(user): add user registration endpoint"
    exit 1
fi

echo "Commit message format is valid."
```

## トラブルシューティング

### よくある Git 問題と解決方法

#### 1. マージコンフリクト解決

```bash
# コンフリクト発生時
git status  # コンフリクトファイル確認

# 手動でコンフリクト解決後
git add conflicted-file.go
git rebase --continue

# リベース中止したい場合
git rebase --abort
```

#### 2. 間違ったコミットの修正

```bash
# 最後のコミットメッセージ修正
git commit --amend -m "correct message"

# 複数のコミットを修正（インタラクティブリベース）
git rebase -i HEAD~3

# 特定のコミットを取り消し
git revert <commit-hash>
```

#### 3. ブランチの強制更新

```bash
# ローカルブランチを origin と同期
git switch feature/branch-name
git reset --hard origin/feature/branch-name

# 危険: 強制プッシュ（十分注意して実行）
git push --force-with-lease origin feature/branch-name
```

## 最終更新

- バージョン: 2025.05.24
- 更新者: AI Assistant
