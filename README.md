# minigo

一个基于 Clean Architecture/DDD 的 Go 工程模板，提供开箱即用的 Web API 开发框架。

## 特性

- **清晰的架构分层** - 遵循 Clean Architecture 和领域驱动设计原则
- **JWT 认证** - 内置用户认证和角色权限管理
- **事务管理** - 优雅的数据库事务处理机制
- **配置管理** - 支持环境变量和配置文件（.env）
- **中间件支持** - CORS、日志、错误处理、限流等
- **Snowflake ID** - 分布式唯一 ID 生成
- **PostgreSQL** - 使用 Bun ORM 操作数据库

## 技术栈

- **Web 框架**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [Bun](https://github.com/uptrace/bun) with PostgreSQL
- **认证**: [JWT](https://github.com/golang-jwt/jwt)
- **配置**: [Viper](https://github.com/spf13/viper) + [gotenv](https://github.com/subosito/gotenv)
- **日志**: [logrus](https://github.com/sirupsen/logrus)
- **ID 生成**: [Snowflake](https://github.com/bwmarrin/snowflake)

## 快速开始

### 环境要求

- Go 1.23+
- PostgreSQL 12+

### 安装依赖

```bash
go mod download
```

### 配置环境变量

复制示例配置文件并修改：

```bash
cp .env.example .env
```

编辑 `.env` 文件，配置数据库连接等信息：

```env
ENV=dev
PORT=8808
DB_DSN=postgres://postgres:postgres@localhost:5432/dbname?sslmode=disable
JWT_SECRET=your_secret_key_change_me
JWT_EXPIRE_DURATION=24h
LOG_LEVEL=debug
```

### 初始化数据库

```bash
# 创建数据库
psql -U postgres -c "CREATE DATABASE dbname;"

# 运行迁移脚本
psql -U postgres -d dbname -f migrations/001_init.sql
```

### 运行应用

```bash
# 直接运行
go run ./cmd/server

# 或使用 make
make run
```

应用将在 `http://localhost:8808` 启动。

## 开发命令

```bash
# 运行应用
make run

# 构建二进制文件（输出到 build/ 目录）
make build

# 运行测试
make test

# 整理依赖
make tidy

# 清理构建产物
make clean
```

## 项目结构

```
minigo/
├── cmd/
│   └── server/              # 应用入口
│       └── main.go
├── internal/
│   ├── domain/              # 域层（业务逻辑核心）
│   │   ├── entity/          # 实体模型
│   │   ├── repository/      # 仓储接口
│   │   ├── service/         # 域服务
│   │   └── errors/          # 域错误
│   ├── application/         # 应用层（用例编排）
│   │   └── service/         # 应用服务
│   ├── infrastructure/      # 基础设施层
│   │   ├── repository/      # 仓储实现
│   │   ├── config/          # 配置管理
│   │   ├── auth/            # JWT 认证
│   │   ├── tx/              # 事务管理
│   │   ├── logging/         # 日志
│   │   └── id/              # ID 生成器
│   └── interfaces/          # 接口层
│       ├── http/            # HTTP 处理器
│       ├── dto/             # 数据传输对象
│       ├── middleware/      # 中间件
│       └── response/        # 响应格式化
├── pkg/                     # 公共工具包
│   └── utils/
├── migrations/              # 数据库迁移文件
└── vendor/                  # 依赖包（vendored）
```

## API 端点

### 健康检查

```
GET /api/health
```

### 认证

```
POST /api/auth/login
```

请求体：
```json
{
  "phone": "13800138000",
  "password": "password123"
}
```

## 核心概念

### 架构分层

项目采用 Clean Architecture，分为四层：

1. **Domain（域层）** - 包含业务实体、仓储接口和域服务，不依赖其他层
2. **Application（应用层）** - 编排业务用例，调用域服务和仓储
3. **Infrastructure（基础设施层）** - 实现技术细节（数据库、认证、配置等）
4. **Interfaces（接口层）** - 处理外部交互（HTTP、DTO、中间件）

### 事务管理

使用上下文传播事务：

```go
txManager.InTx(ctx, func(ctx context.Context) error {
    // 在此函数内的所有仓储操作都在同一事务中
    user, err := userRepo.GetByID(ctx, userID)
    if err != nil {
        return err
    }
    user.Name = "New Name"
    return userRepo.Update(ctx, user)
})
```

### 添加新功能

参考 `CLAUDE.md` 文件中的详细指南，了解如何添加新实体和端点。

## 配置说明

所有配置通过环境变量或 `.env` 文件管理：

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `ENV` | 运行环境（dev/prod） | `prod` |
| `PORT` | 服务端口 | `8808` |
| `DB_DSN` | PostgreSQL 连接串 | - |
| `LOG_LEVEL` | 日志级别 | `info` |
| `JWT_SECRET` | JWT 密钥 | `dev_secret_change_me` |
| `JWT_EXPIRE_DURATION` | Token 过期时间 | `1h` |

## 测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test -v ./internal/domain/service

# 运行单个测试
go test -v ./internal/domain/service -run TestExampleService
```

## 部署

### 构建

```bash
make build
```

生成的二进制文件位于 `build/` 目录。

### 使用部署脚本

```bash
# 启动应用
./deploy.sh start

# 停止应用
./deploy.sh stop

# 重启应用
./deploy.sh restart

# 健康检查
./deploy.sh health-check
```

## 常见问题

**Q: 如何修改默认端口？**

A: 在 `.env` 文件中设置 `PORT=你的端口号`

**Q: 如何开启调试日志？**

A: 设置 `ENV=dev` 和 `LOG_LEVEL=debug`

**Q: 数据库连接失败？**

A: 检查 `DB_DSN` 配置是否正确，确保 PostgreSQL 服务正在运行

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License
