# REST Client 自测结果

测试时间：2026-05-19
测试环境：本地 MySQL + Redis
启动命令：go run cmd/main.go

## 1. 商品模块

| 用例 | 结果 | 备注 |
|---|---|---|
| 创建商品 | 通过 | 新建 product_id=13，status=2 |
| 商品上架 | 通过 | `PATCH /products/13/on-sale` 成功 |
| 商品下架 | 通过 | `PATCH /products/13/off-sale` 成功 |

## 2. 库存模块

| 用例 | 结果 | 备注 |
|---|---|---|
| 初始化库存 | 通过 | `POST /inventory/init` 成功，stock_logs 记录 `biz_type=1` |
| 增加库存 | 通过 | `POST /inventory/add` 成功，stock_logs 记录 `biz_type=2` |
| 重复初始化 | 通过 | 返回业务错误码 `2001`（库存已初始化） |

## 3. 订单模块

| 用例 | 结果 | 备注 |
|---|---|---|
| 创建订单 | 通过 | 新建 order_id=4，扣减库存成功，stock_logs 有 `biz_type=3` |
| 支付订单 | 通过 | `pending -> paid` |
| 完成订单 | 通过 | `paid -> finished` |
| 取消订单 | 通过 | 针对 order_id=5 取消成功，库存回滚，stock_logs 有 `biz_type=4` |
| 重复取消 | 通过 | 再次取消 order_id=5 成功，无异常 |

## 4. Redis 缓存

| 用例 | 结果 | 备注 |
|---|---|---|
| 第一次查询商品详情 | 通过 | `GET /products/13` 成功（首次查询） |
| 第二次查询商品详情 | 通过 | `GET /products/13` 成功（重复查询） |
| 商品上下架后删除缓存 | 通过 | 下架后再次查询成功，缓存失效后可重建 |

## 5. Demo Flow 正常链路演示（demo_flow.http）

执行时间：2026-05-19 16:30:50

| 步骤 | 结果 | 备注 |
|---|---|---|
| 健康检查 `GET /ping` | 通过 | code=0 |
| 创建商品 | 通过 | product_id=2 |
| 商品上架 | 通过 | code=0 |
| 查询商品详情 | 通过 | code=0 |
| 初始化库存 | 通过 | stock_quantity=100 |
| 增加库存 | 通过 | 增加后库存=120 |
| 创建订单 | 通过 | order_id=5 |
| 支付订单 | 通过 | code=0 |
| 完成订单 | 通过 | code=0 |
| 查询订单详情 | 通过 | status=3（finished） |
| 复核库存变化 | 通过 | 下单2件后库存=118 |
