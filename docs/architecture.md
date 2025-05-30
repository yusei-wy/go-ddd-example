# アーキテクチャ設計書

## 概要

このプロジェクトは、Go言語を使用してDDD（Domain-Driven Design）、CQRS（Command Query Responsibility Segregation）、FSD（Feature Sliced Design）パターンを採用したWebアプリケーションのサンプルです。

## アーキテクチャパターン

### 1. Domain-Driven Design (DDD)

DDDの戦術的パターンを実装し、ビジネスロジックをドメイン層に集約しています。

- **ドメイン層**: `feature/<domain>/domain` にビジネスロジックを集約
- **エンティティ**: ドメインモデルの識別子を持つオブジェクト
- **バリューオブジェクト**: 不変の値を表現するオブジェクト
- **ドメインサービス**: エンティティや値オブジェクトに属さないビジネスルール
- **リポジトリパターン**: データアクセスの抽象化

### 2. Command Query Responsibility Segregation (CQRS)

読み取り（Query）と書き込み（Command）の責務を分離しています。

- **コマンド**: `feature/<domain>/domain/model/command.go` で書き込み用モデル
- **クエリ**: `feature/<domain>/domain/model/query.go` で読み取り用モデル
- **分離されたユースケース**: 作成・更新用と読み取り用のユースケースを分離

### 3. Feature Sliced Design (FSD)

機能ごとにディレクトリを分割し、スケーラブルな構造を実現しています。

```
feature/
└── <domain_name>/
    ├── domain/         # ビジネスロジック
    │   ├── model/
    │   │   ├── command/
    │   │   ├── query/
    │   │   └── value_object/
    │   ├── repository.go
    │   └── service.go
    ├── infra/          # インフラストラクチャ
    │   ├── model.go
    │   └── psql_repository.go
    ├── usecase/        # アプリケーションロジック
    └── handler.go      # プレゼンテーション
```

## レイヤードアーキテクチャ

### プレゼンテーション層

HTTPリクエストの受信とレスポンス生成を担当します。

- **ハンドラー**: `feature/<domain>/handler.go` - HTTPリクエスト処理
- **サーバー設定**: `server/handlers.go` - ルーティング設定
- **エラーハンドリング**: `server/response_error.go` - カスタムエラーレスポンス

### アプリケーション層

ビジネスユースケースの実行を制御します。

- **ユースケース**: `feature/<domain>/usecase/` - アプリケーションロジック
- **DTOパターン**: Input/Output構造体でデータ転送
- **依存関係の制御**: ドメイン層のサービスとリポジトリを利用

### ドメイン層

ビジネスロジックの中核を担います。

- **エンティティ・値オブジェクト**: `feature/<domain>/domain/model/`
- **ドメインサービス**: `feature/<domain>/domain/service.go`
- **リポジトリインターフェース**: `feature/<domain>/domain/repository.go`

### インフラストラクチャ層

外部システムとの連携を担当します。

- **データベースアクセス**: `feature/<domain>/infra/psql_repository.go`
- **データモデル**: `feature/<domain>/infra/model.go`
- **外部API連携**: 必要に応じて追加

## バリデーション戦略

レイヤー別に適切なバリデーションを実装し、データの整合性を保証します。

### 1. 入力バリデーション（プレゼンテーション層）

```go
// HTTPリクエストの基本的な形式チェック
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=1,max=100"`
    Email string `json:"email" validate:"required,email"`
}
```

### 2. ビジネスルールバリデーション（ドメイン層）

```go
// 値オブジェクトでのビジネスルール検証
func NewEmail(value string) (Email, error) {
    if !emailRegex.MatchString(value) {
        return Email{}, NewModelError(ErrInvalidEmail, "invalid email format")
    }
    return Email{value: value}, nil
}
```

### 3. データ整合性バリデーション（ユースケース層）

```go
// ユニーク制約などの外部依存を伴う検証
func (u *UserUsecase) CreateUser(ctx context.Context, input CreateUserInput) error {
    if exists := u.userRepo.ExistsByEmail(ctx, input.Email); exists {
        return NewUsecaseError(ErrEmailAlreadyExists, "email already exists")
    }
    // ...existing logic...
}
```

## 命名規則

プロジェクトでは、Go言語の公式ガイドライン「Effective Go」に従った命名規則を採用しています。

詳細な命名規則については、**[コーディング規約](./rules/coding-standards.md#命名規則)** を参照してください。

## テスト戦略（TDD重視）

Test-Driven Development（TDD）を採用し、品質の高いコードを実現します。

詳細なテスト実装ガイドについては、**[テストガイドライン](./rules/testing-guidelines.md)** を参照してください。

### テスト分類

- **ユニットテスト**: 各レイヤーの単体機能をテスト
- **統合テスト**: 複数のコンポーネント間の連携をテスト  
- **E2Eテスト**: APIエンドポイントの全体的な動作をテスト
- **目標カバレッジ**: 80%以上

## 実装例

### ユーザー管理機能の実装例

#### 1. ドメインモデル

```go
// User エンティティ
type User struct {
    id    UserID
    name  UserName
    email Email
}

// Email 値オブジェクト
type Email struct {
    value string
}

func NewEmail(value string) (Email, error) {
    if !emailRegex.MatchString(value) {
        return Email{}, NewModelError(ErrInvalidEmail, "invalid email format")
    }
    return Email{value: value}, nil
}
```

#### 2. リポジトリインターフェース

```go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id UserID) (*User, error)
    FindByEmail(ctx context.Context, email Email) (*User, error)
    ExistsByEmail(ctx context.Context, email Email) bool
}
```

#### 3. ユースケース

```go
type CreateUserUsecase struct {
    userRepo UserRepository
    logger   Logger
}

func (u *CreateUserUsecase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
    // バリデーション
    if u.userRepo.ExistsByEmail(ctx, input.Email) {
        return nil, NewUsecaseError(ErrEmailAlreadyExists, "email already exists")
    }
    
    // ドメインオブジェクト生成
    user := NewUser(input.Name, input.Email)
    
    // 永続化
    if err := u.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    return &CreateUserOutput{ID: user.ID()}, nil
}
```

#### 4. ハンドラー

```go
func (h *UserHandler) CreateUser(c echo.Context) error {
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        return NewHandlerError(ErrInvalidRequest, "invalid request format")
    }
    
    input := CreateUserInput{
        Name:  req.Name,
        Email: req.Email,
    }
    
    output, err := h.createUserUsecase.Execute(c.Request().Context(), input)
    if err != nil {
        return err
    }
    
    return c.JSON(http.StatusCreated, CreateUserResponse{ID: output.ID})
}
```

## セキュリティ対策

Webアプリケーションの一般的な脅威に対する防御策を実装します。

詳細なセキュリティ実装については、**[セキュリティガイドライン](./rules/security-guidelines.md)** を参照してください。

## ログ戦略

構造化ログによる運用監視とトラブルシューティングを支援します。

### 1. ログレベル

- **ERROR**: システムエラー、障害
- **WARN**: 警告、非推奨機能使用
- **INFO**: 重要な処理の開始・終了
- **DEBUG**: デバッグ情報（開発環境のみ）

### 2. ログ出力例

```go
type Logger interface {
    Error(ctx context.Context, msg string, fields ...Field)
    Warn(ctx context.Context, msg string, fields ...Field)
    Info(ctx context.Context, msg string, fields ...Field)
    Debug(ctx context.Context, msg string, fields ...Field)
}

// 使用例
logger.Info(ctx, "user created",
    Field("user_id", user.ID()),
    Field("email", user.Email()),
    Field("duration", time.Since(start)),
)
```

## パフォーマンス考慮事項

### 1. データベース最適化

```go
// インデックス設計
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_orders_user_id ON orders(user_id);

// コネクションプール設定
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

### 2. キャッシュ戦略

```go
// Redis キャッシュ例
type CacheUserRepository struct {
    repo  UserRepository
    cache Cache
    ttl   time.Duration
}

func (r *CacheUserRepository) FindByID(ctx context.Context, id UserID) (*User, error) {
    key := fmt.Sprintf("user:%s", id.String())
    if cached := r.cache.Get(key); cached != nil {
        return cached.(*User), nil
    }
    
    user, err := r.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    r.cache.Set(key, user, r.ttl)
    return user, nil
}
```

## 設定管理

環境別の設定を適切に管理します。

### 1. 設定構造体

```go
type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
    Logger   LoggerConfig   `yaml:"logger"`
}

type ServerConfig struct {
    Port         int    `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
    ReadTimeout  int    `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT" env-default:"30"`
    WriteTimeout int    `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT" env-default:"30"`
}
```

### 2. 環境別設定ファイル

```yaml
# config/development.yaml
server:
  port: 8080
  read_timeout: 30
  write_timeout: 30

database:
  host: localhost
  port: 5432
  name: app_development
  
logger:
  level: debug
  format: json
```

## エラーハンドリング戦略

レイヤー別の専用エラー型を定義し、適切なエラー処理を実現しています。

詳細なエラーハンドリングの実装方法については、**[エラーハンドリング規約](./rules/error-handling.md)** を参照してください。

## 共有コンポーネント

複数のドメインで使用される共通機能を提供します。

```
share/
├── custom_error/           # カスタムエラー定義
├── domain/
│   └── model/
│       └── valueobject/    # 基底バリューオブジェクト
├── usecase/               # 共通ユースケース
├── logger/                # ログ機能
├── validator/             # バリデーション機能
└── config/                # 設定管理
```

## データベース設計

- **マイグレーション**: `db/migrations/` でスキーマ管理
- **Dockerコンテナ**: 開発環境での一貫性確保

## 開発・運用

### 利用可能なコマンド

| コマンド | 説明 |
|---------|------|
| `make dev` | 開発サーバー起動 |
| `make reset-db` | DB初期化・マイグレーション実行 |
| `make lint` | コード品質チェック |
| `make format` | コードフォーマット |
| `make test` | テスト実行 |
| `make test-coverage` | カバレッジ付きテスト実行 |
| `make test-integration` | 統合テスト実行 |

### TDD開発フロー

1. **要件定義**: ユーザーストーリーの作成
2. **テスト作成**: 失敗するテストケースの実装
3. **実装**: テストを通す最小限のコード作成
4. **リファクタリング**: コード品質の向上
5. **統合**: 他機能との連携確認

### 今後の拡張予定

- [ ] ロギング
- [ ] 依存性注入の実装
- [ ] OpenAPIによるリクエスト・レスポンス型生成
- [ ] メトリクス収集機能
- [ ] 分散トレーシング（メトリクス収集と連携したObservability強化）
- [ ] CI/CDパイプライン

## 設計原則

1. **関心の分離**: 各レイヤーが明確な責務を持つ
2. **依存関係逆転**: 上位レイヤーが下位レイヤーのインターフェースに依存
3. **テスタビリティ**: モックやスタブによる単体テストが容易
4. **拡張性**: 新機能追加時の影響範囲を最小化
5. **保守性**: コードの可読性と変更容易性を重視
6. **セキュリティファースト**: セキュリティ要件を設計段階から考慮
7. **パフォーマンス意識**: スケーラビリティを考慮した設計
8. **運用性**: 監視・ログ・デバッグの容易さを重視
