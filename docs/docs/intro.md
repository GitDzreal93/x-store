---
sidebar_position: 1
slug: /intro
title: 项目简介
---

# X-Store 项目简介

**X-Store** 是一个开源的数字商品自动化售卖系统，支持卡密自动发货、多渠道支付、第三方 OAuth 登录等功能。适用于虚拟商品（激活码、会员账号、礼品卡等）的在线销售场景。

## 核心特性

| 特性 | 说明 |
|------|------|
| 🛒 **商品管理** | 多分类商品、图片封面、标签筛选、富文本详情 |
| 🔑 **卡密管理** | 批量导入、自动分配、库存预警 |
| 💳 **多渠道支付** | 支付宝、微信、Stripe、PayPal 等 10 种支付方式 |
| 🔐 **OAuth 登录** | GitHub / Google 第三方登录，后台动态开关 |
| 📊 **数据统计** | GMV 趋势、订单分析、TOP 商品排行 |
| 📧 **邮件通知** | 支付成功自动邮件发货 |
| 🛡️ **安全机制** | JWT 认证、防重放攻击、接口限流 |

## 技术栈一览

```
后端:  Go 1.21 + Gin + GORM + PostgreSQL + Redis
前端:  Next.js 14 + React 18 + TailwindCSS + shadcn/ui
后台:  React 18 + Ant Design 5 + Vite
文档:  Docusaurus 3 (本站)
```

## 项目结构

```
x-store/
├── backend/           # Go 后端服务
├── frontend-store/    # C 端买家前端 (Next.js)
├── admin-panel/       # B 端管理后台 (React + Antd)
├── docs/              # 项目文档 (Docusaurus)
└── start.sh           # 一键启动脚本
```

## 适用场景

- 数字商品（激活码、CDK、会员账号）在线售卖
- 虚拟礼品卡 / 充值卡自动化发货
- 个人开发者 / 小型团队的轻量电商需求
- 学习全栈开发的实战项目

## 下一步

👉 前往 [环境准备](/docs/getting-started/prerequisites) 开始搭建你的 X-Store。
