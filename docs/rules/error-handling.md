# エラーハンドリング規約

## 概要

Go DDD プロジェクトにおけるエラーハンドリングの統一的な実装方法とベストプラクティスを定義します。

## 基本原則

### エラーハンドリング原則
- **レイヤー固有**: 各レイヤーで適切なエラー型を使用
- **コンテキスト保持**: エラーの発生箇所と原因を明確化
- **一貫性**: 統一されたエラー処理パターン
- **ログとの連携**: 適切なログレベルでエラー情報を記録

## エラー型階層

### レイヤー別エラー定義

```go
// 基底エラー型
type BaseError struct {
    Code    string
    Message string
    Cause   error
}

func (e BaseError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Cause)
    }
    return e.Message
}

func (e BaseError) Unwrap() error {
    return e.Cause
}

// ハンドラーエラー
type HandlerError struct {
    BaseError
    StatusCode int
}

func NewHandlerError(code string, message string, statusCode int) HandlerError {
    return HandlerError{
        BaseError:  BaseError{Code: code, Message: message},
        StatusCode: statusCode,
    }
}

func WrapHandlerError(code string, message string, statusCode int, cause error) HandlerError {
    return HandlerError{
        BaseError:  BaseError{Code: code, Message: message, Cause: cause},
        StatusCode: statusCode,
    }
}

// ユースケースエラー
type UsecaseError struct {
    BaseError
}

func NewUsecaseError(code string, message string) UsecaseError {
    return UsecaseError{
        BaseError: BaseError{Code: code, Message: message},
    }
}

func WrapUsecaseError(code string, message string, cause error) UsecaseError {
    return UsecaseError{
        BaseError: BaseError{Code: code, Message: message, Cause: cause},
    }
}

// サービスエラー
type ServiceError struct {
    BaseError
}

func NewServiceError(code string, message string) ServiceError {
    return ServiceError{
        BaseError: BaseError{Code: code, Message: message},
    }
}

func WrapServiceError(code string, message string, cause error) ServiceError {
    return ServiceError{
        BaseError: BaseError{Code: code, Message: message, Cause: cause},
    }
}

// リポジトリエラー
type RepositoryError struct {
    BaseError
}

func NewRepositoryError(code string, message string) RepositoryError {
    return RepositoryError{
        BaseError: BaseError{Code: code, Message: message},
    }
}

func WrapRepositoryError(code string, message string, cause error) RepositoryError {
    return RepositoryError{
        BaseError: BaseError{Code: code, Message: message, Cause: cause},
    }
}

// モデルエラー
type ModelError struct {
    BaseError
}

func NewModelError(code string, message string) ModelError {
    return ModelError{
        BaseError: BaseError{Code: code, Message: message},
    }
}
```

### エラーコード定数

```go
// ハンドラーエラーコード
const (
    ErrInvalidRequest     = "INVALID_REQUEST"
    ErrInvalidPathParam   = "INVALID_PATH_PARAM"
    ErrInvalidQueryParam  = "INVALID_QUERY_PARAM"
    ErrMissingHeader      = "MISSING_HEADER"
    ErrInvalidContentType = "INVALID_CONTENT_TYPE"
)

// ユースケースエラーコード
const (
    ErrUserNotFound         = "USER_NOT_FOUND"
    ErrEmailAlreadyExists   = "EMAIL_ALREADY_EXISTS"
    ErrInsufficientPermission = "INSUFFICIENT_PERMISSION"
    ErrResourceConflict     = "RESOURCE_CONFLICT"
    ErrBusinessRuleViolation = "BUSINESS_RULE_VIOLATION"
)

// サービスエラーコード
const (
    ErrDomainLogicError   = "DOMAIN_LOGIC_ERROR"
    ErrInvalidOperation   = "INVALID_OPERATION"
    ErrServiceUnavailable = "SERVICE_UNAVAILABLE"
)

// リポジトリエラーコード
const (
    ErrDatabaseConnection = "DATABASE_CONNECTION_ERROR"
    ErrDatabaseQuery      = "DATABASE_QUERY_ERROR"
    ErrDuplicateKey       = "DUPLICATE_KEY_ERROR"
    ErrForeignKeyViolation = "FOREIGN_KEY_VIOLATION"
    ErrOptimisticLock     = "OPTIMISTIC_LOCK_ERROR"
)

// モデルエラーコード
const (
    ErrInvalidEmail    = "INVALID_EMAIL"
    ErrInvalidUserName = "INVALID_USER_NAME"
    ErrInvalidPassword = "INVALID_PASSWORD"
    ErrInvalidUserID   = "INVALID_USER_ID"
)
```

## レイヤー別エラーハンドリング

### ハンドラー層

```go
type UserHandler struct {
    logger            Logger
    createUserUsecase CreateUserUsecase
}

func (h *UserHandler) CreateUser(c echo.Context) error {
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        handlerErr := WrapHandlerError(
            ErrInvalidRequest,
            "Invalid request format",
            http.StatusBadRequest,
            err,
        )
        h.logError(c.Request().Context(), handlerErr)
        return h.mapToHTTPError(handlerErr)
    }
    
    input := CreateUserInput{
        Name:  req.Name,
        Email: req.Email,
    }
    
    output, err := h.createUserUsecase.Execute(c.Request().Context(), input)
    if err != nil {
        handlerErr := h.mapUsecaseError(err)
        h.logError(c.Request().Context(), handlerErr)
        return h.mapToHTTPError(handlerErr)
    }
    
    response := CreateUserResponse{
        ID:   output.ID,
        Name: output.Name,
    }
    
    return c.JSON(http.StatusCreated, response)
}

func (h *UserHandler) mapUsecaseError(err error) HandlerError {
    var usecaseErr UsecaseError
    if errors.As(err, &usecaseErr) {
        switch usecaseErr.Code {
        case ErrUserNotFound:
            return NewHandlerError(ErrUserNotFound, "User not found", http.StatusNotFound)
        case ErrEmailAlreadyExists:
            return NewHandlerError(ErrEmailAlreadyExists, "Email already exists", http.StatusConflict)
        default:
            return NewHandlerError("USECASE_ERROR", "Business logic error", http.StatusBadRequest)
        }
    }
    
    return NewHandlerError("INTERNAL_ERROR", "Internal server error", http.StatusInternalServerError)
}

func (h *UserHandler) mapToHTTPError(err HandlerError) error {
    return echo.NewHTTPError(err.StatusCode, map[string]interface{}{
        "error": map[string]interface{}{
            "code":    err.Code,
            "message": err.Message,
        },
    })
}

func (h *UserHandler) logError(ctx context.Context, err error) {
    h.logger.Error(ctx, "handler error",
        Field("error", err.Error()),
        Field("type", fmt.Sprintf("%T", err)),
    )
}
```

### ユースケース層

```go
type CreateUserUsecase struct {
    userRepo UserRepository
    logger   Logger
}

func (u *CreateUserUsecase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
    // 重複チェック
    exists, err := u.userRepo.ExistsByEmail(ctx, input.Email)
    if err != nil {
        return nil, WrapUsecaseError(
            "EMAIL_CHECK_FAILED",
            "Failed to check email existence",
            err,
        )
    }
    
    if exists {
        return nil, NewUsecaseError(
            ErrEmailAlreadyExists,
            fmt.Sprintf("Email %s already exists", input.Email),
        )
    }
    
    // ドメインオブジェクト生成
    user, err := NewUser(input.Name, input.Email)
    if err != nil {
        return nil, WrapUsecaseError(
            "USER_CREATION_FAILED",
            "Failed to create user domain object",
            err,
        )
    }
    
    // 永続化
    if err := u.userRepo.Create(ctx, user); err != nil {
        return nil, WrapUsecaseError(
            "USER_SAVE_FAILED",
            "Failed to save user",
            err,
        )
    }
    
    u.logger.Info(ctx, "user created successfully",
        Field("user_id", user.ID().String()),
        Field("email", user.Email().String()),
    )
    
    return &CreateUserOutput{
        ID:   user.ID().String(),
        Name: user.Name().String(),
    }, nil
}
```

### サービス層

```go
type UserService struct {
    userRepo UserRepository
    logger   Logger
}

func (s *UserService) ValidateUserPermission(ctx context.Context, userID UserID, resource string) error {
    user, err := s.userRepo.FindByID(ctx, userID)
    if err != nil {
        return WrapServiceError(
            "USER_FETCH_FAILED",
            "Failed to fetch user for permission check",
            err,
        )
    }
    
    if !user.HasPermission(resource) {
        return NewServiceError(
            ErrInsufficientPermission,
            fmt.Sprintf("User %s does not have permission for %s", userID.String(), resource),
        )
    }
    
    return nil
}

func (s *UserService) DeactivateUser(ctx context.Context, userID UserID) error {
    user, err := s.userRepo.FindByID(ctx, userID)
    if err != nil {
        return WrapServiceError(
            "USER_FETCH_FAILED",
            "Failed to fetch user for deactivation",
            err,
        )
    }
    
    if err := user.Deactivate(); err != nil {
        return WrapServiceError(
            "USER_DEACTIVATION_FAILED",
            "Failed to deactivate user",
            err,
        )
    }
    
    if err := s.userRepo.Update(ctx, user); err != nil {
        return WrapServiceError(
            "USER_UPDATE_FAILED",
            "Failed to update user status",
            err,
        )
    }
    
    s.logger.Info(ctx, "user deactivated",
        Field("user_id", userID.String()),
    )
    
    return nil
}
```

### リポジトリ層

```go
type postgresUserRepository struct {
    db *sql.DB
}

func (r *postgresUserRepository) Create(ctx context.Context, user *User) error {
    query := `
        INSERT INTO users (id, name, email, password_hash, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
    
    _, err := r.db.ExecContext(ctx, query,
        user.ID().String(),
        user.Name().String(),
        user.Email().String(),
        user.PasswordHash(),
        user.CreatedAt(),
        user.UpdatedAt(),
    )
    
    if err != nil {
        if isUniqueViolation(err) {
            return NewRepositoryError(
                ErrDuplicateKey,
                "User with this email already exists",
            )
        }
        return WrapRepositoryError(
            ErrDatabaseQuery,
            "Failed to create user in database",
            err,
        )
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
            return nil, NewRepositoryError(
                ErrUserNotFound,
                fmt.Sprintf("User with ID %s not found", id.String()),
            )
        }
        return nil, WrapRepositoryError(
            ErrDatabaseQuery,
            "Failed to query user from database",
            err,
        )
    }
    
    user, err := userModel.ToDomain()
    if err != nil {
        return nil, WrapRepositoryError(
            "USER_RECONSTRUCTION_FAILED",
            "Failed to reconstruct user domain object",
            err,
        )
    }
    
    return user, nil
}

// データベースエラー判定ヘルパー
func isUniqueViolation(err error) bool {
    var pgErr *pq.Error
    return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func isForeignKeyViolation(err error) bool {
    var pgErr *pq.Error
    return errors.As(err, &pgErr) && pgErr.Code == "23503"
}
```

### ドメイン層

```go
// ユーザー値オブジェクト
func NewEmail(value string) (Email, error) {
    if value == "" {
        return Email{}, NewModelError(
            ErrInvalidEmail,
            "Email cannot be empty",
        )
    }
    
    if !emailRegex.MatchString(value) {
        return Email{}, NewModelError(
            ErrInvalidEmail,
            fmt.Sprintf("Invalid email format: %s", value),
        )
    }
    
    if len(value) > MaxEmailLength {
        return Email{}, NewModelError(
            ErrInvalidEmail,
            fmt.Sprintf("Email too long: %d characters (max %d)", len(value), MaxEmailLength),
        )
    }
    
    return Email{value: value}, nil
}

func NewUserName(value string) (UserName, error) {
    if value == "" {
        return UserName{}, NewModelError(
            ErrInvalidUserName,
            "User name cannot be empty",
        )
    }
    
    if len(value) > MaxUserNameLength {
        return UserName{}, NewModelError(
            ErrInvalidUserName,
            fmt.Sprintf("User name too long: %d characters (max %d)", len(value), MaxUserNameLength),
        )
    }
    
    if containsInvalidChars(value) {
        return UserName{}, NewModelError(
            ErrInvalidUserName,
            "User name contains invalid characters",
        )
    }
    
    return UserName{value: value}, nil
}

// ユーザーエンティティ
func (u *User) ChangeEmail(newEmail Email) error {
    if u.email.Equals(newEmail) {
        return NewModelError(
            "EMAIL_UNCHANGED",
            "New email is the same as current email",
        )
    }
    
    u.email = newEmail
    u.updatedAt = time.Now()
    return nil
}

func (u *User) Deactivate() error {
    if !u.isActive {
        return NewModelError(
            "USER_ALREADY_INACTIVE",
            "User is already inactive",
        )
    }
    
    u.isActive = false
    u.updatedAt = time.Now()
    return nil
}
```

## エラーラッピングとアンラッピング

### エラーチェーンの管理

```go
// エラーの階層的ラッピング例
func (u *UpdateUserUsecase) Execute(ctx context.Context, input UpdateUserInput) error {
    user, err := u.userRepo.FindByID(ctx, input.ID)
    if err != nil {
        // リポジトリエラーをユースケースエラーでラップ
        return WrapUsecaseError(
            "USER_FETCH_FAILED",
            "Failed to fetch user for update",
            err, // 元のRepositoryErrorを保持
        )
    }
    
    if input.Email != "" {
        email, err := NewEmail(input.Email)
        if err != nil {
            // モデルエラーをユースケースエラーでラップ
            return WrapUsecaseError(
                "EMAIL_VALIDATION_FAILED",
                "Email validation failed",
                err, // 元のModelErrorを保持
            )
        }
        
        if err := user.ChangeEmail(email); err != nil {
            return WrapUsecaseError(
                "EMAIL_CHANGE_FAILED",
                "Failed to change user email",
                err,
            )
        }
    }
    
    if err := u.userRepo.Update(ctx, user); err != nil {
        return WrapUsecaseError(
            "USER_UPDATE_FAILED",
            "Failed to update user",
            err,
        )
    }
    
    return nil
}

// エラーチェーンの確認
func handleError(err error) {
    // 特定のエラー型チェック
    var modelErr ModelError
    if errors.As(err, &modelErr) {
        log.Printf("Model validation error: %s", modelErr.Code)
    }
    
    var repoErr RepositoryError
    if errors.As(err, &repoErr) {
        log.Printf("Repository error: %s", repoErr.Code)
    }
    
    // 特定のエラーコードチェック
    if errors.Is(err, NewModelError(ErrInvalidEmail, "")) {
        log.Println("Invalid email error detected")
    }
    
    // エラーチェーン全体の出力
    log.Printf("Full error chain: %+v", err)
}
```

## ログとの連携

### エラーログ出力

```go
type ErrorLogger struct {
    logger Logger
}

func (el *ErrorLogger) LogError(ctx context.Context, err error, additionalFields ...Field) {
    fields := []Field{
        Field("error", err.Error()),
        Field("error_type", fmt.Sprintf("%T", err)),
    }
    
    // エラー型別の追加情報
    switch e := err.(type) {
    case HandlerError:
        fields = append(fields,
            Field("status_code", e.StatusCode),
            Field("error_code", e.Code),
        )
    case UsecaseError:
        fields = append(fields,
            Field("error_code", e.Code),
        )
    case RepositoryError:
        fields = append(fields,
            Field("error_code", e.Code),
        )
    case ModelError:
        fields = append(fields,
            Field("error_code", e.Code),
        )
    }
    
    // エラーチェーンの情報
    if err := errors.Unwrap(err); err != nil {
        fields = append(fields, Field("cause", err.Error()))
    }
    
    fields = append(fields, additionalFields...)
    
    el.logger.Error(ctx, "error occurred", fields...)
}

// 使用例
func (h *UserHandler) CreateUser(c echo.Context) error {
    // ... 実装 ...
    
    if err != nil {
        h.errorLogger.LogError(c.Request().Context(), err,
            Field("endpoint", "POST /users"),
            Field("user_ip", c.RealIP()),
            Field("user_agent", c.Request().UserAgent()),
        )
        return h.mapToHTTPError(err)
    }
    
    // ...
}
```

## テストでのエラーハンドリング

### エラーケーステスト

```go
func TestCreateUserUsecase_EmailAlreadyExists_ReturnsError(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    usecase := NewCreateUserUsecase(mockRepo, logger)
    
    input := CreateUserInput{
        Name:  "Test User",
        Email: "test@example.com",
    }
    
    // リポジトリが重複を返すよう設定
    mockRepo.On("ExistsByEmail", mock.Anything, input.Email).Return(true, nil)
    
    // Act
    output, err := usecase.Execute(context.Background(), input)
    
    // Assert
    assert.Nil(t, output)
    assert.Error(t, err)
    
    var usecaseErr UsecaseError
    assert.True(t, errors.As(err, &usecaseErr))
    assert.Equal(t, ErrEmailAlreadyExists, usecaseErr.Code)
    assert.Contains(t, usecaseErr.Message, "already exists")
}

func TestUserHandler_CreateUser_InvalidRequest_ReturnsError(t *testing.T) {
    // Arrange
    handler := setupTestHandler(t)
    invalidJSON := `{"name": "", "email": "invalid"}`
    
    req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(invalidJSON))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := echo.New().NewContext(req, rec)
    
    // Act
    err := handler.CreateUser(c)
    
    // Assert
    var handlerErr HandlerError
    assert.True(t, errors.As(err, &handlerErr))
    assert.Equal(t, http.StatusBadRequest, handlerErr.StatusCode)
}
```

## エラーハンドリングチェックリスト

### 設計段階
- [ ] レイヤー別エラー型定義
- [ ] エラーコード体系設計
- [ ] エラーメッセージ規約
- [ ] ログ出力戦略

### 実装段階
- [ ] 適切なエラー型使用
- [ ] エラーラッピング実装
- [ ] ログ出力実装
- [ ] エラーマッピング実装

### テスト段階
- [ ] エラーケーステスト
- [ ] エラー型検証
- [ ] エラーメッセージ検証
- [ ] ログ出力検証

### 運用段階
- [ ] エラー監視設定
- [ ] アラート設定
- [ ] エラー分析ダッシュボード
- [ ] エラーレート監視

## 参考資料

- [Go Error Handling Best Practices](https://go.dev/blog/error-handling-and-go)
- [Effective Error Handling in Go](https://earthly.dev/blog/golang-errors/)
- [Error Wrapping in Go 1.13](https://go.dev/blog/go1.13-errors)
