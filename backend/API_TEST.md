# X-Store 支付功能测试指南

## 前置准备

### 1. 启动服务

```bash
# 启动 MySQL 和 Redis
cd /Volumes/dz/code/x-store
docker-compose up -d

# 等待 MySQL 完全启动（约 15 秒）
sleep 15

# 导入测试数据
mysql -h 127.0.0.1 -u root -proot x_store < backend/init_test_data.sql

# 启动后端服务
cd backend
go run cmd/main.go
```

### 2. 验证服务启动

```bash
curl http://localhost:8080/api/health
# 预期返回: {"status":"ok"}
```

---

## 完整支付流程测试

### 步骤 1: 查询商品列表

```bash
curl http://localhost:8080/api/products
```

**预期返回：** 包含 6 个商品的列表

### 步骤 2: 查询支付渠道

```bash
curl http://localhost:8080/api/payment-channels
```

**预期返回：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "模拟支付（测试）",
      "provider_type": "mock",
      "channel_type": "mockpay",
      "is_active": true
    }
  ]
}
```

### 步骤 3: 创建订单

**注意：** 订单创建接口需要防重放签名，暂时跳过签名验证测试

```bash
# 临时方案：直接在数据库插入测试订单
mysql -h 127.0.0.1 -u root -proot x_store -e "
INSERT INTO orders (order_no, product_id, email, amount, status, expire_at, created_at, updated_at) 
VALUES ('XS20260314001', 1, 'test@example.com', 120.00, 0, DATE_ADD(NOW(), INTERVAL 15 MINUTE), NOW(), NOW());
"
```

### 步骤 4: 创建支付

```bash
curl -X POST http://localhost:8080/api/payments \
  -H "Content-Type: application/json" \
  -H "X-Timestamp: $(date +%s)" \
  -H "X-Nonce: test-nonce-$(date +%s)" \
  -H "X-Signature: test-signature" \
  -d '{
    "order_no": "XS20260314001",
    "channel_id": 1
  }'
```

**预期返回：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "payment_id": 1,
    "pay_url": "/mock-pay?trade_no=MOCK1710396000001&order_no=XS20260314001&amount=120.00&payment_id=1",
    "qr_code": "/mock-pay?trade_no=MOCK1710396000001&order_no=XS20260314001&amount=120.00&payment_id=1",
    "trade_no": "MOCK1710396000001",
    "expire_at": 1710396900
  }
}
```

### 步骤 5: 查询支付状态

```bash
curl http://localhost:8080/api/payments/1/status
```

**预期返回：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "order_id": 1,
    "trade_no": "MOCK1710396000001",
    "pay_method": "mockpay",
    "amount": 120.00,
    "status": 0,
    "created_at": "2026-03-14T15:40:00+08:00"
  }
}
```

### 步骤 6: 模拟支付回调（支付成功）

```bash
curl -X POST http://localhost:8080/api/webhooks/payment/1 \
  -H "Content-Type: application/json" \
  -d '{
    "trade_no": "MOCK1710396000001",
    "order_no": "XS20260314001",
    "payment_id": 1,
    "amount": 120.00,
    "status": "success",
    "paid_at": "2026-03-14T15:40:30+08:00",
    "sign": "mock_signature_here"
  }'
```

**预期返回：**
```json
{
  "code": 0,
  "message": "回调处理成功"
}
```

### 步骤 7: 验证订单状态

```bash
curl http://localhost:8080/api/orders/XS20260314001
```

**预期返回：** 订单状态应为 `2` (已发货)，包含卡密信息

### 步骤 8: 验证卡密已分配

```bash
mysql -h 127.0.0.1 -u root -proot x_store -e "
SELECT ck.id, ck.product_id, ck.content, ck.status, ck.order_id 
FROM card_keys ck 
WHERE ck.order_id = 1;
"
```

**预期结果：** 应该有 1 条卡密记录，状态为 `2` (已售出)

---

## 测试场景

### 场景 1: 正常购买流程

1. ✅ 查询商品
2. ✅ 创建订单
3. ✅ 创建支付
4. ✅ 支付成功回调
5. ✅ 自动发货
6. ✅ 订单完成

### 场景 2: 支付失败流程

```bash
# 模拟支付失败回调
curl -X POST http://localhost:8080/api/webhooks/payment/1 \
  -H "Content-Type: application/json" \
  -d '{
    "trade_no": "MOCK1710396000001",
    "order_no": "XS20260314001",
    "payment_id": 1,
    "amount": 120.00,
    "status": "failed",
    "paid_at": "2026-03-14T15:40:30+08:00",
    "sign": "mock_signature_here"
  }'
```

### 场景 3: 重复回调（幂等性测试）

多次调用步骤 6 的回调接口，验证：
- ✅ 不会重复扣减库存
- ✅ 不会重复分配卡密
- ✅ 订单状态保持一致

### 场景 4: 库存不足

1. 将商品库存设置为 0
2. 尝试创建订单
3. 预期：返回"库存不足"错误

---

## 数据库验证

### 查看订单状态

```sql
SELECT id, order_no, status, amount, paid_at, created_at 
FROM orders 
ORDER BY id DESC 
LIMIT 10;
```

### 查看支付记录

```sql
SELECT id, order_id, trade_no, pay_method, amount, status, completed_at 
FROM payments 
ORDER BY id DESC 
LIMIT 10;
```

### 查看卡密分配情况

```sql
SELECT product_id, 
       COUNT(*) as total,
       SUM(CASE WHEN status = 0 THEN 1 ELSE 0 END) as available,
       SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) as locked,
       SUM(CASE WHEN status = 2 THEN 1 ELSE 0 END) as sold
FROM card_keys 
GROUP BY product_id;
```

### 查看 Redis 库存

```bash
redis-cli -h 127.0.0.1 -p 6379
> GET stock:1
> GET stock:2
```

---

## 下一步：对接真实支付渠道

### Stripe 支付

1. 注册 Stripe 账号
2. 获取 API 密钥
3. 实现 `pkg/payment/stripe.go`
4. 添加 Stripe Webhook 处理

### 支付宝

1. 申请支付宝商户
2. 获取应用密钥
3. 实现 `pkg/payment/alipay.go`
4. 添加支付宝异步通知处理

### 微信支付

1. 申请微信商户
2. 获取商户密钥
3. 实现 `pkg/payment/wechatpay.go`
4. 添加微信支付回调处理

---

## 故障排查

### 问题 1: 订单创建失败

**检查：**
- MySQL 是否正常运行
- 商品库存是否充足
- Redis 库存是否同步

### 问题 2: 支付回调失败

**检查：**
- 签名验证是否通过
- 订单状态是否正确
- 支付金额是否匹配

### 问题 3: 卡密未发货

**检查：**
- 卡密库存是否充足
- 订单状态是否为已支付
- 数据库事务是否正常提交

---

## API 接口清单

| 方法 | 路径 | 说明 | 状态 |
|------|------|------|------|
| GET | `/api/products` | 商品列表 | ✅ |
| GET | `/api/payment-channels` | 支付渠道列表 | ✅ |
| POST | `/api/orders` | 创建订单 | ✅ |
| POST | `/api/payments` | 创建支付 | ✅ |
| GET | `/api/payments/:id/status` | 查询支付状态 | ✅ |
| POST | `/api/webhooks/payment/:channel_id` | 支付回调 | ✅ |
| GET | `/api/orders/:order_no` | 查询订单 | ✅ |
