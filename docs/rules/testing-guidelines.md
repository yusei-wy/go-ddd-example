# テストガイドライン

## 概要

このドキュメントは、Go DDD プロジェクトにおけるテスト戦略と Test-Driven Development (TDD) の実践ガイドラインを定義します。

## TDD 実践ガイド

### TDD サイクル

```
Red → Green → Refactor
┌─────────────────────────────────────────┐
│ 1. Red: 失敗するテストを書く              │
│    ↓                                   │
│ 2. Green: テストが通る最小限のコードを書く │
│    ↓                                   │
│ 3. Refactor: コードをリファクタリングする  │
│    ↓                                   │
│ (繰り返し)                              │
└─────────────────────────────────────────┘
```

### TDD の実践手順

#### 1. Red フェーズ: 失敗するテストを書く

```go
// まず、失敗するテストを書く
func TestNewEmail_ValidEmail_ReturnsEmail(t *testing.T) {
    // Given
    validEmail := "test@example.com"
    
    // When
    email, err := NewEmail(validEmail)  // この関数はまだ存在しない
    
    // Then
    assert.NoError(t, err)
    assert.Equal(t, validEmail, email.String())
}

func TestNewEmail_InvalidEmail_ReturnsError(t *testing.T) {
    // Given
    invalidEmail := "invalid-email"
    
    // When
    _, err := NewEmail(invalidEmail)
    
    // Then
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "invalid email format")
}
```

#### 2. Green フェーズ: テストが通る最小限のコードを書く

```go
// テストを通すための最小限の実装
import "regexp"

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Email struct {
    value string
}

func NewEmail(value string) (Email, error) {
    if !emailRegex.MatchString(value) {
        return Email{}, errors.New("invalid email format")
    }
    return Email{value: value}, nil
}

func (e Email) String() string {
    return e.value
}
```

#### 3. Refactor フェーズ: コードの改善

```go
// リファクタリング後の改善版
func NewEmail(value string) (Email, error) {
    if value == "" {
        return Email{}, NewModelError(ErrEmailEmpty, "email cannot be empty")
    }
    
    if !emailRegex.MatchString(value) {
        return Email{}, NewModelError(ErrInvalidEmail, "invalid email format")
    }
    
    if len(value) > MaxEmailLength {
        return Email{}, NewModelError(ErrEmailTooLong, "email too long")
    }
    
    return Email{value: strings.ToLower(value)}, nil
}
```

## テスト分類と戦略

### 1. ユニットテスト（Unit Tests）

各コンポーネントの単体機能をテストします。

#### ドメインモデルのテスト

```go
func TestUser_ChangeName_ValidName_UpdatesName(t *testing.T) {
    // Given
    user := createTestUser()
    newName, _ := NewUserName("New Name")
    
    // When
    err := user.ChangeName(newName)
    
    // Then
    assert.NoError(t, err)
    assert.Equal(t, newName, user.Name())
}

func TestUser_ChangeName_InvalidName_ReturnsError(t *testing.T) {
    // Given
    user := createTestUser()
    
    // When
    err := user.ChangeName(UserName{}) // 無効な名前
    
    // Then
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "name cannot be empty")
}
```

#### 値オブジェクトのテスト

```go
func TestEmail_Equals_SameEmail_ReturnsTrue(t *testing.T) {
    // Given
    email1, _ := NewEmail("test@example.com")
    email2, _ := NewEmail("test@example.com")
    
    // When & Then
    assert.True(t, email1.Equals(email2))
}

func TestEmail_Equals_DifferentEmail_ReturnsFalse(t *testing.T) {
    // Given
    email1, _ := NewEmail("test1@example.com")
    email2, _ := NewEmail("test2@example.com")
    
    // When & Then
    assert.False(t, email1.Equals(email2))
}
```

#### ドメインサービスのテスト

```go
func TestUserService_CreateUser_ValidData_ReturnsUser(t *testing.T) {
    // Given
    service := NewUserService()
    cmd := CreateUserCommand{
        Name:  "Test User",
        Email: "test@example.com",
    }
    
    // When
    user, err := service.CreateUser(context.Background(), cmd)
    
    // Then
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, cmd.Name, user.Name().String())
}
```

### 2. 統合テスト（Integration Tests）

複数のコンポーネント間の連携をテストします。

#### リポジトリ統合テスト

```go
func TestUserRepository_Create_ValidUser_SavesUser(t *testing.T) {
    // Given
    db := setupTestDB(t)
    repo := NewUserRepository(db)
    user := createTestUser()
    
    // When
    err := repo.Create(context.Background(), user)
    
    // Then
    assert.NoError(t, err)
    
    // データベースから取得して検証
    found, err := repo.FindByID(context.Background(), user.ID())
    assert.NoError(t, err)
    assert.Equal(t, user.Email(), found.Email())
}

func TestUserRepository_FindByEmail_ExistingUser_ReturnsUser(t *testing.T) {
    // Given
    db := setupTestDB(t)
    repo := NewUserRepository(db)
    user := createTestUser()
    repo.Create(context.Background(), user)
    
    // When
    found, err := repo.FindByEmail(context.Background(), user.Email())
    
    // Then
    assert.NoError(t, err)
    assert.Equal(t, user.ID(), found.ID())
}
```

#### ユースケース統合テスト

```go
func TestCreateUserUsecase_ValidInput_CreatesUser(t *testing.T) {
    // Given
    db := setupTestDB(t)
    repo := NewUserRepository(db)
    usecase := NewCreateUserUsecase(repo, mockLogger)
    
    input := CreateUserInput{
        Name:  "Test User",
        Email: "test@example.com",
    }
    
    // When
    output, err := usecase.Execute(context.Background(), input)
    
    // Then
    assert.NoError(t, err)
    assert.NotEmpty(t, output.ID)
    
    // データベースから検証
    user, err := repo.FindByID(context.Background(), UserID(output.ID))
    assert.NoError(t, err)
    assert.Equal(t, input.Email, user.Email().String())
}
```

### 3. E2E テスト（End-to-End Tests）

APIエンドポイントの全体的な動作をテストします。

#### HTTP API テスト

```go
func TestCreateUserAPI_ValidRequest_ReturnsCreated(t *testing.T) {
    // Given
    server := setupTestServer(t)
    defer server.Close()
    
    request := CreateUserRequest{
        Name:  "Test User",
        Email: "test@example.com",
    }
    
    // When
    resp, err := http.Post(
        server.URL+"/users",
        "application/json",
        toJSONReader(request),
    )
    
    // Then
    assert.NoError(t, err)
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
    
    var response CreateUserResponse
    json.NewDecoder(resp.Body).Decode(&response)
    assert.NotEmpty(t, response.ID)
}

func TestCreateUserAPI_InvalidEmail_ReturnsBadRequest(t *testing.T) {
    // Given
    server := setupTestServer(t)
    defer server.Close()
    
    request := CreateUserRequest{
        Name:  "Test User",
        Email: "invalid-email",
    }
    
    // When
    resp, err := http.Post(
        server.URL+"/users",
        "application/json",
        toJSONReader(request),
    )
    
    // Then
    assert.NoError(t, err)
    assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    
    var errorResp ErrorResponse
    json.NewDecoder(resp.Body).Decode(&errorResp)
    assert.Contains(t, errorResp.Message, "invalid email")
}
```

## テストファイル構成

### ディレクトリ構造

```
feature/user/
├── domain/
│   ├── model/
│   │   ├── user.go
│   │   ├── user_test.go          # ユニットテスト
│   │   ├── email.go
│   │   └── email_test.go         # ユニットテスト
│   ├── service.go
│   └── service_test.go           # ユニットテスト
├── usecase/
│   ├── create_user.go
│   ├── create_user_test.go       # ユニットテスト
│   └── create_user_integration_test.go  # 統合テスト
├── infra/
│   ├── psql_repository.go
│   └── psql_repository_test.go   # 統合テスト
├── handler.go
├── handler_test.go               # ユニットテスト
└── handler_e2e_test.go          # E2Eテスト
```

### テストファイル命名規則

| テストタイプ | 命名規則 | 例 |
|--------------|----------|-----|
| ユニットテスト | `*_test.go` | `user_test.go` |
| 統合テスト | `*_integration_test.go` | `repository_integration_test.go` |
| E2Eテスト | `*_e2e_test.go` | `api_e2e_test.go` |

## テスト用ヘルパー関数

### テストデータ生成

```go
// testdata/user.go
func CreateTestUser() *User {
    id := NewUserID()
    name, _ := NewUserName("Test User")
    email, _ := NewEmail("test@example.com")
    
    return &User{
        id:    id,
        name:  name,
        email: email,
    }
}

func CreateTestUserWithEmail(emailStr string) *User {
    id := NewUserID()
    name, _ := NewUserName("Test User")
    email, _ := NewEmail(emailStr)
    
    return &User{
        id:    id,
        name:  name,
        email: email,
    }
}

func CreateTestUsers(count int) []*User {
    users := make([]*User, count)
    for i := 0; i < count; i++ {
        users[i] = CreateTestUserWithEmail(fmt.Sprintf("test%d@example.com", i))
    }
    return users
}
```

### データベーステストセットアップ

```go
// testutil/database.go
func SetupTestDB(t *testing.T) *sql.DB {
    t.Helper()
    
    // テスト用データベース接続
    dsn := "postgres://user:pass@localhost/test_db?sslmode=disable"
    db, err := sql.Open("postgres", dsn)
    require.NoError(t, err)
    
    // マイグレーション実行
    err = runMigrations(db)
    require.NoError(t, err)
    
    // テスト終了時にクリーンアップ
    t.Cleanup(func() {
        cleanupTestDB(t, db)
        db.Close()
    })
    
    return db
}

func CleanupTestDB(t *testing.T, db *sql.DB) {
    t.Helper()
    
    tables := []string{"users", "orders", "products"}
    for _, table := range tables {
        _, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
        require.NoError(t, err)
    }
}
```

### モックオブジェクト

```go
// mock/user_repository.go
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id UserID) (*User, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email Email) bool {
    args := m.Called(ctx, email)
    return args.Bool(0)
}

// テストでの使用例
func TestCreateUserUsecase_EmailExists_ReturnsError(t *testing.T) {
    // Given
    mockRepo := new(MockUserRepository)
    usecase := NewCreateUserUsecase(mockRepo, mockLogger)
    
    email, _ := NewEmail("existing@example.com")
    mockRepo.On("ExistsByEmail", mock.Anything, email).Return(true)
    
    input := CreateUserInput{Email: "existing@example.com"}
    
    // When
    _, err := usecase.Execute(context.Background(), input)
    
    // Then
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "email already exists")
    mockRepo.AssertExpectations(t)
}
```

## テスト実行とカバレッジ

### テスト実行コマンド

```bash
# すべてのテスト実行
make test

# 特定パッケージのテスト実行
go test ./feature/user/domain/model/...

# 詳細出力付きテスト実行
go test -v ./...

# 統合テストのみ実行
go test -tags=integration ./...

# E2Eテストのみ実行
go test -tags=e2e ./...

# 並列実行数制限
go test -p 4 ./...

# タイムアウト設定
go test -timeout 30s ./...
```

### カバレッジ測定

```bash
# カバレッジ測定
go test -cover ./...

# 詳細なカバレッジレポート生成
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# パッケージ別カバレッジ
go test -coverprofile=coverage.out -coverpkg=./... ./...
go tool cover -func=coverage.out

# 目標カバレッジチェック
go test -cover ./... | grep -E "coverage: [0-9]+\.[0-9]+%" | awk -F"coverage: " '{print $2}' | awk -F"%" '{if($1 < 80) print "Coverage " $1 "% is below 80%"}'
```

### Makefile テストタスク

```makefile
# Makefile
.PHONY: test test-unit test-integration test-e2e test-coverage

test:
	go test ./...

test-unit:
	go test -short ./...

test-integration:
	go test -tags=integration ./...

test-e2e:
	go test -tags=e2e ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

test-coverage-check:
	@echo "Checking test coverage..."
	@go test -cover ./... | grep -E "coverage: [0-9]+\.[0-9]+%" | \
	awk -F"coverage: " '{print $$2}' | awk -F"%" '{if($$1 < 80) {print "❌ Coverage " $$1 "% is below 80%"; exit 1} else print "✅ Coverage " $$1 "% meets requirement"}'

bench:
	go test -bench=. ./...

test-race:
	go test -race ./...
```

## テスト品質ガイドライン

### テストケース設計原則

#### 1. AAA パターン（Arrange-Act-Assert）

```go
func TestCreateUser_ValidInput_ReturnsUser(t *testing.T) {
    // Arrange（準備）
    service := NewUserService()
    input := CreateUserInput{
        Name:  "Test User",
        Email: "test@example.com",
    }
    
    // Act（実行）
    result, err := service.CreateUser(context.Background(), input)
    
    // Assert（検証）
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, input.Name, result.Name().String())
}
```

#### 2. 境界値テスト

```go
func TestUserName_Validation_BoundaryValues(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        shouldError bool
    }{
        {"empty string", "", true},
        {"single character", "a", false},
        {"max length", strings.Repeat("a", 100), false},
        {"over max length", strings.Repeat("a", 101), true},
        {"unicode characters", "テストユーザー", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := NewUserName(tt.input)
            if tt.shouldError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

#### 3. エラーケーステスト

```go
func TestUserRepository_FindByID_NotFound_ReturnsError(t *testing.T) {
    // Given
    db := setupTestDB(t)
    repo := NewUserRepository(db)
    nonExistentID := NewUserID()
    
    // When
    user, err := repo.FindByID(context.Background(), nonExistentID)
    
    // Then
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.True(t, errors.Is(err, ErrUserNotFound))
}
```

### テスト可読性向上

#### 1. 意味のあるテスト名

```go
// ✅ 推奨: メソッド_条件_期待結果 パターン
func TestCreateUser_ValidEmail_ReturnsUser(t *testing.T) {}
func TestCreateUser_DuplicateEmail_ReturnsError(t *testing.T) {}
func TestCreateUser_InvalidEmailFormat_ReturnsValidationError(t *testing.T) {}

// ❌ 非推奨: 抽象的な名前
func TestCreateUser1(t *testing.T) {}
func TestCreateUserError(t *testing.T) {}
```

#### 2. Given-When-Then コメント

```go
func TestCreateUser_ValidInput_ReturnsUser(t *testing.T) {
    // Given（前提条件）
    repo := NewMockUserRepository()
    service := NewUserService(repo)
    input := CreateUserInput{Name: "Test", Email: "test@example.com"}
    
    // When（実行内容）
    result, err := service.CreateUser(context.Background(), input)
    
    // Then（期待結果）
    assert.NoError(t, err)
    assert.Equal(t, input.Name, result.Name().String())
}
```

## パフォーマンステスト

### ベンチマークテスト

```go
func BenchmarkUserRepository_Create(b *testing.B) {
    db := setupBenchmarkDB(b)
    repo := NewUserRepository(db)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        user := createTestUserWithID(i)
        repo.Create(context.Background(), user)
    }
}

func BenchmarkEmail_Validation(b *testing.B) {
    emails := generateTestEmails(1000)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        email := emails[i%len(emails)]
        NewEmail(email)
    }
}
```

### ロードテスト

```go
func TestUserAPI_ConcurrentRequests_HandlesLoad(t *testing.T) {
    server := setupTestServer(t)
    defer server.Close()
    
    concurrency := 100
    requestsPerGoroutine := 10
    
    var wg sync.WaitGroup
    errors := make(chan error, concurrency*requestsPerGoroutine)
    
    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func(routineID int) {
            defer wg.Done()
            
            for j := 0; j < requestsPerGoroutine; j++ {
                request := CreateUserRequest{
                    Name:  fmt.Sprintf("User%d-%d", routineID, j),
                    Email: fmt.Sprintf("user%d-%d@example.com", routineID, j),
                }
                
                resp, err := http.Post(
                    server.URL+"/users",
                    "application/json",
                    toJSONReader(request),
                )
                
                if err != nil || resp.StatusCode != http.StatusCreated {
                    errors <- fmt.Errorf("request failed: %v", err)
                }
            }
        }(i)
    }
    
    wg.Wait()
    close(errors)
    
    errorCount := 0
    for err := range errors {
        t.Logf("Error: %v", err)
        errorCount++
    }
    
    successRate := float64(concurrency*requestsPerGoroutine-errorCount) / float64(concurrency*requestsPerGoroutine)
    assert.Greater(t, successRate, 0.95, "Success rate should be above 95%")
}
```

## 継続的テスト

### GitHub Actions でのテスト自動化

```yaml
# .github/workflows/test.yml
name: Test

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:13
        env:
          POSTGRES_PASSWORD: password
          POSTGRES_DB: test_db
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run unit tests
      run: make test-unit
    
    - name: Run integration tests
      run: make test-integration
      env:
        DB_HOST: localhost
        DB_PASSWORD: password
    
    - name: Generate coverage report
      run: make test-coverage
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
```

### テストカバレッジ目標

| コンポーネント | 目標カバレッジ | 説明 |
|----------------|----------------|------|
| ドメイン層 | 90%以上 | ビジネスロジックの品質確保 |
| ユースケース層 | 85%以上 | アプリケーションロジックの信頼性 |
| インフラ層 | 70%以上 | 外部依存の基本動作確認 |
| プレゼンテーション層 | 75%以上 | API契約の遵守確認 |
| 全体 | 80%以上 | プロジェクト全体の品質基準 |

## 最終更新

- バージョン: 2025.05.24
- 更新者: AI Assistant
