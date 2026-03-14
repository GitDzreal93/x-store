-- X-Store PostgreSQL 初始化脚本
-- 执行方式: psql -h localhost -p 5432 -U admin -d x_store -f init_postgres.sql

-- 清理旧数据（可选，谨慎使用）
-- TRUNCATE TABLE card_keys, payments, orders, product_details, products, categories, users, payment_channels RESTART IDENTITY CASCADE;

-- ==================== 插入管理员账号 ====================
INSERT INTO users (username, email, password, nickname, role, status, oauth_provider, oauth_id, created_at, updated_at) 
VALUES (
    'admin', 
    'admin@x-store.com', 
    '$2a$10$l4NijSUwgSfF3rMRrzl17uEMhtRV1ssOXD1hGL.zSsgAyJQIsT1mW', -- 密码: admin123
    '超级管理员', 
    'admin', 
    1,
    '',
    '',
    NOW(), 
    NOW()
) ON CONFLICT (username) DO NOTHING;

-- ==================== 插入分类数据 ====================
INSERT INTO categories (name, icon, sort, status, created_at, updated_at) VALUES
('AI 大模型', '🤖', 100, 1, NOW(), NOW()),
('流媒体账号', '🎬', 90, 1, NOW(), NOW()),
('游戏充值', '🎮', 80, 1, NOW(), NOW()),
('礼品卡', '🎁', 70, 1, NOW(), NOW()),
('办公软件', '💼', 60, 1, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- ==================== 插入商品数据 ====================
INSERT INTO products (category_id, title, cover, price, original_price, stock, sales, delivery_type, tags, sort, status, is_new, is_hot, created_at, updated_at) VALUES
(1, 'ChatGPT Plus 独享账号-月付', 'https://images.unsplash.com/photo-1677442136019-21780ecad995?w=400', 35.00, 40.00, 50, 128, 'auto', '["自动发货","独享资源","7天质保"]', 100, 1, true, true, NOW(), NOW()),
(1, 'Claude Pro 订阅-年付', 'https://images.unsplash.com/photo-1676299081847-1d0c46c6e04a?w=400', 199.00, 240.00, 30, 56, 'auto', '["自动发货","年付优惠"]', 90, 1, false, true, NOW(), NOW()),
(2, 'Netflix 4K 高级会员-30天', 'https://images.unsplash.com/photo-1574375927938-d5a98e8efe85?w=400', 25.00, 35.00, 200, 512, 'auto', '["自动发货","4K画质"]', 80, 1, true, false, NOW(), NOW()),
(2, 'Disney+ 高级会员-月付', 'https://images.unsplash.com/photo-1611162617474-5b21e879e113?w=400', 18.00, 25.00, 150, 234, 'auto', '["自动发货","全球内容"]', 75, 1, false, false, NOW(), NOW()),
(3, 'Steam 余额卡密-100元', 'https://images.unsplash.com/photo-1612287230217-8c7c6c300a8c?w=400', 95.00, 100.00, 80, 89, 'auto', '["自动发货","官方卡密"]', 70, 1, false, false, NOW(), NOW()),
(3, 'PlayStation Plus 会员-12个月', 'https://images.unsplash.com/photo-1606144042614-b2417e99c4e3?w=400', 298.00, 398.00, 45, 67, 'auto', '["自动发货","年度会员"]', 65, 1, true, true, NOW(), NOW()),
(4, 'Apple Gift Card $50', 'https://images.unsplash.com/photo-1611532736579-6b16e2b50449?w=400', 350.00, 380.00, 60, 123, 'auto', '["自动发货","美区可用"]', 60, 1, false, true, NOW(), NOW()),
(5, 'Microsoft Office 365 家庭版-年付', 'https://images.unsplash.com/photo-1633419461186-7d40a38105ec?w=400', 299.00, 398.00, 100, 456, 'auto', '["自动发货","6人共享"]', 55, 1, true, true, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- ==================== 插入商品详情 ====================
INSERT INTO product_details (product_id, description, notice) VALUES
(1, 'ChatGPT Plus 独享账号，支持 GPT-4、DALL-E、代码解释器等全部功能。账号为独享使用，不与他人共享，确保使用体验。', '1. 账号为独享，禁止共享给他人使用\n2. 支持 7 天质保，如有问题请及时联系客服\n3. 请妥善保管账号信息，不要随意修改密码\n4. 建议使用前先测试功能是否正常'),
(2, 'Claude Pro 年度订阅，无限 Claude 3.5 Sonnet 使用权限。年付套餐更优惠，一次性支付享受 12 个月服务。', '1. 年付套餐更优惠，平均每月仅需 16.6 元\n2. 一次性支付享受 12 个月服务\n3. 支持所有 Claude 3.5 功能\n4. 账号独享，不与他人共享'),
(3, 'Netflix 4K 高级会员，支持 4 个屏幕同时观看，杜比视界+杜比全景声。账号已开通高级会员，直接使用即可。', '1. 支持 4K 超高清画质\n2. 可同时在 4 个设备上观看\n3. 支持杜比视界和杜比全景声\n4. 账号已开通，收到后直接登录使用'),
(4, 'Disney+ 高级会员，畅享迪士尼、漫威、星球大战等全球热门内容。', '1. 包含迪士尼、漫威、星球大战等内容\n2. 支持高清画质\n3. 多设备同时观看\n4. 月付套餐，灵活续费'),
(5, 'Steam 官方充值卡密，中国区可用，直接充值到钱包余额。卡密为官方正版，安全可靠。', '1. 官方正版卡密，安全可靠\n2. 中国区可用，直接充值到钱包\n3. 卡密一经售出不退不换\n4. 请确认Steam账号地区为中国区'),
(6, 'PlayStation Plus 会员年度订阅，畅玩 PS4/PS5 在线多人游戏，每月免费游戏。', '1. 支持 PS4 和 PS5 平台\n2. 每月免费游戏\n3. 在线多人游戏\n4. 云存档功能'),
(7, 'Apple Gift Card 美区礼品卡，可用于购买 App、游戏、音乐、电影等。', '1. 美区 Apple ID 可用\n2. 可购买 App、游戏、订阅服务\n3. 卡密一经售出不退不换\n4. 请确认账号地区为美国'),
(8, 'Microsoft Office 365 家庭版年度订阅，支持 6 人共享，包含 Word、Excel、PowerPoint、OneDrive 等。', '1. 支持 6 人共享使用\n2. 包含全套 Office 应用\n3. 每人 1TB OneDrive 云存储\n4. 支持 Windows、Mac、移动端')
ON CONFLICT DO NOTHING;

-- ==================== 插入支付渠道 ====================
INSERT INTO payment_channels (name, provider_type, channel_type, interaction_mode, config_json, fee_rate, fixed_fee, is_active, sort, created_at, updated_at) VALUES
-- 1. 模拟支付（测试用）
('模拟支付（测试）', 'mock', 'mockpay', 'redirect', '{"secret_key":"mock-secret-key-demo","notify_url":"http://localhost:8082/api/webhooks/payment/1","auto_notify":true}', 0.00, 0.00, true, 100, NOW(), NOW()),

-- 2. 支付宝当面付（扫码）
('支付宝当面付', 'alipay', 'alipay_f2f', 'qrcode', '{"app_id":"your_alipay_app_id","private_key":"your_private_key","ali_public_key":"your_ali_public_key","notify_url":"http://localhost:8082/api/webhooks/payment/2","return_url":"http://localhost:3000","is_sandbox":true}', 0.60, 0.00, false, 95, NOW(), NOW()),

-- 3. 支付宝电脑网站支付
('支付宝PC支付', 'alipay', 'alipay_pc', 'redirect', '{"app_id":"your_alipay_app_id","private_key":"your_private_key","ali_public_key":"your_ali_public_key","notify_url":"http://localhost:8082/api/webhooks/payment/3","return_url":"http://localhost:3000","is_sandbox":true}', 0.60, 0.00, false, 94, NOW(), NOW()),

-- 4. 支付宝手机网站支付
('支付宝H5支付', 'alipay', 'alipay_wap', 'h5', '{"app_id":"your_alipay_app_id","private_key":"your_private_key","ali_public_key":"your_ali_public_key","notify_url":"http://localhost:8082/api/webhooks/payment/4","return_url":"http://localhost:3000","is_sandbox":true}', 0.60, 0.00, false, 93, NOW(), NOW()),

-- 5. 微信Native扫码支付
('微信扫码支付', 'wechat', 'wechat_native', 'qrcode', '{"app_id":"your_wx_appid","mch_id":"your_mch_id","api_key":"your_api_key","api_key_v3":"","serial_no":"","notify_url":"http://localhost:8082/api/webhooks/payment/5"}', 0.60, 0.00, false, 90, NOW(), NOW()),

-- 6. 微信H5支付
('微信H5支付', 'wechat', 'wechat_h5', 'h5', '{"app_id":"your_wx_appid","mch_id":"your_mch_id","api_key":"your_api_key","api_key_v3":"","serial_no":"","notify_url":"http://localhost:8082/api/webhooks/payment/6"}', 0.60, 0.00, false, 89, NOW(), NOW()),

-- 7. Stripe 国际信用卡
('Stripe 信用卡', 'stripe', 'stripe', 'redirect', '{"secret_key":"sk_test_xxx","publishable_key":"pk_test_xxx","webhook_secret":"whsec_xxx"}', 2.90, 0.30, false, 85, NOW(), NOW()),

-- 8. PayPal 国际支付
('PayPal', 'paypal', 'paypal', 'redirect', '{"client_id":"your_paypal_client_id","client_secret":"your_paypal_secret","is_sandbox":true,"return_url":"http://localhost:3000/order/success","cancel_url":"http://localhost:3000"}', 3.49, 0.49, false, 80, NOW(), NOW()),

-- 9. USDT/EPUSDT 加密货币
('USDT(TRC20)', 'epusdt', 'epusdt', 'qrcode', '{"api_url":"http://your-epusdt-server","token":"your_epusdt_token","notify_url":"http://localhost:8082/api/webhooks/payment/9","return_url":"http://localhost:3000"}', 0.00, 0.00, false, 75, NOW(), NOW()),

-- 10. 易支付-支付宝
('易支付-支付宝', 'yipay', 'yipay_alipay', 'redirect', '{"api_url":"https://pay.example.com","pid":"your_pid","key":"your_key","notify_url":"http://localhost:8082/api/webhooks/payment/10","return_url":"http://localhost:3000"}', 1.50, 0.00, false, 70, NOW(), NOW()),

-- 11. 易支付-微信
('易支付-微信', 'yipay', 'yipay_wechat', 'redirect', '{"api_url":"https://pay.example.com","pid":"your_pid","key":"your_key","notify_url":"http://localhost:8082/api/webhooks/payment/11","return_url":"http://localhost:3000"}', 1.50, 0.00, false, 69, NOW(), NOW()),

-- 12. Payjs 微信个人
('Payjs微信', 'payjs', 'payjs', 'qrcode', '{"mchid":"your_payjs_mchid","key":"your_payjs_key","notify_url":"http://localhost:8082/api/webhooks/payment/12"}', 2.38, 0.00, false, 65, NOW(), NOW()),

-- 13. 虎皮椒-微信
('虎皮椒-微信', 'xunhupay', 'xunhupay_wechat', 'qrcode', '{"appid":"your_xunhu_appid","appsecret":"your_xunhu_secret","wechat_api_url":"https://api.xunhupay.com/payment/do.html","alipay_api_url":"","notify_url":"http://localhost:8082/api/webhooks/payment/13","return_url":"http://localhost:3000"}', 1.50, 0.00, false, 60, NOW(), NOW()),

-- 14. 虎皮椒-支付宝
('虎皮椒-支付宝', 'xunhupay', 'xunhupay_alipay', 'redirect', '{"appid":"your_xunhu_appid","appsecret":"your_xunhu_secret","wechat_api_url":"","alipay_api_url":"https://api.xunhupay.com/payment/do.html","notify_url":"http://localhost:8082/api/webhooks/payment/14","return_url":"http://localhost:3000"}', 1.50, 0.00, false, 59, NOW(), NOW()),

-- 15. V免签-微信
('V免签-微信', 'vmqpay', 'vmqpay_wechat', 'qrcode', '{"api_url":"http://your-vmq-server","key":"your_vmq_key","notify_url":"http://localhost:8082/api/webhooks/payment/15","return_url":"http://localhost:3000"}', 0.00, 0.00, false, 55, NOW(), NOW()),

-- 16. V免签-支付宝
('V免签-支付宝', 'vmqpay', 'vmqpay_alipay', 'qrcode', '{"api_url":"http://your-vmq-server","key":"your_vmq_key","notify_url":"http://localhost:8082/api/webhooks/payment/16","return_url":"http://localhost:3000"}', 0.00, 0.00, false, 54, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- ==================== 插入 OAuth 提供商配置 ====================
INSERT INTO oauth_providers (provider, name, enabled, client_id, client_secret, redirect_url, sort, created_at, updated_at) VALUES
('github', 'GitHub', false, 'your_github_client_id', 'your_github_client_secret', 'http://localhost:8082/api/oauth/github/callback', 100, NOW(), NOW()),
('google', 'Google', false, 'your_google_client_id', 'your_google_client_secret', 'http://localhost:8082/api/oauth/google/callback', 90, NOW(), NOW())
ON CONFLICT (provider) DO NOTHING;

-- ==================== 插入卡密数据 ====================
INSERT INTO card_keys (product_id, content, status, created_at, updated_at) VALUES
-- ChatGPT Plus 卡密
(1, 'CHATGPT-PLUS-DEMO-001-ABCD1234', 0, NOW(), NOW()),
(1, 'CHATGPT-PLUS-DEMO-002-EFGH5678', 0, NOW(), NOW()),
(1, 'CHATGPT-PLUS-DEMO-003-IJKL9012', 0, NOW(), NOW()),
(1, 'CHATGPT-PLUS-DEMO-004-MNOP3456', 0, NOW(), NOW()),
(1, 'CHATGPT-PLUS-DEMO-005-QRST7890', 0, NOW(), NOW()),

-- Claude Pro 卡密
(2, 'CLAUDE-PRO-YEAR-DEMO-001-UVWX1234', 0, NOW(), NOW()),
(2, 'CLAUDE-PRO-YEAR-DEMO-002-YZAB5678', 0, NOW(), NOW()),
(2, 'CLAUDE-PRO-YEAR-DEMO-003-CDEF9012', 0, NOW(), NOW()),

-- Netflix 卡密
(3, 'NETFLIX-4K-DEMO-001-GHIJ3456', 0, NOW(), NOW()),
(3, 'NETFLIX-4K-DEMO-002-KLMN7890', 0, NOW(), NOW()),
(3, 'NETFLIX-4K-DEMO-003-OPQR1234', 0, NOW(), NOW()),
(3, 'NETFLIX-4K-DEMO-004-STUV5678', 0, NOW(), NOW()),
(3, 'NETFLIX-4K-DEMO-005-WXYZ9012', 0, NOW(), NOW()),

-- Disney+ 卡密
(4, 'DISNEY-PLUS-DEMO-001-ABCD3456', 0, NOW(), NOW()),
(4, 'DISNEY-PLUS-DEMO-002-EFGH7890', 0, NOW(), NOW()),
(4, 'DISNEY-PLUS-DEMO-003-IJKL1234', 0, NOW(), NOW()),

-- Steam 卡密
(5, 'STEAM-100CNY-DEMO-001-MNOP5678', 0, NOW(), NOW()),
(5, 'STEAM-100CNY-DEMO-002-QRST9012', 0, NOW(), NOW()),
(5, 'STEAM-100CNY-DEMO-003-UVWX3456', 0, NOW(), NOW()),
(5, 'STEAM-100CNY-DEMO-004-YZAB7890', 0, NOW(), NOW()),

-- PlayStation Plus 卡密
(6, 'PS-PLUS-YEAR-DEMO-001-CDEF1234', 0, NOW(), NOW()),
(6, 'PS-PLUS-YEAR-DEMO-002-GHIJ5678', 0, NOW(), NOW()),

-- Apple Gift Card 卡密
(7, 'APPLE-GIFT-50USD-DEMO-001-KLMN9012', 0, NOW(), NOW()),
(7, 'APPLE-GIFT-50USD-DEMO-002-OPQR3456', 0, NOW(), NOW()),
(7, 'APPLE-GIFT-50USD-DEMO-003-STUV7890', 0, NOW(), NOW()),

-- Office 365 卡密
(8, 'OFFICE365-FAMILY-DEMO-001-WXYZ1234', 0, NOW(), NOW()),
(8, 'OFFICE365-FAMILY-DEMO-002-ABCD5678', 0, NOW(), NOW()),
(8, 'OFFICE365-FAMILY-DEMO-003-EFGH9012', 0, NOW(), NOW()),
(8, 'OFFICE365-FAMILY-DEMO-004-IJKL3456', 0, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- ==================== 插入测试用户 ====================
INSERT INTO users (username, email, password, nickname, role, status, created_at, updated_at) 
VALUES (
    'testuser', 
    'test@example.com', 
    '$2a$10$l4NijSUwgSfF3rMRrzl17uEMhtRV1ssOXD1hGL.zSsgAyJQIsT1mW', -- 密码: admin123
    '测试用户', 
    'buyer', 
    1, 
    NOW(), 
    NOW()
) ON CONFLICT (username) DO NOTHING;

-- ==================== 查看数据统计 ====================
SELECT '数据初始化完成！' as message;
SELECT '用户数量: ' || COUNT(*) FROM users;
SELECT '分类数量: ' || COUNT(*) FROM categories;
SELECT '商品数量: ' || COUNT(*) FROM products;
SELECT '卡密数量: ' || COUNT(*) FROM card_keys;
SELECT '支付渠道: ' || COUNT(*) FROM payment_channels;
