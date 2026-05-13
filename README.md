# go-order-inventory (轻量级订单库存管理系统项目)

# 1.项目简介

本项目基于 Go + Gin + Gorm + MySQL 实现轻量级订单库存管理系统，支持商品创建、商品查询、商品上架、商品下架、库存初始化、获取商品库存等等功能

项目采用 handler/ service/ dao/ model/ request/	response/ 分层结构，使用环境变量管理数据库配置，并通过统一请求、响应和业务逻辑错误映射提升接口规范性。 

# 2.技术栈

- GO

- Gin

- MySQL

- Gorm

- godotenv

- Yaml配置

- Redis (待接入)

# 3.核心功能

- 商品模块

	- 创建商品

	- 查询商品

	- 商品上架 / 下架	

- 库存模块

	- 库存初始化

	- 查询库存

	- 添加库存

- 查询商品库存流水

- 订单模块

	- 创建订单

	- 查询订单

	- 查询订单详情

	- 支付订单 / 取消订单 / 完成订单

- 统一响应结构

- 环境变量配置

- README + 接口文档 (待完成)

- 事务锁 + 行锁 (待完成)

- Redis 缓存商品详情 (待完成)

# 4. 核心业务流程

# 5. 业务规则说明

# 6.项目结构

```

cmd/	项目启动的主入口

config/	Yaml 配置加载

docs/	项目文档，存放项目手动测试 .http 和初始化 SQL 的文件夹

docs/http 项目接口的测试文件，用于手动测试项目各类接口

docs/sql 项目用测试数据的初始化 SQL 文件

global/	全局资源，如 DB

pkg/	引用外部的资源，如Mysql

router/	全局路由

internal/	资源接口 REST API

internal/handler	HTTP 接口层，负责请求处理、参数绑定、响应数据

internal/service	业务逻辑层，包含商品创建、商品查询、商品上架\下架等相关业务逻辑

internal/dao	数据库操作层，负责数据库操作

internal/model	model 层，用于定义数据结构实体

internal/response	response 层，负责响应数据，返回数据结果

internal/request	request 层，负责处理客户端发送的请求的数据类型

```

# 7.数据表说明

1、product 商品表

```
CREATE TABLE `products` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `name` varchar(100) NOT NULL,
    `description` varchar(500) NOT NULL DEFAULT '',
    `price_fen` bigint NOT NULL,
    `status` tinyint NOT NULL DEFAULT '2',
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_products_status` (`status`)
) ENGINE = InnoDB AUTO_INCREMENT = 17 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci
```

2、product_inventories 商品库存表

```
CREATE TABLE `product_inventories` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `product_id` bigint NOT NULL,
    `stock_quantity` bigint NOT NULL DEFAULT '0',
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_inventory_product_id` (`product_id`),
    CONSTRAINT `chk_product_inventories_stock_quantity` CHECK ((`stock_quantity` >= 0))
) ENGINE = InnoDB AUTO_INCREMENT = 21 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci
```

3、stock_logs 商品库存流水记录表

```
CREATE TABLE `stock_logs` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `product_id` bigint NOT NULL,
    `change_quantity` bigint NOT NULL,
    `before_quantity` bigint NOT NULL,
    `after_quantity` bigint NOT NULL,
    `biz_type` tinyint NOT NULL,
    `biz_id` bigint DEFAULT NULL,
    `remark` varchar(255) NOT NULL DEFAULT '',
    `created_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_stock_logs_product_id` (`product_id`),
    KEY `idx_stock_logs_biz_type` (`biz_type`),
    KEY `idx_stock_logs_biz_id` (`biz_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 84 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci
```

4、order 订单表

```
CREATE TABLE `orders` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `order_no` varchar(255) NOT NULL,
    `total_amount_fen` bigint NOT NULL,
    `status` tinyint NOT NULL DEFAULT '1',
    `paid_at` datetime DEFAULT NULL,
    `completed_at` datetime DEFAULT NULL,
    `cancelled_at` datetime DEFAULT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_orders_order_no` (`order_no`),
    KEY `idx_orders_status` (`status`)
) ENGINE = InnoDB AUTO_INCREMENT = 26 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci
```

5、order_items 订单详情表

```
CREATE TABLE `order_items` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `order_id` bigint NOT NULL,
    `product_id` bigint NOT NULL,
    `product_name` varchar(100) NOT NULL,
    `product_price_fen` bigint NOT NULL,
    `quantity` bigint NOT NULL,
    `subtotal_fen` bigint NOT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_order_items_order_id` (`order_id`),
    KEY `idx_order_items_product_id` (`product_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 22 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci
```
# 8.环境变量

```

DB_USER=root

DB_PASSWORD=your_password

DB_URL=127.0.0.1

DB_PORT=3306

DB_NAME=go-order-inventory

```

# 9.启动方式

```

go mod tidy 

go run cmd/main.go

```

# 10.接口清单

GET /ping

响应

```
{
    "code": 0,
	"msg": "success",
    "data": {
        "message": "success"
    }
}
```

POST /api/v1/products

GET /api/v1/products

GET /api/v1/products/:id

PATCH /api/v1/products/:id/on-sale

PATCH /api/v1/products/:id/off-sale

POST /api/v1/inventory/init

POST /api/v1/inventory/add

GET /api/v1/inventory/products/:product_id

GET /api/v1/stock-logs

POST /api/v1/orders

GET /api/v1/orders/:id

GET /api/v1/orders

PATCH /api/v1/orders/:id/cancel

PATCH /api/v1/orders/:id/pay

PATCH /api/v1/orders/:id/finish


# 11.手动测试说明

本项目使用 VS Code REST Client 插件维护手动接口测试文件。

测试文件位置：

```text
docs/http/*.http
```

# 12.核心自测清单

## 商品模块

- [ ] 正常创建订单

## 库存模块

- [ ] 内容

## 订单模块

- [ ] 内容


# 13.设计与实现要点


## 1. 分层结构设计

项目采用 handler / service / dao / model / request / response 分层结构：

- handler：负责 HTTP 请求处理、参数绑定、错误映射和统一响应
- service：负责业务规则、状态流转、事务控制和跨表操作
- dao：负责数据库 CRUD 和条件更新
- model：负责 GORM 实体映射
- request：负责接口入参结构和参数校验
- response：负责统一响应结构

这种拆分避免了 HTTP 处理、业务规则和数据库操作混在一起，方便后续扩展订单状态流转、库存扣减、库存流水和 Redis 缓存。

## 2. 商品设计

商品创建后默认下架，状态值：

- 1：上架
- 2：下架

创建商品时会校验商品名称和价格，价格使用 price_fen 存储，避免浮点精度问题。

## 3. 库存设计

库存表使用 product_inventories，每个商品只能有一条库存记录，通过 product_id 唯一索引保证。

库存变化不直接只改库存表，而是同时写入 stock_logs，记录：

- before_quantity：变更前库存
- change_quantity：变更数量
- after_quantity：变更后库存
- biz_type：业务类型
- biz_id：关联业务 ID

## 4. 订单事务设计

创建订单时在同一个数据库事务中完成：

1. 创建订单主表
2. 校验商品存在且已上架
3. 校验库存存在且充足
4. 扣减库存
5. 创建订单项
6. 写入库存流水
7. 更新订单总金额

任一环节失败时，整个事务回滚，避免出现订单创建成功但库存未扣减，或库存已扣减但订单失败的问题。

## 5. 订单状态机

订单状态：

- 1：待支付
- 2：已支付
- 3：已完成
- 4：已取消

允许的状态流转：

- 待支付 -> 已支付
- 已支付 -> 已完成
- 待支付 -> 已取消

不允许的状态流转：

- 已支付订单不能取消
- 已完成订单不能取消
- 已取消订单不能支付或完成
- 未支付订单不能完成

# 14. 当前问题与后续演进
# 15. 简历项目描述
