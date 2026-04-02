-- X-Store 测试数据初始化脚本 (PostgreSQL 版本)
-- 使用方法: docker exec postgresql psql -U admin -d x_store < init_test_data_pg.sql

-- 0. 插入管理员账号 (密码: admin123)
INSERT INTO users (username, email, password, nickname, role, status, created_at, updated_at)
VALUES ('admin', 'admin@x-store.com', '$2a$10$l4NijSUwgSfF3rMRrzl17uEMhtRV1ssOXD1hGL.zSsgAyJQIsT1mW', '超级管理员', 'admin', 1, NOW(), NOW())
ON CONFLICT (username) DO NOTHING;

-- 1. 插入分类数据
INSERT INTO categories (name, icon, sort, status, created_at, updated_at)
VALUES
  ('AI 大模型', '🤖', 100, 1, NOW(), NOW()),
  ('流媒体账号', '🎬', 90, 1, NOW(), NOW()),
  ('游戏充值', '🎮', 80, 1, NOW(), NOW()),
  ('礼品卡', '🎁', 70, 1, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- 2. 插入商品数据
INSERT INTO products (category_id, title, cover, price, original_price, stock, sales, delivery_type, tags, sort, status, is_new, is_hot, created_at, updated_at)
VALUES
  (1, 'ChatGPT Plus 独享账号', '', 120.00, 150.00, 100, 328, 'auto', '["自动发货","独享资源","7天质保"]'::jsonb, 100, 1, true, true, NOW(), NOW()),
  (1, 'Claude Pro 会员账号', '', 98.00, 128.00, 50, 156, 'auto', '["自动发货","独享资源"]'::jsonb, 90, 1, true, false, NOW(), NOW()),
  (2, 'Netflix 高级会员 1个月', '', 35.00, 45.00, 200, 892, 'auto', '["自动发货","4K画质"]'::jsonb, 80, 1, false, true, NOW(), NOW()),
  (2, 'Spotify Premium 家庭版', '', 28.00, 35.00, 150, 445, 'auto', '["自动发货","6人共享"]'::jsonb, 70, 1, false, false, NOW(), NOW()),
  (3, 'Steam 充值卡 100元', '', 95.00, 100.00, 300, 1203, 'auto', '["自动发货","秒到账"]'::jsonb, 60, 1, false, true, NOW(), NOW()),
  (4, 'Apple Gift Card $50', '', 320.00, 350.00, 80, 234, 'auto', '["自动发货","全球通用"]'::jsonb, 50, 1, false, false, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- 3. 插入商品详情
INSERT INTO product_details (product_id, description, notice, created_at, updated_at)
VALUES
  (1, '# ChatGPT Plus 独享账号

## 产品特点
- 独享账号，无需共享
- 支持 GPT-4 模型
- 无限次数使用
- 7天质保服务

## 使用说明
1. 购买后自动发货到邮箱
2. 使用账号密码登录
3. 如有问题联系客服', '1. 请勿修改密码
2. 请勿分享账号
3. 7天内如有问题可联系客服', NOW(), NOW()),
  (2, '# Claude Pro 会员账号

## 产品特点
- Anthropic 官方会员
- 支持 Claude 3 Opus
- 独享账号使用

## 使用说明
购买后自动发货', '请勿修改密码', NOW(), NOW()),
  (3, '# Netflix 高级会员

## 产品特点
- 4K 超高清画质
- 支持4个设备同时观看
- 全球影视资源

## 使用说明
自动发货账号密码', '请勿修改密码', NOW(), NOW()),
  (4, '# Spotify Premium 家庭版

## 产品特点
- 6人共享
- 无广告播放
- 离线下载

## 使用说明
自动发货邀请链接', '请按说明加入家庭组', NOW(), NOW()),
  (5, '# Steam 充值卡

## 产品特点
- 官方充值卡
- 秒到账
- 全区通用

## 使用说明
自动发货卡密', '请在Steam客户端兑换', NOW(), NOW()),
  (6, '# Apple Gift Card

## 产品特点
- 官方礼品卡
- 全球通用
- 可购买App/音乐/游戏

## 使用说明
自动发货卡密', '请在Apple账户中兑换', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- 4. 插入支付渠道配置（模拟支付）
INSERT INTO payment_channels (name, provider_type, channel_type, interaction_mode, config_json, fee_rate, fixed_fee, is_active, sort, created_at, updated_at)
VALUES ('模拟支付（测试）', 'mock', 'mockpay', 'redirect', '{"secret":"mock-secret-key-2026","notify_url":"http://localhost:8082/api/webhooks/payment/1","return_url":"http://localhost:3000","auto_notify":true}'::jsonb, 0.00, 0.00, true, 100, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- 5. 插入测试卡密（为商品1-6添加卡密）
INSERT INTO card_keys (product_id, content, status, created_at, updated_at)
VALUES
  -- ChatGPT Plus (product_id = 1)
  (1, 'CHATGPT-PLUS-KEY-001', 0, NOW(), NOW()),
  (1, 'CHATGPT-PLUS-KEY-002', 0, NOW(), NOW()),
  (1, 'CHATGPT-PLUS-KEY-003', 0, NOW(), NOW()),
  (1, 'CHATGPT-PLUS-KEY-004', 0, NOW(), NOW()),
  (1, 'CHATGPT-PLUS-KEY-005', 0, NOW(), NOW()),
  -- Claude Pro (product_id = 2)
  (2, 'CLAUDE-PRO-KEY-001', 0, NOW(), NOW()),
  (2, 'CLAUDE-PRO-KEY-002', 0, NOW(), NOW()),
  (2, 'CLAUDE-PRO-KEY-003', 0, NOW(), NOW()),
  -- Netflix (product_id = 3)
  (3, 'NETFLIX-ACCOUNT-001', 0, NOW(), NOW()),
  (3, 'NETFLIX-ACCOUNT-002', 0, NOW(), NOW()),
  -- Spotify (product_id = 4)
  (4, 'SPOTIFY-INVITE-001', 0, NOW(), NOW()),
  (4, 'SPOTIFY-INVITE-002', 0, NOW(), NOW()),
  -- Steam (product_id = 5)
  (5, 'STEAM-CARD-001', 0, NOW(), NOW()),
  (5, 'STEAM-CARD-002', 0, NOW(), NOW()),
  -- Apple Gift Card (product_id = 6)
  (6, 'APPLE-CARD-001', 0, NOW(), NOW()),
  (6, 'APPLE-CARD-002', 0, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- 查询插入结果
SELECT '=== 分类数据 ===' as info;
SELECT id, name, icon, status FROM categories;

SELECT '=== 商品数据 ===' as info;
SELECT id, title, price, stock, sales, status FROM products;

SELECT '=== 支付渠道 ===' as info;
SELECT id, name, provider_type, channel_type, is_active FROM payment_channels;

SELECT '=== 卡密库存 ===' as info;
SELECT product_id, COUNT(*) as total,
       SUM(CASE WHEN status = 0 THEN 1 ELSE 0 END) as available,
       SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) as locked,
       SUM(CASE WHEN status = 2 THEN 1 ELSE 0 END) as sold
FROM card_keys
GROUP BY product_id;
