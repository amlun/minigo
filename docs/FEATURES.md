# 框架功能说明

本文档介绍 minigo 框架的核心功能和使用方法。

## 1. 统一响应格式

所有 API 响应都遵循统一的 JSON 格式：

### 成功响应
```json
{
  "success": true,
  "code": "SUCCESS",
  "message": "成功",
  "data": { ... }
}
```

### 错误响应
```json
{
  "success": false,
  "code": "USER_001",
  "message": "用户不存在"
}
```

### 分页响应
```json
{
  "success": true,
  "code": "SUCCESS",
  "message": "成功",
  "data": {
    "items": [...],
    "total": 100,
    "page": 1,
    "size": 10
  }
}
```

### 使用方法

```go
// 普通成功响应
response.Ok(c, userData)

// 自定义消息
response.OkWithMessage(c, "创建成功", userData)

// 分页响应
response.OkWithPage(c, userList, total, page, pageSize)

// 错误响应（自动识别错误类型）
response.HandleError(c, err)
```

## 2. 分页查询

### 使用内置分页 DTO

```go
import "minigo/internal/interfaces/dto"

type ListUserRequest struct {
    dto.PageSortRequest  // 嵌入分页和排序参数
    Keyword string `form:"keyword" json:"keyword"`
}

func (h *Handler) List(c *gin.Context) {
    var req ListUserRequest
    if !middleware.ValidateAndBindQuery(c, &req) {
        return
    }

    // 自动处理默认值和上下限
    page := req.GetPage()         // 默认 1
    pageSize := req.GetPageSize() // 默认 10，最大 100
    offset := req.GetOffset()     // 计算偏移量

    // 查询数据
    users, total, err := service.List(ctx, page, pageSize)

    // 返回分页响应
    response.OkWithPage(c, users, total, page, pageSize)
}
```

### 请求示例

```
GET /api/users?page=2&pageSize=20&sortBy=created_at&sortOrder=desc&keyword=test
```

## 3. 查询过滤器

框架提供了强大的查询构建器，支持链式调用。

```go
import "minigo/pkg/query"

// 基础用法
qb := query.NewQueryBuilder().
    Where("status", "=", 1).
    Like("name", keyword).
    DateRange("created_at", startTime, endTime).
    Order("created_at", true). // true = DESC
    Paginate(page, pageSize)

// 应用到查询
query := db.NewSelect().Model(&users)
query = qb.Apply(query)
err := query.Scan(ctx)
```

### 可用的过滤器

- `Where(column, operator, value)` - WHERE 条件
- `In(column, values)` - IN 条件
- `Like(column, value)` - LIKE 模糊查询
- `DateRange(column, start, end)` - 日期范围
- `Range(column, min, max)` - 数值范围
- `Order(column, desc)` - 排序
- `Paginate(page, size)` - 分页

## 4. 错误处理

### 错误类型

框架定义了五种错误类型：

- `SystemError` - 系统错误（数据库、网络等）
- `BusinessError` - 业务错误（余额不足、重复等）
- `ValidationError` - 验证错误（参数格式等）
- `AuthError` - 认证错误（未登录、权限不足）
- `NotFoundError` - 资源不存在

### 创建错误

```go
import apperrors "minigo/internal/domain/errors"

// 使用预定义的错误
return apperrors.ErrUserNotFound

// 创建自定义错误
return apperrors.NewBusinessError("INSUFFICIENT_BALANCE", "余额不足")
return apperrors.NewValidationError("INVALID_PHONE", "手机号格式错误")
return apperrors.NewAuthError("TOKEN_EXPIRED", "令牌已过期")
```

### 处理错误

```go
// 在 Handler 中
if err != nil {
    response.HandleError(c, err)  // 自动识别错误类型并返回相应的HTTP状态码
    return
}

// 或者使用更简洁的方式
if !middleware.MustNotError(c, err) {
    return
}

// 中止请求并返回错误
middleware.AbortWithError(c, err)
```

## 5. 请求追踪

每个请求自动分配一个唯一的 Request ID，用于日志追踪和问题排查。

### 自动处理

Request ID 会：
- 自动生成（UUID v4）
- 记录在所有日志中
- 通过 `X-Request-ID` 响应头返回给客户端
- 可通过请求头传入继续使用

### 使用方法

```go
// 在 Handler 中获取请求 ID
requestID := middleware.GetRequestID(c)

// 在日志中会自动包含
logging.L().WithFields(map[string]interface{}{
    "request_id": requestID,
}).Info("processing request")
```

### 客户端使用

```bash
# 服务端会返回 Request ID
curl -i http://localhost:8808/api/users
# X-Request-ID: 550e8400-e29b-41d4-a716-446655440000

# 可以在后续请求中传入
curl -H "X-Request-ID: 550e8400-e29b-41d4-a716-446655440000" \
     http://localhost:8808/api/users/1
```

## 6. 参数验证

### 内置验证器

```go
import "minigo/pkg/validator"

// 格式验证
validator.IsEmail(email)
validator.IsPhone(phone)
validator.IsIDCard(idCard)
validator.IsURL(url)

// 长度验证
validator.LengthBetween(username, 3, 20)
validator.MinLength(password, 6)
validator.MaxLength(bio, 500)

// 类型验证
validator.IsNumeric(code)
validator.IsAlpha(name)
validator.IsAlphanumeric(username)

// 范围验证
validator.InRange(age, 18, 100)
validator.InSlice(status, []string{"active", "inactive"})

// 密码强度
validator.IsStrongPassword(password) // 至少8位，包含大小写字母和数字
validator.ValidatePassword(password, 8, true, true, true, false)
```

### Struct 标签验证

```go
type CreateUserRequest struct {
    Phone    string `json:"phone" binding:"required"`
    Password string `json:"password" binding:"required,min=6"`
    Name     string `json:"name" binding:"required,min=2,max=50"`
    Email    string `json:"email" binding:"omitempty,email"`
    Age      int    `json:"age" binding:"omitempty,min=18,max=120"`
}

// 在 Handler 中
if !middleware.ValidateAndBindJSON(c, &req) {
    return // 自动返回验证错误
}
```

### 辅助函数

```go
// 验证并绑定 JSON
middleware.ValidateAndBindJSON(c, &req)

// 验证并绑定 Query 参数
middleware.ValidateAndBindQuery(c, &req)

// 验证并绑定 URI 参数
middleware.ValidateAndBindURI(c, &req)

// 验证 ID 参数
id, ok := middleware.ValidateIDParam(c, "id")
if !ok {
    return
}

// 验证必需字符串
if !middleware.ValidateRequiredString(c, name, "name") {
    return
}
```

## 7. 工具函数

### 类型转换

```go
import "minigo/pkg/utils"

// 字符串转换
i64, err := utils.ToInt64("12345")
i64 := utils.ToInt64OrDefault("invalid", 0)
i := utils.ToIntOrDefault("123", 0)
f := utils.ToFloat64OrDefault("3.14", 0.0)
b := utils.ToBoolOrDefault("true", false)

// 转为字符串
s := utils.Int64ToString(12345)
s := utils.IntToString(123)
s := utils.ToString(anyValue)

// 切片转换
ids, err := utils.SplitToInt64Slice("1,2,3,4", ",")
ints, err := utils.SplitToIntSlice("1,2,3", ",")
str := utils.Int64SliceToString([]int64{1,2,3}, ",")

// Struct 和 Map 互转
m, err := utils.StructToMap(user)
err := utils.MapToStruct(m, &user)
```

### 切片操作

```go
import "minigo/pkg/utils"

// 基础操作
exists := utils.Contains(slice, item)
unique := utils.Unique(slice)
filtered := utils.Filter(slice, func(item T) bool { return item > 0 })
mapped := utils.Map(slice, func(item T) R { return transform(item) })

// 查找
item, found := utils.Find(slice, func(item T) bool { return item.ID == 1 })
index := utils.FindIndex(slice, predicate)

// 判断
hasAny := utils.Any(slice, predicate)
allMatch := utils.All(slice, predicate)

// 转换
chunks := utils.Chunk(slice, 10)         // 分块
reversed := utils.Reverse(slice)          // 反转
flattened := utils.Flatten(nested)        // 展平

// 集合操作
intersection := utils.Intersection(s1, s2)
difference := utils.Difference(s1, s2)
union := utils.Union(s1, s2)

// 分组
grouped := utils.GroupBy(slice, func(item T) K { return item.Category })
trueSlice, falseSlice := utils.Partition(slice, predicate)
```

### 时间处理

```go
import "minigo/pkg/utils"

// 时间范围
tr := utils.TimeRange{Start: start, End: end}
err := tr.Validate()
empty := tr.IsEmpty()

// 解析日期
tr, err := utils.ParseTimeRange("2024-01-01", "2024-12-31")
date, err := utils.ParseDate("2024-01-01")
```

## 8. 事务管理

```go
// 在 Service 层使用事务
err := s.txManager.InTx(ctx, func(ctx context.Context) error {
    // 这个上下文中的所有仓储操作都在同一事务中

    // 更新用户
    if err := s.userRepo.Update(ctx, user); err != nil {
        return err // 自动回滚
    }

    // 创建订单
    if err := s.orderRepo.Create(ctx, order); err != nil {
        return err // 自动回滚
    }

    // 返回 nil 表示提交
    return nil
})
```

## 9. 日志记录

```go
import "minigo/internal/infrastructure/logging"

// 使用全局 logger
logging.L().Info("user logged in")
logging.L().WithFields(map[string]interface{}{
    "user_id": userID,
    "action": "login",
}).Info("user action")

logging.L().WithError(err).Error("operation failed")

// 不同日志级别
logging.L().Debug("debug info")
logging.L().Info("info message")
logging.L().Warn("warning")
logging.L().Error("error occurred")
```

## 10. 完整示例

### 创建一个 CRUD 接口

```go
// 1. 定义 DTO
type CreateUserRequest struct {
    Phone    string `json:"phone" binding:"required"`
    Password string `json:"password" binding:"required,min=6"`
    Name     string `json:"name" binding:"required,min=2,max=50"`
}

type ListUserRequest struct {
    dto.PageSortRequest
    Keyword string `form:"keyword"`
    Status  *int   `form:"status"`
}

// 2. Handler
func (h *UserHandler) Create(c *gin.Context) {
    var req CreateUserRequest
    if !middleware.ValidateAndBindJSON(c, &req) {
        return
    }

    // 额外验证
    if !validator.IsPhone(req.Phone) {
        middleware.AbortWithBusinessError(c, "INVALID_PHONE", "手机号格式错误")
        return
    }

    user, err := h.service.Create(c.Request.Context(), &req)
    if !middleware.MustNotError(c, err) {
        return
    }

    response.OkWithMessage(c, "创建成功", user)
}

func (h *UserHandler) List(c *gin.Context) {
    var req ListUserRequest
    if !middleware.ValidateAndBindQuery(c, &req) {
        return
    }

    users, total, err := h.service.List(c.Request.Context(), &req)
    if !middleware.MustNotError(c, err) {
        return
    }

    response.OkWithPage(c, users, total, req.GetPage(), req.GetPageSize())
}

func (h *UserHandler) GetByID(c *gin.Context) {
    id, ok := middleware.ValidateIDParam(c, "id")
    if !ok {
        return
    }

    user, err := h.service.GetByID(c.Request.Context(), id)
    if !middleware.MustNotError(c, err) {
        return
    }

    response.Ok(c, user)
}

// 3. Service
func (s *UserService) List(ctx context.Context, req *ListUserRequest) ([]*entity.User, int64, error) {
    // 构建查询
    qb := query.NewQueryBuilder()

    if req.Keyword != "" {
        qb.Like("name", req.Keyword)
    }
    if req.Status != nil {
        qb.Where("status", "=", *req.Status)
    }

    qb.Order(req.GetSortBy("created_at"), req.GetSortOrder() == dto.SortOrderDesc)
    qb.Paginate(req.GetPage(), req.GetPageSize())

    return s.userRepo.List(ctx, qb)
}

// 4. Repository
func (r *UserRepository) List(ctx context.Context, qb *query.QueryBuilder) ([]*entity.User, int64, error) {
    db := dbctx.FromCtx(ctx, r.DB)

    var users []*entity.User
    query := db.NewSelect().Model(&users)

    // 应用过滤器
    query = qb.Apply(query)

    // 获取总数（分页前）
    total, err := query.Count(ctx)
    if err != nil {
        return nil, 0, ConvertQueryError(err)
    }

    // 执行查询
    err = query.Scan(ctx)
    if err != nil {
        return nil, 0, ConvertQueryError(err)
    }

    return users, int64(total), nil
}
```

## 总结

minigo 框架提供了一套完整的工具集，让你可以快速开发高质量的 Web API：

- **统一的响应格式** - 标准化接口输出
- **分页和过滤** - 开箱即用的查询功能
- **错误处理** - 类型化的错误管理
- **请求追踪** - 完整的日志追踪链路
- **参数验证** - 丰富的验证工具
- **工具函数** - 常用操作的辅助函数
- **事务管理** - 优雅的事务处理
- **日志记录** - 结构化的日志输出

只需专注于业务逻辑，框架会处理好其他一切！
