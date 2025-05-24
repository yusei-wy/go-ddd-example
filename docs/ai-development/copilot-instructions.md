# GitHub Copilot 開発指示書

## 概要

このドキュメントは、GitHub Copilot を使用した Go DDD プロジェクトの開発における AI アシスタントへの指示とガイドラインです。

## 📚 関連ドキュメント参照

**重要**: 詳細な実装ガイドは以下のドキュメントを参照してください：
- **[コーディング規約](../rules/coding-standards.md)** - 命名規則・スタイル詳細
- **[エラーハンドリング](../rules/error-handling.md)** - レイヤー別エラー処理詳細
- **[API設計ルール](../rules/api-design-rules.md)** - REST API設計標準
- **[テストガイドライン](../rules/testing-guidelines.md)** - TDD実践・テスト戦略詳細
- **[セキュリティガイドライン](../rules/security-guidelines.md)** - セキュリティ実装詳細
- **[パフォーマンス規則](../rules/performance-rules.md)** - 最適化・ベンチマーク詳細
- **[データベース規約](../rules/database-conventions.md)** - DB設計・マイグレーション詳細

## ⚡ Copilot クイックリファレンス

### 即座に適用すべき重要ルール

#### 1. 命名規則（即座に適用・必須）
```go
// ✅ 必須: Go標準準拠
func CreateUser(userName string) error {}
type UserRepository interface {}
var userService UserService

// ❌ 禁止: アンダースコア使用
func create_user(user_name string) error {}
type User_Repository interface {}
var user_service UserService

// パッケージ名: 全て小文字、アンダースコアなし
package customerror  // ✅
package custom_error // ❌
```

#### 2. エラーハンドリング（即座に適用・必須）
```go
// ✅ 必須: レイヤー別カスタムエラー
return NewHandlerError(ErrInvalidRequest, "invalid request", 400)
return NewUsecaseError(ErrBusinessRule, "email exists")
return NewModelError(ErrInvalidEmail, "invalid format")
return NewRepositoryError(ErrDatabaseQuery, "query failed")

// ✅ 必須: エラーラップ
if err := someOperation(); err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

#### 3. セキュリティ（即座に適用・必須）
```go
// ✅ 必須: PreparedStatement使用
query := "SELECT * FROM users WHERE email = $1"
row := r.db.QueryRowContext(ctx, query, email.String())

// ❌ 禁止: 文字列結合（SQLインジェクション脆弱性）
query := "SELECT * FROM users WHERE email = '" + email + "'"

// ✅ 必須: 入力値検証
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=1,max=100"`
    Email string `json:"email" validate:"required,email"`
}
```

#### 4. DDD構造（即座に適用・必須）
```go
// ✅ 必須: プライベート不変ID
type User struct {
    id    UserID    // プライベート不変
    name  UserName  // 値オブジェクト
    email Email     // 値オブジェクト
}

// ✅ 必須: Getterメソッド
func (u *User) ID() UserID { return u.id }

// ✅ 必須: バリデーション付きコンストラクタ
func NewEmail(value string) (Email, error) {
    if !emailRegex.MatchString(value) {
        return Email{}, NewModelError(ErrInvalidEmail, "invalid format")
    }
    return Email{value: value}, nil
}
```

## プロジェクト固有の指示

### アーキテクチャ原則

GitHub Copilot は以下のアーキテクチャパターンに従ってコードを生成してください：

1. **Domain-Driven Design (DDD)**: ビジネスロジックをドメイン層に集約
2. **CQRS (Command Query Responsibility Segregation)**: 読み取りと書き込みの責務分離
3. **Feature Sliced Design (FSD)**: 機能ごとのディレクトリ分割
4. **レイヤードアーキテクチャ**: 明確な層分離と依存関係制御

### ディレクトリ構造

コード生成時は以下の構造に従ってください：

```
feature/<domain_name>/
├── domain/                 # ビジネスロジック層
│   ├── model/
│   │   ├── command/       # 書き込み用モデル
│   │   ├── query/         # 読み取り用モデル
│   │   └── value_object/  # 値オブジェクト
│   ├── repository.go      # リポジトリインターフェース
│   └── service.go         # ドメインサービス
├── infra/                 # インフラ層
│   ├── model.go          # データベースモデル
│   └── psql_repository.go # リポジトリ実装
├── usecase/              # アプリケーション層
└── handler.go            # プレゼンテーション層
```

## 命名規則の指示

### Go言語公式ガイドライン準拠

**重要**: Go言語では MixedCaps または mixedCaps を使用し、アンダースコアは避けてください。

```go
// ✅ Copilot が生成すべきコード
func CreateUser(userName string) error { }
type UserRepository interface { }
var userService UserService

// ❌ Copilot が生成してはいけないコード
func create_user(user_name string) error { }
type User_Repository interface { }
var user_service UserService
```

### パッケージ名規則

```go
// ✅ 推奨 - 全て小文字、アンダースコアなし
package customerror
package valueobject
package userservice

// ❌ 非推奨
package custom_error
package valueObject
package UserService
```

### ドメイン固有命名

```go
// エンティティ
type User struct { }
type Order struct { }

// 値オブジェクト
type Email struct { }
type UserId struct { }

// コマンド・クエリ
type CreateUserCommand struct { }
type UserFindQuery struct { }

// エラー
var ErrUserNotFound = errors.New("user not found")
```

## コード生成パターン

### 1. エンティティ生成パターン

```go
// GitHub Copilot 生成テンプレート
type <EntityName> struct {
    id    <EntityName>ID    // プライベート不変ID
    // 他のフィールドは値オブジェクト
}

// コンストラクタ
func New<EntityName>(/* パラメータ */) (*<EntityName>, error) {
    // バリデーション
    // エンティティ生成
    return &<EntityName>{
        id: New<EntityName>ID(),
        // 初期化
    }, nil
}

// Getter メソッド
func (e *<EntityName>) ID() <EntityName>ID {
    return e.id
}

// ビジネスメソッド
func (e *<EntityName>) <BusinessMethod>(/* パラメータ */) error {
    // ビジネスルール検証
    // 状態変更
    return nil
}
```

### 2. 値オブジェクト生成パターン

```go
// GitHub Copilot 生成テンプレート
type <ValueObjectName> struct {
    value string // プライベート、不変
}

func New<ValueObjectName>(value string) (<ValueObjectName>, error) {
    // バリデーション
    if /* バリデーション条件 */ {
        return <ValueObjectName>{}, NewModelError(Err<Error>, "error message")
    }

    return <ValueObjectName>{value: value}, nil
}

func (vo <ValueObjectName>) String() string {
    return vo.value
}

func (vo <ValueObjectName>) Equals(other <ValueObjectName>) bool {
    return vo.value == other.value
}
```

### 3. リポジトリ実装パターン

```go
// インターフェース
type <EntityName>Repository interface {
    Create(ctx context.Context, entity *<EntityName>) error
    FindByID(ctx context.Context, id <EntityName>ID) (*<EntityName>, error)
    Update(ctx context.Context, entity *<EntityName>) error
    Delete(ctx context.Context, id <EntityName>ID) error
}

// 実装
type psql<EntityName>Repository struct {
    db *sql.DB
}

func NewPsql<EntityName>Repository(db *sql.DB) <EntityName>Repository {
    return &psql<EntityName>Repository{db: db}
}

func (r *psql<EntityName>Repository) Create(ctx context.Context, entity *<EntityName>) error {
    query := "INSERT INTO <table_name> (id, ...) VALUES ($1, ...)"
    _, err := r.db.ExecContext(ctx, query, entity.ID(), /* 他のフィールド */)
    if err != nil {
        return fmt.Errorf("failed to create <entity_name>: %w", err)
    }
    return nil
}
```

### 4. ユースケース生成パターン

```go
// GitHub Copilot 生成テンプレート
type <Action><EntityName>Usecase struct {
    <entityName>Repo <EntityName>Repository
    logger        Logger
}

type <Action><EntityName>Input struct {
    // 入力フィールド
}

type <Action><EntityName>Output struct {
    // 出力フィールド
}

func (u *<Action><EntityName>Usecase) Execute(ctx context.Context, input <Action><EntityName>Input) (*<Action><EntityName>Output, error) {
    // 1. 入力バリデーション
    if err := u.validateInput(input); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // 2. ビジネスロジック実行
    // 3. 永続化
    // 4. 結果返却

    return &<Action><EntityName>Output{
        // 出力設定
    }, nil
}
```

### 5. ハンドラー生成パターン

```go
// GitHub Copilot 生成テンプレート
type <EntityName>Handler struct {
    <action><EntityName>Usecase *<Action><EntityName>Usecase
}

type <Action><EntityName>Request struct {
    // リクエストフィールド (JSON/バリデーションタグ付き)
}

type <Action><EntityName>Response struct {
    // レスポンスフィールド (JSON タグ付き)
}

func (h *<EntityName>Handler) <Action><EntityName>(c echo.Context) error {
    var req <Action><EntityName>Request
    if err := c.Bind(&req); err != nil {
        return NewHandlerError(ErrInvalidRequest, "invalid request format")
    }

    input := <Action><EntityName>Input{
        // リクエストから入力への変換
    }

    output, err := h.<action><EntityName>Usecase.Execute(c.Request().Context(), input)
    if err != nil {
        return err
    }

    response := <Action><EntityName>Response{
        // 出力からレスポンスへの変換
    }

    return c.JSON(http.StatusOK, response)
}
```

## エラーハンドリングの指示

### カスタムエラー使用

```go
// GitHub Copilot はレイヤー別のカスタムエラーを使用してください

// ハンドラー層
return NewHandlerError(ErrInvalidRequest, "request validation failed")

// ユースケース層
return NewUsecaseError(ErrBusinessRuleViolation, "email already exists")

// ドメイン層
return NewModelError(ErrInvalidEmail, "invalid email format")

// インフラ層
return NewRepositoryError(ErrDatabaseConnection, "failed to connect to database")
```

### エラーラップ

```go
// GitHub Copilot は常にエラーをラップしてください
if err := someOperation(); err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

## テスト生成の指示

### テスト命名規則

```go
// GitHub Copilot が生成すべきテスト名パターン
func Test<MethodName>_<Condition>_<ExpectedResult>(t *testing.T) {
    // テスト実装
}

// 例
func TestCreateUser_ValidInput_ReturnsUser(t *testing.T) { }
func TestCreateUser_DuplicateEmail_ReturnsError(t *testing.T) { }
func TestNewEmail_InvalidFormat_ReturnsError(t *testing.T) { }
```

### AAA パターン

```go
// GitHub Copilot は AAA パターンでテストを生成してください
func TestSomeMethod_Condition_Result(t *testing.T) {
    // Arrange（準備）
    // Given

    // Act（実行）
    // When

    // Assert（検証）
    // Then
}
```

### テストデータ生成

```go
// GitHub Copilot はテストヘルパー関数を活用してください
func createTestUser() *User {
    id := NewUserID()
    name, _ := NewUserName("Test User")
    email, _ := NewEmail("test@example.com")
    return &User{id: id, name: name, email: email}
}
```

## セキュリティ要件

### SQLインジェクション対策

```go
// ✅ GitHub Copilot が生成すべきコード - プリペアドステートメント使用
func (r *UserRepository) FindByEmail(ctx context.Context, email Email) (*User, error) {
    query := "SELECT id, name, email FROM users WHERE email = $1"
    row := r.db.QueryRowContext(ctx, query, email.String())
    // ...
}

// ❌ GitHub Copilot が生成してはいけないコード - 文字列結合
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
    query := "SELECT id, name, email FROM users WHERE email = '" + email + "'"
    row := r.db.QueryRowContext(ctx, query)
    // ...
}
```

### 入力検証

```go
// GitHub Copilot は階層的なバリデーションを実装してください

// 1. 構造的バリデーション（ハンドラー層）
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=1,max=100"`
    Email string `json:"email" validate:"required,email"`
}

// 2. ビジネスルールバリデーション（ドメイン層）
func NewEmail(value string) (Email, error) {
    if !emailRegex.MatchString(value) {
        return Email{}, NewModelError(ErrInvalidEmail, "invalid email format")
    }
    return Email{value: value}, nil
}
```

## パフォーマンス要件

### データベースアクセス最適化

```go
// GitHub Copilot は以下のパフォーマンスベストプラクティスに従ってください

// 1. バッチ処理
func (r *UserRepository) CreateBatch(ctx context.Context, users []*User) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    stmt, err := tx.PrepareContext(ctx, "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, user := range users {
        _, err := stmt.ExecContext(ctx, user.ID(), user.Name(), user.Email())
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

// 2. スライス容量事前確保
func processUsers(users []User) []ProcessedUser {
    result := make([]ProcessedUser, 0, len(users)) // 容量を事前確保
    for _, user := range users {
        result = append(result, processUser(user))
    }
    return result
}
```

## 禁止事項

### GitHub Copilot が生成してはいけないコード

1. **Global変数の使用**
   ```go
   // ❌ 禁止
   var globalDB *sql.DB
   var globalLogger Logger
   ```

2. **パニックの使用**（ライブラリ内部以外）
   ```go
   // ❌ 禁止
   func CreateUser(name string) User {
       if name == "" {
           panic("name is required")
       }
   }
   ```

3. **アンダースコア変数名**（Go言語コード内）
   ```go
   // ❌ 禁止
   func create_user(user_name string) error { }
   var user_repository UserRepository
   ```

4. **ドメイン層での外部依存**
   ```go
   // ❌ 禁止: ドメインサービスでDB直接アクセス
   func (s *UserService) CreateUser(name string) error {
       db, _ := sql.Open("postgres", "...")  // 禁止
   }
   ```

5. **文字列結合によるSQL**
   ```go
   // ❌ 禁止
   query := "SELECT * FROM users WHERE name = '" + name + "'"
   ```

## コンテキスト情報の活用

### プロジェクト情報

GitHub Copilot は以下のプロジェクト情報を考慮してコードを生成してください：

- **言語**: Go 1.21以上
- **フレームワーク**: Echo (Web), GORM (ORM), Testify (Testing)
- **データベース**: PostgreSQL
- **アーキテクチャ**: DDD + CQRS + FSD

### 既存コードパターンの学習

プロジェクト内の既存コードから以下のパターンを学習し、一貫性を保ってください：

1. **エラーハンドリングパターン**: カスタムエラー型の使用方法
2. **テストパターン**: AAA パターンとテストヘルパーの使用
3. **ログパターン**: 構造化ログの出力方法
4. **バリデーションパターン**: 階層的バリデーションの実装

### ファイル間の依存関係

GitHub Copilot は以下の依存関係制約を守ってコードを生成してください：

```
Handler → Usecase → Domain Service → Repository Interface
   ↓         ↓            ↓              ↓
   ↓         ↓            ↓        Infrastructure
   ↓         ↓            ↓              ↓
   ↓         ↓      Domain Model   Repository Impl
   ↓         ↓            ↓
   ↓    Application      ↓
   ↓         ↓            ↓
Presentation    Value Objects
```

## AI 開発フロー

### 1. コード生成前の確認事項

GitHub Copilot は以下を確認してからコードを生成してください：

- [ ] 生成するコンポーネントの層（ドメイン/アプリケーション/インフラ/プレゼンテーション）
- [ ] 既存の類似コンポーネントのパターン
- [ ] 必要なインポート文
- [ ] エラーハンドリング要件

### 2. 生成コードの品質確保

GitHub Copilot が生成するコードは以下の品質を満たしてください：

- [ ] 適切なエラーハンドリング
- [ ] ゴルーチンセーフティ（必要に応じて）
- [ ] リソース管理（defer によるクリーンアップ）
- [ ] 適切なログ出力
- [ ] バリデーション実装

### 3. テスト自動生成

GitHub Copilot は実装コードと同時にテストコードも生成してください：

- [ ] ユニットテスト（正常系・異常系）
- [ ] モックオブジェクトの使用
- [ ] テストヘルパー関数の活用
- [ ] エッジケースのテスト

## 設定ファイル連携

### VS Code 設定

GitHub Copilot は以下の VS Code 設定を参考にしてください：

```json
{
    "github.copilot.enable": {
        "*": true,
        "go": true
    },
    "github.copilot.advanced": {
        "debug.overrideEngine": "codex",
        "debug.testOverrideProxyUrl": "",
        "debug.overrideProxyUrl": ""
    }
}
```

### Go 言語固有設定

```json
{
    "go.lintTool": "golangci-lint",
    "go.formatTool": "goimports",
    "go.useLanguageServer": true,
    "go.buildOnSave": "workspace",
    "go.testOnSave": true
}
```

## 継続的改善

### フィードバックループ

GitHub Copilot の提案に対して以下の観点でフィードバックを行ってください：

1. **アーキテクチャ準拠**: DDD パターンの正しい実装
2. **命名規則**: Go言語標準の遵守
3. **セキュリティ**: セキュアなコード実装
4. **パフォーマンス**: 効率的なコード生成
5. **テスタビリティ**: テストしやすい設計

### 学習データ更新

プロジェクトの進化に合わせて、この指示書も更新してください：

- [ ] 新しいパターンの追加
- [ ] 非推奨パターンの削除
- [ ] セキュリティ要件の更新
- [ ] パフォーマンス最適化の追加

## 最終更新

- バージョン: 2025.05.24
- 更新者: AI Assistant
