---
sidebar_position: 3
title: 商品与分类
---

# 商品与分类 API

## 分类列表

```
GET /api/categories
```

**响应**：

```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "name": "AI 大模型",
      "icon": "🤖",
      "sort": 100,
      "status": 1
    },
    {
      "id": 2,
      "name": "流媒体账号",
      "icon": "🎬",
      "sort": 90,
      "status": 1
    }
  ]
}
```

## 分类详情

```
GET /api/categories/:id
```

## 商品列表

```
GET /api/products
```

**查询参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| `category_id` | int | 按分类筛选 |
| `keyword` | string | 搜索关键词（标题模糊匹配） |
| `page` | int | 页码，默认 1 |
| `page_size` | int | 每页数量，默认 20 |

**示例**：

```bash
# 获取 AI 大模型分类下的商品
GET /api/products?category_id=1

# 搜索商品
GET /api/products?keyword=ChatGPT

# 分页
GET /api/products?page=1&page_size=10
```

**响应**：

```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "category_id": 1,
        "title": "ChatGPT Plus 一个月",
        "cover": "https://example.com/cover.jpg",
        "price": 29.90,
        "original_price": 49.90,
        "stock": 100,
        "sales": 256,
        "delivery_type": "auto",
        "tags": "热销,限时",
        "status": 1,
        "is_new": false,
        "is_hot": true
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 10
  }
}
```

## 商品详情

```
GET /api/products/:id
```

**响应**：

```json
{
  "code": 0,
  "data": {
    "id": 1,
    "category_id": 1,
    "title": "ChatGPT Plus 一个月",
    "cover": "https://example.com/cover.jpg",
    "price": 29.90,
    "original_price": 49.90,
    "stock": 100,
    "sales": 256,
    "delivery_type": "auto",
    "tags": "热销,限时",
    "status": 1,
    "is_new": false,
    "is_hot": true,
    "detail": {
      "content": "<p>商品详细描述 HTML</p>"
    }
  }
}
```
