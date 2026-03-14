---
sidebar_position: 2
title: 生产环境
---

# 生产环境部署

本指南介绍在生产环境中部署 X-Store 的注意事项和最佳实践。

## 安全配置清单

部署前务必修改以下配置：

```yaml title="backend/config.yaml"
server:
  mode: release          # ❗ 切换为 release 模式

jwt:
  secret: <随机生成的强密钥>   # ❗ 必须修改
  expire_hours: 24

signature:
  secret: <随机生成的强密钥>   # ❗ 必须修改
  expire_seconds: 60
```

生成随机密钥：

```bash
openssl rand -hex 32
```

## Nginx 反向代理

推荐使用 Nginx 统一代理所有服务：

```nginx title="/etc/nginx/conf.d/x-store.conf"
# C 端前端
server {
    listen 80;
    server_name store.yourdomain.com;

    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

# 后端 API
server {
    listen 80;
    server_name api.yourdomain.com;

    location / {
        proxy_pass http://127.0.0.1:8082;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}

# 管理后台
server {
    listen 80;
    server_name admin.yourdomain.com;

    root /var/www/admin-panel/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

## HTTPS 配置

使用 Let's Encrypt 免费证书：

```bash
# 安装 certbot
sudo apt install certbot python3-certbot-nginx

# 自动配置 HTTPS
sudo certbot --nginx -d store.yourdomain.com
sudo certbot --nginx -d api.yourdomain.com
sudo certbot --nginx -d admin.yourdomain.com
```

## 进程管理

使用 systemd 管理后端服务：

```ini title="/etc/systemd/system/x-store.service"
[Unit]
Description=X-Store Backend
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/x-store/backend
ExecStart=/opt/x-store/backend/x-store-backend
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable x-store
sudo systemctl start x-store
sudo systemctl status x-store
```

## 数据库维护

### 备份

```bash
# 每日自动备份
pg_dump -h localhost -U admin x_store | gzip > /backup/x_store_$(date +%Y%m%d).sql.gz
```

### 定时清理

添加 crontab 定时备份：

```bash
# 每天凌晨 3 点备份
0 3 * * * pg_dump -h localhost -U admin x_store | gzip > /backup/x_store_$(date +\%Y\%m\%d).sql.gz
# 保留最近 30 天
0 4 * * * find /backup -name "x_store_*.sql.gz" -mtime +30 -delete
```

## 监控

### 健康检查

```bash
# 后端健康检查
curl -s http://localhost:8082/api/categories | jq .code
# 返回 0 表示正常
```

### 日志管理

后端日志输出到 stdout，systemd 自动收集：

```bash
# 查看实时日志
journalctl -u x-store -f

# 查看最近 100 行
journalctl -u x-store -n 100
```

## 性能优化

1. **数据库连接池**：调整 `max_idle_conns` 和 `max_open_conns`
2. **Redis 缓存**：热点数据缓存到 Redis
3. **静态资源 CDN**：前端构建产物部署到 CDN
4. **Gzip 压缩**：Nginx 启用 gzip
5. **HTTP/2**：Nginx 配置 HTTP/2 支持

```nginx
server {
    listen 443 ssl http2;
    # ...
    gzip on;
    gzip_types text/plain application/json application/javascript text/css;
}
```
