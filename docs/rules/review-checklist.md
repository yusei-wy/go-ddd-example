# コードレビューチェックリスト

## 概要

このドキュメントは、Go DDD プロジェクトにおけるコードレビューの品質と一貫性を確保するためのチェックリストです。

## レビュー前チェック（PR作成者）

### 基本チェック

- [ ] **テスト実行**: すべてのテストが通ることを確認
  ```bash
  make test
  make test-integration
  ```

- [ ] **リンターチェック**: コード品質チェックをパス
  ```bash
  make lint
  ```

- [ ] **フォーマット**: コードフォーマットが適用済み
  ```bash
  make format
  ```

- [ ] **カバレッジ**: テストカバレッジが基準を満たしている
  ```bash
  make test-coverage
  # ドメイン層: 90%以上、その他: 80%以上
  ```

### コミット品質

- [ ] **コミットメッセージ**: Conventional Commits 規約に準拠
- [ ] **コミット単位**: 論理的な単位でコミットが分割されている
- [ ] **履歴**: 不要なマージコミットがない（リベース実行済み）

### PR説明

- [ ] **変更概要**: 何を、なぜ変更したかが明確
- [ ] **影響範囲**: 変更による影響範囲の説明
- [ ] **テスト計画**: テストケースと確認項目の記載
- [ ] **関連Issue**: 関連するIssueへのリンク

## アーキテクチャレビュー

### DDD パターン遵守

#### ドメイン層

- [ ] **エンティティ**: 識別子を持ち、ビジネス不変条件を満たしている
  ```go
  // ✅ 適切な例
  type User struct {
      id    UserID    // 不変の識別子
      name  UserName  // 値オブジェクト
      email Email     // 値オブジェクト
  }
  
  func (u *User) ChangeName(newName UserName) error {
      // ビジネスルールによる検証
      if err := u.validateNameChange(newName); err != nil {
          return err
      }
      u.name = newName
      return nil
  }
  ```

- [ ] **値オブジェクト**: 不変性とバリデーションが適切に実装されている
  ```go
  // ✅ 適切な例
  type Email struct {
      value string
  }
  
  func NewEmail(value string) (Email, error) {
      if !emailRegex.MatchString(value) {
          return Email{}, ErrInvalidEmail
      }
      return Email{value: strings.ToLower(value)}, nil
  }
  
  func (e Email) Equals(other Email) bool {
      return e.value == other.value
  }
  ```

- [ ] **ドメインサービス**: エンティティや値オブジェクトに属さないビジネスロジック
  ```go
  // ✅ 適切な例
  type UserDomainService struct {
      userRepo UserRepository
  }
  
  func (s *UserDomainService) IsDuplicateEmail(ctx context.Context, email Email) bool {
      return s.userRepo.ExistsByEmail(ctx, email)
  }
  ```

#### 依存関係

- [ ] **依存方向**: 上位層が下位層に依存している（循環依存なし）
- [ ] **インターフェース分離**: 必要最小限のメソッドのみを含む
- [ ] **外部依存**: ドメイン層が外部ライブラリに直接依存していない

### CQRS パターン

- [ ] **コマンド**: 書き込み用モデルが適切に分離されている
- [ ] **クエリ**: 読み取り用モデルが適切に分離されている
- [ ] **責務分離**: 読み取りと書き込みのユースケースが分離されている

### レイヤー責務

- [ ] **プレゼンテーション層**: HTTPリクエスト/レスポンス処理のみ
- [ ] **アプリケーション層**: ユースケースの制御と調整のみ
- [ ] **ドメイン層**: ビジネスロジックのみ
- [ ] **インフラ層**: 外部システムとの通信のみ

## コード品質レビュー

### 命名規則

- [ ] **Go言語規約**: MixedCaps/mixedCaps を使用（アンダースコア避ける）
  ```go
  // ✅ 適切
  func CreateUser(userName string) error
  var userRepository UserRepository
  
  // ❌ 不適切
  func create_user(user_name string) error
  var user_repository UserRepository
  ```

- [ ] **パッケージ名**: 小文字のみ、アンダースコアなし
  ```go
  // ✅ 適切
  package customerror
  package valueobject
  
  // ❌ 不適切
  package custom_error
  package valueObject
  ```

- [ ] **意味のある名前**: 機能や目的が名前から理解できる
- [ ] **一貫性**: プロジェクト全体で命名パターンが統一されている

### 関数・メソッド設計

- [ ] **単一責任**: 一つの関数が一つの責務を持つ
- [ ] **適切な長さ**: 関数は適切な長さ（~20行が目安）
- [ ] **引数数**: 引数は適切な数（4個以下が目安）
- [ ] **戻り値**: エラーハンドリングが適切

```go
// ✅ 適切な例
func (s *UserService) CreateUser(ctx context.Context, cmd CreateUserCommand) (*User, error) {
    if err := s.validateCommand(cmd); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    user := s.buildUser(cmd)
    if err := s.userRepo.Save(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to save user: %w", err)
    }
    
    return user, nil
}
```

### エラーハンドリング

- [ ] **適切なラップ**: エラーの原因と文脈が保持されている
  ```go
  // ✅ 適切
  if err := userRepo.Create(ctx, user); err != nil {
      return fmt.Errorf("failed to create user: %w", err)
  }
  ```

- [ ] **カスタムエラー**: ビジネスエラーは適切なカスタムエラー型を使用
- [ ] **エラーレベル**: 適切なレイヤーでエラーがハンドリングされている

### メモリ・パフォーマンス

- [ ] **メモリ効率**: 不要なメモリ使用がない
- [ ] **ゴルーチンリーク**: ゴルーチンのリークがない
- [ ] **リソース管理**: ファイル、DB接続等の適切なクローズ

```go
// ✅ 適切な例
func (r *UserRepository) FindUsers(ctx context.Context) ([]*User, error) {
    rows, err := r.db.QueryContext(ctx, "SELECT ...")
    if err != nil {
        return nil, err
    }
    defer rows.Close() // 必須: リソースのクローズ
    
    users := make([]*User, 0, 100) // 容量を事前確保
    for rows.Next() {
        user := &User{}
        if err := rows.Scan(&user.ID, &user.Name); err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}
```

## セキュリティレビュー

### 入力検証

- [ ] **バリデーション**: すべての外部入力が適切に検証されている
- [ ] **サニタイゼーション**: XSS対策のためのHTMLエスケープ
- [ ] **SQLインジェクション**: プリペアドステートメントの使用

```go
// ✅ セキュアな例
func (r *UserRepository) FindByEmail(ctx context.Context, email Email) (*User, error) {
    query := "SELECT id, name, email FROM users WHERE email = $1"
    row := r.db.QueryRowContext(ctx, query, email.String())
    // ...
}

// ❌ 脆弱な例
func (r *UserRepository) FindByEmailUnsafe(ctx context.Context, email string) (*User, error) {
    query := "SELECT id, name, email FROM users WHERE email = '" + email + "'"
    row := r.db.QueryRowContext(ctx, query) // SQLインジェクション脆弱性
    // ...
}
```

### 認証・認可

- [ ] **認証チェック**: 適切な認証チェックが実装されている
- [ ] **認可制御**: 適切な権限チェックが実装されている
- [ ] **トークン管理**: JWTトークンの適切な検証

### データ保護

- [ ] **機密データ**: パスワード等の機密データが平文で保存されていない
- [ ] **ログ出力**: 機密情報がログに出力されていない
- [ ] **暗号化**: 必要に応じて暗号化が実装されている

## テストレビュー

### テストカバレッジ

- [ ] **カバレッジ**: 適切なテストカバレッジが確保されている
- [ ] **境界値**: 境界値のテストケースが含まれている
- [ ] **エラーケース**: 異常系のテストケースが含まれている

### テスト品質

- [ ] **テスト名**: テストの意図が名前から理解できる
  ```go
  // ✅ 適切
  func TestCreateUser_DuplicateEmail_ReturnsError(t *testing.T)
  
  // ❌ 不適切
  func TestCreateUser2(t *testing.T)
  ```

- [ ] **AAA パターン**: Arrange-Act-Assert の構造が明確
- [ ] **独立性**: テストケース間の依存関係がない
- [ ] **モック使用**: 外部依存が適切にモック化されている

### テストデータ

- [ ] **テストデータ**: 意味のあるテストデータを使用
- [ ] **データクリーンアップ**: テスト後のデータクリーンアップが実装されている
- [ ] **並行実行**: テストが並行実行に対応している

## パフォーマンスレビュー

### データベースアクセス

- [ ] **N+1問題**: N+1問題の回避策が実装されている
- [ ] **インデックス**: 必要なインデックスが設計されている
- [ ] **トランザクション**: 適切なトランザクション境界が設定されている

```go
// ✅ N+1問題を回避する例
func (r *OrderRepository) FindOrdersWithItems(ctx context.Context) ([]*Order, error) {
    query := `
        SELECT o.id, o.total, oi.id, oi.product_id, oi.quantity
        FROM orders o
        LEFT JOIN order_items oi ON o.id = oi.order_id
    `
    // JOINを使用して一度のクエリで取得
}
```

### 同期処理

- [ ] **ゴルーチン**: 適切なゴルーチンの使用
- [ ] **チャネル**: チャネルの適切な使用とクローズ
- [ ] **ミューテックス**: 必要に応じた排他制御

### キャッシュ戦略

- [ ] **キャッシュ**: 適切なキャッシュ戦略が実装されている
- [ ] **TTL**: 適切なTTL設定
- [ ] **キャッシュ無効化**: 適切なキャッシュ無効化戦略

## ドキュメントレビュー

### コードコメント

- [ ] **必要なコメント**: 複雑なロジックに適切なコメント
- [ ] **関数コメント**: パブリック関数にGoDocコメント
- [ ] **最新性**: コメントとコードの整合性

```go
// ✅ 適切なコメント例
// CreateUser creates a new user with the provided information.
// It validates the user data and ensures email uniqueness before creation.
//
// Parameters:
//   - ctx: context for cancellation and timeout
//   - cmd: command containing user creation data
//
// Returns:
//   - *User: created user entity with generated ID
//   - error: validation error or creation failure
func (s *UserService) CreateUser(ctx context.Context, cmd CreateUserCommand) (*User, error) {
    // ビジネスルールバリデーション
    if err := s.validateUserCreation(cmd); err != nil {
        return nil, fmt.Errorf("user validation failed: %w", err)
    }
    // ... 実装
}
```

### API ドキュメント

- [ ] **OpenAPI**: APIの変更がOpenAPI仕様に反映されている
- [ ] **例**: リクエスト/レスポンスの例が含まれている
- [ ] **エラーコード**: エラーレスポンスが文書化されている

## レビューフィードバックガイド

### フィードバックの種類

#### 必須修正（Must Fix）
```
🔴 Must Fix: SQLインジェクション脆弱性があります
文字列結合ではなくプリペアドステートメントを使用してください。

// 現在のコード
query := "SELECT * FROM users WHERE email = '" + email + "'"

// 推奨修正
query := "SELECT * FROM users WHERE email = $1"
row := db.QueryRowContext(ctx, query, email)
```

#### 推奨修正（Should Fix）
```
🟡 Should Fix: エラーハンドリングの改善
エラーの文脈情報を保持するためにfmt.Errorfでラップしてください。

// 現在のコード
return err

// 推奨修正
return fmt.Errorf("failed to create user: %w", err)
```

#### 提案（Consider）
```
🟢 Consider: パフォーマンス最適化の提案
大量データの処理では、バッチ処理の実装を検討してください。

// 提案するアプローチ
func (r *Repository) CreateUsersBatch(ctx context.Context, users []*User) error {
    // バッチ挿入の実装
}
```

#### 質問（Question）
```
❓ Question: 設計判断について
この値オブジェクトでEqualsメソッドが必要な理由を教えてください。
現在の実装では使用されていないようですが、将来的な用途がありますか？
```

### フィードバックのベストプラクティス

- [ ] **具体的**: 具体的な改善案を提示
- [ ] **理由明記**: なぜ変更が必要かを説明
- [ ] **サンプル提供**: 修正例のコードを提示
- [ ] **建設的**: 批判ではなく改善提案として記載

## レビューツール設定

### GitHub PR テンプレート

```markdown
<!-- .github/pull_request_template.md -->
## 変更概要
<!-- この PR の概要を記載 -->

## チェックリスト

### コード品質
- [ ] テスト実行済み (`make test`)
- [ ] リンターチェック通過 (`make lint`)
- [ ] カバレッジ確認済み (`make test-coverage`)

### アーキテクチャ
- [ ] DDD パターン遵守
- [ ] レイヤー責務の分離
- [ ] 依存関係の方向性確認

### セキュリティ
- [ ] 入力検証実装済み
- [ ] SQLインジェクション対策済み
- [ ] 機密情報の適切な取り扱い

### パフォーマンス
- [ ] N+1問題の確認
- [ ] リソースの適切な管理
- [ ] 必要に応じたインデックス設計

## レビュー依頼事項
<!-- 特に注意してレビューしてほしい点を記載 -->
```

### コードオーナー設定

```
# .github/CODEOWNERS
# Global reviewers
* @tech-lead @senior-developer

# Domain specific
/feature/user/ @domain-expert @user-team-lead
/feature/order/ @order-team-lead
/docs/ @tech-writer @architect

# Infrastructure
/infra/ @infra-team @sre-team
/db/ @database-team
/.github/ @devops-team
```

## 最終更新

- バージョン: 2025.05.24
- 更新者: AI Assistant
