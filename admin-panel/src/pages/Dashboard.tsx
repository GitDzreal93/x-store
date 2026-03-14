import React, { useEffect, useState } from 'react';
import { Card, Col, Row, Statistic, Table, Tag, Spin } from 'antd';
import {
  ShoppingCartOutlined,
  DollarOutlined,
  ShoppingOutlined,
  KeyOutlined,
} from '@ant-design/icons';
import { getProducts } from '../api/product';
import { getCategories } from '../api/category';
import type { Product, Category } from '../types';

const Dashboard: React.FC = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [prodRes, catRes] = await Promise.all([
          getProducts({ page: 1, size: 100 }),
          getCategories(true),
        ]);
        setProducts(prodRes.data.list || []);
        setCategories(catRes.data || []);
      } catch {
        // handled by interceptor
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  const totalProducts = products.length;
  const activeProducts = products.filter((p) => p.status === 1).length;
  const totalStock = products.reduce((sum, p) => sum + p.stock, 0);
  const totalSales = products.reduce((sum, p) => sum + p.sales, 0);
  const lowStockProducts = products.filter((p) => p.stock > 0 && p.stock < 10);
  const outOfStockProducts = products.filter((p) => p.stock === 0);

  const topProducts = [...products].sort((a, b) => b.sales - a.sales).slice(0, 5);

  if (loading) {
    return <div style={{ textAlign: 'center', padding: 100 }}><Spin size="large" /></div>;
  }

  return (
    <div>
      <h2 style={{ marginBottom: 24 }}>数据概览</h2>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="商品总数"
              value={totalProducts}
              prefix={<ShoppingOutlined style={{ color: '#1890ff' }} />}
              suffix={<span style={{ fontSize: 14, color: '#999' }}>/ {activeProducts} 在售</span>}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="总库存"
              value={totalStock}
              prefix={<KeyOutlined style={{ color: '#52c41a' }} />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="累计销量"
              value={totalSales}
              prefix={<ShoppingCartOutlined style={{ color: '#faad14' }} />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="分类数"
              value={categories.length}
              prefix={<DollarOutlined style={{ color: '#722ed1' }} />}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
        <Col xs={24} lg={12}>
          <Card title="热销商品 TOP 5" size="small">
            <Table
              dataSource={topProducts}
              rowKey="id"
              pagination={false}
              size="small"
              columns={[
                { title: '商品', dataIndex: 'title', ellipsis: true },
                { title: '价格', dataIndex: 'price', width: 80, render: (v: number) => `¥${v.toFixed(2)}` },
                { title: '销量', dataIndex: 'sales', width: 70, sorter: (a: Product, b: Product) => a.sales - b.sales },
                {
                  title: '库存',
                  dataIndex: 'stock',
                  width: 70,
                  render: (v: number) => (
                    <span style={{ color: v === 0 ? '#ff4d4f' : v < 10 ? '#faad14' : '#52c41a' }}>{v}</span>
                  ),
                },
              ]}
            />
          </Card>
        </Col>
        <Col xs={24} lg={12}>
          <Card title="库存预警" size="small">
            {outOfStockProducts.length > 0 && (
              <div style={{ marginBottom: 16 }}>
                <h4 style={{ color: '#ff4d4f', marginBottom: 8 }}>缺货商品 ({outOfStockProducts.length})</h4>
                {outOfStockProducts.map((p) => (
                  <Tag color="red" key={p.id} style={{ marginBottom: 4 }}>{p.title}</Tag>
                ))}
              </div>
            )}
            {lowStockProducts.length > 0 && (
              <div>
                <h4 style={{ color: '#faad14', marginBottom: 8 }}>低库存商品 ({lowStockProducts.length})</h4>
                {lowStockProducts.map((p) => (
                  <Tag color="orange" key={p.id} style={{ marginBottom: 4 }}>{p.title} (剩{p.stock})</Tag>
                ))}
              </div>
            )}
            {outOfStockProducts.length === 0 && lowStockProducts.length === 0 && (
              <p style={{ color: '#52c41a' }}>所有商品库存充足 ✓</p>
            )}
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard;
