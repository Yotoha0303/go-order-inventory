# 项目测试说明

## 1. 测试类型

本项目当前采用三类测试方式：

1. REST Client 手动接口测试
2. 核心业务流程自测
3. 业务规则检查清单

后续可补充 Go 自动化测试，包括 service 层单元测试和 handler 层接口测试。

## 2. REST Client 接口测试

测试文件位置：

```text
docs/http/products.http
docs/http/inventory.http
docs/http/stock_logs.http
docs/http/orders.http
```

执行方式：

1. 安装 VS Code REST Client 插件
2. 启动项目：`go run cmd/main.go`
3. 打开对应 `.http` 文件
4. 点击每个请求上方的 `Send Request`
5. 对比响应结果和数据库变化

## 3. 商品模块自测

- [ ] 创建商品成功
- [ ] 创建商品后 `status = 2`
- [ ] `price_fen <= 0` 返回参数错误
- [ ] `name` 为空返回参数错误
- [ ] 查询商品列表成功
- [ ] 查询商品详情成功
- [ ] 查询不存在商品返回错误
- [ ] 商品上架成功
- [ ] 商品下架成功

## 4. 库存模块自测

- [ ] 存在商品可以初始化库存
- [ ] 不存在商品不能初始化库存
- [ ] 重复初始化库存失败
- [ ] `stock_quantity = 0` 可以初始化
- [ ] 初始化库存后 `product_inventories` 有记录
- [ ] 初始化库存后 `stock_logs` 有 `biz_type = 1` 记录
- [ ] 已初始化库存的商品可以增加库存
- [ ] 未初始化库存的商品不能增加库存
- [ ] `quantity <= 0` 返回参数错误
- [ ] 增加库存后 `stock_quantity` 正确变化
- [ ] 增加库存后 `stock_logs` 有 `biz_type = 2` 记录

## 5. 库存流水自测

- [ ] 不传 `product_id` 可以查询全部流水
- [ ] 传 `product_id` 可以查询指定商品流水
- [ ] `product_id` 非法返回参数错误
- [ ] 初始化库存后能查到 `biz_type = 1`
- [ ] 增加库存后能查到 `biz_type = 2`
- [ ] 创建订单后能查到 `biz_type = 3`
- [ ] 取消订单后能查到 `biz_type = 4`
- [ ] `before_quantity / change_quantity / after_quantity` 正确

## 6. 订单状态机测试

创建订单

- [ ] 正常创建订单成功
- [ ] 商品不存在时创建订单失败
- [ ] 商品下架时创建订单失败
- [ ] 库存不存在时创建订单失败
- [ ] 库存不足时创建订单失败
- [ ] 创建订单成功后 `orders` 有记录
- [ ] 创建订单成功后 `order_items` 有记录
- [ ] 创建订单成功后 `product_inventories` 库存扣减
- [ ] 创建订单成功后 `stock_logs` 有 `biz_type = 3` 记录

支付订单

- [ ] 待支付订单可以支付
- [ ] 已支付订单重复支付失败
- [ ] 已取消订单支付失败
- [ ] 已完成订单支付失败
- [ ] 不存在订单支付失败

完成订单

- [ ] 已支付订单可以完成
- [ ] 未支付订单完成失败
- [ ] 已取消订单完成失败
- [ ] 已完成订单重复完成失败
- [ ] 不存在订单完成失败

取消订单

- [ ] 待支付订单可以取消
- [ ] 取消订单后库存回滚
- [ ] 取消订单后 `stock_logs` 有 `biz_type = 4` 记录
- [ ] 已支付订单取消失败
- [ ] 已完成订单取消失败
- [ ] 不存在订单取消失败
