# 项目问题清单与复盘

## 1. 已解决问题

### 1. Response 放置位置问题

问题：最初 Response 结构体放在 model 层，导致 model 层混入 HTTP 响应结构。

修正：将 Response 结构和 Success / Fail 方法统一放入 internal/response。

结果：model 层只保留数据库实体，response 层负责统一响应。

### 2. 商品创建返回 nil 问题

问题：创建商品成功后只返回 success，不返回商品 ID。

修正：创建成功后返回 product。

结果：后续库存初始化可以直接使用返回的 product_id。

### 3. 商品 ID 校验问题

问题：最初只解析 id 是否为数字，没有判断 id <= 0。

修正：增加 parsePositiveProductID，统一校验 id 必须为正整数。

结果：避免 /products/0、/orders/0 等非法 ID 进入 service 层。

### 4. 库存初始化事务问题

问题：初始化库存和写库存流水最初不在同一个事务内。

修正：使用 global.DB.Transaction，并在事务内使用 tx 创建 inventory 和 stock_log。

结果：库存记录和库存流水可以保持一致。

### 5. 增加库存覆盖库存问题

问题：增加库存最初把 stock_quantity 更新为 quantity，而不是累加。

修正：先读取库存 before_quantity，再计算 after_quantity = before + quantity。

结果：库存数量和库存流水保持一致。

### 6. 订单创建事务问题

问题：订单创建涉及订单主表、订单项、库存扣减、库存流水，必须保证一致性。

修正：将订单创建、库存扣减、订单项创建、库存流水写入放在同一个事务内。

结果：任一环节失败时事务回滚。

## 2. 当前仍需优化问题

### 1. 并发扣库存时 before/after 流水准确性

当前通过条件扣减避免库存扣成负数，但并发场景下 stock_logs 的 before_quantity / after_quantity 可能不完全准确。

后续优化：

- 在事务中使用 SELECT ... FOR UPDATE 锁住库存行
- 基于锁定后的库存计算 before/after
- 再更新库存并写流水

### 2. Handler 错误码需要进一步统一

部分状态流转错误需要统一使用 409 Conflict，而不是 404。

### 3. 自动化测试缺失

当前主要依赖 REST Client 手动测试，后续需要补充：

- service 层业务规则测试
- handler 层参数错误测试
- 订单状态机测试
- 库存事务测试