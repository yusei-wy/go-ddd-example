# AI生成コードレビュー基準

## 概要

AI（GitHub Copilot、ChatGPT等）によって生成されたコードをレビューする際の専用チェックリストです。AI生成コードの特性を理解し、品質とセキュリティを確保するためのガイドラインを提供します。

## AI生成コードの特徴と注意点

### 一般的な特徴
- コード品質が一般的に高い
- ベストプラクティスに準拠している場合が多い
- 文法エラーが少ない
- パターンマッチングが得意

### 注意すべき点
- ビジネスロジックの誤解
- セキュリティ要件の見落とし
- プロジェクト固有のルールへの不適合
- 過度に複雑なソリューション
- ライセンス・著作権の問題

## レイヤー別レビューチェックリスト

### 1. プレゼンテーション層（Handler）

#### ✅ 基本チェック項目
- [ ] HTTPステータスコードが適切に設定されているか
- [ ] レスポンス構造がAPI設計ルールに準拠しているか
- [ ] エラーハンドリングが適切に実装されているか
- [ ] ログ出力が必要な箇所で実装されているか

#### ✅ セキュリティチェック
- [ ] 入力バリデーションが実装されているか
- [ ] SQLインジェクション対策が適切か
- [ ] XSS対策が実装されているか
- [ ] 認証・認可チェックが適切か
- [ ] レート制限が考慮されているか

#### ✅ プロジェクト固有チェック
- [ ] カスタムエラー型（HandlerError）を使用しているか
- [ ] 命名規則（PascalCase/camelCase）に準拠しているか
- [ ] DIパターンに従った実装になっているか

```go
// ✅ 良い例 - AI生成コードの適切な修正
func (h *UserHandler) CreateUser(c echo.Context) error {
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        return NewHandlerError(ErrInvalidRequest, "invalid request format")
    }
    
    // バリデーション実装
    if err := h.validator.Validate(req); err != nil {
        return NewHandlerError(ErrValidationFailed, err.Error())
    }
    
    // ログ出力
    h.logger.Info(c.Request().Context(), "create user request received",
        Field("email", req.Email),
    )
    
    // ユースケース実行
    output, err := h.createUserUsecase.Execute(c.Request().Context(), CreateUserInput{
        Name:  req.Name,
        Email: req.Email,
    })
    if err != nil {
        return err
    }
    
    return c.JSON(http.StatusCreated, CreateUserResponse{ID: output.ID})
}

// ❌ 避けるべき例 - AI生成コードの一般的な問題
func (h *UserHandler) CreateUser(c echo.Context) error {
    var req CreateUserRequest
    c.Bind(&req) // エラーハンドリング不足
    
    // 直接データベースアクセス（レイヤー違反）
    user := &User{Name: req.Name, Email: req.Email}
    h.db.Create(user) // エラーハンドリング不足
    
    return c.JSON(200, user) // ステータスコード不適切
}
```

### 2. アプリケーション層（Usecase）

#### ✅ 基本チェック項目
- [ ] ビジネスロジックがドメイン層に委譲されているか
- [ ] トランザクション管理が適切に実装されているか
- [ ] エラーハンドリングが適切か
- [ ] 入力・出力の型定義が明確か

#### ✅ CQRS準拠チェック
- [ ] Command/Queryの責務分離が適切か
- [ ] 読み取り専用操作で副作用がないか
- [ ] 書き込み操作で適切なバリデーションが実装されているか

#### ✅ プロジェクト固有チェック
- [ ] カスタムエラー型（UsecaseError）を使用しているか
- [ ] DI用のインターフェース依存になっているか
- [ ] ログ出力が適切に実装されているか

```go
// ✅ 良い例
func (u *CreateUserUsecase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
    // ビジネスルールバリデーション
    if exists := u.userRepo.ExistsByEmail(ctx, input.Email); exists {
        return nil, NewUsecaseError(ErrEmailAlreadyExists, "email already exists")
    }
    
    // ドメインオブジェクト生成（ドメイン層に委譲）
    user, err := NewUser(input.Name, input.Email)
    if err != nil {
        return nil, err
    }
    
    // 永続化
    if err := u.userRepo.Create(ctx, user); err != nil {
        return nil, NewUsecaseError(ErrUserCreationFailed, "failed to create user")
    }
    
    u.logger.Info(ctx, "user created successfully",
        Field("user_id", user.ID()),
        Field("email", user.Email()),
    )
    
    return &CreateUserOutput{ID: user.ID()}, nil
}
```

### 3. ドメイン層

#### ✅ 基本チェック項目
- [ ] ビジネスルールが適切にモデル化されているか
- [ ] 不変条件が保証されているか
- [ ] 外部依存がインターフェースで抽象化されているか
- [ ] 値オブジェクトが適切に実装されているか

#### ✅ DDD準拠チェック
- [ ] エンティティ・値オブジェクトの区別が適切か
- [ ] ドメインサービスの責務が明確か
- [ ] アグリゲート境界が適切に設計されているか

```go
// ✅ 良い例 - 値オブジェクトの適切な実装
type Email struct {
    value string
}

func NewEmail(value string) (Email, error) {
    if value == "" {
        return Email{}, NewModelError(ErrEmailRequired, "email is required")
    }
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

### 4. インフラストラクチャ層

#### ✅ 基本チェック項目
- [ ] リポジトリインターフェースが適切に実装されているか
- [ ] データベースアクセスが最適化されているか
- [ ] エラーハンドリングが適切か
- [ ] セキュリティ対策が実装されているか

#### ✅ セキュリティチェック
- [ ] SQLインジェクション対策（prepared statement）が実装されているか
- [ ] 接続プールが適切に設定されているか
- [ ] トランザクション管理が適切か

## AI特有の問題パターンと対策

### 1. 過度に複雑なソリューション

#### 問題例
```go
// ❌ AI生成コードでよくある過度に複雑な実装
func (r *UserRepository) FindByComplexCriteria(ctx context.Context, criteria SearchCriteria) ([]*User, error) {
    query := "SELECT * FROM users WHERE 1=1"
    args := []interface{}{}
    
    if criteria.Name != "" {
        query += " AND name ILIKE $" + strconv.Itoa(len(args)+1)
        args = append(args, "%"+criteria.Name+"%")
    }
    
    if criteria.Email != "" {
        query += " AND email = $" + strconv.Itoa(len(args)+1)
        args = append(args, criteria.Email)
    }
    
    // 過度に複雑なクエリ構築...
}
```

#### 改善案
```go
// ✅ シンプルで保守しやすい実装
func (r *UserRepository) FindByName(ctx context.Context, name string) ([]*User, error) {
    query := "SELECT id, name, email FROM users WHERE name ILIKE $1"
    // シンプルな実装
}

func (r *UserRepository) FindByEmail(ctx context.Context, email Email) (*User, error) {
    query := "SELECT id, name, email FROM users WHERE email = $1"
    // 目的に特化した実装
}
```

### 2. プロジェクト固有ルールの無視

#### よくある問題
- [ ] カスタムエラー型を使用していない
- [ ] 命名規則に準拠していない
- [ ] ディレクトリ構造ルールに従っていない
- [ ] ログ出力が不適切

#### 対策
1. **プロンプトエンジニアリング**: プロジェクト固有のルールを明示
2. **段階的レビュー**: 生成→修正→再レビューのサイクル
3. **テンプレート活用**: `prompt-templates.md`の活用

### 3. セキュリティホールの見落とし

#### チェックポイント
- [ ] 入力値検証の実装
- [ ] SQLインジェクション対策
- [ ] 認証・認可の実装
- [ ] ログ出力での機密情報漏洩防止
- [ ] エラーメッセージでの情報漏洩防止

## テストコードのレビュー

### ✅ AI生成テストの一般的チェック
- [ ] テストケースの網羅性が十分か
- [ ] エッジケースがカバーされているか
- [ ] モック・スタブが適切に使用されているか
- [ ] テストデータが適切に管理されているか

### ✅ TDD準拠チェック
- [ ] Given-When-Then構造になっているか
- [ ] テスト名が仕様を表現しているか
- [ ] テストが独立して実行可能か

```go
// ✅ 良い例 - AI生成テストの適切な構造
func TestNewEmail_ValidEmail_ReturnsEmail(t *testing.T) {
    // Given
    validEmail := "test@example.com"
    
    // When
    email, err := NewEmail(validEmail)
    
    // Then
    assert.NoError(t, err)
    assert.Equal(t, validEmail, email.String())
}

func TestNewEmail_InvalidFormat_ReturnsError(t *testing.T) {
    // Given
    invalidEmail := "invalid-email"
    
    // When
    _, err := NewEmail(invalidEmail)
    
    // Then
    assert.Error(t, err)
    assert.IsType(t, ModelError{}, err)
}
```

## コードメトリクス

### 自動チェック項目
- [ ] 循環複雑度（10以下を目標）
- [ ] 関数の行数（50行以下を目標）
- [ ] ネストレベル（4以下を目標）
- [ ] テストカバレッジ（80%以上を目標）

### ツール活用
```bash
# コード品質チェック
make lint

# テストカバレッジ測定
make test-coverage

# セキュリティスキャン
make security-scan
```

## 承認フロー

### 1. 自動チェック
- [ ] CI/CDパイプラインでの自動テスト
- [ ] 静的解析ツールでの品質チェック
- [ ] セキュリティスキャン

### 2. 人的レビュー
- [ ] ビジネスロジックの正確性確認
- [ ] アーキテクチャルール準拠確認
- [ ] セキュリティ要件確認

### 3. 最終承認
- [ ] テクニカルリードによる最終確認
- [ ] デプロイ承認

## 継続的改善

### レビュー結果の記録
- よくある問題パターンの蓄積
- プロンプト改善のためのフィードバック
- ルール更新の必要性評価

### AI支援ツールの活用
- GitHub Copilot の設定最適化
- カスタムプロンプトテンプレートの改善
- 自動化可能な箇所の特定

## まとめ

AI生成コードは強力なツールですが、プロジェクト固有の要件やビジネスロジックの理解には限界があります。このチェックリストを活用して、AI の利点を最大化しつつ、品質とセキュリティを確保したコードレビューを実施してください。

## 参考資料

- [プロンプトテンプレート](./prompt-templates.md)
- [Copilot利用指針](./copilot-instructions.md)
- [コーディング規約](../rules/coding-standards.md)
- [セキュリティガイドライン](../rules/security-guidelines.md)
