# セキュリティガイドライン

## 概要

Go DDD プロジェクトにおけるセキュリティ要件とベストプラクティスを定義します。

## 基本原則

### セキュリティバイデザイン
- 設計段階からセキュリティを考慮
- 最小権限の原則
- 多層防御の実装
- 機密情報の適切な管理

## 入力検証・サニタイゼーション

### SQLインジェクション対策

```go
// ✅ 正しい例：prepared statement使用
func (r *UserRepository) FindByEmail(ctx context.Context, email Email) (*User, error) {
    query := "SELECT id, name, email FROM users WHERE email = $1"
    row := r.db.QueryRowContext(ctx, query, email.String())
    // ...
}

// ❌ 危険な例：文字列結合
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
    query := "SELECT id, name, email FROM users WHERE email = '" + email + "'"
    // SQLインジェクションの脆弱性
}
```

### XSS対策

```go
import "html"

// HTMLエスケープ処理
func sanitizeInput(input string) string {
    return html.EscapeString(input)
}

// レスポンス時のエスケープ
type UserResponse struct {
    Name string `json:"name"`
}

func (u User) ToResponse() UserResponse {
    return UserResponse{
        Name: html.EscapeString(u.Name()),
    }
}
```

### 入力値検証

```go
// 値オブジェクトでの検証強化
func NewEmail(value string) (Email, error) {
    // 基本形式チェック
    if !emailRegex.MatchString(value) {
        return Email{}, NewModelError(ErrInvalidEmail, "invalid email format")
    }
    
    // 長さ制限
    if len(value) > MaxEmailLength {
        return Email{}, NewModelError(ErrEmailTooLong, "email too long")
    }
    
    // 危険文字チェック
    if containsDangerousChars(value) {
        return Email{}, NewModelError(ErrDangerousChars, "contains dangerous characters")
    }
    
    return Email{value: value}, nil
}
```

## 認証・認可

### JWT認証実装

```go
type AuthService struct {
    secretKey []byte
    expiry    time.Duration
}

func (s *AuthService) GenerateToken(userID UserID) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID.String(),
        "exp":     time.Now().Add(s.expiry).Unix(),
        "iat":     time.Now().Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.secretKey)
}

func (s *AuthService) ValidateToken(tokenString string) (*UserClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return s.secretKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return parseUserClaims(claims)
    }
    
    return nil, fmt.Errorf("invalid token")
}
```

### 認可ミドルウェア

```go
func AuthorizeMiddleware(requiredRole Role) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            user := getUserFromContext(c)
            if user == nil {
                return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
            }
            
            if !user.HasRole(requiredRole) {
                return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
            }
            
            return next(c)
        }
    }
}
```

## レート制限

### APIレート制限

```go
func RateLimitMiddleware() echo.MiddlewareFunc {
    store := middleware.NewRateLimiterMemoryStore(10) // 10 requests per second
    
    return middleware.RateLimiter(store)
}

// IPベースレート制限
func IPRateLimitMiddleware() echo.MiddlewareFunc {
    store := middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
        Rate:      10,
        Burst:     20,
        ExpiresIn: time.Minute,
    })
    
    return middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
        Store: store,
        IdentifierExtractor: func(c echo.Context) (string, error) {
            return c.RealIP(), nil
        },
    })
}
```

## 機密情報管理

### 環境変数による設定

```go
type SecurityConfig struct {
    JWTSecret       string `env:"JWT_SECRET,required"`
    DBPassword      string `env:"DB_PASSWORD,required"`
    APIKey          string `env:"API_KEY,required"`
    EncryptionKey   string `env:"ENCRYPTION_KEY,required"`
}

// 機密情報のマスキング
func (c SecurityConfig) String() string {
    return fmt.Sprintf("SecurityConfig{JWTSecret: %s, DBPassword: %s}", 
        maskSecret(c.JWTSecret), maskSecret(c.DBPassword))
}

func maskSecret(secret string) string {
    if len(secret) <= 4 {
        return "****"
    }
    return secret[:2] + "****" + secret[len(secret)-2:]
}
```

### 暗号化処理

```go
type EncryptionService struct {
    key []byte
}

func (s *EncryptionService) Encrypt(plaintext string) (string, error) {
    block, err := aes.NewCipher(s.key)
    if err != nil {
        return "", err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }
    
    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}
```

## セキュリティヘッダー

### HTTPセキュリティヘッダー

```go
func SecurityHeadersMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // XSS Protection
            c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
            
            // Content Type Options
            c.Response().Header().Set("X-Content-Type-Options", "nosniff")
            
            // Frame Options
            c.Response().Header().Set("X-Frame-Options", "DENY")
            
            // HSTS
            c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
            
            // CSP
            c.Response().Header().Set("Content-Security-Policy", "default-src 'self'")
            
            return next(c)
        }
    }
}
```

## ログセキュリティ

### 機密情報のログ除外

```go
type SecureLogger struct {
    logger Logger
}

func (l *SecureLogger) LogUserAction(ctx context.Context, action string, user User) {
    // パスワードやトークンなどの機密情報を除外
    sanitizedUser := struct {
        ID    string `json:"id"`
        Email string `json:"email"`
        // パスワードフィールドは除外
    }{
        ID:    user.ID().String(),
        Email: maskEmail(user.Email().String()),
    }
    
    l.logger.Info(ctx, "user action",
        Field("action", action),
        Field("user", sanitizedUser),
    )
}

func maskEmail(email string) string {
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return "****@****"
    }
    return parts[0][:1] + "****@" + parts[1]
}
```

## 依存関係セキュリティ

### 依存関係の脆弱性チェック

```bash
# Makefileに追加
.PHONY: security-check
security-check:
	go list -json -deps ./... | nancy sleuth
	gosec ./...
	go mod verify
```

### セキュアな依存関係管理

```go
// go.mod でのバージョン固定
require (
    github.com/golang-jwt/jwt/v5 v5.0.0
    golang.org/x/crypto v0.10.0
    // 脆弱性のないバージョンを明示的に指定
)
```

## セキュリティテスト

### セキュリティテストの実装

```go
func TestPasswordHashing_SecureImplementation(t *testing.T) {
    password := "testpassword123"
    
    // パスワードハッシュ化テスト
    hasher := NewPasswordHasher()
    hashed, err := hasher.Hash(password)
    
    assert.NoError(t, err)
    assert.NotEqual(t, password, hashed)
    assert.True(t, hasher.Verify(password, hashed))
    
    // 同じパスワードでも異なるハッシュになることを確認（salt使用）
    hashed2, _ := hasher.Hash(password)
    assert.NotEqual(t, hashed, hashed2)
}

func TestJWTToken_SecurityValidation(t *testing.T) {
    authService := NewAuthService("test-secret", time.Hour)
    userID := NewUserID()
    
    // トークン生成
    token, err := authService.GenerateToken(userID)
    assert.NoError(t, err)
    
    // 有効なトークンの検証
    claims, err := authService.ValidateToken(token)
    assert.NoError(t, err)
    assert.Equal(t, userID.String(), claims.UserID)
    
    // 無効なトークンの検証
    _, err = authService.ValidateToken("invalid-token")
    assert.Error(t, err)
}
```

## 脆弱性対応プロセス

### 脆弱性発見時の対応

1. **即座の影響評価**
   - 脆弱性の深刻度評価（CVSS スコア）
   - 影響範囲の特定
   - 緊急度の判定

2. **緊急対応**
   - セキュリティパッチの適用
   - 一時的な緩和策の実施
   - インシデント記録

3. **恒久対策**
   - 根本原因分析
   - セキュリティテストの追加
   - プロセス改善

### セキュリティ監査チェックリスト

- [ ] 入力値検証の実装確認
- [ ] SQL インジェクション対策確認
- [ ] XSS 対策確認
- [ ] 認証・認可の実装確認
- [ ] 機密情報の適切な管理確認
- [ ] ログの機密情報除外確認
- [ ] セキュリティヘッダーの設定確認
- [ ] 依存関係の脆弱性チェック
- [ ] セキュリティテストの実装確認

## 参考資料

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go Security Checklist](https://github.com/Checkmarx/Go-SCP)
- [JWT Best Practices](https://datatracker.ietf.org/doc/html/draft-ietf-oauth-jwt-bcp)

## 最終更新

- 日付: 2025.05.24
- 更新者: AI Assistant
