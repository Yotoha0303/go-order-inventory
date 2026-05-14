# go-order-inventory

轻量级订单库存管理系统，一个面向 Go 后端求职实战的项目。项目围绕商品、库存、库存流水和订单状态流转展开，重点练习 Gin 接口开发、GORM 数据建模、MySQL 事务、库存扣减一致性和分层代码组织。

## 1. 项目简介

本项目基于 Go + Gin + GORM + MySQL 实现，提供商品管理、库存管理、库存流水查询、订单创建、订单支付、订单完成和订单取消等能力。

项目目标不是堆功能，而是把常见后端工程能力做扎实：

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
- Redis（已预留接入目录）
- godotenv
- YAML 配置

## 3. 核心功能

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

## 4. 项目结构

```text
cmd/                  项目启动入口
config/               配置加载
docs/                 项目文档、接口测试文件、SQL 脚本
docs/http/            REST Client 手动接口测试文件
docs/sql/             初始化和测试 SQL
global/               全局资源，如 DB
internal/dao/         数据库访问层
internal/handler/     HTTP 接口层
internal/model/       GORM 数据模型
internal/request/     请求参数结构
internal/response/    响应结构
internal/service/     业务逻辑层
pkg/database/         MySQL 初始化
pkg/redis/            Redis 初始化预留
router/               路由注册
```

## 5. 分层说明

项目采用简单的企业后端分层方式：

- handler：负责 HTTP 请求处理、参数绑定、错误映射和统一响应
- service：负责业务规则、状态流转、事务控制和跨表操作
- dao：负责数据库 CRUD、条件查询和条件更新
- model：负责数据库表结构映射
- request：负责接口入参结构和校验规则
- response：负责接口响应结构

核心原则：handler 不写业务规则，service 不直接拼 HTTP 响应，dao 不处理业务状态。

## 6. 数据表设计

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

## 7. 核心业务规则

### 商品规则

- 商品名称不能为空
- 商品价格 price_fen 必须大于 0
- 商品创建后默认下架，status = 2
- 商品上架后 status = 1
- 商品下架后 status = 2

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

详细规则见：[docs/business_rules.md](docs/business_rules.md)

## 8. 订单状态机

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

## 9. 接口清单

### 健康检查

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | /ping | 健康检查 |

### 商品接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | /api/v1/products | 创建商品 |
| GET | /api/v1/products | 查询商品列表 |
| GET | /api/v1/products/:id | 查询商品详情 |
| PATCH | /api/v1/products/:id/on-sale | 商品上架 |
| PATCH | /api/v1/products/:id/off-sale | 商品下架 |

### 库存接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | /api/v1/inventory/init | 初始化库存 |
| POST | /api/v1/inventory/add | 增加库存 |
| GET | /api/v1/inventory/products/:product_id | 查询商品库存 |
| GET | /api/v1/stock-logs | 查询库存流水 |

### 订单接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | /api/v1/orders | 创建订单 |
| GET | /api/v1/orders | 查询订单列表 |
| GET | /api/v1/orders/:id | 查询订单详情 |
| PATCH | /api/v1/orders/:id/pay | 支付订单 |
| PATCH | /api/v1/orders/:id/finish | 完成订单 |
| PATCH | /api/v1/orders/:id/cancel | 取消订单 |

完整接口说明见：[docs/api_list.md](docs/api_list.md)

## 10. 环境变量

项目通过 `.env` 读取数据库配置：

```env
DB_USER=root
DB_PASSWORD=your_password
DB_URL=127.0.0.1
DB_PORT=3306
DB_NAME=go-order-inventory
```

服务端口配置在 `config.yml`：

```yaml
server:
  port: 8082
```

## 11. 启动方式

安装依赖：

```bash
go mod tidy
```

启动服务：

```bash
go run cmd/main.go
```

默认访问地址：

```text
http://localhost:8082
```

健康检查：

```bash
curl http://localhost:8082/ping
```

## 12. 测试方式

运行 Go 测试：

```bash
go test ./...
```

手动接口测试：

```text
docs/http/products.http
docs/http/inventory.http
docs/http/stock_logs.http
docs/http/orders.http
```

测试计划见：[docs/test_plan.md](docs/test_plan.md)

## 13. 项目文档

- [docs/api_list.md](docs/api_list.md)：接口清单
- [docs/business_rules.md](docs/business_rules.md)：业务规则
- [docs/table_design.md](docs/table_design.md)：数据表设计
- [docs/test_plan.md](docs/test_plan.md)：测试计划
- [docs/project_review.md](docs/project_review.md)：项目复盘
- [docs/project_evolution.md](docs/project_evolution.md)：后续演进

## 14. 当前可复盘亮点

- 订单创建使用事务包住订单、库存、订单项和库存流水
- 库存扣减使用条件更新，避免库存不足时继续扣减
- 库存变化有 stock_logs，可追踪业务来源
- 订单状态通过 service 层统一控制，避免接口层散落状态判断
- 项目文档、接口测试文件和测试清单逐步补齐，便于面试讲解

## 15. 后续演进方向

- 补充更多 service 层单元测试
- 增加 handler 层接口测试
- 接入 Redis 缓存商品详情
- 增加订单幂等控制，避免重复下单或重复取消
- 优化错误码文档和接口返回示例
- 增加 Docker Compose，降低本地启动成本

