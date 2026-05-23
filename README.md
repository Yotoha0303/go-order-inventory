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

## 2. 技术栈

- Go
- Gin
- GORM
- MySQL
- Redis
- godotenv
- YAML 配置

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
- 更完整的 service 层自动化测试
- 本地启动脚本优化（依赖就绪检测与一键初始化）

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
cmd/                  项目启动入口
config/               配置加载
docs/                 项目文档、接口测试文件、SQL 脚本
docs/http/            REST Client 手动接口测试文件
docs/sql/             初始化和测试 SQL
global/               全局资源，如 DB
internal/apperror/    业务错误定义与错误码映射
internal/bizcache/    数据缓存
internal/dao/         数据库访问层
internal/handler/     HTTP 接口层
internal/model/       GORM 数据模型
internal/request/     请求参数结构
internal/response/    响应结构
internal/service/     业务逻辑层
pkg/database/         MySQL 初始化
pkg/redis/            Redis 初始化
router/               路由注册
docker-compose.yml.example 本地依赖服务编排示例
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

## 14. 环境变量

项目通过 `.env` 读取数据库配置：

```env
# runtime
MYSQL_PASSWORD=your-password
REDIS_PASSWORD=
# test
MYSQL_TEST_PASSWORD=your-password
MYSQL_TEST_DATABASE=go_order_inventory_test
```

服务端口配置在 `config.yml`：

```yaml
server:
  port: 8082

mysql:
  user: root
  host: 127.0.0.1
  port: "3306"
  database: go_order_inventory

redis:
  addr: 127.0.0.1:6379
  db: 0

```

项目可通过 Docker 进行部署， `docker-compose.yml` ：

```
services:
  mysql:
    image: mysql:${MYSQL_IMAGE_VERSION:-8.4.8}
    container_name: go-order-inventory-mysql
    restart: unless-stopped
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_PASSWORD:-your-password}
      MYSQL_DATABASE: ${MYSQL_DATABASE:-go_order_inventory_demo}
      TZ: Asia/Shanghai
    command:
      - --character-set-server=utf8mb4
      - --collation-server=utf8mb4_general_ci
    volumes:
      - mysql_data:/var/lib/mysql
      - ./docs/sql/sql_table_data_init.sql:/docker-entrypoint-initdb.d/01_init.sql:ro
    healthcheck:
      test: ["CMD-SHELL", "mysqladmin ping -h 127.0.0.1 -uroot -p$$MYSQL_ROOT_PASSWORD --silent"]
      interval: 10s
      timeout: 5s
      retries: 10

  redis:
    image: redis:${REDIS_IMAGE_VERSION:-7.2-alpine}
    container_name: go-order-inventory-redis
    restart: unless-stopped
    ports:
      - "6379:6379"
    command: ["redis-server", "--appendonly", "yes"]
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 10

volumes:
  mysql_data:
  redis_data:

```

## 15. 启动方式

安装依赖：

```bash
go mod tidy
```

启动服务：

```bash
go run cmd/main.go
```

Docker 启动：

本项目依赖 MySQL 和 Redis 的依赖，需要安装依赖后才能启动项目

```bash
docker compose up -d
```

默认访问地址：

```text
http://localhost:8082
```

健康检查：

```bash
curl http://localhost:8082/ping
```

## 16. 测试方式

运行 Go 测试：

```bash
go test -v ./...
```

手动接口测试：

```text
docs/http/products.http
docs/http/inventory.http
docs/http/stock_logs.http
docs/http/orders.http
docs/http/redis.http
```

Redis 集成测试：

```bash
RUN_REDIS_TEST=1 go test -v ./internal/bizcache
```

测试计划见：[docs/test_plan.md](docs/test_plan.md)

## 17. 项目文档

- [docs/api_list.md](docs/api_list.md)：接口清单
- [docs/business_rules.md](docs/business_rules.md)：业务规则
- [docs/table_design.md](docs/table_design.md)：数据表设计
- [docs/test_plan.md](docs/test_plan.md)：测试计划
- [docs/test_result.md](docs/test_result.md)：测试结果记录
- [docs/project_review.md](docs/project_review.md)：项目复盘
- [docs/project_evolution.md](docs/project_evolution.md)：后续演进
- [docs/interview_guide.md](docs/interview_guide.md)：面试讲解提纲

## 18. 当前可复盘亮点

- 使用 handler / service / dao / model 分层组织代码，避免业务逻辑散落在接口层
- 创建订单使用事务保证 orders、order_items、product_inventories、stock_logs 多表一致性
- 库存扣减使用库存行锁 + 条件更新，避免库存不足时继续扣减
- order_items 保存商品名称和价格快照，避免商品后续修改影响历史订单
- stock_logs 记录库存变更前后数量、业务类型和业务 ID，便于排查库存异常
- 订单状态机限制待支付、已支付、已完成、已取消之间的非法流转
- 取消待支付订单时回滚库存，并记录 biz_type=4 的库存流水
- 商品详情使用 Redis cache-aside 缓存，商品上下架时删除缓存
- Redis 异常时降级走 MySQL，不影响主业务流程
- 使用 apperror 统一业务错误、HTTP 状态码和业务 code

## 19. 后续演进方向

- service 层增加更多的测试内容
- 优化 Docker Compose 与本地启动脚本，统一依赖启动与初始化流程
- 优化错误码文档和接口返回示例
- 订单中使用雪花 ID 代替时间戳生成 orderNO
- 创建订单时可引入 client_order_no / idempotency_key，避免重复下单
