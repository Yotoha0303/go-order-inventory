# 接口清单

## 1. 健康检查

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | /ping | 健康检查 |

## 2. 商品模块

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | /api/v1/products | 创建商品 |
| GET | /api/v1/products | 查询商品列表，当前默认查询下架商品 |
| GET | /api/v1/products/:id | 查询商品详情 |
| PATCH | /api/v1/products/:id/on-sale | 商品上架 |
| PATCH | /api/v1/products/:id/off-sale | 商品下架 |

## 3. 库存模块

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | /api/v1/inventory/init | 初始化库存 |
| POST | /api/v1/inventory/add | 增加库存 |
| GET | /api/v1/inventory/products/:product_id | 查询商品库存 |

## 4. 库存流水模块

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | /api/v1/stock-logs | 查询库存流水，product_id 可选 |

## 5. 订单模块

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | /api/v1/orders | 创建订单 |
| GET | /api/v1/orders | 查询订单列表 |
| GET | /api/v1/orders/:id | 查询订单详情 |
| PATCH | /api/v1/orders/:id/pay | 支付订单 |
| PATCH | /api/v1/orders/:id/finish | 完成订单 |
| PATCH | /api/v1/orders/:id/cancel | 取消订单 |