---
sidebar_position: 1
title: 环境准备
---

# 环境准备

在开始之前，请确保你的开发环境满足以下要求。

## 必需软件

| 软件 | 最低版本 | 说明 |
|------|----------|------|
| **Go** | 1.21+ | 后端开发语言 |
| **Node.js** | 18+ | 前端构建工具 |
| **npm** | 9+ | 包管理器 |
| **PostgreSQL** | 14+ | 主数据库 |
| **Redis** | 6+ | 缓存 & 防重放 |
| **Git** | 2.0+ | 版本控制 |

## 安装指南

### macOS (Homebrew)

```bash
# Go
brew install go

# Node.js (推荐使用 nvm)
brew install nvm
nvm install 18

# PostgreSQL & Redis
brew install postgresql@16 redis
brew services start postgresql@16
brew services start redis
```

### Ubuntu / Debian

```bash
# Go
wget https://go.dev/dl/go1.21.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Node.js (nvm)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install 18

# PostgreSQL & Redis
sudo apt update
sudo apt install -y postgresql redis-server
sudo systemctl enable postgresql redis-server
```

### Windows

推荐使用 [WSL2](https://learn.microsoft.com/zh-cn/windows/wsl/) + Ubuntu，参考上方 Ubuntu 安装步骤。

## 验证安装

```bash
go version          # go version go1.21.x ...
node -v             # v18.x.x
npm -v              # 9.x.x
psql --version      # psql (PostgreSQL) 16.x
redis-cli ping      # PONG
```

## 创建数据库

```bash
# 进入 PostgreSQL
psql -U postgres

# 创建数据库和用户
CREATE USER admin WITH PASSWORD 'admin123';
CREATE DATABASE x_store OWNER admin;
GRANT ALL PRIVILEGES ON DATABASE x_store TO admin;
\q
```

:::tip
如果你已经有 PostgreSQL 运行中，只需确保创建了 `x_store` 数据库即可。数据库表会在后端首次启动时自动创建（GORM AutoMigrate）。
:::

## 下一步

👉 前往 [安装部署](/docs/getting-started/installation) 克隆项目并安装依赖。
