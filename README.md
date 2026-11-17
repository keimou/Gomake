# gmkdb

基于 Gin + GORM 的 Go 后端服务，参考 Gin-Vue-Admin 分层架构，支持 JWT 鉴权、Casbin RBAC、多数据库（SQLite/MySQL）、Zap 日志、Viper 配置热重载及 Swagger 文档。

## 技术栈

| 组件 | 版本 | 用途 |
|------|------|------|
| Go | 1.27 | 语言 |
| Gin | v1.9.1 | Web 框架 |
| GORM | v1.30.0 | ORM |
| SQLite / MySQL | - | 数据库 |
| Zap | v1.17.0 | 日志 |
| Viper | v1.12.0 | 配置管理 |
| Redis | v8.11.5 | 缓存 / 单点登录 |
| JWT | v4.4.2 | Token 鉴权 |
| Casbin | - | RBAC 权限控制 |
| Swagger | v1.8.12 | API 文档 |

## 目录结构

```
gmkdb/
├── main.go              # 程序入口，调用 core.RunServer()
├── config.yaml          # YAML 配置文件（jwt/zap/redis/sysdb/db-list/cors 等）
├── go.mod / go.sum      # 依赖管理
├── ninex.sqlite         # SQLite 默认数据库文件
│
├── api/                 # 控制层（Handler）
│   └── v1/
│       ├── system/      # 系统模块接口（用户、登录、验证码等）
│       └── mall/        # 商城模块接口
│       └── enter.go     # API 分组聚合
│
├── config/              # 配置结构体（与 config.yaml 映射）
│   ├── config.go        # Server 顶层配置
│   ├── system.go        # 系统配置
│   ├── jwt.go / zap.go / redis.go / db_mysql.go / db_list.go / cors.go / captcha.go / email.go
│
├── core/                # 核心组件初始化
│   ├── alice.go         # init() 入口：Viper() → Zap() → initalice.Preload()
│   ├── server.go        # RunServer()：启动 HTTP/TLS 服务，优雅关闭
│   ├── server_win.go    # Windows/通用服务实现（initServer, initServerTLS）
│   ├── server_other.go  # 其他平台服务实现（预留）
│   ├── viper.go         # Viper 配置加载与热重载
│   ├── zap.go           # Zap 日志初始化
│   └── internal/        # 内部常量与日志核心
│       ├── constant.go  # 配置文件路径常量
│       ├── zap.go       # Zap Core 组装
│       └── file_rotatelogs.go
│
├── global/              # 全局单例对象
│   ├── aorta.go         # GIA_CONFIG, GIA_DB, GIA_LOG, GIA_DBList, GIA_REDIS, JWT_Concurrency_Control 等
│   ├── gormdb.go        # 数据库列表获取、连接关闭、DB 配置
│   ├── model.go         # GIA_MODEL / TS_MODEL 基础模型
│   └── system.go        # 系统资源清理 Shutdown
│
├── initalice/           # 初始化层
│   ├── aorta.go         # Preload()：激活数据库与通道
│   ├── router.go        # Routers()：Gin 路由注册，中间件挂载，Swagger
│   ├── gorm.go          # GORM 实例初始化
│   ├── db_sqlite.go / db_mysql.go / db_list.go / db_artery.go / redis.go
│   └── internal/        # gorm/logger 内部实现
│
├── routers/             # 路由层
│   ├── enter.go         # RouterGroupApp 聚合（System / Mall 路由组）
│   ├── system/
│   │   ├── enter.go     # 系统路由组定义
│   │   └── sys_base.go  # 基础路由（登录 / 验证码，公开访问）
│   └── mall/
│       └── enter.go     # 商城路由组定义
│
├── middlewares/         # Gin 中间件
│   ├── jwt.go           # JWT 鉴权（Header: sbl-token，黑名单，续期）
│   ├── cors.go          # CORS 跨域（白名单规则）
│   ├── casbin_rbac.go   # RBAC 权限拦截
│   └── ginlog.go        # Gin 请求日志中间件
│
├── models/              # 数据模型
│   ├── system/          # 系统模块实体
│   │   ├── sys_user.go              # 用户实体
│   │   ├── sys_jwt_blacklist.go     # JWT 黑名单
│   │   ├── sys_base_menu.go         # 菜单
│   │   ├── sys_auhority.go          # 角色/权限
│   │   ├── sys_authority_menu.go    # 角色菜单关联
│   │   ├── sys_menu_btn.go          # 菜单按钮
│   │   ├── sys_operation_record.go  # 操作记录
│   │   ├── request/                 # 入参 DTO（Login / Register 等）
│   │   └── response/                # 出参 DTO（LoginResponse 等）
│   └── common/          # 通用模型
│       ├── request/  & response/    # 统一响应 Response {code, data, msg, ts, count}
│
├── services/            # 业务逻辑层
│   ├── enter.go         # ServiceGroupApp 聚合（SystemSvcGroup 等）
│   ├── system/
│   │   ├── sys_user.go      # 用户业务（Register / Login）
│   │   ├── jwt.black_list.go # JWT 黑名单服务
│   │   └── enter.go
│   └── mall/
│       └── enter.go
│
├── utils/               # 工具函数
│   ├── jwt.go           # JWT 生成/解析
│   ├── hash.go          # Bcrypt 密码加密
│   ├── net.go           # 公网/本地 IP 获取
│   ├── validator.go     # 参数校验
│   ├── verify.go        # 验证码校验
│   ├── md.go / struct.go / value.go
│
├── docs/                # Swagger 生成的 API 文档
│   └── docs.go
│
├── log/                 # 运行日志（按日期分目录：debug/info/warn/error）
│   └── 2026-10-28/
│
└── config.yaml          # 运行时配置
```

## 架构分层

```
客户端请求
    ↓
Gin 中间件链：CorsByRules → GinLogger → GinRecovery → JWTAuth
    ↓
路由分发：routers/system（PublicGroup 公开 / /binance 私有组）
    ↓
控制层：api/v1/*（参数绑定 → 调用 Service → 返回统一响应）
    ↓
业务层：services/*（核心业务逻辑）
    ↓
数据层：GORM 操作 global.GIA_DB / global.GIA_DBList
    ↓
数据库：SQLite / MySQL
```

### 各层职责

| 层 | 职责 | 示例 |
|---|---|---|
| `api/` | HTTP 请求入参校验、调用 Service、返回 JSON | `api/v1/system/sys_user.go:Login` |
| `services/` | 业务逻辑、数据库增删改查 | `services/system/sys_user.go:UserService` |
| `models/system/` | GORM 实体（对应数据表） | `models/system/sys_user.go:SysUser` |
| `models/system/request/` | 前端入参 DTO | `request/sys_user.go:Login` |
| `models/system/response/` | 前端出参 DTO | `response/sys_user.go:LoginResponse` |
| `routers/` | Gin 路由组注册 | `routers/system/sys_base.go:InitBaseRouter` |
| `middlewares/` | 跨域、鉴权、RBAC、日志 | `middlewares/jwt.go:JWTAuth` |
| `core/` | 框架核心初始化（Viper/Zap/Server） | `core/alice.go:init()` |
| `initalice/` | 数据库、路由等模块初始化 | `initalice/router.go:Routers()` |
| `global/` | 全局单例容器 | `global/aorta.go` |
| `config/` | 配置结构体定义 | `config/config.go:Server` |

## 启动流程

1. `main()` 调用 `core.RunServer()`
2. 包 `core` 的 `init()`（`core/alice.go`）优先自动执行：
   - `setMode()`：根据命令行参数设置 Gin 运行模式（debug/release/test）
   - `initialize()`：
     - `Viper()`：读取 `config.yaml`，支持热重载
     - `Zap()`：初始化 Zap 日志，挂载到 `global.GIA_LOG`
     - `initalice.Preload()`：预加载数据库连接
3. `RunServer()`：
   - `initalice.Routers()`：构建 Gin Engine，挂载中间件与路由
   - `initServerPriority()`：根据配置选择 HTTP 或 HTTPS 服务
   - 监听端口，等待信号量优雅关闭

## 配置说明（config.yaml）

主要配置项：

- `jwt`：签名密钥、过期时间、缓冲时间、签发人
- `zap`：日志级别、输出目录、格式、是否替换 Gin 默认日志
- `system`：运行环境、监听端口、数据库类型（sqlite/mysql）、是否使用 Redis/TLS/单点登录
- `sysdb`：主数据库连接信息（路径、端口、账号密码等）
- `db-list`：多实例数据库列表，支持 `disable`、`type`、`alias-name` 等字段
- `redis`：Redis 连接信息（用于 JWT 黑名单与单点登录）
- `cors`：跨域白名单规则（allow-all / whitelist / strict-whitelist）
- `captcha`：图形验证码长度与尺寸

## 本地运行

```bash
# 进入项目根目录
cd /home/codes/gero

# 下载依赖
go mod tidy

# 运行（默认 debug 模式，读取 config.yaml）
go run main.go

# 指定运行模式
go run main.go release
go run main.go test

# 指定配置文件
go run main.go -c config.yaml
```

服务启动后：
- 健康检查：`GET http://localhost:8080/health`
- Swagger 文档：`http://localhost:8080/swagger/index.html`
- 登录接口：`POST http://localhost:8080/base/login`

## 注意事项

1. `main.go` 中 `core` 初始化代码已被注释，实际初始化由 `core/alice.go` 的 `init()` 自动完成，无需手动取消注释。
2. 当前 `routers/system/sys_base.go` 中登录/验证码路径已注册，其余业务路由在 `initalice/router.go` 中多半被注释，按需取消注释并注册。
3. `config.yaml` 中 `db-type` 默认为 `sqlite`，`sysdb` 和 `db-list` 配置的是 MySQL 地址，如使用 SQLite 请确保 `ninex.sqlite` 可用，或将 `db-type` 改为 `mysql` 并提供正确的 MySQL 连接。
4. JWT 默认从请求头 `sbl-token` 读取 Token（可在 `config.yaml` 的 `jwt.header-key` 中修改）。
5. 日志按日期自动分目录存储在 `log/` 下。