# パフォーマンスルール

## 概要

Go DDDプロジェクトにおけるパフォーマンス最適化のガイドラインと実装パターンを定義します。

## 基本原則

### パフォーマンス設計原則
- **測定ファースト**: プロファイリング結果に基づく最適化
- **ボトルネック特定**: 実際の性能問題箇所の特定
- **段階的最適化**: 大きな改善から小さな最適化へ
- **可読性とのバランス**: 可読性を犠牲にしない最適化

## データベースパフォーマンス

### コネクションプール設定

```go
func ConfigureDB(config DatabaseConfig) (*sql.DB, error) {
    db, err := sql.Open("postgres", config.DSN)
    if err != nil {
        return nil, err
    }
    
    // コネクションプール設定
    db.SetMaxOpenConns(25)           // 最大接続数
    db.SetMaxIdleConns(5)            // アイドル接続数
    db.SetConnMaxLifetime(5 * time.Minute) // 接続最大存続時間
    db.SetConnMaxIdleTime(1 * time.Minute) // アイドル最大時間
    
    return db, nil
}
```

### インデックス設計

```sql
-- ユーザーテーブルのインデックス
CREATE INDEX CONCURRENTLY idx_users_email ON users(email);
CREATE INDEX CONCURRENTLY idx_users_created_at ON users(created_at);
CREATE INDEX CONCURRENTLY idx_users_active ON users(is_active) WHERE is_active = true;

-- 複合インデックス
CREATE INDEX CONCURRENTLY idx_orders_user_status ON orders(user_id, status);
CREATE INDEX CONCURRENTLY idx_orders_created_at_status ON orders(created_at, status);
```

### クエリ最適化

```go
// ✅ 効率的なページネーション
type PaginationQuery struct {
    Limit  int
    Offset int
    SortBy string
    Order  string
}

func (r *UserRepository) FindWithPagination(ctx context.Context, query PaginationQuery) ([]*User, error) {
    // インデックスを活用したクエリ
    sqlQuery := `
        SELECT id, name, email, created_at 
        FROM users 
        WHERE is_active = true
        ORDER BY created_at DESC, id DESC
        LIMIT $1 OFFSET $2
    `
    
    rows, err := r.db.QueryContext(ctx, sqlQuery, query.Limit, query.Offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    return r.scanUsers(rows)
}

// ✅ バッチ処理による効率化
func (r *UserRepository) CreateBatch(ctx context.Context, users []*User) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO users (id, name, email, created_at) 
        VALUES ($1, $2, $3, $4)
    `)
    if err != nil {
        return err
    }
    defer stmt.Close()
    
    for _, user := range users {
        _, err := stmt.ExecContext(ctx, 
            user.ID(), user.Name(), user.Email(), user.CreatedAt())
        if err != nil {
            return err
        }
    }
    
    return tx.Commit()
}
```

### N+1問題の回避

```go
// ❌ N+1問題のある実装
func (s *OrderService) GetOrdersWithUsers(ctx context.Context) ([]*OrderWithUser, error) {
    orders, err := s.orderRepo.FindAll(ctx)
    if err != nil {
        return nil, err
    }
    
    var result []*OrderWithUser
    for _, order := range orders {
        // 各注文に対してユーザー情報を個別取得（N+1問題）
        user, err := s.userRepo.FindByID(ctx, order.UserID())
        if err != nil {
            return nil, err
        }
        result = append(result, NewOrderWithUser(order, user))
    }
    
    return result, nil
}

// ✅ JOIN またはバッチ取得による解決
func (s *OrderService) GetOrdersWithUsers(ctx context.Context) ([]*OrderWithUser, error) {
    orders, err := s.orderRepo.FindAll(ctx)
    if err != nil {
        return nil, err
    }
    
    // ユーザーIDを収集
    userIDs := make([]UserID, 0, len(orders))
    for _, order := range orders {
        userIDs = append(userIDs, order.UserID())
    }
    
    // バッチでユーザー情報取得
    users, err := s.userRepo.FindByIDs(ctx, userIDs)
    if err != nil {
        return nil, err
    }
    
    // マップ化してO(1)でアクセス
    userMap := make(map[string]*User)
    for _, user := range users {
        userMap[user.ID().String()] = user
    }
    
    var result []*OrderWithUser
    for _, order := range orders {
        user := userMap[order.UserID().String()]
        result = append(result, NewOrderWithUser(order, user))
    }
    
    return result, nil
}
```

## キャッシュ戦略

### メモリキャッシュ実装

```go
type CacheService struct {
    cache sync.Map
    ttl   time.Duration
}

type CacheItem struct {
    Value     interface{}
    ExpiresAt time.Time
}

func (c *CacheService) Set(key string, value interface{}) {
    item := CacheItem{
        Value:     value,
        ExpiresAt: time.Now().Add(c.ttl),
    }
    c.cache.Store(key, item)
}

func (c *CacheService) Get(key string) (interface{}, bool) {
    item, ok := c.cache.Load(key)
    if !ok {
        return nil, false
    }
    
    cacheItem := item.(CacheItem)
    if time.Now().After(cacheItem.ExpiresAt) {
        c.cache.Delete(key)
        return nil, false
    }
    
    return cacheItem.Value, true
}
```

### Redisキャッシュ実装

```go
type RedisCache struct {
    client *redis.Client
    ttl    time.Duration
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    
    return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
    data, err := r.client.Get(ctx, key).Result()
    if err != nil {
        return err
    }
    
    return json.Unmarshal([]byte(data), dest)
}

// キャッシュ戦略の実装例
type CachedUserRepository struct {
    repo  UserRepository
    cache Cache
    ttl   time.Duration
}

func (r *CachedUserRepository) FindByID(ctx context.Context, id UserID) (*User, error) {
    key := fmt.Sprintf("user:%s", id.String())
    
    // キャッシュから取得試行
    var user User
    if err := r.cache.Get(ctx, key, &user); err == nil {
        return &user, nil
    }
    
    // リポジトリから取得
    user, err := r.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // キャッシュに保存
    _ = r.cache.Set(ctx, key, user)
    
    return user, nil
}
```

## HTTPパフォーマンス

### レスポンス最適化

```go
// gzip圧縮
func GzipMiddleware() echo.MiddlewareFunc {
    return middleware.GzipWithConfig(middleware.GzipConfig{
        Level: 5, // 圧縮レベル調整
    })
}

// レスポンスサイズ最適化
type UserListResponse struct {
    Users []UserSummary `json:"users"`
    Total int           `json:"total"`
}

type UserSummary struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    // 不要なフィールドは除外
}

func (u User) ToSummary() UserSummary {
    return UserSummary{
        ID:    u.ID().String(),
        Name:  u.Name().String(),
        Email: u.Email().String(),
    }
}
```

### コネクション最適化

```go
func NewHTTPServer(config ServerConfig) *echo.Echo {
    e := echo.New()
    
    // HTTP サーバー設定
    e.Server.ReadTimeout = time.Duration(config.ReadTimeout) * time.Second
    e.Server.WriteTimeout = time.Duration(config.WriteTimeout) * time.Second
    e.Server.IdleTimeout = time.Duration(config.IdleTimeout) * time.Second
    
    // Keep-Alive設定
    e.Server.SetKeepAlivesEnabled(true)
    
    return e
}
```

## Goランタイム最適化

### ガベージコレクション最適化

```go
// メモリプール使用例
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 0, 1024)
    },
}

func processData(data []byte) ([]byte, error) {
    buffer := bufferPool.Get().([]byte)
    defer bufferPool.Put(buffer[:0])
    
    // データ処理
    result := append(buffer, data...)
    // 処理結果をコピーして返す
    return append([]byte(nil), result...), nil
}

// ガベージコレクション調整
func init() {
    // GC頻度調整（本番環境では慎重に設定）
    debug.SetGCPercent(100)
}
```

### Goroutine管理

```go
// ワーカープール実装
type WorkerPool struct {
    workerCount int
    jobQueue    chan Job
    quit        chan bool
}

type Job func() error

func NewWorkerPool(workerCount int, bufferSize int) *WorkerPool {
    return &WorkerPool{
        workerCount: workerCount,
        jobQueue:    make(chan Job, bufferSize),
        quit:        make(chan bool),
    }
}

func (p *WorkerPool) Start() {
    for i := 0; i < p.workerCount; i++ {
        go p.worker()
    }
}

func (p *WorkerPool) worker() {
    for {
        select {
        case job := <-p.jobQueue:
            job()
        case <-p.quit:
            return
        }
    }
}

func (p *WorkerPool) Submit(job Job) {
    p.jobQueue <- job
}

func (p *WorkerPool) Stop() {
    close(p.quit)
}
```

## パフォーマンス測定

### ベンチマークテスト

```go
func BenchmarkUserRepository_FindByID(b *testing.B) {
    repo := setupTestRepository(b)
    userID := createTestUser(b, repo)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := repo.FindByID(context.Background(), userID)
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkUserService_CreateUser_Parallel(b *testing.B) {
    service := setupTestService(b)
    
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            input := CreateUserInput{
                Name:  "Test User",
                Email: fmt.Sprintf("test%d@example.com", rand.Int()),
            }
            _, err := service.CreateUser(context.Background(), input)
            if err != nil {
                b.Fatal(err)
            }
        }
    })
}
```

### プロファイリング

```go
import _ "net/http/pprof"

func main() {
    // pprof エンドポイント有効化（開発環境のみ）
    if os.Getenv("ENV") == "development" {
        go func() {
            log.Println(http.ListenAndServe("localhost:6060", nil))
        }()
    }
    
    // アプリケーション起動
    startServer()
}

// CPUプロファイリング実行例
func enableCPUProfiling() {
    f, err := os.Create("cpu.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
}
```

## パフォーマンステスト

### 負荷テスト

```go
func TestConcurrentUserCreation(t *testing.T) {
    service := setupTestService(t)
    concurrency := 100
    requestCount := 1000
    
    var wg sync.WaitGroup
    errors := make(chan error, requestCount)
    
    start := time.Now()
    
    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            
            for j := 0; j < requestCount/concurrency; j++ {
                input := CreateUserInput{
                    Name:  fmt.Sprintf("User%d-%d", workerID, j),
                    Email: fmt.Sprintf("user%d-%d@example.com", workerID, j),
                }
                
                _, err := service.CreateUser(context.Background(), input)
                if err != nil {
                    errors <- err
                    return
                }
            }
        }(i)
    }
    
    wg.Wait()
    close(errors)
    
    duration := time.Since(start)
    
    // エラーチェック
    for err := range errors {
        t.Error(err)
    }
    
    // パフォーマンス検証
    throughput := float64(requestCount) / duration.Seconds()
    t.Logf("Throughput: %.2f requests/second", throughput)
    
    if throughput < 100 { // 最低要求性能
        t.Errorf("Throughput too low: %.2f requests/second", throughput)
    }
}
```

## 監視・メトリクス

### パフォーマンスメトリクス収集

```go
type PerformanceMetrics struct {
    RequestDuration   prometheus.HistogramVec
    RequestCount      prometheus.CounterVec
    ActiveConnections prometheus.Gauge
    DBConnectionPool  prometheus.GaugeVec
}

func NewPerformanceMetrics() *PerformanceMetrics {
    return &PerformanceMetrics{
        RequestDuration: *prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "http_request_duration_seconds",
                Help: "HTTP request duration in seconds",
            },
            []string{"method", "endpoint", "status"},
        ),
        RequestCount: *prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "http_requests_total",
                Help: "Total number of HTTP requests",
            },
            []string{"method", "endpoint", "status"},
        ),
    }
}

// メトリクス収集ミドルウェア
func (m *PerformanceMetrics) Middleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            start := time.Now()
            
            err := next(c)
            
            duration := time.Since(start)
            status := c.Response().Status
            
            m.RequestDuration.WithLabelValues(
                c.Request().Method,
                c.Path(),
                fmt.Sprintf("%d", status),
            ).Observe(duration.Seconds())
            
            m.RequestCount.WithLabelValues(
                c.Request().Method,
                c.Path(),
                fmt.Sprintf("%d", status),
            ).Inc()
            
            return err
        }
    }
}
```

## パフォーマンス改善チェックリスト

### データベース層
- [ ] 適切なインデックス設計
- [ ] クエリ実行計画の確認
- [ ] N+1問題の回避
- [ ] コネクションプール設定
- [ ] バッチ処理の活用

### アプリケーション層
- [ ] キャッシュ戦略の実装
- [ ] Goroutine の適切な管理
- [ ] メモリプールの活用
- [ ] 不要なメモリアロケーション削減

### HTTP層
- [ ] レスポンス圧縮
- [ ] Keep-Alive設定
- [ ] レスポンスサイズ最適化
- [ ] 適切なタイムアウト設定

### 監視・測定
- [ ] ベンチマークテスト実装
- [ ] プロファイリング設定
- [ ] メトリクス収集
- [ ] 負荷テスト実行

## 参考資料

- [Go Performance Tips](https://github.com/golang/go/wiki/Performance)
- [Database Performance Best Practices](https://use-the-index-luke.com/)
- [Redis Performance Guidelines](https://redis.io/docs/management/optimization/)

## 最終更新

- 日付: 2025.05.24
- 更新者: AI Assistant
