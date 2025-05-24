# コーディング規約

## 概要

このドキュメントは、Go言語でのDDD（Domain-Driven Design）プロジェクトにおけるコーディング規約を定義します。

## 命名規則

### Go言語公式ガイドライン準拠

Go言語の公式ガイドライン「Effective Go」に従い、以下の規則を適用します：

> "the convention in Go is to use MixedCaps or mixedCaps rather than underscores to write multiword names."

#### 変数・関数・型名

- **パブリック**: `PascalCase` (例: `CreateUser`, `UserRepository`)
- **プライベート**: `camelCase` (例: `validateEmail`, `userRepo`)
- **定数**: `PascalCase` (例: `MaxUserNameLength`, `DefaultTimeout`)
- **インターフェース**: 動作を表す名詞 + `er` (例: `UserRepository`, `EmailSender`)

```go
// ✅ 推奨
func CreateUser(userName string) error { }
type UserRepository interface { }
const MaxRetryCount = 3

// ❌ 非推奨
func create_user(user_name string) error { }
type User_Repository interface { }
const MAX_RETRY_COUNT = 3
```

#### ファイル・ディレクトリ名

- **ファイル名**: `snake_case.go` (例: `user_service.go`, `value_object.go`)
- **ディレクトリ名**: `snake_case` (例: `value_object/`, `custom_error/`)
- **テストファイル**: `*_test.go` (例: `user_service_test.go`)

#### パッケージ名

Go言語公式推奨事項に従い、以下の規則を適用します：

- 全て小文字
- アンダースコアや大文字を使用しない
- 短く、簡潔な名前
- 複数の単語は連結

```go
// ✅ 推奨
package customerror
package valueobject
package userservice

// ❌ 非推奨
package custom_error
package valueObject
package UserService
```

### ドメイン固有命名

#### エンティティ・値オブジェクト

```go
// エンティティ: ドメインオブジェクト名
type User struct { }
type Order struct { }

// 値オブジェクト: 値の種類を表す名詞
type Email struct { }
type UserId struct { }
type Money struct { }
```

#### CQRS パターン

```go
// コマンド: <動詞><対象>Command
type CreateUserCommand struct { }
type UpdateOrderCommand struct { }

// クエリ: <対象><動詞>Query
type UserFindQuery struct { }
type OrderListQuery struct { }
```

#### エラー定義

```go
// エラー: Err<エラー内容>
var (
    ErrUserNotFound      = errors.New("user not found")
    ErrInvalidEmail      = errors.New("invalid email format")
    ErrEmailAlreadyExists = errors.New("email already exists")
)
```

## ディレクトリ構造規約

### Feature Sliced Design (FSD) 準拠

```
feature/
└── <domain_name>/
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
    │   ├── create_user.go    # 作成用ユースケース
    │   └── find_user.go      # 取得用ユースケース
    └── handler.go            # プレゼンテーション層
```

### 共有コンポーネント

```
share/
├── custom_error/           # カスタムエラー定義
├── domain/
│   └── model/
│       └── value_object/   # 基底バリューオブジェクト
├── usecase/               # 共通ユースケース
├── logger/                # ログ機能
├── validator/             # バリデーション機能
└── config/                # 設定管理
```

## コーディングスタイル

### インポート順序

```go
import (
    // 標準ライブラリ
    "context"
    "fmt"
    "time"

    // サードパーティライブラリ
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"

    // プロジェクト内パッケージ
    "project/feature/user/domain"
    "project/share/custom_error"
)
```

### 構造体定義

```go
// ✅ 推奨: フィールドの型揃え
type User struct {
    ID        UserID    `json:"id"`
    Name      UserName  `json:"name"`
    Email     Email     `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// ✅ 推奨: バリデーションタグの整理
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=1,max=100"`
    Email string `json:"email" validate:"required,email"`
}
```

### 関数定義

```go
// ✅ 推奨: 戻り値の型を明記
func (s *UserService) CreateUser(ctx context.Context, cmd CreateUserCommand) (*User, error) {
    // 実装
}

// ✅ 推奨: エラーハンドリング
func (r *UserRepository) FindByID(ctx context.Context, id UserID) (*User, error) {
    user, err := r.db.QueryUser(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %w", err)
    }
    return user, nil
}
```

### エラーハンドリング

```go
// ✅ 推奨: ラップエラーで詳細情報を保持
func (u *UserUsecase) CreateUser(ctx context.Context, input CreateUserInput) error {
    if err := u.validator.Validate(input); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    user, err := u.userService.CreateUser(ctx, input.ToCommand())
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}

// ✅ 推奨: カスタムエラー型の活用
func NewModelError(err error, message string) *ModelError {
    return &ModelError{
        BaseError: err,
        Message:   message,
        Code:      "MODEL_ERROR",
    }
}
```

## DDD パターン実装規約

### エンティティ

```go
// ✅ 推奨: イミュータブルなID
type User struct {
    id    UserID    // プライベート、不変
    name  UserName  // 値オブジェクト
    email Email     // 値オブジェクト
}

// Getter メソッド
func (u *User) ID() UserID { return u.id }
func (u *User) Name() UserName { return u.name }
func (u *User) Email() Email { return u.email }

// ビジネスロジック
func (u *User) ChangeName(newName UserName) error {
    if err := newName.Validate(); err != nil {
        return err
    }
    u.name = newName
    return nil
}
```

### 値オブジェクト

```go
// ✅ 推奨: 不変性とバリデーション
type Email struct {
    value string
}

func NewEmail(value string) (Email, error) {
    if !emailRegex.MatchString(value) {
        return Email{}, NewModelError(ErrInvalidEmail, "invalid email format")
    }
    return Email{value: value}, nil
}

func (e Email) String() string {
    return e.value
}

func (e Email) Equals(other Email) bool {
    return e.value == other.value
}
```

### リポジトリパターン

```go
// ✅ 推奨: インターフェース分離
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id UserID) (*User, error)
    FindByEmail(ctx context.Context, email Email) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id UserID) error
}

// クエリ専用インターフェース
type UserQueryRepository interface {
    FindByID(ctx context.Context, id UserID) (*User, error)
    FindByEmail(ctx context.Context, email Email) (*User, error)
    List(ctx context.Context, query UserListQuery) ([]*User, error)
}
```

## コメント規約

### 関数・メソッドコメント

```go
// CreateUser creates a new user with the given name and email.
// It validates the input data and returns an error if validation fails.
//
// Parameters:
//   - ctx: context for cancellation and timeout
//   - name: user's display name (required, 1-100 characters)
//   - email: user's email address (required, valid email format)
//
// Returns:
//   - *User: created user entity
//   - error: validation or creation error
func CreateUser(ctx context.Context, name string, email string) (*User, error) {
    // 実装
}
```

### 構造体コメント

```go
// User represents a user entity in the domain model.
// It encapsulates user identity and basic profile information.
//
// Business Rules:
//   - Email must be unique across the system
//   - Name cannot be empty and must be 1-100 characters
//   - ID is immutable once created
type User struct {
    id    UserID   // 一意識別子
    name  UserName // 表示名（1-100文字）
    email Email    // メールアドレス（システム内でユニーク）
}
```

## パフォーマンス規約

### データベースアクセス

```go
// ✅ 推奨: プリペアドステートメント使用
func (r *UserRepository) FindByEmail(ctx context.Context, email Email) (*User, error) {
    query := "SELECT id, name, email FROM users WHERE email = $1"
    row := r.db.QueryRowContext(ctx, query, email.String())
    // ...
}

// ✅ 推奨: バッチ処理
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
        if _, err := stmt.ExecContext(ctx, user.ID(), user.Name(), user.Email()); err != nil {
            return err
        }
    }

    return tx.Commit()
}
```

### メモリ効率

```go
// ✅ 推奨: スライスの事前容量確保
func (s *UserService) ProcessUsers(users []User) []ProcessedUser {
    result := make([]ProcessedUser, 0, len(users)) // 容量を事前確保
    for _, user := range users {
        result = append(result, s.processUser(user))
    }
    return result
}

// ✅ 推奨: ポインタを適切に使用
func (r *UserRepository) FindAll(ctx context.Context) ([]*User, error) {
    // 大きな構造体はポインタで返す
    return users, nil
}
```

## セキュリティ規約

### 入力検証

```go
// ✅ 推奨: 階層的バリデーション
func (h *UserHandler) CreateUser(c echo.Context) error {
    var req CreateUserRequest
    
    // 1. 構造的バリデーション
    if err := c.Bind(&req); err != nil {
        return NewHandlerError(ErrInvalidRequest, "invalid request format")
    }
    
    // 2. ビジネスルールバリデーション
    if err := h.validator.Validate(req); err != nil {
        return NewHandlerError(ErrValidationFailed, err.Error())
    }
    
    // 3. サニタイゼーション
    req.Name = html.EscapeString(req.Name)
    
    return h.usecase.Execute(c.Request().Context(), req.ToInput())
}
```

### SQL インジェクション対策

```go
// ✅ 推奨: プリペアドステートメント使用
func (r *UserRepository) Search(ctx context.Context, keyword string) ([]*User, error) {
    query := "SELECT id, name, email FROM users WHERE name ILIKE $1"
    rows, err := r.db.QueryContext(ctx, query, "%"+keyword+"%")
    // ...
}

// ❌ 非推奨: 文字列結合
func (r *UserRepository) SearchUnsafe(ctx context.Context, keyword string) ([]*User, error) {
    query := "SELECT id, name, email FROM users WHERE name ILIKE '%" + keyword + "%'"
    rows, err := r.db.QueryContext(ctx, query) // SQL インジェクション脆弱性
    // ...
}
```

## 禁止事項

### 一般的な禁止事項

1. **Global変数の使用禁止**
   ```go
   // ❌ 禁止
   var globalDB *sql.DB
   var globalLogger Logger
   ```

2. **パニックの使用禁止**（ライブラリ内部でのリカバリー不可能なエラー以外）
   ```go
   // ❌ 禁止
   func CreateUser(name string) User {
       if name == "" {
           panic("name is required")
       }
   }

   // ✅ 推奨
   func CreateUser(name string) (*User, error) {
       if name == "" {
           return nil, ErrNameRequired
       }
   }
   ```

3. **アンダースコア変数名の使用禁止**（Go言語コード内）
   ```go
   // ❌ 禁止
   func create_user(user_name string) error { }
   var user_repository UserRepository

   // ✅ 推奨
   func createUser(userName string) error { }
   var userRepository UserRepository
   ```

### DDD固有の禁止事項

1. **ドメイン層でのインフラ依存禁止**
   ```go
   // ❌ 禁止: ドメインサービスでDB直接アクセス
   func (s *UserService) CreateUser(name string) error {
       db, _ := sql.Open("postgres", "...")  // 禁止
   }
   ```

2. **エンティティでの外部API呼び出し禁止**
   ```go
   // ❌ 禁止: エンティティでHTTP通信
   func (u *User) SendWelcomeEmail() error {
       http.Post("https://api.email.com/send", ...) // 禁止
   }
   ```

3. **値オブジェクトの可変性禁止**
   ```go
   // ❌ 禁止: 値オブジェクトのセッター
   func (e *Email) SetValue(value string) {
       e.value = value // 禁止
   }
   ```

## 最終更新

- バージョン: 2025.05.24
- 更新者: AI Assistant
