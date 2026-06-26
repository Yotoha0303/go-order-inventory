# go-order-inventory

一个轻量级订单库存管理系统，重点展示 Go 后端业务分层、事务一致性、库存扣减、订单状态机、Redis 缓存和测试能力。

## 1. 项目简介

本项目基于 Go + Gin + GORM + MySQL + Redis 实现，提供商品管理、库存管理、库存流水、订单创建、订单状态流转和商品详情缓存等能力。

项目的目标不是堆功能，而是把常见后端工程能力做扎实：

- 清晰的 handler / service / dao / model 分层
- 统一请求参数校验和响应结构
- 使用事务保证订单创建和库存扣减一致
- 使用库存流水追踪每一次库存变化
- 使用订单状态机限制非法状态流转
- 通过文档和测试清单支撑项目复盘

### 工程化能力

- HTTP Server 设置 ReadTimeout、WriteTimeout、IdleTimeout、ReadHeaderTimeout 和 MaxHeaderBytes
- MySQL 初始化时配置连接池：MaxOpenConns、MaxIdleConns、ConnMaxLifetime、ConnMaxIdleTime
- 启动时使用 PingContext 检查 MySQL 连通性
- 请求层使用 Request ID、Access Log 和 Recovery；HTTP Server 外层强制请求超时，向下游传递 deadline，超时返回 `503 / 5002` 并隔离后续响应写入
- Redis 不可用时商品详情缓存自动降级，不影响主流程
- CI 覆盖 go test、go test -race、go vet、golangci-lint、goose validate 和 go build
- Dockerfile 使用多阶段构建和非 root 用户运行应用
- Docker Compose 编排应用、MySQL 和 Redis，并通过健康检查控制依赖启动顺序
- Goose 管理数据库版本，Makefile 统一封装开发、测试、Docker 和迁移命令

## 2. 技术栈

- Go 1.25.5
- Gin + GORM
- MySQL 8.4 + Redis 7.2
- Goose 数据库迁移
- YAML + godotenv 配置
- Docker + Docker Compose
- golangci-lint + GitHub Actions

## 3. 当前实现进度

已完成：

- 商品创建、查询、上下架
- 库存初始化、增加、查询
- 库存流水记录
- 创建订单时扣减库存
- 库存不足时事务回滚
- 订单支付、完成、取消
- 取消待支付订单时回滚库存
- 订单状态机限制非法流转
- 商品详情 Redis cache-aside 缓存
- 商品上架 / 下架时删除缓存
- Redis 不可用时不影响主流程

未实现，作为后续演进：

- 创建订单 request_id / idempotency_key 幂等控制
- 除健康检查外的 handler 业务接口测试、DAO 层测试与更完整的异常分支覆盖
- 并发下单防超卖测试，以及多商品下单中途失败后的事务完整回滚测试
- 自动迁移或独立 migration job，进一步简化首次启动流程

## 4. 核心功能

### 商品模块

- 创建商品
- 查询商品列表
- 查询商品详情
- 商品上架
- 商品下架

### 库存模块

- 初始化商品库存
- 增加商品库存
- 查询商品库存
- 记录库存变更流水

### 订单模块

- 创建订单
- 查询订单列表
- 查询订单详情
- 支付订单
- 完成订单
- 取消订单
- 取消订单时回滚库存

### Redis缓存

- 商品详情缓存 cache-aside 
- 商品上架 / 下架时删除商品详情缓存

## 5. 项目结构

```text
.github/workflows/ci.yml  持续集成配置
cmd/                      项目启动入口
config/                   YAML 加载、环境变量覆盖和配置校验
docs/                     设计文档、REST Client 请求和验证证据
internal/apperror/        业务错误定义与错误码映射
internal/app/             依赖装配、HTTP Server 和优雅退出
internal/bizcache/        Redis 业务缓存
internal/dao/             数据库访问层
internal/handler/         HTTP 接口层
internal/middleware/      请求 ID、日志、超时和恢复中间件
internal/model/           GORM 数据模型
internal/request/         请求参数和校验规则
internal/response/        统一响应结构
internal/service/         业务规则、状态机和事务
migrations/               Goose SQL 迁移
pkg/database/             MySQL 初始化与连接池
pkg/redis/                Redis 客户端初始化
router/                   路由注册
compose.yml               应用、MySQL、Redis 编排
Dockerfile                应用镜像多阶段构建
Makefile                  开发、测试、Docker 和迁移命令入口
```

## 6. 分层说明

项目采用简单的企业后端分层方式：

- handler：负责 HTTP 请求处理、参数绑定、错误映射和统一响应
- service：负责业务规则、状态流转、事务控制和跨表操作
- dao：负责数据库 CRUD、条件查询和条件更新
- model：负责数据库表结构映射
- request：负责接口入参结构和校验规则
- response：负责接口响应结构
- bizcache：负责业务缓存读写、缓存 key 规则和缓存失效
- apperror：负责业务错误定义、错误码和错误信息封装

核心原则：handler 不写业务规则，service 不直接拼 HTTP 响应，dao 不处理业务状态。

## 7. 数据表设计

当前核心表：

- products：商品表
- product_inventories：商品库存表
- stock_logs：库存流水表
- orders：订单主表
- order_items：订单明细表

关键设计点：

- 商品价格使用 price_fen，单位为分，避免浮点精度问题
- 商品创建后默认下架，避免未准备库存的商品直接下单
- product_inventories 通过 product_id 唯一索引保证一个商品只有一条库存记录
- stock_logs 记录 before_quantity、change_quantity、after_quantity，便于追踪库存变化
- orders 使用状态机控制待支付、已支付、已完成、已取消
- order_items 保存下单时的商品名称和价格快照

详细表结构见：[docs/table_design.md](docs/table_design.md)

## 8. 核心业务规则

### 商品规则

- 商品名称不能为空
- 商品价格 price_fen 必须大于 0
- 商品创建后默认下架，status = 2
- 商品上架后 status = 1
- 商品下架后 status = 2
- 查询商品时，默认查询下架的商品, status = 2

### 库存规则

- 初始化库存前商品必须存在
- 一个商品只能初始化一次库存
- 增加库存前库存记录必须存在
- 库存变更必须写入 stock_logs
- 库存流水 biz_type：1 初始化库存，2 手动入库，3 订单扣减，4 取消订单回滚

### 订单规则

- 创建订单时 items 不能为空
- 下单商品必须存在且已上架
- 商品库存必须存在且充足
- 创建订单、扣减库存、创建订单项、写库存流水必须在同一个事务内完成
- 取消待支付订单时需要回滚库存

### Redis 缓存规则

- 查询商品详情时，设置商品缓存
- 上架/下架 商品时，删除商品缓存

详细规则见：[docs/business_rules.md](docs/business_rules.md)

## 9. 订单创建事务设计

当前项目在订单创建与取消场景中，使用数据库事务保证订单、库存和库存流水一致性，并通过库存行级锁控制并发扣减。

### 订单创建

订单创建流程

- 通过当前时间戳生成 orderNO，并创建订单
- 遍历订单商品项
- 对每个商品库存记录使用行级锁（`FOR UPDATE`）读取并计算调整前后库存
- 减去需要扣减的商品库存
- 创建订单明细 order_items
- 创建并记录商品库存调整流水
- 任一步骤失败时，事务整体回滚，避免出现部分写入

### 订单取消

订单取消流程

- 查询需要取消的订单
- 判断订单状态，仅允许取消待支付订单；已取消订单按幂等直接返回
- 查询订单下已订购的商品
- 遍历订单商品项，并回滚商品库存
- 创建并记录商品库存调整流水
- 任一步骤失败时，事务整体回滚，避免库存回滚不完整


## 10. 订单状态机

订单状态：

- 1：待支付
- 2：已支付
- 3：已完成
- 4：已取消

允许的状态流转：

- 待支付 -> 已支付
- 已支付 -> 已完成
- 待支付 -> 已取消

禁止的状态流转：

- 已支付订单不能取消
- 已完成订单不能取消
- 已取消订单不能支付或完成
- 未支付订单不能完成

## 11. Redis 商品详情缓存设计

当前项目对商品详情接口增加了 cache-aside 缓存。

### 缓存 key

```text
product:detail:{product_id}
```

### 查询流程

- 查询商品详情时，优先读取 Redis

- 如果 Redis 命中，直接返回缓存数据

- 如果 Redis 未命中，查询 MySQL

- MySQL 查询成功后，将商品详情写入 Redis

- Redis 不可用时，不影响 MySQL 主流程

### 缓存删除

- 商品状态变化时删除缓存：

- 商品上架：删除商品详情缓存

- 商品下架：删除商品详情缓存

### 当前边界

- 当前仅对商品详情接口实现 cache-aside 缓存
- 当前通过“状态变更时主动删缓存”保证基础一致性，不包含延迟双删等增强策略
- 当前未引入缓存击穿保护（如互斥锁、逻辑过期），后续可按流量特征演进
- Redis 不可用时直接降级到 MySQL 主流程，优先保证业务可用性



## 12. 幂等设计说明

本项目只对状态设置类接口做轻量级幂等处理：

- 商品已上架时，再次调用上架接口直接成功
- 商品已下架时，再次调用下架接口直接成功
- 已取消订单再次取消时直接成功，但不会重复回滚库存

以下接口当前不做幂等：

- 创建订单：重复调用会创建多笔订单
- 增加库存：重复调用会多次增加库存
- 初始化库存：重复调用会返回库存已初始化错误
- 支付订单：重复支付返回状态冲突
- 完成订单：重复完成返回状态冲突

## 13. 接口清单

接口说明详见：[docs/api_list.md](docs/api_list.md)

## 14. 配置与环境变量

应用启动时先加载 `.env`，再读取 [config.yml](config.yml)。环境变量会覆盖 YAML 中适合按环境变化的连接配置。

### 14.1 基础环境变量

```env
MYSQL_PASSWORD=your-password
REDIS_PASSWORD=
```

- `MYSQL_PASSWORD`：必填，应用、Docker Compose 和 Goose 共用。
- `REDIS_PASSWORD`：可选；当前 Compose 中的 Redis 未启用密码认证，保持为空即可。
- 不要提交真实的 `.env`，可从 [.env.example](.env.example) 复制后修改。

### 14.2 YAML 配置覆盖

| 环境变量 | 覆盖的配置 |
| --- | --- |
| `APP_PORT` | `server.port` |
| `DB_HOST` | `mysql.host` |
| `DB_PORT` | `mysql.port` |
| `DB_USER` | `mysql.user` |
| `DB_NAME` | `mysql.database` |
| `REDIS_ADDR` | `redis.addr` |
| `REDIS_DB` | `redis.db` |

本地运行默认连接 `127.0.0.1:3306` 和 `127.0.0.1:6379`。Compose 会为应用容器设置 `mysql:3306` 和 `redis:6379`，无需修改 `config.yml`。

## 15. 启动方式

### 15.1 前置依赖

- Go 1.25.5
- GNU Make
- Docker 与 Docker Compose
- Goose v3.27.1

```bash
go mod download
go install github.com/pressly/goose/v3/cmd/goose@v3.27.1
```

### 15.2 本地运行应用

PowerShell 示例：

```powershell
Copy-Item .env.example .env
$env:MYSQL_PASSWORD = "your-password"

make infra-up
make migrate-up
make run
```

首次启动必须执行 `make migrate-up` 建表。之后可用 `make dev` 启动 MySQL、Redis 并运行应用。

### 15.3 Docker 运行完整服务

```powershell
$env:MYSQL_PASSWORD = "your-password"

make infra-up
make migrate-up
make docker-up
```

迁移由宿主机上的 Goose 执行，默认连接 Compose 暴露的 `127.0.0.1:3306`。当前 `docker-up` 不会自动执行数据库迁移。

常用 Docker 命令：

| 命令 | 作用 |
| --- | --- |
| `make compose-config` | 校验 Compose 配置 |
| `make infra-up` | 仅启动 MySQL 和 Redis，并等待健康 |
| `make infra-down` | 停止 Compose 项目 |
| `make docker-build` | 构建应用镜像 |
| `make docker-up` | 构建并启动应用、MySQL、Redis |
| `make docker-down` | 停止并移除容器，保留数据卷 |
| `make docker-ps` | 查看服务状态 |
| `make docker-logs` | 持续查看全部服务日志 |

常用迁移命令：

| 命令 | 作用 |
| --- | --- |
| `make migrate-validate` | 静态校验迁移文件 |
| `make migrate-status` | 查看数据库迁移状态 |
| `make migrate-up` | 执行全部待处理迁移 |
| `make migrate-up-one` | 只执行下一条迁移 |
| `make migrate-up-to VERSION=5` | 迁移到指定版本 |
| `make migrate-down` | 回滚最近一条迁移 |
| `make migrate-down-to VERSION=3` | 回滚到指定版本 |
| `make migrate-redo` | 重做最近一条迁移 |
| `make migrate-create NAME=add_sku` | 创建顺序编号的 SQL 迁移 |

默认访问地址为 `http://localhost:8082`，健康检查为：

```bash
curl http://localhost:8082/ping
```

## 16. 测试方式

service 测试会清理所连接数据库中的业务表。必须使用独立测试库，禁止将 `DB_NAME` 指向含有开发数据或生产数据的数据库。

### 16.1 自动化测试覆盖现状

| 测试项 | 当前状态 | 检查结论 |
| --- | --- | --- |
| 核心业务测试 | 已有 | `internal/service/*_test.go` 已覆盖商品、库存、订单创建和订单状态机；另有健康检查 handler 测试与 Redis 缓存测试 |
| 并发下单测试 | 待补充 | 当前没有多个 goroutine 同时购买同一商品并校验成功订单数、最终库存和防超卖结果的测试 |
| 创建订单事务回滚测试 | 部分覆盖 | `TestCreateOrder_InsufficientStock` 已校验库存不足错误，但没有构造“前一商品已扣减、后一商品失败”的场景，也没有断言订单、订单项、库存和库存流水全部回滚 |

取消待支付订单后的库存恢复已有 `TestCancelOrder_Success` 和重复取消幂等测试覆盖；这是取消业务的库存补偿，不等同于创建订单失败时的数据库事务回滚测试。

建议补充的用例名称：

- `TestCreateOrder_ConcurrentStock_NoOversell`
- `TestCreateOrder_RollbackWhenLaterItemInsufficient`

假设已创建 `go_order_inventory_test`：

```powershell
$env:MYSQL_PASSWORD = "your-password"
$env:DB_NAME = "go_order_inventory_test"

make migrate-up
make test
```

常用测试和质量命令：

| 命令 | 作用 |
| --- | --- |
| `make test` | 运行全部 Go 测试 |
| `make test-service` | 运行 service 测试 |
| `make test-redis` | 运行 Redis 集成测试 |
| `make test-all` | 运行普通测试和 Redis 集成测试 |
| `make test-race` | 使用 race detector 运行测试 |
| `make coverage` | 生成 `coverage.out` |
| `make coverage-html` | 生成 `coverage.html` |
| `make check` | 执行格式化、模块校验、vet 和测试 |

Redis 集成测试前需保证 Redis 已启动，可先执行 `make infra-up`。

手动接口测试文件位于 [docs/http](docs/http)，完整业务链路见 [docs/http/demo_flow.http](docs/http/demo_flow.http)。测试计划见 [docs/test_plan.md](docs/test_plan.md)。

## 17. 项目文档

- [docs/api_list.md](docs/api_list.md)：接口清单
- [docs/business_rules.md](docs/business_rules.md)：业务规则
- [docs/table_design.md](docs/table_design.md)：数据表设计
- [docs/test_plan.md](docs/test_plan.md)：测试计划
- [docs/test_result.md](docs/test_result.md)：测试结果记录
- [docs/project_evolution.md](docs/project_evolution.md)：后续演进
- [docs/evidence](docs/evidence)：项目运行、测试与关键业务截图证据

### 17.1 项目证据链（docs/evidence/）

本目录用于保存项目运行、测试和关键业务链路截图，便于项目展示和面试讲解。

- 创建订单成功：`create_order_success_2026-05-23_17-26-48.png`
- 创建订单库存不足并回滚：`create_order_insufficient_inventory_rollback_2026-05-25_00-02-52.png`
- 取消订单后库存回滚：`order_cancel_inventory_rollback_2026-05-25_00-10-18.png`
- Redis 商品详情缓存命中：`redis_get_product_cache_success_2026-05-25_00-16-01.png`
- 商品上架/下架后缓存删除：`redis_on_or_off_sale_product_cache_delete_success_2026-05-25_00-16-01.png`
- Redis 集成测试执行成功：`redis_test_execute_success_2026-05-23_17-23-36.png`
- 自动化测试运行结果（分段截图）：
  `test_run_success_part_1_2026-05-23_17-19-26.png`
  `test_run_success_part_2_2026-05-23_17-19-26.png`
  `test_run_success_part_3_2026-05-23_17-19-26.png`

## 18. 当前可复盘亮点

- 使用 handler / service / dao / model 分层组织代码，避免业务逻辑散落在接口层
- 创建订单使用事务保证 orders、order_items、product_inventories、stock_logs 多表一致性
- 库存扣减使用库存行锁 + 条件更新，避免库存不足时继续扣减；并发防超卖效果仍待专门的并发测试验证
- order_items 保存商品名称和价格快照，避免商品后续修改影响历史订单
- stock_logs 记录库存变更前后数量、业务类型和业务 ID，便于排查库存异常
- 订单状态机限制待支付、已支付、已完成、已取消之间的非法流转
- 取消待支付订单时回滚库存，并记录 biz_type=4 的库存流水
- 商品详情使用 Redis cache-aside 缓存，商品上下架时删除缓存
- Redis 异常时降级走 MySQL，不影响主业务流程
- 使用 AppError 统一业务错误、HTTP 状态码和业务 code，减少 handler 层重复错误判断
- 配置支持 YAML 默认值、环境变量覆盖和启动参数校验
- HTTP Server 配置超时、请求 ID、访问日志、panic 恢复和优雅退出
- Docker 使用多阶段构建、非 root 用户、健康检查和 Compose 依赖编排
- Goose 管理数据库版本，CI 自动执行 lint、test、race、vet、build 和迁移校验

## 19. 后续演进方向

- 增加 handler 业务接口测试和 DAO 测试
- 增加同一商品并发下单测试，验证不超卖、成功订单数和最终库存一致
- 增加多商品下单中途失败测试，验证订单、订单项、库存和库存流水整体回滚
- 在 Compose 或部署流水线中加入独立 migration job
- 增加指标、链路追踪和结构化日志字段规范
- 优化错误码文档和接口返回示例
- 订单中使用雪花 ID 代替时间戳生成 orderNO
- 创建订单时可引入 client_order_no / idempotency_key，避免重复下单
