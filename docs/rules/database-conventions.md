# データベース規約

## 概要

Go DDD プロジェクトにおけるデータベース設計・実装の規約とベストプラクティスを定義します。

## 基本原則

### データベース設計原則
- **正規化**: 適切な正規化による冗長性排除
- **パフォーマンス**: インデックス戦略とクエリ最適化
- **スケーラビリティ**: 将来の拡張を考慮した設計
- **整合性**: 制約とトランザクションによるデータ整合性確保

## テーブル設計規約

### 命名規則

```sql
-- テーブル名：スネークケース、複数形
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 中間テーブル：関連するテーブル名を結合
CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(id),
    role_id UUID NOT NULL REFERENCES roles(id),
    PRIMARY KEY (user_id, role_id)
);

-- インデックス名：idx_テーブル名_カラム名
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);
```

### カラム規約

```sql
-- 主キー：UUIDを使用
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

-- 外部キー：参照先テーブル名_id
user_id UUID NOT NULL REFERENCES users(id),
order_id UUID NOT NULL REFERENCES orders(id),

-- 論理削除：deleted_at
deleted_at TIMESTAMP WITH TIME ZONE,

-- 作成・更新日時：必須
created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

-- バージョン管理：楽観的ロック
version INTEGER NOT NULL DEFAULT 1,

-- 列挙型：CHECK制約またはENUM
status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'processing', 'completed', 'cancelled')),

-- 金額：DECIMAL使用
price DECIMAL(10, 2) NOT NULL,
```

## マイグレーション戦略

### マイグレーションファイル構造

```
db/
├── migrations/
│   ├── 001_create_users_table.up.sql
│   ├── 001_create_users_table.down.sql
│   ├── 002_create_orders_table.up.sql
│   ├── 002_create_orders_table.down.sql
│   └── ...
├── seeds/
│   ├── development/
│   │   ├── users.sql
│   │   └── roles.sql
│   └── test/
│       └── test_data.sql
└── schema.sql
```

### マイグレーションファイル例

```sql
-- 001_create_users_table.up.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- インデックス作成
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_active ON users(is_active) WHERE is_active = true;
CREATE INDEX idx_users_created_at ON users(created_at);

-- 更新日時自動更新トリガー
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
```

```sql
-- 001_create_users_table.down.sql
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS users;
```

### マイグレーション実行

```go
// マイグレーション実行関数
func RunMigrations(db *sql.DB, migrationsPath string) error {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return err
    }
    
    m, err := migrate.NewWithDatabaseInstance(
        fmt.Sprintf("file://%s", migrationsPath),
        "postgres",
        driver,
    )
    if err != nil {
        return err
    }
    
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }
    
    return nil
}
```

## リポジトリパターン実装

### リポジトリインターフェース

```go
// ドメイン層でのリポジトリインターフェース定義
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id UserID) (*User, error)
    FindByEmail(ctx context.Context, email Email) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id UserID) error
    ExistsByEmail(ctx context.Context, email Email) bool
    FindWithPagination(ctx context.Context, query PaginationQuery) ([]*User, int, error)
}

// クエリ専用リポジトリ（CQRS）
type UserQueryRepository interface {
    FindByID(ctx context.Context, id UserID) (*UserQueryModel, error)
    Search(ctx context.Context, query UserSearchQuery) ([]*UserQueryModel, error)
    GetStatistics(ctx context.Context) (*UserStatistics, error)
}
```

### PostgreSQL実装

```go
type postgresUserRepository struct {
    db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
    return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *User) error {
    query := `
        INSERT INTO users (id, name, email, password_hash, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
    
    _, err := r.db.ExecContext(ctx, query,
        user.ID().String(),
        user.Name().String(),
        user.Email().String(),
        user.PasswordHash(),
        user.IsActive(),
        user.CreatedAt(),
        user.UpdatedAt(),
    )
    
    if err != nil {
        if isUniqueViolation(err) {
            return NewRepositoryError(ErrUserEmailDuplicate, "email already exists")
        }
        return NewRepositoryError(ErrDatabaseError, "failed to create user")
    }
    
    return nil
}

func (r *postgresUserRepository) FindByID(ctx context.Context, id UserID) (*User, error) {
    query := `
        SELECT id, name, email, password_hash, is_active, created_at, updated_at, version
        FROM users 
        WHERE id = $1 AND deleted_at IS NULL
    `
    
    row := r.db.QueryRowContext(ctx, query, id.String())
    
    var userModel UserModel
    err := row.Scan(
        &userModel.ID,
        &userModel.Name,
        &userModel.Email,
        &userModel.PasswordHash,
        &userModel.IsActive,
        &userModel.CreatedAt,
        &userModel.UpdatedAt,
        &userModel.Version,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, NewRepositoryError(ErrUserNotFound, "user not found")
        }
        return nil, NewRepositoryError(ErrDatabaseError, "failed to find user")
    }
    
    return userModel.ToDomain()
}

func (r *postgresUserRepository) Update(ctx context.Context, user *User) error {
    query := `
        UPDATE users 
        SET name = $2, email = $3, is_active = $4, updated_at = $5, version = version + 1
        WHERE id = $1 AND version = $6 AND deleted_at IS NULL
    `
    
    result, err := r.db.ExecContext(ctx, query,
        user.ID().String(),
        user.Name().String(),
        user.Email().String(),
        user.IsActive(),
        time.Now(),
        user.Version(),
    )
    
    if err != nil {
        return NewRepositoryError(ErrDatabaseError, "failed to update user")
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return NewRepositoryError(ErrDatabaseError, "failed to get affected rows")
    }
    
    if rowsAffected == 0 {
        return NewRepositoryError(ErrOptimisticLock, "user was modified by another process")
    }
    
    return nil
}
```

## データモデル実装

### インフラ層データモデル

```go
// インフラ層のデータモデル
type UserModel struct {
    ID           string    `db:"id"`
    Name         string    `db:"name"`
    Email        string    `db:"email"`
    PasswordHash string    `db:"password_hash"`
    IsActive     bool      `db:"is_active"`
    CreatedAt    time.Time `db:"created_at"`
    UpdatedAt    time.Time `db:"updated_at"`
    DeletedAt    *time.Time `db:"deleted_at"`
    Version      int       `db:"version"`
}

// ドメインモデルへの変換
func (m UserModel) ToDomain() (*User, error) {
    id, err := NewUserID(m.ID)
    if err != nil {
        return nil, err
    }
    
    name, err := NewUserName(m.Name)
    if err != nil {
        return nil, err
    }
    
    email, err := NewEmail(m.Email)
    if err != nil {
        return nil, err
    }
    
    return ReconstructUser(id, name, email, m.PasswordHash, m.IsActive, m.CreatedAt, m.UpdatedAt, m.Version), nil
}

// ドメインモデルからの変換
func NewUserModel(user *User) UserModel {
    return UserModel{
        ID:           user.ID().String(),
        Name:         user.Name().String(),
        Email:        user.Email().String(),
        PasswordHash: user.PasswordHash(),
        IsActive:     user.IsActive(),
        CreatedAt:    user.CreatedAt(),
        UpdatedAt:    user.UpdatedAt(),
        Version:      user.Version(),
    }
}
```

## トランザクション管理

### トランザクション実装

```go
type TransactionManager interface {
    WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type postgresTransactionManager struct {
    db *sql.DB
}

func NewPostgresTransactionManager(db *sql.DB) TransactionManager {
    return &postgresTransactionManager{db: db}
}

func (tm *postgresTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
    tx, err := tm.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()
    
    // コンテキストにトランザクションを設定
    ctx = context.WithValue(ctx, "tx", tx)
    
    if err := fn(ctx); err != nil {
        tx.Rollback()
        return err
    }
    
    return tx.Commit()
}

// トランザクション対応リポジトリ
func (r *postgresUserRepository) getDB(ctx context.Context) database {
    if tx, ok := ctx.Value("tx").(*sql.Tx); ok {
        return tx
    }
    return r.db
}

type database interface {
    ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
    QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
    QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
```

### ユースケースでのトランザクション使用

```go
func (u *CreateOrderUsecase) Execute(ctx context.Context, input CreateOrderInput) (*CreateOrderOutput, error) {
    var output *CreateOrderOutput
    
    err := u.txManager.WithTransaction(ctx, func(ctx context.Context) error {
        // 在庫確認
        product, err := u.productRepo.FindByID(ctx, input.ProductID)
        if err != nil {
            return err
        }
        
        if !product.HasStock(input.Quantity) {
            return NewUsecaseError(ErrInsufficientStock, "insufficient stock")
        }
        
        // 注文作成
        order := NewOrder(input.UserID, input.ProductID, input.Quantity, product.Price())
        if err := u.orderRepo.Create(ctx, order); err != nil {
            return err
        }
        
        // 在庫減少
        product.ReduceStock(input.Quantity)
        if err := u.productRepo.Update(ctx, product); err != nil {
            return err
        }
        
        output = &CreateOrderOutput{ID: order.ID()}
        return nil
    })
    
    return output, err
}
```

## インデックス戦略

### インデックス設計

```sql
-- 単一カラムインデックス
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_orders_created_at ON orders(created_at);

-- 複合インデックス（クエリパターンに応じて）
CREATE INDEX idx_orders_user_status ON orders(user_id, status);
CREATE INDEX idx_products_category_price ON products(category_id, price);

-- 部分インデックス（WHERE句付き）
CREATE INDEX idx_users_active ON users(email) WHERE is_active = true;
CREATE INDEX idx_orders_pending ON orders(created_at) WHERE status = 'pending';

-- 関数インデックス
CREATE INDEX idx_users_lower_email ON users(lower(email));

-- 全文検索インデックス
CREATE INDEX idx_products_search ON products USING gin(to_tsvector('english', name || ' ' || description));
```

### インデックス監視

```sql
-- 未使用インデックスの確認
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
WHERE idx_tup_read = 0 AND idx_tup_fetch = 0;

-- インデックスサイズ確認
SELECT 
    indexname,
    pg_size_pretty(pg_relation_size(indexname::regclass)) as index_size
FROM pg_indexes
WHERE tablename = 'users';
```

## クエリ最適化

### 効率的なクエリパターン

```go
// ページネーション（OFFSET/LIMIT）
func (r *postgresUserRepository) FindWithPagination(ctx context.Context, query PaginationQuery) ([]*User, int, error) {
    // 件数取得
    countQuery := `SELECT COUNT(*) FROM users WHERE is_active = true AND deleted_at IS NULL`
    var total int
    err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
    if err != nil {
        return nil, 0, err
    }
    
    // データ取得
    dataQuery := `
        SELECT id, name, email, created_at, updated_at
        FROM users 
        WHERE is_active = true AND deleted_at IS NULL
        ORDER BY created_at DESC, id DESC
        LIMIT $1 OFFSET $2
    `
    
    rows, err := r.db.QueryContext(ctx, dataQuery, query.Limit, query.Offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()
    
    users, err := r.scanUsers(rows)
    return users, total, err
}

// カーソルベースページネーション（大量データ向け）
func (r *postgresUserRepository) FindWithCursor(ctx context.Context, cursor string, limit int) ([]*User, string, error) {
    var query string
    var args []interface{}
    
    if cursor == "" {
        query = `
            SELECT id, name, email, created_at
            FROM users 
            WHERE is_active = true AND deleted_at IS NULL
            ORDER BY created_at DESC, id DESC
            LIMIT $1
        `
        args = []interface{}{limit + 1} // 次のページ存在確認用
    } else {
        cursorTime, cursorID, err := parseCursor(cursor)
        if err != nil {
            return nil, "", err
        }
        
        query = `
            SELECT id, name, email, created_at
            FROM users 
            WHERE is_active = true AND deleted_at IS NULL
            AND (created_at < $1 OR (created_at = $1 AND id < $2))
            ORDER BY created_at DESC, id DESC
            LIMIT $3
        `
        args = []interface{}{cursorTime, cursorID, limit + 1}
    }
    
    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, "", err
    }
    defer rows.Close()
    
    users, err := r.scanUsers(rows)
    if err != nil {
        return nil, "", err
    }
    
    var nextCursor string
    if len(users) > limit {
        lastUser := users[limit-1]
        nextCursor = createCursor(lastUser.CreatedAt(), lastUser.ID())
        users = users[:limit]
    }
    
    return users, nextCursor, nil
}
```

## データベース接続管理

### 接続プール設定

```go
type DatabaseConfig struct {
    Host            string `env:"DB_HOST" env-default:"localhost"`
    Port            int    `env:"DB_PORT" env-default:"5432"`
    Name            string `env:"DB_NAME,required"`
    User            string `env:"DB_USER,required"`
    Password        string `env:"DB_PASSWORD,required"`
    SSLMode         string `env:"DB_SSL_MODE" env-default:"disable"`
    MaxOpenConns    int    `env:"DB_MAX_OPEN_CONNS" env-default:"25"`
    MaxIdleConns    int    `env:"DB_MAX_IDLE_CONNS" env-default:"5"`
    ConnMaxLifetime int    `env:"DB_CONN_MAX_LIFETIME" env-default:"300"` // seconds
    ConnMaxIdleTime int    `env:"DB_CONN_MAX_IDLE_TIME" env-default:"60"`  // seconds
}

func NewDatabase(config DatabaseConfig) (*sql.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        config.Host, config.Port, config.User, config.Password, config.Name, config.SSLMode,
    )
    
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }
    
    // 接続プール設定
    db.SetMaxOpenConns(config.MaxOpenConns)
    db.SetMaxIdleConns(config.MaxIdleConns)
    db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
    db.SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime) * time.Second)
    
    // 接続確認
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    return db, nil
}
```

## エラーハンドリング

### データベースエラー分類

```go
func classifyDatabaseError(err error) error {
    var pgErr *pq.Error
    if errors.As(err, &pgErr) {
        switch pgErr.Code {
        case "23505": // unique_violation
            return NewRepositoryError(ErrDuplicateKey, "duplicate key violation")
        case "23503": // foreign_key_violation
            return NewRepositoryError(ErrForeignKeyViolation, "foreign key violation")
        case "23514": // check_violation
            return NewRepositoryError(ErrCheckViolation, "check constraint violation")
        case "23502": // not_null_violation
            return NewRepositoryError(ErrNotNullViolation, "not null constraint violation")
        }
    }
    
    return NewRepositoryError(ErrDatabaseError, "database error")
}

func isUniqueViolation(err error) bool {
    var pgErr *pq.Error
    return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func isForeignKeyViolation(err error) bool {
    var pgErr *pq.Error
    return errors.As(err, &pgErr) && pgErr.Code == "23503"
}
```

## テストデータベース

### テスト用データベース設定

```go
func SetupTestDB(t *testing.T) *sql.DB {
    config := DatabaseConfig{
        Host:     "localhost",
        Port:     5432,
        Name:     "test_db",
        User:     "test_user",
        Password: "test_password",
        SSLMode:  "disable",
    }
    
    db, err := NewDatabase(config)
    if err != nil {
        t.Fatalf("Failed to connect to test database: %v", err)
    }
    
    // マイグレーション実行
    if err := RunMigrations(db, "../../db/migrations"); err != nil {
        t.Fatalf("Failed to run migrations: %v", err)
    }
    
    // テスト後のクリーンアップ
    t.Cleanup(func() {
        cleanupTestDB(t, db)
        db.Close()
    })
    
    return db
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
    tables := []string{"user_roles", "users", "orders", "products"}
    
    for _, table := range tables {
        _, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
        if err != nil {
            t.Logf("Failed to truncate table %s: %v", table, err)
        }
    }
}
```

## データベース規約チェックリスト

### 設計段階
- [ ] 適切な正規化レベル
- [ ] 主キー・外部キー設計
- [ ] インデックス戦略
- [ ] 制約定義
- [ ] 命名規則準拠

### 実装段階
- [ ] マイグレーションファイル作成
- [ ] リポジトリパターン実装
- [ ] トランザクション管理
- [ ] エラーハンドリング
- [ ] テストデータベース設定

### パフォーマンス
- [ ] クエリ実行計画確認
- [ ] インデックス効果測定
- [ ] N+1問題回避
- [ ] 接続プール最適化

### 運用・保守
- [ ] バックアップ戦略
- [ ] 監視・アラート設定
- [ ] パフォーマンス監視
- [ ] 容量管理

## 参考資料

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Go Database/SQL Tutorial](https://go.dev/doc/tutorial/database-access)
- [SQL Performance Explained](https://use-the-index-luke.com/)

## 最終更新

- 日付: 2025.05.24
- 更新者: AI Assistant
