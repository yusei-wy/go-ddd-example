# API設計ルール

## 概要

Go DDD プロジェクトにおけるREST API設計の規約とベストプラクティスを定義します。

## 基本原則

### API設計原則
- **一貫性**: 統一されたAPI設計パターン
- **直感性**: 開発者が理解しやすいエンドポイント設計
- **セキュリティ**: セキュリティ要件を満たした設計
- **拡張性**: 将来の機能追加を考慮した設計
- **パフォーマンス**: 効率的なデータ転送

## REST API規約

### HTTPメソッド使用規則

```
GET    /users          # ユーザー一覧取得
GET    /users/{id}     # 特定ユーザー取得
POST   /users          # ユーザー作成
PUT    /users/{id}     # ユーザー全体更新
PATCH  /users/{id}     # ユーザー部分更新
DELETE /users/{id}     # ユーザー削除
```

### エンドポイント命名規則

```
✅ 推奨パターン
GET    /api/v1/users
GET    /api/v1/users/{id}
GET    /api/v1/users/{id}/orders
POST   /api/v1/users/{id}/orders
GET    /api/v1/orders/{id}/items

❌ 避けるべきパターン
GET    /api/v1/getUsers
POST   /api/v1/createUser
GET    /api/v1/user-list
```

### リソース設計

```go
// ユーザーリソース
type UserResource struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// ユーザー一覧レスポンス
type UsersResponse struct {
    Users      []UserResource `json:"users"`
    Total      int            `json:"total"`
    Page       int            `json:"page"`
    PerPage    int            `json:"per_page"`
    TotalPages int            `json:"total_pages"`
}

// エラーレスポンス
type ErrorResponse struct {
    Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}
```

## リクエスト・レスポンス設計

### リクエスト構造

```go
// 作成リクエスト
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=1,max=100"`
    Email string `json:"email" validate:"required,email,max=255"`
}

// 更新リクエスト
type UpdateUserRequest struct {
    Name     *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
    Email    *string `json:"email,omitempty" validate:"omitempty,email,max=255"`
    IsActive *bool   `json:"is_active,omitempty"`
}

// クエリパラメータ
type UserListQuery struct {
    Page    int    `query:"page" validate:"min=1"`
    PerPage int    `query:"per_page" validate:"min=1,max=100"`
    Sort    string `query:"sort" validate:"oneof=name email created_at"`
    Order   string `query:"order" validate:"oneof=asc desc"`
    Search  string `query:"search" validate:"max=100"`
}
```

### レスポンス構造

```go
// 成功レスポンス
type SuccessResponse struct {
    Data    interface{} `json:"data"`
    Message string      `json:"message,omitempty"`
}

// ページネーション情報
type PaginationMeta struct {
    Page       int `json:"page"`
    PerPage    int `json:"per_page"`
    Total      int `json:"total"`
    TotalPages int `json:"total_pages"`
}

// 一覧レスポンス基底構造
type ListResponse struct {
    Data       interface{}    `json:"data"`
    Pagination PaginationMeta `json:"pagination"`
}
```

## HTTPステータスコード

### 標準ステータスコード使用規則

```go
const (
    // 2xx Success
    StatusOK            = 200 // GET, PUT, PATCH成功
    StatusCreated       = 201 // POST成功（リソース作成）
    StatusNoContent     = 204 // DELETE成功、PUT/PATCH成功（レスポンスボディなし）
    
    // 4xx Client Error
    StatusBadRequest          = 400 // リクエスト形式エラー
    StatusUnauthorized        = 401 // 認証エラー
    StatusForbidden           = 403 // 認可エラー
    StatusNotFound           = 404 // リソース未発見
    StatusMethodNotAllowed   = 405 // HTTPメソッドエラー
    StatusConflict           = 409 // リソース競合（重複等）
    StatusUnprocessableEntity = 422 // バリデーションエラー
    StatusTooManyRequests    = 429 // レート制限
    
    // 5xx Server Error
    StatusInternalServerError = 500 // サーバー内部エラー
    StatusBadGateway         = 502 // 上流サーバーエラー
    StatusServiceUnavailable = 503 // サービス利用不可
)
```

### エラーレスポンス実装

```go
type APIError struct {
    Code       string `json:"code"`
    Message    string `json:"message"`
    Details    string `json:"details,omitempty"`
    StatusCode int    `json:"-"`
}

func (e APIError) Error() string {
    return e.Message
}

// エラーマッピング
func MapDomainErrorToAPI(err error) APIError {
    switch {
    case errors.Is(err, domain.ErrUserNotFound):
        return APIError{
            Code:       "USER_NOT_FOUND",
            Message:    "User not found",
            StatusCode: http.StatusNotFound,
        }
    case errors.Is(err, domain.ErrEmailAlreadyExists):
        return APIError{
            Code:       "EMAIL_ALREADY_EXISTS",
            Message:    "Email address already exists",
            StatusCode: http.StatusConflict,
        }
    case errors.Is(err, domain.ErrInvalidEmail):
        return APIError{
            Code:       "INVALID_EMAIL",
            Message:    "Invalid email format",
            Details:    err.Error(),
            StatusCode: http.StatusBadRequest,
        }
    default:
        return APIError{
            Code:       "INTERNAL_ERROR",
            Message:    "Internal server error",
            StatusCode: http.StatusInternalServerError,
        }
    }
}
```

## バリデーション

### リクエストバリデーション

```go
import (
    "github.com/go-playground/validator/v10"
)

type RequestValidator struct {
    validator *validator.Validate
}

func NewRequestValidator() *RequestValidator {
    v := validator.New()
    
    // カスタムバリデーター登録
    v.RegisterValidation("user_id", validateUserID)
    v.RegisterValidation("order_status", validateOrderStatus)
    
    return &RequestValidator{validator: v}
}

func (rv *RequestValidator) Validate(i interface{}) error {
    if err := rv.validator.Struct(i); err != nil {
        return rv.mapValidationError(err)
    }
    return nil
}

func (rv *RequestValidator) mapValidationError(err error) error {
    var validationErrors []string
    
    for _, err := range err.(validator.ValidationErrors) {
        switch err.Tag() {
        case "required":
            validationErrors = append(validationErrors, 
                fmt.Sprintf("%s is required", err.Field()))
        case "email":
            validationErrors = append(validationErrors, 
                fmt.Sprintf("%s must be a valid email", err.Field()))
        case "min":
            validationErrors = append(validationErrors, 
                fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param()))
        default:
            validationErrors = append(validationErrors, 
                fmt.Sprintf("%s is invalid", err.Field()))
        }
    }
    
    return APIError{
        Code:       "VALIDATION_ERROR",
        Message:    "Validation failed",
        Details:    strings.Join(validationErrors, "; "),
        StatusCode: http.StatusUnprocessableEntity,
    }
}

// カスタムバリデーター
func validateUserID(fl validator.FieldLevel) bool {
    userID := fl.Field().String()
    _, err := uuid.Parse(userID)
    return err == nil
}
```

## ハンドラー実装

### ベースハンドラー

```go
type BaseHandler struct {
    logger    Logger
    validator *RequestValidator
}

func NewBaseHandler(logger Logger, validator *RequestValidator) BaseHandler {
    return BaseHandler{
        logger:    logger,
        validator: validator,
    }
}

// 成功レスポンス
func (h BaseHandler) Success(c echo.Context, statusCode int, data interface{}, message ...string) error {
    response := SuccessResponse{Data: data}
    if len(message) > 0 {
        response.Message = message[0]
    }
    return c.JSON(statusCode, response)
}

// エラーレスポンス
func (h BaseHandler) Error(c echo.Context, err error) error {
    apiErr := MapDomainErrorToAPI(err)
    
    h.logger.Error(c.Request().Context(), "API error",
        Field("error", err.Error()),
        Field("code", apiErr.Code),
        Field("path", c.Request().URL.Path),
        Field("method", c.Request().Method),
    )
    
    return c.JSON(apiErr.StatusCode, ErrorResponse{
        Error: ErrorDetail{
            Code:    apiErr.Code,
            Message: apiErr.Message,
            Details: apiErr.Details,
        },
    })
}

// リクエストバインド&バリデーション
func (h BaseHandler) BindAndValidate(c echo.Context, req interface{}) error {
    if err := c.Bind(req); err != nil {
        return APIError{
            Code:       "INVALID_REQUEST_FORMAT",
            Message:    "Invalid request format",
            StatusCode: http.StatusBadRequest,
        }
    }
    
    return h.validator.Validate(req)
}
```

### CRUDハンドラー実装

```go
type UserHandler struct {
    BaseHandler
    createUserUsecase CreateUserUsecase
    getUserUsecase    GetUserUsecase
    updateUserUsecase UpdateUserUsecase
    deleteUserUsecase DeleteUserUsecase
    listUsersUsecase  ListUsersUsecase
}

func NewUserHandler(
    logger Logger,
    validator *RequestValidator,
    createUserUsecase CreateUserUsecase,
    getUserUsecase GetUserUsecase,
    updateUserUsecase UpdateUserUsecase,
    deleteUserUsecase DeleteUserUsecase,
    listUsersUsecase ListUsersUsecase,
) *UserHandler {
    return &UserHandler{
        BaseHandler:       NewBaseHandler(logger, validator),
        createUserUsecase: createUserUsecase,
        getUserUsecase:    getUserUsecase,
        updateUserUsecase: updateUserUsecase,
        deleteUserUsecase: deleteUserUsecase,
        listUsersUsecase:  listUsersUsecase,
    }
}

// POST /users
func (h *UserHandler) CreateUser(c echo.Context) error {
    var req CreateUserRequest
    if err := h.BindAndValidate(c, &req); err != nil {
        return h.Error(c, err)
    }
    
    input := CreateUserInput{
        Name:  req.Name,
        Email: req.Email,
    }
    
    output, err := h.createUserUsecase.Execute(c.Request().Context(), input)
    if err != nil {
        return h.Error(c, err)
    }
    
    response := CreateUserResponse{
        ID:        output.ID,
        Name:      output.Name,
        Email:     output.Email,
        CreatedAt: output.CreatedAt,
    }
    
    return h.Success(c, http.StatusCreated, response, "User created successfully")
}

// GET /users/{id}
func (h *UserHandler) GetUser(c echo.Context) error {
    userID := c.Param("id")
    if userID == "" {
        return h.Error(c, APIError{
            Code:       "MISSING_USER_ID",
            Message:    "User ID is required",
            StatusCode: http.StatusBadRequest,
        })
    }
    
    input := GetUserInput{ID: userID}
    
    output, err := h.getUserUsecase.Execute(c.Request().Context(), input)
    if err != nil {
        return h.Error(c, err)
    }
    
    response := UserResource{
        ID:        output.ID,
        Name:      output.Name,
        Email:     output.Email,
        IsActive:  output.IsActive,
        CreatedAt: output.CreatedAt,
        UpdatedAt: output.UpdatedAt,
    }
    
    return h.Success(c, http.StatusOK, response)
}

// GET /users
func (h *UserHandler) ListUsers(c echo.Context) error {
    var query UserListQuery
    if err := c.Bind(&query); err != nil {
        return h.Error(c, err)
    }
    
    // デフォルト値設定
    if query.Page == 0 {
        query.Page = 1
    }
    if query.PerPage == 0 {
        query.PerPage = 20
    }
    if query.Sort == "" {
        query.Sort = "created_at"
    }
    if query.Order == "" {
        query.Order = "desc"
    }
    
    if err := h.validator.Validate(query); err != nil {
        return h.Error(c, err)
    }
    
    input := ListUsersInput{
        Page:    query.Page,
        PerPage: query.PerPage,
        Sort:    query.Sort,
        Order:   query.Order,
        Search:  query.Search,
    }
    
    output, err := h.listUsersUsecase.Execute(c.Request().Context(), input)
    if err != nil {
        return h.Error(c, err)
    }
    
    users := make([]UserResource, len(output.Users))
    for i, user := range output.Users {
        users[i] = UserResource{
            ID:        user.ID,
            Name:      user.Name,
            Email:     user.Email,
            IsActive:  user.IsActive,
            CreatedAt: user.CreatedAt,
            UpdatedAt: user.UpdatedAt,
        }
    }
    
    response := ListResponse{
        Data: users,
        Pagination: PaginationMeta{
            Page:       output.Page,
            PerPage:    output.PerPage,
            Total:      output.Total,
            TotalPages: output.TotalPages,
        },
    }
    
    return h.Success(c, http.StatusOK, response)
}
```

## セキュリティ

### 認証・認可

```go
// JWT認証ミドルウェア
func JWTAuthMiddleware(secret []byte) echo.MiddlewareFunc {
    return middleware.JWTWithConfig(middleware.JWTConfig{
        SigningKey: secret,
        Claims:     &UserClaims{},
        TokenLookup: "header:Authorization:Bearer ",
        ErrorHandler: func(err error) error {
            return APIError{
                Code:       "INVALID_TOKEN",
                Message:    "Invalid or expired token",
                StatusCode: http.StatusUnauthorized,
            }
        },
    })
}

// 権限チェックミドルウェア
func RequireRole(requiredRole string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            user := c.Get("user").(*jwt.Token)
            claims := user.Claims.(*UserClaims)
            
            if !contains(claims.Roles, requiredRole) {
                return APIError{
                    Code:       "INSUFFICIENT_PERMISSIONS",
                    Message:    "Insufficient permissions",
                    StatusCode: http.StatusForbidden,
                }
            }
            
            return next(c)
        }
    }
}

// リソース所有者チェック
func RequireResourceOwner() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            user := c.Get("user").(*jwt.Token)
            claims := user.Claims.(*UserClaims)
            resourceUserID := c.Param("id")
            
            if claims.UserID != resourceUserID && !contains(claims.Roles, "admin") {
                return APIError{
                    Code:       "ACCESS_DENIED",
                    Message:    "Access denied to this resource",
                    StatusCode: http.StatusForbidden,
                }
            }
            
            return next(c)
        }
    }
}
```

### レート制限

```go
func RateLimitMiddleware() echo.MiddlewareFunc {
    config := middleware.RateLimiterConfig{
        Store: middleware.NewRateLimiterMemoryStoreWithConfig(
            middleware.RateLimiterMemoryStoreConfig{
                Rate:      10,              // 10 requests per second
                Burst:     20,              // burst of 20 requests
                ExpiresIn: time.Minute,     // expire entries after 1 minute
            },
        ),
        IdentifierExtractor: func(c echo.Context) (string, error) {
            // IPアドレスベースの制限
            return c.RealIP(), nil
        },
        ErrorHandler: func(c echo.Context, err error) error {
            return APIError{
                Code:       "RATE_LIMIT_EXCEEDED",
                Message:    "Rate limit exceeded",
                StatusCode: http.StatusTooManyRequests,
            }
        },
    }
    
    return middleware.RateLimiterWithConfig(config)
}
```

## ドキュメント生成

### OpenAPI仕様

```go
// Swagger annotation例
// @Summary      Create user
// @Description  Create a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body CreateUserRequest true "User creation request"
// @Success      201 {object} CreateUserResponse
// @Failure      400 {object} ErrorResponse
// @Failure      409 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /users [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
    // 実装
}

// @Summary      Get user by ID
// @Description  Get user information by ID
// @Tags         users
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200 {object} UserResource
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c echo.Context) error {
    // 実装
}
```

## API設計チェックリスト

### 設計段階
- [ ] RESTful設計原則の準拠
- [ ] 一貫したエンドポイント命名
- [ ] 適切なHTTPメソッド使用
- [ ] 統一されたレスポンス構造
- [ ] エラーハンドリング戦略

### 実装段階
- [ ] リクエスト・レスポンス型定義
- [ ] バリデーション実装
- [ ] ハンドラー実装
- [ ] 認証・認可実装
- [ ] レート制限実装

### セキュリティ
- [ ] 入力値検証
- [ ] 認証実装
- [ ] 認可実装
- [ ] セキュリティヘッダー設定
- [ ] レート制限設定

### ドキュメント
- [ ] OpenAPI仕様作成
- [ ] APIドキュメント生成
- [ ] 利用例の記載
- [ ] エラーコード一覧

### テスト
- [ ] ハンドラーテスト
- [ ] 統合テスト
- [ ] セキュリティテスト
- [ ] パフォーマンステスト

## 参考資料

- [REST API Design Guidelines](https://restfulapi.net/)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Echo Framework Guide](https://echo.labstack.com/guide/)

## 最終更新

- 日付: 2025.05.24
- 更新者: AI Assistant
