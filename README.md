go-order-inventory (轻量级订单库存管理系统项目)

# 1.项目简介

本项目基于 Go + Gin + Gorm + MySQL 实现轻量级订单库存管理系统，支出商品创建、商品查询、商品上架、商品下架...

项目采用 handler/ service/ dao/ model/ request/	response/ 分层结构，使用环境变量管理数据库配置，并通过统一请求、响应和业务逻辑错误，提升接口规范性。 

# 2.技术栈

- GO

- Gin

- MySQL

- Gorm

- godotenv

- Yaml配置

# 3.核心功能

- 创建商品

- 查询商品

- 商品上架 / 下架

# 4.项目结构

```

cmd/	项目启动的主入口

config/	Yaml 配置加载

docs/	存放项目手动测试 http 和初始化 SQL 的文件夹

global/	全局资源，如 DB

pkg/	外部引用的资源，如Mysql

router/	全局路由

internal/	资源接口 REST API

internal/handler	HTTP 接口层，负责请求处理、数据返回

internal/service	业务逻辑层，包含商品创建、商品查询、商品上架\下架等相关业务逻辑

internal/dao	数据库操作层，负责数据库操作

internal/model	model 层，负责存放数据结构实体

internal/response	response 层，负责响应数据，返回数据结果

internal/request	request 层，负责处理客户端请求的数据类型

```

# 5.SQL结构

# 6.环境变量

```

DB_USER=root

DB_PASSWORD=your_password

DB_URL=127.0.0.1

DB_PORT=3306

DB_NAME=go-order-inventory

```

# 7.启动方式

```

go mod tidy 

go run cmd/main.go

```

# 8.接口说明

# 9.手动测试流程

本项目使用 VS Code REST Client 插件维护手动接口测试文件。

测试文件位置：

```text
docs/http/products.http
```

# 10.最终自测清单

# 11.设计与实现要点

