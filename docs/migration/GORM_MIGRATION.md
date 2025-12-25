# BaseInfo GORM 实现更新说明

## 更新时间
2025年12月25日

## 更新内容

### 1. 数据层改用 GORM 框架

将 `internal/data/baseinfo.go` 从原生 SQL 改为使用 GORM ORM 框架实现。

### 2. 查询接口优化

将 `FindByCodeAndDate(code, date)` 改为 `FindByCode(code)`，查询逻辑优化为：
- 只需要传入股票代码
- 自动返回该股票最新的一条记录（按 trade_date 降序）

### 3. 依赖新增

添加了以下依赖包：

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

## 技术实现

### Data 层双数据库支持

`internal/data/stock.go` 中的 `Data` 结构体现在同时支持原生 SQL 和 GORM：

```go
type Data struct {
    db     *sql.DB      // 原生 SQL（用于 stock）
    gormDB *gorm.DB     // GORM（用于 baseinfo）
    log    *log.Helper
}
```

### BaseInfo GORM 模型

```go
type BaseInfoModel struct {
    ID                         int64     `gorm:"column:id;primaryKey;autoIncrement"`
    StockName                  string    `gorm:"column:stock_name;type:varchar(100)"`
    LatestPrice                float64   `gorm:"column:latest_price;type:decimal(10,2)"`
    // ... 其他字段
    TradeDate                  time.Time `gorm:"column:trade_date;type:date;index:idx_stock_date"`
    Code                       string    `gorm:"column:code;type:varchar(50);uniqueIndex:uk_code_date"`
    CreateTime                 time.Time `gorm:"column:create_time;type:datetime;autoCreateTime"`
    UpdateTime                 time.Time `gorm:"column:update_time;type:datetime;autoUpdateTime"`
    // ... 更多字段
}

func (BaseInfoModel) TableName() string {
    return "base_info"
}
```

### FindByCode 实现

```go
func (r *baseInfoRepo) FindByCode(ctx context.Context, code string) (*biz.BaseInfo, error) {
    var model BaseInfoModel
    err := r.data.gormDB.WithContext(ctx).
        Where("code = ?", code).
        Order("trade_date DESC").  // 按日期降序，获取最新记录
        First(&model).Error

    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil  // 未找到记录返回 nil
        }
        return nil, err
    }

    return r.toBiz(&model), nil
}
```

### BatchSave GORM Upsert

使用 GORM 的 `Clauses` 实现批量 upsert：

```go
err := r.data.gormDB.WithContext(ctx).
    Clauses(clause.OnConflict{
        Columns: []clause.Column{{Name: "code"}, {Name: "trade_date"}},
        DoUpdates: clause.AssignmentColumns([]string{
            "stock_name", "latest_price", "auction_change_rate", 
            // ... 更多字段
        }),
    }).
    CreateInBatches(models, 100).Error
```

## 优势对比

### 使用 GORM 的优势

1. **代码简洁**: 减少了大量 SQL 语句编写
2. **类型安全**: 编译时检查字段类型
3. **自动映射**: 自动处理 struct 和数据库字段的映射
4. **链式调用**: 更优雅的查询构建方式
5. **上下文支持**: 原生支持 context
6. **事务处理**: 更简单的事务管理
7. **错误处理**: 统一的错误类型（如 `gorm.ErrRecordNotFound`）

### 原生 SQL vs GORM

| 特性 | 原生 SQL | GORM |
|------|---------|------|
| 代码量 | 较多 | 较少 |
| 性能 | 最优 | 略低（可忽略） |
| 灵活性 | 最高 | 高 |
| 类型安全 | 需手动处理 | 自动处理 |
| 学习曲线 | 低 | 中 |
| 维护成本 | 高 | 低 |

## 接口变更

### Biz 层接口

**变更前：**
```go
type BaseInfoRepo interface {
    FindByCodeAndDate(ctx context.Context, code string, date time.Time) (*BaseInfo, error)
}
```

**变更后：**
```go
type BaseInfoRepo interface {
    FindByCode(ctx context.Context, code string) (*BaseInfo, error)
}
```

## 使用示例

### 查询最新基础数据

```go
// 按股票代码查询最新记录
info, err := baseInfoRepo.FindByCode(ctx, "000001")
if err != nil {
    log.Error("查询失败:", err)
    return
}

if info == nil {
    log.Info("未找到记录")
    return
}

log.Infof("股票: %s, 最新价: %.2f, 日期: %s", 
    info.StockName, 
    info.LatestPrice, 
    info.TradeDate.Format("2006-01-02"))
```

### 批量保存数据

```go
infos := []*biz.BaseInfo{
    {Code: "000001", StockName: "平安银行", TradeDate: time.Now()},
    {Code: "000002", StockName: "万科A", TradeDate: time.Now()},
}

err := baseInfoRepo.BatchSave(ctx, infos)
if err != nil {
    log.Error("批量保存失败:", err)
}
```

## 数据库连接管理

### 初始化

程序启动时会创建两个数据库连接：

```go
// 原生 SQL 连接（用于 stock 表操作）
db, err := sql.Open(c.Database.Driver, c.Database.Source)

// GORM 连接（用于 base_info 表操作）
gormDB, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
```

### 清理

程序关闭时会同时关闭两个连接：

```go
cleanup := func() {
    db.Close()                // 关闭原生 SQL 连接
    sqlDB, _ := gormDB.DB()
    if sqlDB != nil {
        sqlDB.Close()         // 关闭 GORM 连接
    }
}
```

## 迁移指南

如果你有现有代码使用了 `FindByCodeAndDate`，需要进行以下调整：

### 变更前

```go
date := time.Now()
info, err := repo.FindByCodeAndDate(ctx, "000001", date)
```

### 变更后

```go
// 现在只需要传入 code，自动返回最新记录
info, err := repo.FindByCode(ctx, "000001")
```

## 性能优化建议

1. **索引优化**: 确保 `code` 和 `trade_date` 字段有合适的索引
2. **批量操作**: 使用 `CreateInBatches` 时，批量大小建议 100-1000
3. **查询优化**: 只查询需要的字段，使用 `.Select()` 方法
4. **连接池**: GORM 会自动管理连接池，默认配置已足够

## 注意事项

1. **数据库兼容**: 当前实现仅支持 MySQL
2. **时区处理**: 确保配置中包含 `parseTime=True&loc=Local`
3. **NULL 值**: GORM 会自动处理 NULL 值为零值
4. **软删除**: 当前未启用 GORM 软删除功能
5. **自动迁移**: 未使用 `AutoMigrate`，表结构通过 SQL 文件管理

## 测试

编译成功确认：

```bash
$ go build -o cmd/goweicai/goweicai cmd/goweicai/main.go
# 无错误输出，编译成功
```

## 后续优化方向

1. 考虑将 stock 表也迁移到 GORM
2. 添加更多查询方法（如按日期范围查询）
3. 实现分页查询支持
4. 添加查询缓存层
5. 性能监控和慢查询日志
